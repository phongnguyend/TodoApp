<?php

namespace Database\Factories;

use App\Models\TodoItem;
use Illuminate\Database\Eloquent\Factories\Factory;

class TodoItemFactory extends Factory
{
    protected $model = TodoItem::class;

    public function definition(): array
    {
        return [
            'title'        => $this->faker->sentence(4),
            'description'  => $this->faker->optional()->paragraph(),
            'is_completed' => false,
        ];
    }

    public function completed(): static
    {
        return $this->state(['is_completed' => true]);
    }
}
