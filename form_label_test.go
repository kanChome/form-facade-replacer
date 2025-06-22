package main

import (
	"testing"
)

func TestFormLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic label with text",
			input:    "{{ Form::label('name', 'Name') }}",
			expected: `<label for="name">{!! 'Name' !!}</label>`,
		},
		{
			name:     "Label with attributes",
			input:    "{{ Form::label('name', 'Full Name', ['class' => 'form-label', 'for' => 'user-name']) }}",
			expected: `<label for="user-name" class="form-label">{!! 'Full Name' !!}</label>`,
		},
		{
			name:     "Label without text",
			input:    "{{ Form::label('email') }}",
			expected: `<label for="email">{!! 'email' !!}</label>`,
		},
		{
			name:     "Label with empty text",
			input:    "{{ Form::label('password', '') }}",
			expected: `<label for="password">{!! '' !!}</label>`,
		},
		{
			name:     "Label with null text",
			input:    "{{ Form::label('password', null) }}",
			expected: `<label for="password">{!! null !!}</label>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormLabel(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}