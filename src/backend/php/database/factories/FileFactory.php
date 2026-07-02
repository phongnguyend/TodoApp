<?php

namespace Database\Factories;

use App\Models\File;
use Illuminate\Database\Eloquent\Factories\Factory;

class FileFactory extends Factory
{
    protected $model = File::class;

    public function definition(): array
    {
        $name = $this->faker->slug() . '.txt';

        return [
            'name'         => $name,
            'extension'    => 'txt',
            'size'         => $this->faker->numberBetween(10, 10000),
            'content_type' => 'text/plain',
            'location'     => storage_path('app/uploads/' . uniqid('', true) . '_' . $name),
        ];
    }
}
