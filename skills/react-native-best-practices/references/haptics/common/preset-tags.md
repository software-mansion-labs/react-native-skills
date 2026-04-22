## Preset Tags

Every preset is tagged with exactly four attributes — one from each category — that describe its haptic character. Use tags to find alternatives within a similar feel, filter by intensity or duration, or understand what a preset will feel like before playing it.

### Intensity
The overall weight and forcefulness of the haptic.

| Tag | Meaning | Usage |
|---|---|---|
| `Gentle` | Haptics with low amplitude — barely perceptible vibration that stays in the background. | Good for subtle confirmations, hover effects, or any interaction where the feedback should not interrupt the user. Common situations: cursor hover, background sync, silent notification, passive scroll boundary nudge, accessibility hints. |
| `Substantial` | Haptics with medium amplitude — a clear, balanced vibration that is easy to notice without being disruptive. | Ideal for most standard UI interactions such as button presses, toggles, or form confirmations. Common situations: button tap, toggle switch, form validation, menu item selection, pull-to-refresh trigger. |
| `Bold` | Haptics with high amplitude — a strong, attention-commanding vibration that is immediately felt. | Best for critical alerts, error states, or high-impact moments that require the user's immediate attention. Common situations: payment failure, security alert, authentication error, destructive action confirmation, game hit or collision. |

### Sharpness
How the haptic energy is shaped at the tactile level.

| Tag | Meaning | Usage |
|---|---|---|
| `Soft` | Haptics with low frequency — a smooth, rounded vibration with a gentle, cushioned feel. | Perfect for calm, ambient feedback, wellness interactions, or any context where a soft touch is preferred. Common situations: sleep reminder, gentle onboarding hint, ambient soundscape interaction. |
| `Flexible` | Haptics with medium frequency — a balanced vibration that sits between soft and rigid. | Suitable for general-purpose UI feedback, notifications, and interactions that need a neutral tactile character. Common situations: incoming message, push notification, content scroll snap, date picker tick, standard in-app alert. |
| `Rigid` | Haptics with high frequency — a crisp, precise vibration with a sharp, mechanical feel. | Great for snappy UI elements, keyboard-like taps, or any interaction that should feel precise and definitive. Common situations: virtual keyboard key press, numeric keypad tap, rotary dial click, picker wheel snap, PIN entry digit confirmation. |

### Shape
The envelope or waveform pattern of the haptic over time.

| Tag | Meaning | Usage |
|---|---|---|
| `Peak` | Haptics with a single amplitude peak — a quick rise and fall in intensity. | Ideal for single-event feedback, such as button presses or selection confirmations. Common situations: like or heart button tap, photo shutter release, item selection confirmation, swipe action completion, quick reply send. |
| `Ramp` | Haptics with a ramp-shaped amplitude pattern — amplitude increases or decreases linearly over the duration. | Suited for fade-in or fade-out effects, swipe feedback, or any interaction that should feel like a gradual build or release. Common situations: volume or brightness slider, swipe-to-dismiss gesture, pinch zoom, countdown timer nearing zero, pull-down refresh building tension. |
| `Saw` | Haptics with a sawtooth-shaped amplitude pattern — a sharp rise followed by an abrupt drop, or vice versa. | Effective for mechanical-feeling interactions, ratchet effects, or any feedback with a sharp, asymmetric edge. Common situations: ratchet scroll, slot machine reel spin, rotary dial simulation, drag-and-drop snap into position, file shredding animation. |
| `Impulses` | Haptics with a discrete pattern — short, distinct pulses separated by silence. | Useful for click-like feedback, Morse-style cues, or sequences of distinct tactile events. Common situations: step counter tick, quantity increment, metronome cue, item added to cart, typing indicator in chat. |
| `Bumps` | Haptics with multiple amplitude peaks — a series of rhythmic rises and falls. | Perfect for multi-step feedback, rhythmic notifications, or interactions that need a repeating pulse feel. Common situations: in-app achievement or badge unlock, multi-item batch selection, dice roll, streak milestone reached, coin or reward collection. |
| `Pattern` | Haptics with a custom, often repeating amplitude pattern — a structured sequence that defines a unique rhythm. | Best for branded feedback signatures, complex notifications, or effects that carry a recognizable tactile identity. Common situations: branded notification signature, custom incoming call alert, game character footstep cycle, sound-to-haptic mapping, recurring rhythm in a music app. |
| `Solid` | Haptics with a long continuous pattern at a constant amplitude — a steady, uniform vibration. | Good for indicating ongoing processes, loading states, or sustained alerts that need consistent presence. Common situations: file upload or download in progress, active voice recording, hold-to-confirm gesture, persistent alarm, live activity tracking. |

