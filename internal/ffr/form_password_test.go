package ffr

import (
	"testing"
)

func TestFormPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic password field",
			input:    `{{ Form::password('password') }}`,
			expected: `<input type="password" name="password" value="">`,
		},
		{
			name:     "Password with attributes",
			input:    `{{ Form::password('password', ['class' => 'form-control', 'placeholder' => 'Enter password']) }}`,
			expected: `<input type="password" name="password" value="" placeholder="Enter password" class="form-control">`,
		},
		{
			name:     "Password with double exclamation marks",
			input:    `{!! Form::password('current_password', ['class' => 'form-control', 'required' => 'required']) !!}`,
			expected: `<input type="password" name="current_password" value="" class="form-control" required>`,
		},
		{
			name: "Multi-line password field",
			input: `{{ Form::password('new_password', [
    'class' => 'form-control password-field',
    'placeholder' => 'New password',
    'id' => 'new-password',
    'minlength' => '8'
]) }}`,
			expected: `<input type="password" name="new_password" value="" placeholder="New password" class="form-control password-field" id="new-password">`,
		},
		{
			name:     "Password with complex name",
			input:    `{{ Form::password('user[password]', ['class' => 'user-password']) }}`,
			expected: `<input type="password" name="user[password]" value="" class="user-password">`,
		},
		{
			name:     "Password with PHP concatenation",
			input:    `{{ Form::password('passwords[' . $index . '][value]', ['class' => 'dynamic-password']) }}`,
			expected: `<input type="password" name="passwords[{{ $index }}][value]" value="" class="dynamic-password">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormPassword(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
