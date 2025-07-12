# Laravel Form Facade Replacer

A Go tool that converts Laravel Form Facade syntax to pure HTML tags, enabling framework-independent Blade templates with standard HTML forms.

## 🌐 Language / 言語
- **English** (Current) - You are here
- [日本語 (Japanese)](README.ja.md) - 日本語版はこちら

## Overview

This tool analyzes Blade template files (.blade.php) that use Laravel's Form Facade and automatically converts them to corresponding HTML form elements. The conversion eliminates Laravel framework dependencies while generating pure HTML + Blade syntax template files.

## Features

- **Complete Form Facade Support**: Supports 24 types of Form Facade methods
- **Dynamic Attribute Processing**: Handles conditional disabled attributes and complex ternary operators
- **String Concatenation Processing**: Automatically converts PHP string concatenation to appropriate Blade syntax
- **HTML5 Compliance**: Generated HTML adheres to HTML5 standards
- **AttributeProcessor System**: Ensures unified attribute processing and consistent output ordering
- **Event Handler Support**: Properly processes JavaScript attributes like onClick and onChange (with JavaScript string literal conversion)
- **Blade Syntax Preservation**: Maintains Laravel-specific Blade syntax (@if, @foreach, etc.)
- **CSRF Protection**: Automatically adds CSRF protection for POST/PUT/PATCH/DELETE requests
- **High Performance**: Fast processing through regex caching system
- **Comprehensive Test Coverage**: Thorough test suite with 3,557 lines of code

## Installation

### Requirements
- Go 1.21 or higher

### Installation Steps

```bash
# Clone the repository
git clone https://github.com/ryohirano/form-facade-replacer.git
cd form-facade-replacer

# Initialize Go modules (if needed)
go mod tidy

# Development build
go build -o form_facade_replacer form_facade_replacer.go

# Release build (with version information)
go build -ldflags "-X main.version=v2.1.0 -X main.buildDate=$(date +'%Y-%m-%d')" -o form_facade_replacer form_facade_replacer.go
```

## Usage

### Basic Usage

```bash
# Process a single file
go run form_facade_replacer.go path/to/file.blade.php

# Process a directory recursively
go run form_facade_replacer.go path/to/views/directory

# Display help
go run form_facade_replacer.go --help

# Show version information
go run form_facade_replacer.go --version
```

### Usage Examples

```bash
# Process all files in Laravel's views directory
go run form_facade_replacer.go resources/views

# Process a specific file only
go run form_facade_replacer.go resources/views/user/create.blade.php
```

## Supported Features

### Form::open / Form::close

**Before:**
```php
{!! Form::open(['route' => 'user.store', 'method' => 'POST', 'class' => 'user-form']) !!}
{!! Form::close() !!}
```

**After:**
```html
<form action="{{ route('user.store') }}" method="POST" class="user-form">
{{ csrf_field() }}
</form>
```

### Form::text / Form::number

**Before:**
```php
{{ Form::text('name', $user->name, ['class' => 'form-control', 'placeholder' => 'Your Name']) }}
{{ Form::number('age', null, ['min' => 0, 'max' => 120]) }}
```

**After:**
```html
<input type="text" name="name" value="{{ $user->name }}" placeholder="Your Name" class="form-control">
<input type="number" name="age" min="0" max="120">
```

### Form::textarea

**Before:**
```php
{{ Form::textarea('message', 'Default Message', ['rows' => 5, 'class' => 'form-control']) }}
```

**After:**
```html
<textarea name="message" rows="5" class="form-control">{{ 'Default Message' }}</textarea>
```

### Form::select

**Before:**
```php
{{ Form::select('category', $categories, $selected, ['class' => 'form-select']) }}
```

**After:**
```html
<select name="category" class="form-select">
@foreach($categories as $key => $value)
<option value="{{ $key }}" @if($key == $selected) selected @endif>{{ $value }}</option>
@endforeach
</select>
```

### Form::checkbox

