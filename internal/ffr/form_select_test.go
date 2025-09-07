package ffr

import (
	"testing"
)

func TestFormSelect(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Basic select",
			input: "{{ Form::select('country', ['jp' => 'Japan', 'us' => 'United States', 'uk' => 'United Kingdom']) }}",
			expected: `<select name="country">
@foreach(['jp' => 'Japan', 'us' => 'United States', 'uk' => 'United Kingdom'] as $key => $value)
<option value="{{ $key }}" @if($key == ) selected @endif>{{ $value }}</option>
@endforeach
</select>`,
		},
		{
			name:  "Select with selected value",
			input: "{{ Form::select('country', ['jp' => 'Japan', 'us' => 'United States', 'uk' => 'United Kingdom'], 'us') }}",
			expected: `<select name="country">
@foreach(['jp' => 'Japan', 'us' => 'United States', 'uk' => 'United Kingdom'] as $key => $value)
<option value="{{ $key }}" @if($key == 'us') selected @endif>{{ $value }}</option>
@endforeach
</select>`,
		},
		{
			name:  "Select with attributes",
			input: "{{ Form::select('country', ['jp' => 'Japan', 'us' => 'United States'], 'jp', ['class' => 'form-control', 'multiple' => 'multiple']) }}",
			expected: `<select name="country" class="form-control">
@foreach(['jp' => 'Japan', 'us' => 'United States'] as $key => $value)
<option value="{{ $key }}" @if($key == 'jp') selected @endif>{{ $value }}</option>
@endforeach
</select>`,
		},
		{
			name:  "Select with empty options",
			input: "{{ Form::select('empty', []) }}",
			expected: `<select name="empty">
@foreach([] as $key => $value)
<option value="{{ $key }}" @if($key == ) selected @endif>{{ $value }}</option>
@endforeach
</select>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormSelect(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
