package main

import (
	"testing"
)

func TestFormButton(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic button",
			input:    "{{ Form::button('Click me') }}",
			expected: `<button>{!! Click me !!}</button>`,
		},
		{
			name:     "Button with attributes",
			input:    "{{ Form::button('Submit', ['type' => 'submit', 'class' => 'btn btn-primary', 'id' => 'submit-btn']) }}",
			expected: `<button type="submit" class="btn btn-primary" id="submit-btn">{!! 'Submit' !!}</button>`,
		},
		{
			name:     "Disabled button",
			input:    "{{ Form::button('Disabled Button', ['disabled' => 'disabled']) }}",
			expected: `<button disabled>{!! 'Disabled Button' !!}</button>`,
		},
		{
			name:     "Button with empty text",
			input:    "{{ Form::button('') }}",
			expected: `<button>{!!  !!}</button>`,
		},
		{
			name:     "Button with HTML content",
			input:    "{{ Form::button('<span>HTML Button</span>') }}",
			expected: `<button>{!! <span>HTML Button</span> !!}</button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormButton(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}