**Before:**
```php
{{ Form::checkbox('newsletter', 'yes', true, ['class' => 'form-check-input']) }}
```

**After:**
```html
<input type="checkbox" name="newsletter" value="{{ 'yes' }}" @if(true) checked @endif class="form-check-input">
```

### Form::button / Form::submit

**Before:**
```php
{{ Form::button('Click Me') }}
{{ Form::submit('Submit', ['class' => 'btn btn-primary']) }}
```

**After:**
```html
<button type="button">{!! 'Click Me' !!}</button>
<button type="submit" class="btn btn-primary">Submit</button>
```

### Form::label

**Before:**
```php
{{ Form::label('name', 'Your Name', ['class' => 'form-label']) }}
```

**After:**
```html
<label for="name" class="form-label">{!! 'Your Name' !!}</label>
```

### Form::hidden

**Before:**
```php
{!! Form::hidden('user_id', $user->id) !!}
```

**After:**
```html
<input type="hidden" name="user_id" value="{{ $user->id }}">
```

### Dynamic Attribute Processing (Conditional Attributes)

**Before:**
```php
{!! Form::button('Use This', [
    'class' => 'btn btn-info', 
    $status ? 'disabled' : '' => $status ? 'disabled' : null
]) !!}
```

**After:**
```html
<button class="btn btn-info" {{ $status ? 'disabled' : '' }}="{{ $status ? 'disabled' : null }}">{!! 'Use This' !!}</button>
```

### String Concatenation Processing

**Before:**
```php
{!! Form::checkbox('items[]', $item->id, false, [
    'id' => 'item-' . $item->id,
    'class' => 'item-checkbox'
]) !!}
```

**After:**
```html
<input type="checkbox" name="items[]" value="{{ $item->id }}" @if(in_array($item->id, (array)false)) checked @endif id="{{ 'item-' . $item->id }}" class="item-checkbox">
```

### Event Handler Processing

**Before:**
```php
{!! Form::checkbox('notifications[]', 'email', old('notifications'), [
    'onClick' => 'toggleNotification(this)',
    'onChange' => 'updateSettings()',
    'class' => 'notification-toggle'
]) !!}
```

**After:**
```html
<input type="checkbox" name="notifications[]" value="{{ 'email' }}" @if(in_array('email', (array)old('notifications'))) checked @endif class="notification-toggle" onClick="toggleNotification(this)" onChange="updateSettings()">
```

### JavaScript String Literal Conversion

**Before:**
```php
{!! Form::file('image', ['onchange' => 'previewImage(this, "uploads", "preview")']) !!}
```

**After:**
```html
<input type="file" name="image" onchange="previewImage(this, 'uploads', "preview")">
```

This feature enables partial conversion of JavaScript string literals to single quotes, achieving proper separation between HTML and JavaScript.

## Supported Form Facade Methods (24 Types)

### Basic Form Elements
1. **Form::open** - Form opening tag (with automatic CSRF protection)
2. **Form::close** - Form closing tag
3. **Form::text** - Text input field
4. **Form::textarea** - Textarea element
5. **Form::hidden** - Hidden input field
6. **Form::label** - Label element

### Selection & Check Elements
7. **Form::checkbox** - Checkbox (supports arrays, dynamic attributes)
8. **Form::radio** - Radio button
9. **Form::select** - Select box (generates foreach loops)

### Button Elements
10. **Form::button** - General button (supports dynamic attributes)
11. **Form::submit** - Submit button

### Input Type-Specific Elements
12. **Form::number** - Number input field
13. **Form::email** - Email input field
14. **Form::password** - Password input field
15. **Form::url** - URL input field
16. **Form::tel** - Telephone input field
17. **Form::search** - Search input field
18. **Form::file** - File input field

### Date, Color & Range Elements
19. **Form::date** - Date input field
20. **Form::time** - Time input field
21. **Form::datetime** - DateTime input field
22. **Form::range** - Range input field
23. **Form::color** - Color picker field
24. **Form::input** - Generic input handler

## Supported Parameter Patterns

