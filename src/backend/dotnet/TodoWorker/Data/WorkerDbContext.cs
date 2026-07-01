using Microsoft.EntityFrameworkCore;
using TodoShared.Data;

namespace TodoWorker.Data;

public class WorkerDbContext(DbContextOptions<WorkerDbContext> options) : TodoDbContext(options)
{
}
