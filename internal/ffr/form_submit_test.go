package ffr

import (
	"testing"
)

func TestFormSubmit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic submit button",
			input:    "{{ Form::submit('Submit') }}",
			expected: `<button type="submit">Submit</button>`,
		},
		{
			name:     "Submit with attributes",
			input:    "{{ Form::submit('Save Changes', ['class' => 'btn btn-success', 'id' => 'save-btn']) }}",
			expected: `<button type="submit" class="btn btn-success" id="save-btn">Save Changes</button>`,
		},
		{
			name:     "Disabled submit button",
			input:    "{{ Form::submit('Process', ['disabled' => 'disabled']) }}",
			expected: `<button type="submit" disabled>Process</button>`,
		},
		{
			name:     "Submit with empty value",
			input:    "{{ Form::submit('') }}",
			expected: `<button type="submit"></button>`,
		},
		{
			name:     "Submit with null value",
			input:    "{{ Form::submit(null) }}",
			expected: `<button type="submit"></button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormSubmit(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
