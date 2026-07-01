# Build stage
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY TodoApi/TodoApi.csproj TodoApi/
COPY TodoShared/TodoShared.csproj TodoShared/
RUN dotnet restore TodoApi/TodoApi.csproj

COPY TodoApi/ TodoApi/
COPY TodoShared/ TodoShared/
RUN dotnet publish TodoApi/TodoApi.csproj -c Release -o /app/publish --no-restore

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime
WORKDIR /app

COPY --from=build /app/publish .

ENV ASPNETCORE_URLS=http://+:8080
ENV ASPNETCORE_ENVIRONMENT=Production

EXPOSE 8080

ENTRYPOINT ["dotnet", "TodoApi.dll"]
