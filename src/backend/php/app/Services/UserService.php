<?php

namespace App\Services;

use App\Exceptions\InvalidPasswordException;
use App\Exceptions\InvalidPasswordResetTokenException;
use App\Exceptions\UserConflictException;
use App\Http\Requests\ChangePasswordRequest;
use App\Http\Requests\ConfirmPasswordResetRequest;
use App\Http\Requests\CreateTokenRequest;
use App\Http\Requests\CreateUserRequest;
use App\Http\Requests\ResetPasswordRequest;
use App\Http\Requests\SignUpRequest;
use App\Http\Requests\UpdateProfileRequest;
use App\Http\Requests\UpdateUserRequest;
use App\Models\User;
use App\Repositories\Contracts\UserRepositoryInterface;
use App\Services\Contracts\UserServiceInterface;
use Firebase\JWT\JWT;
use Illuminate\Contracts\Encryption\DecryptException;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Pagination\LengthAwarePaginator;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Hash;
use JsonException;

class UserService implements UserServiceInterface
{
    private ?string $dummyPasswordHash = null;

    public function __construct(
        private readonly UserRepositoryInterface $repository
    ) {}

    public function getAll(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginate(max(1, $page), min(100, max(1, $perPage)));
    }

    public function getById(int $id): User
    {
        $user = $this->repository->findById($id);
        if (! $user instanceof User) {
            throw new ModelNotFoundException("User {$id} not found.");
        }

        return $user;
    }

    public function create(CreateUserRequest $request): User
    {
        $data = $request->validated();
        $username = trim($data['username']);
        $email = $this->normalizeEmail($data['email']);
        $this->ensureUnique($username, $email);

        $user = $this->repository->create([
            'username' => $username,
            'email' => $email,
            'password_hash' => Hash::make($data['password']),
            'is_active' => $data['is_active'] ?? true,
        ]);

        /** @var User $user */
        return $user;
    }

    public function update(int $id, UpdateUserRequest $request): User
    {
        $user = $this->getById($id);
        $data = $request->validated();
        $username = array_key_exists('username', $data) ? trim($data['username']) : $user->username;
        $email = array_key_exists('email', $data) ? $this->normalizeEmail($data['email']) : $user->email;
        $this->ensureUnique($username, $email, $id);

        $changes = ['username' => $username, 'email' => $email, 'updated_at' => now()];
        if (array_key_exists('password', $data)) {
            $changes['password_hash'] = Hash::make($data['password']);
        }

        /** @var User $updated */
        $updated = $this->repository->update($user, $changes);

        return $updated;
    }

    public function setActive(int $id, bool $isActive): User
    {
        $user = $this->getById($id);
        /** @var User $updated */
        $updated = $this->repository->update($user, ['is_active' => $isActive, 'updated_at' => now()]);

        return $updated;
    }

    public function signup(SignUpRequest $request): User
    {
        $data = $request->validated();
        $username = trim($data['username']);
        $email = $this->normalizeEmail($data['email']);
        $this->ensureUnique($username, $email);

        /** @var User $user */
        $user = $this->repository->create([
            'username' => $username,
            'email' => $email,
            'password_hash' => Hash::make($data['password']),
            'is_active' => true,
        ]);

        return $user;
    }

    public function getProfile(int $userId): User
    {
        return $this->getById($userId);
    }

    public function updateProfile(int $userId, UpdateProfileRequest $request): User
    {
        $user = $this->getById($userId);
        $data = $request->validated();
        $username = array_key_exists('username', $data) ? trim($data['username']) : $user->username;
        $email = array_key_exists('email', $data) ? $this->normalizeEmail($data['email']) : $user->email;
        $this->ensureUnique($username, $email, $userId);

        /** @var User $updated */
        $updated = $this->repository->update($user, [
            'username' => $username,
            'email' => $email,
            'updated_at' => now(),
        ]);

        return $updated;
    }

