# When to Use React Native SVG

## Decision Guide

```
Do you need to animate or interactively control parts of the SVG?
│
├── YES: Use react-native-svg + Reanimated
│
└── NO
    │
    ├── Is it an icon?
    │   └── YES: Use @expo/vector-icons
    │
    └── Is it a static SVG image?
        └── YES: Use expo-image
```

## Alternative Comparison

| Tool | Best for | Key limitation |
|---|---|---|
| `@expo/vector-icons` | Icon sets | limited to available icon fonts |
| `expo-image` | Static SVG images | loads async = can blink, no filter support |
| `react-native-skia` | Complex SVGs with filters | heavy `Canvas` objects |
| Lottie / Rive | Animated vector graphics | converting assets from SVG format |
| WebView | - | use only as last resort |
| `react-native-vector-image` | Static SVGs as native assets | build-time asset generation |

---
