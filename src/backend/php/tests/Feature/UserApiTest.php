<?php

namespace Tests\Feature;

use App\Models\EmailLog;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Hash;
use Tests\TestCase;

class UserApiTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->authenticateRequests();
    }

    public function test_index_returns_paginated_users_without_password_hashes(): void
    {
        User::factory()->count(3)->create();

        $this->getJson('/api/users?page_size=2')
            ->assertOk()
            ->assertJsonCount(2, 'data')
            ->assertJsonStructure(['data', 'meta', 'links'])
            ->assertJsonMissingPath('data.0.password_hash');
    }

    public function test_create_normalizes_user_and_hashes_password(): void
    {
        $response = $this->postJson('/api/users', [
            'username' => '  alice  ',
            'email' => 'Alice@Example.COM',
            'password' => 'correct-horse',
        ]);

        $response->assertCreated()
            ->assertJsonPath('data.username', 'alice')
            ->assertJsonPath('data.email', 'alice@example.com')
            ->assertJsonMissingPath('data.password_hash');
        $user = User::where('email', 'alice@example.com')->firstOrFail();
        $this->assertTrue(Hash::check('correct-horse', $user->password_hash));
    }

    public function test_create_rejects_duplicate_username_case_insensitively(): void
    {
        User::factory()->create(['username' => 'Alice']);

        $this->postJson('/api/users', [
            'username' => 'alice',
            'email' => 'another@example.com',
            'password' => 'password123',
        ])->assertConflict();
    }

    public function test_user_management_routes_update_activation_and_missing_users(): void
    {
        $user = User::factory()->inactive()->create();

        $this->getJson("/api/users/{$user->id}")->assertOk()->assertJsonPath('data.is_active', false);
        $this->putJson("/api/users/{$user->id}", ['username' => 'updated'])
            ->assertOk()->assertJsonPath('data.username', 'updated');
        $this->patchJson("/api/users/{$user->id}/activate")
            ->assertOk()->assertJsonPath('data.is_active', true);
        $this->patchJson("/api/users/{$user->id}/deactivate")
            ->assertOk()->assertJsonPath('data.is_active', false);
        $this->getJson('/api/users/9999')->assertNotFound();
    }

    public function test_signup_always_creates_an_active_user(): void
    {
        $this->withHeader('Authorization', '');
        $this->postJson('/api/users/signup', [
            'username' => 'new-user',
            'email' => 'new@example.com',
            'password' => 'password123',
        ])->assertCreated()->assertJsonPath('data.is_active', true);
    }

    public function test_profile_routes_require_authentication(): void
    {
        $this->withHeader('Authorization', '');
        $this->getJson('/api/users/profile')->assertUnauthorized();
        $this->putJson('/api/users/profile', ['username' => 'new-name'])->assertUnauthorized();
        $this->postJson('/api/users/password/change', [
            'current_password' => 'password123',
            'new_password' => 'new-password123',
        ])->assertUnauthorized();
    }

    public function test_authenticated_user_can_read_update_profile_and_change_password(): void
    {
        config(['users.jwt_secret' => 'test-secret-at-least-32-bytes-long']);
        $user = User::factory()->create(['password_hash' => Hash::make('old-password')]);
        $headers = ['Authorization' => 'Bearer '.$this->jwtFor($user->id)];

        $this->withHeaders($headers)->getJson('/api/users/profile')
            ->assertOk()->assertJsonPath('data.id', $user->id);
        $this->withHeaders($headers)->putJson('/api/users/profile', ['username' => 'profile-name'])
            ->assertOk()->assertJsonPath('data.username', 'profile-name');
        $this->withHeaders($headers)->postJson('/api/users/password/change', [
            'current_password' => 'old-password',
            'new_password' => 'new-password123',
        ])->assertNoContent();

        $this->assertTrue(Hash::check('new-password123', $user->fresh()->password_hash));
    }

    public function test_password_change_rejects_incorrect_current_password(): void
    {
        config(['users.jwt_secret' => 'test-secret-at-least-32-bytes-long']);
        $user = User::factory()->create();

        $this->withHeader('Authorization', 'Bearer '.$this->jwtFor($user->id))
            ->postJson('/api/users/password/change', [
                'current_password' => 'wrong-password',
                'new_password' => 'new-password123',
            ])->assertBadRequest();
    }

    public function test_password_reset_does_not_reveal_missing_accounts(): void
    {
        $this->withHeader('Authorization', '');
        $this->postJson('/api/users/password/reset', ['email' => 'missing@example.com'])
            ->assertAccepted();
        $this->assertDatabaseCount('email_logs', 0);
    }

    public function test_password_reset_token_can_be_confirmed_only_once(): void
    {
        $this->withHeader('Authorization', '');
        $user = User::factory()->create(['email' => 'alice@example.com']);
        $this->postJson('/api/users/password/reset', ['email' => 'Alice@Example.com'])
            ->assertAccepted();

        $log = EmailLog::firstOrFail();
        preg_match('/[?&]token=([^\s]+)/', $log->body, $matches);
        $token = rawurldecode($matches[1]);
        $payload = ['token' => $token, 'new_password' => 'reset-password123'];

        $this->postJson('/api/users/password/confirm', $payload)->assertNoContent();
        $this->assertTrue(Hash::check('reset-password123', $user->fresh()->password_hash));
        $this->postJson('/api/users/password/confirm', $payload)->assertBadRequest();
    }

    public function test_token_endpoint_issues_a_signed_jwt_for_active_user(): void
    {
        $this->withHeader('Authorization', '');
        config(['users.jwt_secret' => 'test-secret-at-least-32-bytes-long', 'users.jwt_token_lifetime_minutes' => 60]);
        $user = User::factory()->create([
            'email' => 'alice@example.com',
            'password_hash' => Hash::make('password123'),
            'is_active' => true,
        ]);

        $response = $this->postJson('/api/tokens', [
            'email' => ' Alice@Example.com ',
            'password' => 'password123',
        ])->assertOk()
            ->assertHeader('Cache-Control')
            ->assertHeader('Pragma', 'no-cache')
            ->assertJsonPath('token_type', 'Bearer')
            ->assertJsonPath('expires_in', 3600);

        $this->assertStringContainsString('no-store', $response->headers->get('Cache-Control'));

        [$header, $payload, $signature] = explode('.', $response->json('access_token'));
        $this->assertSame('HS256', json_decode($this->decodeBase64Url($header), true)['alg']);
        $this->assertSame((string) $user->id, json_decode($this->decodeBase64Url($payload), true)['sub']);
        $this->assertNotEmpty($signature);
    }

    public function test_token_endpoint_does_not_disclose_authentication_failure(): void
    {
        $this->withHeader('Authorization', '');
        $this->postJson('/api/tokens', ['email' => 'missing@example.com', 'password' => 'wrong'])
            ->assertUnauthorized()
            ->assertHeader('WWW-Authenticate', 'Bearer')
            ->assertExactJson(['error' => 'Invalid email or password.']);
    }

    private function decodeBase64Url(string $value): string
    {
        return base64_decode(strtr($value, '-_', '+/').str_repeat('=', (4 - strlen($value) % 4) % 4), true);
    }

    private function jwtFor(int $userId): string
    {
        $encode = static fn (string $value): string => rtrim(strtr(base64_encode($value), '+/', '-_'), '=');
        $header = $encode(json_encode(['alg' => 'HS256', 'typ' => 'JWT'], JSON_THROW_ON_ERROR));
        $payload = $encode(json_encode(['sub' => (string) $userId, 'exp' => time() + 300], JSON_THROW_ON_ERROR));
        $signature = $encode(hash_hmac('sha256', "{$header}.{$payload}", 'test-secret-at-least-32-bytes-long', true));

        return "{$header}.{$payload}.{$signature}";
    }
}
