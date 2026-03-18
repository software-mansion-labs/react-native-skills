# Scroll, Events, and Utilities

Patterns for scroll-driven animations, event-based reactions, frame callbacks, measurement, and value mapping utilities in Reanimated 4.

---

## Scroll-Driven Animations

### useAnimatedScrollHandler

Respond to scroll events with multiple handlers:

```tsx
const scrollHandler = useAnimatedScrollHandler({
  onScroll: (event, context) => {
    offset.value = event.contentOffset.y;
  },
  onBeginDrag: (event, context) => {
    context.startY = event.contentOffset.y;
  },
  onEndDrag: (event, context) => {
    if (event.contentOffset.y - context.startY > 100) {
      offset.value = withSpring(200);
    }
  },
});

<Animated.ScrollView onScroll={scrollHandler}>
  {children}
</Animated.ScrollView>
```

Available handlers: `onScroll`, `onBeginDrag`, `onEndDrag`, `onMomentumBegin`, `onMomentumEnd`.

The `context` object is shared between all handlers for the same component, letting you pass state across events (e.g., save drag start position in `onBeginDrag`, read it in `onEndDrag`).

Passing a single function instead of an object is treated as `onScroll`.

**Gotchas:**
- Must use `Animated.ScrollView`, not plain `ScrollView`.
- On Web, only `onScroll` fires. Other events are iOS/Android only.

### useScrollOffset

Simpler alternative when you only need the scroll position as a shared value:

```tsx
const animatedRef = useAnimatedRef<Animated.ScrollView>();
const scrollOffset = useScrollOffset(animatedRef);

const headerStyle = useAnimatedStyle(() => ({
  opacity: interpolate(scrollOffset.value, [0, 100], [1, 0]),
}));
```

Automatically detects horizontal or vertical scroll. Works with `ScrollView`, `FlatList`, and `FlashList`. The ref can be changed at runtime.

### scrollTo

Programmatic scrolling from the UI thread:

```tsx
const animatedRef = useAnimatedRef<Animated.ScrollView>();

// In a worklet or useDerivedValue
scrollTo(animatedRef, 0, targetY, true); // (ref, x, y, animated)
```

Can only be called from the UI thread. Wrap with `runOnUI()` when calling from JS-thread event handlers.

---

## Value Mapping Utilities

### interpolate

Maps a numeric value from one range to another:

```tsx
const opacity = interpolate(scrollOffset.value, [0, 100, 200], [1, 0.5, 0]);
```

```tsx
function interpolate(
  value: number,
  input: number[],
  output: number[],
  extrapolation?: Extrapolation | { extrapolateLeft?, extrapolateRight? }
): number;
```

**Extrapolation modes** (behavior for values outside the input range):
- `Extrapolation.EXTEND` (default) — linear extrapolation beyond the range
- `Extrapolation.CLAMP` — clamps to the nearest edge of the output range
- `Extrapolation.IDENTITY` — returns the input value as-is

You can set different modes per edge:

```tsx
interpolate(value, [0, 100], [0, 1], {
  extrapolateLeft: Extrapolation.CLAMP,
  extrapolateRight: Extrapolation.EXTEND,
});
```

Input values must be in increasing order.

### interpolateColor

Maps a numeric value to a color, producing smooth color transitions:

```tsx
const color = interpolateColor(
  progress.value,
  [0, 1],
  ['#ff0000', '#0000ff'],
  'RGB'
);
```

**Color spaces:**
- `'RGB'` (default) — linear RGB with gamma correction (default gamma: 2.2)
- `'HSV'` — hue-saturation-value; smooth hue transitions
- `'LAB'` — Oklab; perceptually uniform color differences

Set `gamma: 1` to disable gamma correction. Use `useCorrectedHSVInterpolation: true` (default) to prevent long hue paths (e.g., red-to-blue going through green).

Returns color in `rgba(r, g, b, a)` format.

### clamp

```tsx
const clamped = clamp(value, min, max);
```

Constrains a number between `min` and `max`. Use with scroll offsets, touch positions, or derived animation values to prevent out-of-bounds behavior.

---

## Event Reactions

### useAnimatedReaction

React to shared value changes with access to both current and previous values:

```tsx
useAnimatedReaction(
  () => Math.floor(scrollOffset.value / PAGE_HEIGHT),
  (currentPage, previousPage) => {
    if (previousPage !== null && currentPage !== previousPage) {
      scheduleOnRN(onPageChanged, currentPage);
    }
  }
);
```