### Duration
The rough playback length of the preset.

| Tag | Meaning | Usage |
|---|---|---|
| `Impulse` | Extremely short haptic lasting less than 100 ms — an instantaneous tactile click. | Best for keyboard taps, quick confirmations, or any interaction that requires an instant, minimal response. Common situations: virtual keyboard key press, quick tap micro-interaction, toggle switch flip, checkbox tick, cursor click simulation. |
| `Short` | Brief haptic lasting between 100 ms and 250 ms — long enough to be clearly felt without lingering. | Ideal for button presses, toggle switches, and standard UI element interactions. Common situations: primary action button press, navigation tab switch, swipe gesture acknowledgement, photo filter selection, card flip or reveal. |
| `Extended` | Medium-length haptic lasting between 250 ms and 600 ms — provides richer, more expressive feedback. | Good for notifications, multi-step confirmations, or interactions that benefit from a more deliberate tactile moment. Common situations: incoming push notification, payment or purchase success, pull-to-refresh completion, form submission success, app rating prompt. |
| `Long` | Prolonged haptic lasting 600 ms or more — creates an immersive, sustained tactile experience. | Best for complex animations, gaming effects, elaborate alerts, or experiences where the haptic plays a central role. Common situations: in-app achievement or level completion, onboarding celebration, game victory screen, subscription or reward unlock, end-of-session summary celebration. |

### Using tags to choose a preset

Tags let you browse by feel. For example:
- **Gentle + Soft + Peak + Impulse** — the quietest, most invisible feedback possible.
- **Bold + Rigid + Impulses + Short** — sharp, decisive single or double taps.
- **Bold + Rigid + Pattern + Extended** — relentless, hard, rhythmic — error / alarm territory.
- **Substantial + Flexible + Bumps + Long** — natural-feeling, mid-weight, rolling sequence.

---

### Confirmations & Completions

**Preferred spec:** `Substantial` + `Impulses` or `Peak` + `Short` — the moment should feel decisive but proportionate. Brief and high-contrast.

**Avoid:** `Gentle` when the confirmation must register consciously; `Extended` duration on a single-tap event — it will feel sluggish rather than satisfying.

**Context you need before choosing:**
- *Emotional register* — warm/organic completion (personal save, wellness check-in) → `Soft` texture; crisp system confirmation (form submit, payment) → `Rigid` texture.
- *Reversibility* — irreversible or destructive action → `Bold` + `Rigid`; routine save or dismiss → `Substantial`.
- *Who triggered it* — user-initiated completion → `Impulses`/`Peak`/`Short`; system-completed task → `Bumps`/`Long`.

### Achievements & Celebrations

**Preferred spec:** `Bold` + `Extended` for major milestones; `Substantial` + `Long` for minor rewards.

**Avoid:** `Gentle` for any celebratory moment — it won't feel like a reward. `Short` or `Impulse` duration for a major achievement — it will feel like a tap, not a fanfare.

**Context you need before choosing:**
- *Milestone weight* — minor reward (streak day, points earned) → `Substantial`/`Long`; major achievement (level-up, rank, unlock) → `Bold`/`Extended`.
- *Emotional tone* — warm/personal achievement → `Soft` or `Flexible` texture; triumphant/game-like → `Flexible` or `Rigid` with `Ramp` shape.

### Errors & Rejections

**Preferred spec:** `Bold` + `Rigid` + `Saw` + `Short` or `Extended` — errors should feel sharp and interruptive. The `Saw` shape (sharp rise, rapid decay) communicates rejection clearly.

**Avoid:** `Gentle` or `Soft` for any error — the signal will not read as a rejection. `Peak` shape (smooth arc) softens the error and loses its urgency.

**Context you need before choosing:**
- *Severity* — soft validation hint → `Substantial`/`Flexible`/`Saw`/`Short`; hard block or access denied → `Bold`/`Rigid`/`Saw`/`Extended`.
- *Reversibility* — destructive/irreversible action that was blocked → escalate to `Bold`; recoverable error → `Substantial`.

### Notifications & Alerts

**Preferred spec:** `Bumps` shape at `Substantial` or `Bold` intensity — rounded pulses announce attention without demanding it. `Bumps` is the most natural shape for incoming signals.

**Avoid:** `Gentle` when the alert must be noticed; `Impulse` or `Short` duration for a notification — it will feel like a UI tap rather than an incoming signal.

