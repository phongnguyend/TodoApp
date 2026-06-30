using TodoApi.Data;
using TodoApi.Models;

namespace TodoApi.Repositories;

public class EmailLogRepository(AppDbContext db) : BaseRepository<EmailLog>(db), IEmailLogRepository
{
}
