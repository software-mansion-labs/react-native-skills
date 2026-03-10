---
name: get-library-description
description: "Best practices for using the get_library_description tool in Radon IDE. Returns a detailed description of an npm library and its use cases. Use when evaluating whether to use a specific library, understanding what a dependency does, or recommending libraries for a particular task. Trigger on: 'what does X package do', 'should I use', 'library description', 'library info', 'describe library', 'library use cases', or any request to understand the purpose and capabilities of a specific npm package."
---

# get_library_description

Returns a description of an npm library and its use cases from the Radon AI backend.

## Input schema:

```
{ library_npm_name: "<npm package name>" }
```

## Key rules

- Use the **exact npm package name** as published on the registry (e.g., `"react-native-reanimated"`, `"@react-navigation/native"`).
- This tool is for **evaluation, not implementation**. For API docs and usage examples, follow up with `query_documentation`.
- Scoped to the React Native / Expo ecosystem.

## Typical workflow

1. `get_library_description` - understand what the library does and whether it fits.
2. `query_documentation` - get specific API docs and usage patterns.

## Error handling

- **Invalid/missing license:** ensure the user has an active Radon IDE license.
- **Network failure:** ensure internet connectivity.
