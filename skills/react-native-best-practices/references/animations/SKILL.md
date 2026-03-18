---
name: animations
description: "Production animation patterns for React Native using Reanimated 4, WebGPU, and TypeGPU. Covers CSS transitions, CSS animations, shared value animations, GPU shader animations, layout animations, scroll-driven animations, interpolation, particle systems, procedural noise, SDF rendering, performance tuning, and accessibility. Trigger on: Reanimated, useSharedValue, useAnimatedStyle, withSpring, withTiming, withDecay, withRepeat, withSequence, CSS transition, CSS animation, layout animation, FadeIn, SlideIn, ZoomIn, LinearTransition, keyframe, interpolate, scrollTo, useFrameCallback, WebGPU, react-native-wgpu, TypeGPU, GPU shader, WGSL, particle system, Perlin noise, SDF, Three.js, react-three-fiber, animation performance, or any request to animate UI in React Native."
---

# Animations

Software Mansion's production animation patterns for React Native on Reanimated 4 and the New Architecture.

Load at most one reference file per question.

## References

| File | When to read |
|------|-------------|
| `animations.md` | Choosing between CSS transitions, CSS animations, and shared value animations; CSS transition and CSS animation property details; animating text; infinite animation cleanup; `scheduleOnRN` |
| `animation-functions.md` | Core hooks (`useSharedValue`, `useAnimatedStyle`, `useAnimatedProps`, `useDerivedValue`, `createAnimatedComponent`, `cancelAnimation`, `useAnimatedRef`); animation functions (`withTiming`, `withSpring`, `withDecay`); animation modifiers (`withDelay`, `withRepeat`, `withSequence`, `withClamp`); composing animations |
| `layout-animations.md` | Entering/exiting animations (`FadeIn`, `SlideIn`, `BounceIn`, `ZoomIn`, etc.); layout transitions (`LinearTransition`, `FadingTransition`, `SequencedTransition`); keyframe animations; list item animations (`itemLayoutAnimation`); shared element transitions |
| `scroll-and-events.md` | Scroll-driven animations; `useAnimatedScrollHandler`; `scrollTo`; `useScrollOffset`; `interpolate`; `interpolateColor`; `clamp`; `useAnimatedReaction`; `useFrameCallback`; `measure`; `setNativeProps`; device sensors (`useAnimatedSensor`) |
| `gpu-animations.md` | GPU shader animations; `react-native-wgpu` Canvas and device setup; TypeGPU typed pipelines; Reanimated + WebGPU worklet integration; compute pipelines for particle systems, physics, and simulations; `@typegpu/noise` (Perlin noise, PRNG); `@typegpu/sdf` (signed distance shapes); Three.js / React Three Fiber for 3D |
| `animations-performance.md` | Performance tuning; 120fps setup; feature flags; FPS drop fixes; simultaneous animation limits; accessibility (`useReducedMotion`, `ReducedMotionConfig`); worklet closure optimization; debug vs release builds |
