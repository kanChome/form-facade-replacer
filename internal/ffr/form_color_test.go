package ffr

import (
	"testing"
)

func TestFormColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic color field",
			input:    `{{ Form::color('theme_color') }}`,
			expected: `<input type="color" name="theme_color" value="">`,
		},
		{
			name:     "Color with value",
			input:    `{{ Form::color('brand_color', '#ff0000') }}`,
			expected: `<input type="color" name="brand_color" value="{{ #ff0000 }}">`,
		},
		{
			name:     "Color with old() helper",
			input:    `{{ Form::color('background_color', old('background_color')) }}`,
			expected: `<input type="color" name="background_color" value="{{ old('background_color') }}">`,
		},
		{
			name:     "Color with attributes",
			input:    `{{ Form::color('accent_color', old('accent_color', '#3498db'), ['class' => 'form-control color-picker']) }}`,
			expected: `<input type="color" name="accent_color" value="{{ old('accent_color', '#3498db') }}" class="form-control color-picker">`,
		},
		{
			name:     "Color with double exclamation marks",
			input:    `{!! Form::color('primary_color', $theme->primary_color, ['class' => 'color-input', 'required' => 'required']) !!}`,
			expected: `<input type="color" name="primary_color" value="{{ $theme->primary_color }}" class="color-input" required>`,
		},
		{
			name: "Multi-line color field",
			input: `{!! Form::color('custom_color', old('custom_color'), [
    'class' => 'form-control custom-color-picker',
    'id' => 'custom-color-input',
    'title' => 'Choose custom color'
]) !!}`,
			expected: `<input type="color" name="custom_color" value="{{ old('custom_color') }}" class="form-control custom-color-picker" id="custom-color-input">`,
		},
		{
			name:     "Color with default value",
			input:    `{{ Form::color('ui_color', $settings->ui_color ?? '#ffffff', ['class' => 'ui-color-picker']) }}`,
			expected: `<input type="color" name="ui_color" value="{{ $settings->ui_color ?? '#ffffff' }}" class="ui-color-picker">`,
		},
		{
			name:     "Color with complex name",
			input:    `{{ Form::color('colors[' . $section . '][primary]', old('colors[' . $section . '][primary]'), ['class' => 'section-color']) }}`,
			expected: `<input type="color" name="colors[{{ $section }}][primary]" value="{{ old('colors[' . $section . '][primary]') }}" class="section-color">`,
		},
		{
			name:     "Color with RGB value",
			input:    `{{ Form::color('text_color', '#333333', ['class' => 'text-color-picker']) }}`,
			expected: `<input type="color" name="text_color" value="{{ #333333 }}" class="text-color-picker">`,
		},
		{
			name:     "Color with list attribute",
			input:    `{{ Form::color('palette_color', old('palette_color'), ['list' => 'color_presets', 'class' => 'palette-picker']) }}`,
			expected: `<input type="color" name="palette_color" value="{{ old('palette_color') }}" class="palette-picker">`,
		},
		{
			name:     "Color with onchange event",
			input:    `{{ Form::color('preview_color', '#000000', ['onchange' => 'updatePreview(this.value)', 'class' => 'preview-color']) }}`,
			expected: `<input type="color" name="preview_color" value="{{ #000000 }}" class="preview-color">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormColor(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
