package main

import (
	"testing"
)

func TestFormEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic email field",
			input:    `{{ Form::email('email') }}`,
			expected: `<input type="email" name="email" value="">`,
		},
		{
			name:     "Email with value",
			input:    `{{ Form::email('email', 'user@example.com') }}`,
			expected: `<input type="email" name="email" value="{{ 'user@example.com' }}">`,
		},
		{
			name:     "Email with old() helper",
			input:    `{{ Form::email('email', old('email')) }}`,
			expected: `<input type="email" name="email" value="{{ old('email') }}">`,
		},
		{
			name:     "Email with attributes",
			input:    `{{ Form::email('email', old('email'), ['placeholder' => 'Enter email', 'class' => 'form-control']) }}`,
			expected: `<input type="email" name="email" value="{{ old('email') }}" placeholder="Enter email" class="form-control">`,
		},
		{
			name:     "Email with double exclamation marks",
			input:    `{!! Form::email('contact_email', $user->email, ['class' => 'email-input', 'required' => 'required']) !!}`,
			expected: `<input type="email" name="contact_email" value="{{ $user->email }}" class="email-input" required>`,
		},
		{
			name: "Multi-line email field",
			input: `{!! Form::email('business_email', old('business_email'), [
    'placeholder' => 'business@company.com',
    'class' => 'form-control email-field',
    'id' => 'business-email'
]) !!}`,
			expected: `<input type="email" name="business_email" value="{{ old('business_email') }}" placeholder="business@company.com" class="form-control email-field" id="business-email">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormEmail(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
