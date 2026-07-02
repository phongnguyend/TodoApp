# Build stage
FROM composer:2 AS vendor
WORKDIR /app

COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --no-autoloader --prefer-dist

COPY . .
RUN composer dump-autoload --optimize --no-dev

# Runtime stage
FROM php:8.4-fpm-alpine
WORKDIR /var/www/html

RUN apk add --no-cache \
    nginx \
    sqlite-dev \
    && docker-php-ext-install pdo pdo_sqlite opcache \
    && rm -rf /var/cache/apk/*

COPY --from=vendor /app /var/www/html
COPY --from=vendor /app/vendor ./vendor

RUN mkdir -p storage/logs storage/framework/cache storage/framework/sessions storage/framework/views storage/app/uploads \
    && chmod -R 775 storage bootstrap/cache \
    && chown -R www-data:www-data storage bootstrap/cache

COPY build/api/nginx.conf /etc/nginx/http.d/default.conf

ENV APP_ENV=production
ENV APP_DEBUG=false

EXPOSE 8080

CMD ["sh", "-c", "php artisan migrate --force && php artisan config:cache && php-fpm -D && nginx -g 'daemon off;'"]
