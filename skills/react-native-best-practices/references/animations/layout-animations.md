# Layout Animations

Animations for components entering, exiting, or changing position in the view hierarchy. Reanimated 4, New Architecture required.

---

## Entering and Exiting Animations

Animate elements when they are added to or removed from the view hierarchy:

```tsx
import Animated, { FadeIn, FadeOut } from 'react-native-reanimated';

{visible && (
  <Animated.View entering={FadeIn} exiting={FadeOut}>
    <Text>Hello</Text>
  </Animated.View>
)}
```

### Predefined Animation Families

Each family has directional variants (e.g., `FadeInRight`, `FadeInLeft`, `FadeInUp`, `FadeInDown`).

| Family | Entering | Exiting | Default Duration |
|--------|----------|---------|-----------------|
| **Fade** | `FadeIn`, `FadeInRight/Left/Up/Down` | `FadeOut`, `FadeOutRight/Left/Up/Down` | 300ms |
| **Slide** | `SlideInRight/Left/Up/Down` | `SlideOutRight/Left/Up/Down` | 300ms |
| **Zoom** | `ZoomIn`, `ZoomInDown/Up/Left/Right`, `ZoomInEasyDown/Up`, `ZoomInRotate` | `ZoomOut` + same variants | 300ms |
| **Bounce** | `BounceIn`, `BounceInRight/Left/Up/Down` | `BounceOut` + same variants | 600ms |
| **Flip** | `FlipInEasyX/Y`, `FlipInXDown/Up`, `FlipInYLeft/Right` | `FlipOutEasyX/Y` + same variants | 300ms |
| **Stretch** | `StretchInX/Y` | `StretchOutX/Y` | 300ms |
| **Roll** | `RollInRight/Left` | `RollOutRight/Left` | 300ms |
| **Rotate** | `RotateInDownLeft/Right`, `RotateInUpLeft/Right` | same + `Out` | 300ms |
| **LightSpeed** | `LightSpeedInRight/Left` | `LightSpeedOutRight/Left` | 300ms |
| **Pinwheel** | `PinwheelIn` | `PinwheelOut` | 300ms |

### Modifiers

Chain modifiers on any predefined animation:

```tsx
entering={FadeIn.duration(500).delay(200).springify().damping(15)}
```

**Time-based** (incompatible with `.springify()`):
- `.duration(ms)` — animation length
- `.easing(fn)` — defaults to `Easing.inOut(Easing.quad)`

**Spring-based** (after calling `.springify()`):
- `.damping(value)` — default 10 (higher = faster to rest)
- `.mass(value)` — default 1 (lower = faster)
- `.stiffness(value)` — default 100 (higher = less bouncy)
- `.overshootClamping(boolean)` — default false

**Common:**
- `.delay(ms)` — delay before start
- `.randomDelay()` — random delay between 0 and provided delay (default 1000ms)
- `.reduceMotion(ReduceMotion)` — accessibility
- `.withInitialValues(styleProps)` — override default initial config
- `.withCallback((finished) => {})` — fires when animation ends

### Gotchas

- **`nativeID` conflict (New Architecture)**: Reanimated uses `nativeID` internally for entering animations. Overwriting it breaks the animation. Wrap animated children in a plain `View` to work around this, especially with `TouchableWithoutFeedback`.
- **View flattening**: Removing a non-animated parent triggers exiting animations in its children, but the parent will not wait for children to finish. Add `collapsable={false}` to the parent to prevent this.
- **Spring-based animations**: Not yet available on the web platform.
- **Performance**: Define animation builders outside of components or wrap with `useMemo`.

---

## Layout Transitions

Smooth animations when a component's position or size changes due to state updates:

```tsx
import Animated, { LinearTransition } from 'react-native-reanimated';

<Animated.View layout={LinearTransition}>
  {items.map((item) => (
    <Item key={item.id} {...item} />
  ))}
</Animated.View>
```

### Predefined Transitions

