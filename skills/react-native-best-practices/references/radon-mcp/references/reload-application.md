---
name: reload-application
description: "Best practices for using the reload_application tool in Radon IDE. Triggers a reload of the app running in the development emulator with three methods: reloadJs (JS bundle reload), restartProcess (native process restart), and rebuild (full native rebuild). Use when debugging state issues, after code changes, when the app crashes, or when native changes require a rebuild. Trigger on: 'reload app', 'restart app', 'rebuild', 'hot reload', 'reset state', 'app crashed', 'app frozen', 'refresh', 'reloadJs', 'restartProcess', or any request to restart or reload the running application."
---

# reload_application

Reloads the app in the Radon IDE emulator. Three methods with escalating scope:

## Input schema:

```
{ reloadMethod: "reloadJs" | "restartProcess" | "rebuild" }
```

## Reload methods

| Method           | Speed   | Use when                                                                   |
| ---------------- | ------- | -------------------------------------------------------------------------- |
| `reloadJs`       | Seconds | JS/TS changes, reset React state, HMR missed a change                      |
| `restartProcess` | Seconds | Native module in bugged state, need cold start, `reloadJs` didn't fix it   |
| `rebuild`        | Minutes | Changed native code, added/removed native dependency, changed build config |

**Always start with the lightest method and escalate, unless you are already certain that a more powerful restart method is required:** `reloadJs` -> `restartProcess` -> `rebuild`.

## Key behaviors

- The tool **waits until the app is running** before returning - safe to call `view_screenshot` or other tools immediately after.
- After reloading, use `view_screenshot` or `view_application_logs` to confirm the app is healthy.
- `reloadJs` resets all React state, context providers, and in-memory stores.

## Error handling

- **Reload failure:** check `view_application_logs` for build or runtime errors.
- **Device off / Radon IDE not launched:** request the user to turn on the Radon IDE emulator.
