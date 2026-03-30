# Location Services Migration

Migrating from `expo-location` to `expo-horizon-location` for Meta Quest compatibility.

For the full README, webfetch the [expo-horizon-location README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-location/README.md).

---

## Migration Steps

1. Install `expo-horizon-core` first (see `setup.md`)

2. Replace the package:

```bash
npx expo install expo-horizon-location
npm uninstall expo-location
```

3. Update your Expo config -- replace `expo-location` plugin with `expo-horizon-location`:

```json
{
  "plugins": [
    ["expo-horizon-core", { "supportedDevices": "quest2|quest3|quest3s" }],
    "expo-horizon-location"
  ]
}
```

4. Update all imports:

```typescript
// Before
import * as Location from 'expo-location';

// After
import * as Location from 'expo-horizon-location';
```

5. Run `npx expo prebuild --clean` to regenerate native files.

---

## How It Works

`expo-horizon-location` is a drop-in replacement for `expo-location`. It provides two implementations selected automatically by build variant:

- **`mobile` build** -- Standard `expo-location` behavior using Google Play Services
- **`quest` build** -- Meta Horizon-compatible implementation without Google Play Services, using the Android network provider

On iOS, behavior is identical to `expo-location`.

---

## Quest Location Limitations

Meta Quest devices lack GPS hardware and certain sensors. The network provider is the only available location source.

### Key constraints

- **No GPS** -- `getCurrentPositionAsync` and `watchPositionAsync` use the network provider regardless of the accuracy setting. Network updates occur no more frequently than every ~10 minutes.
- **No heading/compass** -- `watchHeadingAsync` and `getHeadingAsync` are unsupported (no magnetic or accelerometer sensors).
- **No geocoding** -- `geocodeAsync` and `reverseGeocodeAsync` are unsupported (Android `Geocoder` is not present on Quest).
- **No background location** -- `startLocationUpdatesAsync`, `stopLocationUpdatesAsync`, geofencing APIs are unsupported. Meta Horizon Store prohibits `ACCESS_BACKGROUND_LOCATION`.

### Guard unsupported features

```typescript
import ExpoHorizon from 'expo-horizon-core';
import * as Location from 'expo-horizon-location';

// Heading is only available on mobile
if (!ExpoHorizon.isHorizonDevice) {
  const subscription = await Location.watchHeadingAsync((heading) => {
    console.log('Heading:', heading.magHeading);
  });
}

// Geocoding is only available on mobile
if (!ExpoHorizon.isHorizonDevice) {
  const results = await Location.geocodeAsync('1 Hacker Way, Menlo Park, CA');
}
```

---

## Feature Support Matrix

| Function | Android | Quest | Notes |
|---|---|---|---|
| `enableNetworkProviderAsync` | Yes | Yes | |
| `getProviderStatusAsync` | Yes | Yes | |
| `hasServicesEnabledAsync` | Yes | Yes | |
| `requestForegroundPermissionsAsync` | Yes | Yes | |
| `requestBackgroundPermissionsAsync` | Yes | Yes | |
| `getForegroundPermissionsAsync` | Yes | Yes | |
| `getBackgroundPermissionsAsync` | Yes | Yes | |
| `getCurrentPositionAsync` | Yes | Yes | Network provider only on Quest, ~10 min update interval |
| `watchPositionAsync` | Yes | Yes | Network provider only on Quest, ~10 min update interval |
| `getLastKnownPositionAsync` | Yes | Yes | |
| `watchHeadingAsync` | Yes | No | No magnetic/accelerometer sensors on Quest |
| `getHeadingAsync` | Yes | No | No magnetic/accelerometer sensors on Quest |
| `geocodeAsync` | Yes | No | No Geocoder on Quest |
| `reverseGeocodeAsync` | Yes | No | No Geocoder on Quest |
| `startGeofencingAsync` | Yes | No | No `ACCESS_BACKGROUND_LOCATION` on Quest |
| `stopGeofencingAsync` | Yes | No | No `ACCESS_BACKGROUND_LOCATION` on Quest |
| `startLocationUpdatesAsync` | Yes | No | No `ACCESS_BACKGROUND_LOCATION` on Quest |
| `stopLocationUpdatesAsync` | Yes | No | No `ACCESS_BACKGROUND_LOCATION` on Quest |

---

## Version Compatibility

| `expo-horizon-location` | `expo-location` | Expo SDK |
|---|---|---|
| 55.0.0 | 55.1.2 | 55 |
| 0.0.4-0.0.5 | 18.1.17 | 54 |
