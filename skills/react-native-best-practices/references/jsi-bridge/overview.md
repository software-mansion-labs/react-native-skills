# JSI Overview

JSI (JavaScript Interface) is a C++ header-only API that provides a language-agnostic abstraction over a JavaScript engine. It replaces the old React Native async message bridge with direct, synchronous C++ calls into the engine.

---

## Where JSI Sits

```
Your C++ code
      │
      ▼
  jsi::Runtime  ← the abstraction layer (jsi.h)
      │
      ▼
JS Engine (Hermes / V8 / JSC)
      │
      ▼
  Your JavaScript
```

`jsi::Runtime` is the entry point. Every JSI operation goes through a `Runtime&` reference. The concrete implementation (Hermes, V8, JSC) is behind the abstraction — your binding code is engine-agnostic.

---

## Why Everything Goes Through `rt.global()`

A fresh `Runtime` is a blank engine with no bindings. It doesn't know about your C++ code. To make a C++ function or object callable from JS, you must install it explicitly on the global object:

```cpp
// Make a C++ function available as globalThis.myNativeFunction
auto fn = jsi::Function::createFromHostFunction(
    rt,
    jsi::PropNameID::forAscii(rt, "myNativeFunction"),
    1, // paramCount hint
    [](jsi::Runtime& rt, const jsi::Value&, const jsi::Value* args, size_t) {
      // ...
      return jsi::Value::undefined();
    });

rt.global().setProperty(rt, "myNativeFunction", std::move(fn));
```

Nothing is injected automatically. Every binding is opt-in via `rt.global()`.

---

## Sync vs Async Execution Model

JSI calls are **synchronous**. When C++ calls into JS (or JS calls into C++ via a HostFunction), both sides execute on the same thread and control returns to the caller before anything else runs. This is fundamentally different from the old RN bridge, which queued messages asynchronously across threads.

Consequences:
- JSI bindings are faster — no serialization, no thread hops for the call itself.
- The JS thread can be blocked by a slow HostFunction. Don't do heavy work synchronously in a HostFunction; return a Promise and use `CallInvoker` to resolve it from a background thread.
- `evaluateJavaScript` is also synchronous — it blocks until the script finishes.

---

## Why Some Methods Take `Runtime&` and Others Don't

Types like `jsi::Value`, `jsi::Object`, `jsi::String`, and `jsi::PropNameID` are thin wrappers around an opaque `PointerValue*` handle. The handle alone doesn't contain the data — it's a reference into the engine's heap.

Any operation that needs to read, write, or allocate within the engine must pass `Runtime&` so JSI can dispatch to the correct engine implementation. Simple operations that only manipulate the wrapper object in C++ memory (like moving or destroying it) don't need `Runtime&`.

```cpp
jsi::Value v = rt.global().getProperty(rt, "x"); // Runtime& required: reads from engine

bool isNum = v.isNumber();   // no Runtime& needed: inspects the C++ kind tag
double n   = v.getNumber();  // no Runtime& needed: reads from C++ union
```

---

## Pure Runtime Has No Event Loop

`jsi::Runtime` is a JavaScript engine, not a full JS environment. It has no event loop, no timer queue, no I/O. Microtasks (Promises) don't drain on their own unless you drain them explicitly.

### `queueMicrotask`

Enqueues a JS function as a microtask in the engine's internal Job queue:

```cpp
rt.queueMicrotask(callbackFunction);
```

### `drainMicrotasks`

Runs pending microtasks (Promise jobs). Returns `true` when the queue is fully drained, `false` when there's more work but the hint limit was reached:

```cpp
// Drain all pending microtasks
while (!rt.drainMicrotasks()) {}

// Or drain at most 10 microtasks per call
rt.drainMicrotasks(10);
```

React Native's event loop calls `drainMicrotasks` on your behalf. You only need to call it yourself if you're hosting a bare `Runtime` outside of RN (e.g., in tests or a standalone C++ app).

---

## `setRuntimeData` / `getRuntimeData`

A runtime can carry arbitrary C++ state keyed by UUID. This is the JSI equivalent of a thread-local: store shared resources (connection pools, config objects) on the runtime once, retrieve them anywhere that has a `Runtime&`.

```cpp
// Define a stable UUID for your data
static constexpr jsi::UUID kMyConfigUUID{
    0x12345678, 0x1234, 0x5678, 0x9abc, 0xdef012345678};

// Store
rt.setRuntimeData(kMyConfigUUID, std::make_shared<MyConfig>(config));

// Retrieve
auto config = std::static_pointer_cast<MyConfig>(rt.getRuntimeData(kMyConfigUUID));
```

When the runtime is destroyed or the key is overwritten, it releases ownership of the stored object.

---

## `evaluateJavaScript` — Use Sparingly

`evaluateJavaScript` compiles and runs a JS buffer. It's the right tool for bootstrapping a bundle, but expensive for anything you could do directly through the JSI API:

```cpp
// Slow: parses and executes JS to call a function
rt.evaluateJavaScript(
    std::make_shared<jsi::StringBuffer>("JSON.stringify(value)"), "");

// Fast: call the function directly through JSI
auto jsonStringify = rt.global()
    .getPropertyAsObject(rt, "JSON")
    .getPropertyAsFunction(rt, "stringify");
auto result = jsonStringify.call(rt, value);
```

Rule of thumb: if you already know the function you want to call, use JSI APIs. Reserve `evaluateJavaScript` for code that can't be expressed as a direct API call.

---

## Prototype Manipulation

JSI exposes prototype operations for cases where you need to set up a prototype chain in C++:

```cpp
// Create an object with a custom prototype
jsi::Object obj = jsi::Object::create(rt, protoValue);

// Read/write an existing object's prototype
jsi::Value proto = obj.getPrototype(rt);
obj.setPrototype(rt, newProtoValue); // throws if unsuccessful
```

These are rarely needed in typical bindings but useful when implementing class hierarchies that need to be visible to JS `instanceof` checks.
