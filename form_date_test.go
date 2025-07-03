package main

import (
	"testing"
)

func TestFormDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic date field",
			input:    `{{ Form::date('birth_date') }}`,
			expected: `<input type="date" name="birth_date" value="">`,
		},
		{
			name:     "Date with value",
			input:    `{{ Form::date('birth_date', '1990-01-01') }}`,
			expected: `<input type="date" name="birth_date" value="{{ '1990-01-01' }}">`,
		},
		{
			name:     "Date with old() helper",
			input:    `{{ Form::date('birth_date', old('birth_date')) }}`,
			expected: `<input type="date" name="birth_date" value="{{ old('birth_date') }}">`,
		},
		{
			name:     "Date with attributes",
			input:    `{{ Form::date('birth_date', old('birth_date'), ['class' => 'form-control', 'max' => '2023-12-31']) }}`,
			expected: `<input type="date" name="birth_date" value="{{ old('birth_date') }}" class="form-control">`,
		},
		{
			name:     "Date with double exclamation marks",
			input:    `{!! Form::date('event_date', $event->date, ['class' => 'date-picker', 'required' => 'required']) !!}`,
			expected: `<input type="date" name="event_date" value="{{ $event->date }}" class="date-picker" required>`,
		},
		{
			name: "Multi-line date field",
			input: `{!! Form::date('appointment_date', old('appointment_date'), [
    'class' => 'form-control appointment-date',
    'id' => 'appointment-date-picker',
    'min' => '2023-01-01'
]) !!}`,
			expected: `<input type="date" name="appointment_date" value="{{ old('appointment_date') }}" class="form-control appointment-date" id="appointment-date-picker">`,
		},
		{
			name:     "Date with Carbon object",
			input:    `{{ Form::date('created_at', $model->created_at->format('Y-m-d'), ['class' => 'read-only-date']) }}`,
			expected: `<input type="date" name="created_at" value="{{ $model->created_at->format('Y-m-d') }}" class="read-only-date">`,
		},
		{
			name:     "Date with complex name",
			input:    `{{ Form::date('dates[' . $index . '][start]', old('dates[' . $index . '][start]'), ['class' => 'start-date']) }}`,
			expected: `<input type="date" name="dates[{{ $index }}][start]" value="{{ old('dates[' . $index . '][start]') }}" class="start-date">`,
		},
		{
			name:     "Date with null value",
			input:    `{{ Form::date('expiry_date', null, ['class' => 'expiry-date']) }}`,
			expected: `<input type="date" name="expiry_date" value="" class="expiry-date">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormDate(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
