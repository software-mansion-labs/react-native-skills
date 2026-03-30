# Building and Deploying for Meta Quest

Build variants, running on Quest hardware, and Meta Horizon Store publishing.

---

## Build Flavors

The `expo-horizon-core` config plugin creates two Android build flavors:

| Flavor | Target | Description |
|---|---|---|
| `mobile` | Standard Android phones/tablets | Uses Google Play Services, Firebase, standard manifest |
| `quest` | Meta Quest devices | Uses Horizon OS manifest, Meta push service, no Google Play Services |

Each flavor has debug and release variants:

| Variant | Use case |
|---|---|
| `mobileDebug` | Development on Android phones/tablets |
| `mobileRelease` | Production for Google Play Store |
| `questDebug` | Development on Meta Quest devices |
| `questRelease` | Production for Meta Horizon Store |

---

## Package.json Scripts

Add these scripts for convenience:

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

---

## Running on Quest Hardware

### Prerequisites

1. Enable Developer Mode on your Quest device via the Meta Quest app on your phone
2. Install [Meta Quest Developer Hub](https://developers.meta.com/horizon/documentation/android-apps/meta-quest-developer-hub) for device management, video casting, storage access, and app sideloading
3. Connect your Quest via USB or Wi-Fi ADB

### Development Build

```bash
# First time: generate native project
npx expo prebuild --clean

# Run debug build on connected Quest
npx expo run:android --variant questDebug
```

### Release Build

```bash
npx expo run:android --variant questRelease
```

---

## Prebuild Workflow

The config plugin modifies native project files during prebuild. Follow this workflow:

```
Change plugin config in app.json / app.config.ts
    │
    ├── Run: npx expo prebuild --clean
    │   (regenerates android/ and ios/ directories)
    │
    └── Run: npx expo run:android --variant questDebug
        (builds and installs on connected Quest)
```

Always run `npx expo prebuild --clean` after:
- Adding or removing expo-horizon plugins
- Changing any plugin option (horizonAppId, supportedDevices, panel size, etc.)
- Upgrading expo-horizon package versions

---

## Meta Horizon Store Publishing

### Manifest Requirements

The config plugin automatically handles most manifest requirements. Ensure you have configured:

- `supportedDevices` -- pipe-separated list of target Quest devices
- `allowBackup` set to `false` (recommended by Meta for sensitive data)
- VR headtracking feature declaration (enabled by default)

For the full manifest checklist, webfetch [Meta Horizon Store manifest requirements](https://developers.meta.com/horizon/resources/publish-mobile-manifest/).

### Build the Release APK/AAB

```bash
npx expo run:android --variant questRelease
```

The output APK/AAB is in `android/app/build/outputs/`.

### Key Publishing Notes

- Only `questRelease` builds should be submitted to the Meta Horizon Store
- `mobileRelease` builds go to the Google Play Store as usual
- The same codebase produces both builds -- no code duplication needed
- Test on actual Quest hardware before submission; the Quest emulator has limited fidelity

---

## Dual-Platform Development Tips

```
Running the same app on both platforms?
│
├── Use isHorizonBuild for build-time branching
│   (different native code paths, e.g., Firebase vs Meta push)
│
├── Use isHorizonDevice for runtime branching
│   (same build, different behavior on Quest hardware)
│
├── Gate unavailable features behind device checks
│   (heading, geocoding, badge counts, Expo Push Service)
│
└── Test both variants regularly
    ├── npm run android   (mobile build on phone)
    └── npm run quest     (quest build on Quest)
```
