# Animation Functions and Core Hooks

API reference for Reanimated 4 core hooks, animation functions, and animation modifiers.

For choosing between animation approaches (CSS vs shared value), see **`animations.md`**.

---

## Core Hooks

### useSharedValue

```tsx
const sv = useSharedValue(initialValue);
sv.value = 42;                    // read/write via .value
sv.set(v => v + 1);              // updater function (React Compiler safe)
const current = sv.get();         // read (React Compiler safe)
```

**Gotchas:**
- Never destructure: `const { value } = sv` breaks reactivity.
- For objects, reassign the entire value: `sv.value = { ...sv.value, x: 50 }`. Direct mutation (`sv.value.x = 50`) loses reactivity.
- For large arrays/objects, use `.modify()` to mutate in place: `sv.modify(arr => { arr.push(item); return arr; })`.
- Reading `.value` on the JS thread blocks until the UI thread syncs. Minimize cross-thread reads.
- Never read/modify during component render. Access only in callbacks (`useAnimatedStyle`, event handlers, `useEffect`).
- Use `.get()`/`.set()` methods instead of `.value` for React Compiler compatibility.

### useAnimatedStyle

```tsx
const animatedStyle = useAnimatedStyle(() => ({
  transform: [{ translateX: withSpring(offset.value) }],
}));

<Animated.View style={[styles.box, animatedStyle]} />
```

**Rules:**
- Keep static styles in `StyleSheet.create()`. Only put dynamic parts in `useAnimatedStyle`.
- Animated styles override static styles in the style array.
- Removing an animated style does not unset its values. Manually set properties to `undefined` to clear them.
- Never mutate shared values inside the updater (e.g., `sv.value = withTiming(1)` in the callback). This causes infinite loops.
- The callback runs on the JS thread first, then immediately on the UI thread. Use `global._WORKLET` to guard thread-specific code.

### useAnimatedProps

```tsx
const animatedProps = useAnimatedProps(() => ({
  text: String(Math.round(progress.value)),
}));

<AnimatedTextInput animatedProps={animatedProps} />
```

For animating component properties (not styles). Use adapters for third-party components whose prop names differ between JS and native:

```tsx
import { SVGAdapter } from 'react-native-reanimated';

const props = useAnimatedProps(() => ({ cx: x.value }), null, SVGAdapter);
```

Custom color properties require manual `processColor()` wrapping. Define adapters outside the component body to avoid recalculation.

### useDerivedValue

```tsx
const doubled = useDerivedValue(() => sv.value * 2);
```

Creates a read-only shared value that recomputes when its dependencies change. Runs on the UI thread automatically. The `.set()` method is deprecated and will be removed.

If you need access to the previous value, use `useAnimatedReaction` instead.

### createAnimatedComponent

```tsx
const AnimatedTextInput = Animated.createAnimatedComponent(TextInput);
```

Wraps a React Native component to accept animated styles and props. Function components **must** be wrapped with `React.forwardRef()`. Class components work directly.

Built-in animated components: `Animated.View`, `Animated.Text`, `Animated.Image`, `Animated.ScrollView`, `Animated.FlatList`.

### useAnimatedRef

```tsx
const animatedRef = useAnimatedRef<Animated.View>();
<Animated.View ref={animatedRef} />
```

Returns a ref usable with `measure()`, `scrollTo()`, and `useScrollOffset()`. The ref value (`current`) is `null` until the component mounts. It is only accessible from the JS thread, so do not read it inside worklets.

### cancelAnimation

```tsx
cancelAnimation(sharedValue);
```

Stops a running animation. The shared value retains its current position. Safe to call on non-animated values (no-op). To resume, assign a new animation: `sv.value = withSpring(target)`.

---

## Animation Functions

### withTiming

```tsx
sv.value = withTiming(toValue, config?, callback?);
```

| Config | Type | Default |
|--------|------|---------|
| `duration` | number (ms) | `300` |
| `easing` | EasingFunction | `Easing.inOut(Easing.quad)` |
| `reduceMotion` | ReduceMotion | `System` |

