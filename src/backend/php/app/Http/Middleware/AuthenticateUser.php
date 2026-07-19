<?php

namespace App\Http\Middleware;

use App\Models\User;
use Closure;
use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;
use Throwable;

class AuthenticateUser
{
    public function handle(Request $request, Closure $next): Response
    {
        try {
            $authenticated = $request->user();
        } catch (Throwable) {
            $authenticated = null;
        }

        if ($authenticated instanceof User && $authenticated->getKey() > 0) {
            $request->attributes->set('authenticated_user_id', (int) $authenticated->getKey());

            return $next($request);
        }

        $header = $request->header('Authorization', '');
        if (! preg_match('/^Bearer\s+(\S+)$/i', $header, $matches)) {
            return $this->unauthorized('Authentication required.');
        }

        $payload = $this->decodeToken($matches[1]);
        $userId = filter_var($payload['sub'] ?? null, FILTER_VALIDATE_INT);
        if ($payload === null || $userId === false || $userId < 1) {
            return $this->unauthorized('Invalid or expired token.');
        }

        $request->attributes->set('authenticated_user_id', $userId);

        return $next($request);
    }

    private function decodeToken(string $token): ?array
    {
        try {
            $payload = (array) JWT::decode($token, new Key((string) config('users.jwt_secret'), 'HS256'));
        } catch (Throwable) {
            return null;
        }

        return $payload;
    }

    private function unauthorized(string $message): JsonResponse
    {
        return response()->json(['message' => $message], 401);
    }
}
