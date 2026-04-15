# C++ Memory Patterns for JSI

JavaScript hides memory management. C++ exposes it. Every `&`, `*`, and smart pointer in JSI code is an explicit statement about ownership — who created this data, who is responsible for freeing it, and how long it should live. This reference maps those statements to concepts you already know from JavaScript.

---

## Stack vs Heap

In JavaScript, you never choose where a value lives. The engine decides. In C++, you do.

**Stack** — fast, automatic, scoped. Variables on the stack are destroyed when the enclosing function returns. No GC, no cleanup code. The closing brace is the finalizer.

```cpp
void process() {
    int count = 0;            // stack: destroyed when process() returns
    std::string key = "abc";  // stack: same (string contents may be heap-allocated
                              //        internally, but that's the std::string's problem)
    double result = count * 2.0;
} // ← count, key, result are all freed here. Deterministic. Instant.
```

**Heap** — manual lifetime, required when data must outlive its creating function. In C++, heap allocation without a smart pointer means you own the memory and must `delete` it:

```cpp
// Raw new/delete — never write this in modern JSI code
void bad() {
    int* data = new int[1024]; // heap: lives until you delete it
    // ... if an exception throws here, delete never runs, bytes leak forever
    delete[] data;
}
```

The JS analogy: stack ≈ primitive values (scoped to the expression), heap ≈ objects (reference-counted by GC). The difference: C++ heap is freed *immediately and deterministically* when ownership drops to zero. JS heap is freed *eventually* during a GC pass.

For JSI modules: local variables and function arguments live on the stack. HostObjects, callbacks, and anything passed to JavaScript live on the heap, managed by smart pointers.

---

## `unique_ptr` — Sole Ownership

`unique_ptr` owns its heap object exclusively. One owner at a time. When the `unique_ptr` is destroyed, the object is freed. You cannot copy it — only move it (transfer ownership).

```cpp
#include <memory>

// Creating
auto buffer = std::make_unique<AudioBuffer>(1024); // heap-allocated, owned by `buffer`
buffer->fill(0.0f);

// Transferring ownership
auto other = std::move(buffer); // `buffer` is now nullptr, `other` owns the AudioBuffer
// auto copy = buffer;          // COMPILER ERROR — unique_ptr is not copyable

} // `other` destroyed → AudioBuffer freed automatically
```

Think of `unique_ptr` like a non-shareable `const` binding in JS — conceptually, only one thing holds the value at a time. If you want to pass it somewhere, you hand it over entirely and lose access.

**Use `unique_ptr` when:**
- You have one clear owner and no one else needs the object
- The object lives inside another C++ object (member variable)
- You're moving through a pipeline (produce → transform → consume)

**Not for JSI HostObjects** — JavaScript's GC and your C++ code both need access, which requires shared ownership. Use `shared_ptr` for anything exposed to JS.

---

## `shared_ptr` — Reference-Counted Ownership

`shared_ptr` allows multiple owners. It maintains an internal reference count — every copy of the `shared_ptr` increments it, every destruction decrements it. When the count reaches zero, the object is freed.

```cpp
#include <memory>

auto config = std::make_shared<AppConfig>(); // refcount = 1

{
    auto a = config; // refcount = 2 — both `config` and `a` own AppConfig
    auto b = config; // refcount = 3

    // `a` and `b` go out of scope here → refcount drops to 1
}

// `config` goes out of scope → refcount = 0 → AppConfig destructor runs
```

This is the closest C++ equivalent to JavaScript's GC semantics: the object lives as long as anyone holds a reference. The difference is that `shared_ptr` uses *deterministic* reference counting (freed immediately when count hits zero) rather than a tracing collector.

**`.get()` for borrowing** — when a function only needs to use the object, not share ownership, pass the raw pointer with `.get()`. The caller does not extend the object's lifetime:

```cpp
void render(AudioBuffer* buf); // takes a raw pointer — doesn't own, doesn't extend lifetime

auto buffer = std::make_shared<AudioBuffer>(1024);
render(buffer.get()); // borrows without incrementing refcount
// buffer still owns the AudioBuffer after render() returns
```

Raw `.get()` is safe only when the `shared_ptr` is guaranteed to outlive the function call. Never store a raw pointer obtained from `.get()` beyond the immediate call.

