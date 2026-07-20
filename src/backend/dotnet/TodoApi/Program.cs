using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Scalar.AspNetCore;
using TodoApi.Data;
using TodoApi.Repositories;
using TodoApi.Services;

var builder = WebApplication.CreateBuilder(args);

// ── Database (Entity Framework Core) ──────────────────────────────────────────
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlite(builder.Configuration.GetConnectionString("DefaultConnection")
        ?? "Data Source=todo.db"));

// ── Repository & Service (Dependency Injection) ────────────────────────────────
builder.Services.AddScoped<ITodoItemRepository, TodoItemRepository>();
builder.Services.AddScoped<ITodoItemService, TodoItemService>();
builder.Services.AddScoped<ITodoItemAttachmentRepository, TodoItemAttachmentRepository>();
builder.Services.AddScoped<ITodoItemAttachmentService, TodoItemAttachmentService>();
builder.Services.AddScoped<IFileRepository, FileRepository>();
builder.Services.AddScoped<IFileService, FileService>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddScoped<Microsoft.AspNetCore.Identity.IPasswordHasher<TodoShared.Models.User>,
    Microsoft.AspNetCore.Identity.PasswordHasher<TodoShared.Models.User>>();
builder.Services.AddDataProtection();
builder.Services.AddHttpContextAccessor();

var jwtSecret = builder.Configuration["JWT_SECRET_KEY"]
    ?? builder.Configuration["Authentication:Secret"]
    ?? throw new InvalidOperationException("JWT signing secret is not configured.");
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.MapInboundClaims = false;
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = false,
            ValidateAudience = false,
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtSecret)),
            ValidateLifetime = true,
            ClockSkew = TimeSpan.Zero,
            NameClaimType = "sub"
        };
    });
builder.Services.AddAuthorization(options =>
{
    options.FallbackPolicy = new Microsoft.AspNetCore.Authorization.AuthorizationPolicyBuilder()
        .RequireAuthenticatedUser()
        .Build();
});

// ── Controllers & OpenAPI / Swagger ────────────────────────────────────────────
builder.Services.AddControllers();
builder.Services.AddOpenApi();

var app = builder.Build();

// ── Auto-apply EF migrations on startup (development convenience) ──────────────
if (app.Environment.IsDevelopment())
{
    using var scope = app.Services.CreateScope();
    var dbContext = scope.ServiceProvider.GetRequiredService<AppDbContext>();
    dbContext.Database.Migrate();

    app.MapOpenApi();
    // Scalar API reference UI at /scalar/v1
    app.MapScalarApiReference();
}

app.UseHttpsRedirection();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();

app.MapGet("/", () => new { status = "healthy", app = "Todo API", version = "1.0.0" })
   .ExcludeFromDescription();

app.Run();
