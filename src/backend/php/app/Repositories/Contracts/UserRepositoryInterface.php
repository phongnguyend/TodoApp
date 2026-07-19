<?php

namespace App\Repositories\Contracts;

use App\Models\EmailLog;
use App\Models\User;

interface UserRepositoryInterface extends RepositoryInterface
{
    public function findByEmail(string $email): ?User;

    public function usernameExists(string $username, ?int $excludingId = null): bool;

    public function emailExists(string $email, ?int $excludingId = null): bool;

    public function createEmailLog(array $data): EmailLog;
}
