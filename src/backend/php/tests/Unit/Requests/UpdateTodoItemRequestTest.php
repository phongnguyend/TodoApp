<?php

namespace Tests\Unit\Requests;

use App\Http\Requests\UpdateTodoItemRequest;
use Illuminate\Support\Facades\Validator;
use Tests\TestCase;

class UpdateTodoItemRequestTest extends TestCase
{
    private function validate(array $data): \Illuminate\Validation\Validator
    {
        return Validator::make($data, (new UpdateTodoItemRequest())->rules());
    }

    public function test_empty_payload_is_valid(): void
    {
        $v = $this->validate([]);
        $this->assertFalse($v->fails());
    }

    public function test_all_fields_provided_passes(): void
    {
        $v = $this->validate([
            'title'        => 'Updated title',
            'description'  => 'Updated desc',
            'is_completed' => true,
        ]);
        $this->assertFalse($v->fails());
    }

    public function test_title_cannot_be_empty_string_when_present(): void
    {
        $v = $this->validate(['title' => '']);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('title', $v->errors()->toArray());
    }

    public function test_title_cannot_exceed_200_characters(): void
    {
        $v = $this->validate(['title' => str_repeat('x', 201)]);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('title', $v->errors()->toArray());
    }

    public function test_description_cannot_exceed_2000_characters(): void
    {
        $v = $this->validate(['description' => str_repeat('x', 2001)]);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('description', $v->errors()->toArray());
    }

    public function test_is_completed_must_be_boolean(): void
    {
        $v = $this->validate(['is_completed' => 'not-a-bool']);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('is_completed', $v->errors()->toArray());
    }

    public function test_is_completed_accepts_false(): void
    {
        $v = $this->validate(['is_completed' => false]);
        $this->assertFalse($v->fails());
    }
}
