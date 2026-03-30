# Setup and Configuration

Installation, config plugin setup, runtime API, and native module integration for expo-horizon-core.

For the full README, webfetch the [expo-horizon-core README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-core/README.md).

---

## Prerequisites

- Expo SDK 54 or later (`expo` package version 54.0.13+)
- Android development environment configured
- Meta Quest developer account (for publishing)
- Expo Go is **not supported** -- use custom development builds

## Installation

```bash
npx expo install expo-horizon-core
```

---

## Config Plugin

Add the plugin to your `app.json` or `app.config.[js|ts]`:

```json
{
  "expo": {
    "plugins": [
      [
        "expo-horizon-core",
        {
          "horizonAppId": "your-horizon-app-id",
          "defaultHeight": "640dp",
          "defaultWidth": "1024dp",
          "supportedDevices": "quest2|quest3|quest3s",
          "disableVrHeadtracking": false,
          "allowBackup": false
        }
      ]
    ]
  }
}
```

After adding or changing the plugin config, regenerate native files:

```bash
npx expo prebuild --clean
```

### Plugin Options

| Option | Type | Required | Default | Description |
|---|---|---|---|---|
| `horizonAppId` | `string` | No | `""` | Meta Horizon application ID. Required by `expo-horizon-notifications` for push tokens. |
| `defaultHeight` | `string` | No | Not added | Default panel height in dp (e.g., `"640dp"`). |
| `defaultWidth` | `string` | No | Not added | Default panel width in dp (e.g., `"1024dp"`). |
| `supportedDevices` | `string` | Yes | None | Pipe-separated list of Quest devices: `"quest2\|quest3\|quest3s"`. |
| `disableVrHeadtracking` | `boolean` | No | `false` | Set `true` to omit the `android.hardware.vr.headtracking` manifest entry. |
| `allowBackup` | `boolean` | No | `false` | Set `true` to allow Android backup in the Quest build. Meta recommends `false` for sensitive data. |

### Panel Sizing

If your app renders with black bars after setting `defaultWidth` or `defaultHeight`, make sure the `orientation` value in your Expo config matches the specified dimensions. For landscape panels, use `"landscape"`.

For sizing guidelines, webfetch [Meta Panel Sizing docs](https://developers.meta.com/horizon/documentation/android-apps/panel-sizing).

---

## Build Scripts

After installing expo-horizon-core, add build scripts to your `package.json` so you can target both mobile and Quest from the same project:

```json
{
  "scripts": {
    "android": "expo run:android --variant mobileDebug",
    "android:release": "expo run:android --variant mobileRelease",
    "quest": "expo run:android --variant questDebug",
    "quest:release": "expo run:android --variant questRelease",
    "ios": "expo run:ios",
    "web": "expo start --web"
  }
}
```

Then run with:

```bash
npm run android    # Standard Android phone/tablet
npm run quest      # Meta Quest device
```

The `quest` variants include Horizon-specific manifest settings and native code paths. The `mobile` variants behave identically to a standard Expo Android build.

---

## What the Plugin Does

The config plugin automatically:

1. Adds `quest` and `mobile` build flavors to `build.gradle`
2. Creates a Horizon-specific `AndroidManifest.xml` with required permissions and features
3. Configures panel dimensions on the main activity
4. Adds `com.oculus.supportedDevices` meta-data to the manifest
5. Sets up `android.hardware.vr.headtracking` uses-feature (unless disabled)
6. Passes `horizonAppId` as a Gradle property for native module access

---

## Runtime API

```typescript
import ExpoHorizon from 'expo-horizon-core';
```

| Property | Type | Description |
|---|---|---|
| `isHorizonDevice` | `boolean` | `true` if running on a physical Meta Quest device. |
| `isHorizonBuild` | `boolean` | `true` if built with the `quest` build flavor. |
| `horizonAppId` | `string \| null` | The Horizon App ID from config, or `null` if not set. |

### Conditional Feature Pattern

```typescript
import ExpoHorizon from 'expo-horizon-core';

// Gate features unavailable on Quest
if (!ExpoHorizon.isHorizonDevice) {
  const heading = await Location.watchHeadingAsync(callback);
}

// Quest-specific UI
if (ExpoHorizon.isHorizonBuild) {
  // Render VR-optimized layout with larger touch targets
}
```

Use `isHorizonDevice` for runtime hardware checks (physical device detection). Use `isHorizonBuild` for build-time feature gating (which native code is included).

---

## Accessing Horizon App ID in Native Modules

To read the Horizon App ID from custom Kotlin native modules:

1. Add the config field to your module's `build.gradle`:

```gradle
def horizonAppIdConfigField = "\"${project.findProperty('horizonAppId') ?: ''}\""

android {
  defaultConfig {
    buildConfigField "String", "META_HORIZON_APP_ID", horizonAppIdConfigField
  }
}
```

2. Access it in Kotlin:

```kotlin
val horizonAppId = BuildConfig.META_HORIZON_APP_ID
```

---

## Version Compatibility

| `expo-horizon-core` | Expo SDK |
|---|---|
| 55.0.0 | 55 |
| 1.0.7 | 54 |
