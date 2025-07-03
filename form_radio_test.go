package main

import (
	"testing"
)

func TestFormRadio(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "User's specific example with onchange and style",
			input:    `{!! Form::radio('body_type', $bodyTypeHtml, old('body_type'), ['id' => 'body-type-html', 'style' => 'transform: scale(1.2); margin-right: 8px;', 'onchange' => 'onClickBodyIsHtml();']) !!}`,
			expected: `<input type="radio" name="body_type" value="{{ $bodyTypeHtml }}" @if(old('body_type')) checked @endif id="body-type-html" style="transform: scale(1.2); margin-right: 8px;" onchange="onClickBodyIsHtml();">`,
		},
		{
			name:     "Basic radio button",
			input:    `{{ Form::radio('gender', 'male') }}`,
			expected: `<input type="radio" name="gender" value="{{ 'male' }}">`,
		},
		{
			name:     "Radio with variable value",
			input:    `{{ Form::radio('status', $value) }}`,
			expected: `<input type="radio" name="status" value="{{ $value }}">`,
		},
		{
			name:     "Radio with checked condition",
			input:    `{{ Form::radio('gender', 'male', $user->gender == 'male') }}`,
			expected: `<input type="radio" name="gender" value="{{ 'male' }}" @if($user->gender == 'male') checked @endif>`,
		},
		{
			name:     "Radio with old() helper checked",
			input:    `{{ Form::radio('option', 'yes', old('option') == 'yes') }}`,
			expected: `<input type="radio" name="option" value="{{ 'yes' }}" @if(old('option') == 'yes') checked @endif>`,
		},
		{
			name:     "Radio with simple old() checked",
			input:    `{{ Form::radio('enabled', '1', old('enabled')) }}`,
			expected: `<input type="radio" name="enabled" value="{{ '1' }}" @if(old('enabled')) checked @endif>`,
		},
		{
			name:     "Radio with null checked (not checked)",
			input:    `{{ Form::radio('disabled_option', 'no', null) }}`,
			expected: `<input type="radio" name="disabled_option" value="{{ 'no' }}">`,
		},
		{
			name:     "Radio with false checked (not checked)",
			input:    `{{ Form::radio('inactive', '0', false) }}`,
			expected: `<input type="radio" name="inactive" value="{{ '0' }}">`,
		},
		{
			name:     "Radio with class attribute",
			input:    `{{ Form::radio('type', 'premium', old('type'), ['class' => 'form-check-input']) }}`,
			expected: `<input type="radio" name="type" value="{{ 'premium' }}" @if(old('type')) checked @endif class="form-check-input">`,
		},
		{
			name:     "Radio with id and class",
			input:    `{{ Form::radio('plan', 'basic', $user->plan == 'basic', ['id' => 'plan-basic', 'class' => 'form-radio']) }}`,
			expected: `<input type="radio" name="plan" value="{{ 'basic' }}" @if($user->plan == 'basic') checked @endif id="plan-basic" class="form-radio">`,
		},
		{
			name:     "Radio with style attribute",
			input:    `{{ Form::radio('color', 'red', false, ['style' => 'color: red;']) }}`,
			expected: `<input type="radio" name="color" value="{{ 'red' }}" style="color: red;">`,
		},
		{
			name:     "Radio with onchange event",
			input:    `{{ Form::radio('category', 'tech', old('category'), ['onchange' => 'updateCategory();']) }}`,
			expected: `<input type="radio" name="category" value="{{ 'tech' }}" @if(old('category')) checked @endif onchange="updateCategory();">`,
		},
		{
			name:     "Radio with disabled attribute",
			input:    `{{ Form::radio('readonly', 'value', false, ['disabled' => 'disabled']) }}`,
			expected: `<input type="radio" name="readonly" value="{{ 'value' }}" disabled>`,
		},
		{
			name:     "Radio with all attributes",
			input:    `{{ Form::radio('field', 'value', true, ['id' => 'field-id', 'class' => 'radio-input', 'style' => 'margin: 5px;', 'onchange' => 'handleChange();', 'disabled' => '']) }}`,
			expected: `<input type="radio" name="field" value="{{ 'value' }}" @if(true) checked @endif id="field-id" class="radio-input" style="margin: 5px;" onchange="handleChange();" disabled>`,
		},
		{
			name:     "Radio with double exclamation marks",
			input:    `{!! Form::radio('content', $content->value, old('content'), ['class' => 'content-radio']) !!}`,
			expected: `<input type="radio" name="content" value="{{ $content->value }}" @if(old('content')) checked @endif class="content-radio">`,
		},
		{
			name:     "Radio with PHP string concatenation in name",
			input:    `{{ Form::radio('items[' . $i . '][selected]', '1', old('items[' . $i . '][selected]'), ['class' => 'item-radio']) }}`,
			expected: `<input type="radio" name="items[{{ $i }}][selected]" value="{{ '1' }}" @if(old('items[' . $i . '][selected]')) checked @endif class="item-radio">`,
		},
		{
			name:     "Radio with complex PHP string concatenation",
			input:    `{{ Form::radio('data[' . $row['id'] . '][type]', $types[$index], $selected[$row['id']] ?? false, ['id' => 'type-' . $row['id']]) }}`,
			expected: `<input type="radio" name="data[{{ $row['id'] }}][type]" value="{{ $types[$index] }}" @if($selected[$row['id']] ?? false) checked @endif id="type-">`,
		},
		{
			name: "Multi-line radio button",
			input: `{!! Form::radio('notification_type', 'email', old('notification_type') == 'email', [
    'id' => 'notification-email',
    'class' => 'form-check-input',
    'onchange' => 'toggleEmailSettings();'
]) !!}`,
			expected: `<input type="radio" name="notification_type" value="{{ 'email' }}" @if(old('notification_type') == 'email') checked @endif id="notification-email" class="form-check-input" onchange="toggleEmailSettings();">`,
		},
		{
			name: "Multi-line radio with complex attributes",
			input: `{{ Form::radio('theme', 'dark', $user->preferences['theme'] == 'dark', [
    'id' => 'theme-dark',
    'class' => 'theme-selector custom-radio',
    'style' => 'transform: scale(1.5); margin: 10px;',
    'onchange' => 'applyTheme("dark");'
]) }}`,
			expected: `<input type="radio" name="theme" value="{{ 'dark' }}" @if($user->preferences['theme'] == 'dark') checked @endif id="theme-dark" class="theme-selector custom-radio" style="transform: scale(1.5); margin: 10px;" onchange="applyTheme("dark");">`,
		},
		{
			name:     "Radio with numeric value",
			input:    `{{ Form::radio('priority', 1, old('priority') == 1) }}`,
			expected: `<input type="radio" name="priority" value="{{ 1 }}" @if(old('priority') == 1) checked @endif>`,
		},
		{
			name:     "Radio with boolean value",
			input:    `{{ Form::radio('active', true, $model->active) }}`,
			expected: `<input type="radio" name="active" value="{{ true }}" @if($model->active) checked @endif>`,
		},
		{
			name:     "Radio with complex ternary checked condition",
			input:    `{{ Form::radio('visibility', 'public', isset($post->visibility) ? $post->visibility == 'public' : true, ['class' => 'visibility-radio']) }}`,
			expected: `<input type="radio" name="visibility" value="{{ 'public' }}" @if(isset($post->visibility) ? $post->visibility == 'public' : true) checked @endif class="visibility-radio">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormRadio(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestProcessFormRadio(t *testing.T) {
	tests := []struct {
		name     string
		params   []string
		expected string
	}{
		{
			name:     "User's example parameters",
			params:   []string{"'body_type'", "$bodyTypeHtml", "old('body_type')", "['id' => 'body-type-html', 'style' => 'transform: scale(1.2); margin-right: 8px;', 'onchange' => 'onClickBodyIsHtml();']"},
			expected: `<input type="radio" name="body_type" value="{{ $bodyTypeHtml }}" @if(old('body_type')) checked @endif id="body-type-html" style="transform: scale(1.2); margin-right: 8px;" onchange="onClickBodyIsHtml();">`,
		},
		{
			name:     "Basic radio parameters",
			params:   []string{"'gender'", "'male'"},
			expected: `<input type="radio" name="gender" value="{{ 'male' }}">`,
		},
		{
			name:     "Radio with checked condition",
			params:   []string{"'status'", "'active'", "$user->status == 'active'"},
			expected: `<input type="radio" name="status" value="{{ 'active' }}" @if($user->status == 'active') checked @endif>`,
		},
		{
			name:     "Radio with attributes but no checked",
			params:   []string{"'type'", "'premium'", "", "['class' => 'premium-radio', 'id' => 'type-premium']"},
			expected: `<input type="radio" name="type" value="{{ 'premium' }}" id="type-premium" class="premium-radio">`,
		},
		{
			name:     "Radio with null checked (should not show checked)",
			params:   []string{"'option'", "'yes'", "null", "['class' => 'option-radio']"},
			expected: `<input type="radio" name="option" value="{{ 'yes' }}" class="option-radio">`,
		},
		{
			name:     "Radio with false checked (should not show checked)",
			params:   []string{"'enabled'", "'1'", "false", "['id' => 'enabled-radio']"},
			expected: `<input type="radio" name="enabled" value="{{ '1' }}" id="enabled-radio">`,
		},
		{
			name:     "Radio with all supported attributes",
			params:   []string{"'field'", "'value'", "true", "['id' => 'field-id', 'class' => 'field-class', 'style' => 'color: blue;', 'onchange' => 'doSomething();', 'disabled' => 'disabled']"},
			expected: `<input type="radio" name="field" value="{{ 'value' }}" @if(true) checked @endif id="field-id" class="field-class" style="color: blue;" onchange="doSomething();" disabled>`,
		},
		{
			name:     "Radio with insufficient parameters",
			params:   []string{"'name'"},
			expected: ``,
		},
		{
			name:     "Radio with empty parameters",
			params:   []string{},
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormRadio(tt.params)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
