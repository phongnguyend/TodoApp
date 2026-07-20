using Microsoft.EntityFrameworkCore;
using TodoShared.Models;

namespace TodoShared.Data;

public abstract class TodoDbContext(DbContextOptions options) : DbContext(options)
{
    public DbSet<TodoItem> TodoItems => Set<TodoItem>();
    public DbSet<TodoItemAttachment> TodoItemAttachments => Set<TodoItemAttachment>();
    public DbSet<EmailLog> EmailLogs => Set<EmailLog>();
    public DbSet<FileEntity> Files => Set<FileEntity>();
    public DbSet<User> Users => Set<User>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<TodoItem>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Title).IsRequired().HasMaxLength(200);
            entity.Property(e => e.Description).HasMaxLength(2000);
            entity.Property(e => e.IsCompleted).HasDefaultValue(false);
            entity.Property(e => e.CreatedAt).IsRequired();
        });

        modelBuilder.Entity<EmailLog>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Recipient).IsRequired().HasMaxLength(255);
            entity.Property(e => e.Subject).IsRequired().HasMaxLength(500);
            entity.Property(e => e.Body).IsRequired();
            entity.Property(e => e.Status).IsRequired().HasMaxLength(50).HasDefaultValue("pending");
            entity.Property(e => e.CreatedAt).IsRequired();
        });

        modelBuilder.Entity<FileEntity>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Name).IsRequired().HasMaxLength(255);
            entity.Property(e => e.Extension).IsRequired().HasMaxLength(20);
            entity.Property(e => e.Size).IsRequired();
            entity.Property(e => e.ContentType).HasMaxLength(100);
            entity.Property(e => e.Location).IsRequired().HasMaxLength(500);
            entity.Property(e => e.CreatedAt).IsRequired();
        });

        modelBuilder.Entity<TodoItemAttachment>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.CreatedAt).IsRequired();
            entity.HasIndex(e => new { e.TodoItemId, e.FileId }).IsUnique();
            entity.HasOne(e => e.TodoItem)
                .WithMany()
                .HasForeignKey(e => e.TodoItemId)
                .OnDelete(DeleteBehavior.Cascade);
            entity.HasOne(e => e.File)
                .WithMany()
                .HasForeignKey(e => e.FileId)
                .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Username).IsRequired().HasMaxLength(50);
            entity.Property(e => e.Email).IsRequired().HasMaxLength(255);
            entity.Property(e => e.PasswordHash).IsRequired().HasMaxLength(255);
            entity.Property(e => e.IsActive).HasDefaultValue(true);
            entity.Property(e => e.CreatedAt).IsRequired();
            entity.HasIndex(e => e.Username).IsUnique();
            entity.HasIndex(e => e.Email).IsUnique();
        });

        foreach (var entityType in new[]
        {
            typeof(TodoItem), typeof(TodoItemAttachment), typeof(EmailLog), typeof(FileEntity), typeof(User)
        })
        {
            modelBuilder.Entity(entityType)
                .HasOne(typeof(User), null)
                .WithMany()
                .HasForeignKey(nameof(TodoItem.CreatedByUserId))
                .OnDelete(DeleteBehavior.SetNull);
            modelBuilder.Entity(entityType)
                .HasOne(typeof(User), null)
                .WithMany()
                .HasForeignKey(nameof(TodoItem.UpdatedByUserId))
                .OnDelete(DeleteBehavior.SetNull);
        }
    }
}
