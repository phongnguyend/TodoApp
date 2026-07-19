using System.Text;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Options;
using TodoShared.Models;
using TodoWorker.Data;

namespace TodoWorker.Services;

public class WorkerSettings
{
    public int IntervalMinutes { get; set; } = 5;
    public string RecipientEmail { get; set; } = "admin@todo.app";
}

public class WorkerService(
    IServiceScopeFactory scopeFactory,
    IOptions<WorkerSettings> settings,
    ILogger<WorkerService> logger) : BackgroundService
{
    private readonly WorkerSettings _settings = settings.Value;

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        logger.LogInformation("Worker started. Interval: {Interval} minute(s)", _settings.IntervalMinutes);

        // Run once immediately on startup
        await ProcessAsync(stoppingToken);

        using var timer = new PeriodicTimer(TimeSpan.FromMinutes(_settings.IntervalMinutes));

        while (await timer.WaitForNextTickAsync(stoppingToken))
        {
            await ProcessAsync(stoppingToken);
        }
    }

    private async Task ProcessAsync(CancellationToken ct)
    {
        logger.LogInformation("Checking for incomplete todos...");

        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<WorkerDbContext>();
        var emailService = scope.ServiceProvider.GetRequiredService<IEmailService>();

        await SendPendingEmailsAsync(db, emailService, ct);

        List<TodoItem> incompleteTodos;
        try
        {
            incompleteTodos = await db.TodoItems
                .Where(t => !t.IsCompleted)
                .OrderBy(t => t.CreatedAt)
                .ToListAsync(ct);
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "Failed to query incomplete todos");
            return;
        }

        if (incompleteTodos.Count == 0)
        {
            logger.LogInformation("No incomplete todos found - skipping email");
            return;
        }

        logger.LogInformation("Found {Count} incomplete todo(s) - preparing email", incompleteTodos.Count);

        var subject = $"Incomplete Todos - {incompleteTodos.Count} item(s) pending";
        var body = BuildEmailBody(incompleteTodos);

        var emailLog = new EmailLog
        {
            Recipient = _settings.RecipientEmail,
            Subject = subject,
            Body = body,
            Status = "pending",
            CreatedAt = DateTime.UtcNow
        };

        db.EmailLogs.Add(emailLog);
        await db.SaveChangesAsync(ct);

        try
        {
            await emailService.SendAsync(_settings.RecipientEmail, subject, body, ct);
            emailLog.Status = "sent";
            emailLog.SentAt = DateTime.UtcNow;
            logger.LogInformation("Email sent successfully to {Recipient}", _settings.RecipientEmail);
        }
        catch (Exception ex)
        {
            emailLog.Status = "failed";
            emailLog.ErrorMessage = ex.Message;
            logger.LogError(ex, "Failed to send email to {Recipient}", _settings.RecipientEmail);
        }

        await db.SaveChangesAsync(ct);
    }

    private async Task SendPendingEmailsAsync(WorkerDbContext db, IEmailService emailService, CancellationToken ct)
    {
        var pendingEmails = await db.EmailLogs
            .Where(email => email.Status == "pending")
            .OrderBy(email => email.CreatedAt)
            .ToListAsync(ct);

        foreach (var email in pendingEmails)
        {
            try
            {
                await emailService.SendAsync(email.Recipient, email.Subject, email.Body, ct);
                email.Status = "sent";
                email.SentAt = DateTime.UtcNow;
                email.ErrorMessage = null;
            }
            catch (Exception ex)
            {
                email.Status = "failed";
                email.ErrorMessage = ex.Message;
                logger.LogError(ex, "Failed to send queued email {EmailLogId} to {Recipient}", email.Id, email.Recipient);
            }
        }

        if (pendingEmails.Count > 0)
            await db.SaveChangesAsync(ct);
    }

    private static string BuildEmailBody(IEnumerable<TodoItem> items)
    {
        var sb = new StringBuilder();
        sb.AppendLine("The following todo items are still incomplete:");
        sb.AppendLine();

        foreach (var item in items)
        {
            sb.AppendLine($"• [{item.Id}] {item.Title}");
            if (!string.IsNullOrWhiteSpace(item.Description))
                sb.AppendLine($"  {item.Description}");
        }

        sb.AppendLine();
        sb.AppendLine($"Generated at: {DateTime.UtcNow:yyyy-MM-dd HH:mm:ss} UTC");
        return sb.ToString();
    }
}
