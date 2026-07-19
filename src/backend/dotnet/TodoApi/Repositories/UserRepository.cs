using Microsoft.EntityFrameworkCore;
using TodoApi.Data;
using TodoShared.Models;

namespace TodoApi.Repositories;

public class UserRepository(AppDbContext db) : BaseRepository<User>(db), IUserRepository
{
    public Task<User?> GetByEmailAsync(string email, CancellationToken ct = default) =>
        Db.Users.FirstOrDefaultAsync(user => user.Email == email.ToLower(), ct);

    public Task<bool> UsernameExistsAsync(string username, int? excludingId = null, CancellationToken ct = default) =>
        Db.Users.AnyAsync(user => user.Username.ToLower() == username.ToLower() && (!excludingId.HasValue || user.Id != excludingId), ct);

    public Task<bool> EmailExistsAsync(string email, int? excludingId = null, CancellationToken ct = default) =>
        Db.Users.AnyAsync(user => user.Email == email.ToLower() && (!excludingId.HasValue || user.Id != excludingId), ct);
}
