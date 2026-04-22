## Custom Patterns with PatternComposer

When no built-in preset fits your use case, compose a custom haptic pattern using `PatternComposer`. A pattern combines **discrete** events (individual taps at specific moments) and a **continuous** envelope (sustained vibration with evolving amplitude and frequency).

Get a new `PatternComposer` instance from `pulsar.getPatternComposer()`. Each call returns a fresh instance — you can hold multiple composers for different patterns independently.

```swift
let composer = pulsar.getPatternComposer()
```

### Types

#### `PatternData`

The top-level container for a complete haptic pattern.

```swift
class PatternData: NSObject, Codable {
  let continuousPattern: ContinuousPattern
  let discretePattern: [DiscretePoint]

  init(
    continuousPattern: ContinuousPattern,
    discretePattern: [DiscretePoint]
  )
}
```

#### `ContinuousPattern`

Two envelope curves — one for amplitude, one for frequency — that shape the sustained part of the haptic.

```swift
class ContinuousPattern: NSObject, Codable {
  let amplitude: [ValuePoint]
  let frequency: [ValuePoint]

  init(amplitude: [ValuePoint], frequency: [ValuePoint])
}
```

#### `ValuePoint`

A single control point on a continuous curve.

```swift
class ValuePoint: NSObject, Codable {
  let time: Double    // Milliseconds from pattern start
  let value: Float    // Normalized value 0–1

  init(time: Double, value: Float)
}
```

#### `DiscretePoint`

A single discrete (transient) haptic event.

```swift
class DiscretePoint: NSObject, Codable {
  let time: Double      // Milliseconds from pattern start
  let amplitude: Float  // Intensity 0–1
  let frequency: Float  // Sharpness 0–1 (0 = round/soft, 1 = crisp)

  init(time: Double, amplitude: Float, frequency: Float)
}
```

### PatternComposer Methods

| Method | Description |
|---|---|
| `parsePattern(hapticsData:)` | Parse a `PatternData` and prepare it for playback |
| `playPattern(hapticsData:)` | Parse and immediately play a pattern (combines `parsePattern` + `play`) |
| `play()` | Play the previously parsed pattern |
| `stop()` | Stop active playback |

### Example

```swift
import Pulsar

let pulsar = Pulsar()
let composer = pulsar.getPatternComposer()

let pattern = PatternData(
  continuousPattern: ContinuousPattern(
    amplitude: [
      ValuePoint(time: 0,   value: 0.0),
      ValuePoint(time: 50,  value: 1.0),
      ValuePoint(time: 300, value: 0.0),
    ],
    frequency: [
      ValuePoint(time: 0,   value: 0.3),
      ValuePoint(time: 300, value: 0.8),
    ]
  ),
  discretePattern: [
    DiscretePoint(time: 0,   amplitude: 1.0, frequency: 0.8),  // crisp opening tap
    DiscretePoint(time: 100, amplitude: 0.5, frequency: 0.4),  // softer follow-up
    DiscretePoint(time: 200, amplitude: 0.2, frequency: 0.2),  // gentle close
  ]
)

// Parse once, play many times
composer.parsePattern(hapticsData: pattern)
composer.play()

// Or parse and play in one call
composer.playPattern(hapticsData: pattern)
```

### Using PatternComposer in SwiftUI

```swift
import SwiftUI
import Pulsar

struct CustomHapticView: View {
  private let pulsar = Pulsar()

  private let pattern = PatternData(
    continuousPattern: ContinuousPattern(
      amplitude: [
        ValuePoint(time: 0,   value: 0.0),
        ValuePoint(time: 200, value: 0.8),
        ValuePoint(time: 400, value: 0.0),
      ],
      frequency: [
        ValuePoint(time: 0,   value: 0.2),
        ValuePoint(time: 400, value: 0.6),
      ]
    ),
    discretePattern: [
      DiscretePoint(time: 0, amplitude: 1.0, frequency: 0.7),
    ]
  )

  var body: some View {
    Button("Custom Haptic") {
      // playPattern parses and plays in one call — fine for one-shot use
      pulsar.getPatternComposer().playPattern(hapticsData: pattern)
    }
  }
}
```

### Pattern Design Tips

See the "Custom Pattern Parameters" section in `../common/design-principles.md` for amplitude and frequency range guidance, and when to use discrete vs. continuous patterns.
