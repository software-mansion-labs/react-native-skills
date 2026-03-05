# Gesture Composition

All gestures wrapped in `useMemo`. Compositions also wrapped in `useMemo`.

## Same-Component

```tsx
// Both active at once
const composed = useMemo(() => Gesture.Simultaneous(pan, pinch), [pan, pinch]);

// First to activate wins
const composed = useMemo(() => Gesture.Race(swipeLeft, swipeRight), [swipeLeft, swipeRight]);

// Priority order — first arg wins, others activate only after it fails
const composed = useMemo(() => Gesture.Exclusive(doubleTap, singleTap), [doubleTap, singleTap]);
```

## Cross-Component Relations

Both gestures must share the same `GestureHandlerRootView`. Use builder methods:

```tsx
// Won't cancel each other
.simultaneousWithExternalGesture(otherGesture)

// Wait for other to fail first
.requireExternalGestureToFail(otherGesture)

// Block other while active
.blocksExternalGesture(otherGesture)
```

## Pan Inside ScrollView

```tsx
const scrollRef = useRef(null);
const nativeScroll = useMemo(() => Gesture.Native().withRef(scrollRef), []);
const pan = useMemo(() =>
  Gesture.Pan()
    .simultaneousWithExternalGesture(nativeScroll)
    .onUpdate((e) => { ... }),
[nativeScroll]);
const composed = useMemo(() => Gesture.Simultaneous(pan, nativeScroll), [pan, nativeScroll]);

<ScrollView ref={scrollRef}>
  <GestureDetector gesture={composed}><Animated.View /></GestureDetector>
</ScrollView>
```

## Modals

Gestures inside a modal need their own `GestureHandlerRootView`:

```tsx
<Modal>
  <GestureHandlerRootView style={{ flex: 1 }}>
    {/* gestures here */}
  </GestureHandlerRootView>
</Modal>
```
