package ffr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Complete user registration form",
			input: `<!DOCTYPE html>
<html>
<head>
    <title>User Registration</title>
</head>
<body>
    <div class="container">
        {!! Form::open(['route' => 'users.store', 'method' => 'POST', 'class' => 'registration-form']) !!}
            <div class="form-group">
                {!! Form::label('name', 'Full Name') !!}
                {!! Form::text('name', old('name'), ['class' => 'form-control', 'placeholder' => 'Enter your name']) !!}
            </div>
            <div class="form-group">
                {!! Form::label('email', 'Email Address') !!}
                {!! Form::text('email', old('email'), ['class' => 'form-control', 'placeholder' => 'Enter your email']) !!}
            </div>
            <div class="form-group">
                {!! Form::label('age', 'Age') !!}
                {!! Form::number('age', old('age'), ['class' => 'form-control', 'min' => 18]) !!}
            </div>
            <div class="form-group">
                {!! Form::label('bio', 'Biography') !!}
                {!! Form::textarea('bio', old('bio'), ['class' => 'form-control', 'rows' => 4]) !!}
            </div>
            <div class="form-group">
                {!! Form::label('country', 'Country') !!}
                {!! Form::select('country', ['jp' => 'Japan', 'us' => 'USA'], old('country'), ['class' => 'form-control']) !!}
            </div>
            <div class="form-group">
                {!! Form::checkbox('newsletter', 1, old('newsletter'), ['class' => 'form-check-input']) !!}
                {!! Form::label('newsletter', 'Subscribe to newsletter') !!}
            </div>
            {!! Form::hidden('source', 'web') !!}
            <div class="form-actions">
                {!! Form::submit('Register', ['class' => 'btn btn-primary']) !!}
                {!! Form::button('Cancel', ['type' => 'button', 'class' => 'btn btn-secondary']) !!}
            </div>
        {!! Form::close() !!}
    </div>
</body>
</html>`,
			expected: `<!DOCTYPE html>
<html>
<head>
    <title>User Registration</title>
</head>
<body>
    <div class="container">
        <form action="{{ route('users.store') }}" method="POST" class="registration-form">
{{ csrf_field() }}
            <div class="form-group">
                <label for="name">{!! 'Full Name' !!}</label>
                <input type="text" name="name" value="{{ old('name') }}" placeholder="Enter your name" class="form-control">
            </div>
            <div class="form-group">
                <label for="email">{!! 'Email Address' !!}</label>
                <input type="text" name="email" value="{{ old('email') }}" placeholder="Enter your email" class="form-control">
            </div>
            <div class="form-group">
                <label for="age">{!! 'Age' !!}</label>
                <input type="number" name="age" value="{{ old('age') }}" class="form-control" min="18">
            </div>
            <div class="form-group">
                <label for="bio">{!! 'Biography' !!}</label>
                <textarea name="bio" rows="4" class="form-control">{{ old('bio') }}</textarea>
            </div>
            <div class="form-group">
                <label for="country">{!! 'Country' !!}</label>
                <select name="country" class="form-control">
@foreach(['jp' => 'Japan', 'us' => 'USA'] as $key => $value)
<option value="{{ $key }}" @if($key == old('country')) selected @endif>{{ $value }}</option>
@endforeach
</select>
            </div>
            <div class="form-group">
                <input type="checkbox" name="newsletter" value="{{ 1 }}" @if(old('newsletter')) checked @endif class="form-check-input">
                <label for="newsletter">{!! 'Subscribe to newsletter' !!}</label>
            </div>
            <input type="hidden" name="source" value="{{ 'web' }}">
            <div class="form-actions">
                <button type="submit" class="btn btn-primary">Register</button>
                <button type="button" class="btn btn-secondary">{!! 'Cancel' !!}</button>
            </div>
        </form>
    </div>
</body>
</html>`,
		},
		{
			name: "Mixed form elements with various patterns",
			input: `{{ Form::open(['url' => '/contact', 'method' => 'POST']) }}
    {{ Form::text('name', $user->name) }}
    {!! Form::hidden('user_id', $user->id) !!}
    {{ Form::textarea('message', old('message'), ['placeholder' => 'Your message']) }}
    {!! Form::checkbox('urgent', 1, false, ['id' => 'urgent-check']) !!}
    {{ Form::select('department', $departments, null, ['class' => 'form-select']) }}
{{ Form::close() }}`,
			expected: `<form action="'/contact'" method="POST">
{{ csrf_field() }}
    <input type="text" name="name" value="{{ $user->name }}">
    <input type="hidden" name="user_id" value="{{ $user->id }}">
    <textarea name="message" placeholder="Your message">{{ old('message') }}</textarea>
    <input type="checkbox" name="urgent" value="{{ 1 }}" @if(false) checked @endif id="urgent-check">
    <select name="department" class="form-select">
@foreach($departments as $key => $value)
<option value="{{ $key }}" @if($key == null) selected @endif>{{ $value }}</option>
@endforeach
</select>
</form>`,
		},
		{
			name: "Form with array fields and complex attributes",
			input: `{!! Form::open(['route' => ['survey.store'], 'method' => 'POST']) !!}
    {!! Form::checkbox('interests[]', 'tech', old('interests'), ['class' => 'interest-check']) !!}
    {!! Form::checkbox('interests[]', 'sports', old('interests'), ['class' => 'interest-check']) !!}
    {!! Form::text('skills[]', '', ['class' => 'skill-input']) !!}
    {!! Form::hidden('responses[0][question_id]', 1) !!}
    {!! Form::textarea('responses[0][answer]', '', ['rows' => 3]) !!}
{!! Form::close() !!}`,
			expected: `<form action="{{ route('survey.store') }}" method="POST">
{{ csrf_field() }}
    <input type="checkbox" name="interests[]" value="{{ tech }}" @if(in_array(tech, (array)old('interests'))) checked @endif class="interest-check">
    <input type="checkbox" name="interests[]" value="{{ sports }}" @if(in_array(sports, (array)old('interests'))) checked @endif class="interest-check">
    <input type="text" name="skills[]" value="" class="skill-input">
    <input type="hidden" name="responses[0][question_id]" value="{{ is_array(1) ? implode(',', 1) : 1 }}">
    <textarea name="responses[0][answer]" rows="3"></textarea>
</form>`,
		},
		{
			name: "Form with file upload and mixed elements",
			input: `{!! Form::open(['route' => 'profile.update', 'method' => 'POST', 'files' => true]) !!}
    {!! Form::label('name', 'Full Name') !!}
    {!! Form::text('name', old('name'), ['class' => 'form-control']) !!}
    {!! Form::label('avatar', 'Profile Picture') !!}
    {!! Form::file('avatar', ['accept' => 'image/*', 'class' => 'form-control']) !!}
    {!! Form::label('documents', 'Documents') !!}
    {!! Form::file('documents[]', ['multiple' => true, 'accept' => '.pdf,.doc,.docx']) !!}
    {!! Form::textarea('bio', old('bio'), ['rows' => 4, 'placeholder' => 'Tell us about yourself']) !!}
    {!! Form::submit('Update Profile', ['class' => 'btn btn-primary']) !!}
{!! Form::close() !!}`,
			expected: `<form action="{{ route('profile.update') }}" method="POST">
{{ csrf_field() }}
    <label for="name">{!! 'Full Name' !!}</label>
    <input type="text" name="name" value="{{ old('name') }}" class="form-control">
    <label for="avatar">{!! 'Profile Picture' !!}</label>
    <input type="file" name="avatar" accept="image/*" class="form-control">
    <label for="documents">{!! 'Documents' !!}</label>
    <input type="file" name="documents[]" accept=".pdf,.doc,.docx" multiple>
    <textarea name="bio" rows="4" placeholder="Tell us about yourself">{{ old('bio') }}</textarea>
    <button type="submit" class="btn btn-primary">Update Profile</button>
</form>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormPatternsString(tt.input)
			if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
				t.Errorf("Integration test failed.\nGot:\n%s\nWant:\n%s", result, tt.expected)
			}
		})
	}
}

// replaceFormPatternsString テスト用のヘルパー関数（文字列処理版）
func replaceFormPatternsString(text string) string {
	text = replaceFormOpen(text)
	text = replaceFormClose(text)
	text = replaceFormHidden(text)
	text = replaceFormButton(text)
	text = replaceFormTextarea(text)
	text = replaceFormLabel(text)
	text = replaceFormText(text)
	text = replaceFormNumber(text)
	text = replaceFormSelect(text)
	text = replaceFormCheckbox(text)
	text = replaceFormSubmit(text)
	text = replaceFormFile(text)
	return text
}

func TestFileProcessingIntegration(t *testing.T) {
	// 一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用のBladeファイルを作成
	testContent := `@extends('layouts.app')

@section('content')
<div class="container">
    {!! Form::open(['route' => 'posts.store', 'method' => 'POST']) !!}
        <div class="form-group">
            {!! Form::label('title', 'Title') !!}
            {!! Form::text('title', old('title'), ['class' => 'form-control']) !!}
        </div>
        <div class="form-group">
            {!! Form::label('content', 'Content') !!}
            {!! Form::textarea('content', old('content'), ['class' => 'form-control', 'rows' => 10]) !!}
        </div>
        {!! Form::submit('Save Post', ['class' => 'btn btn-primary']) !!}
    {!! Form::close() !!}
</div>
@endsection`

	expectedContent := `@extends('layouts.app')

@section('content')
<div class="container">
    <form action="{{ route('posts.store') }}" method="POST">
{{ csrf_field() }}
        <div class="form-group">
            <label for="title">{!! 'Title' !!}</label>
            <input type="text" name="title" value="{{ old('title') }}" class="form-control">
        </div>
        <div class="form-group">
            <label for="content">{!! 'Content' !!}</label>
            <textarea name="content" rows="10" class="form-control">{{ old('content') }}</textarea>
        </div>
        <button type="submit" class="btn btn-primary">Save Post</button>
    </form>
</div>
@endsection`

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.blade.php")
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// ファイル処理を実行
	err = replaceFormPatterns(testFile)
	if err != nil {
		t.Fatalf("Failed to process file: %v", err)
	}

	// 処理後のファイル内容を読み込み
	result, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read processed file: %v", err)
	}

	// 結果を検証
	if strings.TrimSpace(string(result)) != strings.TrimSpace(expectedContent) {
		t.Errorf("File processing integration test failed.\nGot:\n%s\nWant:\n%s", string(result), expectedContent)
	}
}

func TestDirectoryProcessingIntegration(t *testing.T) {
	// 一時ディレクトリを作成
	tempDir := t.TempDir()

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "components")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// 複数のBladeファイルを作成
	files := map[string]string{
		"form1.blade.php": `{!! Form::open(['url' => '/test']) !!}
{!! Form::text('name', '') !!}
{!! Form::close() !!}`,
		"form2.blade.php": `{{ Form::open(['method' => 'POST']) }}
{{ Form::number('age', 25) }}
{{ Form::close() }}`,
		"components/input.blade.php": `{!! Form::label('field', 'Label') !!}
{!! Form::text('field', $value) !!}`,
	}

	expected := map[string]string{
		"form1.blade.php": `<form action="'/test'" method="GET">
<input type="text" name="name" value="">
</form>`,
		"form2.blade.php": `<form action="" method="POST">
{{ csrf_field() }}
<input type="number" name="age" value="{{ 25 }}">
</form>`,
		"components/input.blade.php": `<label for="field">{!! 'Label' !!}</label>
<input type="text" name="field" value="{{ $value }}">`,
	}

	// テストファイルを作成
	for filename, content := range files {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// 設定を作成
	config := &ReplacementConfig{
		TargetPath:     tempDir,
		IsFile:         false,
		ProcessedFiles: make([]string, 0),
		FileCount:      0,
	}

	// ディレクトリ処理を実行
	err = processBladeFiles(config)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	// 処理されたファイル数を検証
	if config.FileCount != len(files) {
		t.Errorf("Expected %d processed files, got %d", len(files), config.FileCount)
	}

	// 各ファイルの内容を検証
	for filename, expectedContent := range expected {
		filePath := filepath.Join(tempDir, filename)
		result, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read processed file %s: %v", filename, err)
		}

		if strings.TrimSpace(string(result)) != strings.TrimSpace(expectedContent) {
			t.Errorf("File %s processing failed.\nGot:\n%s\nWant:\n%s", filename, string(result), expectedContent)
		}
	}
}