    public function changePassword(int $userId, ChangePasswordRequest $request): void
    {
        $user = $this->getById($userId);
        if (! $user->is_active) {
            throw new InvalidPasswordException('The user account is inactive.');
        }
        if (! Hash::check($request->validated('current_password'), $user->password_hash)) {
            throw new InvalidPasswordException('The current password is incorrect.');
        }

        $this->repository->update($user, [
            'password_hash' => Hash::make($request->validated('new_password')),
            'updated_at' => now(),
        ]);
    }

    public function requestPasswordReset(ResetPasswordRequest $request): void
    {
        $user = $this->repository->findByEmail($this->normalizeEmail($request->validated('email')));
        if ($user === null || ! $user->is_active) {
            return;
        }

        $lifetime = max(1, (int) config('users.password_reset_lifetime_minutes', 60));
        $token = Crypt::encryptString(json_encode([
            'user_id' => $user->id,
            'expires_at' => now()->addMinutes($lifetime)->getTimestamp(),
            'password' => hash('sha256', $user->password_hash),
        ], JSON_THROW_ON_ERROR));
        $baseUrl = (string) config('users.password_reset_confirmation_url', '/reset-password');
        $separator = str_contains($baseUrl, '?') ? '&' : '?';
        $url = $baseUrl.$separator.'token='.rawurlencode($token);

        $this->repository->createEmailLog([
            'recipient' => $user->email,
            'subject' => 'Reset your Todo API password',
            'body' => "Use this link to reset your password: {$url}\n\nThis link expires in {$lifetime} minutes.",
            'status' => 'pending',
        ]);
    }

    public function confirmPasswordReset(ConfirmPasswordResetRequest $request): void
    {
        try {
            $payload = json_decode(Crypt::decryptString($request->validated('token')), true, flags: JSON_THROW_ON_ERROR);
            $user = isset($payload['user_id']) ? $this->repository->findById((int) $payload['user_id']) : null;
            $valid = $user instanceof User
                && $user->is_active
                && isset($payload['expires_at'], $payload['password'])
                && (int) $payload['expires_at'] >= now()->getTimestamp()
                && hash_equals((string) $payload['password'], hash('sha256', $user->password_hash));
            if (! $valid) {
                throw new InvalidPasswordResetTokenException('The password reset token is invalid or expired.');
            }
        } catch (DecryptException|JsonException $exception) {
            throw new InvalidPasswordResetTokenException('The password reset token is invalid or expired.', previous: $exception);
        }

        $this->repository->update($user, [
            'password_hash' => Hash::make($request->validated('new_password')),
            'updated_at' => now(),
        ]);
    }

    public function createToken(CreateTokenRequest $request): ?array
    {
        $email = strtolower(trim((string) $request->validated('email')));
        $password = (string) $request->validated('password');
        $user = $this->repository->findByEmail($email);
        $this->dummyPasswordHash ??= Hash::make('not-a-real-password');
        $passwordValid = Hash::check($password, $user?->password_hash ?? $this->dummyPasswordHash);

        if ($user === null || ! $passwordValid || ! $user->is_active) {
            return null;
        }

        $issuedAt = time();
        $expiresIn = max(1, (int) config('users.jwt_token_lifetime_minutes', 60)) * 60;
        $token = JWT::encode([
            'sub' => (string) $user->getKey(),
            'iat' => $issuedAt,
            'exp' => $issuedAt + $expiresIn,
        ], (string) config('users.jwt_secret'), 'HS256');

        return [
            'access_token' => $token,
            'token_type' => 'Bearer',
            'expires_in' => $expiresIn,
        ];
    }

    private function ensureUnique(string $username, string $email, ?int $excludingId = null): void
    {
        if ($this->repository->usernameExists($username, $excludingId)) {
            throw new UserConflictException('Username is already in use.');
        }
        if ($this->repository->emailExists($email, $excludingId)) {
            throw new UserConflictException('Email is already in use.');
        }
    }

    private function normalizeEmail(string $email): string
    {
        return strtolower(trim($email));
    }
}
