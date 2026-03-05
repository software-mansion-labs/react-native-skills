# Animations

Production-quality animation patterns for React Native apps using Reanimated 4 on the New Architecture.

For performance tuning and feature flags, see **`animations-performance.md`**.

---

## Decision Tree

Pick the animation type based on what drives the animation and what it needs to compute.

```
Is the animation driven by a state change (not a gesture or continuous input)?
├── YES → Can it be expressed as a simple A→B property transition?
│   ├── YES → Use CSS Transition (transitionProperty)
│   └── NO  → Does it need a defined keyframe sequence?
│       ├── YES → Use CSS Animation (animationName + keyframes)
│       └── NO  → Use CSS Transition with multiple properties
└── NO  → Is it gesture-driven, or does it need math / trig / layout reads?
    └── YES → Use Shared Value Animation (useSharedValue + useAnimatedStyle)
```

Default to CSS transitions and CSS animations. They are declarative, easier to read, and remove the overhead of worklet execution. Reach for shared values only when the animation requires programmatic control that CSS cannot express.

---

## CSS Transitions

Use when a component's style should animate smoothly whenever a state-driven prop changes.

```tsx
<Animated.View
  style={{
    width: isExpanded ? 200 : 100,
    transitionProperty: 'width',
    transitionDuration: 300,
    transitionTimingFunction: Easing.out(Easing.quad),
  }}
/>
```

Animate multiple properties by passing arrays:

```tsx
transitionProperty: ['width', 'opacity', 'backgroundColor'],
transitionDuration: [300, 200, 150],
```

Avoid `transitionProperty: 'all'` — it forces evaluation of every style property on each frame and degrades performance.

CSS transitions cannot animate discrete properties like `flexDirection` or `justifyContent`. Use Layout Animations for those cases.

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
    animationTimingFunction: Easing.inOut(Easing.ease),
  }}
/>
```

Reanimated takes the current element state as the implicit `0%` keyframe, so you only need to define the frames that differ.

For infinite CSS animations, set `animationIterationCount: 'infinite'`. The animation is tied to the component's mount state — it stops automatically when the component unmounts, no manual cleanup needed.

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

## Layout Animations

Use `LinearTransition` for animating position and size changes in response to state updates. The generic `Layout` shorthand from older Reanimated versions is deprecated.

```tsx
import { LinearTransition } from 'react-native-reanimated';

<Animated.View layout={LinearTransition}>
  {items.map((item) => (
    <Item key={item.id} {...item} />
  ))}
</Animated.View>
```

Other available transitions: `FadingTransition`, `SequencedTransition`, `EntryExitTransition`.

---

## Prefer Non-Layout Properties

Animating layout properties (`top`, `left`, `width`, `height`) forces a layout pass on every frame, which is expensive and causes jank.

Prefer:
- `transform: [{ translateX }, { translateY }, { scale }, { rotate }]`
- `opacity`
- `backgroundColor`

If a design requires what looks like a size change, consider `scale` transforms — same visual effect without triggering layout.

---

## Threading: scheduleOnRN instead of runOnJS

`runOnJS` is removed in Reanimated 4. Use `scheduleOnRN` to call JS-thread functions from a worklet. Arguments are passed directly, not curried:

```tsx
// Reanimated 3 (removed)
runOnJS(setCount)(newCount);

// Reanimated 4
scheduleOnRN(setCount, newCount);
```

`scheduleOnRN` schedules the call asynchronously on the React Native runtime. For a synchronous return to the UI thread, use `runOnUISync`.
