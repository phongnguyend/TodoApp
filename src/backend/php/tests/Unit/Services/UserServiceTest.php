<?php

namespace Tests\Unit\Services;

use App\Exceptions\InvalidPasswordException;
use App\Exceptions\InvalidPasswordResetTokenException;
use App\Exceptions\UserConflictException;
use App\Http\Requests\ChangePasswordRequest;
use App\Http\Requests\ConfirmPasswordResetRequest;
use App\Http\Requests\CreateUserRequest;
use App\Http\Requests\ResetPasswordRequest;
use App\Http\Requests\SignUpRequest;
use App\Http\Requests\UpdateProfileRequest;
use App\Http\Requests\UpdateUserRequest;
use App\Models\EmailLog;
use App\Models\User;
use App\Repositories\Contracts\UserRepositoryInterface;
use App\Services\UserService;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Pagination\LengthAwarePaginator;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Hash;
use Mockery;
use Tests\TestCase;

class UserServiceTest extends TestCase
{
    private UserRepositoryInterface $repository;

    private UserService $service;

    protected function setUp(): void
    {
        parent::setUp();
        $this->repository = Mockery::mock(UserRepositoryInterface::class);
        $this->service = new UserService($this->repository);
    }

    protected function tearDown(): void
    {
        Mockery::close();
        parent::tearDown();
    }

    public function test_get_all_clamps_pagination_and_delegates_to_repository(): void
    {
        $paginator = Mockery::mock(LengthAwarePaginator::class);
        $this->repository->shouldReceive('paginate')->with(1, 100)->once()->andReturn($paginator);

        $this->assertSame($paginator, $this->service->getAll(0, 500));
    }

    public function test_get_by_id_returns_user_when_found(): void
    {
        $user = $this->user();
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);

