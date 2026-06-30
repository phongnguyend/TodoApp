using Microsoft.EntityFrameworkCore;
using TodoWorker.Data;
using TodoWorker.Services;

var builder = Host.CreateApplicationBuilder(args);

// ── Database (Entity Framework Core) ──────────────────────────────────────────
builder.Services.AddDbContext<WorkerDbContext>(options =>
    options.UseSqlite(builder.Configuration.GetConnectionString("DefaultConnection")
        ?? "Data Source=todo.db"));

// ── Settings ───────────────────────────────────────────────────────────────────
builder.Services.Configure<SmtpSettings>(builder.Configuration.GetSection("Smtp"));
builder.Services.Configure<WorkerSettings>(builder.Configuration.GetSection("Worker"));

// ── Email & Background Service ─────────────────────────────────────────────────
builder.Services.AddScoped<IEmailService, SmtpEmailService>();
builder.Services.AddHostedService<WorkerService>();

var host = builder.Build();
await host.RunAsync();
