# Push Notifications Migration

Migrating from `expo-notifications` to `expo-horizon-notifications` for Meta Quest compatibility.

For the full README, webfetch the [expo-horizon-notifications README](https://github.com/software-mansion-labs/expo-horizon/blob/main/expo-horizon-notifications/README.md).

---

## Migration Steps

1. Install `expo-horizon-core` first with a `horizonAppId` (see `setup.md`). Push notifications **require** `horizonAppId`.

2. Replace the package:

```bash
npx expo install expo-horizon-notifications
npm uninstall expo-notifications
```

3. Update your Expo config:

```json
{
  "plugins": [
    [
      "expo-horizon-core",
      {
        "horizonAppId": "your-horizon-app-id",
        "supportedDevices": "quest2|quest3|quest3s"
      }
    ],
    "expo-horizon-notifications"
  ]
}
```

4. Update all imports:

```typescript
// Before
import * as Notifications from 'expo-notifications';

// After
import * as Notifications from 'expo-horizon-notifications';
```

5. Run `npx expo prebuild --clean` to regenerate native files.

---

## How It Works

`expo-horizon-notifications` is a drop-in replacement for `expo-notifications`. The implementation is selected by build variant:

- **`mobile` build** -- Standard `expo-notifications` behavior using Firebase Cloud Messaging
- **`quest` build** -- Meta Horizon push notification service

On iOS, behavior is identical to `expo-notifications`.

---

## Push Token Differences

On Quest, `getDevicePushTokenAsync` returns a token with type `"horizon"` instead of the standard `"android"`:

```typescript
import * as Notifications from 'expo-horizon-notifications';

const token = await Notifications.getDevicePushTokenAsync();
// On Quest: { type: 'horizon', data: '...' }
// On mobile: { type: 'android', data: '...' }
```

Send the token to your backend server, which delivers push notifications via Meta's push service. For server-side implementation, webfetch the [Horizon OS push notification docs](https://developers.meta.com/horizon/documentation/android-apps/ps-user-notifications/).

---

## Quest Notification Limitations

- **No Expo Push Service** -- `getExpoPushTokenAsync` is unsupported on Quest. Use `getDevicePushTokenAsync` with Meta's push service directly.
- **No badge counts** -- `getBadgeCountAsync` and `setBadgeCountAsync` are unsupported (underlying ShortcutBadger library lacks Quest support).
- **No unregister** -- `unregisterForNotificationsAsync` is unsupported on Quest.
- **Notification channels and interactive categories** -- Not yet tested on Quest.

### Guard unsupported features

```typescript
import ExpoHorizon from 'expo-horizon-core';
import * as Notifications from 'expo-horizon-notifications';

// Use device push token (works on both platforms)
const token = await Notifications.getDevicePushTokenAsync();

// Expo Push Token is mobile-only
if (!ExpoHorizon.isHorizonDevice) {
  const expoPushToken = await Notifications.getExpoPushTokenAsync();
}

// Badge counts are mobile-only
if (!ExpoHorizon.isHorizonDevice) {
  await Notifications.setBadgeCountAsync(5);
}
```

---

## Feature Support Matrix

| Function | Quest | Notes |
|---|---|---|
| `addPushTokenListener` | Yes | Requires `horizonAppId` |
| `getDevicePushTokenAsync` | Yes | Returns `{ type: 'horizon', data: '...' }`. Requires `horizonAppId` |
| `getExpoPushTokenAsync` | No | Expo Push Service not supported |
| `addNotificationReceivedListener` | Yes | |
| `addNotificationResponseReceivedListener` | Yes | |
| `addNotificationsDroppedListener` | Yes | |
| `useLastNotificationResponse` | Yes | |
| `setNotificationHandler` | Yes | |
| `registerTaskAsync` | Yes | |
| `unregisterTaskAsync` | Yes | |
| `getPermissionsAsync` | Yes | |
| `requestPermissionsAsync` | Yes | |
| `getBadgeCountAsync` | No | ShortcutBadger lacks Quest support |
| `setBadgeCountAsync` | No | ShortcutBadger lacks Quest support |
| `scheduleNotificationAsync` | Yes | |
| `cancelScheduledNotificationAsync` | Yes | |
| `cancelAllScheduledNotificationsAsync` | Yes | |
| `getAllScheduledNotificationsAsync` | Yes | |
| `getNextTriggerDateAsync` | Yes | |
| `dismissNotificationAsync` | Yes | |
| `dismissAllNotificationsAsync` | Yes | |
| `getPresentedNotificationsAsync` | Yes | |
| `getLastNotificationResponse` | Yes | |
| `clearLastNotificationResponse` | Yes | |
| `unregisterForNotificationsAsync` | No | |

---

## Version Compatibility

| `expo-horizon-notifications` | `expo-notifications` | Expo SDK |
|---|---|---|
| 55.0.0 | 55.0.10 | 55 |
| 0.0.9-0.0.11 | 19.0.7 | 54 |
