# EnrichedTextInput API Reference

Complete API for the `EnrichedTextInput` component from `react-native-enriched`.

## Table of Contents

- [Props](#props)
- [Ref Methods](#ref-methods)
- [HtmlStyle Customization](#htmlstyle-customization)
- [Context Menu Items](#context-menu-items)

## Props

### `autoFocus`

If `true`, focuses the input on mount.

| Type | Default | Platform |
|------|---------|----------|
| `bool` | `false` | Both |

### `autoCapitalize`

Controls automatic capitalization: `'characters'`, `'words'`, `'sentences'`, or `'none'`.

| Type | Default | Platform |
|------|---------|----------|
| `'none' \| 'sentences' \| 'words' \| 'characters'` | `'sentences'` | Both |

### `contextMenuItems`

Custom items for the native text editing menu. See [Context Menu Items](#context-menu-items).

| Type | Default | Platform |
|------|---------|----------|
| `ContextMenuItem[]` | `[]` | Both |

On iOS, items appear before system items (Copy/Paste/Cut). On Android, items may appear in a submenu depending on the device.

### `cursorColor`

Color of the cursor/caret.

| Type | Default | Platform |
|------|---------|----------|
| `color` | system default | Android |

### `defaultValue`

Initial HTML value. If the string is valid HTML from `EnrichedTextInput` output (or compatible HTML), styles will be applied.

| Type | Default | Platform |
|------|---------|----------|
| `string` | - | Both |

### `editable`

If `false`, disables user interaction. Programmatic changes via ref methods still work.

| Type | Default | Platform |
|------|---------|----------|
| `bool` | `true` | Both |

### `htmlStyle`

Customizes the appearance of formatted elements. See [HtmlStyle Customization](#htmlstyle-customization).

| Type | Default | Platform |
|------|---------|----------|
| `HtmlStyle` | defaults | Both |

### `linkRegex`

Custom regex for link detection. Pass `null` to disable auto-detection entirely.

| Type | Default | Platform |
|------|---------|----------|
| `RegExp \| null` | platform default | Both |

### `mentionIndicators`

Characters that trigger mention creation. Each must be a single character.

| Type | Default | Platform |
|------|---------|----------|
| `string[]` | `['@']` | Both |

### `placeholder`

Placeholder text shown when input is empty.

| Type | Default | Platform |
|------|---------|----------|
| `string` | `''` | Both |

### `placeholderTextColor`

Color of the placeholder text.

| Type | Default | Platform |
|------|---------|----------|
| `color` | input's color | Both |

### `selectionColor`

Color of the selection rectangle. On iOS, the cursor also uses this color.

| Type | Default | Platform |
|------|---------|----------|
| `color` | system default | Both |

### `style`

Accepts most `ViewStyle` props plus these `TextStyle` props: `color`, `fontFamily`, `fontSize`, `fontWeight`, `lineHeight`, `fontStyle` (Android only), `lineHeight` (iOS only).

### `androidExperimentalSynchronousEvents` (EXPERIMENTAL)

If `true`, uses experimental synchronous events on Android to prevent input flickering when updating component size.

| Type | Default | Platform |
|------|---------|----------|
| `bool` | `false` | Android |

### `useHtmlNormalizer` (EXPERIMENTAL)

If `true`, external HTML pasted into the input (from Google Docs, Word, web pages) is normalized into the tag subset the enriched parser understands.

| Type | Default | Platform |
|------|---------|----------|
| `bool` | `false` | Both |

## Event Callbacks

### `onBlur`

Fires when the input loses focus.

```ts
() => void
```

### `onFocus`

Fires when the input gains focus.

```ts
() => void
```

### `onChangeText`

Returns the plain text value on each change. Omit this callback if you don't need plain text, as continuous text extraction has a performance cost.

```ts
(event: NativeSyntheticEvent<{ value: string }>) => void
```

### `onChangeHtml`

Returns the HTML string on each change. Omit this callback if you only need HTML at save time (use `getHTML()` ref method instead), as continuous HTML parsing is expensive.

```ts
(event: NativeSyntheticEvent<{ value: string }>) => void
```

### `onChangeState`

Fires when style state changes. The payload has one entry per style, each with `isActive`, `isBlocking`, and `isConflicting` booleans. Styles tracked: `bold`, `italic`, `underline`, `strikeThrough`, `inlineCode`, `h1`-`h6`, `codeBlock`, `blockQuote`, `orderedList`, `unorderedList`, `checkboxList`, `link`, `image`, `mention`.

```ts
(event: NativeSyntheticEvent<OnChangeStateEvent>) => void
```

### `onChangeSelection`

Fires on selection or cursor changes.

```ts
interface OnChangeSelectionEvent {
  start: number;  // Selection start index
  end: number;    // First index after selection end (start === end for cursor)
  text: string;   // Text within the selection
}
```

### `onStartMention`

Fires when mention editing begins (indicator typed or cursor moved back to unfinished mention).

```ts
(indicator: string) => void
```

### `onChangeMention`

Fires when the user types or deletes characters after a mention indicator.

```ts
(event: { indicator: string; text: string }) => void
```

### `onEndMention`

Fires when the user leaves mention editing (moved cursor away or typed a space).

```ts
(indicator: string) => void
```

### `onLinkDetected`

Fires when a link is added or the cursor moves near one.

```ts
(event: { text: string; url: string; start: number; end: number }) => void
```

### `onMentionDetected`

Fires when a mention is added or the cursor moves near one.

```ts
(event: { text: string; indicator: string; attributes: Record<string, string> }) => void
```

### `onKeyPress`

Fires on key press. Follows React Native TextInput's `onKeyPress` spec.

```ts
(event: NativeSyntheticEvent<{ key: string }>) => void
```

### `onPasteImages`

Fires when images/GIFs are pasted. Returns an array of image details.

```ts
(event: NativeSyntheticEvent<{
  images: { uri: string; type: string; width: number; height: number }[]
}>) => void
```

## Ref Methods

All methods are called on `ref.current` where `ref` is `useRef<EnrichedTextInputInstance>(null)`.

### Style toggles

```ts
toggleBold(): void
toggleItalic(): void
toggleUnderline(): void
toggleStrikeThrough(): void
toggleInlineCode(): void
toggleH1(): void
toggleH2(): void
toggleH3(): void
toggleH4(): void
toggleH5(): void
toggleH6(): void
toggleCodeBlock(): void
toggleBlockQuote(): void
toggleOrderedList(): void
toggleUnorderedList(): void
toggleCheckboxList(checked: boolean): void
```

### Content methods

```ts
// Get HTML on-demand (preferred over onChangeHtml for save-time use)
getHTML(): Promise<string>

// Set content (HTML string or plain text)
setValue(value: string): void
```

### Selection

```ts
setSelection(start: number, end: number): void
```

### Focus

```ts
focus(): void
blur(): void
```

### Links

```ts
// Set a link. Replaces text between start and end with the link text.
setLink(start: number, end: number, text: string, url: string): void

// Remove link styling from range (preserves text content)
removeLink(start: number, end: number): void
```

### Mentions

```ts
// Start a new mention at cursor/selection position
startMention(indicator: string): void

// Complete the current mention
setMention(indicator: string, text: string, attributes?: Record<string, string>): void
```

### Images

```ts
// Insert image at cursor (or replace selection). You must provide correct dimensions.
setImage(src: string, width: number, height: number): void
```

## Context Menu Items

Extend the native text editing menu with custom items. Supported on Android and iOS 16+.

```tsx
<EnrichedTextInput
  ref={ref}
  contextMenuItems={[
    {
      text: 'Paste Link',
      onPress: ({ text, selection, styleState }) => {
        if (!styleState.link.isBlocking) {
          ref.current?.setLink(selection.start, selection.end, text, url);
        }
      },
      visible: true,
    },
  ]}
/>
```

```ts
interface ContextMenuItem {
  text: string;       // Title displayed in menu
  onPress: (args: {
    text: string;
    selection: { start: number; end: number };
    styleState: OnChangeStateEvent;
  }) => void;
  visible?: boolean;  // Defaults to true
}
```

## HtmlStyle Customization

Customize the appearance of formatted elements via the `htmlStyle` prop.

### Headings (h1-h6)

| Property | Type | Default |
|----------|------|---------|
| `fontSize` | `number` | H1: 32, H2: 24, H3: 20, H4: 16, H5: 14, H6: 12 |
| `bold` | `boolean` | `false` |

### Blockquote

| Property | Type | Default |
|----------|------|---------|
| `borderColor` | `color` | `'darkgray'` |
| `borderWidth` | `number` | `4` |
| `gapWidth` | `number` | `16` |
| `color` | `color` | input's color |

### Code block

| Property | Type | Default |
|----------|------|---------|
| `color` | `color` | `'black'` |
| `borderRadius` | `number` | `8` |
| `backgroundColor` | `color` | `'darkgray'` |

### Inline code

| Property | Type | Default |
|----------|------|---------|
| `color` | `color` | `'red'` |
| `backgroundColor` | `color` | `'darkgray'` |

### Links (a)

| Property | Type | Default |
|----------|------|---------|
| `color` | `color` | `'blue'` |
| `textDecorationLine` | `'underline' \| 'none'` | `'underline'` |

### Mentions

Set a single config for all mention types, or a record keyed by indicator for per-indicator styling.

| Property | Type | Default |
|----------|------|---------|
| `color` | `color` | `'blue'` |
| `backgroundColor` | `color` | `'yellow'` |
| `textDecorationLine` | `'underline' \| 'none'` | `'underline'` |

```tsx
htmlStyle={{
  mention: {
    '@': { color: 'blue', backgroundColor: '#E3F2FD' },
    '#': { color: 'green', backgroundColor: '#E8F5E9' },
  },
}}
```

### Ordered list (ol)

| Property | Type | Default |
|----------|------|---------|
| `gapWidth` | `number` | `16` |
| `marginLeft` | `number` | `16` |
| `markerFontWeight` | `fontWeight` | input's fontWeight |
| `markerColor` | `color` | input's color |

### Unordered list (ul)

| Property | Type | Default |
|----------|------|---------|
| `bulletColor` | `color` | `'black'` |
| `bulletSize` | `number` | `8` |
| `marginLeft` | `number` | `16` |
| `gapWidth` | `number` | `16` |

### Checkbox list (ulCheckbox)

| Property | Type | Default |
|----------|------|---------|
| `boxColor` | `color` | `'blue'` |
| `boxSize` | `number` | `24` |
| `marginLeft` | `number` | `16` |
| `gapWidth` | `number` | `16` |
