package ffr

import (
	"testing"
)

func TestFormTel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic tel field",
			input:    `{{ Form::tel('phone') }}`,
			expected: `<input type="tel" name="phone" value="">`,
		},
		{
			name:     "Tel with value",
			input:    `{{ Form::tel('phone', '+1234567890') }}`,
			expected: `<input type="tel" name="phone" value="{{ '+1234567890' }}">`,
		},
		{
			name:     "Tel with old() helper",
			input:    `{{ Form::tel('phone', old('phone')) }}`,
			expected: `<input type="tel" name="phone" value="{{ old('phone') }}">`,
		},
		{
			name:     "Tel with attributes",
			input:    `{{ Form::tel('phone', old('phone'), ['placeholder' => '+1 (555) 123-4567', 'class' => 'form-control']) }}`,
			expected: `<input type="tel" name="phone" value="{{ old('phone') }}" placeholder="+1 (555) 123-4567" class="form-control">`,
		},
		{
			name:     "Tel with double exclamation marks",
			input:    `{!! Form::tel('mobile', $user->mobile, ['class' => 'phone-input', 'required' => 'required']) !!}`,
			expected: `<input type="tel" name="mobile" value="{{ $user->mobile }}" class="phone-input" required>`,
		},
		{
			name: "Multi-line tel field",
			input: `{!! Form::tel('emergency_contact', old('emergency_contact'), [
    'placeholder' => 'Emergency contact number',
    'class' => 'form-control emergency-phone',
    'id' => 'emergency-contact'
]) !!}`,
			expected: `<input type="tel" name="emergency_contact" value="{{ old('emergency_contact') }}" placeholder="Emergency contact number" class="form-control emergency-phone" id="emergency-contact">`,
		},
		{
			name:     "Tel with complex name",
			input:    `{{ Form::tel('contacts[' . $i . '][phone]', old('contacts[' . $i . '][phone]'), ['class' => 'contact-phone']) }}`,
			expected: `<input type="tel" name="contacts[{{ $i }}][phone]" value="{{ old('contacts[' . $i . '][phone]') }}" class="contact-phone">`,
		},
		{
			name:     "Tel with pattern attribute",
			input:    `{{ Form::tel('phone_number', old('phone_number'), ['pattern' => '[0-9]{3}-[0-9]{3}-[0-9]{4}', 'class' => 'formatted-phone']) }}`,
			expected: `<input type="tel" name="phone_number" value="{{ old('phone_number') }}" class="formatted-phone">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormTel(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
