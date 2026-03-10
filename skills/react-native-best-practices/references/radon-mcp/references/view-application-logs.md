---
name: view-application-logs
description: "Best practices for using the view_application_logs tool in Radon IDE. Returns all build, bundling, and runtime logs from the running app. Use when the user has any issue with the app, when builds are failing, when there are runtime errors or crashes, or whenever logs would help diagnose a problem. Trigger on: 'app logs', 'build logs', 'console errors', 'app crash', 'build failure', 'Metro error', 'runtime error', 'debug logs', 'what went wrong', 'native logs', 'JS logs', or any debugging scenario where log output would be useful."
---

# view_application_logs

Returns all build, bundling, and runtime logs plus an optional PNG screenshot. **Call this first when the user reports any issue.**

## Log sections

| Section                | Look here for                                                                       |
| ---------------------- | ----------------------------------------------------------------------------------- |
| `BUILD PROCESS LOGS`   | Compilation errors, linking failures, missing native deps, code signing (iOS)       |
| `JS PACKAGER LOGS`     | Dependency installation issues, package resolution errors                           |
| `METRO LOGS`           | Bundle compilation errors, module resolution, syntax errors, transform errors       |
| `NATIVE-SIDE APP LOGS` | Native crash reports, native module init errors, memory warnings                    |
| `JS-SIDE APP LOGS`     | Application exceptions, React errors, state management issues, `console.log` output |

Missing sections mean no output on that channel - this is normal.

## Key rules

- **Read all sections before diagnosing** - issues often span layers (e.g., native crash triggered by JS error).
- Build-time issues: focus on BUILD PROCESS + JS PACKAGER. Runtime: focus on NATIVE-SIDE + JS-SIDE.
- Based on findings, follow up with `view_component_tree` (layout), `view_network_logs` (API), or `reload_application` (reset).

## Error handling

- **No build run:** select a project and device in the Radon IDE panel.
- **Device off / Radon IDE not launched:** request the user to turn on the Radon IDE emulator.
