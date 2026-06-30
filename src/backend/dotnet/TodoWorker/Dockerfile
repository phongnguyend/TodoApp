# Build stage
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY TodoWorker.csproj .
RUN dotnet restore TodoWorker.csproj

COPY . .
RUN dotnet publish TodoWorker.csproj -c Release -o /app/publish --no-restore

# Runtime stage
FROM mcr.microsoft.com/dotnet/runtime:10.0 AS runtime
WORKDIR /app

COPY --from=build /app/publish .

ENV DOTNET_ENVIRONMENT=Production

# The worker shares the same database volume as the API container.
# Mount the same named volume at /app so todo.db is shared.
ENTRYPOINT ["dotnet", "TodoWorker.dll"]
