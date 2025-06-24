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
		{
			name:     "Checkbox with array name",
			input:    "{{ Form::checkbox('tags[]', 'php', true) }}",
			expected: `<input type="checkbox" name="tags[]" value="{{ php }}" @if(in_array(php, (array)true)) checked @endif>`,
		},
		{
			name:     "Checkbox with array name and old() helper",
			input:    "{{ Form::checkbox('categories[]', 'tech', old('categories')) }}",
			expected: `<input type="checkbox" name="categories[]" value="{{ tech }}" @if(in_array(tech, (array)old('categories'))) checked @endif>`,
		},
		{
			name:     "Checkbox with array name and attributes",
			input:    "{{ Form::checkbox('skills[]', 'javascript', false, ['class' => 'skill-checkbox', 'id' => 'skill-js']) }}",
			expected: `<input type="checkbox" name="skills[]" value="{{ javascript }}" @if(in_array(javascript, (array)false)) checked @endif class="skill-checkbox" id="skill-js">`,
		},
		{
			name:     "Checkbox with array name and session() helper",
			input:    "{{ Form::checkbox('preferences[]', 'email', session('user_prefs')) }}",
			expected: `<input type="checkbox" name="preferences[]" value="{{ email }}" @if(in_array(email, (array)session('user_prefs'))) checked @endif>`,
		},
		{
			name:     "Checkbox with array name and complex attributes",
			input:    "{{ Form::checkbox('hobbies[]', 'reading', old('hobbies'), ['class' => 'hobby-check', 'style' => 'margin:5px', 'disabled' => '']) }}",
			expected: `<input type="checkbox" name="hobbies[]" value="{{ reading }}" @if(in_array(reading, (array)old('hobbies'))) checked @endif class="hobby-check" style="margin:5px" disabled>`,
		},
		{
			name:     "Checkbox with double exclamation marks and array name",
			input:    "{!! Form::checkbox('languages[]', 'go', $user->languages) !!}",
			expected: `<input type="checkbox" name="languages[]" value="{{ go }}" @if(in_array(go, (array)$user->languages)) checked @endif>`,
		},
		{
			name:     "Checkbox with complex array name structure",
			input:    "{{ Form::checkbox('users[0][roles][]', 'admin', false) }}",
			expected: `<input type="checkbox" name="users[0][roles][]" value="{{ admin }}" @if(in_array(admin, (array)false)) checked @endif>`,
		},
		{
			name:     "Checkbox without array suffix but with array-like checked value",
			input:    "{{ Form::checkbox('single_option', 'value1', ['value1', 'value2']) }}",
			expected: `<input type="checkbox" name="single_option" value="{{ value1 }}" @if(['value1', 'value2']) checked @endif>`,
		},
		{
			name:     "Multiple checkboxes with same array name",
			input:    "{{ Form::checkbox('colors[]', 'red', old('colors')) }} {{ Form::checkbox('colors[]', 'blue', old('colors')) }}",
			expected: `<input type="checkbox" name="colors[]" value="{{ red }}" @if(in_array(red, (array)old('colors'))) checked @endif> <input type="checkbox" name="colors[]" value="{{ blue }}" @if(in_array(blue, (array)old('colors'))) checked @endif>`,
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