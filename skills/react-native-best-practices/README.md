# React Native Best Practices Skill

Production patterns for React Native apps on the New Architecture, by [Software Mansion](https://swmansion.com/).

Add this skill to give your AI coding agent accurate, current guidance for Software Mansion's React Native libraries: Reanimated, Gesture Handler, React Native SVG, ExecuTorch, Audio API, and more.

## Sub-skills

| Sub-skill | Covers | Status |
|-----------|--------|--------|
| [Animations](./references/animations/) | Reanimated 4, CSS transitions, CSS animations, shared values, layout animations, 120fps, performance flags | Complete |
| [Gestures](./references/gestures/) | Gesture Handler: tap, pan, pinch, rotation, swipe, long press, drag | Stub |
| [SVG](./references/svg/) | React Native SVG: vector graphics, icons, charts, illustrations | Stub |
| [Haptics](./references/haptics/) | Tactile feedback, vibration patterns, impact/notification haptics | Stub |
| [On-device AI](./references/on-device-ai/) | React Native ExecuTorch: LLMs, computer vision, OCR, speech, embeddings | Complete |
| [Multimedia](./references/multimedia/) | Video playback, camera, image picking, media recording and processing | Stub |
| [Rich Text](./references/rich-text/) | Formatted text, markdown, HTML rendering, text editors | Stub |
| [Multithreading](./references/multithreading/) | Worklets, background processing, offloading computation | Stub |
| [Audio](./references/audio/) | React Native Audio API: playback, recording, visualization, sessions | Complete |

**Complete** = full reference documentation with code examples. **Stub** = frontmatter and description only, reference content coming soon.

## Structure

```
react-native-best-practices/
├── SKILL.md                          # Entry point: routing table for sub-skills
└── references/
    ├── animations/
    │   ├── SKILL.md                  # When to use, what references to read
    │   ├── animations.md             # Decision tree, CSS transitions/animations, shared values
    │   └── animations-performance.md # 120fps, feature flags, simultaneous animation limits
    ├── gestures/SKILL.md
    ├── svg/SKILL.md
    ├── haptics/SKILL.md
    ├── on-device-ai/
    │   ├── SKILL.md                  # Use cases, capabilities overview, getting started
    │   └── references/
    │       ├── reference-llms.md     # LLM hooks, tool calling, structured output
    │       ├── reference-cv.md       # Classification, detection, segmentation
    │       ├── reference-cv-2.md     # Style transfer, text-to-image, image embeddings
    │       ├── reference-ocr.md      # Horizontal/vertical text recognition
    │       ├── reference-audio.md    # Speech-to-text, text-to-speech, VAD
    │       ├── reference-nlp.md      # Text embeddings, tokenization
    │       ├── reference-models.md   # Model catalog, loading strategies, device constraints
    │       └── core-utilities.md     # ResourceFetcher, error handling, custom models
    ├── rich-text/SKILL.md
    ├── multithreading/SKILL.md
    └── audio/SKILL.md                # AudioContext singletons, buffer state, visualizations
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
- [React Native ExecuTorch](https://docs.swmansion.com/react-native-executorch/)
- [React Native Audio API](https://docs.swmansion.com/react-native-audio-api/)
