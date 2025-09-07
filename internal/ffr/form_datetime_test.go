package ffr

import (
	"testing"
)

func TestFormDatetime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic datetime field",
			input:    `{{ Form::datetime('created_at') }}`,
			expected: `<input type="datetime-local" name="created_at" value="">`,
		},
		{
			name:     "Datetime with value",
			input:    `{{ Form::datetime('event_start', '2023-12-25T10:00') }}`,
			expected: `<input type="datetime-local" name="event_start" value="{{ '2023-12-25T10:00' }}">`,
		},
		{
			name:     "Datetime with old() helper",
			input:    `{{ Form::datetime('appointment', old('appointment')) }}`,
			expected: `<input type="datetime-local" name="appointment" value="{{ old('appointment') }}">`,
		},
		{
			name:     "Datetime with attributes",
			input:    `{{ Form::datetime('meeting_time', old('meeting_time'), ['class' => 'form-control', 'step' => '60']) }}`,
			expected: `<input type="datetime-local" name="meeting_time" value="{{ old('meeting_time') }}" class="form-control">`,
		},
		{
			name:     "Datetime with double exclamation marks",
			input:    `{!! Form::datetime('deadline', $task->deadline, ['class' => 'datetime-picker', 'required' => 'required']) !!}`,
			expected: `<input type="datetime-local" name="deadline" value="{{ $task->deadline }}" class="datetime-picker" required>`,
		},
		{
			name: "Multi-line datetime field",
			input: `{!! Form::datetime('event_datetime', old('event_datetime'), [
    'class' => 'form-control event-datetime',
    'id' => 'event-datetime-picker',
    'min' => '2023-01-01T00:00'
]) !!}`,
			expected: `<input type="datetime-local" name="event_datetime" value="{{ old('event_datetime') }}" class="form-control event-datetime" id="event-datetime-picker">`,
		},
		{
			name:     "Datetime with Carbon format",
			input:    `{{ Form::datetime('published_at', $post->published_at->format('Y-m-d\\TH:i'), ['class' => 'publish-time']) }}`,
			expected: `<input type="datetime-local" name="published_at" value="{{ $post->published_at->format('Y-m-d\\TH:i') }}" class="publish-time">`,
		},
		{
			name:     "Datetime with complex name",
			input:    `{{ Form::datetime('events[' . $i . '][start_time]', old('events[' . $i . '][start_time]'), ['class' => 'event-start']) }}`,
			expected: `<input type="datetime-local" name="events[{{ $i }}][start_time]" value="{{ old('events[' . $i . '][start_time]') }}" class="event-start">`,
		},
		{
			name:     "Datetime with null value",
			input:    `{{ Form::datetime('reminder_at', null, ['class' => 'reminder-datetime']) }}`,
			expected: `<input type="datetime-local" name="reminder_at" value="" class="reminder-datetime">`,
		},
		{
			name:     "Datetime with min and max constraints",
			input:    `{{ Form::datetime('booking_time', old('booking_time'), ['min' => '2023-01-01T09:00', 'max' => '2023-12-31T17:00', 'class' => 'booking-datetime']) }}`,
			expected: `<input type="datetime-local" name="booking_time" value="{{ old('booking_time') }}" class="booking-datetime">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormDatetime(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