**`std::make_shared` vs `new`** — always prefer `make_shared`. It allocates the object and its control block (the refcount) in a single allocation, which is faster and avoids a subtle exception-safety hole with raw `new`.

```cpp
// Preferred
auto db = std::make_shared<Database>(path);

// Avoid — two allocations, exception-unsafe in some call contexts
auto db = std::shared_ptr<Database>(new Database(path));
```

---

## `std::move` — Transferring Ownership Without Copying

In JavaScript, assigning an object doesn't copy it — both variables point to the same object. In C++, assignment copies by default. `std::move` opts out of copying and transfers the object's internal state to a new owner, leaving the source in a valid but empty state.

```cpp
std::vector<uint8_t> data = {1, 2, 3, 4, 5};

// Copy — allocates new memory, duplicates 5 bytes
std::vector<uint8_t> copy = data;    // data still intact, copy is independent

// Move — transfers the internal buffer pointer, no allocation
std::vector<uint8_t> moved = std::move(data); // data is now empty, moved owns the buffer
```

In JSI code, `std::move` appears in two recurring patterns:

**Installing a function or object into the runtime:**

```cpp
auto fn = jsi::Function::createFromHostFunction(rt, name, 0, callback);
rt.global().setProperty(rt, "myFunc", std::move(fn));
// fn is now empty — the runtime owns the function
```

**Moving a `shared_ptr` into a lambda capture** (transfers the pointer without incrementing refcount, then immediately decrements the source):

```cpp
auto db = std::make_shared<Database>(path); // refcount = 1

auto hostFn = jsi::Function::createFromHostFunction(rt, name, 1,
    [db = std::move(db)](jsi::Runtime& rt, ...) -> jsi::Value { // refcount still = 1
        return db->query(...);
    });
// db (local) is now nullptr. The lambda owns the only reference.
```

`std::move` does not actually move anything at the call site — it's a cast to an rvalue reference that tells the compiler "the move constructor is permitted here." The actual work happens in the receiving type's move constructor.

---

## Lambda Captures

C++ lambdas are closures, but unlike JavaScript closures, you explicitly choose what to capture and how.

```cpp
[capture list](parameters) -> return_type { body }
```

The four capture forms:

```cpp
int x = 10;
std::string name = "JSI";

auto byValue   = [x]()    { return x; };        // snapshot of x at capture time
auto byRef     = [&x]()   { return x; };        // live alias — sees mutations to x
auto allByVal  = [=]()    { return x + 1; };    // copies every variable referenced in body
auto allByRef  = [&]()    { return name; };     // references everything — closest to JS default
```

| Syntax | What you get | JS analogy |
|--------|-------------|------------|
| `[x]` | Copy of `x` at the moment of capture | `const snapshot = x` before the closure |
| `[&x]` | Reference to `x` — sees changes, can modify | Shared variable binding (but can dangle — see below) |
| `[this]` | Raw pointer to the enclosing class instance | Unguarded `this` in a callback |
| `[=]` | Copies all referenced variables | No direct analogy |
| `[&]` | References all referenced variables | Default JS closure behavior |

**Which is safe for async / GC contexts:**

`[x]` (capture by value) is the safe default for JSI lambdas stored beyond the current call. When `x` is a `shared_ptr`, capturing by value copies the pointer and increments the refcount — the lambda extends the object's lifetime.

`[&x]` is only safe when the lambda executes synchronously before `x` goes out of scope. Never capture a local variable by reference in a lambda that is stored, passed to a background thread, or returned to JavaScript.

`[this]` is unsafe in any function returned from `HostObject::get`. JavaScript can detach the method, drop the object, and later call the detached function — at which point `this` is a dangling pointer into freed memory. Use `weak_from_this()` instead (see the GC Boundary section below).

```cpp
// WRONG — [&db] dangles after install() returns
void install(jsi::Runtime& rt, std::shared_ptr<Database> db) {
    auto fn = jsi::Function::createFromHostFunction(rt, name, 1,
        [&db](...) { return db->query(...); }); // db is a local — gone after install() returns
    rt.global().setProperty(rt, "query", std::move(fn));
} // ← db destroyed here; lambda now holds a dangling reference

// CORRECT — [db] copies the shared_ptr, lambda keeps it alive
void install(jsi::Runtime& rt, std::shared_ptr<Database> db) {
    auto fn = jsi::Function::createFromHostFunction(rt, name, 1,
        [db](...) { return db->query(...); }); // lambda owns a share of db
    rt.global().setProperty(rt, "query", std::move(fn));
}
```