| Transition | Behavior | Default Duration |
|------------|----------|-----------------|
| `LinearTransition` | Position and size change uniformly | 300ms |
| `SequencedTransition` | x-position/width first, then y-position/height | 500ms |
| `FadingTransition` | Fades out at old position, fades in at new position | 500ms |
| `JumpingTransition` | Components "jump" to new position | 300ms |
| `CurvedTransition` | Per-dimension easing (`.easingX()`, `.easingY()`, `.easingWidth()`, `.easingHeight()`) | 300ms |
| `EntryExitTransition` | Combines entering/exiting animations (`.entering()`, `.exiting()`) | Sum of both |

The generic `Layout` transition from older Reanimated versions is deprecated. Use `LinearTransition`.

### Transition Modifiers

Same modifier system as entering/exiting: `.duration()`, `.delay()`, `.springify()`, `.damping()`, `.mass()`, `.stiffness()`, `.reduceMotion()`, `.withCallback()`.

`SequencedTransition` also supports `.reverse()` to animate y first, then x.

**Spring config modes**: Use either physics-based (`damping`/`stiffness`) or duration-based (`duration`/`dampingRatio`), never both. If both are provided, duration-based overrides.

---

## Keyframe Animations

For complex multi-step entering/exiting animations beyond what presets offer:

```tsx
import { Keyframe } from 'react-native-reanimated';

const enteringAnimation = new Keyframe({
  0: { opacity: 0, transform: [{ scale: 0.5 }, { rotate: '-45deg' }] },
  50: {
    opacity: 1,
    transform: [{ scale: 1.2 }, { rotate: '0deg' }],
    easing: Easing.out(Easing.quad),
  },
  100: { transform: [{ scale: 1 }, { rotate: '0deg' }] },
});

<Animated.View entering={enteringAnimation.duration(600)} />
```

### Rules

- Keyframe `0` (or `from`) is **required**. Provide initial values for all properties you intend to animate.
- Keyframe `100` (or `to`) is optional.
- Do not provide both `0` and `from`, or both `100` and `to`.
- Easing is assigned to the second keyframe in a pair. Never provide easing to keyframe `0`.
- Default easing between keyframes is `Easing.linear`.
- **All properties in the transform array must appear in the same order across all keyframes.**

### Modifiers

`.duration(ms)` (default 500), `.delay(ms)`, `.reduceMotion()`, `.withCallback()`.

---

## List Layout Animations

Animate item layout changes in `FlatList` when items are added, removed, or reordered:

```tsx
<Animated.FlatList
  data={data}
  renderItem={renderItem}
  itemLayoutAnimation={LinearTransition}
/>
```

### Rules

- Only works with single-column `FlatList`. `numColumns` cannot be greater than 1.
- Items must have a `key` or `id` property (or provide a custom `keyExtractor`).
- Set `itemLayoutAnimation` to `undefined` to disable at runtime.
- Use `.skipEnteringExitingAnimations` to prevent entering/exiting animations on initial mount and unmount of the FlatList.

---

## LayoutAnimationConfig

Skip entering/exiting animations for a subtree:

```tsx
import { LayoutAnimationConfig } from 'react-native-reanimated';

<LayoutAnimationConfig skipEntering skipExiting>
  {children}
</LayoutAnimationConfig>
```

Can be nested. For FlatLists, use the `.skipEnteringExitingAnimations` modifier on `itemLayoutAnimation` instead.

---

## Shared Element Transitions

**Status: Experimental. Not recommended for production.**

Animates a view between two screens during navigation:

```tsx
<Animated.Image
  sharedTransitionTag="hero-image"
  sharedTransitionStyle={SharedTransition.duration(550).springify()}
/>
```

- Requires React Navigation native stack navigator. Tab navigator and `transparentModal` (iOS) are not supported.
- Tags must be unique per screen. Add the same tag to matching components on both screens.
- Default duration: 500ms. Animates width, height, position, transform, backgroundColor, opacity.
- iOS supports progress-based (swipe gesture) transitions. Android uses timing-based transitions only.
- Custom animation functions are not yet supported.
