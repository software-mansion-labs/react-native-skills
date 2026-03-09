# Tap Handling

## Decision

| What to use | When |
|---|---|
| RectButton | Inside ScrollView / FlatList |
| RectButton + GestureDetector(Gesture.Tap) | Inside ScrollView / FlatList + need opacity/scale animation on press |
| Gesture.Tap() + useMemo | Custom animation, multi-tap, double-tap |
| Pressable (from RNGH) | Outside scroll containers, simple case |
| RN Pressable / Touchables | Avoid — conflicts with RNGH, causes double-tap bugs |

**RNGH Pressable vs RectButton in lists:** Pressable highlights items immediately on scroll start (bad native feel). RectButton doesn't - prefer RectButton inside scroll containers.

**Don't mix RN touch + RNGH** - causes double-tap bugs. Pick one per app.

## RectButton in Scroll Containers

Import both `RectButton` and the scroll container from `react-native-gesture-handler` - mixing RNGH buttons with React Native's `ScrollView`/`FlatList` causes gesture conflicts:

```tsx
import { FlatList, RectButton, BorderlessButton } from 'react-native-gesture-handler';

<RectButton onPress={handlePress} style={styles.row}>
  <Text>{item.title}</Text>
  <BorderlessButton onPress={handleDelete}>
    <DeleteIcon />
  </BorderlessButton>
</RectButton>
```

## RectButton + UI Thread Animation

RectButton alone can't run worklets on the UI thread. Wrap it with GestureDetector + Gesture.Tap. 

```tsx
const opacity = useSharedValue(1);
const tap = useMemo(() =>
  Gesture.Tap()
    .onBegin(() => { opacity.value = withTiming(0.7); })
    .onFinalize(() => { opacity.value = withTiming(1); }),
[]);

<GestureDetector gesture={tap}>
  <Animated.View style={{ opacity }}>
    <RectButton onPress={handlePress}>
      <Text>{item.title}</Text>
    </RectButton>
  </Animated.View>
</GestureDetector>
```

## Gesture.Tap() - Custom Feedback / Multi-tap

```tsx
const tap = useMemo(() =>
  Gesture.Tap()
    .onBegin(() => { scale.value = withTiming(0.95); })
    .onFinalize(() => { scale.value = withTiming(1); }),
[]);
```

Double-tap with single-tap fallback - `Gesture.Exclusive` gives doubleTap priority:

```tsx
const doubleTap = useMemo(() => Gesture.Tap().numberOfTaps(2).onEnd(() => { runOnJS(onDouble)(); }), []);
const singleTap = useMemo(() => Gesture.Tap().numberOfTaps(1).onEnd(() => { runOnJS(onSingle)(); }), []);
const composed = useMemo(() => Gesture.Exclusive(doubleTap, singleTap), [doubleTap, singleTap]);
```

`hitSlop` - increases touch target area - use on small/icon components.