The `prepare` function transforms/filters shared values before comparison. The `react` function fires when the prepared value changes.

**Critical:** Do not mutate the same shared value in `react` that you track in `prepare`. This causes an infinite loop.

Use `prepare` to reduce callback frequency (e.g., `Math.floor()` to react only on whole page changes instead of every pixel).

---

## Frame Callbacks

### useFrameCallback

Run logic on every frame (60Hz or 120Hz depending on the device):

```tsx
const frameCallback = useFrameCallback((frameInfo) => {
  // frameInfo.timestamp — system time in ms
  // frameInfo.timeSincePreviousFrame — ms since last frame (null on first frame)
  // frameInfo.timeSinceFirstFrame — ms since callback activated
  progress.value += (frameInfo.timeSincePreviousFrame ?? 0) * speed;
});

// Pause/resume
frameCallback.setActive(false);
frameCallback.setActive(true);
```

- `autostart` (second parameter, default `true`) controls whether the callback begins immediately.
- Always memoize the callback with `useCallback` to avoid recreation on every render.
- Use time deltas (`timeSincePreviousFrame`) for frame-rate-independent animations.

---

## Measurement

### measure

Synchronously get a view's dimensions and position on the UI thread:

```tsx
const animatedRef = useAnimatedRef<Animated.View>();

const animatedStyle = useAnimatedStyle(() => {
  if (!_WORKLET) return {}; // Guard: first evaluation runs on JS thread

  const measurements = measure(animatedRef);
  if (measurements === null) return {};

  return {
    transform: [{ translateY: -measurements.height }],
  };
});
```

Returns `{ x, y, width, height, pageX, pageY }` or `null` if the component is unmounted or off-screen (e.g., recycled FlatList items).

**Rules:**
- Always check for `null` before using measurements.
- In `useAnimatedStyle`, guard with `if (!_WORKLET) return {}` because the first evaluation runs on the JS thread where `measure` is unavailable.
- Wrap with `runOnUI()` when calling from JS-thread event handlers.
- Not available with Remote JS Debugger (use Chrome DevTools).

### setNativeProps

Imperatively update a component's properties from the UI thread:

```tsx
setNativeProps(animatedRef, { backgroundColor: 'red' });
```

Runs on the UI thread only. Designed for gesture handlers where you need instant updates without going through `useAnimatedStyle`. Prefer `useAnimatedStyle` and `useAnimatedProps` for most cases.

### dispatchCommand

Call native component commands from the UI thread:

```tsx
dispatchCommand(animatedRef, 'focus');
dispatchCommand(animatedRef, 'scrollToEnd', [true]);
```

Available commands vary by component (e.g., `focus`, `blur`, `clear` for TextInput; `scrollToEnd` for ScrollView). Android and iOS only.

---

## Device Sensors

### useAnimatedSensor

Track device motion for parallax, tilt, or orientation-based animations:

```tsx
import { useAnimatedSensor, SensorType } from 'react-native-reanimated';

const sensor = useAnimatedSensor(SensorType.ROTATION, {
  interval: 'auto', // match screen refresh rate
});

const style = useAnimatedStyle(() => ({
  transform: [
    { rotateX: `${sensor.sensor.value.pitch}rad` },
    { rotateY: `${sensor.sensor.value.roll}rad` },
  ],
}));
```

**Sensor types:** `ACCELEROMETER`, `GYROSCOPE`, `GRAVITY`, `MAGNETIC_FIELD`, `ROTATION`.

**Data formats:**
- Accelerometer/Gyroscope/Gravity/Magnetic Field: `{ x, y, z, interfaceOrientation }`
- Rotation: `{ pitch, roll, yaw, qw, qx, qy, qz, interfaceOrientation }`

**Config:** `interval` (`'auto'` or ms), `adjustToInterfaceOrientation` (default `true`), `iosReferenceFrame`.

**Gotchas:**
- iOS requires location services enabled (Settings > Privacy > Location Services).
- Web requires HTTPS.
- Most sensors operate at up to 100Hz.
- Use `sensor.unregister()` to stop listening.

---

## Keyboard (Deprecated)

`useAnimatedKeyboard` is deprecated in Reanimated 4. Use `react-native-keyboard-controller` for keyboard-aware animations.
