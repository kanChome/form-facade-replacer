package main

import (
	"testing"
)

func TestFormUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic URL field",
			input:    `{{ Form::url('website') }}`,
			expected: `<input type="url" name="website" value="">`,
		},
		{
			name:     "URL with value",
			input:    `{{ Form::url('website', 'https://example.com') }}`,
			expected: `<input type="url" name="website" value="{{ 'https://example.com' }}">`,
		},
		{
			name:     "URL with old() helper",
			input:    `{{ Form::url('website', old('website')) }}`,
			expected: `<input type="url" name="website" value="{{ old('website') }}">`,
		},
		{
			name:     "URL with attributes",
			input:    `{{ Form::url('website', old('website'), ['placeholder' => 'https://example.com', 'class' => 'form-control']) }}`,
			expected: `<input type="url" name="website" value="{{ old('website') }}" placeholder="https://example.com" class="form-control">`,
		},
		{
			name:     "URL with double exclamation marks",
			input:    `{!! Form::url('company_url', $company->website, ['class' => 'url-input', 'required' => 'required']) !!}`,
			expected: `<input type="url" name="company_url" value="{{ $company->website }}" class="url-input" required>`,
		},
		{
			name: "Multi-line URL field",
			input: `{!! Form::url('social_link', old('social_link'), [
    'placeholder' => 'https://social.com/profile',
    'class' => 'form-control social-url',
    'id' => 'social-link-input'
]) !!}`,
			expected: `<input type="url" name="social_link" value="{{ old('social_link') }}" placeholder="https://social.com/profile" class="form-control social-url" id="social-link-input">`,
		},
		{
			name:     "URL with variable value",
			input:    `{{ Form::url('portfolio_url', $user->portfolio_url, ['class' => 'portfolio-link']) }}`,
			expected: `<input type="url" name="portfolio_url" value="{{ $user->portfolio_url }}" class="portfolio-link">`,
		},
		{
			name:     "URL with complex nested name",
			input:    `{{ Form::url('links[' . $type . '][url]', old('links[' . $type . '][url]'), ['class' => 'link-input']) }}`,
			expected: `<input type="url" name="links[{{ $type }}][url]" value="{{ old('links[' . $type . '][url]') }}" class="link-input">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormUrl(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
