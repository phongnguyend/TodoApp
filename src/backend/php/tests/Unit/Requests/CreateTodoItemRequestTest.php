<?php

namespace Tests\Unit\Requests;

use App\Http\Requests\CreateTodoItemRequest;
use Illuminate\Support\Facades\Validator;
use Tests\TestCase;

class CreateTodoItemRequestTest extends TestCase
{
    private function validate(array $data): \Illuminate\Validation\Validator
    {
        return Validator::make($data, (new CreateTodoItemRequest())->rules());
    }

    public function test_valid_data_passes(): void
    {
        $v = $this->validate(['title' => 'Buy milk', 'description' => 'From store']);
        $this->assertFalse($v->fails());
    }

    public function test_description_is_optional(): void
    {
        $v = $this->validate(['title' => 'Buy milk']);
        $this->assertFalse($v->fails());
    }

    public function test_title_is_required(): void
    {
        $v = $this->validate(['description' => 'no title']);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('title', $v->errors()->toArray());
    }

    public function test_title_cannot_be_empty_string(): void
    {
        $v = $this->validate(['title' => '']);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('title', $v->errors()->toArray());
    }

    public function test_title_cannot_exceed_200_characters(): void
    {
        $v = $this->validate(['title' => str_repeat('a', 201)]);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('title', $v->errors()->toArray());
    }

    public function test_title_at_max_length_passes(): void
    {
        $v = $this->validate(['title' => str_repeat('a', 200)]);
        $this->assertFalse($v->fails());
    }

    public function test_description_cannot_exceed_2000_characters(): void
    {
        $v = $this->validate(['title' => 'Valid', 'description' => str_repeat('a', 2001)]);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('description', $v->errors()->toArray());
    }

    public function test_description_at_max_length_passes(): void
    {
        $v = $this->validate(['title' => 'Valid', 'description' => str_repeat('a', 2000)]);
        $this->assertFalse($v->fails());
    }

    public function test_description_can_be_null(): void
    {
        $v = $this->validate(['title' => 'Valid', 'description' => null]);
        $this->assertFalse($v->fails());
    }
}
