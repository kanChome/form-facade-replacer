package main

import (
	"testing"
)

func TestFormFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic file field",
			input:    "{{ Form::file('image') }}",
			expected: `<input type="file" name="image">`,
		},
		{
			name:     "File field with accept attribute",
			input:    "{{ Form::file('image', ['accept' => '.jpg,.jpeg,.png']) }}",
			expected: `<input type="file" name="image" accept=".jpg,.jpeg,.png">`,
		},
		{
			name:     "File field with multiple attribute",
			input:    "{{ Form::file('images', ['multiple' => true]) }}",
			expected: `<input type="file" name="images" multiple>`,
		},
		{
			name:     "File field with onchange attribute",
			input:    `{{ Form::file('avatar', ['onchange' => 'previewImage(this)']) }}`,
			expected: `<input type="file" name="avatar" onchange="previewImage(this)">`,
		},
		{
			name:     "File field with complex onchange (user example)",
			input:    `{{ Form::file('image', ['accept' => '.jpg,.jpeg,.png', 'onchange' => 'previewMainImage(this, "event_area")']) }}`,
			expected: `<input type="file" name="image" accept=".jpg,.jpeg,.png" onchange="previewMainImage(this, 'event_area')">`,
		},
		{
			name:     "File field with id and class",
			input:    "{{ Form::file('document', ['id' => 'doc-upload', 'class' => 'file-input']) }}",
			expected: `<input type="file" name="document" class="file-input" id="doc-upload">`,
		},
		{
			name:     "File field with array name",
			input:    "{{ Form::file('documents[]', ['multiple' => true]) }}",
			expected: `<input type="file" name="documents[]" multiple>`,
		},
		{
			name:     "File field with capture attribute",
			input:    "{{ Form::file('photo', ['accept' => 'image/*', 'capture' => 'camera']) }}",
			expected: `<input type="file" name="photo" accept="image/*" capture="camera">`,
		},
		{
			name:     "File field with all common attributes",
			input:    "{{ Form::file('upload', ['accept' => '.pdf,.doc', 'multiple' => true, 'class' => 'file-upload', 'id' => 'file-input']) }}",
			expected: `<input type="file" name="upload" accept=".pdf,.doc" class="file-upload" id="file-input" multiple>`,
		},
		{
			name:     "File field with empty attributes",
			input:    "{{ Form::file('simple', []) }}",
			expected: `<input type="file" name="simple">`,
		},
		{
			name:     "File field with double exclamation marks",
			input:    "{!! Form::file('media', ['accept' => 'video/*,audio/*']) !!}",
			expected: `<input type="file" name="media" accept="video/*,audio/*">`,
		},
		{
			name:     "File field with boolean false",
			input:    "{{ Form::file('optional', ['multiple' => false]) }}",
			expected: `<input type="file" name="optional">`,
		},
		{
			name:     "File field with onchange single string argument",
			input:    `{{ Form::file('photo', ['onchange' => 'uploadFile(this, "temp/photos")']) }}`,
			expected: `<input type="file" name="photo" onchange="uploadFile(this, 'temp/photos')">`,
		},
		{
			name:     "File field with onchange multiple string arguments",
			input:    `{{ Form::file('docs', ['onchange' => 'processFile(this, "uploads", "documents")']) }}`,
			expected: `<input type="file" name="docs" onchange="processFile(this, "uploads", "documents")">`,
		},
		{
			name:     "File field with onchange mixed arguments",
			input:    `{{ Form::file('img', ['onchange' => 'resizeImage(this, "thumb", 300, "jpg")']) }}`,
			expected: `<input type="file" name="img" onchange="resizeImage(this, "thumb", 300, "jpg")">`,
		},
		{
			name:     "File field with onclick double quotes",
			input:    `{{ Form::file('data', ['onclick' => 'showDialog(this, "Upload File", "Select a file to upload")']) }}`,
			expected: `<input type="file" name="data" onclick="showDialog(this, "Upload File", "Select a file to upload")">`,
		},
		{
			name:     "File field with onchange escaped quotes",
			input:    `{{ Form::file('file', ['onchange' => 'alert(this, "Say \"Hello\"")']) }}`,
			expected: `<input type="file" name="file" onchange="alert(this, "Say \"Hello\"")">`,
		},
		{
			name:     "File field with onchange nested quotes in JSON",
			input:    `{{ Form::file('config', ['onchange' => 'parseJSON(this, "{\"key\": \"value\"}")']) }}`,
			expected: `<input type="file" name="config" onchange="parseJSON(this, "{\"key\": \"value\"}")">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormFile(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
