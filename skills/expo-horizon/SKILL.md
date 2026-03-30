---
name: expo-horizon
description: "Software Mansion's guide for migrating Expo SDK apps to Meta Quest using expo-horizon packages. Use when adding Meta Quest or Meta Horizon OS support to an existing Expo or React Native project. Trigger on: Meta Quest, Horizon OS, Quest 2, Quest 3, Quest 3S, VR app, expo-horizon-core, expo-horizon-location, expo-horizon-notifications, build flavors for Quest, panel sizing, VR headtracking, Horizon App ID, quest build variant, isHorizonDevice, isHorizonBuild, migrate expo-location to Quest, migrate expo-notifications to Quest, Meta Horizon Store publishing, or any task involving running an Expo app on Meta Quest hardware."
---

# Expo Horizon: Migrating Expo SDK to Meta Quest

Software Mansion's production guide for adding Meta Quest support to Expo apps using the [expo-horizon](https://github.com/software-mansion-labs/expo-horizon) packages.

Read the relevant reference for the topic at hand. All references are in `references/`.

## Decision Tree

```
What do you need to do?
│
├── Starting from scratch or adding Quest support to an existing Expo app?
│   └── See setup.md
│       ├── Install expo-horizon-core
│       ├── Configure the config plugin
│       ├── Set up build scripts
│       └── Add runtime device detection
│
├── Need location services on Quest?
│   └── See location.md
│       ├── Replace expo-location with expo-horizon-location
│       ├── Understand Quest location limitations (no GPS, no geocoding)
│       └── Handle feature parity differences
│
├── Need push notifications on Quest?
│   └── See notifications.md
│       ├── Replace expo-notifications with expo-horizon-notifications
│       ├── Configure horizonAppId for push tokens
│       └── Handle feature parity differences
│
└── Need to build, run, or publish for Quest?
    └── See build-and-deploy.md
        ├── Build variants (questDebug, questRelease, mobileDebug, mobileRelease)
        ├── Running on Quest hardware
        └── Meta Horizon Store publishing requirements
```

## Critical Rules

- **Always install `expo-horizon-core` first.** It is required by all other expo-horizon packages and sets up the `quest`/`mobile` build flavors that other packages depend on.

- **Use `quest` build variants only on Meta Quest devices.** Running `questDebug` or `questRelease` builds on standard Android phones is unsupported and will behave unexpectedly.

- **Set `supportedDevices` in the config plugin.** This is required for Meta Horizon Store submission. Use pipe-separated values: `"quest2|quest3|quest3s"`.

- **Run `npx expo prebuild --clean` after any plugin config change.** The config plugin modifies native project files at prebuild time. Stale native projects will not reflect your changes.

- **Replace imports, not just packages.** When migrating from `expo-location` or `expo-notifications`, update all import statements to use the new package names (`expo-horizon-location`, `expo-horizon-notifications`).

- **Quest has no GPS, magnetic sensors, or Geocoder.** Features like heading, geocoding, reverse geocoding, and geofencing are unavailable on Quest. Guard these calls with `ExpoHorizon.isHorizonDevice` or `ExpoHorizon.isHorizonBuild`.

- **Push notifications require `horizonAppId`.** Without it, `getDevicePushTokenAsync` will not return a valid token on Quest devices.

- **Expo Go is not supported.** You must use custom development builds (`npx expo prebuild`).

## References

| File | When to read |
|------|-------------|
| `references/setup.md` | Installing expo-horizon-core, configuring the Expo config plugin, plugin options (horizonAppId, panel sizing, supportedDevices, VR headtracking, allowBackup), adding quest/mobile build scripts to package.json, runtime detection API (isHorizonDevice, isHorizonBuild, horizonAppId), accessing Horizon App ID from native modules |
| `references/location.md` | Migrating from expo-location, Quest location limitations (no GPS, no geocoding, no heading, no geofencing, no background location), network provider behavior, feature support matrix |
| `references/notifications.md` | Migrating from expo-notifications, configuring push notifications for Quest, Horizon push token type, Firebase vs Meta push service, feature support matrix, unsupported features (Expo Push Service, badge counts) |
| `references/build-and-deploy.md` | Build flavors (quest/mobile), build variants (questDebug/questRelease/mobileDebug/mobileRelease), running on Quest hardware, package.json scripts, Meta Horizon Store requirements, Meta Quest Developer Hub |