**Context you need before choosing:**
- *Criticality* — must-not-miss → `Bold`/`Extended`; standard notification → `Substantial`/`Long`; passive background signal → `Gentle`/`Bumps`.
- *Source* — personal/social (message, call) → warmer texture (`Flexible`/`Soft`); system/service alert → crisper texture (`Rigid`/`Flexible`).

### Warnings & Tension

**Preferred spec:** `Bold` + `Extended` with `Pattern` or `Ramp` shape — warnings should build or persist, not just tap once. Sustained duration communicates ongoing risk.

**Avoid:** `Short` or `Impulse` for a warning state — it reads as a notification, not a caution. `Soft` texture weakens the urgency.

**Context you need before choosing:**
- *Tension arc* — static/persistent warning → `Pattern`; escalating urgency → `Ramp`; repeated alarm-like pulse → `Saw`/`Pattern`.
- *Severity relative to errors in the same app* — warnings should be perceptibly lighter than hard errors; if errors use `Bold`/`Rigid`, warnings can use `Substantial`/`Flexible` unless they are critical.

### UI Interactions

**Preferred spec:** `Substantial` + `Rigid` + `Impulse` for primary controls; `Gentle` + `Impulse` for secondary or ghost-button interactions (treat as progressive enhancement).

**Avoid:** `Extended` or `Long` for any single tap — duration must match interaction speed. `Bold` for routine taps — it will make every interaction feel heavy and alarming.

**Context you need before choosing:**
- *Control weight* — primary CTA, toggle, mechanical control (shutter, lock) → `Substantial`/`Rigid`; secondary button, chip, icon tap → `Gentle`/`Impulse`.
- *Is the haptic essential or enhancement?* — if it can be skipped without confusing the user, use `Gentle`; if it confirms a state change the user must perceive, use `Substantial`.

### Ambient & Background

**Preferred spec:** `Gentle` or `Substantial` + `Solid` or `Pattern` + `Extended` — ambient haptics should feel like texture, not events. Continuous shapes work best.

**Avoid:** `Bold` for ambient states — it will feel like an alert, not an atmosphere. `Impulse` or `Short` duration is meaningless for a looping state. `Saw` shape reads as an error, not background texture.

**Context you need before choosing:**
- *Is this looping or one-shot?* — looping ambient → `Solid`/`Extended` or `Pattern`/`Extended`; one-shot ambient cue (attention nudge) → `Bumps`/`Long`.
- *Emotional quality* — calm/meditative → `Soft`/`Gentle`; active/monitoring → `Flexible`/`Substantial` with `Pattern`.
- *Does it layer with other haptics?* — if discrete events fire over the ambient loop, the ambient layer must be `Gentle` enough not to mask them.

### Games & Physical Metaphors

**Preferred spec:** `Bold` + `Peak` or `Ramp` for impact events; `Solid`/`Extended` for continuous physical states (engine, heartbeat).

**Avoid:** `Soft` texture for sharp physical impacts — it will feel cushioned rather than physical. `Impulse` duration for explosions or collisions — the event needs mass, which requires `Long` or `Extended`.

**Context you need before choosing:**
- *Event type* — discrete impact (hit, collision, tap) → `Bold`/`Peak` or `Ramp`/`Short` or `Long`; continuous physical state (engine running, heartbeat) → `Solid` or `Pattern`/`Extended`.
- *Physical metaphor weight* — a pebble drop vs. a boulder collapse should differ in both Intensity and Duration, not just amplitude.
- *Texture* — mechanical/digital game event → `Rigid`; organic/physical world event → `Flexible` or `Soft`.

### Keyboard & Typing

**Preferred spec:** `Rigid` + `Impulse` — keyboard haptics must be near-instantaneous and mechanically crisp. Every keystroke fires at typing speed, so `Impulse` duration is required to avoid blur.

**Avoid:** `Short` or longer duration — it will overlap with the next keypress at normal typing speed. `Soft` or `Flexible` texture weakens the mechanical character that makes keyboard haptics useful.

**Context you need before choosing:**
- *Key type differentiation* — if the design distinguishes key types (mechanical vs. membrane feel), use amplitude contrast between them rather than relying on texture alone; texture differences are subtle and easily missed.
- *Typing speed* — high-frequency input means even `Short` can overlap; confirm the call site fires per keypress, not per word.

### Selection & Scrolling

**Preferred spec (per-item tick):** `Rigid` + `Impulse` — scroll ticks must be imperceptible in isolation but accumulate into a tactile rhythm. Brief and crisp.

