## Full Types Reference

### PatternData

```kotlin
data class PatternData(
    val continuousPattern: ContinuousPattern,
    val discretePattern: List<ConfigPoint>
)
```

Convenience constructor also available for raw float arrays (used internally):

```kotlin
PatternData(
    rawContinuousPattern: List<List<List<Float>>> = listOf(listOf(), listOf()),
    rawDiscretePattern: List<List<Float>> = listOf()
)
```

### ContinuousPattern

```kotlin
data class ContinuousPattern(
    val amplitude: List<ValuePoint>,  // Amplitude envelope control points
    val frequency: List<ValuePoint>   // Frequency envelope control points
)
```

### ValuePoint

```kotlin
data class ValuePoint(
    val time: Long,   // Milliseconds from pattern start
    val value: Float  // Normalized value (0–1)
)
```

### ConfigPoint

```kotlin
data class ConfigPoint(
    val time: Long,        // Milliseconds from pattern start
    val amplitude: Float,  // Intensity (0–1)
    val frequency: Float   // Sharpness (0–1); ignored on devices without envelope support
)
```

### Preset (interface)

```kotlin
interface Preset {
    fun play()
}
```

All built-in preset objects implement `Preset`. `presets.getByName(name)` returns a `Preset?`.

### CompatibilityMode

```kotlin
enum class CompatibilityMode {
    NO_SUPPORT,
    MINIMAL_SUPPORT,
    LIMITED_SUPPORT,
    STANDARD_SUPPORT,
    ADVANCED_SUPPORT,
}
```

Comparable — `CompatibilityMode.STANDARD_SUPPORT >= CompatibilityMode.LIMITED_SUPPORT` is `true`.

### RealtimeComposerStrategy

```kotlin
enum class RealtimeComposerStrategy {
    ENVELOPE,
    PRIMITIVE_TICK,
    PRIMITIVE_COMPLEX,
    ENVELOPE_WITH_DISCRETE_PRIMITIVES,
}
```
