using System.Net;
using System.Net.Mail;
using Microsoft.Extensions.Options;

namespace TodoWorker.Services;

public class SmtpSettings
{
    public string Host { get; set; } = "localhost";
    public int Port { get; set; } = 25;
    public string? Username { get; set; }
    public string? Password { get; set; }
    public bool EnableSsl { get; set; } = false;
    public string From { get; set; } = "noreply@todo.app";
}

public class SmtpEmailService(IOptions<SmtpSettings> settings, ILogger<SmtpEmailService> logger) : IEmailService
{
    private readonly SmtpSettings _settings = settings.Value;

    public async Task SendAsync(string to, string subject, string body, CancellationToken ct = default)
    {
        using var client = new SmtpClient(_settings.Host, _settings.Port)
        {
            EnableSsl = _settings.EnableSsl,
            Credentials = string.IsNullOrEmpty(_settings.Username)
                ? null
                : new NetworkCredential(_settings.Username, _settings.Password)
        };

        using var message = new MailMessage(_settings.From, to, subject, body);
        await client.SendMailAsync(message, ct);

        logger.LogInformation("Email sent to {Recipient} with subject '{Subject}'", to, subject);
    }
}
