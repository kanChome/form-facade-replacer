# Laravel Form Facade Replacer

Laravel Form Facade を純粋なHTMLタグに変換するGoツールです。Laravel依存を削除し、標準的なBladeテンプレートとHTMLを使用できるようにします。

## 🌐 Language / 言語
- [English](README.md) - English version
- **日本語** (現在) - このページです

## 概要

このツールは、`laravelCollective/html`を用いて生成されたForm Facadeを使用したBladeテンプレートファイル（.blade.php）を解析し、対応するHTMLフォーム要素に自動変換します。変換により、Laravelフレームワークに依存しない、純粋なHTML + Blade構文のテンプレートファイルを生成できます。

## 特徴

- **完全なForm Facade対応**: 24種類のForm Facadeメソッドをサポート
- **動的属性処理**: 条件付きdisabled属性や複雑な三項演算子をサポート
- **文字列連結処理**: PHP文字列連結を適切なBlade構文に自動変換
- **HTML5準拠**: 生成されるHTMLはHTML5標準に準拠
- **AttributeProcessorシステム**: 属性の統一された処理と出力順序の一貫性を保証
- **イベントハンドラー対応**: onClick、onChange等のJavaScript属性を適切に処理（JavaScript文字列リテラル変換機能付き）
- **Blade構文保持**: Laravel独自のBlade構文（@if、@foreach等）を適切に保持
- **CSRF保護**: POST/PUT/PATCH/DELETEリクエストには自動でCSRF保護を追加
- **高性能**: 正規表現キャッシュシステムによる高速処理

## インストール

### 要件
- Go 1.21以上

### インストール手順

```bash
# リポジトリをクローン
git clone https://github.com/ryohirano/form-facade-replacer.git
cd form-facade-replacer

# Go modules初期化（必要に応じて）
go mod tidy

# 開発ビルド
go build -o form-facade-replacer ./cmd/form-facade-replacer

# リリースビルド（バージョン情報付き）
go build -ldflags "-X form-facade-replacer/internal/ffr.version=v2.1.0 -X form-facade-replacer/internal/ffr.buildDate=$(date +'%Y-%m-%d')" -o form-facade-replacer ./cmd/form-facade-replacer
```

## 使用方法

### 基本的な使用法

```bash
# 単一ファイルを処理
./form-facade-replacer path/to/file.blade.php

# ディレクトリを再帰的に処理
./form-facade-replacer path/to/views/directory

# ヘルプを表示
./form-facade-replacer --help
```

### 実行例

```bash
# Laravelのviews配下すべてを処理
./form-facade-replacer resources/views

# 特定のファイルのみ処理
./form-facade-replacer resources/views/user/create.blade.php
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
<input type="checkbox" name="newsletter" value="{{ 'yes' }}" @if(true) checked @endif class="form-check-input">
```

### Form::button / Form::submit

**変換前:**
```php
{{ Form::button('クリック') }}
{{ Form::submit('送信', ['class' => 'btn btn-primary']) }}
```

