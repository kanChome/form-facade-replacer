package main

import (
	"testing"
)

func TestFormHidden(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic hidden field",
			input:    "{{ Form::hidden('user_id') }}",
			expected: `<input type="hidden" name="user_id" value="{{  }}">`,
		},
		{
			name:     "Hidden field with value",
			input:    "{{ Form::hidden('user_id', $user->id) }}",
			expected: `<input type="hidden" name="user_id" value="{{ $user->id }}">`,
		},
		{
			name:     "Hidden field with attributes",
			input:    "{{ Form::hidden('user_id', $user->id, ['id' => 'user-id', 'class' => 'hidden-field']) }}",
			expected: `<input type="hidden" name="user_id" value="{{ $user->id }}" id="user-id" class="hidden-field">`,
		},
		{
			name:     "Hidden field with array name",
			input:    "{!! Form::hidden('time_ids[]', $timeId) !!}",
			expected: `<input type="hidden" name="time_ids[]" value="{{ $timeId }}">`,
		},
		{
			name:     "Hidden field with old() helper",
			input:    "{!! Form::hidden('time_ids', old('time_ids')) !!}",
			expected: `<input type="hidden" name="time_ids" value="{{ old('time_ids') }}">`,
		},
		{
			name:     "Hidden field with array name and old() helper",
			input:    "{!! Form::hidden('time_ids[]', old('time_ids')) !!}",
			expected: `<input type="hidden" name="time_ids[]" value="{{ old('time_ids') }}">`,
		},
		{
			name:     "Hidden field with session() helper",
			input:    "{{ Form::hidden('session_data', session('data')) }}",
			expected: `<input type="hidden" name="session_data" value="{{ session('data') }}">`,
		},
		{
			name:     "Hidden field with request() helper",
			input:    "{{ Form::hidden('request_data', request('data')) }}",
			expected: `<input type="hidden" name="request_data" value="{{ request('data') }}">`,
		},
		{
			name:     "Hidden field with input() helper",
			input:    "{{ Form::hidden('input_data', input('data')) }}",
			expected: `<input type="hidden" name="input_data" value="{{ input('data') }}">`,
		},
		{
			name:     "Hidden field with null value",
			input:    "{{ Form::hidden('user_id', null) }}",
			expected: `<input type="hidden" name="user_id" value="{{ null }}">`,
		},
		{
			name:     "Hidden field with empty string value",
			input:    "{{ Form::hidden('user_id', '') }}",
			expected: `<input type="hidden" name="user_id" value="{{ '' }}">`,
		},
		{
			name:     "Hidden field with uppercase OLD() helper",
			input:    "{{ Form::hidden('data', OLD('data')) }}",
			expected: `<input type="hidden" name="data" value="{{ OLD('data') }}">`,
		},
		{
			name:     "Hidden field with spaced old() helper",
			input:    "{{ Form::hidden('data', old ('data')) }}",
			expected: `<input type="hidden" name="data" value="{{ old ('data') }}">`,
		},
		{
			name:     "Hidden field with complex array name",
			input:    "{!! Form::hidden('users[0][name]', $user->name) !!}",
			expected: `<input type="hidden" name="users[0][name]" value="{{ $user->name }}">`,
		},
		{
			name:     "Hidden field with old() and array name",
			input:    "{!! Form::hidden('users[]', old('users.0.name')) !!}",
			expected: `<input type="hidden" name="users[]" value="{{ old('users.0.name') }}">`,
		},
		{
			name:     "Hidden field with multiple attributes and old()",
			input:    "{!! Form::hidden('data[]', old('data'), ['class' => 'hidden-data', 'id' => 'data-field']) !!}",
			expected: `<input type="hidden" name="data[]" value="{{ old('data') }}" id="data-field" class="hidden-data">`,
		},
		{
			name:     "Hidden field with old() and default value",
			input:    "{!! Form::hidden('time_ids[]', old('time_ids', $timeIds)) !!}",
			expected: `<input type="hidden" name="time_ids[]" value="{{ old('time_ids', $timeIds) }}">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormHidden(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
