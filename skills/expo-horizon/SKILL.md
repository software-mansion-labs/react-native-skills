---
name: expo-horizon
description: "Software Mansion's guide for migrating Expo SDK apps to Meta Quest using expo-horizon packages. Use when adding Meta Quest or Meta Horizon OS support to an existing Expo or React Native project. Trigger on: Meta Quest, Horizon OS, Quest 2, Quest 3, Quest 3S, VR app, expo-horizon-core, expo-horizon-location, expo-horizon-notifications, build flavors for Quest, panel sizing, VR headtracking, Horizon App ID, quest build variant, isHorizonDevice, isHorizonBuild, migrate expo-location to Quest, migrate expo-notifications to Quest, Meta Horizon Store publishing, or any task involving running an Expo app on Meta Quest hardware."
---

# Expo Horizon: Migrating Expo SDK to Meta Quest

Software Mansion's production guide for adding Meta Quest support to Expo apps using the [expo-horizon](https://github.com/software-mansion-labs/expo-horizon) packages.

**This skill does not bundle a copy of the docs.** For any task below, always webfetch the linked official README or Meta documentation page to get up-to-date installation steps, plugin options, API surface, and feature matrices. This skill only captures the decision tree, critical rules, and non-obvious gotchas that agents routinely miss.

## Decision Tree

```
What do you need to do?
│
├── Starting from scratch or adding Quest support to an existing Expo app?
│   └── Follow the "Setup Workflow" below (do NOT auto-install location or
│       notifications packages) and webfetch: expo-horizon-core README
│       ├── Install expo-horizon-core
│       ├── Configure the config plugin (horizonAppId, panel size, supportedDevices)
│       ├── Add quest/mobile build scripts
│       ├── Add runtime device detection (isHorizonDevice, isHorizonBuild)
│       └── Detect expo-location / expo-notifications, then ASK before migrating
│
├── Need location services on Quest?
│   └── Webfetch: expo-horizon-location README
│       ├── Replace expo-location with expo-horizon-location
│       ├── Review the feature support matrix
│       └── Guard unsupported calls (heading, geocoding, geofencing, background)
│
├── Need push notifications on Quest?
│   └── Webfetch: expo-horizon-notifications README
│       ├── Replace expo-notifications with expo-horizon-notifications
│       ├── Configure horizonAppId in expo-horizon-core
│       ├── Use getDevicePushTokenAsync (Expo Push Service is not supported)
│       └── Skip badge counts (not supported on Quest)
│
└── Need to build, run, or publish for Quest?
    └── Webfetch: expo-horizon-core README (build variants) + Meta docs below
        ├── Build variants: questDebug, questRelease, mobileDebug, mobileRelease
        ├── Meta Quest Developer Hub (device management, sideloading)
        └── Meta Horizon Store manifest requirements
```

## Setup Workflow (adding Quest support to an existing Expo app)

Follow these steps **in order** when the user asks to add Meta Quest support. Do not combine steps 2 and 3 into a single install command — the sibling packages require explicit user confirmation.

1. **Install and configure `expo-horizon-core`.**
   - Run `npx expo install expo-horizon-core`.
   - **Ask the user for the config plugin values before writing them.** Present all options in a single prompt so the user can paste custom values or accept defaults in one pass. Show each option on its own line with its default in brackets, e.g.:
     > I'll add the `expo-horizon-core` config plugin to `app.json`. Please confirm or override each value (press Enter / reply "default" to accept the bracketed default):
     >
     > - `supportedDevices` [`quest2|quest3|quest3s`] — pipe-separated Quest devices your app supports (**required** for Meta Horizon Store submission).
     > - `horizonAppId` [empty] — Meta Horizon application ID. Leave empty unless you plan to use push notifications; required by `expo-horizon-notifications` to issue device push tokens.
     > - `defaultWidth` [`1024dp`] — Default panel width. Leave blank to omit.
     > - `defaultHeight` [`640dp`] — Default panel height. Leave blank to omit. If you set width/height, make sure your Expo `orientation` matches (use `"landscape"` for wide panels).
     > - `disableVrHeadtracking` [`false`] — Set `true` to omit the `android.hardware.vr.headtracking` manifest entry.
     > - `allowBackup` [`false`] — Meta recommends `false` for sensitive data; set `true` only if you need Android backup in the Quest build.
   - Only write the plugin config after the user replies. Omit any option the user left blank so the package's own default applies (don't write empty strings for `horizonAppId`, `defaultWidth`, or `defaultHeight`).
   - Add `quest` / `mobile` build scripts to `package.json`.
   - Run `npx expo prebuild --clean`.
   - Webfetch the [expo-horizon-core README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-core/README.md) for current option names and defaults before asking — the defaults above can change between releases.

2. **Detect existing location / notification packages. Do NOT install the horizon equivalents yet.**
   - Read the project's `package.json`.
   - Check `dependencies` and `devDependencies` for `expo-location` and `expo-notifications`.
   - If neither is present, skip the rest of this workflow — the user has no migration to do.

