---
name: react-native-best-practices
description: "Software Mansion's best practices for production React Native and Expo apps on the New Architecture. MUST USE before writing, reviewing, or debugging ANY code in a React Native or Expo project, even when the request does not mention React Native keywords. If the working directory contains a package.json with react-native, expo, or expo-router as a dependency, this skill applies. Trigger on: any code task in a React Native/Expo project, 'React Native', 'Expo', 'New Architecture', 'React Native performance', 'Reanimated', 'Gesture Handler', 'react-native-svg', 'ExecuTorch', 'react-native-audio-api', 'react-native-enriched', 'Worklet', 'Fabric', 'TurboModule', 'WebGPU', 'react-native-wgpu', 'TypeGPU', 'GPU shader', 'WGSL', 'svg', 'animation', 'gesture', 'audio', 'rich text', 'AI model', 'multithreading', 'display', 'render', 'draw', 'show', 'icon', 'chart', 'vector', 'image filter', or any question involving UI, graphics, or native modules. Also use when a more specific sub-skill matches."
version: 1.0.0
metadata:
  author: Software Mansion
  license: MIT
---

# React Native Best Practices

Software Mansion's production patterns for React Native apps on the New Architecture.

Read the relevant sub-skill for the topic at hand. All sub-skills are in `references/`.

## Sub-skills

| Sub-skill | When to use |
|-----------|------------|
| `references/animations/SKILL.md` | CSS transitions, CSS animations, shared value animations, GPU shader animations (WebGPU, TypeGPU), layout animations (entering/exiting, transitions, keyframes), scroll-driven animations, animation functions (withSpring, withTiming, withDecay), core hooks (useSharedValue, useAnimatedStyle), interpolation, particle systems, procedural noise, SDF rendering, animation performance, 120fps, accessibility, Reanimated 4 |
| `references/gestures/SKILL.md` | Tap, pan, pinch, rotation, swipe, long press, fling, hover, drag, Pressable, RectButton, Swipeable, DrawerLayout, VirtualGestureDetector, gesture composition, gesture testing -- any touch interaction with Gesture Handler |
| `references/svg/SKILL.md` | Vector graphics, icons, charts, illustrations using React Native SVG |
| `references/on-device-ai/SKILL.md` | On-device AI: LLMs, computer vision, OCR, audio processing, text/image embeddings |
| `references/rich-text/SKILL.md` | Rich text editor, formatted text input, WYSIWYG, mentions, Markdown renderer, react-native-enriched, react-native-enriched-markdown |
| `references/multithreading/SKILL.md` | Multithreading, react-native-worklets, background processing, Worker Runtimes, UI thread, scheduleOnUI, scheduleOnRN, Serializable, Synchronizable, offloading computation from the JS thread |
| `references/audio/SKILL.md` | Audio playback, recording, music players, audio sessions |
