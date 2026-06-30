/**
 * incomplete-todos-email.job.ts
 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * Job: query all incomplete todo items, build an email digest, persist an
 * EmailLog record, then send it via SMTP.
 *
 * Flow
 * ----
 * 1. Query incomplete todos from the database.
 * 2. Build plain-text and HTML body.
 * 3. INSERT an EmailLog row with status="pending".
 * 4. Send the email via nodemailer (STARTTLS or SSL depending on SMTP_SECURE).
 * 5. UPDATE the EmailLog row to status="sent" (or "failed" + errorMessage).
 */

import * as nodemailer from 'nodemailer';
import type { PrismaClient, TodoItem } from '@prisma/client';

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function buildEmail(todos: TodoItem[]): {
  subject: string;
  textBody: string;
  htmlBody: string;
} {
  const count = todos.length;
  const subject = `Incomplete Todos Digest — ${count} item(s) pending`;

  // ── plain text ─────────────────────────────────────────────────────────────
  const lines: string[] = [`You have ${count} incomplete todo item(s):\n`];
  todos.forEach((todo, i) => {
    lines.push(`${i + 1}. [${todo.id}] ${todo.title}`);
    if (todo.description) lines.push(`   ${todo.description}`);
    lines.push(`   Created: ${todo.createdAt.toISOString().replace('T', ' ').substring(0, 16)} UTC`);
    lines.push('');
  });
  const textBody = lines.join('\n');

  // ── HTML ────────────────────────────────────────────────────────────────────
  const rows = todos
    .map(
      (todo) =>
        `<tr>` +
        `<td>${todo.id}</td>` +
        `<td>${todo.title}</td>` +
        `<td>${todo.description ?? ''}</td>` +
        `<td>${todo.createdAt.toISOString().replace('T', ' ').substring(0, 16)} UTC</td>` +
        `</tr>`,
    )
    .join('');

  const htmlBody =
    `<html><body>` +
    `<p>You have <strong>${count}</strong> incomplete todo item(s):</p>` +
    `<table border="1" cellpadding="4" cellspacing="0">` +
    `<thead><tr><th>ID</th><th>Title</th><th>Description</th><th>Created At</th></tr></thead>` +
    `<tbody>${rows}</tbody>` +
    `</table></body></html>`;

  return { subject, textBody, htmlBody };
}

// ---------------------------------------------------------------------------
// Public job entry-point
// ---------------------------------------------------------------------------

export async function sendIncompleteTodosEmail(prisma: PrismaClient): Promise<void> {
  try {
    const todos = await prisma.todoItem.findMany({
      where: { isCompleted: false },
      orderBy: { createdAt: 'asc' },
    });

    if (todos.length === 0) {
      console.info('[worker] No incomplete todos — skipping email digest.');
      return;
    }

    console.info(`[worker] Found ${todos.length} incomplete todo(s); preparing email digest.`);

    const { subject, textBody, htmlBody } = buildEmail(todos);

    const recipient = process.env.EMAIL_RECIPIENT ?? 'admin@example.com';
    const sender = process.env.EMAIL_SENDER ?? 'noreply@example.com';

    // Persist as "pending" before attempting delivery
    const emailLog = await prisma.emailLog.create({
      data: {
        recipient,
        subject,
        body: textBody,
        status: 'pending',
      },
    });

    try {
      const transporter = nodemailer.createTransport({
        host: process.env.SMTP_HOST ?? 'localhost',
        port: parseInt(process.env.SMTP_PORT ?? '587', 10),
        secure: process.env.SMTP_SECURE === 'true',
        auth:
          process.env.SMTP_USERNAME
            ? { user: process.env.SMTP_USERNAME, pass: process.env.SMTP_PASSWORD ?? '' }
            : undefined,
      });

      await transporter.sendMail({ from: sender, to: recipient, subject, text: textBody, html: htmlBody });

      await prisma.emailLog.update({
        where: { id: emailLog.id },
        data: { status: 'sent', sentAt: new Date() },
      });

      console.info(`[worker] Email digest sent (emailLog.id=${emailLog.id}).`);
    } catch (smtpError) {
      console.error(`[worker] SMTP delivery failed (emailLog.id=${emailLog.id}):`, smtpError);
      await prisma.emailLog.update({
        where: { id: emailLog.id },
        data: {
          status: 'failed',
          errorMessage: smtpError instanceof Error ? smtpError.message : String(smtpError),
        },
      });
    }
  } catch (err) {
    console.error('[worker] Unexpected error in sendIncompleteTodosEmail:', err);
  }
}
