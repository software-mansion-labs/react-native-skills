# Setup and Library Templates

---

## Installing a JSI Binding in React Native

A JSI binding must be installed into the `Runtime` **synchronously** during module initialization — before any JS runs. The installation point differs by platform.

### Android — `JSIModulePackage`

Override `getJSIModulePackage` in your `ReactApplication` (or use a `TurboReactPackage`):

```java
// MainApplication.java / MainApplication.kt
@Override
protected ReactNativeHost createReactNativeHost() {
  return new DefaultReactNativeHost(this) {
    @Override
    public boolean isNewArchEnabled() { return BuildConfig.IS_NEW_ARCHITECTURE_ENABLED; }
  };
}
```

On the C++ side, implement `JSIModuleProvider`:

```cpp
// MyJSIModule.cpp
void install(jsi::Runtime& rt) {
  rt.global().setProperty(
      rt, "myNativeModule",
      jsi::Object::createFromHostObject(rt, std::make_shared<MyHostObject>()));
}
```

For JNI glue between Java/Kotlin and your C++, use **FBJNI** — Meta's JNI helper library that's already a dependency of React Native:

```cpp
#include <fbjni/fbjni.h>

struct MyModule : public facebook::jni::JavaClass<MyModule> {
  static constexpr auto kJavaDescriptor = "Lcom/example/MyModule;";

  static void install(jsi::Runtime& rt) {
    // ...
  }
};
```

FBJNI handles JNI reference management, exception translation, and type-safe method lookup — prefer it over raw JNI.

### iOS — ObjC++ (`.mm`)

Install from your native module's `install` method, called from the `RCTCxxBridge` setup:

```objc
// MyModule.mm  (must be .mm for ObjC++ to mix with C++)
#import <React/RCTBridge+Private.h>
#import <jsi/jsi.h>

@implementation MyModule

RCT_EXPORT_MODULE()

- (void)setBridge:(RCTBridge *)bridge {
  RCTCxxBridge *cxxBridge = (RCTCxxBridge *)bridge;
  if (!cxxBridge.runtime) return;

  auto &rt = *(jsi::Runtime *)cxxBridge.runtime;
  MyBinding::install(rt);
}

@end
```

Use `.mm` extension (ObjC++) so you can mix Objective-C with C++ headers. Pure `.m` files cannot include C++ headers.

---

## Library Templates

| Template | Command / URL | Notes |
|----------|--------------|-------|
| **create-react-native-library** *(recommended)* | `npx create-react-native-library@latest MyLib` → choose *"C++ for both iOS & Android"* | Actively maintained, sets up TurboModule + JSI boilerplate, CMakeLists, podspec |
| mrousavy/react-native-jsi-library-template | GitHub: `mrousavy/react-native-jsi-library-template` | Minimal JSI-only template, good for learning the wiring |
| ospfranco/react-native-jsi-template | GitHub: `ospfranco/react-native-jsi-template` | Another minimal template, sqlite-focused examples |
| ammarahm-ed/react-native-jsi-template | GitHub: `ammarahm-ed/react-native-jsi-template` | Simple JSI setup without TurboModule overhead |

Start with `create-react-native-library` unless you have a specific reason to use a minimal template — it handles the platform glue that's easy to get wrong.

---

## JSI vs TurboModules vs Nitro Modules

| | **Raw JSI** | **TurboModules** | **Nitro Modules** |
|-|------------|-----------------|------------------|
| **What it is** | Bare C++ engine API | RN's official native module system, built on JSI | Third-party alternative to TurboModules (mrousavy) |
| **Type safety** | Manual — you check types at runtime | Schema-generated from TypeScript spec | Schema-generated from TypeScript; stricter types |
| **Codegen** | None | Yes (`react-native-codegen`) | Yes (Nitrogen codegen) |
| **Async support** | Manual (Promise + CallInvoker) | Built-in via Promise return type | Built-in |
| **Performance** | Maximum — you control everything | Good | Slightly faster than TurboModules (fewer layers) |
| **Maintenance burden** | High — all boilerplate manual | Low — mostly generated | Low — mostly generated |
| **Best for** | Libraries with unusual requirements, learning JSI internals | Standard native modules shipping with React Native | Third-party libraries targeting maximum performance |

For most library authors, **TurboModules** (via `create-react-native-library`) is the right choice. Use raw JSI only when you need capabilities that TurboModules don't expose — such as `HostObject` with intercepted property access, custom `ArrayBuffer` providers, or a custom runtime decorator.
