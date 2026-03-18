# EnrichedMarkdownText API Reference

Complete API for the `EnrichedMarkdownText` component from `react-native-enriched-markdown`.

## Table of Contents

- [Props](#props)
- [Style Properties](#style-properties)
- [Supported Markdown Elements](#supported-markdown-elements)
- [LaTeX Math](#latex-math)
- [Accessibility](#accessibility)
- [RTL Support](#rtl-support)
- [Image Caching](#image-caching)
- [Copy Options](#copy-options)
- [Disabling LaTeX Math](#disabling-latex-math)

## Props

### `markdown` (required)

The Markdown content to render.

| Type | Platform |
|------|----------|
| `string` | Both |

### `flavor`

Markdown flavor. `'github'` enables GFM tables, task lists, and block math.

| Type | Default | Platform |
|------|---------|----------|
| `'commonmark' \| 'github'` | `'commonmark'` | Both |

Layout differences:
- `'commonmark'`: All content renders as a single TextView. Text selection covers all content.
- `'github'`: Content is split into segments (text blocks + table views). Text selection is per-segment.

### `markdownStyle`

Style configuration for Markdown elements. See [Style Properties](#style-properties).

| Type | Default | Platform |
|------|---------|----------|
| `MarkdownStyle` | `{}` | Both |

Memoize with `useMemo` to avoid unnecessary re-renders.

### `containerStyle`

Style for the outer container view.

| Type | Default | Platform |
|------|---------|----------|
| `ViewStyle` | - | Both |

### `selectable`

Whether text can be selected.

| Type | Default | Platform |
|------|---------|----------|
| `boolean` | `true` | Both |

### `allowFontScaling`

Whether fonts scale with Text Size accessibility settings.

| Type | Default | Platform |
|------|---------|----------|
| `boolean` | `true` | Both |

### `maxFontSizeMultiplier`

Maximum font scale when `allowFontScaling` is enabled.

| Type | Default | Platform |
|------|---------|----------|
| `number` | `undefined` | Both |

### `allowTrailingMargin`

Whether to keep the bottom margin of the last block element.

| Type | Default | Platform |
|------|---------|----------|
| `boolean` | `false` | Both |

### `enableLinkPreview`

Controls native link preview on long press (iOS). Automatically `false` when `onLinkLongPress` is provided.

| Type | Default | Platform |
|------|---------|----------|
| `boolean` | `true` | iOS |

### `md4cFlags`

Configuration for the md4c parser.

| Type | Default | Platform |
|------|---------|----------|
| `Md4cFlags` | `{ underline: false }` | Both |

Properties:
- `underline`: When `true`, `_text_` renders as underline instead of italic. Only `*text*` works for italic.
- `latexMath`: When `false`, `$` is treated as plain text (disables math rendering).

### `onLinkPress`

Fires when a link is tapped.

```ts
(event: { url: string }) => void
```

### `onLinkLongPress`

Fires when a link is long-pressed. On iOS, providing this automatically disables the system link preview.

```ts
(event: { url: string }) => void
```

### `onTaskListItemPress`

Fires when a task list checkbox is tapped (requires `flavor="github"`).

```ts
(event: { index: number; checked: boolean; text: string }) => void
```

- `index`: 0-based position of the task item
- `checked`: New state after toggling
- `text`: The item's text content

## Style Properties

All styles are set via the `markdownStyle` prop.

### Style inheritance model

Block elements (paragraph, headings, lists, blockquote, code blocks) establish their own typography context with `fontSize`, `fontFamily`, `fontWeight`, `color`, `marginTop`, `marginBottom`, and `lineHeight`.

Inline elements (strong, emphasis, links, inline code) inherit the parent block's typography and apply additional styling on top.

### Platform defaults

| Property | iOS | Android |
|----------|-----|---------|
| System font | SF Pro | Roboto |
| Monospace font | Menlo | monospace |
| Line height | Tighter (0.75x multiplier) | Standard |

### Block styles (paragraph, h1-h6, blockquote, list, codeBlock)

| Property | Type | Description |
|----------|------|-------------|
| `fontSize` | `number` | Font size in points |
| `fontFamily` | `string` | Font family name |
| `fontWeight` | `string` | Font weight |
| `color` | `string` | Text color |
| `marginTop` | `number` | Top margin |
| `marginBottom` | `number` | Bottom margin |
| `lineHeight` | `number` | Line height |

### Paragraph and heading-specific (paragraph, h1-h6)

| Property | Type | Description |
|----------|------|-------------|
| `textAlign` | `'auto' \| 'left' \| 'right' \| 'center' \| 'justify'` | Text alignment (default: `'left'`) |

### Blockquote-specific

| Property | Type | Description |
|----------|------|-------------|
| `borderColor` | `string` | Left border color |
| `borderWidth` | `number` | Left border width |
| `gapWidth` | `number` | Gap between border and text |
| `backgroundColor` | `string` | Background color |

### List-specific

| Property | Type | Description |
|----------|------|-------------|
| `bulletColor` | `string` | Bullet color (unordered lists) |
| `bulletSize` | `number` | Bullet size |
| `markerColor` | `string` | Number marker color (ordered lists) |
| `markerFontWeight` | `string` | Number marker font weight |
| `gapWidth` | `number` | Gap between marker and text |
| `marginLeft` | `number` | Left margin for nesting |

### Code block-specific

| Property | Type | Description |
|----------|------|-------------|
| `backgroundColor` | `string` | Background color |
| `borderColor` | `string` | Border color |
| `borderRadius` | `number` | Corner radius |
| `borderWidth` | `number` | Border width |
| `padding` | `number` | Inner padding |

### Inline code-specific

| Property | Type | Description |
|----------|------|-------------|
| `fontFamily` | `string` | Font family (defaults to platform monospace: SF Mono / monospace) |
| `fontSize` | `number` | Font size (defaults to parent block's size) |
| `color` | `string` | Text color |
| `backgroundColor` | `string` | Background color |
| `borderColor` | `string` | Border color |

### Link-specific

| Property | Type | Description |
|----------|------|-------------|
| `fontFamily` | `string` | Font family (overrides parent block's) |
| `color` | `string` | Link text color |
| `underline` | `boolean` | Show underline |

### Strong-specific

| Property | Type | Description |
|----------|------|-------------|
| `fontFamily` | `string` | Font family (when set, replaces parent block's font) |
| `fontWeight` | `'bold' \| 'normal'` | Set `'normal'` to use `fontFamily` as-is without bold trait |
| `color` | `string` | Bold text color |

### Emphasis-specific

| Property | Type | Description |
|----------|------|-------------|
| `fontFamily` | `string` | Font family (when set, replaces parent block's font) |
| `fontStyle` | `'italic' \| 'normal'` | Set `'normal'` to use `fontFamily` as-is without italic trait |
| `color` | `string` | Italic text color |

### Strikethrough-specific

| Property | Type | Description |
|----------|------|-------------|
| `color` | `string` | Strikethrough line color (iOS only) |

### Underline-specific

| Property | Type | Description |
|----------|------|-------------|
| `color` | `string` | Underline color (iOS only) |

### Image-specific

| Property | Type | Description |
|----------|------|-------------|
| `height` | `number` | Image height |
| `borderRadius` | `number` | Corner radius |
| `marginTop` | `number` | Top margin |
| `marginBottom` | `number` | Bottom margin |

### Inline image-specific

| Property | Type | Description |
|----------|------|-------------|
| `size` | `number` | Image size (square) |

### Thematic break (horizontal rule)

| Property | Type | Description |
|----------|------|-------------|
| `color` | `string` | Line color |
| `height` | `number` | Line thickness |
| `marginTop` | `number` | Top margin |
| `marginBottom` | `number` | Bottom margin |

### Table-specific (requires `flavor="github"`)

Tables inherit base block styles and add:

| Property | Type | Description |
|----------|------|-------------|
| `headerFontFamily` | `string` | Header cell font family |
| `headerBackgroundColor` | `string` | Header row background |
| `headerTextColor` | `string` | Header row text color |
| `rowEvenBackgroundColor` | `string` | Even data row background |
| `rowOddBackgroundColor` | `string` | Odd data row background |
| `borderColor` | `string` | Grid line color |
| `borderWidth` | `number` | Grid line width |
| `borderRadius` | `number` | Table container corner radius |
| `cellPaddingHorizontal` | `number` | Horizontal cell padding |
| `cellPaddingVertical` | `number` | Vertical cell padding |

### Task list-specific (requires `flavor="github"`)

| Property | Type | Description |
|----------|------|-------------|
| `checkedColor` | `string` | Checked checkbox background |
| `borderColor` | `string` | Unchecked checkbox border |
| `checkmarkColor` | `string` | Checkmark color |
| `checkboxSize` | `number` | Checkbox size (defaults to 90% of list font size) |
| `checkboxBorderRadius` | `number` | Checkbox corner radius |
| `checkedTextColor` | `string` | Text color for checked items |
| `checkedStrikethrough` | `boolean` | Strikethrough on checked items |

### Math block-specific (requires `flavor="github"`)

| Property | Type | Description |
|----------|------|-------------|
| `fontSize` | `number` | Equation font size |
| `color` | `string` | Equation text color |
| `backgroundColor` | `string` | Block background |
| `padding` | `number` | Inner padding |
| `marginTop` | `number` | Top margin |
| `marginBottom` | `number` | Bottom margin |
| `textAlign` | `'left' \| 'center' \| 'right'` | Alignment (default: `'center'`) |

### Inline math-specific

| Property | Type | Description |
|----------|------|-------------|
| `color` | `string` | Equation text color |

## Supported Markdown Elements

### Block elements

| Element | Syntax | Style Property |
|---------|--------|----------------|
| Headings | `# H1` to `###### H6` | `h1` - `h6` |
| Paragraphs | Plain text | `paragraph` |
| Blockquotes | `> Quote` | `blockquote` |
| Code blocks | ` ``` code ``` ` | `codeBlock` |
| Unordered lists | `- Item`, `* Item`, `+ Item` | `list` |
| Ordered lists | `1. Item` | `list` |
| Task lists | `- [x] Done`, `- [ ] Todo` | `taskList` |
| Thematic break | `---`, `***`, `___` | `thematicBreak` |
| Images | `![alt](url)` | `image` |
| Tables | `\| col \| col \|` | `table` |
| Math block | `$$...$$` | `math` |

### Inline elements

| Element | Syntax | Style Property |
|---------|--------|----------------|
| Bold | `**text**` or `__text__` | `strong` |
| Italic | `*text*` or `_text_` | `em` |
| Underline | `_text_` (with `md4cFlags`) | `underline` |
| Strikethrough | `~~text~~` | `strikethrough` |
| Links | `[text](url)` | `link` |
| Inline code | `` `code` `` | `code` |
| Inline images | `![alt](url)` (within text) | `inlineImage` |
| Inline math | `$...$` | `inlineMath` |

Images are automatically categorized as block or inline based on context: standalone images are block-level, images within text paragraphs are inline.

Lists and blockquotes support unlimited nesting depth with automatic indentation.

## LaTeX Math

- **Inline math** (`$...$`): rendered within text flow, works with both flavors
- **Block math** (`$$...$$`): rendered as standalone display element, requires `flavor="github"`, must be on its own line

Backslashes in LaTeX commands need escaping in JS strings. Use `String.raw` or double backslashes:

```tsx
// Double backslashes
const math = '$$x = \\frac{-b \\pm \\sqrt{b^2-4ac}}{2a}$$';

// String.raw
const math = String.raw`$$x = \frac{-b \pm \sqrt{b^2-4ac}}{2a}$$`;
```

## Accessibility

### iOS (VoiceOver)

Custom rotors for navigating headings, links, and images. Headings use `UIAccessibilityTraitHeader` with level info. Links are activatable with `UIAccessibilityTraitLink`. Images announce alt text.

### Android (TalkBack)

Reading controls for headings, links, and images. List items announce position and type. Nested list items get a "Nested" prefix. Images announce alt text or "Image" as fallback.

### Announcements

| Element | Announcement |
|---------|-------------|
| Headings | "Welcome to Markdown, heading" |
| Links | "React Native, link" |
| Images (with alt) | Alt text content |
| Images (no alt) | "Image" |
| Unordered list | "bullet point" |
| Ordered list | "list item 1", "list item 2" |

## RTL Support

RTL is automatic on Android. On iOS, call `I18nManager.forceRTL(true)` before the root component mounts and restart the app.

```tsx
import { I18nManager } from 'react-native';
I18nManager.forceRTL(true); // Call before app renders
```

This affects the entire app layout on iOS.

RTL behavior by element:
- Paragraphs, headings: right-aligned with RTL direction
- Lists: bullets/numbers on the right, text indented from right
- Task lists: checkboxes on the right
- Blockquotes: border on the right
- Tables: columns right-to-left
- Code blocks: always LTR
- Copy as HTML: includes `dir="rtl"`

## Image Caching

Three-tier caching with no configuration needed:

| Layer | Android | iOS | Size |
|-------|---------|-----|------|
| Originals (memory) | `LruCache` | `NSCache` | 20 MB |
| Processed variants (memory) | `LruCache` | `NSCache` | 30 MB |
| Disk | OkHttp `Cache` | `NSURLCache` | 100 MB |

- **Original cache**: decoded images keyed by URL. Android downsamples large images to screen width.
- **Processed cache**: scaled/clipped variants keyed by URL + dimensions + border radius. Repeated layouts with the same geometry skip all processing.
- **Disk cache**: raw HTTP responses respecting cache headers.
- **Request deduplication**: multiple simultaneous requests for the same URL make only one network call.

## Copy Options

When text is selected, the context menu provides:

### Smart Copy (default "Copy" action)

**iOS**: Copies in multiple formats simultaneously (plain text, Markdown, HTML, RTF, RTFD). Receiving apps pick the richest format they support.

**Android**: Copies as plain text and HTML.

### Copy as Markdown

Dedicated menu item that copies only the Markdown source text.

### Copy Image URL

Appears when selection contains images. Copies the image source URL. On Android with multiple images, all URLs are copied (one per line).

## Disabling LaTeX Math

To reduce bundle size, remove the native math libraries:

### JS level

```tsx
<EnrichedMarkdownText md4cFlags={{ latexMath: false }} markdown="Price is $5" />
```

### iOS native

Add to Podfile, then re-run `pod install`:
```ruby
ENV['ENRICHED_MARKDOWN_ENABLE_MATH'] = '0'
```
Removes iosMath (~2.5 MB).

### Android native

Add to `gradle.properties`:
```properties
enrichedMarkdown.enableMath=false
```

### Expo

```json
{
  "expo": {
    "plugins": [
      ["react-native-enriched-markdown", { "enableMath": false }]
    ]
  }
}
```

Run `npx expo prebuild` after adding. To re-enable, remove the plugin or set `enableMath: true` and run `npx expo prebuild --clean`.
