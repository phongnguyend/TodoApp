using System.Security.Claims;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Moq;
using TodoApi.Controllers;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Tests.Controllers;

public class UsersControllerTests
{
    private readonly Mock<IUserService> _service = new();

    [Fact]
    public async Task Create_ReturnsConflictForDuplicateUser()
    {
        var request = new CreateUserRequest("alice", "alice@example.com", "password123");
        _service.Setup(service => service.CreateAsync(request, It.IsAny<CancellationToken>()))
            .ThrowsAsync(new UserConflictException("Username is already in use."));
        var controller = new UsersController(_service.Object);

        var result = await controller.Create(request);

        Assert.IsType<ConflictObjectResult>(result);
    }

    [Fact]
    public async Task Activate_ForwardsActiveStateAndReturnsUser()
    {
        var response = new UserResponse(3, "alice", "alice@example.com", true, DateTime.UtcNow, DateTime.UtcNow);
        _service.Setup(service => service.SetActiveAsync(3, true, It.IsAny<CancellationToken>())).ReturnsAsync(response);
        var controller = new UsersController(_service.Object);

        var result = await controller.Activate(3);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(response, ok.Value);
    }

    [Fact]
    public async Task GetProfile_ReturnsUnauthorizedWithoutUserClaim()
    {
        var controller = new UsersController(_service.Object);

        var result = await controller.GetProfile();

        Assert.IsType<UnauthorizedResult>(result);
        _service.Verify(service => service.GetProfileAsync(It.IsAny<int>(), It.IsAny<CancellationToken>()), Times.Never);
    }

    [Fact]
    public async Task GetProfile_UsesNameIdentifierClaim()
    {
        var response = new UserResponse(7, "alice", "alice@example.com", true, DateTime.UtcNow, null);
        _service.Setup(service => service.GetProfileAsync(7, It.IsAny<CancellationToken>())).ReturnsAsync(response);
        var controller = new UsersController(_service.Object)
        {
            ControllerContext = new ControllerContext
            {
                HttpContext = new DefaultHttpContext
                {
                    User = new ClaimsPrincipal(new ClaimsIdentity(
                        [new Claim(ClaimTypes.NameIdentifier, "7")], "test"))
                }
            }
        };

        var result = await controller.GetProfile();

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(response, ok.Value);
    }

    [Fact]
    public async Task RequestPasswordReset_AlwaysReturnsAccepted()
    {
        var request = new ResetPasswordRequest("missing@example.com");
        var controller = new UsersController(_service.Object);

        var result = await controller.RequestPasswordReset(request);

        Assert.IsType<AcceptedResult>(result);
    }
}
