using TodoShared.Models;

namespace TodoApi.Repositories;

public interface IUserRepository : IRepository<User>
{
    Task<User?> GetByEmailAsync(string email, CancellationToken ct = default);
    Task<bool> UsernameExistsAsync(string username, int? excludingId = null, CancellationToken ct = default);
    Task<bool> EmailExistsAsync(string email, int? excludingId = null, CancellationToken ct = default);
}
