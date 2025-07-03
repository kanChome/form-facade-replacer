package main

import (
	"testing"
)

func TestFormInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic time input (user example)",
			input:    `{{ Form::input('time', 'start_at_time', old('start_at_time', '00:00'), ['placeholder' => '時間を入力', 'class' => 'form-control col-4']) }}`,
			expected: `<input type="time" name="start_at_time" value="{{ old('start_at_time', '00:00') }}" placeholder="時間を入力" class="form-control col-4">`,
		},
		{
			name:     "Email input type",
			input:    `{{ Form::input('email', 'email', old('email'), ['placeholder' => 'Enter email', 'class' => 'form-control']) }}`,
			expected: `<input type="email" name="email" value="{{ old('email') }}" placeholder="Enter email" class="form-control">`,
		},
		{
			name:     "Password input type",
			input:    `{{ Form::input('password', 'password', '', ['class' => 'form-control']) }}`,
			expected: `<input type="password" name="password" value="" class="form-control">`,
		},
		{
			name:     "Date input type",
			input:    `{{ Form::input('date', 'birth_date', old('birth_date'), ['class' => 'form-control']) }}`,
			expected: `<input type="date" name="birth_date" value="{{ old('birth_date') }}" class="form-control">`,
		},
		{
			name:     "URL input type",
			input:    `{{ Form::input('url', 'website', old('website'), ['placeholder' => 'https://example.com']) }}`,
			expected: `<input type="url" name="website" value="{{ old('website') }}" placeholder="https://example.com">`,
		},
		{
			name:     "Tel input type",
			input:    `{{ Form::input('tel', 'phone', old('phone'), ['placeholder' => '+1234567890']) }}`,
			expected: `<input type="tel" name="phone" value="{{ old('phone') }}" placeholder="+1234567890">`,
		},
		{
			name:     "Search input type",
			input:    `{{ Form::input('search', 'query', old('query'), ['placeholder' => 'Search...']) }}`,
			expected: `<input type="search" name="query" value="{{ old('query') }}" placeholder="Search...">`,
		},
		{
			name:     "Range input type",
			input:    `{{ Form::input('range', 'volume', '50', ['min' => '0', 'max' => '100']) }}`,
			expected: `<input type="range" name="volume" value="{{ 50 }}">`,
		},
		{
			name:     "Color input type",
			input:    `{{ Form::input('color', 'color', '#ff0000', ['class' => 'color-picker']) }}`,
			expected: `<input type="color" name="color" value="{{ #ff0000 }}" class="color-picker">`,
		},
		{
			name:     "Input with double exclamation marks",
			input:    `{!! Form::input('datetime-local', 'datetime', old('datetime'), ['class' => 'form-control']) !!}`,
			expected: `<input type="datetime-local" name="datetime" value="{{ old('datetime') }}" class="form-control">`,
		},
		{
			name:     "Input with null value",
			input:    `{{ Form::input('text', 'name', null, ['placeholder' => 'Enter name']) }}`,
			expected: `<input type="text" name="name" value="" placeholder="Enter name">`,
		},
		{
			name:     "Input with empty string value",
			input:    `{{ Form::input('text', 'title', '', ['class' => 'form-control']) }}`,
			expected: `<input type="text" name="title" value="" class="form-control">`,
		},
		{
			name:     "Input with variable value",
			input:    `{{ Form::input('text', 'username', $user->username, ['id' => 'username']) }}`,
			expected: `<input type="text" name="username" value="{{ $user->username }}" id="username">`,
		},
		{
			name:     "Input with PHP string concatenation",
			input:    `{{ Form::input('text', 'items[' . $i . '][name]', old('items[' . $i . '][name]'), ['class' => 'item-name']) }}`,
			expected: `<input type="text" name="items[{{ $i }}][name]" value="{{ old('items[' . $i . '][name]') }}" class="item-name">`,
		},
		{
			name:     "Multi-line input field",
			input:    `{!! Form::input('time', 'start_time', old('start_time'), [
    'placeholder' => '開始時間',
    'class' => 'form-control',
    'id' => 'start-time-input'
]) !!}`,
			expected: `<input type="time" name="start_time" value="{{ old('start_time') }}" placeholder="開始時間" class="form-control" id="start-time-input">`,
		},
		{
			name:     "Input with all common attributes",
			input:    `{{ Form::input('email', 'contact_email', old('contact_email'), ['placeholder' => 'your@email.com', 'class' => 'form-control email-input', 'id' => 'contact-email-field']) }}`,
			expected: `<input type="email" name="contact_email" value="{{ old('contact_email') }}" placeholder="your@email.com" class="form-control email-input" id="contact-email-field">`,
		},
		{
			name:     "Input with complex nested value",
			input:    `{{ Form::input('text', 'data[user][profile][name]', $userData['profile']['name'] ?? '', ['class' => 'profile-field']) }}`,
			expected: `<input type="text" name="data[user][profile][name]" value="{{ $userData['profile']['name'] ?? '' }}" class="profile-field">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormInput(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestProcessFormInputDynamic(t *testing.T) {
	tests := []struct {
		name     string
		params   []string
		expected string
	}{
		{
			name:     "Time input with user's example",
			params:   []string{"'time'", "'start_at_time'", "old('start_at_time', '00:00')", "['placeholder' => '時間を入力', 'class' => 'form-control col-4']"},
			expected: `<input type="time" name="start_at_time" value="{{ old('start_at_time', '00:00') }}" placeholder="時間を入力" class="form-control col-4">`,
		},
		{
			name:     "Email input type",
			params:   []string{"'email'", "'email'", "$user->email"},
			expected: `<input type="email" name="email" value="{{ $user->email }}">`,
		},
		{
			name:     "Password input type",
			params:   []string{"'password'", "'password'"},
			expected: `<input type="password" name="password" value="">`,
		},
		{
			name:     "Date input with attributes",
			params:   []string{"'date'", "'birth_date'", "old('birth_date')", "['class' => 'form-control', 'placeholder' => 'YYYY-MM-DD']"},
			expected: `<input type="date" name="birth_date" value="{{ old('birth_date') }}" placeholder="YYYY-MM-DD" class="form-control">`,
		},
		{
			name:     "URL input type",
			params:   []string{"'url'", "'website'", "$company->website", "['placeholder' => 'https://example.com']"},
			expected: `<input type="url" name="website" value="{{ $company->website }}" placeholder="https://example.com">`,
		},
		{
			name:     "Invalid params (less than 2)",
			params:   []string{"'text'"},
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormInputDynamic(tt.params)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}