        $this->assertSame($user, $this->service->getById(1));
    }

    public function test_get_by_id_throws_when_user_is_missing(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getById(99);
    }

    public function test_create_normalizes_values_hashes_password_and_uses_default_active_state(): void
    {
        $request = Mockery::mock(CreateUserRequest::class);
        $request->shouldReceive('validated')->once()->andReturn([
            'username' => '  Alice  ',
            'email' => ' Alice@Example.COM ',
            'password' => 'password123',
        ]);
        $this->expectUniqueChecks('Alice', 'alice@example.com');

        $created = $this->user(username: 'Alice', email: 'alice@example.com');
        $this->repository->shouldReceive('create')->once()->with(Mockery::on(
            fn (array $data): bool => $data['username'] === 'Alice'
                && $data['email'] === 'alice@example.com'
                && $data['is_active'] === true
                && Hash::check('password123', $data['password_hash'])
        ))->andReturn($created);

        $this->assertSame($created, $this->service->create($request));
    }

    public function test_create_rejects_a_duplicate_username_before_creating(): void
    {
        $request = Mockery::mock(CreateUserRequest::class);
        $request->shouldReceive('validated')->once()->andReturn([
            'username' => 'Alice', 'email' => 'alice@example.com', 'password' => 'password123',
        ]);
        $this->repository->shouldReceive('usernameExists')->with('Alice', null)->once()->andReturn(true);
        $this->repository->shouldNotReceive('emailExists');
        $this->repository->shouldNotReceive('create');

        $this->expectException(UserConflictException::class);
        $this->service->create($request);
    }

    public function test_create_rejects_a_duplicate_email_before_creating(): void
    {
        $request = Mockery::mock(CreateUserRequest::class);
        $request->shouldReceive('validated')->once()->andReturn([
            'username' => 'Alice', 'email' => 'alice@example.com', 'password' => 'password123',
        ]);
        $this->repository->shouldReceive('usernameExists')->with('Alice', null)->once()->andReturn(false);
        $this->repository->shouldReceive('emailExists')->with('alice@example.com', null)->once()->andReturn(true);
        $this->repository->shouldNotReceive('create');

        $this->expectException(UserConflictException::class);
        $this->service->create($request);
    }

    public function test_update_normalizes_changed_fields_and_hashes_a_new_password(): void
    {
        $existing = $this->user();
        $updated = $this->user(username: 'New name', email: 'new@example.com');
        $request = Mockery::mock(UpdateUserRequest::class);
        $request->shouldReceive('validated')->once()->andReturn([
            'username' => ' New name ', 'email' => ' New@Example.com ', 'password' => 'new-password123',
        ]);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($existing);
        $this->expectUniqueChecks('New name', 'new@example.com', 1);
        $this->repository->shouldReceive('update')->once()->with($existing, Mockery::on(
            fn (array $data): bool => $data['username'] === 'New name'
                && $data['email'] === 'new@example.com'
                && isset($data['updated_at'])
                && Hash::check('new-password123', $data['password_hash'])
        ))->andReturn($updated);

        $this->assertSame($updated, $this->service->update(1, $request));
    }

    public function test_set_active_updates_the_requested_state(): void
    {
        $user = $this->user(active: false);
        $updated = $this->user(active: true);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);
        $this->repository->shouldReceive('update')->once()->with($user, Mockery::on(
            fn (array $data): bool => $data['is_active'] === true && isset($data['updated_at'])
        ))->andReturn($updated);

        $this->assertSame($updated, $this->service->setActive(1, true));
    }

    public function test_signup_always_creates_an_active_user(): void
    {
        $request = Mockery::mock(SignUpRequest::class);
        $request->shouldReceive('validated')->once()->andReturn([
            'username' => 'alice', 'email' => 'alice@example.com', 'password' => 'password123',
        ]);
        $this->expectUniqueChecks('alice', 'alice@example.com');
        $created = $this->user();
        $this->repository->shouldReceive('create')->once()->with(Mockery::on(
            fn (array $data): bool => $data['is_active'] === true
                && Hash::check('password123', $data['password_hash'])
        ))->andReturn($created);

        $this->assertSame($created, $this->service->signup($request));
    }

    public function test_update_profile_changes_only_profile_fields(): void
    {
        $existing = $this->user();
        $updated = $this->user(username: 'new-name');
        $request = Mockery::mock(UpdateProfileRequest::class);
        $request->shouldReceive('validated')->once()->andReturn(['username' => ' new-name ']);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($existing);
        $this->expectUniqueChecks('new-name', 'alice@example.com', 1);
        $this->repository->shouldReceive('update')->once()->with($existing, Mockery::on(
            fn (array $data): bool => $data['username'] === 'new-name'
                && $data['email'] === 'alice@example.com'
                && ! array_key_exists('password_hash', $data)
        ))->andReturn($updated);

        $this->assertSame($updated, $this->service->updateProfile(1, $request));
    }

    public function test_change_password_verifies_current_password_and_stores_a_new_hash(): void
    {
        $user = $this->user(password: 'old-password');
        $request = $this->changePasswordRequest('old-password', 'new-password123');
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);
        $this->repository->shouldReceive('update')->once()->with($user, Mockery::on(
            fn (array $data): bool => Hash::check('new-password123', $data['password_hash'])
                && isset($data['updated_at'])
        ))->andReturn($user);

        $this->service->changePassword(1, $request);
        $this->addToAssertionCount(1);
    }

    public function test_change_password_rejects_an_inactive_user(): void
    {
        $user = $this->user(active: false);
        $request = $this->changePasswordRequest('password123', 'new-password123');
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);
        $this->repository->shouldNotReceive('update');

        $this->expectException(InvalidPasswordException::class);
        $this->service->changePassword(1, $request);
    }

    public function test_change_password_rejects_an_incorrect_current_password(): void
    {
        $user = $this->user(password: 'correct-password');
        $request = $this->changePasswordRequest('wrong-password', 'new-password123');
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);
        $this->repository->shouldNotReceive('update');

        $this->expectException(InvalidPasswordException::class);
        $this->service->changePassword(1, $request);
    }

    public function test_request_password_reset_does_nothing_for_an_unknown_account(): void
    {
        $request = Mockery::mock(ResetPasswordRequest::class);
        $request->shouldReceive('validated')->with('email')->once()->andReturn(' Missing@Example.com ');
        $this->repository->shouldReceive('findByEmail')->with('missing@example.com')->once()->andReturn(null);
        $this->repository->shouldNotReceive('createEmailLog');

        $this->service->requestPasswordReset($request);
        $this->addToAssertionCount(1);
    }

    public function test_request_password_reset_queues_an_encrypted_expiring_token(): void
    {
        config([
            'users.password_reset_lifetime_minutes' => 30,
            'users.password_reset_confirmation_url' => 'https://example.test/reset',
        ]);
        $user = $this->user();
        $request = Mockery::mock(ResetPasswordRequest::class);
        $request->shouldReceive('validated')->with('email')->once()->andReturn('alice@example.com');
        $this->repository->shouldReceive('findByEmail')->with('alice@example.com')->once()->andReturn($user);
        $this->repository->shouldReceive('createEmailLog')->once()->with(Mockery::on(function (array $data) use ($user): bool {
            preg_match('/[?&]token=([^\s]+)/', $data['body'], $matches);
            $payload = json_decode(Crypt::decryptString(rawurldecode($matches[1] ?? '')), true);

            return $data['recipient'] === 'alice@example.com'
                && $data['status'] === 'pending'
                && $payload['user_id'] === 1
                && $payload['password'] === hash('sha256', $user->password_hash)
                && $payload['expires_at'] > time();
        }))->andReturn(new EmailLog);

        $this->service->requestPasswordReset($request);
        $this->addToAssertionCount(1);
    }

    public function test_confirm_password_reset_updates_password_for_a_valid_token(): void
    {
        $user = $this->user();
        $token = $this->resetTokenFor($user);
        $request = $this->confirmPasswordRequest($token, 'reset-password123');
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($user);
        $this->repository->shouldReceive('update')->once()->with($user, Mockery::on(
            fn (array $data): bool => Hash::check('reset-password123', $data['password_hash'])
                && isset($data['updated_at'])
        ))->andReturn($user);

        $this->service->confirmPasswordReset($request);
        $this->addToAssertionCount(1);
    }

    public function test_confirm_password_reset_rejects_a_malformed_token(): void
    {
        $request = $this->confirmPasswordRequest('not-a-token', 'reset-password123');
        $this->repository->shouldNotReceive('update');

        $this->expectException(InvalidPasswordResetTokenException::class);
        $this->service->confirmPasswordReset($request);
    }

    public function test_confirm_password_reset_rejects_a_token_after_the_password_changed(): void
    {
        $original = $this->user(password: 'old-password');
        $changed = $this->user(password: 'different-password');
        $request = $this->confirmPasswordRequest($this->resetTokenFor($original), 'reset-password123');
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($changed);
        $this->repository->shouldNotReceive('update');

        $this->expectException(InvalidPasswordResetTokenException::class);
        $this->service->confirmPasswordReset($request);
    }

    private function expectUniqueChecks(string $username, string $email, ?int $excludingId = null): void
    {
        $this->repository->shouldReceive('usernameExists')->with($username, $excludingId)->once()->andReturn(false);
        $this->repository->shouldReceive('emailExists')->with($email, $excludingId)->once()->andReturn(false);
    }

    private function user(
        string $username = 'alice',
        string $email = 'alice@example.com',
        string $password = 'password123',
        bool $active = true,
    ): User {
        $user = new User([
            'username' => $username,
            'email' => $email,
            'password_hash' => Hash::make($password),
            'is_active' => $active,
        ]);
        $user->id = 1;

        return $user;
    }

    private function changePasswordRequest(string $current, string $new): ChangePasswordRequest
    {
        $request = Mockery::mock(ChangePasswordRequest::class);
        $request->shouldReceive('validated')->with('current_password')->andReturn($current);
        $request->shouldReceive('validated')->with('new_password')->andReturn($new);

        return $request;
    }

    private function confirmPasswordRequest(string $token, string $newPassword): ConfirmPasswordResetRequest
    {
        $request = Mockery::mock(ConfirmPasswordResetRequest::class);
        $request->shouldReceive('validated')->with('token')->andReturn($token);
        $request->shouldReceive('validated')->with('new_password')->andReturn($newPassword);

        return $request;
    }

    private function resetTokenFor(User $user): string
    {
        return Crypt::encryptString(json_encode([
            'user_id' => $user->id,
            'expires_at' => now()->addHour()->getTimestamp(),
            'password' => hash('sha256', $user->password_hash),
        ], JSON_THROW_ON_ERROR));
    }
}