**Preferred spec (threshold / snap-to):** `Substantial` + `Rigid` + `Short` — a snap-point or pull-to-refresh threshold is a one-shot event; it should feel like a definitive click, not a tick.

**Avoid:** Substituting a heavier spec for a `Gentle` tick at high firing frequency — it will feel intrusive rather than textured. If `Gentle`/`Impulse` is too subtle, omit ticks rather than replacing them with something heavier.

**Context you need before choosing:**
- *Firing frequency* — per-detent scroll tick → `Impulse` only; one-shot threshold → `Short` acceptable.
- *Interaction type* — continuous scroll picker → `Rigid`/`Impulse`; pull-to-refresh → `Substantial`/`Rigid`/`Short`; dial complete/unlock ceremony → `Substantial` or `Bold`/`Extended`.

### Drag & Drop

**Preferred spec:** `Substantial` + `Rigid` + `Short` for snap-to-target and drop placement — the moment of placement should feel like something locking in.

**Avoid:** `Gentle` for drop events — the signal will be lost during the physical motion of dragging. `Extended` duration for a drop landing — it should be instantaneous.

**Context you need before choosing:**
- *Event type* — snap-to-target mid-drag → lighter spec (`Substantial`/`Short`); final drop placement → equal or heavier (`Substantial` or `Bold`/`Short`); destructive/irreversible drop → `Bold`/`Rigid`.
- *Texture* — mechanical list reorder → `Rigid`; organic drag (photo collage, freeform canvas) → `Flexible`.

### Navigation & Transitions

**Preferred spec:** `Substantial` + `Short` for modal present/dismiss; `Gentle` + `Impulse` for lightweight navigation (tab switch, back).

**Avoid:** `Bold` for routine navigation — it makes every screen feel like a critical event. `Extended` for any transition — it outlasts the animation and feels slow.

**Context you need before choosing:**
- *Transition weight* — modal or sheet (heavier) → `Substantial`/`Short`; tab switch, back, push/pop → `Gentle`/`Impulse`.
- *Is the haptic essential or enhancement?* — navigation haptics are almost always enhancement. If it would feel odd without haptics, it's enhancement; design accordingly and do not escalate beyond `Substantial`.
- *Direction* — presenting/opening → slightly heavier; dismissing/closing → slightly lighter or omit entirely.

### Media & Playback

**Preferred spec:** `Substantial` + `Rigid` + `Short` for discrete playback events (play/pause, chapter marker); `Gentle` + `Rigid` + `Impulse` for high-frequency scrub ticks.

**Avoid:** `Bold` for routine play/pause — it should feel like a control tap, not an alert. `Short` or longer for scrub ticks — at drag speed, anything longer than `Impulse` turns into continuous noise.

**Context you need before choosing:**
- *Firing frequency* — scrub ticks fire continuously while the user drags; use `Impulse` duration only. One-shot events (play, pause, chapter boundary) can use `Short` or `Long`.
- *Event weight* — chapter boundary or major marker → `Substantial`/`Long`; routine play/pause → `Substantial`/`Short`; scrub tick → `Gentle`/`Impulse`.

### Authentication & Biometrics

**Preferred spec:** `Substantial` + `Rigid` + `Impulses` or `Peak` + `Short` for success; `Bold` + `Rigid` + `Saw` + `Short` for failure/access denied.

**Avoid:** `Soft` texture for authentication feedback — the signal needs to feel definitive, not cushioned. `Gentle` for any auth event — this is a trust signal; it must register.

**Context you need before choosing:**
- *Success vs. failure* — success should feel settled and final (`Impulses`/`Peak`); failure should feel sharp and interruptive (`Saw`).
- *App tone* — if the app is explicitly warm or wellness-oriented, `Flexible` texture is acceptable for success; otherwise default to `Rigid` for clarity.

### Onboarding & Tutorials

**Preferred spec:** `Substantial` + `Bumps` + `Short` for step-completion moments; `Gentle` + `Bumps`/`Solid` + `Extended` for ambient guidance cues (only when continuous looping is appropriate).

**Avoid:** `Bold` for guidance cues — onboarding should feel supportive, not alarming. `Impulse` duration for step milestones — the completion moment deserves more than a tap.

**Context you need before choosing:**
- *Step importance* — primary milestone (first action completed, setup done) → `Substantial`/`Bumps`/`Short` or `Long`; subtle ambient guide → `Gentle`/`Soft`.
- *Continuous vs. one-shot* — spotlight/attention nudge that pulses → `Bumps`/`Extended`; step-completion confirmation → `Bumps`/`Short` or `Long`.