3. **For each detected package, ask the user before migrating.**
   - If `expo-location` is found, ask:
     > "I found `expo-location` in your project. Do you want me to replace it with `expo-horizon-location` so location works on Meta Quest? (Quest has no GPS, heading, geocoding, or geofencing — unsupported calls will need to be guarded with `ExpoHorizon.isHorizonDevice`.)"
   - If `expo-notifications` is found, ask:
     > "I found `expo-notifications` in your project. Do you want me to replace it with `expo-horizon-notifications` so push notifications work on Meta Quest? This requires a `horizonAppId` in the core config plugin, uses Meta's push service (not the Expo Push Service), and does not support badge counts or `getExpoPushTokenAsync`."
   - Present the questions together if both packages are present.
   - **Wait for an explicit answer before running any install or edit for these packages.**

4. **Only after the user confirms**, perform the migration for the approved package(s):
   - Install the horizon equivalent (`npx expo install expo-horizon-location` or `expo-horizon-notifications`).
   - Uninstall the original (`npm uninstall expo-location` or `expo-notifications`).
   - Update all `import` statements to the horizon package name.
   - For notifications: ensure `horizonAppId` is set in the `expo-horizon-core` plugin config and add `expo-horizon-notifications` to the plugins array.
   - Run `npx expo prebuild --clean` again.
   - Webfetch the relevant README (location or notifications) for the full feature support matrix and guard unsupported calls behind `ExpoHorizon.isHorizonDevice`.

5. **If the user declines a migration**, leave the original package untouched and note in the summary that feature X (location / notifications) will not work on the `quest` build until migrated.

## Critical Rules

- **Always install `expo-horizon-core` first.** It is required by all other expo-horizon packages and sets up the `quest`/`mobile` build flavors that other packages depend on.

- **Never auto-install `expo-horizon-location` or `expo-horizon-notifications`.** When adding Quest support, detect existing `expo-location` / `expo-notifications` dependencies in `package.json` and ask the user whether to migrate each one. Install and configure only the packages the user explicitly approves. See the Setup Workflow above.

- **Use `quest` build variants only on Meta Quest devices.** Running `questDebug` or `questRelease` builds on standard Android phones is unsupported and will behave unexpectedly.

- **Set `supportedDevices` in the config plugin.** This is required for Meta Horizon Store submission. Use pipe-separated values: `"quest2|quest3|quest3s"`.

- **Run `npx expo prebuild --clean` after any plugin config change.** The config plugin modifies native project files at prebuild time. Stale native projects will not reflect your changes.

- **Replace imports, not just packages.** When migrating from `expo-location` or `expo-notifications`, update all import statements to use the new package names (`expo-horizon-location`, `expo-horizon-notifications`).

- **Quest has no GPS, magnetic sensors, or Geocoder.** Features like heading, geocoding, reverse geocoding, and geofencing are unavailable on Quest. Guard these calls with `ExpoHorizon.isHorizonDevice` or `ExpoHorizon.isHorizonBuild`.

- **Push notifications require `horizonAppId`.** Without it, `getDevicePushTokenAsync` will not return a valid token on Quest devices. Use `getDevicePushTokenAsync` (not `getExpoPushTokenAsync`) on Quest; send the returned `{ type: 'horizon', data }` token to your backend and deliver via Meta's push service.

- **`isHorizonDevice` vs `isHorizonBuild`.** Use `isHorizonDevice` for runtime hardware checks (physical Quest detection). Use `isHorizonBuild` for build-time feature gating (which native code was compiled in).

- **Expo Go is not supported.** You must use custom development builds via `npx expo prebuild`.

## Official References

Always webfetch the raw markdown (`raw.githubusercontent.com/...`) if the HTML view does not render the source; the raw URL is the source of truth.

| Topic | Official source |
|-------|-----------------|
| Repo overview and package list | [expo-horizon README](https://github.com/software-mansion-labs/expo-horizon/blob/main/README.md) |
| Install, config plugin options, runtime API, native module access | [expo-horizon-core README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-core/README.md) |
| Location migration, limitations, feature support matrix | [expo-horizon-location README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-location/README.md) |
| Push notifications migration, token types, feature support matrix | [expo-horizon-notifications README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-notifications/README.md) |
| Example app wiring for all three packages | [expo-horizon example README](https://github.com/software-mansion-labs/expo-horizon/blob/main/example/README.md) |
| Panel sizing guidelines (dp values, orientation, letterboxing) | [Meta Panel Sizing](https://developers.meta.com/horizon/documentation/android-apps/panel-sizing) |
| Meta Horizon Store manifest checklist for publishing | [Publish Mobile Manifest](https://developers.meta.com/horizon/resources/publish-mobile-manifest/) |
| Device management, casting, sideloading, ADB | [Meta Quest Developer Hub](https://developers.meta.com/horizon/documentation/android-apps/meta-quest-developer-hub) |
| Server-side push delivery via Meta's push service | [Horizon OS push notifications](https://developers.meta.com/horizon/documentation/android-apps/ps-user-notifications/) |
