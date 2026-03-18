# Animations

Production-quality animation patterns for React Native apps using Reanimated 4 on the New Architecture.

For animation function APIs and core hooks, see **`animation-functions.md`**.
For entering/exiting and layout transition animations, see **`layout-animations.md`**.
For scroll-driven animations and event-based patterns, see **`scroll-and-events.md`**.
For GPU shader animations (particles, noise, SDF, physics, 3D), see **`gpu-animations.md`**.
For performance tuning and feature flags, see **`animations-performance.md`**.

---

## Decision Tree

Pick the animation type based on what drives the animation and what it needs to compute.

```
Does the effect require GPU-level computation?
(Particle systems, fluid/physics sims, procedural noise, SDF shapes, 3D scenes,
 or more simultaneously animated elements than Reanimated can handle)
├── YES → Use GPU Shaders (react-native-wgpu + TypeGPU)   → see gpu-animations.md
└── NO  → Is the animation driven by a state change (not a gesture or continuous input)?
    ├── YES → Can it be expressed as a simple A→B property transition?
    │   ├── YES → Use CSS Transition (transitionProperty)
    │   └── NO  → Does it need a defined keyframe sequence?
    │       ├── YES → Use CSS Animation (animationName + keyframes)
    │       └── NO  → Use CSS Transition with multiple properties
    └── NO  → Is it gesture-driven, or does it need math / trig / layout reads?
        └── YES → Use Shared Value Animation (useSharedValue + useAnimatedStyle)
```

Default to CSS transitions and CSS animations. They are declarative, easier to read, and remove the overhead of worklet execution. Reach for shared values only when the animation requires programmatic control that CSS cannot express. Reach for GPU shaders when the animation involves per-pixel computation, hundreds/thousands of independently animated elements, physics simulations, or 3D rendering that operates outside the React Native view hierarchy.

---

## CSS Transitions

Use when a component's style should animate smoothly whenever a state-driven prop changes.

```tsx
<Animated.View
  style={{
    width: isExpanded ? 200 : 100,
    transitionProperty: 'width',
    transitionDuration: 300,
    transitionTimingFunction: 'ease-out',
  }}
/>
```

### Properties

| Property | Values | Default |
|----------|--------|---------|
| `transitionProperty` | Style property name, array of names, `'all'`, `'none'` | — |
| `transitionDuration` | Number (ms), `'300ms'`, `'0.3s'`, or array | `0` |
| `transitionDelay` | Number (ms), string, or array. Supports negative values | `0` |
| `transitionTimingFunction` | `'linear'`, `'ease'`, `'ease-in'`, `'ease-out'`, `'ease-in-out'`, `cubicBezier(x1,y1,x2,y2)`, `steps(n, modifier)`, `linear(...points)`, or array | `'ease'` |
| `transitionBehavior` | `'normal'`, `'allow-discrete'` | `'normal'` |

When using arrays, the order must match the `transitionProperty` array:

```tsx
transitionProperty: ['width', 'opacity', 'backgroundColor'],
transitionDuration: [300, 200, 150],
transitionTimingFunction: ['ease-out', 'linear', 'ease-in-out'],
```

### Timing functions

Predefined: `'linear'`, `'ease'`, `'ease-in'`, `'ease-out'`, `'ease-in-out'`, `'step-start'`, `'step-end'`.

Parametrized:
- `cubicBezier(0.25, 0.1, 0.5, 2)` — custom Bezier curve
- `steps(4, 'jump-end')` — discrete steps; modifiers: `'jump-start'`, `'jump-end'` (default), `'jump-none'`, `'jump-both'`
- `linear(0, [0.25, '75%'], 1)` — polygonal chain with control points

### Discrete properties

Properties like `flexDirection`, `justifyContent`, and `alignItems` cannot be smoothly animated. By default, they change instantly. To make them flip at the animation midpoint, set:

```tsx
transitionBehavior: 'allow-discrete',
```

The `display` property flips at animation start (0%) instead of the midpoint. For smoother transitions of discrete properties, use Layout Animations instead.

### Rules

- Avoid `transitionProperty: 'all'` — it forces evaluation of every style property on each frame and degrades performance.
- Negative delays start the transition partway through (e.g., `'-5s'` on a 10s transition starts at 50%).
- CSS transitions cannot animate discrete properties smoothly without `transitionBehavior: 'allow-discrete'`.

---

## CSS Animations

Use when the animation follows a predefined keyframe sequence independent of external state — loaders, pulse effects, entrance choreography.

```tsx
const pulse = {
  '0%':   { opacity: 1 },
  '50%':  { opacity: 0.4 },
  '100%': { opacity: 1 },
};

<Animated.View
  style={{
    animationName: pulse,
    animationDuration: '1200ms',
    animationIterationCount: 'infinite',
    animationTimingFunction: 'ease-in-out',
  }}
/>
```

Reanimated uses the current element state as the implicit `0%` keyframe, so you only need to define the frames that differ. At minimum, one keyframe is required.

### Properties

| Property | Values | Default |
|----------|--------|---------|
| `animationName` | Keyframes object, array of keyframes objects, `'none'` | — |
| `animationDuration` | Number (ms), `'300ms'`, `'0.3s'`, or array | `0` |
| `animationDelay` | Number (ms), string, or array. Supports negative values | `0` |
| `animationTimingFunction` | Same as transition timing functions, or array | `'ease'` |
| `animationIterationCount` | `'infinite'`, number (supports fractions like `0.5`), or array | `1` |
| `animationDirection` | `'normal'`, `'reverse'`, `'alternate'`, `'alternate-reverse'`, or array | `'normal'` |
| `animationFillMode` | `'none'`, `'forwards'`, `'backwards'`, `'both'`, or array | `'none'` |
| `animationPlayState` | `'running'`, `'paused'`, or array | `'running'` |