### Form::open
- `['route' => 'route.name']` - Route specification
- `['url' => '/path']` - Direct URL specification
- `['method' => 'POST/GET/PUT/PATCH/DELETE']` - HTTP method
- `['class' => 'css-class', 'id' => 'element-id']` - HTML attributes

### Common Input Elements
- `(name)` - Name only
- `(name, value)` - Name and value
- `(name, value, [attributes])` - Name, value, and attributes

### Checkbox & Radio Button
- `(name, value, checked, [attributes])` - Name, value, checked state, attributes
- **Array Support**: Supports array fields with `name[]` format
- **Dynamic Check State**: Supports `old()`, `session()`, `$user->settings`, etc.

### Advanced Attribute Patterns
- **Dynamic Attributes**: `$condition ? 'disabled' : '' => $condition ? 'disabled' : null`
- **String Concatenation**: `'id' => 'prefix-' . $variable . '-suffix'`
- **Event Handlers**: `'onClick' => 'function()', 'onChange' => 'update()'`

## Technical Features

### AttributeProcessor System
A unified attribute processing system that ensures consistent attribute handling and output ordering for all Form elements.
- **Fixed Ordering**: Manages attribute output order through arrays, ensuring consistent output
- **Dynamic Attribute Support**: Handles dynamic attributes with ternary operators and complex conditional expressions
- **String Concatenation Processing**: Automatically converts PHP string concatenation to appropriate Blade syntax

### Regex Caching System
Implements a RegexCache system for high-performance processing by caching frequently used regular expressions.
- **Concurrent Safety**: Thread-safe implementation using `sync.RWMutex`
- **Memory Efficiency**: Reuses compiled regular expressions
- **High Performance**: Improves performance when processing large numbers of files

### HTML5 Compliance and Accessibility
- **Boolean Attributes**: Outputs boolean attributes like `disabled`, `required` without values
- **Invalid Value Handling**: Automatically omits unnecessary attributes for null or empty values
- **W3C Standards**: Generates elements and accessibility attributes in full compliance with HTML5 specifications
- **Form Validation**: Enables automatic validation through appropriate input type attributes

### Complex Syntax Support
- **Nested Arrays**: Proper handling of multi-dimensional array parameters
- **PHP Function Calls**: Preserves functions like `old()`, `session()`, `route()`
- **Logical Operators**: Handles complex conditional expressions with `&&`, `||`
- **Method Chaining**: Supports method chaining like `$user->settings->get('key')`

### JavaScript String Literal Processing
- **Adaptive Conversion**: Appropriately converts string literals within JavaScript attributes
- **Partial Conversion**: Converts only the first string literal to single quotes, maintaining safety for complex JavaScript code
- **Event Handler Optimization**: Automatically applied to event handler attributes like onClick, onChange
- **Non-Greedy Matching**: Achieves precise attribute boundary detection through regex, enabling accurate processing of multiple attributes

### CSRF Protection and Security
Automatically adds `{{ csrf_field() }}` for non-GET HTTP methods (POST, PUT, PATCH, DELETE), maintaining Laravel's security features.

## Testing

### Running Tests
```bash
# Run all tests
go test -v

# Run tests with coverage
go test -cover -v

# Run specific tests
go test -run TestFormOpen -v
go test -run TestFormCheckbox -v
```

### Comprehensive Test Coverage (3,557 lines)
This project provides a thorough test suite covering all 24 Form Facade methods:

#### Basic Form Element Tests
- `form_open_test.go` - Form::open functionality (routes, URLs, HTTP methods)
- `form_close_test.go` - Form::close functionality
- `form_text_test.go` - Form::text functionality
- `form_textarea_test.go` - Form::textarea functionality
- `form_hidden_test.go` - Form::hidden functionality
- `form_label_test.go` - Form::label functionality

#### Selection & Check Element Tests  
- `form_checkbox_test.go` - Form::checkbox functionality (array support, dynamic attributes, event handlers)
- `form_radio_test.go` - Form::radio functionality
- `form_select_test.go` - Form::select functionality (foreach loop generation)

