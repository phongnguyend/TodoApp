using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Moq;
using TodoApi.Controllers;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Tests.Controllers;

public class TokensControllerTests
{
    private readonly Mock<IUserService> _service = new();

    private TokensController Controller() => new(_service.Object)
    {
        ControllerContext = new ControllerContext { HttpContext = new DefaultHttpContext() }
    };

    [Fact]
    public async Task Create_ReturnsTokenAndNoStoreHeaders()
    {
        var request = new TokenRequest("alice@example.com", "password123");
        var response = new TokenResponse("header.payload.signature", "Bearer", 3600);
        _service.Setup(service => service.CreateTokenAsync(request, It.IsAny<CancellationToken>()))
            .ReturnsAsync(response);
        var controller = Controller();

        var result = await controller.Create(request);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(response, ok.Value);
        Assert.Equal("no-store", controller.Response.Headers.CacheControl);
        Assert.Equal("no-cache", controller.Response.Headers.Pragma);
    }

    [Fact]
    public async Task Create_ReturnsNonDisclosingUnauthorizedResponse()
    {
        var request = new TokenRequest("missing@example.com", "wrong");
        _service.Setup(service => service.CreateTokenAsync(request, It.IsAny<CancellationToken>()))
            .ThrowsAsync(new InvalidCredentialsException("Invalid email or password."));
        var controller = Controller();

        var result = await controller.Create(request);

        Assert.IsType<UnauthorizedObjectResult>(result);
        Assert.Equal("Bearer", controller.Response.Headers.WWWAuthenticate);
    }
}