---

## RAII — Destructors as Cleanup

RAII (Resource Acquisition Is Initialization) is C++'s answer to `try/finally`. The constructor acquires a resource; the destructor releases it. Because destructors run automatically when an object leaves scope — including during exception unwinding — cleanup is guaranteed without any manual effort.

```cpp
// JavaScript: you must remember to clean up, even under exceptions
function readFile(path) {
    const handle = openFile(path);
    try {
        return handle.read();
    } finally {
        handle.close(); // forget this and the handle leaks
    }
}
```

```cpp
// C++: cleanup is structural, not procedural
class FileHandle {
    FILE* file_;
public:
    explicit FileHandle(const char* path) : file_(fopen(path, "r")) {
        if (!file_) throw std::runtime_error("failed to open");
    }
    ~FileHandle() { fclose(file_); } // runs automatically when FileHandle leaves scope

    std::string read() { /* ... */ }
};

std::string readFile(const char* path) {
    FileHandle handle(path); // constructor opens the file
    return handle.read();
}  // ← ~FileHandle() runs here, file closed — even if read() threw
```

The closing brace is the `finally` block. It runs unconditionally.

**Why RAII matters for JSI:** native modules manage resources the JavaScript GC knows nothing about — audio sessions, database connections, GPU handles, thread pools. RAII ensures these are cleaned up when the HostObject is destroyed, not "eventually" when someone manually calls a cleanup method that might not get called.

```cpp
class AudioSessionHostObject : public jsi::HostObject {
    AVAudioSession* session_;
public:
    AudioSessionHostObject() {
        session_ = [AVAudioSession sharedInstance];
        [session_ setActive:YES error:nil]; // acquire
    }
    ~AudioSessionHostObject() {
        [session_ setActive:NO error:nil];  // release — guaranteed to run
    }
};
```

When JavaScript's GC collects the object and the `shared_ptr` refcount drops to zero, `~AudioSessionHostObject()` runs and the audio session is deactivated. No manual lifecycle management from the JS side required.

---

## Circular Ownership with `shared_ptr` — the Memory Leak Trap

`shared_ptr` reference counting cannot break cycles. If object A holds a `shared_ptr` to B, and B holds a `shared_ptr` to A, both refcounts stay above zero even after all *external* references are dropped. The GC collects the JavaScript proxy objects, releasing their `shared_ptr` handles — but the C++ objects keep each other alive indefinitely.

```cpp
// WRONG — creates a cycle that never frees
class Parent : public jsi::HostObject {
    std::shared_ptr<Child> child_; // Parent holds Child
};

class Child : public jsi::HostObject {
    std::shared_ptr<Parent> parent_; // Child holds Parent — CYCLE
};

// Even after JS drops both objects, refcounts stay at 1. Neither destructor runs.
```

```
Parent ──shared_ptr──▶ Child
  ▲                      │
  └──────shared_ptr───────┘

JS drops both → external refcounts go to 0
But internal cycle keeps both counts at 1 → neither is freed
```

**`weak_ptr` breaks the cycle.** A `weak_ptr` observes an object without owning it — it does not increment the refcount. To use the object, call `.lock()`, which returns a `shared_ptr` if the object is still alive, or an empty `shared_ptr` if it has been freed.

```cpp
// CORRECT — child holds a weak reference back to parent
class Child : public jsi::HostObject {
    std::weak_ptr<Parent> parent_; // doesn't extend Parent's lifetime

public:
    void notifyParent() {
        auto parent = parent_.lock(); // try to get a shared_ptr
        if (!parent) return;          // parent was already destroyed — safe to ignore
        parent->onChildEvent();
    }
};
```

Rule: in any parent-child or observer relationship implemented with `shared_ptr`, the "upward" or "back" reference should be `weak_ptr`. The owning direction (parent → child) uses `shared_ptr`; the back reference (child → parent) uses `weak_ptr`.

---

## The GC Boundary: Who Frees What

JavaScript and C++ run separate heaps with separate rules. Understanding which side owns which memory determines whether your module is safe.

