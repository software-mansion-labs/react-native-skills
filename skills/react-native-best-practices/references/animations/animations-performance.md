# Animations Performance

Reanimated 4 requires the New Architecture (Fabric). All guidance here assumes that.

---

## 120fps Support

Enable ProMotion display support on iOS by adding to `Info.plist`:

```xml
<key>CADisableMinimumFrameDurationOnPhone</key>
<true/>
```

Without this flag, iOS caps animations at 60fps even on ProMotion devices.

---

## Feature Flags

Reanimated 4 exposes feature flags to opt into fixes for known New Architecture issues. Enable them early in your app entry point, before any Reanimated code runs.

```tsx
import { configureReanimatedLogger } from 'react-native-reanimated';
```

The flags live in `ReanimatedFeatureFlags` (exact import varies by Reanimated patch version — check the changelog for your version).

### Flickering / Jittering While Scrolling

Animated components like sticky headers flicker during `FlatList` or `ScrollView` scrolling on the New Architecture.

**Fix:** Upgrade to React Native 0.81+ and enable:
- `preventShadowTreeCommitExhaustion` (experimental release-level flag in RN)
- `DISABLE_COMMIT_PAUSING_MECHANISM` (Reanimated feature flag)

### FPS Drops During Scrolling

FPS drops when many animated components are visible during scroll.

**Fix:** Upgrade to React Native 0.80+ and Reanimated 4.2.0+, then enable:
- `USE_COMMIT_HOOK_ONLY_FOR_REACT_COMMITS`

### Low FPS with Many Simultaneous Animations

**Fix:** Enable platform-specific synchronous UI update flags:
- `ANDROID_SYNCHRONOUSLY_UPDATE_UI_PROPS` (available since 4.0.0)
- `IOS_SYNCHRONOUSLY_UPDATE_UI_PROPS` (available since 4.2.0)

Note: these flags may interfere with touch detection on animated `transform` elements. Prefer `Pressable` from `react-native-gesture-handler` over the core `Pressable` when using these flags.

---

## Simultaneous Animation Limits

Reanimated can handle many animated components, but performance degrades at scale:

| Platform       | Practical limit |
|----------------|-----------------|
| iOS            | ~500 components |
| Low-end Android | ~100 components |

For lists with many animated items, consider reducing animation complexity on low-end devices using `useReducedMotion`.

---

## Avoid Reading Shared Values on the JS Thread

Reading `sv.value` inside React render, event handlers, or `useEffect` triggers a synchronization from the UI thread to the JS thread, which can block the JS thread.

Instead, use `useDerivedValue` to transform shared values and `useAnimatedStyle` to consume them — both run on the UI thread.

---

## Memoize Callbacks and Gesture Objects

Frame callbacks and gesture objects are re-created on every render by default. Wrap them:

```tsx
const frameCallback = useFrameCallback(
  useCallback((frameInfo) => {
    // runs on UI thread every frame
  }, [])
);

const gesture = useMemo(() =>
  Gesture.Pan().onUpdate((e) => {
    offset.value = e.translationX;
  }),
  []
);
```

---

## Debug vs. Release Builds

Always profile animations in a release build. Debug builds add significant JS overhead (Metro bundler, Hermes debug mode, dev warnings) that makes animations appear slower than they are in production. Use `npx react-native run-android --mode=release` or the equivalent iOS scheme.
