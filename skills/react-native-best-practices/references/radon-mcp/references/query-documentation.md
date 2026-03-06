---
name: query-documentation
description: "Best practices for using the query_documentation tool in Radon IDE. Returns documentation snippets relevant to a provided query from a curated React Native and Expo knowledge base. Use when accurate, up-to-date documentation is needed for React Native APIs, Expo modules, or related libraries. Trigger on: 'how to use', 'documentation for', 'API reference', 'React Native docs', 'Expo docs', 'how does X work in React Native', 'what is the API for', or any question about React Native or Expo library usage that benefits from authoritative documentation."
---

# query_documentation

Returns documentation snippets from a curated React Native / Expo knowledge base on the Radon AI backend.

```
query_documentation({ text: "<query string>" })
```

## Key rules

- Write **specific, focused queries** targeting a concrete API or feature.
- Include the **library name** when relevant (e.g., `"react-native-reanimated useSharedValue hook"`).
- **Trust returned docs over training data** for version-specific details — RN/Expo APIs change between versions.
- Query docs **before implementing** features that rely on RN/Expo APIs.
- For high-level library evaluation, use `get_library_description` instead.

## Good vs bad queries

```
// Good
query_documentation({ text: "React Navigation stack navigator configuration options" })
query_documentation({ text: "expo-image-picker launchImageLibraryAsync options" })

// Too vague
query_documentation({ text: "React Native" })
query_documentation({ text: "navigation" })
```

## Error handling

- **Invalid/missing license:** returns a license-required message. Ensure the user has an active Radon IDE license.
- **Network failure:** ensure internet connectivity.
