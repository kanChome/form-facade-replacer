package main

import (
	"testing"
)

func TestFormText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic text field",
			input:    "{{ Form::text('name') }}",
			expected: `<input type="text" name="name" value="">`,
		},
		{
			name:     "Text field with value",
			input:    "{{ Form::text('name', 'John Doe') }}",
			expected: `<input type="text" name="name" value="{{ 'John Doe' }}">`,
		},
		{
			name:     "Text field with variable value",
			input:    "{{ Form::text('name', $user->name) }}",
			expected: `<input type="text" name="name" value="{{ $user->name }}">`,
		},
		{
			name:     "Text field with attributes",
			input:    "{{ Form::text('name', 'John Doe', ['class' => 'form-control', 'placeholder' => 'Enter your name']) }}",
			expected: `<input type="text" name="name" value="{{ 'John Doe' }}" placeholder="Enter your name" class="form-control">`,
		},
		{
			name:     "Text field with all attributes",
			input:    "{{ Form::text('username', $user->username, ['placeholder' => 'Username', 'class' => 'form-input', 'id' => 'username-field']) }}",
			expected: `<input type="text" name="username" value="{{ $user->username }}" placeholder="Username" class="form-input" id="username-field">`,
		},
		{
			name:     "Text field with old() helper",
			input:    "{{ Form::text('name', old('name')) }}",
			expected: `<input type="text" name="name" value="{{ old('name') }}">`,
		},
		{
			name:     "Text field with session() helper",
			input:    "{{ Form::text('data', session('form_data')) }}",
			expected: `<input type="text" name="data" value="{{ session('form_data') }}">`,
		},
		{
			name:     "Text field with request() helper",
			input:    "{{ Form::text('search', request('q')) }}",
			expected: `<input type="text" name="search" value="{{ request('q') }}">`,
		},
		{
			name:     "Text field with input() helper",
			input:    "{{ Form::text('query', input('search_term')) }}",
			expected: `<input type="text" name="query" value="{{ input('search_term') }}">`,
		},
		{
			name:     "Text field with uppercase OLD() helper",
			input:    "{{ Form::text('name', OLD('name')) }}",
			expected: `<input type="text" name="name" value="{{ OLD('name') }}">`,
		},
		{
			name:     "Text field with spaced old() helper",
			input:    "{{ Form::text('email', old ('email')) }}",
			expected: `<input type="text" name="email" value="{{ old ('email') }}">`,
		},
		{
			name:     "Text field with old() and attributes",
			input:    "{{ Form::text('name', old('name'), ['class' => 'form-control', 'placeholder' => 'Your name']) }}",
			expected: `<input type="text" name="name" value="{{ old('name') }}" placeholder="Your name" class="form-control">`,
		},
		{
			name:     "Text field with null value",
			input:    "{{ Form::text('description', null) }}",
			expected: `<input type="text" name="description" value="">`,
		},
		{
			name:     "Text field with empty string value",
			input:    "{{ Form::text('title', '') }}",
			expected: `<input type="text" name="title" value="">`,
		},
		{
			name:     "Text field with double exclamation marks",
			input:    "{!! Form::text('content', $post->content) !!}",
			expected: `<input type="text" name="content" value="{{ $post->content }}">`,
		},
		{
			name:     "Text field with double exclamation marks and old()",
			input:    "{!! Form::text('tags', old('tags')) !!}",
			expected: `<input type="text" name="tags" value="{{ old('tags') }}">`,
		},
		{
			name:     "Text field with double exclamation marks and attributes",
			input:    "{!! Form::text('keywords', $data->keywords, ['class' => 'tag-input', 'id' => 'keywords']) !!}",
			expected: `<input type="text" name="keywords" value="{{ $data->keywords }}" class="tag-input" id="keywords">`,
		},
		{
			name:     "Text field with array name",
			input:    "{{ Form::text('tags[]', $tag) }}",
			expected: `<input type="text" name="tags[]" value="{{ $tag }}">`,
		},
		{
			name:     "Text field with complex array name",
			input:    "{{ Form::text('users[0][name]', $user->name) }}",
			expected: `<input type="text" name="users[0][name]" value="{{ $user->name }}">`,
		},
		{
			name:     "Text field with old() and array name",
			input:    "{{ Form::text('items[]', old('items.0')) }}",
			expected: `<input type="text" name="items[]" value="{{ old('items.0') }}">`,
		},
		{
			name:     "Text field with PHP string concatenation",
			input:    `{{ Form::text('user_data[' . $i . '][name]', old('user_data[' . $i . '][name]'), ['class' => 'form-control']) }}`,
			expected: `<input type="text" name="user_data[{{ $i }}][name]" value="{{ old('user_data[' . $i . '][name]') }}" class="form-control">`,
		},
		{
			name:     "Text field with complex PHP string concatenation",
			input:    `{{ Form::text('data[' . $row['id'] . '][value]', $values[$row['id']] ?? '', ['placeholder' => 'Enter value']) }}`,
			expected: `<input type="text" name="data[{{ $row['id'] }}][value]" value="{{ $values[$row['id']] ?? '' }}" placeholder="Enter value">`,
		},
		{
			name:     "Multi-line text field (user example)",
			input:    `{!! Form::text('date', old('date', $contentsData['targetData']['date']), [
'placeholder' => '日付を入力',
'class' => 'form-control col-4',
]) !!}`,
			expected: `<input type="text" name="date" value="{{ old('date', $contentsData['targetData']['date']) }}" placeholder="日付を入力" class="form-control col-4">`,
		},
		{
			name:     "Multi-line text field with complex attributes",
			input:    `{{ Form::text('name', old('name'), [
    'placeholder' => 'Enter your name',
    'class' => 'form-control',
    'id' => 'name-field'
]) }}`,
			expected: `<input type="text" name="name" value="{{ old('name') }}" placeholder="Enter your name" class="form-control" id="name-field">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormText(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestProcessFormInput(t *testing.T) {
	tests := []struct {
		name      string
		inputType string
		params    []string
		expected  string
	}{
		{
			name:      "Basic text input",
			inputType: "text",
			params:    []string{"'name'"},
			expected:  `<input type="text" name="name" value="">`,
		},
		{
			name:      "Text input with value",
			inputType: "text",
			params:    []string{"'email'", "$user->email"},
			expected:  `<input type="text" name="email" value="{{ $user->email }}">`,
		},
		{
			name:      "Text input with attributes",
			inputType: "text",
			params:    []string{"'phone'", "$user->phone", "['class' => 'form-input', 'placeholder' => 'Phone number']"},
			expected:  `<input type="text" name="phone" value="{{ $user->phone }}" placeholder="Phone number" class="form-input">`,
		},
		{
			name:      "Text input with old() helper",
			inputType: "text",
			params:    []string{"'name'", "old('name')"},
			expected:  `<input type="text" name="name" value="{{ old('name') }}">`,
		},
		{
			name:      "Email input type",
			inputType: "email",
			params:    []string{"'email'", "$user->email"},
			expected:  `<input type="email" name="email" value="{{ $user->email }}">`,
		},
		{
			name:      "Password input type",
			inputType: "password",
			params:    []string{"'password'"},
			expected:  `<input type="password" name="password" value="">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormInput(tt.inputType, tt.params)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
