package main

import (
	"testing"
)

func TestFormNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic number field",
			input:    "{{ Form::number('age') }}",
			expected: `<input type="number" name="age">`,
		},
		{
			name:     "Number field with value",
			input:    "{{ Form::number('age', 25) }}",
			expected: `<input type="number" name="age" value="{{ 25 }}">`,
		},
		{
			name:     "Number field with attributes",
			input:    "{{ Form::number('age', 25, ['class' => 'form-control', 'min' => 18, 'max' => 100]) }}",
			expected: `<input type="number" name="age" value="{{ 25 }}" class="form-control" min="18" max="100">`,
		},
		{
			name:     "Number field with null value",
			input:    "{{ Form::number('age', null) }}",
			expected: `<input type="number" name="age">`,
		},
		{
			name:     "Number field with empty string value",
			input:    "{{ Form::number('age', '') }}",
			expected: `<input type="number" name="age">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormNumber(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}