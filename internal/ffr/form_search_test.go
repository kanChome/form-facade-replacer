package ffr

import (
	"testing"
)

func TestFormSearch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic search field",
			input:    `{{ Form::search('query') }}`,
			expected: `<input type="search" name="query" value="">`,
		},
		{
			name:     "Search with value",
			input:    `{{ Form::search('query', 'Laravel') }}`,
			expected: `<input type="search" name="query" value="{{ 'Laravel' }}">`,
		},
		{
			name:     "Search with old() helper",
			input:    `{{ Form::search('search', old('search')) }}`,
			expected: `<input type="search" name="search" value="{{ old('search') }}">`,
		},
		{
			name:     "Search with attributes",
			input:    `{{ Form::search('search', old('search'), ['placeholder' => 'Search...', 'class' => 'form-control']) }}`,
			expected: `<input type="search" name="search" value="{{ old('search') }}" placeholder="Search..." class="form-control">`,
		},
		{
			name:     "Search with double exclamation marks",
			input:    `{!! Form::search('product_search', request('search'), ['class' => 'search-input', 'autocomplete' => 'off']) !!}`,
			expected: `<input type="search" name="product_search" value="{{ request('search') }}" class="search-input">`,
		},
		{
			name: "Multi-line search field",
			input: `{!! Form::search('global_search', old('global_search'), [
    'placeholder' => 'Search products, categories...',
    'class' => 'form-control global-search',
    'id' => 'global-search-input'
]) !!}`,
			expected: `<input type="search" name="global_search" value="{{ old('global_search') }}" placeholder="Search products, categories..." class="form-control global-search" id="global-search-input">`,
		},
		{
			name:     "Search with results attribute",
			input:    `{{ Form::search('site_search', old('site_search'), ['results' => '10', 'class' => 'site-search']) }}`,
			expected: `<input type="search" name="site_search" value="{{ old('site_search') }}" class="site-search">`,
		},
		{
			name:     "Search with complex name",
			input:    `{{ Form::search('filters[' . $category . '][search]', old('filters[' . $category . '][search]'), ['class' => 'filter-search']) }}`,
			expected: `<input type="search" name="filters[{{ $category }}][search]" value="{{ old('filters[' . $category . '][search]') }}" class="filter-search">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormSearch(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
