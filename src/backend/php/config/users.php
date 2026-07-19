<?php

return [
    'jwt_secret' => env('JWT_SECRET_KEY', env('APP_KEY', 'change-me')),
    'password_reset_lifetime_minutes' => env('PASSWORD_RESET_TOKEN_LIFETIME_MINUTES', 60),
    'password_reset_confirmation_url' => env('PASSWORD_RESET_CONFIRMATION_URL', '/reset-password'),
];
