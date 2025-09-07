package ffr

import (
	"testing"
)

func TestFormButton(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic button",
			input:    "{{ Form::button('Click me') }}",
			expected: `<button>{!! Click me !!}</button>`,
		},
		{
			name:     "Button with attributes",
			input:    "{{ Form::button('Submit', ['type' => 'submit', 'class' => 'btn btn-primary', 'id' => 'submit-btn']) }}",
			expected: `<button type="submit" class="btn btn-primary" id="submit-btn">{!! 'Submit' !!}</button>`,
		},
		{
			name:     "Disabled button",
			input:    "{{ Form::button('Disabled Button', ['disabled' => 'disabled']) }}",
			expected: `<button disabled>{!! 'Disabled Button' !!}</button>`,
		},
		{
			name:     "Button with empty text",
			input:    "{{ Form::button('') }}",
			expected: `<button>{!!  !!}</button>`,
		},
		{
			name:     "Button with HTML content",
			input:    "{{ Form::button('<span>HTML Button</span>') }}",
			expected: `<button>{!! <span>HTML Button</span> !!}</button>`,
		},
		{
			name:     "Dynamic disabled attribute (basic)",
			input:    `{!! Form::button('使用する', ['class' => 'btn btn-info', $isDisabled ? 'disabled' : '' => $isDisabled ? 'disabled' : null]) !!}`,
			expected: `<button class="btn btn-info" {{ $isDisabled ? 'disabled' : '' }}="{{ $isDisabled ? 'disabled' : null }}">{!! '使用する' !!}</button>`,
		},
		{
			name:     "Dynamic disabled attribute (user example)",
			input:    `{!! Form::button('使用する', ['class' => 'btn btn-info button-text-color', 'data-toggle' => 'modal', 'data-target' => '#modal', $status ? 'disabled' : '' => $status ? 'disabled' : null]) !!}`,
			expected: `<button class="btn btn-info button-text-color" data-toggle="modal" data-target="#modal" {{ $status ? 'disabled' : '' }}="{{ $status ? 'disabled' : null }}">{!! '使用する' !!}</button>`,
		},
		{
			name:     "Complex dynamic attribute",
			input:    `{!! Form::button('Submit', [$user->isActive() && $user->hasPermission('edit') ? 'disabled' : 'data-action' => $user->isActive() ? 'disabled' : 'edit']) !!}`,
			expected: `<button {{ $user->isActive() && $user->hasPermission('edit') ? 'disabled' : 'data-action' }}="{{ $user->isActive() ? 'disabled' : 'edit' }}">{!! 'Submit' !!}</button>`,
		},
		{
			name:     "Multiple dynamic attributes",
			input:    `{!! Form::button('Test', [$condition1 ? 'disabled' : 'id' => $condition1 ? 'disabled' : 'my-button', $condition2 ? 'data-active' : 'data-inactive' => $condition2 ? 'true' : 'false']) !!}`,
			expected: `<button {{ $condition1 ? 'disabled' : 'id' }}="{{ $condition1 ? 'disabled' : 'my-button' }}" {{ $condition2 ? 'data-active' : 'data-inactive' }}="{{ $condition2 ? 'true' : 'false' }}">{!! 'Test' !!}</button>`,
		},
		{
			name:     "Mixed static and dynamic attributes",
			input:    `{!! Form::button('Mixed', ['class' => 'btn btn-primary', $isDisabled ? 'disabled' : '' => $isDisabled ? 'disabled' : null, 'type' => 'submit']) !!}`,
			expected: `<button type="submit" class="btn btn-primary" {{ $isDisabled ? 'disabled' : '' }}="{{ $isDisabled ? 'disabled' : null }}">{!! 'Mixed' !!}</button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormButton(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
