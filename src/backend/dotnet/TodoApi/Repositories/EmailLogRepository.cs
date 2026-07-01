using TodoApi.Data;
using TodoShared.Models;

namespace TodoApi.Repositories;

public class EmailLogRepository(AppDbContext db) : BaseRepository<EmailLog>(db), IEmailLogRepository
{
}
