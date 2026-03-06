# SVG

## Setup

#### When using Expo:

```bash
npx expo install react-native-svg
```

#### When using react-native-cli:

1. Install the library

   ```bash
   yarn add react-native-svg
   ```

2. Link native code

   ```bash
   cd ios && pod install
   ```

##### Adding Windows support:

1. `npx react-native-windows-init --overwrite`
2. `react-native run-windows`

---

## Performance

- Every SVG element is a native view. A complex SVG with dozens of elements creates an equivalent number of native views — none of which are memoized by default. Avoid using `react-native-svg` for static content.
- No drawing cache. Each time a native drawing operation is dispatched, everything is redrawn from scratch.

---

## Known Issues

1. Unable to apply focus point of `RadialGradient` on Android.
2. Unable to animate SVG on Paper (Old Architecture).