#### Button Element Tests
- `form_button_test.go` - Form::button functionality (supports dynamic disabled attributes)
- `form_submit_test.go` - Form::submit functionality

#### Input Type-Specific Element Tests
- `form_number_test.go` - Form::number functionality
- `form_email_test.go` - Form::email functionality
- `form_password_test.go` - Form::password functionality
- `form_url_test.go` - Form::url functionality
- `form_tel_test.go` - Form::tel functionality
- `form_search_test.go` - Form::search functionality
- `form_file_test.go` - Form::file functionality

#### Date, Color & Range Element Tests
- `form_date_test.go` - Form::date functionality
- `form_time_test.go` - Form::time functionality
- `form_datetime_test.go` - Form::datetime functionality
- `form_range_test.go` - Form::range functionality
- `form_color_test.go` - Form::color functionality

#### Integration & Special Feature Tests
- `integration_test.go` - Real-world use case integration tests
- `dynamic_attribute_detector_test.go` - Dynamic attribute detection functionality tests
- `dynamic_disabled_test.go` - Dynamic disabled attribute tests

### Test Case Categories
Each test file includes the following types of test cases:
- **Basic Functionality Tests**: Standard parameter patterns
- **Attribute Processing Tests**: Accurate HTML attribute handling and ordering
- **Dynamic Attribute Tests**: Conditional attribute processing
- **String Concatenation Tests**: PHP concatenation to Blade syntax conversion
- **Edge Case Tests**: Handling of null values, empty strings, special characters
- **Array Field Tests**: Support for `name[]` format arrays
- **Event Handler Tests**: Proper processing of JavaScript attributes

### testdata Directory
Provides sample sets of complex Blade files and expected HTML output used in real projects:
- `testdata/blades/` - Test Blade files
- `testdata/expected/` - Expected HTML output
- `testdata/run_tests.sh` - Batch test execution script

## Limitations

- **Laravel Support**: Supports Laravel 5.x~8.x Form Facade syntax
- **Custom Methods**: Custom Form Facade methods are not supported
- **Extremely Complex Syntax**: Arrays with 5 or more levels of deep nesting have some limitations
- **Dynamic Method Names**: Dynamic method name resolution using variables (`Form::$method(...)`) is not supported

## Known Issues

- Evaluation of very complex PHP expressions may not be complete in some cases
- Array structures with 5 or more levels of deep nesting have some limitations
- JavaScript string literal conversion only applies to the first string (design specification)

These limitations are planned to be gradually resolved in future versions. However, the partial conversion of JavaScript string literals is a design specification that balances safety and performance considerations.

## Contributing

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## License

This project is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## Version History

### v2.1.0 (Latest)
- **Support for 24 Form Facade methods** (major expansion from 12 methods in v1.0.0)
- **JavaScript String Literal Conversion feature** added (properly converts strings within event handlers)
- **Advanced Event Handler Processing** (accurate processing of multiple attributes like onClick, onChange)
- **Non-greedy Matching Regular Expressions** implementation (enhanced support for complex attribute structures)
- **Dynamic Attribute Processing feature** added (conditional disabled attributes, ternary operator support)
- **String Concatenation Processing feature** added (automatic conversion of PHP string concatenation to Blade syntax)
- **AttributeProcessor System** implementation (unified attribute processing and order guarantee)
- **Regex Caching System** implementation (performance improvement)
- **TDD Development Methodology** adoption (Test-Driven Development)
- **Comprehensive Test Suite** (3,557 lines of test code)
- **Array Field Enhancement** (complete support for `name[]` format)
- **Complex Syntax Support** (nested arrays, logical operators, method chaining)

### v1.0.0
- Initial release of Form Facade conversion functionality
- Support for 12 Form Facade methods
- HTML5 compliance and CSRF protection features
- Basic test suite

## Support

For issues or questions, please report them at [GitHub Issues](https://github.com/ryohirano/form-facade-replacer/issues).
