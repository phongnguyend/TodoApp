<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\ChangePasswordRequest;
use App\Http\Requests\ConfirmPasswordResetRequest;
use App\Http\Requests\CreateTokenRequest;
use App\Http\Requests\CreateUserRequest;
use App\Http\Requests\ResetPasswordRequest;
use App\Http\Requests\SignUpRequest;
use App\Http\Requests\UpdateProfileRequest;
use App\Http\Requests\UpdateUserRequest;
use App\Http\Resources\UserResource;
use App\Services\Contracts\UserServiceInterface;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;
use Illuminate\Http\Response;

class UserController extends Controller
{
    public function __construct(
        private readonly UserServiceInterface $service
    ) {}

    public function index(Request $request): AnonymousResourceCollection
    {
        return UserResource::collection($this->service->getAll(
            (int) $request->query('page', 1),
            (int) $request->query('page_size', 20),
        ));
    }

    public function show(int $id): UserResource
    {
        return new UserResource($this->service->getById($id));
    }

    public function store(CreateUserRequest $request): JsonResponse
    {
        return (new UserResource($this->service->create($request)))
            ->response()
            ->setStatusCode(201);
    }

    public function update(UpdateUserRequest $request, int $id): UserResource
    {
        return new UserResource($this->service->update($id, $request));
    }

    public function activate(int $id): UserResource
    {
        return new UserResource($this->service->setActive($id, true));
    }

    public function deactivate(int $id): UserResource
    {
        return new UserResource($this->service->setActive($id, false));
    }

    public function signup(SignUpRequest $request): JsonResponse
    {
        return (new UserResource($this->service->signup($request)))
            ->response()
            ->setStatusCode(201);
    }

    public function profile(Request $request): UserResource
    {
        return new UserResource($this->service->getProfile($this->currentUserId($request)));
    }

    public function updateProfile(UpdateProfileRequest $request): UserResource
    {
        return new UserResource($this->service->updateProfile($this->currentUserId($request), $request));
    }

    public function changePassword(ChangePasswordRequest $request): Response
    {
        $this->service->changePassword($this->currentUserId($request), $request);

        return response()->noContent();
    }

    public function requestPasswordReset(ResetPasswordRequest $request): JsonResponse
    {
        $this->service->requestPasswordReset($request);

        return response()->json([
            'message' => 'If the account exists, a password reset email has been queued.',
        ], 202);
    }

    public function confirmPasswordReset(ConfirmPasswordResetRequest $request): Response
    {
        $this->service->confirmPasswordReset($request);

        return response()->noContent();
    }

    public function createToken(CreateTokenRequest $request): JsonResponse
    {
        $token = $this->service->createToken($request);
        if ($token === null) {
            return response()->json(['error' => 'Invalid email or password.'], 401)
                ->header('WWW-Authenticate', 'Bearer');
        }

        return response()->json($token)
            ->header('Cache-Control', 'no-store')
            ->header('Pragma', 'no-cache');
    }

    private function currentUserId(Request $request): int
    {
        return (int) $request->attributes->get('authenticated_user_id');
    }
}
