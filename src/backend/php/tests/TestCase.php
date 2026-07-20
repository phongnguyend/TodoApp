<?php

namespace Tests;

use App\Models\User;
use Firebase\JWT\JWT;
use Illuminate\Foundation\Testing\TestCase as BaseTestCase;

abstract class TestCase extends BaseTestCase
{
    protected function authenticateRequests(): User
    {
        config(['users.jwt_secret' => 'test-secret-at-least-32-bytes-long']);
        $user = User::factory()->create();
        $token = JWT::encode([
            'sub' => (string) $user->id,
            'iat' => time(),
            'exp' => time() + 3600,
        ], (string) config('users.jwt_secret'), 'HS256');
        $this->withHeader('Authorization', 'Bearer '.$token);

        return $user;
    }
}
