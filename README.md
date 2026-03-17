# React Native Skills

> **Work in progress.** This repository is actively being developed. Some sub-skills are stubs and more content is on the way.

Production-ready patterns for React Native development, packaged as [Claude Code skills](https://docs.anthropic.com/en/docs/claude-code/skills). Maintained by [Software Mansion](https://swmansion.com/).

Add a skill to your project and your AI coding agent gets up-to-date guidance for animations, gestures, on-device AI, audio, and other React Native features. Works with Claude Code, Cursor, Windsurf, and other tools that support the skill format.

## Available Skills

### [react-native-best-practices](./skills/react-native-best-practices/)

Production patterns for React Native apps on the New Architecture, covering:

| Topic | What it covers |
|-------|---------------|
| **Animations** | CSS transitions, CSS animations, shared value animations, layout animations, Reanimated 4, 120fps, performance tuning |
| **Gestures** | Tap, pan, pinch, swipe, long press, drag with Gesture Handler |
| **SVG** | Vector graphics, icons, charts, illustrations with React Native SVG |
| **On-device AI** | LLMs, computer vision, OCR, audio processing, embeddings with React Native ExecuTorch |
| **Multimedia** | Video playback and streaming |
| **Rich Text** | Rich text editing with react-native-enriched: formatting toolbar, mentions, links |
| **Multithreading** | Background processing, Worklets, offloading computation from the JS thread |
| **Audio** | Playback, recording, visualization, session management with React Native Audio API |

### [radon-mcp](./skills/radon-mcp/)

Best practices for using Radon IDE's MCP tools when developing, debugging, and inspecting React Native and Expo apps. Covers viewing screenshots, reading logs, inspecting the component tree, debugging network requests, reloading the app, and querying React Native documentation.

## Repository Structure

```
react-native-skills/
└── skills/
    ├── react-native-best-practices/
    │   ├── SKILL.md                        # Main skill entry point
    │   └── references/
    │       ├── animations/
    │       │   ├── SKILL.md                # Animation sub-skill
    │       │   ├── animations.md           # Core animation patterns
    │       │   └── animations-performance.md
    │       ├── gestures
    │       │   ├── SKILL.md                # Gestures sub-skill
    │       │   ├── gestures-composition.md
    │       │   ├── reanimated-patterns.md  # Reanimated+Gestures integration
    │       │   └── tap-handling.md
    │       ├── svg/
    │       │   ├── SKILL.md                # SVG sub-skill
    │       │   ├── svg.md                  # Setup, issues and performance
    │       │   ├── animation-patterns.md   # SVG animation patterns
    │       │   └── when-to-use.md
    │       ├── haptics/SKILL.md
    │       ├── on-device-ai/
    │       │   ├── SKILL.md                # On-device AI sub-skill
    │       │   └── references/             # Detailed API references
    │       ├── multimedia/SKILL.md
    │       ├── rich-text/SKILL.md
    │       ├── multithreading/SKILL.md
    │       └── audio/SKILL.md
    └── radon-mcp/
        ├── SKILL.md
        └── references/
            ├── view-application-logs.md
            ├── view-screenshot.md
            ├── view-component-tree.md
            ├── view-network-logs.md
            ├── view-network-request-details.md
            ├── reload-application.md
            ├── query-documentation.md
            └── get-library-description.md
```

The top-level `SKILL.md` acts as a table of contents. Reference files load only when relevant to the current task, keeping the context window focused.

## Contributing

### Adding a new sub-skill

1. Create a directory under `skills/react-native-best-practices/references/<topic>/`
2. Add a `SKILL.md` with YAML frontmatter (`name`, `description`) and a references table
3. Place detailed documentation in sibling `.md` files
4. Register the sub-skill in `skills/react-native-best-practices/SKILL.md`

### Writing guidelines

- Keep `SKILL.md` files under 500 lines. Split longer content into reference files.
- Write descriptions that include trigger terms (what phrases should activate this skill).
- Use concrete code examples over abstract explanations.
- Target the New Architecture and current library versions.

## License

MIT
