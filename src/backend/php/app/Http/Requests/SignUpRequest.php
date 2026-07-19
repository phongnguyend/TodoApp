<?php

namespace App\Http\Requests;

class SignUpRequest extends CreateUserRequest
{
    public function rules(): array
    {
        $rules = parent::rules();
        unset($rules['is_active']);

        return $rules;
    }
}
