---
name: view-network-request-details
description: "Best practices for using the view_network_request_details tool in Radon IDE. Returns full details of a specific network request including headers, body, and metadata. Use after view_network_logs to drill into a specific request for debugging. Trigger on: 'request details', 'request headers', 'response body', 'request body', 'inspect request', 'API response', or any follow-up to viewing network logs where full details of a particular request are needed."
---

# view_network_request_details

Returns full details of a specific network request (headers, body, metadata). **Always call `view_network_logs` first** to obtain the `requestId`.

```
view_network_request_details({ requestId: "<string>" })
```

## Workflow

1. `view_network_logs({ pageIndex: "latest" })` — find the problematic request.
2. Note the `requestId` from the output.
3. `view_network_request_details({ requestId: "<id>" })` — inspect full details.

## Key behaviors

- **Sensitive headers are redacted** (case-insensitive match on: `auth`, `cookie`, `token`, `secret`, `key`, `session`, `credential`). You can still check header presence/absence and non-sensitive headers.
- **Response bodies > 1000 chars are truncated** with a placeholder showing MIME type and original size. For full payload, direct the user to the Radon IDE Network panel.

## Error handling

- **Request ID not found:** verify it was copied correctly from `view_network_logs`.
- **Device off:** turn on the Radon IDE emulator.
