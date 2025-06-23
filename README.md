# Laravel Form Facade Replacer

Laravel Form Facade を純粋なHTMLタグに変換するGoツールです。Laravel依存を削除し、標準的なBladeテンプレートとHTMLを使用できるようにします。

## 概要

このツールは、LaravelのForm Facadeを使用したBladeテンプレートファイル（.blade.php）を解析し、対応するHTMLフォーム要素に自動変換します。変換により、Laravelフレームワークに依存しない、純粋なHTML + Blade構文のテンプレートファイルを生成できます。

## 特徴

- **完全なForm Facade対応**: 12種類のForm Facadeメソッドをサポート
- **HTML5準拠**: 生成されるHTMLはHTML5標準に準拠
- **属性順序の一貫性**: 出力されるHTML属性の順序が常に一致
- **Blade構文保持**: Laravel独自のBlade構文（@if、@foreach等）を適切に保持
- **CSRF保護**: POST/PUT/PATCH/DELETEリクエストには自動でCSRF保護を追加
- **高性能**: Go言語実装による高速処理

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/your-username/form-facade-replacer.git
cd form-facade-replacer

# Go modules初期化（必要に応じて）
go mod tidy

# ビルド
go build -o form_facade_replacer form_facade_replacer.go
```

## 使用方法

### 基本的な使用法

```bash
# 単一ファイルを処理
go run form_facade_replacer.go path/to/file.blade.php

# ディレクトリを再帰的に処理
go run form_facade_replacer.go path/to/views/directory

# ヘルプを表示
go run form_facade_replacer.go --help
```

### 実行例

```bash
# Laravelのviews配下すべてを処理
go run form_facade_replacer.go resources/views

# 特定のファイルのみ処理
go run form_facade_replacer.go resources/views/user/create.blade.php
```

## 対応機能

### Form::open / Form::close

**変換前:**
```php
{!! Form::open(['route' => 'user.store', 'method' => 'POST', 'class' => 'user-form']) !!}
{!! Form::close() !!}
```

**変換後:**
```html
<form action="{{ route('user.store') }}" method="POST" class="user-form">
{{ csrf_field() }}
</form>
```

### Form::text / Form::number

**変換前:**
```php
{{ Form::text('name', $user->name, ['class' => 'form-control', 'placeholder' => 'お名前']) }}
{{ Form::number('age', null, ['min' => 0, 'max' => 120]) }}
```

**変換後:**
```html
<input type="text" name="name" value="{{ $user->name }}" placeholder="お名前" class="form-control">
<input type="number" name="age" min="0" max="120">
```

### Form::textarea

**変換前:**
```php
{{ Form::textarea('message', 'デフォルトメッセージ', ['rows' => 5, 'class' => 'form-control']) }}
```

**変換後:**
```html
<textarea name="message" rows="5" class="form-control">{{ 'デフォルトメッセージ' }}</textarea>
```

### Form::select

**変換前:**
```php
{{ Form::select('category', $categories, $selected, ['class' => 'form-select']) }}
```

**変換後:**
```html
<select name="category" class="form-select">
@foreach($categories as $key => $value)
<option value="{{ $key }}" @if($key == $selected) selected @endif>{{ $value }}</option>
@endforeach
</select>
```

### Form::checkbox

**変換前:**
```php
{{ Form::checkbox('newsletter', 'yes', true, ['class' => 'form-check-input']) }}
```

**変換後:**
```html
<input type="checkbox" name="newsletter" value="{{ yes }}" @if(true) checked @endif class="form-check-input">
```

### Form::button / Form::submit

**変換前:**
```php
{{ Form::button('クリック') }}
{{ Form::submit('送信', ['class' => 'btn btn-primary']) }}
```

**変換後:**
```html
<button>{!! クリック !!}</button>
<button type="submit" class="btn btn-primary">送信</button>
```

### Form::label

**変換前:**
```php
{{ Form::label('name', 'お名前', ['class' => 'form-label']) }}
```

**変換後:**
```html
<label for="name" class="form-label">{!! 'お名前' !!}</label>
```

### Form::hidden

**変換前:**
```php
{!! Form::hidden('user_id', $user->id) !!}
```

**変換後:**
```html
<input type="hidden" name="user_id" value="{{ $user->id }}">
```

## 対応パラメータパターン

各Form Facadeメソッドは、以下のパラメータパターンに対応しています：

### Form::open
- `['route' => 'route.name']` - ルート指定
- `['url' => '/path']` - URL直接指定
- `['method' => 'POST/GET/PUT/PATCH/DELETE']` - HTTPメソッド
- `['class' => 'css-class', 'id' => 'element-id']` - HTML属性

### Form::text / Form::number
- `(name)` - 名前のみ
- `(name, value)` - 名前と値
- `(name, value, [attributes])` - 名前、値、属性

### Form::textarea
- `(name)` - 名前のみ
- `(name, value)` - 名前と値
- `(name, value, [attributes])` - 名前、値、属性

### Form::button
- `(text)` - テキストのみ
- `(text, [attributes])` - テキストと属性

### その他のメソッド
各メソッドは最大4つのパラメータ（名前、値、チェック状態、属性）をサポートします。

## 技術的特徴

### 属性順序の固定化
Go言語のmapは反復順序が非決定的ですが、本ツールでは属性の出力順序を配列で固定し、常に一貫したHTML出力を保証します。

### HTML5準拠
- `disabled`属性は値なしのブール属性として出力
- 無効な値（null、空文字列）の場合、value属性を省略
- 適切なHTML5要素とアクセシビリティ属性の生成

### CSRF保護
GET以外のHTTPメソッド（POST、PUT、PATCH、DELETE）使用時に自動で`{{ csrf_field() }}`を追加します。

## テスト

### テスト実行
```bash
# 全テスト実行
go test ./*test.go ./form_facade_replacer.go -v

# 特定のテスト実行
go test ./form_open_test.go ./form_facade_replacer.go -v
```

### テストカバレッジ
本プロジェクトは以下のテストファイルを含んでいます：
- `form_open_test.go` - Form::open機能のテスト
- `form_close_test.go` - Form::close機能のテスト
- `form_button_test.go` - Form::button機能のテスト
- `form_textarea_test.go` - Form::textarea機能のテスト
- `form_checkbox_test.go` - Form::checkbox機能のテスト
- `form_submit_test.go` - Form::submit機能のテスト
- `form_number_test.go` - Form::number機能のテスト
- `form_select_test.go` - Form::select機能のテスト
- `form_label_test.go` - Form::label機能のテスト

### サンプルデータ
`testdata/`ディレクトリには、実際のBladeファイルと期待される出力HTMLのサンプルが含まれています。

## 制限事項

- Laravel 5.x〜8.x のForm Facade構文をサポート
- ネストした配列形式の複雑な属性は一部制限があります
- カスタムForm Facadeメソッドには対応していません

## 貢献

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 更新履歴

### v1.0.0
- Form Facade変換機能の初回リリース
- 12種類のForm Facadeメソッドをサポート
- HTML5準拠とCSRF保護機能
- 包括的なテストスイート

## サポート

問題や質問がある場合は、[GitHub Issues](https://github.com/kanChome/form-facade-replacer/issues)で報告してください。