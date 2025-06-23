package main

import (
	"strings"
	"testing"
)

func TestReplaceFormOpen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic GET form with route",
			input:    `{!! Form::open(['route' => 'user.index']) !!}`,
			expected: `<form action="{{ route('user.index') }}" method="GET">`,
		},
		{
			name:     "GET from with none route",
			input:    `{!! Form::open(['method' => 'get']) !!}`,
			expected: `<form action="" method="get">`,
		},
		{
			name:  "POST form with route and CSRF",
			input: `{!! Form::open(['route' => 'user.store', 'method' => 'POST']) !!}`,
			expected: `<form action="{{ route('user.store') }}" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:  "Form with class and id attributes",
			input: `{!! Form::open(['route' => 'user.store', 'method' => 'POST', 'class' => 'user-form', 'id' => 'create-user']) !!}`,
			expected: `<form action="{{ route('user.store') }}" method="POST" class="user-form" id="create-user">
{{ csrf_field() }}`,
		},
		{
			name:  "Form with URL instead of route",
			input: `{!! Form::open(['url' => '/users', 'method' => 'POST']) !!}`,
			expected: `<form action="'/users'" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:     "Form with target attribute",
			input:    `{!! Form::open(['route' => 'user.index', 'target' => '_blank']) !!}`,
			expected: `<form action="{{ route('user.index') }}" method="GET" target="_blank">`,
		},
		{
			name:  "Route with array parameters",
			input: `{!! Form::open(['route' => ['user.update', ['id' => $user->id]], 'method' => 'PUT']) !!}`,
			expected: `<form action="{{ route('user.update', ['id' => $user->id]) }}" method="PUT">
{{ csrf_field() }}`,
		},
		{
			name:  "Blade syntax with double curly braces",
			input: `{{ Form::open(['route' => 'user.store', 'method' => 'POST']) }}`,
			expected: `<form action="{{ route('user.store') }}" method="POST">
{{ csrf_field() }}`,
		},
		{
			name: "Multi-line form definition",
			input: `{!! Form::open([
    'route' => 'user.store',
    'method' => 'POST',
    'class' => 'user-form'
]) !!}`,
			expected: `<form action="{{ route('user.store') }}" method="POST" class="user-form">
{{ csrf_field() }}`,
		},
		{
			name:  "Multiple Form::open in same text",
			input: `{!! Form::open(['route' => 'user.create']) !!} some content {!! Form::open(['route' => 'post.create', 'method' => 'POST']) !!}`,
			expected: `<form action="{{ route('user.create') }}" method="GET"> some content <form action="{{ route('post.create') }}" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:     "Empty form (no parameters)",
			input:    `{!! Form::open([]) !!}`,
			expected: `<form action="" method="GET">`,
		},
		{
			name:     "Double curly braces with extra spaces",
			input:    `{{  Form::open(['route' => 'user.index'])  }}`,
			expected: `<form action="{{ route('user.index') }}" method="GET">`,
		},
		{
			name:  "Double curly braces with complex spacing",
			input: `{{ Form::open( [ 'route' => 'user.store', 'method' => 'POST' ] ) }}`,
			expected: `<form action="{{ route('user.store') }}" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:  "Double curly braces with URL",
			input: `{{ Form::open(['url' => '/test', 'method' => 'POST']) }}`,
			expected: `<form action="'/test'" method="POST">
{{ csrf_field() }}`,
		},
		{
			name: "Double curly braces multi-line",
			input: `{{ Form::open([
    'route' => 'user.update',
    'method' => 'PUT'
]) }}`,
			expected: `<form action="{{ route('user.update') }}" method="PUT">
{{ csrf_field() }}`,
		},
		{
			name:  "Mixed brackets in same text",
			input: `{!! Form::open(['route' => 'user.create']) !!} and {{ Form::open(['route' => 'user.edit', 'method' => 'PUT']) }}`,
			expected: `<form action="{{ route('user.create') }}" method="GET"> and <form action="{{ route('user.edit') }}" method="PUT">
{{ csrf_field() }}`,
		},
		{
			name:  "Double curly braces with array route parameters",
			input: `{{ Form::open(['route' => ['user.update', ['id' => $user->id]], 'method' => 'PATCH']) }}`,
			expected: `<form action="{{ route('user.update', ['id' => $user->id]) }}" method="PATCH">
{{ csrf_field() }}`,
		},
		{
			name:  "Double curly braces with all attributes",
			input: `{{ Form::open(['route' => 'user.store', 'method' => 'POST', 'class' => 'form-control', 'id' => 'main-form', 'target' => '_self']) }}`,
			expected: `<form action="{{ route('user.store') }}" method="POST" class="form-control" id="main-form" target="_self">
{{ csrf_field() }}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormOpen(tt.input)
			if result != tt.expected {
				t.Errorf("replaceFormOpen() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProcessFormOpen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic route",
			input:    `'route' => 'user.index'`,
			expected: `<form action="{{ route('user.index') }}" method="GET">`,
		},
		{
			name:  "Route with POST method",
			input: `'route' => 'user.store', 'method' => 'POST'`,
			expected: `<form action="{{ route('user.store') }}" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:     "URL with GET method",
			input:    `'url' => '/users'`,
			expected: `<form action="'/users'" method="GET">`,
		},
		{
			name:  "URL with POST method",
			input: `'url' => '/users', 'method' => 'POST'`,
			expected: `<form action="'/users'" method="POST">
{{ csrf_field() }}`,
		},
		{
			name:  "Route with class attribute",
			input: `'route' => 'user.store', 'method' => 'POST', 'class' => 'user-form'`,
			expected: `<form action="{{ route('user.store') }}" method="POST" class="user-form">
{{ csrf_field() }}`,
		},
		{
			name:  "Route with id attribute",
			input: `'route' => 'user.store', 'method' => 'POST', 'id' => 'create-user'`,
			expected: `<form action="{{ route('user.store') }}" method="POST" id="create-user">
{{ csrf_field() }}`,
		},
		{
			name:     "Route with target attribute",
			input:    `'route' => 'user.index', 'target' => '_blank'`,
			expected: `<form action="{{ route('user.index') }}" method="GET" target="_blank">`,
		},
		{
			name:  "Route with all attributes",
			input: `'route' => 'user.store', 'method' => 'POST', 'class' => 'user-form', 'id' => 'create-user', 'target' => '_self'`,
			expected: `<form action="{{ route('user.store') }}" method="POST" class="user-form" id="create-user" target="_self">
{{ csrf_field() }}`,
		},
		{
			name:  "Route with array parameters",
			input: `'route' => ['user.update', ['id' => $user->id]], 'method' => 'PUT'`,
			expected: `<form action="{{ route('user.update', ['id' => $user->id]) }}" method="PUT">
{{ csrf_field() }}`,
		},
		{
			name:     "URL with route function",
			input:    `'url' => route('user.index')`,
			expected: `<form action="{{ route('user.index') }}" method="GET">`,
		},
		{
			name:     "Empty content",
			input:    ``,
			expected: `<form action="" method="GET">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processFormOpen(tt.input)
			if result != tt.expected {
				t.Errorf("processFormOpen() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormOpenIntegration(t *testing.T) {
	// 実際のBladeファイルの内容をシミュレート
	input := `
<div class="container">
    {!! Form::open(['route' => 'user.store', 'method' => 'POST', 'class' => 'user-form']) !!}
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" name="name" class="form-control">
        </div>
    {!! Form::close() !!}
</div>
`

	// 修正後の期待値
	expected := `
<div class="container">
    <form action="{{ route('user.store') }}" method="POST" class="user-form">
{{ csrf_field() }}
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" name="name" class="form-control">
        </div>
    {!! Form::close() !!}
</div>
`

	result := replaceFormOpen(input)
	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Integration test failed.\nGot:\n%s\nWant:\n%s", result, expected)
	}
}
