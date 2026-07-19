<?php

namespace App\Http\Middleware;

use App\Models\User;
use Closure;
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
        $parts = explode('.', $token);
        $secret = (string) config('users.jwt_secret');

        if (count($parts) === 3) {
            [$encodedHeader, $encodedPayload, $signature] = $parts;
            $header = $this->decodeJson($encodedHeader);
            if (($header['alg'] ?? null) !== 'HS256') {
                return null;
            }
            $expected = $this->base64UrlEncode(hash_hmac('sha256', "{$encodedHeader}.{$encodedPayload}", $secret, true));
        } elseif (count($parts) === 2) {
            [$encodedPayload, $signature] = $parts;
            $expected = $this->base64UrlEncode(hash_hmac('sha256', $encodedPayload, $secret, true));
        } else {
            return null;
        }

        if (! hash_equals($expected, $signature)) {
            return null;
        }

        $payload = $this->decodeJson($encodedPayload);
        if ($payload === null || (isset($payload['exp']) && (int) $payload['exp'] < time())) {
            return null;
        }

        return $payload;
    }

    private function decodeJson(string $encoded): ?array
    {
        $decoded = base64_decode(strtr($encoded, '-_', '+/').str_repeat('=', (4 - strlen($encoded) % 4) % 4), true);
        if ($decoded === false) {
            return null;
        }
        $value = json_decode($decoded, true);

        return is_array($value) ? $value : null;
    }

    private function base64UrlEncode(string $value): string
    {
        return rtrim(strtr(base64_encode($value), '+/', '-_'), '=');
    }

    private function unauthorized(string $message): JsonResponse
    {
        return response()->json(['message' => $message], 401);
    }
}
