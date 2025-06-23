package main

import (
	"testing"
)

func TestFormTextarea(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic textarea",
			input:    "{{ Form::textarea('message') }}",
			expected: `<textarea name="message"></textarea>`,
		},
		{
			name:     "Textarea with value",
			input:    "{{ Form::textarea('message', 'Default content') }}",
			expected: `<textarea name="message">{{ 'Default content' }}</textarea>`,
		},
		{
			name:     "Textarea with attributes",
			input:    "{{ Form::textarea('message', 'Default content', ['class' => 'form-control', 'rows' => 5, 'cols' => 30]) }}",
			expected: `<textarea name="message" cols="30" rows="5" class="form-control">{{ 'Default content' }}</textarea>`,
		},
		{
			name:     "Textarea with null value and attributes",
			input:    "{{ Form::textarea('message', null, ['class' => 'form-control']) }}",
			expected: `<textarea name="message" class="form-control">{{ null }}</textarea>`,
		},
		{
			name:     "Textarea with empty string value and attributes",
			input:    "{{ Form::textarea('message', '', ['rows' => 4]) }}",
			expected: `<textarea name="message" rows="4">{{ '' }}</textarea>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormTextarea(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}