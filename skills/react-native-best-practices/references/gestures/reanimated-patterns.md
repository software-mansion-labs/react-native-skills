# Reanimated Integration Patterns

All gestures wrapped in `useMemo`. `Animated` from `react-native-reanimated`, not `react-native`.

## Draggable Element

`translationX/Y` is cumulative from gesture start - save offset on `onBegin` to support multi-drag:

```tsx
const offsetX = useSharedValue(0);
const startX = useSharedValue(0);

const pan = useMemo(() =>
  Gesture.Pan()
    .onBegin(() => { startX.value = offsetX.value; })
    .onUpdate((e) => { offsetX.value = startX.value + e.translationX; })
    .onEnd(() => { offsetX.value = withSpring(0); }),
[]);
```

## Pinch to Zoom

`e.scale` resets to 1 each gesture — multiply by `savedScale` to accumulate:

```tsx
const scale = useSharedValue(1);
const savedScale = useSharedValue(1);

const pinch = useMemo(() =>
  Gesture.Pinch()
    .onUpdate((e) => { scale.value = savedScale.value * e.scale; })
    .onEnd(() => { savedScale.value = scale.value; }),
[]);
```

## Pinch + Pan (Photo Viewer)

```tsx
const pinch = useMemo(() => Gesture.Pinch().onUpdate(...).onEnd(...), []);
const pan = useMemo(() => Gesture.Pan().onBegin(...).onUpdate(...), []);
const composed = useMemo(() => Gesture.Simultaneous(pan, pinch), [pan, pinch]);
```

## runOnJS

Use `runOnJS` to call React state setters from gesture callbacks (which run on the UI thread):

```tsx
const tap = useMemo(() =>
  Gesture.Tap().onEnd(() => { runOnJS(setState)(value); }),
[]);
```
