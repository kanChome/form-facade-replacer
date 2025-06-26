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
		{
			name:     "Number field with old() helper",
			input:    "{{ Form::number('quantity', old('quantity')) }}",
			expected: `<input type="number" name="quantity" value="{{ old('quantity') }}">`,
		},
		{
			name:     "Number field with session() helper",
			input:    "{{ Form::number('score', session('user_score')) }}",
			expected: `<input type="number" name="score" value="{{ session('user_score') }}">`,
		},
		{
			name:     "Number field with request() helper",
			input:    "{{ Form::number('amount', request('amount')) }}",
			expected: `<input type="number" name="amount" value="{{ request('amount') }}">`,
		},
		{
			name:     "Number field with input() helper",
			input:    "{{ Form::number('price', input('price')) }}",
			expected: `<input type="number" name="price" value="{{ input('price') }}">`,
		},
		{
			name:     "Number field with uppercase OLD() helper",
			input:    "{{ Form::number('count', OLD('count')) }}",
			expected: `<input type="number" name="count" value="{{ OLD('count') }}">`,
		},
		{
			name:     "Number field with spaced old() helper",
			input:    "{{ Form::number('value', old ('value')) }}",
			expected: `<input type="number" name="value" value="{{ old ('value') }}">`,
		},
		{
			name:     "Number field with old() and attributes",
			input:    "{{ Form::number('age', old('age'), ['class' => 'form-control', 'min' => 0, 'max' => 120]) }}",
			expected: `<input type="number" name="age" value="{{ old('age') }}" class="form-control" min="0" max="120">`,
		},
		{
			name:     "Number field with session() and step attribute",
			input:    "{{ Form::number('decimal', session('decimal_value'), ['step' => 1, 'min' => 0]) }}",
			expected: `<input type="number" name="decimal" value="{{ session('decimal_value') }}" min="0" step="1">`,
		},
		{
			name:     "Number field with double exclamation marks and old()",
			input:    "{!! Form::number('rating', old('rating')) !!}",
			expected: `<input type="number" name="rating" value="{{ old('rating') }}">`,
		},
		{
			name:     "Number field with double exclamation marks and attributes",
			input:    "{!! Form::number('stars', old('stars'), ['class' => 'rating-input', 'min' => 1, 'max' => 5]) !!}",
			expected: `<input type="number" name="stars" value="{{ old('stars') }}" class="rating-input" min="1" max="5">`,
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
