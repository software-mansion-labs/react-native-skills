# React Native Best Practices Skill

Production patterns for React Native apps on the New Architecture, by [Software Mansion](https://swmansion.com/).

Add this skill to give your AI coding agent accurate, current guidance for Software Mansion's React Native libraries: Reanimated, Gesture Handler, React Native SVG, ExecuTorch, Audio API, and more.

## Sub-skills

| Sub-skill | Covers | Status |
|-----------|--------|--------|
| [Animations](./references/animations/) | Reanimated 4, CSS transitions, CSS animations, shared values, GPU shader animations (WebGPU, TypeGPU), layout animations, scroll-driven animations, SVG animations, 120fps, performance flags | Complete |
| [Gestures](./references/gestures/) | Gesture Handler: tap, pan, pinch, rotation, fling, hover, long press, Pressable, RectButton, Swipeable, DrawerLayout, gesture composition, testing | Complete |
| [SVG](./references/svg/) | React Native SVG: when to use vs expo-image/Skia/Lottie/Rive/fonts/WebView, installation, loading (URI/XML/file), touch events, filters, FilterImage, performance pitfalls, iOS memory leaks | Complete |
| [On-device AI](./references/on-device-ai/) | React Native ExecuTorch: LLMs, computer vision, OCR, speech, text/image embeddings, best practices, with pointers to official docs | Complete |
| [Rich Text](./references/rich-text/) | Rich text editing with react-native-enriched, Markdown rendering with react-native-enriched-markdown | Complete |
| [Multithreading](./references/multithreading/) | react-native-worklets: Worker Runtimes, scheduling APIs, shared memory, Serializable, Synchronizable | Complete |
| [Audio](./references/audio/) | React Native Audio API: playback, recording, visualization, audio sessions | Complete |
| [Multimedia](./references/multimedia/) | Video playback and streaming | Stub |

**Complete** = full reference documentation with code examples. **Stub** = frontmatter and description only, reference content coming soon.

## Structure

```
react-native-best-practices/
├── SKILL.md                              # Entry point: routing table for sub-skills
└── references/
    ├── animations/
    │   ├── SKILL.md                      # When to use, what references to read
    │   ├── animations.md                 # Decision tree, CSS transitions/animations, shared values
    │   ├── animation-functions.md        # Core hooks, withSpring, withTiming, withDecay, composition
    │   ├── layout-animations.md          # Entering/exiting, transitions, keyframes
    │   ├── scroll-and-events.md          # Scroll-driven animations, useAnimatedReaction, useFrameCallback
    │   ├── gpu-animations.md             # Shader animations, react-native-wgpu, TypeGPU, particles
    │   ├── svg-animations.md              # Animating SVG elements and paths with Reanimated
    │   └── animations-performance.md     # 120fps, feature flags, simultaneous animation limits
    ├── gestures/
    │   ├── SKILL.md                      # Version decision tree (v2 Builder vs v3 Hook API)
    │   ├── gestures.md                   # Decision tree, lifecycle, threading, SharedValue config
    │   ├── tap-handling.md               # RectButton, Pressable, tap, double-tap, hit slop
    │   ├── continuous-gestures.md        # Pan, Pinch, Rotation, LongPress, Fling, Hover
    │   ├── gesture-composition.md        # Simultaneous, Race, Exclusive, VirtualGestureDetector
    │   ├── swipeable-and-drawer.md       # ReanimatedSwipeable, ReanimatedDrawerLayout
    │   └── testing.md                    # Jest setup, fireGestureHandler, troubleshooting
    ├── svg/
    │   ├── SKILL.md                      # When to use react-native-svg vs alternatives
    │   ├── when-to-use.md               # Choosing between svg, expo-image, icons, Skia, Lottie
    │   └── svg.md                        # Installation, loading (URI/XML/file), touch events, filters, FilterImage, performance
    ├── on-device-ai/
    │   ├── SKILL.md                      # Doc map with webfetch URLs for each feature
    │   └── best-practices.md             # Model loading, device constraints, error handling, processing tips
    ├── rich-text/
    │   └── SKILL.md                      # Editor and renderer patterns, style customization
    ├── multithreading/
    │   ├── SKILL.md                      # Runtime model, API decision tree, critical rules
    │   ├── threading-api.md              # Scheduling APIs, Worker Runtimes, sync/async
    │   ├── shared-memory.md              # Closures, Serializable, Synchronizable
    │   └── setup-and-advanced.md         # Installation, Babel config, Bundle Mode, Jest
    ├── audio/
    │   └── SKILL.md                      # AudioContext singletons, buffer state, visualizations, sessions
    └── multimedia/
        └── SKILL.md                      # Stub: video playback and streaming (coming soon)
```

## Adding a Sub-skill

1. **Create the directory**: `references/<your-topic>/`
2. **Write `SKILL.md`** with frontmatter:
   ```yaml
   ---
   name: your-topic
   description: "What it covers and when to trigger. Include specific keywords users might type."
   ---
   ```
3. **Add reference files** for detailed patterns and code examples
4. **Register it** in the parent `SKILL.md` sub-skills table
5. **Keep SKILL.md under 500 lines**. Move detailed content to reference files and link to them with a "when to read" table.

## Libraries Covered

- [React Native Reanimated](https://docs.swmansion.com/react-native-reanimated/) (v4, New Architecture)
- [React Native Gesture Handler](https://docs.swmansion.com/react-native-gesture-handler/)
- [React Native SVG](https://github.com/software-mansion/react-native-svg)
- [React Native Worklets](https://docs.swmansion.com/react-native-worklets/)
- [React Native ExecuTorch](https://docs.swmansion.com/react-native-executorch/)
- [React Native Audio API](https://docs.swmansion.com/react-native-audio-api/)
- [React Native Enriched](https://github.com/software-mansion/react-native-enriched)
- [React Native WGPU](https://github.com/software-mansion/react-native-wgpu) / [TypeGPU](https://docs.swmansion.com/TypeGPU/)
