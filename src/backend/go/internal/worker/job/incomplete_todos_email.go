package job

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/gorm"
)

// SendIncompleteTodosEmail queries all incomplete todos, builds a digest email,
// persists an EmailLog record, then delivers the message via SMTP.
//
// Flow mirrors the Python worker job:
//  1. Query incomplete todos.
//  2. Build plain-text email body.
//  3. INSERT an EmailLog row with status="pending".
//  4. Send via SMTP.
//  5. UPDATE EmailLog to status="sent" or "failed".
func SendIncompleteTodosEmail(db *gorm.DB, cfg *config.Config) {
	todoRepo := repository.NewTodoItemRepository(db)
	logRepo := repository.NewEmailLogRepository(db)

	// 1. Fetch all incomplete todos (no pagination - digest covers everything).
	result, err := todoRepo.FindIncomplete(0, 10_000)
	if err != nil {
		log.Printf("[worker] failed to query incomplete todos: %v", err)
		return
	}
	if len(result.Items) == 0 {
		log.Println("[worker] no incomplete todos - skipping email digest")
		return
	}

	log.Printf("[worker] found %d incomplete todo(s); preparing email digest", len(result.Items))

	// 2. Build email content.
	subject, body := buildEmail(result.Items)

	// 3. Persist EmailLog with status="pending".
	emailLog := &models.EmailLog{
		Recipient: cfg.EmailRecipient,
		Subject:   subject,
		Body:      body,
		Status:    "pending",
	}
	emailLog, err = logRepo.Create(emailLog)
	if err != nil {
		log.Printf("[worker] failed to create email log: %v", err)
		return
	}

	// 4. Send via SMTP.
	sendErr := sendSMTP(cfg, subject, body)

	// 5. Update EmailLog status.
	if sendErr != nil {
		log.Printf("[worker] SMTP delivery failed: %v", sendErr)
		if markErr := logRepo.MarkFailed(emailLog, sendErr.Error()); markErr != nil {
			log.Printf("[worker] failed to mark email log as failed: %v", markErr)
		}
		return
	}

	if markErr := logRepo.MarkSent(emailLog); markErr != nil {
		log.Printf("[worker] failed to mark email log as sent: %v", markErr)
	}
	log.Printf("[worker] digest email sent to %s (log id=%d)", cfg.EmailRecipient, emailLog.ID)
}

// buildEmail constructs the email subject and plain-text body.
func buildEmail(todos []models.TodoItem) (subject, body string) {
	count := len(todos)
	subject = fmt.Sprintf("Incomplete Todos Digest - %d item(s) pending", count)

	var sb strings.Builder
	fmt.Fprintf(&sb, "You have %d incomplete todo item(s):\n\n", count)
	for i, todo := range todos {
		fmt.Fprintf(&sb, "%d. [%d] %s\n", i+1, todo.ID, todo.Title)
		if todo.Description != nil && *todo.Description != "" {
			fmt.Fprintf(&sb, "   %s\n", *todo.Description)
		}
		fmt.Fprintf(&sb, "   Created: %s\n\n", todo.CreatedAt.UTC().Format(time.RFC3339))
	}

	return subject, sb.String()
}

// sendSMTP delivers a plain-text email using the configured SMTP settings.
// Uses STARTTLS when cfg.SMTPUseTLS is true, otherwise plain SMTP.
func sendSMTP(cfg *config.Config, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	msg := buildMessage(cfg.EmailSender, cfg.EmailRecipient, subject, body)

	var auth smtp.Auth
	if cfg.SMTPUsername != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	}

	if cfg.SMTPUseTLS {
		return sendSTARTTLS(addr, auth, cfg.EmailSender, cfg.EmailRecipient, msg)
	}
	return smtp.SendMail(addr, auth, cfg.EmailSender, []string{cfg.EmailRecipient}, []byte(msg))
}

// buildMessage formats an RFC 2822 email message.
func buildMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)
}

// sendSTARTTLS connects, issues EHLO, upgrades to TLS, then sends the message.
func sendSTARTTLS(addr string, auth smtp.Auth, from, to, msg string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid smtp address %q: %w", addr, err)
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer c.Close() //nolint:errcheck

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	if _, err = fmt.Fprint(w, msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return c.Quit()
}