**Easing catalog:**
- Base: `Easing.linear`, `.quad`, `.cubic`, `.sin`, `.circle`, `.exp`, `.bounce`, `.elastic(bounciness?)`, `.poly(n)`, `.back(s?)`
- Bezier: `Easing.bezier(x1, y1, x2, y2)`
- Modifiers: `Easing.in(fn)`, `Easing.out(fn)`, `Easing.inOut(fn)`
- Default (`Easing.inOut(Easing.quad)`) gives smooth acceleration/deceleration.

### withSpring

```tsx
sv.value = withSpring(toValue, config?, callback?);
```

Two configuration modes (cannot mix):

**Physics-based** (stiffness/damping):

| Config | Type | Default |
|--------|------|---------|
| `stiffness` | number | `900` |
| `damping` | number | `120` |

**Duration-based** (duration/dampingRatio):

| Config | Type | Default |
|--------|------|---------|
| `duration` | number (ms) | `550` |
| `dampingRatio` | number | `1` (critically damped) |

`dampingRatio` values: `< 1` = underdamped (bouncy), `1` = critically damped (no bounce, fastest settle), `> 1` = overdamped (slow, no bounce).

**Shared config** (both modes):

| Config | Type | Default |
|--------|------|---------|
| `mass` | number | `4` |
| `velocity` | number | `0` |
| `overshootClamping` | boolean | `false` |
| `clamp` | { min?, max? } | — |
| `reduceMotion` | ReduceMotion | `System` |

If both physics-based and duration-based configs are provided, duration-based overrides.

### withDecay

```tsx
sv.value = withDecay(config, callback?);
```

Simulates inertial motion (friction). Starts at a velocity and decelerates to a stop. Ideal for gesture flings.

| Config | Type | Default |
|--------|------|---------|
| `velocity` | number | `0` |
| `deceleration` | number | `0.998` |
| `clamp` | [min, max] | — |
| `velocityFactor` | number | `1` |
| `rubberBandEffect` | boolean | `false` |
| `rubberBandFactor` | number | `0.6` |
| `reduceMotion` | ReduceMotion | `System` |

`clamp` is **required** when `rubberBandEffect` is `true`. The rubber band effect makes the animation bounce at clamp boundaries instead of stopping.

---

## Animation Modifiers

Modifiers wrap animation functions to add delay, repetition, sequencing, or clamping.

### withDelay

```tsx
sv.value = withDelay(delayMs, animation, reduceMotion?);
```

Delays the start of an animation. The animation itself is unmodified.

### withRepeat

```tsx
sv.value = withRepeat(animation, numberOfReps?, reverse?, callback?, reduceMotion?);
```

| Param | Type | Default |
|-------|------|---------|
| `numberOfReps` | number | `2` |
| `reverse` | boolean | `false` |

- Non-positive values (`0`, `-1`) repeat infinitely until cancelled or unmounted.
- `reverse: true` creates a ping-pong effect (plays forward, then backward).
- **`reverse` only works with animation functions** (`withSpring`, `withTiming`). It does **not** work with animation modifiers like `withSequence`.

### withSequence

```tsx
sv.value = withSequence(animation1, animation2, ...moreAnimations);
```

Runs animations one after another on the same shared value. Requires at least two animations. The first positional argument can optionally be a `ReduceMotion` value.

### withClamp

```tsx
sv.value = withClamp({ min: 0, max: 100 }, withSpring(target));
```

Limits the animated value range. Designed for `withSpring` to prevent overshoot beyond boundaries. When the spring hits a clamped boundary, its dampingRatio is automatically reduced.

---

## Composition Patterns

Modifiers nest freely:

```tsx
// Staggered entrance
items.forEach((_, i) => {
  sv[i].value = withDelay(i * 100, withSpring(1));
});

// Infinite ping-pong
sv.value = withRepeat(withTiming(1, { duration: 800 }), -1, true);

// Multi-step sequence
sv.value = withSequence(
  withTiming(50, { duration: 200 }),
  withSpring(0),
  withDelay(300, withTiming(100))
);

// Clamped spring
sv.value = withClamp({ min: 0, max: 200 }, withSpring(scrollTarget));
```

### Callback behavior

Callbacks on `withTiming`, `withSpring`, `withDecay`, and `withRepeat` are automatically workletized and run on the UI thread. They receive `(finished: boolean, current: AnimatableValue)` where `finished` is `true` if the animation completed normally, `false` if cancelled.
