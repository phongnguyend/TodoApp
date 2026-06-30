"""
incomplete_todos_email.py
~~~~~~~~~~~~~~~~~~~~~~~~~
Job: query all incomplete todo items, build an email digest, persist an
EmailLog record, then send it via SMTP.

Flow
----
1. Query incomplete todos from the database.
2. Build plain-text and HTML body.
3. INSERT an EmailLog row with status="pending".
4. Send the email via smtplib (STARTTLS or SSL depending on SMTP_USE_TLS).
5. UPDATE the EmailLog row to status="sent" (or "failed" + error_message).
"""

import logging
import smtplib
from datetime import datetime, timezone
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

from sqlalchemy.orm import Session

from app.config import settings
from app.database import SessionLocal
from app.models.email_log import EmailLog
from app.models.todo_item import TodoItem

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------


def _build_email(todos: list[TodoItem]) -> tuple[str, str, str]:
    """Return (subject, plain_text_body, html_body)."""
    count = len(todos)
    subject = f"Incomplete Todos Digest — {count} item(s) pending"

    # ── plain text ──────────────────────────────────────────────────────────
    lines: list[str] = [f"You have {count} incomplete todo item(s):\n"]
    for i, todo in enumerate(todos, start=1):
        lines.append(f"{i}. [{todo.id}] {todo.title}")
        if todo.description:
            lines.append(f"   {todo.description}")
        lines.append(f"   Created: {todo.created_at.strftime('%Y-%m-%d %H:%M UTC')}")
        lines.append("")
    text_body = "\n".join(lines)

    # ── HTML ────────────────────────────────────────────────────────────────
    rows = "".join(
        f"<tr>"
        f"<td>{todo.id}</td>"
        f"<td>{todo.title}</td>"
        f"<td>{todo.description or ''}</td>"
        f"<td>{todo.created_at.strftime('%Y-%m-%d %H:%M UTC')}</td>"
        f"</tr>"
        for todo in todos
    )
    html_body = (
        "<html><body>"
        f"<p>You have <strong>{count}</strong> incomplete todo item(s):</p>"
        "<table border='1' cellpadding='4' cellspacing='0'>"
        "<thead><tr><th>ID</th><th>Title</th><th>Description</th><th>Created At</th></tr></thead>"
        f"<tbody>{rows}</tbody>"
        "</table></body></html>"
    )
    return subject, text_body, html_body


def _send_smtp(subject: str, text_body: str, html_body: str) -> None:
    """Send a multipart/alternative email via SMTP."""
    msg = MIMEMultipart("alternative")
    msg["Subject"] = subject
    msg["From"] = settings.EMAIL_SENDER
    msg["To"] = settings.EMAIL_RECIPIENT
    msg.attach(MIMEText(text_body, "plain", "utf-8"))
    msg.attach(MIMEText(html_body, "html", "utf-8"))

    if settings.SMTP_USE_TLS:
        conn = smtplib.SMTP(settings.SMTP_HOST, settings.SMTP_PORT, timeout=30)
        conn.ehlo()
        conn.starttls()
        conn.ehlo()
    else:
        conn = smtplib.SMTP_SSL(settings.SMTP_HOST, settings.SMTP_PORT, timeout=30)

    with conn:
        if settings.SMTP_USERNAME and settings.SMTP_PASSWORD:
            conn.login(settings.SMTP_USERNAME, settings.SMTP_PASSWORD)
        conn.sendmail(settings.EMAIL_SENDER, [settings.EMAIL_RECIPIENT], msg.as_string())


# ---------------------------------------------------------------------------
# Public job entry-point
# ---------------------------------------------------------------------------


def send_incomplete_todos_email() -> None:
    """
    Scheduled job called by the worker scheduler.

    Fetches every incomplete todo, creates a digest email, persists an
    EmailLog record, and delivers the email via SMTP.
    """
    db: Session = SessionLocal()
    log: EmailLog | None = None
    try:
        todos: list[TodoItem] = (
            db.query(TodoItem).filter(TodoItem.is_completed.is_(False)).all()
        )

        if not todos:
            logger.info("No incomplete todos — skipping email digest.")
            return

        logger.info("Found %d incomplete todo(s); preparing email digest.", len(todos))

        subject, text_body, html_body = _build_email(todos)

        # Persist as "pending" before attempting delivery
        log = EmailLog(
            recipient=settings.EMAIL_RECIPIENT,
            subject=subject,
            body=text_body,
            status="pending",
        )
        db.add(log)
        db.commit()
        db.refresh(log)

        try:
            _send_smtp(subject, text_body, html_body)
            log.status = "sent"
            log.sent_at = datetime.now(timezone.utc)
            db.commit()
            logger.info("Email digest sent (email_log.id=%d).", log.id)
        except Exception as smtp_exc:
            logger.exception("SMTP delivery failed (email_log.id=%d).", log.id)
            log.status = "failed"
            log.error_message = str(smtp_exc)
            db.commit()

    except Exception:
        logger.exception("Unexpected error in send_incomplete_todos_email job.")
    finally:
        db.close()
