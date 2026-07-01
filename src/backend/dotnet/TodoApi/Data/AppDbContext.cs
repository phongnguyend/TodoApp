using Microsoft.EntityFrameworkCore;
using TodoShared.Data;

namespace TodoApi.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : TodoDbContext(options)
{
}
