---
name: view-screenshot
description: "Best practices for using the view_screenshot tool in Radon IDE. Captures a screenshot of the app running in the development viewport. Use to inspect the current visual state of the app, verifying UI changes, checking layout, confirming styling, or diagnosing visual bugs. Trigger on: 'screenshot', 'what does the app look like', 'show me the screen', 'UI preview', 'check the layout', 'see the device', 'view the app', or any request to visually inspect the running application."
---

# view_screenshot

Captures a screenshot of the app running in the Radon IDE development viewport.

## When to use

- After UI/styling changes to confirm they rendered correctly.
- After `reload_application` to verify the app recovered.
- When diagnosing layout or visual issues.

## Key rules

- `view_application_logs` already includes a screenshot if the preview is ready - if logs output is already available, a separate screenshot call is redundant.
- For visual-only checks without logs, `view_screenshot` is the lighter option.
- The screenshot reflects current device orientation, dark/light mode, and content size settings.

## Error handling

- **Device off:** turn on the Radon IDE emulator.
