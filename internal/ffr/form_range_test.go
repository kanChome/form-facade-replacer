package ffr

import (
	"testing"
)

func TestFormRange(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic range field",
			input:    `{{ Form::range('volume') }}`,
			expected: `<input type="range" name="volume" value="">`,
		},
		{
			name:     "Range with value",
			input:    `{{ Form::range('volume', '50') }}`,
			expected: `<input type="range" name="volume" value="{{ 50 }}">`,
		},
		{
			name:     "Range with old() helper",
			input:    `{{ Form::range('brightness', old('brightness')) }}`,
			expected: `<input type="range" name="brightness" value="{{ old('brightness') }}">`,
		},
		{
			name:     "Range with attributes",
			input:    `{{ Form::range('volume', old('volume', '50'), ['min' => '0', 'max' => '100', 'class' => 'form-range']) }}`,
			expected: `<input type="range" name="volume" value="{{ old('volume', '50') }}" class="form-range">`,
		},
		{
			name:     "Range with double exclamation marks",
			input:    `{!! Form::range('temperature', $settings->temperature, ['min' => '16', 'max' => '30', 'step' => '0.5', 'class' => 'temp-slider']) !!}`,
			expected: `<input type="range" name="temperature" value="{{ $settings->temperature }}" class="temp-slider">`,
		},
		{
			name: "Multi-line range field",
			input: `{!! Form::range('opacity', old('opacity', '100'), [
    'min' => '0',
    'max' => '100',
    'step' => '1',
    'class' => 'form-control opacity-slider',
    'id' => 'opacity-range'
]) !!}`,
			expected: `<input type="range" name="opacity" value="{{ old('opacity', '100') }}" class="form-control opacity-slider" id="opacity-range">`,
		},
		{
			name:     "Range with percentage value",
			input:    `{{ Form::range('progress', $task->completion_percentage, ['min' => '0', 'max' => '100', 'class' => 'progress-bar']) }}`,
			expected: `<input type="range" name="progress" value="{{ $task->completion_percentage }}" class="progress-bar">`,
		},
		{
			name:     "Range with complex name",
			input:    `{{ Form::range('settings[' . $category . '][value]', old('settings[' . $category . '][value]'), ['min' => '1', 'max' => '10', 'class' => 'setting-range']) }}`,
			expected: `<input type="range" name="settings[{{ $category }}][value]" value="{{ old('settings[' . $category . '][value]') }}" class="setting-range">`,
		},
		{
			name:     "Range with decimal step",
			input:    `{{ Form::range('rating', old('rating'), ['min' => '0', 'max' => '5', 'step' => '0.1', 'class' => 'rating-slider']) }}`,
			expected: `<input type="range" name="rating" value="{{ old('rating') }}" class="rating-slider">`,
		},
		{
			name:     "Range with data attributes",
			input:    `{{ Form::range('zoom', '1', ['min' => '0.5', 'max' => '3', 'step' => '0.1', 'data-default' => '1', 'class' => 'zoom-control']) }}`,
			expected: `<input type="range" name="zoom" value="{{ 1 }}" class="zoom-control">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormRange(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
