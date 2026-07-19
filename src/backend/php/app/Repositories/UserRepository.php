<?php

namespace App\Repositories;

use App\Models\EmailLog;
use App\Models\User;
use App\Repositories\Contracts\UserRepositoryInterface;
use Illuminate\Database\Eloquent\Builder;

class UserRepository extends BaseRepository implements UserRepositoryInterface
{
    public function __construct(User $model)
    {
        parent::__construct($model);
    }

    public function findByEmail(string $email): ?User
    {
        return User::whereRaw('LOWER(email) = ?', [strtolower($email)])->first();
    }

    public function usernameExists(string $username, ?int $excludingId = null): bool
    {
        return $this->existsIgnoringCase('username', $username, $excludingId);
    }

    public function emailExists(string $email, ?int $excludingId = null): bool
    {
        return $this->existsIgnoringCase('email', $email, $excludingId);
    }

    public function createEmailLog(array $data): EmailLog
    {
        return EmailLog::create($data);
    }

    private function existsIgnoringCase(string $column, string $value, ?int $excludingId): bool
    {
        return User::whereRaw("LOWER({$column}) = ?", [strtolower($value)])
            ->when($excludingId !== null, fn (Builder $query) => $query->whereKeyNot($excludingId))
            ->exists();
    }
}