When using arrays, the order must match the `animationName` array.

### Keyframe selectors

Use percentage strings (`'0%'`, `'50%'`, `'100%'`), aliases (`'from'`/`'to'`), or floats between 0 and 1.

### Multiple animations

```tsx
const fadeInOut = { '0%': { opacity: 0 }, '100%': { opacity: 1 } };
const moveLeft = { '100%': { transform: [{ translateX: -100 }] } };

<Animated.View
  style={{
    animationName: [fadeInOut, moveLeft],
    animationDuration: ['2.5s', '5s'],
    animationIterationCount: ['infinite', 1],
  }}
/>
```

If multiple animations target the same property, the later animation in the array wins.

### Rules

- The timing function on the last keyframe (`100%`, `to`, or `1`) is ignored — there is no subsequent keyframe to animate toward.
- All properties in the `transform` array must appear in the same order across all keyframes.
- Avoid `animationFillMode: 'forwards'` or `'both'` with fractional `animationIterationCount` and relative units (percentages). If the parent resizes after the animation, the child retains stale dimensions.
- For infinite CSS animations, set `animationIterationCount: 'infinite'`. The animation stops automatically on unmount — no manual cleanup needed.
- Negative delays start the animation partway through its cycle.
- Pause and resume with `animationPlayState: 'paused'` / `'running'`.

---

## Shared Value Animations

Use when:
- The animation is driven by a gesture or continuous input (scroll position, drag offset)
- It requires math, trigonometric functions, or interpolation between computed values
- It needs to read layout measurements on each frame (`measure`, `useAnimatedRef`)
- Multiple animated values need to be derived from a single source of truth

```tsx
const offset = useSharedValue(0);

const animatedStyle = useAnimatedStyle(() => ({
  transform: [{ translateX: withSpring(offset.value) }],
}));

// Gesture-driven example
const gesture = Gesture.Pan().onUpdate((e) => {
  offset.value = e.translationX;
});
```

Avoid reading `sharedValue.value` on the JS thread inside React render or event handlers — it causes a synchronization that blocks the JS thread. Derive values from shared values using `useDerivedValue` instead.

---

## Animating Text

Avoid updating `Animated.Text` content by changing state — it triggers a full React re-render for every frame.

For animated numeric counters or any frequently-changing text, use `AnimatedTextInput` with `animatedProps`:

```tsx
import Animated, { useAnimatedProps } from 'react-native-reanimated';
import { TextInput } from 'react-native';

const AnimatedTextInput = Animated.createAnimatedComponent(TextInput);

function Counter({ progress }: { progress: SharedValue<number> }) {
  const animatedProps = useAnimatedProps(() => ({
    text: String(Math.round(progress.value)),
    defaultValue: '0',
  }));

  return (
    <AnimatedTextInput
      animatedProps={animatedProps}
      editable={false}
      style={styles.counter}
    />
  );
}
```

This updates the native text node directly on the UI thread, bypassing React and eliminating re-renders.

---

## Infinite Animations

CSS animations with `animationIterationCount: 'infinite'` clean up automatically on unmount.

For shared value infinite animations, always cancel them in the `useEffect` cleanup:

```tsx
useEffect(() => {
  offset.value = withRepeat(withTiming(1, { duration: 800 }), -1, true);

  return () => {
    cancelAnimation(offset);
  };
}, []);
```

Never start infinite animations outside the component lifecycle (module scope, global timers). They cannot be cleaned up and will leak.

---

## Prefer Non-Layout Properties

Animating layout properties (`top`, `left`, `width`, `height`, `margin`, `padding`) forces a layout pass on every frame, which is expensive and causes jank.

Prefer:
- `transform: [{ translateX }, { translateY }, { scale }, { rotate }]`
- `opacity`
- `backgroundColor`

If a design requires what looks like a size change, consider `scale` transforms — same visual effect without triggering layout.

---

## Supported Style Properties

Most React Native style properties are animatable. Key exceptions and platform notes:

- **`flexBasis`**: Changes are calculated but never applied to the view. Use `flexGrow`/`flexShrink` instead.
- **Shadow properties**: `shadowOffset`, `shadowOpacity`, `shadowRadius` do not work on Android. Use `boxShadow` instead (works on all platforms).
- **Web shadows**: All shadow styles must be specified in every keyframe on Web, or they are lost.
- **`tintColor` on iOS**: Must be present in the initial style when the `Image` component mounts. Adding it later has no effect.
- **Style inheritance**: Not supported. Properties that normally inherit in CSS (e.g., `textDecorationColor` from `color`) must be set explicitly.
- **Mixed-unit margins**: Interpolating between absolute and percentage margins may produce unexpected results when the parent's dimensions are affected by the child's margins.

---

## Threading: scheduleOnRN instead of runOnJS

`runOnJS` is removed in Reanimated 4. Use `scheduleOnRN` to call JS-thread functions from a worklet. Arguments are passed directly, not curried:

```tsx
// Reanimated 3 (removed)
runOnJS(setCount)(newCount);

// Reanimated 4
scheduleOnRN(setCount, newCount);
```

`scheduleOnRN` schedules the call asynchronously on the React Native runtime. Functions passed to `scheduleOnRN` must be defined in JS thread scope (they cannot be created inside worklets or animation callbacks).
