package main

import (
	"testing"
)

func TestFormCheckbox(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic checkbox",
			input:    "{{ Form::checkbox('agree') }}",
			expected: `<input type="checkbox" name="agree" value="{{  }}" @if() checked @endif>`,
		},
		{
			name:     "Checkbox with value and checked",
			input:    "{{ Form::checkbox('agree', 1, true) }}",
			expected: `<input type="checkbox" name="agree" value="{{ 1 }}" @if(true) checked @endif>`,
		},
		{
			name:     "Checkbox with custom value and not checked",
			input:    "{{ Form::checkbox('newsletter', 'yes', false) }}",
			expected: `<input type="checkbox" name="newsletter" value="{{ yes }}" @if(false) checked @endif>`,
		},
		{
			name:     "Checkbox with attributes",
			input:    "{{ Form::checkbox('newsletter', 'yes', false, ['class' => 'form-check-input', 'id' => 'newsletter-check']) }}",
			expected: `<input type="checkbox" name="newsletter" value="{{ yes }}" @if(false) checked @endif class="form-check-input" id="newsletter-check">`,
		},
		{
			name:     "Checkbox with null checked value",
			input:    "{{ Form::checkbox('terms', 1, null) }}",
			expected: `<input type="checkbox" name="terms" value="{{ 1 }}" @if(null) checked @endif>`,
		},
		{
			name:     "Checkbox with disabled attribute",
			input:    "{{ Form::checkbox('newsletter', 'yes', false, ['class' => 'form-check-input', 'disabled' => 'disabled']) }}",
			expected: `<input type="checkbox" name="newsletter" value="{{ yes }}" @if(false) checked @endif class="form-check-input" disabled>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormCheckbox(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}