```
┌───────────────────────────────┐   ┌───────────────────────────────┐
│      GC Heap (Hermes)         │   │      Native Heap (C++)        │
│                               │   │                               │
│  Managed by garbage collector │   │  Managed by smart pointers    │
│  Non-deterministic cleanup    │   │  Deterministic cleanup        │
│  Safe: never frees live objs  │   │  Trusts you: free too early → │
│                               │   │  crash; forget → leak         │
│  JS objects, strings, funcs   │   │  HostObjects, buffers, state  │
└──────────────┬────────────────┘   └────────────────┬──────────────┘
               │                                     │
               │           ◆ BOUNDARY ◆              │
               │    shared_ptr bridges both worlds   │
               └─────────────────────────────────────┘
```

**The bridge: `shared_ptr` as the shared checkout card.** When you create a HostObject, the runtime takes a copy of your `shared_ptr`. Now two owners exist — C++ code and the Hermes runtime (on behalf of JavaScript). The C++ object is freed only when both release their reference:

```cpp
void install(jsi::Runtime& rt) {
    auto store = std::make_shared<KVStoreHostObject>(); // refcount = 1

    // Runtime takes its own copy of the shared_ptr — refcount = 2
    auto obj = jsi::Object::createFromHostObject(rt, store);
    rt.global().setProperty(rt, "NativeKV", std::move(obj));

} // `store` local variable destroyed → refcount = 1
  // KVStoreHostObject survives — runtime still holds a reference
  //
  // Later, when JS drops NativeKV and GC runs:
  //   runtime releases its shared_ptr → refcount = 0 → ~KVStoreHostObject() runs
```

**JSI values live on the JS heap — never cache them past the current call.** `jsi::Value`, `jsi::Object`, `jsi::Function` are handles into GC-managed memory. The GC can move or free the underlying object without telling you. Storing a `jsi::Value*` or raw pointer obtained from `.getArrayBuffer(rt).data(rt)` beyond the synchronous HostFunction call is a use-after-free waiting to happen.

To hold a JS callback beyond the current call, wrap it in a `shared_ptr`:

```cpp
// Capture a JS callback for async use
auto callback = std::make_shared<jsi::Value>(std::move(callbackArg));

asyncWork([callback, invoker]() {
    // Back on the JS thread — only then is it safe to call JSI
    invoker->invokeAsync([callback](jsi::Runtime& rt) {
        callback->asObject(rt).asFunction(rt).call(rt);
    });
});
```

**`[this]` in a HostObject method is unsafe.** JavaScript can extract a method, drop the object, wait for GC, and then call the extracted function. At that point `this` points to freed memory:

```cpp
// WRONG
jsi::Value get(jsi::Runtime& rt, const jsi::PropNameID& name) override {
    if (name.utf8(rt) == "increment") {
        return jsi::Function::createFromHostFunction(rt, name, 0,
            [this](...) { count_++; return jsi::Value(count_); }); // raw this — can dangle
    }
    return jsi::Value::undefined();
}
```

Use `weak_from_this()` (requires inheriting from `std::enable_shared_from_this`):

```cpp
// CORRECT
class CounterHostObject
    : public jsi::HostObject,
      public std::enable_shared_from_this<CounterHostObject> {
public:
    jsi::Value get(jsi::Runtime& rt, const jsi::PropNameID& name) override {
        if (name.utf8(rt) == "increment") {
            auto weak = weak_from_this(); // doesn't extend lifetime
            return jsi::Function::createFromHostFunction(rt, name, 0,
                [weak](...) -> jsi::Value {
                    auto self = weak.lock(); // upgrade: returns empty ptr if already freed
                    if (!self) throw jsi::JSError(rt, "Counter was destroyed");
                    self->count_++;
                    return jsi::Value(self->count_);
                });
        }
        return jsi::Value::undefined();
    }
private:
    int count_ = 0;
};
```

**The decision rule for captures in methods returned to JS:**

| Pattern | Use when | Lifetime effect |
|---------|----------|----------------|
| Raw `this` | Called synchronously, never stored by JS | None — safe only if guaranteed |
| `shared_from_this()` | C++ must keep the object alive (timer, background callback) | Extends lifetime — prevents GC cleanup |
| `weak_from_this()` | Function may be stored by JS after dropping the parent object | None — graceful failure if already freed |

Default to `weak_from_this()` for any function returned from `get`. You cannot control what JavaScript does with it.
