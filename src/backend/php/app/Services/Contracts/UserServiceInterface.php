<?php

namespace App\Services\Contracts;

use App\Http\Requests\ChangePasswordRequest;
use App\Http\Requests\ConfirmPasswordResetRequest;
use App\Http\Requests\CreateTokenRequest;
use App\Http\Requests\CreateUserRequest;
use App\Http\Requests\ResetPasswordRequest;
use App\Http\Requests\SignUpRequest;
use App\Http\Requests\UpdateProfileRequest;
use App\Http\Requests\UpdateUserRequest;
use App\Models\User;
use Illuminate\Pagination\LengthAwarePaginator;

interface UserServiceInterface
{
    public function getAll(int $page = 1, int $perPage = 20): LengthAwarePaginator;

    public function getById(int $id): User;

    public function create(CreateUserRequest $request): User;

    public function update(int $id, UpdateUserRequest $request): User;

    public function setActive(int $id, bool $isActive): User;

    public function signup(SignUpRequest $request): User;

    public function getProfile(int $userId): User;

    public function updateProfile(int $userId, UpdateProfileRequest $request): User;

    public function changePassword(int $userId, ChangePasswordRequest $request): void;

    public function requestPasswordReset(ResetPasswordRequest $request): void;

    public function confirmPasswordReset(ConfirmPasswordResetRequest $request): void;

    public function createToken(CreateTokenRequest $request): ?array;
}