**変換後:**
```html
<button type="button">{!! 'クリック' !!}</button>
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

### 動的属性処理（条件付き属性）

**変換前:**
```php
{!! Form::button('使用する', [
    'class' => 'btn btn-info', 
    $status ? 'disabled' : '' => $status ? 'disabled' : null
]) !!}
```

**変換後:**
```html
<button class="btn btn-info" {{ $status ? 'disabled' : '' }}="{{ $status ? 'disabled' : null }}">{!! '使用する' !!}</button>
```

### 文字列連結処理

**変換前:**
```php
{!! Form::checkbox('items[]', $item->id, false, [
    'id' => 'item-' . $item->id,
    'class' => 'item-checkbox'
]) !!}
```

**変換後:**
```html
<input type="checkbox" name="items[]" value="{{ $item->id }}" @if(in_array($item->id, (array)false)) checked @endif id="{{ 'item-' . $item->id }}" class="item-checkbox">
```

### イベントハンドラー処理

**変換前:**
```php
{!! Form::checkbox('notifications[]', 'email', old('notifications'), [
    'onClick' => 'toggleNotification(this)',
    'onChange' => 'updateSettings()',
    'class' => 'notification-toggle'
]) !!}
```

**変換後:**
```html
<input type="checkbox" name="notifications[]" value="{{ 'email' }}" @if(in_array('email', (array)old('notifications'))) checked @endif class="notification-toggle" onClick="toggleNotification(this)" onChange="updateSettings()">
```

### JavaScript文字列リテラル変換

**変換前:**
```php
{!! Form::file('image', ['onchange' => 'previewImage(this, "uploads", "preview")']) !!}
```

**変換後:**
```html
<input type="file" name="image" onchange="previewImage(this, 'uploads', "preview")">
```

この機能により、JavaScript内の文字列リテラルが部分的にシングルクォートに変換され、HTMLとJavaScriptの適切な分離が実現されます。

## サポートされるForm Facadeメソッド（24種類）

### 基本フォーム要素
1. **Form::open** - フォーム開始タグ（CSRF保護自動追加）
2. **Form::close** - フォーム終了タグ
3. **Form::text** - テキスト入力フィールド
4. **Form::textarea** - テキストエリア
5. **Form::hidden** - 隠し入力フィールド
6. **Form::label** - ラベル要素

### 選択・チェック要素
7. **Form::checkbox** - チェックボックス（配列対応、動的属性対応）
8. **Form::radio** - ラジオボタン
9. **Form::select** - セレクトボックス（foreachループ生成）

### ボタン要素
10. **Form::button** - 汎用ボタン（動的属性対応）
11. **Form::submit** - 送信ボタン

### 入力タイプ別要素
12. **Form::number** - 数値入力フィールド
13. **Form::email** - メール入力フィールド
14. **Form::password** - パスワード入力フィールド
15. **Form::url** - URL入力フィールド
16. **Form::tel** - 電話番号入力フィールド
17. **Form::search** - 検索入力フィールド
18. **Form::file** - ファイル入力フィールド

### 日時・色・範囲要素
19. **Form::date** - 日付入力フィールド
20. **Form::time** - 時間入力フィールド
21. **Form::datetime** - 日時入力フィールド
22. **Form::range** - 範囲入力フィールド
23. **Form::color** - 色選択フィールド

## 対応パラメータパターン

### Form::open
- `['route' => 'route.name']` - ルート指定
- `['url' => '/path']` - URL直接指定
- `['method' => 'POST/GET/PUT/PATCH/DELETE']` - HTTPメソッド
- `['class' => 'css-class', 'id' => 'element-id']` - HTML属性

### 入力要素共通
- `(name)` - 名前のみ
- `(name, value)` - 名前と値
- `(name, value, [attributes])` - 名前、値、属性

### チェックボックス・ラジオボタン
- `(name, value, checked, [attributes])` - 名前、値、チェック状態、属性
- **配列名サポート**: `name[]` 形式で配列フィールドに対応
- **動的チェック状態**: `old()`, `session()`, `$user->settings` 等

### 高度な属性パターン
- **動的属性**: `$condition ? 'disabled' : '' => $condition ? 'disabled' : null`
- **文字列連結**: `'id' => 'prefix-' . $variable . '-suffix'`
- **イベントハンドラー**: `'onClick' => 'function()', 'onChange' => 'update()'`

## 技術的特徴

### AttributeProcessorシステム
統一された属性処理システムにより、すべてのForm要素で一貫した属性の処理と出力順序を保証します。
- **固定順序**: 属性の出力順序を配列で管理し、常に同じ順序で出力
- **動的属性サポート**: 三項演算子や複雑な条件式を含む動的属性の処理
- **文字列連結処理**: PHP文字列連結を適切なBlade構文に自動変換

### 正規表現キャッシュシステム
高性能な処理のため、使用する正規表現をキャッシュするRegexCacheシステムを実装。
- **並行安全**: `sync.RWMutex`によるスレッドセーフな実装
- **メモリ効率**: 一度コンパイルした正規表現の再利用
- **高速処理**: 大量のファイル処理時のパフォーマンス向上

### HTML5準拠とアクセシビリティ
- **ブール属性**: `disabled`、`required`等は値なしのブール属性として出力
- **無効値処理**: null、空文字列の場合、不要な属性を自動省略
- **W3C標準**: HTML5仕様に完全準拠した要素とアクセシビリティ属性の生成
- **フォームバリデーション**: 適切なinput type属性による自動バリデーション

### 複雑な構文サポート
- **ネストした配列**: 多層配列パラメータの適切な処理
- **PHP関数呼び出し**: `old()`, `session()`, `route()`等の関数の保持
- **論理演算子**: `&&`, `||`を含む複雑な条件式の処理
- **メソッドチェーン**: `$user->settings->get('key')`等のメソッドチェーンサポート

### JavaScript文字列リテラル処理
- **適応的変換**: JavaScript属性内の文字列リテラルを適切に変換
- **部分変換機能**: 最初の文字列リテラルのみをシングルクォートに変換し、複雑なJavaScriptコードの安全性を保持
- **イベントハンドラー最適化**: onClick、onChange等のイベントハンドラー属性で自動適用
- **非貪欲マッチング**: 正規表現による精密な属性境界検出で、複数属性の正確な処理を実現

### CSRF保護とセキュリティ
GET以外のHTTPメソッド（POST、PUT、PATCH、DELETE）使用時に自動で`{{ csrf_field() }}`を追加し、Laravelのセキュリティ機能を維持します。

## テスト

### テスト実行
```bash
# 全テスト実行
go test -v

