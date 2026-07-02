using TodoApi.Data;
using TodoShared.Models;

namespace TodoApi.Repositories;

public class FileRepository(AppDbContext db) : BaseRepository<FileEntity>(db), IFileRepository
{
}
