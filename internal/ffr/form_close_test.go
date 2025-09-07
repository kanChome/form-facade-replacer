package ffr

import (
	"testing"
)

func TestReplaceFormClose(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Form::close() with {!! !!}",
			input:    `{!! Form::close() !!}`,
			expected: `</form>`,
		},
		{
			name:     "Basic Form::close() with {{ }}",
			input:    `{{ Form::close() }}`,
			expected: `</form>`,
		},
		{
			name:     "Form::close() with extra spaces inside braces",
			input:    `{!!  Form::close()  !!}`,
			expected: `</form>`,
		},
		{
			name:     "Form::close() with extra spaces in double braces",
			input:    `{{  Form::close()  }}`,
			expected: `</form>`,
		},
		{
			name:     "Form::close() with no spaces",
			input:    `{!!Form::close()!!}`,
			expected: `</form>`,
		},
		{
			name:     "Form::close() in double braces with no spaces",
			input:    `{{Form::close()}}`,
			expected: `</form>`,
		},
		{
			name:     "Multiple Form::close() in same text",
			input:    `{!! Form::close() !!} some content {{ Form::close() }}`,
			expected: `</form> some content </form>`,
		},
		{
			name: "Form::close() in HTML context",
			input: `<div class="form-container">
    {!! Form::open(['route' => 'test']) !!}
        <input type="text" name="field">
    {!! Form::close() !!}
</div>`,
			expected: `<div class="form-container">
    {!! Form::open(['route' => 'test']) !!}
        <input type="text" name="field">
    </form>
</div>`,
		},
		{
			name:     "Inline Form::close()",
			input:    `<p>Form: {!! Form::open(['route' => 'test']) !!}<input type="text">{!! Form::close() !!}</p>`,
			expected: `<p>Form: {!! Form::open(['route' => 'test']) !!}<input type="text"></form></p>`,
		},
		{
			name:     "Mixed brackets Form::close()",
			input:    `{!! Form::open() !!} content {{ Form::close() }} more content {!! Form::close() !!}`,
			expected: `{!! Form::open() !!} content </form> more content </form>`,
		},
		{
			name: "Form::close() with newlines",
			input: `{!! 
    Form::close() 
!!}`,
			expected: `</form>`,
		},
		{
			name:     "Form::close() with tabs and spaces",
			input:    `{!!	 Form::close() 	!!}`,
			expected: `</form>`,
		},
		{
			name:     "No Form::close() in text",
			input:    `<div>Some regular HTML content</div>`,
			expected: `<div>Some regular HTML content</div>`,
		},
		{
			name:     "Text containing 'close' but not Form::close()",
			input:    `<button onclick="closeDialog()">Close</button>`,
			expected: `<button onclick="closeDialog()">Close</button>`,
		},
		{
			name:     "Empty string",
			input:    ``,
			expected: ``,
		},
		{
			name:     "Form::close() at start of text",
			input:    `{!! Form::close() !!} following content`,
			expected: `</form> following content`,
		},
		{
			name:     "Form::close() at end of text",
			input:    `preceding content {!! Form::close() !!}`,
			expected: `preceding content </form>`,
		},
		{
			name:     "Multiple consecutive Form::close()",
			input:    `{!! Form::close() !!}{!! Form::close() !!}{{ Form::close() }}`,
			expected: `</form></form></form>`,
		},
		{
			name:     "Form::close() with extra parentheses (should not match)",
			input:    `{!! Form::close()() !!}`,
			expected: `{!! Form::close()() !!}`,
		},
		{
			name:     "Form::close() with arguments (should not match)",
			input:    `{!! Form::close('arg') !!}`,
			expected: `{!! Form::close('arg') !!}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormClose(tt.input)
			if result != tt.expected {
				t.Errorf("replaceFormClose() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestReplaceFormCloseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Large text with multiple Form::close()",
			input: `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
    <div class="container">
        {!! Form::open(['route' => 'users.store']) !!}
            <div class="form-group">
                <input type="text" name="name">
            </div>
        {!! Form::close() !!}
        
        <hr>
        
        {{ Form::open(['route' => 'posts.store']) }}
            <textarea name="content"></textarea>
        {{ Form::close() }}
    </div>
</body>
</html>`,
			expected: `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
    <div class="container">
        {!! Form::open(['route' => 'users.store']) !!}
            <div class="form-group">
                <input type="text" name="name">
            </div>
        </form>
        
        <hr>
        
        {{ Form::open(['route' => 'posts.store']) }}
            <textarea name="content"></textarea>
        </form>
    </div>
</body>
</html>`,
		},
		{
			name:     "Form::close() with Unicode characters around",
			input:    `前の内容 {!! Form::close() !!} 後の内容`,
			expected: `前の内容 </form> 後の内容`,
		},
		{
			name:     "Form::close() in Blade comments (should still replace)",
			input:    `{{-- Comment: {!! Form::close() !!} --}} {!! Form::close() !!}`,
			expected: `{{-- Comment: </form> --}} </form>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceFormClose(tt.input)
			if result != tt.expected {
				t.Errorf("replaceFormClose() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestReplaceFormClosePerformance(t *testing.T) {
	// パフォーマンステスト用の大きなテキスト
	largeText := ""
	for i := 0; i < 1000; i++ {
		largeText += `<div>Some content {!! Form::close() !!} more content</div>`
	}

	result := replaceFormClose(largeText)

	// 正しく置換されているかチェック
	expectedCount := 1000
	actualCount := 0
	for i := 0; i < len(result)-7; i++ {
		if result[i:i+7] == "</form>" {
			actualCount++
		}
	}

	if actualCount != expectedCount {
		t.Errorf("Expected %d </form> tags, got %d", expectedCount, actualCount)
	}
}

func BenchmarkReplaceFormClose(b *testing.B) {
	input := `<div class="container">
        {!! Form::open(['route' => 'test']) !!}
            <input type="text" name="field">
        {!! Form::close() !!}
    </div>
    <div class="another">
        {{ Form::open() }}
            <textarea name="content"></textarea>
        {{ Form::close() }}
    </div>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceFormClose(input)
	}
}

func BenchmarkReplaceFormCloseLarge(b *testing.B) {
	// 大きなテキストでのベンチマーク
	largeText := ""
	for i := 0; i < 100; i++ {
		largeText += `<div>Content {!! Form::close() !!} more content {{ Form::close() }}</div>`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceFormClose(largeText)
	}
}
