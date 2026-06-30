/**
 * worker/main.ts
 * ~~~~~~~~~~~~~~
 * Entry-point for the background worker process.
 *
 * Behaviour
 * ---------
 * * Runs sendIncompleteTodosEmail() immediately on startup.
 * * Then repeats on the interval defined by WORKER_INTERVAL_MINUTES (default 60).
 * * Handles SIGTERM / SIGINT for clean container shutdown.
 *
 * This process is intentionally kept outside the NestJS application context —
 * it accesses Prisma directly and runs as a plain Node.js process in a
 * separate container (see Dockerfile.worker).
 */

import { PrismaClient } from '@prisma/client';
import { sendIncompleteTodosEmail } from './jobs/incomplete-todos-email.job';

const prisma = new PrismaClient();
const INTERVAL_MINUTES = parseInt(process.env.WORKER_INTERVAL_MINUTES ?? '60', 10);

async function main(): Promise<void> {
  console.info(`[worker] Background worker starting (interval=${INTERVAL_MINUTES} min).`);

  await prisma.$connect();

  // Run once immediately so there is no wait on first boot
  await sendIncompleteTodosEmail(prisma);

  // Then repeat on the configured interval
  const timer = setInterval(async () => {
    await sendIncompleteTodosEmail(prisma);
  }, INTERVAL_MINUTES * 60 * 1000);

  const shutdown = async (): Promise<void> => {
    console.info('[worker] Shutdown signal received; stopping worker.');
    clearInterval(timer);
    await prisma.$disconnect();
    process.exit(0);
  };

  process.on('SIGTERM', () => void shutdown());
  process.on('SIGINT', () => void shutdown());
}

main().catch((err: unknown) => {
  console.error('[worker] Fatal error:', err);
  process.exit(1);
});
