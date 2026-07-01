package com.example.todo.worker;

import com.example.todo.entity.EmailLog;
import com.example.todo.entity.TodoItem;
import com.example.todo.repository.EmailLogRepository;
import com.example.todo.repository.TodoItemRepository;
import jakarta.mail.internet.MimeMessage;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.List;

/**
 * IncompleteTodosEmailJob — the core worker logic.
 *
 * Flow (mirrors the Node.js / Python implementations):
 * 1. Query all incomplete todo items ordered by created_at.
 * 2. Build plain-text and HTML email bodies.
 * 3. INSERT an EmailLog row with status="pending".
 * 4. Send the email via Spring JavaMailSender (SMTP).
 * 5. UPDATE the EmailLog row to status="sent" (or "failed" + errorMessage).
 */
@Component
@RequiredArgsConstructor
@Slf4j
public class IncompleteTodosEmailJob {

    private static final DateTimeFormatter FORMATTER =
            DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm").withZone(ZoneOffset.UTC);

    private final TodoItemRepository todoItemRepository;
    private final EmailLogRepository emailLogRepository;
    private final JavaMailSender mailSender;

    @Value("${worker.email.recipient:admin@example.com}")
    private String recipient;

    @Value("${worker.email.sender:noreply@example.com}")
    private String sender;

    /**
     * Runs one email-digest cycle. The outer {@code @Transactional} covers both the
     * "pending" insert and the final "sent"/"failed" update in a single JPA session.
     * SMTP delivery happens inside the transaction boundary; any mail exception is
     * caught so the transaction always commits with the final status set.
     */
    @Transactional
    public void execute() {
        log.info("[worker] Running incomplete todos email job.");

        List<TodoItem> todos = todoItemRepository.findByCompletedFalseOrderByCreatedAtAsc();

        if (todos.isEmpty()) {
            log.info("[worker] No incomplete todos — skipping email digest.");
            return;
        }

        log.info("[worker] Found {} incomplete todo(s); preparing email digest.", todos.size());

        String subject = String.format("Incomplete Todos Digest — %d item(s) pending", todos.size());
        String textBody = buildTextBody(todos);
        String htmlBody = buildHtmlBody(todos);

        // Persist as "pending" before attempting delivery
        EmailLog emailLog = new EmailLog(recipient, subject, textBody);
        emailLogRepository.save(emailLog);

        try {
            MimeMessage message = mailSender.createMimeMessage();
            MimeMessageHelper helper = new MimeMessageHelper(message, true, "UTF-8");
            helper.setFrom(sender);
            helper.setTo(recipient);
            helper.setSubject(subject);
            // setText(plain, html) creates a multipart/alternative message
            helper.setText(textBody, htmlBody);

            mailSender.send(message);

            emailLog.setStatus("sent");
            emailLog.setSentAt(Instant.now());
            emailLogRepository.save(emailLog);

            log.info("[worker] Email digest sent (emailLog.id={}).", emailLog.getId());
        } catch (Exception e) {
            log.error("[worker] SMTP delivery failed (emailLog.id={}):", emailLog.getId(), e);
            emailLog.setStatus("failed");
            emailLog.setErrorMessage(e.getMessage());
            emailLogRepository.save(emailLog);
        }
    }

    // ── Email body builders ────────────────────────────────────────────────────

    private String buildTextBody(List<TodoItem> todos) {
        StringBuilder sb = new StringBuilder();
        sb.append(String.format("You have %d incomplete todo item(s):%n%n", todos.size()));
        for (int i = 0; i < todos.size(); i++) {
            TodoItem todo = todos.get(i);
            sb.append(String.format("%d. [%d] %s%n", i + 1, todo.getId(), todo.getTitle()));
            if (todo.getDescription() != null) {
                sb.append(String.format("   %s%n", todo.getDescription()));
            }
            sb.append(String.format("   Created: %s UTC%n%n", FORMATTER.format(todo.getCreatedAt())));
        }
        return sb.toString();
    }

    private String buildHtmlBody(List<TodoItem> todos) {
        StringBuilder rows = new StringBuilder();
        for (TodoItem todo : todos) {
            rows.append("<tr>")
                    .append("<td>").append(todo.getId()).append("</td>")
                    .append("<td>").append(escapeHtml(todo.getTitle())).append("</td>")
                    .append("<td>").append(todo.getDescription() != null ? escapeHtml(todo.getDescription()) : "").append("</td>")
                    .append("<td>").append(FORMATTER.format(todo.getCreatedAt())).append(" UTC</td>")
                    .append("</tr>");
        }
        return "<html><body>" +
                "<p>You have <strong>" + todos.size() + "</strong> incomplete todo item(s):</p>" +
                "<table border=\"1\" cellpadding=\"4\" cellspacing=\"0\">" +
                "<thead><tr><th>ID</th><th>Title</th><th>Description</th><th>Created At</th></tr></thead>" +
                "<tbody>" + rows + "</tbody>" +
                "</table></body></html>";
    }

    private static String escapeHtml(String text) {
        return text.replace("&", "&amp;")
                .replace("<", "&lt;")
                .replace(">", "&gt;")
                .replace("\"", "&quot;");
    }
}
