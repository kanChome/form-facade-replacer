package main

import (
	"testing"
)

func TestFormTextarea(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic textarea",
			input:    "{{ Form::textarea('message') }}",
			expected: `<textarea name="message"></textarea>`,
		},
		{
			name:     "Textarea with value",
			input:    "{{ Form::textarea('message', 'Default content') }}",
			expected: `<textarea name="message">{{ 'Default content' }}</textarea>`,
		},
		{
			name:     "Textarea with attributes",
			input:    "{{ Form::textarea('message', 'Default content', ['class' => 'form-control', 'rows' => 5, 'cols' => 30]) }}",
			expected: `<textarea name="message" cols="30" rows="5" class="form-control">{{ 'Default content' }}</textarea>`,
		},
		{
			name:     "Textarea with null value and attributes",
			input:    "{{ Form::textarea('message', null, ['class' => 'form-control']) }}",
			expected: `<textarea name="message" class="form-control">{{ null }}</textarea>`,
		},
		{
			name:     "Textarea with empty string value and attributes",
			input:    "{{ Form::textarea('message', '', ['rows' => 4]) }}",
			expected: `<textarea name="message" rows="4">{{ '' }}</textarea>`,
		},
		{
			name:     "Textarea with old() helper",
			input:    "{{ Form::textarea('description', old('description')) }}",
			expected: `<textarea name="description">{!! json_encode(old('description')) !!}</textarea>`,
		},
		{
			name:     "Textarea with session() helper",
			input:    "{{ Form::textarea('content', session('draft_content')) }}",
			expected: `<textarea name="content">{!! json_encode(session('draft_content')) !!}</textarea>`,
		},
		{
			name:     "Textarea with request() helper",
			input:    "{{ Form::textarea('notes', request('notes')) }}",
			expected: `<textarea name="notes">{!! json_encode(request('notes')) !!}</textarea>`,
		},
		{
			name:     "Textarea with input() helper",
			input:    "{{ Form::textarea('comments', input('comment_text')) }}",
			expected: `<textarea name="comments">{!! json_encode(input('comment_text')) !!}</textarea>`,
		},
		{
			name:     "Textarea with uppercase OLD() helper",
			input:    "{{ Form::textarea('bio', OLD('bio')) }}",
			expected: `<textarea name="bio">{!! json_encode(OLD('bio')) !!}</textarea>`,
		},
		{
			name:     "Textarea with spaced old() helper",
			input:    "{{ Form::textarea('message', old ('message')) }}",
			expected: `<textarea name="message">{!! json_encode(old ('message')) !!}</textarea>`,
		},
		{
			name:     "Textarea with old() and attributes",
			input:    "{{ Form::textarea('description', old('description'), ['class' => 'form-control', 'rows' => 5]) }}",
			expected: `<textarea name="description" rows="5" class="form-control">{!! json_encode(old('description')) !!}</textarea>`,
		},
		{
			name:     "Textarea with session() and attributes",
			input:    "{{ Form::textarea('draft', session('draft'), ['cols' => 50, 'rows' => 8, 'class' => 'draft-area']) }}",
			expected: `<textarea name="draft" cols="50" rows="8" class="draft-area">{!! json_encode(session('draft')) !!}</textarea>`,
		},
		{
			name:     "Textarea with double exclamation marks and old()",
			input:    "{!! Form::textarea('content', old('content')) !!}",
			expected: `<textarea name="content">{!! json_encode(old('content')) !!}</textarea>`,
		},
		{
			name:     "Textarea with double exclamation marks and attributes",
			input:    "{!! Form::textarea('body', old('body'), ['class' => 'editor', 'placeholder' => 'Enter content']) !!}",
			expected: `<textarea name="body" placeholder="Enter content" class="editor">{!! json_encode(old('body')) !!}</textarea>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormTextarea(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}