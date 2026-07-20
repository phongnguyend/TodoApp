namespace TodoShared.Models;

public class TodoItemAttachment
{
    public int Id { get; set; }
    public int TodoItemId { get; set; }
    public int FileId { get; set; }
    public DateTime CreatedAt { get; set; }
    public int? CreatedByUserId { get; set; }
    public DateTime? UpdatedAt { get; set; }
    public int? UpdatedByUserId { get; set; }

    public TodoItem? TodoItem { get; set; }
    public FileEntity? File { get; set; }
}
