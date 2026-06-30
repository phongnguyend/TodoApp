# Build stage
FROM node:22-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run prisma:generate
RUN npm run build

# Runtime stage
FROM node:22-alpine AS runtime
WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules/.prisma ./node_modules/.prisma
COPY prisma ./prisma

ENV NODE_ENV=production

# The worker shares the same database as the API container.
# Mount the same volume / point DATABASE_URL at the shared DB.
CMD ["sh", "-c", "npm run prisma:migrate:deploy && node dist/worker/main"]
