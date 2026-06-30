package com.example.todo.worker;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.stereotype.Component;

/**
 * WorkerRunner — schedules the email-digest job on startup.
 *
 * Behaviour (mirrors the Node.js / Python worker entry-points):
 * - Runs {@link IncompleteTodosEmailJob#execute()} immediately on startup.
 * - Repeats every {@code WORKER_INTERVAL_MINUTES} (default 60) minutes.
 * - Stops cleanly on SIGTERM / SIGINT via Spring's JVM shutdown hook.
 *
 * Because {@code spring.main.web-application-type=none} is set in
 * {@code application.yml}, the embedded web server is not started
 * and this runner's loop becomes the sole long-running task.
 */
@Component
@RequiredArgsConstructor
@Slf4j
public class WorkerRunner implements ApplicationRunner {

    private final IncompleteTodosEmailJob job;

    @Value("${worker.interval-minutes:60}")
    private int intervalMinutes;

    @Override
    public void run(ApplicationArguments args) throws InterruptedException {
        log.info("[worker] Background worker starting (interval={} min).", intervalMinutes);

        // Run once immediately so there is no wait on first boot
        job.execute();

        long intervalMs = (long) intervalMinutes * 60 * 1_000;

        while (!Thread.currentThread().isInterrupted()) {
            try {
                Thread.sleep(intervalMs);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
            job.execute();
        }

        log.info("[worker] Worker stopped.");
    }
}