# カバレッジ付きテスト実行
go test -cover -v

# 特定のテスト実行
go test -run TestFormOpen -v
go test -run TestFormCheckbox -v
```

### テストカバレッジ
本プロジェクトは24種類すべてのForm Facadeメソッドに対応した徹底的なテストスイートを提供します。

#### 基本フォーム要素テスト
- `form_open_test.go` - Form::open機能（ルート、URL、HTTPメソッド）
- `form_close_test.go` - Form::close機能
- `form_text_test.go` - Form::text機能
- `form_textarea_test.go` - Form::textarea機能
- `form_hidden_test.go` - Form::hidden機能
- `form_label_test.go` - Form::label機能

#### 選択・チェック要素テスト  
- `form_checkbox_test.go` - Form::checkbox機能（配列対応、動的属性、イベントハンドラー）
- `form_radio_test.go` - Form::radio機能
- `form_select_test.go` - Form::select機能（foreachループ生成）

#### ボタン要素テスト
- `form_button_test.go` - Form::button機能（動的disabled属性対応）
- `form_submit_test.go` - Form::submit機能

#### 入力タイプ別要素テスト
- `form_number_test.go` - Form::number機能
- `form_email_test.go` - Form::email機能
- `form_password_test.go` - Form::password機能
- `form_url_test.go` - Form::url機能
- `form_tel_test.go` - Form::tel機能
- `form_search_test.go` - Form::search機能
- `form_file_test.go` - Form::file機能

#### 日時・色・範囲要素テスト
- `form_date_test.go` - Form::date機能
- `form_time_test.go` - Form::time機能
- `form_datetime_test.go` - Form::datetime機能
- `form_range_test.go` - Form::range機能
- `form_color_test.go` - Form::color機能

#### 統合テスト
- `integration_test.go` - 実際のユースケース統合テスト

### テストケースの種類
各テストファイルには以下のテストケースが含まれます：
- **基本機能テスト**: 標準的なパラメータパターン
- **属性処理テスト**: HTML属性の正確な処理と順序
- **動的属性テスト**: 条件付き属性の処理
- **文字列連結テスト**: PHP連結のBlade構文変換
- **エッジケーステスト**: null値、空文字列、特殊文字の処理
- **配列フィールドテスト**: `name[]`形式の配列対応
- **イベントハンドラーテスト**: JavaScript属性の適切な処理

## 制限事項

- **Laravel対応**: laravelCollective/htmlを用いて生成したのForm Facade構文をサポート
- **カスタムメソッド**: カスタムForm Facadeメソッドには対応していません
- **極めて複雑な構文**: 5層以上の深いネストした配列は一部制限があります
- **動的メソッド名**: 変数によるメソッド名の動的決定（`Form::$method(...)`）は未対応

## 既知の課題

- 非常に複雑なPHP式の評価については完全ではない場合があります
- 5層以上の深いネストした配列構造は一部制限があります
- JavaScript文字列リテラル変換は最初の文字列のみが対象（設計仕様）

これらの制限事項は将来のバージョンで段階的に解決される予定です。ただし、JavaScript文字列リテラルの部分変換は安全性とパフォーマンスのバランスを考慮した設計仕様です。

## 貢献

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## サポート

問題や質問がある場合は、[GitHub Issues](https://github.com/ryohirano/form-facade-replacer/issues)で報告してください。
