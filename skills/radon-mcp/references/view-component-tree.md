---
name: view-component-tree
description: "Best practices for using the view_component_tree tool in Radon IDE. Displays the React component tree (view hierarchy) of the running app, filtered to user-authored components. Use for understanding mounted component structure, resolving layout issues, finding context providers, or mapping file structure to component hierarchy. Trigger on: 'component tree', 'view hierarchy', 'component structure', 'React tree', 'mounted components', 'layout structure', 'where is this component', 'context provider location', 'component hierarchy', or any request to understand the UI structure of the running app."
---

# view_component_tree

Displays the React component tree of the running app, filtered to user-authored components with source file locations.

## Output format

Indented XML-like tree with component names, source paths, line numbers, and HOC descriptors:

```xml
<MyComponent> (src/screens/HomeScreen.tsx:42)
  <Header> [memo, forwardRef] (src/components/Header.tsx:15)
    <Title> (src/components/Title.tsx:8)
      Hello World
    </Title>
  </Header>
  <Icon /> (src/components/Icon.tsx:5)
</MyComponent>
```

## Key behaviors

- **Only user-authored components** are shown - `node_modules` components (View, Text, ScrollView, third-party) are excluded.
- **Only currently mounted components** are visible - conditionally rendered components with `false` condition won't appear.
- Paths are relative to the workspace root.

## Best practices

- **Pair with `view_screenshot`** for layout debugging: screenshot shows the visual problem, tree reveals the structure causing it.
- Use source file paths from the tree to navigate directly to relevant code.

## Error handling

- **Corrupted tree:** reload the app with `reload_application({ reloadMethod: "reloadJs" })`, then retry.
- **Device off:** request the user to turn on the Radon IDE emulator.
