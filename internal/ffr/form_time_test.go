package ffr

import (
	"testing"
)

func TestFormTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic time field",
			input:    `{{ Form::time('start_time') }}`,
			expected: `<input type="time" name="start_time" value="">`,
		},
		{
			name:     "Time with value",
			input:    `{{ Form::time('start_time', '09:00') }}`,
			expected: `<input type="time" name="start_time" value="{{ '09:00' }}">`,
		},
		{
			name:     "Time with old() helper (user's original example)",
			input:    `{{ Form::time('start_at_time', old('start_at_time', '00:00'), ['placeholder' => '時間を入力', 'class' => 'form-control col-4']) }}`,
			expected: `<input type="time" name="start_at_time" value="{{ old('start_at_time', '00:00') }}" placeholder="時間を入力" class="form-control col-4">`,
		},
		{
			name:     "Time with attributes",
			input:    `{{ Form::time('meeting_time', old('meeting_time'), ['class' => 'form-control', 'step' => '300']) }}`,
			expected: `<input type="time" name="meeting_time" value="{{ old('meeting_time') }}" class="form-control">`,
		},
		{
			name:     "Time with double exclamation marks",
			input:    `{!! Form::time('end_time', $schedule->end_time, ['class' => 'time-picker', 'required' => 'required']) !!}`,
			expected: `<input type="time" name="end_time" value="{{ $schedule->end_time }}" class="time-picker" required>`,
		},
		{
			name: "Multi-line time field",
			input: `{!! Form::time('working_hours_start', old('working_hours_start'), [
    'class' => 'form-control working-time',
    'id' => 'working-hours-start',
    'step' => '900'
]) !!}`,
			expected: `<input type="time" name="working_hours_start" value="{{ old('working_hours_start') }}" class="form-control working-time" id="working-hours-start">`,
		},
		{
			name:     "Time with Carbon format",
			input:    `{{ Form::time('created_time', $log->created_at->format('H:i'), ['class' => 'log-time']) }}`,
			expected: `<input type="time" name="created_time" value="{{ $log->created_at->format('H:i') }}" class="log-time">`,
		},
		{
			name:     "Time with complex name",
			input:    `{{ Form::time('schedules[' . $day . '][start]', old('schedules[' . $day . '][start]'), ['class' => 'schedule-time']) }}`,
			expected: `<input type="time" name="schedules[{{ $day }}][start]" value="{{ old('schedules[' . $day . '][start]') }}" class="schedule-time">`,
		},
		{
			name:     "Time with min and max",
			input:    `{{ Form::time('office_hours', old('office_hours'), ['min' => '08:00', 'max' => '18:00', 'class' => 'office-time']) }}`,
			expected: `<input type="time" name="office_hours" value="{{ old('office_hours') }}" class="office-time">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormTime(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
