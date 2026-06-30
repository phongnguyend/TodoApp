"""
main.py
~~~~~~~
Entry-point for the background worker process.

Behaviour
---------
* Runs send_incomplete_todos_email() immediately on startup.
* Then repeats on the interval defined by WORKER_INTERVAL_MINUTES (default 60).
* Handles SIGTERM / SIGINT for clean container shutdown.
"""

import logging
import signal
import sys

from apscheduler.schedulers.blocking import BlockingScheduler

from shared.config import settings
from worker.jobs.incomplete_todos_email import send_incomplete_todos_email

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger(__name__)


def main() -> None:
    logger.info(
        "Background worker starting (interval=%d min).",
        settings.WORKER_INTERVAL_MINUTES,
    )

    scheduler = BlockingScheduler(timezone="UTC")
    scheduler.add_job(
        send_incomplete_todos_email,
        trigger="interval",
        minutes=settings.WORKER_INTERVAL_MINUTES,
        id="incomplete_todos_email",
        replace_existing=True,
        max_instances=1,
    )

    def _shutdown(signum, frame) -> None:  # noqa: ANN001
        logger.info("Shutdown signal received; stopping scheduler.")
        scheduler.shutdown(wait=False)
        sys.exit(0)

    signal.signal(signal.SIGTERM, _shutdown)
    signal.signal(signal.SIGINT, _shutdown)

    # Run once immediately so there is no wait on first boot
    send_incomplete_todos_email()

    scheduler.start()


if __name__ == "__main__":
    main()
