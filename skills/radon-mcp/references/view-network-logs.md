---
name: view-network-logs
description: "Best practices for using the view_network_logs tool in Radon IDE. Returns a paginated list of all network requests made by the app, including method, URL, status, duration, and size. Use when debugging networking issues, inspecting API calls, or verifying that the app communicates correctly with backend services. Trigger on: 'network logs', 'network requests', 'API calls', 'HTTP requests', 'network inspector', 'fetch requests', 'network traffic', 'request failed', 'status code', '404', '500', 'CORS', or any request to inspect network activity from the running app."
---

# view_network_logs

Returns a paginated list of all network requests (method, URL, status, duration, size). 50 entries per page.

## Input schema:

```
{ pageIndex: "latest" | "<number>" }
```

## Output format

```
=== NETWORK LOGS (page 3/5) ===
{id: abc123} "GET https://api.example.com/users" 200 OK json 1.2kB 150ms
```

The `requestId` (e.g., `abc123`) is used to inspect further details with `view_network_request_details`.

## Key rules

- **Start with `"latest"`** for current issues - almost always the right starting point.
- Use numeric page indexes (`"0"`, `"1"`, ...) to search older requests, with `0` being the oldest requests page.
- **Follow up with `view_network_request_details`** to inspect headers, body, and metadata of a specific request that could be problematic.
- If no logs returned, the network inspector may not be enabled - it must be active for requests to be captured.

## Error handling

- **Network inspector unavailable:** ensure the Network panel is available in Radon IDE.
- **Invalid/out-of-range page index:** use `"latest"` or check the page header for valid range.
- **Device off:** request the user to turn on the Radon IDE emulator.
