## Custom Patterns with PatternComposer

When no built-in preset fits, define a custom `PatternData`. A pattern combines **discrete** events (individual taps/impacts at specific moments) and a **continuous** envelope (sustained vibration whose amplitude and frequency evolve over time).

Get a fresh `PatternComposer` from `pulsar.getPatternComposer()`. Each call returns a new instance — manage multiple patterns independently by holding multiple composers.

### PatternData Types

```kotlin
import com.swmansion.pulsar.types.PatternData
import com.swmansion.pulsar.types.ContinuousPattern
import com.swmansion.pulsar.types.ConfigPoint
import com.swmansion.pulsar.types.ValuePoint

// ConfigPoint — one discrete haptic event
data class ConfigPoint(
    val time: Long,        // Milliseconds from pattern start
    val amplitude: Float,  // Intensity (0–1)
    val frequency: Float   // Sharpness (0–1); 0 = soft/round, 1 = crisp/sharp
)

// ValuePoint — one control point in a continuous envelope
data class ValuePoint(
    val time: Long,   // Milliseconds from pattern start
    val value: Float  // Amplitude or frequency value (0–1)
)

// ContinuousPattern — two time-varying envelopes
data class ContinuousPattern(
    val amplitude: List<ValuePoint>,  // Amplitude envelope
    val frequency: List<ValuePoint>   // Frequency envelope
)

// PatternData — the complete pattern passed to PatternComposer
data class PatternData(
    val continuousPattern: ContinuousPattern,
    val discretePattern: List<ConfigPoint>
)
```

### PatternComposer Methods

```kotlin
val composer = pulsar.getPatternComposer()

fun parsePattern(hapticsData: PatternData)  // Parse and prepare for playback
fun play()                                   // Play the parsed pattern
fun stop()                                   // Stop active playback
```

### Example

```kotlin
import com.swmansion.pulsar.composers.PatternComposer
import com.swmansion.pulsar.types.*

val composer = pulsar.getPatternComposer()

val pattern = PatternData(
    discretePattern = listOf(
        ConfigPoint(time = 0,   amplitude = 1.0f, frequency = 0.8f), // sharp opening tap
        ConfigPoint(time = 100, amplitude = 0.5f, frequency = 0.4f), // softer second tap
        ConfigPoint(time = 200, amplitude = 0.2f, frequency = 0.2f), // gentle close
    ),
    continuousPattern = ContinuousPattern(
        amplitude = listOf(
            ValuePoint(time = 0,   value = 0.3f),
            ValuePoint(time = 150, value = 1.0f),
            ValuePoint(time = 300, value = 0.0f),
        ),
        frequency = listOf(
            ValuePoint(time = 0,   value = 0.2f),
            ValuePoint(time = 300, value = 0.7f),
        )
    )
)

composer.parsePattern(pattern)

// Later — in a click handler or gesture callback
button.setOnClickListener {
    composer.play()
}
```

### Pattern Design Tips

See the "Custom Pattern Parameters" section in `../common/design-principles.md` for amplitude and frequency range guidance, and when to use discrete vs. continuous patterns.

**Android-specific notes:**
- On devices without Envelope API (below API 36), `continuousPattern` frequency is ignored; only amplitude is used.
- On devices below API 26, only the presence/absence of vibration at each time point is meaningful (no amplitude control).
