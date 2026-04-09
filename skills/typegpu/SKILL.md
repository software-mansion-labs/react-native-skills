---
name: typegpu
description: >-
  TypeGPU is type-safe WebGPU in TypeScript. Use whenever the user writes, debugs, or designs TypeGPU code: 'use gpu' shader functions, tgpu.fn, buffers, textures, bind groups, compute and render pipelines, vertex layouts, slots, accessors, and any TypeGPU API. Shader logic and CPU-side resources are tightly coupled - handle both sides here even if the user only mentions one (e.g. "how do I write a shader", "how do I create a buffer"). Trigger on any mention of typegpu, tgpu, "use gpu", TypedGPU, or WebGPU code using TypeScript schemas.
---

# TypeGPU

A single schema (`d.*`) defines a GPU type, CPU buffer layout, and TypeScript type at once - no manual alignment, type mapping, or casting. The build plugin `unplugin-typegpu` transpiles `'use gpu'`-marked TypeScript into WGSL at compile time, enabling type inference and polymorphism across the CPU/GPU boundary.

---

## Setup

```ts
import tgpu, { d, common } from 'typegpu';
import { std } from 'typegpu';

const root = await tgpu.init();             // request a GPU device
const root = tgpu.initFromDevice(device);   // or wrap an existing GPUDevice

const context = root.configureContext({ canvas, alphaMode: 'premultiplied' });
```

Create one root at app startup. Resources from different roots cannot interact. 

---

## Data schemas (`d.*`)

A schema defines memory layout and infers TypeScript types; the same schema is used for buffers, shader signatures, and bind group entries.

### Scalars
```ts
d.f32    d.i32    d.u32    d.f16
// d.bool is NOT host-shareable - use d.u32 in buffers
```

### Vectors and matrices
```ts
d.vec2f  d.vec3f  d.vec4f     // f32
d.vec2i  d.vec3i  d.vec4i     // i32
d.vec2u  d.vec3u  d.vec4u     // u32
d.vec2h  d.vec3h  d.vec4h     // f16

d.mat2x2f   d.mat3x3f   d.mat4x4f
```

Instance types: `d.vec3f()` -> `d.v3f`, `d.mat4x4f()` -> `d.m4x4f`.

**Vector constructors are richly overloaded - use them.** They compose from any mix of scalars and smaller vectors that adds up to the right component count:

```ts
d.vec3f()              // zero-init: (0, 0, 0)
d.vec3f(1)             // broadcast:  (1, 1, 1)
d.vec3f(1, 2, 3)       // individual components
d.vec3f(someVec2, 1)   // vec2 + scalar
d.vec3f(1, someVec2)   // scalar + vec2

d.vec4f()              // zero-init: (0, 0, 0, 0)
d.vec4f(0.5)           // broadcast:  (0.5, 0.5, 0.5, 0.5)
d.vec4f(rgb, 1)        // vec3 + scalar (common: color + alpha)
d.vec4f(v2a, v2b)      // two vec2s
d.vec4f(1, uv, 0)      // scalar + vec2 + scalar
```

Swizzles (`.xy`, `.zw`, `.rgb`, `.ba`, etc.) return vector instances that work as constructor arguments: `d.vec4f(pos.xy, vel.zw)`.

**Prefer these overloads over manual component decomposition.** Instead of `d.vec3f(v.x, v.y, newZ)`, write `d.vec3f(v.xy, newZ)`.

### Compound types
```ts
const Particle = d.struct({
  position: d.vec2f,
  velocity: d.vec2f,
  color:    d.vec4f,
});

const ParticleArray = d.arrayOf(Particle, 1000); // fixed-size
```

**Runtime-sized schemas.** `d.arrayOf(Element)` without a count returns a *function* `(n: number) => WgslArray<Element>`. This dual nature is the key: pass the function itself (unsized) to bind group layouts, call it with a count (sized) for buffer creation.

```ts
// Plain array - arrayOf without count is already a factory:
const layout = tgpu.bindGroupLayout({
  data: { storage: d.arrayOf(d.f32), access: 'mutable' },  // unsized for layout
});
const buf = root.createBuffer(d.arrayOf(d.f32, 1024)).$usage('storage'); // sized for buffer

// Struct with a runtime-sized last field - wrap in a factory function:
const RuntimeStruct = (n: number) =>
  d.struct({
    counter: d.atomic(d.u32),
    items:   d.arrayOf(d.f32, n),  // last field gets the runtime size
  });

const layout2 = tgpu.bindGroupLayout({
  runtimeData: { storage: RuntimeStruct, access: 'mutable' }, // unsized (the function)
});
const buf2 = root.createBuffer(RuntimeStruct(1024)).$usage('storage'); // sized (called)
```

You cannot pass an unsized schema directly to `createBuffer` - size must be known on the CPU.

---

## GPU functions

TypeGPU compiles TypeScript marked with `'use gpu'` into WGSL.

### Plain callback (polymorphic)

No explicit signature; best for helper math and flexible utilities.

```ts
const rotate = (v: d.v2f, angle: number) => {
  'use gpu';
  const c = std.cos(angle);
  const s = std.sin(angle);
  return d.vec2f(c * v.x - s * v.y, s * v.x + c * v.y);
};
```

`number` parameters and unions like `d.v2f | d.v3f` are polymorphic - TypeGPU generates one WGSL overload per unique call-site type combination. Values captured from outer scope are **inlined as WGSL constants**; use buffers/uniforms for anything that changes at runtime.

### `tgpu.fn` (explicit types)

Pinned WGSL signature. Use for library code or when you need a fixed WGSL interface.

```ts
const rotate = tgpu.fn([d.vec2f, d.f32], d.vec2f)((v, angle) => {
  'use gpu';
  // ...
});
```

### Shader entrypoints

```ts
// Compute
const myCompute = tgpu.computeFn({
  workgroupSize: [64],
  in: { gid: d.builtin.globalInvocationId },
})((input) => { 'use gpu'; /* input.gid: d.v3u */ });

// Vertex
const myVertex = tgpu.vertexFn({
  in:  { position: d.vec3f, uv: d.vec2f },
  out: { position: d.builtin.position, fragUv: d.vec2f },
})((input) => {
  'use gpu';
  return { position: d.vec4f(input.position, 1), fragUv: input.uv };
});

// Fragment
const myFragment = tgpu.fragmentFn({
  in: { fragUv: d.vec2f },
  out: d.vec4f,
})((input) => { 'use gpu'; return d.vec4f(input.fragUv, 0, 1); });
```

Vertex `in` may include builtins: `d.builtin.vertexIndex`, `d.builtin.instanceIndex`.

Full shader syntax, branch pruning, the `std` library, and type inference: see `references/shaders.md`.

---

## Values vs references in `'use gpu'` code

**The single most common source of `ResolutionError` in hand-written shaders.**

**Mental model.** Scalars (`d.f32`, `number`, `boolean`) behave as values. Composites (vectors, matrices, structs, arrays) behave as **references** - a variable holding a `d.vec3f` is a handle to a memory location, not a bag of numbers. Reading is free; writing a component (`v.x = 1`) mutates the location.

**The extra rule.** A reference can be aliased with `const`, or copied with its schema constructor, but it cannot bind to `let` or be assigned with `=`. TypeGPU forces the ambiguity ("rebind or mutate through?") to be resolved at the binding site:

```ts
const startPoint = d.vec2f(1, 2);

let endPoint = startPoint;           // BAD: "references cannot be assigned"
let endPoint = d.vec2f(startPoint);  // OK: copy - fresh, independent, reassignable
const endPoint = startPoint;         // OK: alias - same memory, not reassignable
```

Same rule for reassignment and returning a parameter directly:

```ts
outColor = d.vec3f(spriteRgb);   // OK - not: outColor = spriteRgb
return d.vec3f(baseColor);       // OK when baseColor is a parameter
```

**Expression results are ephemeral and safe.** Constructor calls, arithmetic, function returns, and literals aren't named storage, so `let`/`const`/assignment/arguments/returns all work:

```ts
let p = d.vec2f(0);                // constructor result
p = input.uv * 2 + d.vec2f(1, 0);  // expression result
let q = computeOffset();           // function return
```

The "must copy" rule only kicks in when the right-hand side is itself a named reference (local variable, struct field, array element, function parameter). **Fix:** wrap in the schema constructor (`d.vec2f(v)`, `d.mat4x4f(m)`, `MyStruct(other)`), or swap `let` for `const` if you don't need to reassign.

---

## Idiomatic shader code

**Prefer vector operations over component decomposition.** With `tsover`, `+ - * / %` work on vectors and matrices directly, and scalar broadcast applies to all four arithmetic operators. Write expressive vector math instead of pulling apart components:

```ts
// GOOD - clean vector math:
let uv = input.uv * d.vec2f(0.3, 0.2) + 0.5;
let offset = direction * speed;
let color = baseColor * intensity + ambient;

// BAD - unnecessary decomposition:
let uvX = input.uv.x * 0.3 + 0.5;
let uvY = input.uv.y * 0.2 + 0.5;
let uv = d.vec2f(uvX, uvY);
```

**Struct constructors work inside shaders** - use them to minimize global memory traffic. Instead of mutating struct fields one by one on a storage buffer, build the whole struct locally and assign once:

```ts
const Particle = d.struct({ pos: d.vec2f, vel: d.vec2f, life: d.f32 });

// GOOD - single write to global memory:
const newP = Particle({
  pos: Particle(oldP).pos + Particle(oldP).vel * dt,
  vel: d.vec2f(Particle(oldP).vel) * 0.99,
  life: Particle(oldP).life - dt,
});
particles.$[idx] = Particle(newP);

// BAD - multiple global memory round-trips:
particles.$[idx].pos = particles.$[idx].pos + particles.$[idx].vel * dt;
particles.$[idx].vel.x = particles.$[idx].vel.x * 0.99;
particles.$[idx].vel.y = particles.$[idx].vel.y * 0.99;
particles.$[idx].life = particles.$[idx].life - dt;
```

Struct constructor forms: `MyStruct()` (zero-init), `MyStruct({ field: value, ... })` (named fields), `MyStruct(otherInstance)` (copy). These work both on CPU and inside `'use gpu'` callbacks.

---

## Buffers

### Creating

```ts
// Schema only:
const buf = root.createBuffer(d.arrayOf(Particle, 1000)).$usage('storage');

// With typed initial value (only when non-zero — all buffers are zero-initialized by default):
const uBuf = root.createBuffer(Config, { time: 1, scale: 2.0 }).$usage('uniform');

// With an initializer callback - buffer is still mapped (cheapest CPU path):
const buf = root.createBuffer(Schema, (mappedBuffer) => {
  mappedBuffer.write([10, 20], { startOffset: firstChunk.offset });
  mappedBuffer.write([30, 40], { startOffset: secondChunk.offset });
});

// Wrap an existing GPUBuffer (you own its lifecycle and flags):
const buf = root.createBuffer(d.u32, existingGPUBuffer);
buf.write(12);
```

### Usage flags

| Literal | Shader access |
|---|---|
| `'uniform'` | `var<uniform>` |
| `'storage'` | `var<storage, read>` (or `read_write` with `access: 'mutable'`) |
| `'vertex'` | vertex input, paired with `tgpu.vertexLayout` |
| `'index'` | index buffer (`d.u16` or `d.u32` schema only) |
| `'indirect'` | indirect dispatch/draw |

All buffers get `COPY_SRC | COPY_DST` automatically. `$addFlags(GPUBufferUsage.X)` adds any flag not covered by `$usage`.

### Writing

`.write(value)` handles alignment. Four input forms:

| Form | Example (`vec3f`) | Notes |
|---|---|---|
| Typed instance | `d.vec3f(1, 2, 3)` | Allocates a wrapper |
| Plain JS array / tuple | `[1, 2, 3]` | No allocation - prefer in hot paths |
| TypedArray | `new Float32Array([1, 2, 3])` | Bytes copied as-is - must include padding |
| ArrayBuffer | `rawBytes` | Bytes copied as-is |

For `arrayOf(vec3f, N)`, each element is 16 bytes (12 data + 4 padding). Plain arrays handle padding automatically; `TypedArray`/`ArrayBuffer` inputs must include it manually. `mat3x3f` accepts 9 floats as a plain array; 12 with `Float32Array` (4 per column, 4th is padding). WGSL matrices are **column-major**.

**Fast-path ordering (slowest to fastest):** TypeGPU instances, plain tuples, pre-allocated `Float32Array`, raw `ArrayBuffer`. The gap widens with buffer size. Instances are fine for setup-time data, small rarely-written uniforms, and prototypes. **If it runs in `requestAnimationFrame`, it should not allocate TypeGPU wrappers** - cache instances/typed arrays/buffers at setup and reuse every frame. For struct-heavy uniforms see `references/matrices.md`.

**Slice write** - sub-region with byte offsets, using `d.memoryLayoutOf` to avoid hand-calculating:

```ts
const schema = d.arrayOf(d.u32, 6);
const buffer = root.createBuffer(schema, [0, 1, 2, 0, 0, 0]);

const layout = d.memoryLayoutOf(schema, (a) => a[3]);
buffer.write([4, 5, 6], { startOffset: layout.offset }); // leaves [0,1,2] untouched

// Bounded range:
buffer.write(data, { startOffset: startLayout.offset, endOffset: endLayout.offset });
```

**`.patch(data)`** - update specific struct fields or array indices without touching the rest. For structs, provide a subset of fields. For arrays, use an object with numeric index keys (sparse update), a plain array (full replacement), or a `TypedArray` (byte-level replacement). Values accept the same permissive forms as `.write()` — typed instances, plain tuples, `TypedArray`, `ArrayBuffer`:

```ts
planetBuffer.patch({
  mass: 123.1,
  colors: {
    2: [1, 0, 0],          // plain tuple for vec3f
    4: d.vec3f(0, 0, 1),   // typed instance works too
  },
});
```

**Struct-of-arrays write** - when CPU data is separate per-field arrays (common in simulations):

```ts
common.writeSoA(particleBuffer, {
  position: new Float32Array([1, 2, 3,  4, 5, 6]),  // 3 floats per element (packed)
  velocity: new Float32Array([0.1, 0.2, 0.3]),
});
// TypeGPU scatters packed input into the AoS GPU layout, adding padding as needed.
// Optional slice: common.writeSoA(buffer, fields, { startOffset, endOffset })
```

**GPU-side copy:** `destBuffer.copyFrom(srcBuffer)` (schemas must match).

### Reading

```ts
const data = await buffer.read(); // returns a typed JS value matching the schema
```

### Shorthand "fixed" resources

Skip manual bind groups - the buffer is always bound when referenced in any shader:

```ts
const particles = root.createMutable(d.arrayOf(Particle, 1000)); // var<storage, read_write>
const config    = root.createUniform(Config);                     // var<uniform>
const ro        = root.createReadonly(d.arrayOf(d.f32, N));       // var<storage, read>
```

Access inside shaders via `particles.$`, `config.$`. Prefer fixed resources by default; switch to manual bind groups when you need to swap resources per frame, manage `@group` indices, or share layouts across pipelines.

---

## Bind group layouts (manual binding)

```ts
const layout = tgpu.bindGroupLayout({
  config:    { uniform: ConfigSchema },
  particles: { storage: d.arrayOf(Particle), access: 'mutable' },
  mySampler: { sampler: 'filtering' },   // 'filtering' | 'non-filtering' | 'comparison'
  myTexture: { texture: d.texture2d(d.f32) },
});

// Inside shaders: layout.$.config, layout.$.particles, ...

const bindGroup = root.createBindGroup(layout, {
  config:    configBuffer,
  particles: particleBuffer,
  mySampler: tgpuSampler,
  myTexture: textureOrView,
});

pipeline.with(bindGroup).dispatchWorkgroups(N);
```

Explicit `@group` index (only needed when integrating with raw WGSL that hardcodes group indices): `layout.$idx(0)`.

---

## Compute pipelines

```ts
// Standard - you control workgroup sizing
const pipeline = root.createComputePipeline({ compute: myComputeFn });
pipeline.with(bindGroup).dispatchWorkgroups(Math.ceil(N / 64));

// Guarded - TypeGPU handles workgroup sizing and bounds checking automatically.
// The callback's parameter count sets the dimensionality (0D to 3D):
const p0 = root.createGuardedComputePipeline(() => { 'use gpu'; /* runs once */ });
const p1 = root.createGuardedComputePipeline((x: number) => { 'use gpu'; });
const p2 = root.createGuardedComputePipeline((x: number, y: number) => { 'use gpu'; });
const p3 = root.createGuardedComputePipeline((x: number, y: number, z: number) => { 'use gpu'; });

// dispatchThreads matches the callback's arity - pass thread counts, not workgroup counts.
// TypeGPU picks workgroup sizes internally and injects a bounds guard so threads
// outside the requested range are no-ops.
p2.with(bindGroup).dispatchThreads(width, height);

// WGSL builtins like globalInvocationId are NOT available - use the callback parameters instead.
```

---

## WebGPU coordinate conventions

WebGPU matches DirectX/Metal, **not** OpenGL/WebGL. Most shader tutorials online assume OpenGL; porting verbatim introduces subtle bugs (flipped images, clipped geometry, broken depth tests).

### NDC (clip-space, after `xyz/w`)
- `x in [-1, +1]`, `+x` right
- `y in [-1, +1]`, `+y` **up** (same as OpenGL)
- `z in [0, +1]` - **half range, not `[-1, 1]`** (the OpenGL trap)

A `gluPerspective`-style matrix copied from a WebGL tutorial maps the near plane to `z = -1` and gets clipped. Use a WebGPU/DirectX-style projection (or `typegpu/std` / `wgpu-matrix` helpers) targeting `z in [0, 1]`. Reversed-Z (clear `0`, compare `'greater'`, near plane at `z = 1`) gives better depth precision.

### Framebuffer / fragment coordinates
- `(0, 0)` is **top-left**, `+y` is **down** (opposite of OpenGL)
- These are what `d.builtin.position.xy` gives in a fragment shader (pixel-space, not NDC)

### Texture UVs
- `u, v in [0, 1]`, `(0, 0)` is the first texel = **top-left** of an image from `createImageBitmap`
- Do **not** pre-flip `v` (`1 - v`) - that flips the image upside down

### Matrices are column-major
`d.mat4x4f(c0, c1, c2, c3)` takes four **columns**. `M * v` applies `M`'s transform; composition is right-to-left: `projection * view * model * position`. Inside shaders, use `mat.columns[c][r]` - plain `mat[i]` is rejected. `wgpu-matrix` follows the same convention, so its output drops straight into TypeGPU buffers. For the `Float32Array` padding rule, see `references/matrices.md`.

### Porting checklist (WebGL/three.js shader)
1. Swap `[-1, 1]` z-range projection for a `[0, 1]` one.
2. Don't flip texture v - already top-left.
3. Geometry missing near the camera? Near plane is being clipped.
4. Fullscreen image mirrored? UVs or y flipped twice. `common.fullscreenTriangle.uvs` already produces correct top-left-origin UVs.

---

## Render pipelines

```ts
const pipeline = root.createRenderPipeline({
  vertex:   myVertex,
  fragment: myFragment,
  targets:  { format: presentationFormat }, // single target - shorthand
  primitive?:    GPUPrimitiveState,
  depthStencil?: GPUDepthStencilState,
  multisample?:  GPUMultisampleState,
});

pipeline
  .with(bindGroup)
  .withColorAttachment({
    view: context,
    // loadOp/storeOp/clearValue have defaults
  })
  .withDepthStencilAttachment({ /* ... */ })
  .withIndexBuffer(indexBuffer)  // enables .drawIndexed()
  .draw(vertexCount, instanceCount?);
```

Shell-less inline vertex/fragment lambdas are also valid for simple cases.

### Multiple render targets (MRT)

When a fragment outputs to several attachments (deferred shading G-buffers, picking + colour, bloom thresholding), use a **named record** for the fragment `out`, pipeline `targets`, and `withColorAttachment`. Names - not numeric `@location` indices - are how TypeGPU wires them together, and TypeScript enforces that the keys match.

```ts
const gBufferFrag = tgpu.fragmentFn({
  in:  { worldPos: d.vec3f, normal: d.vec3f },
  out: { albedo: d.vec4f, normal: d.vec4f, position: d.vec4f },
})((input) => ({
  albedo:   d.vec4f(0.8, 0.2, 0.2, 1),
  normal:   d.vec4f(input.normal, 0),
  position: d.vec4f(input.worldPos, 1),
}));

const pipeline = root.createRenderPipeline({
  vertex:   myVertex,
  fragment: gBufferFrag,
  targets: {
    albedo:   { format: 'rgba8unorm' },
    normal:   { format: 'rgba16float' },
    position: { format: 'rgba16float' },
  },
});

pipeline
  .with(bindGroup)
  .withColorAttachment({
    albedo:   { view: albedoView },
    normal:   { view: normalView },
    position: { view: positionView },
  })
  .draw(vertexCount);
```

**Shelled entrypoint field names must be plain WGSL identifiers - no `$`-prefixes.** Keys in `in`/`out` records become WGSL struct fields verbatim. A builtin output is just `depth: d.builtin.fragDepth`, not `$fragDepth: ...`. `$`-prefixed keys belong to the shelless/auto-IO path (`TgpuVertexFn.AutoIn`/`AutoOut`) - different feature, do not mix. (`layout.$`, `config.$` is host-side unwrapping - unrelated.)

Per-target blend/writeMask config and the `fragDepth`-as-output footgun: see `references/pipelines.md`.

### Cache bind groups and views

`root.createBindGroup(...)` and `texture.createView(...)` allocate fresh GPU objects each call. Fine for prototypes; for anything you care about, create them once at setup (near the resource they wrap), store handles in `const`s, and reuse. Per-frame allocation isn't slow per se, but it raises GC pressure and introduces stutters. When a view or bind group legitimately varies each frame, cache the small set you cycle through.

For vertex buffer layouts, the attribs spread trick, and the `common.fullscreenTriangle` helper: `references/pipelines.md`.

---

## GPU-scoped variables

Declared at module scope, persistent on the GPU for the shader's lifetime.

```ts
// Shared across all threads in a workgroup (compute only):
const sharedAccum = tgpu.workgroupVar(d.arrayOf(d.f32, 64));

// Thread-private - each thread gets its own copy:
const threadState = tgpu.privateVar(d.vec3f);

// Compile-time constant - embedded as a WGSL literal:
const PI_OVER_2 = tgpu.const(d.f32, Math.PI / 2);
```

Access via `sharedAccum.$`, `threadState.$`, `PI_OVER_2.$`.

---

## Slots

`tgpu.slot<T>()` is a typed placeholder; fill with `.with(slot, value)` at pipeline, root, or function scope. Any type fits: GPU values, functions, callbacks. Slots are the idiomatic way to build configurable/reusable shaders.

```ts
const distFnSlot = tgpu.slot<(pos: d.v3f) => number>();

const rayMarcher = tgpu.computeFn({
  workgroupSize: [64],
  in: { gid: d.builtin.globalInvocationId },
})(({ gid }) => {
  'use gpu';
  const dist = distFnSlot.$(currentPos); // call the injected function
});

root
  .with(distFnSlot, (pos) => {
    'use gpu';
    return std.length(pos - d.vec3f(0, 0, -5)) - 1.0; // sphere SDF
  })
  .createComputePipeline({ compute: rayMarcher });
```

Scalar/vector slot with a default:

```ts
const colorSlot = tgpu.slot<d.v4f>(d.vec4f(1, 0, 0, 1));
pipeline.with(colorSlot, d.vec4f(0, 1, 0, 1)).draw(3);
```

---

## Accessors

`tgpu.accessor(schema, initial?)` is schema-aware - the value can be a buffer binding, a constant, a literal, or a `'use gpu'` function returning one. The shader is agnostic about how the value is sourced.

```ts
const colorAccess = tgpu.accessor(d.vec3f);

// Fill with a uniform buffer:
root.with(colorAccess, colorBuffer.as('uniform')).createComputePipeline(...)

// Fill with a literal (inlined):
root.with(colorAccess, d.vec3f(1, 0, 0)).createComputePipeline(...)

// Fill with a GPU function:
root.with(colorAccess, () => { 'use gpu'; return computeColor(); }).createComputePipeline(...)
```

Write access: `tgpu.mutableAccessor(schema, initial?)`.

---

## Type utilities

```ts
type ParticleInput = d.InferInput<typeof Particle>; // for .write() calls (CPU-side)
type ParticleGPU   = d.InferGPU<typeof Particle>;   // for shader function parameters

import type { AnyData } from 'typegpu'; // broadest schema constraint

function fillBuffer<T extends AnyData>(buf: TgpuBuffer<T>, value: d.InferInput<T>) {
  buf.write(value);
}
```

---

## Common pitfalls

1. **Numeric literals**: `1.0` may strip to `1` before transpilation -> inferred as `i32`. Use `d.f32(1)` when precision matters.
2. **Outer-scope captures are constants**: inlined once at first compilation. Use `createUniform`/`createMutable` for runtime-dynamic values.
3. **TypedArray/ArrayBuffer alignment**: bytes copied verbatim. `vec3f` elements are 16 bytes (12 + 4 padding). Plain arrays handle padding; typed arrays must include it.
4. **Integer division**: `a / b` on primitives is `f32`. Wrap in `d.i32()`/`d.u32()` for integer semantics.
5. **Uninitialised variables**: `let x;` is invalid - always initialise so the type can be inferred: `let x = d.f32(0)`.
6. **Ternary operators**: runtime ternaries aren't supported. Use `std.select(falseVal, trueVal, condition)`.
7. **Fragment output is always `d.vec4f`**, even for fewer-channel formats. A pipeline with `targets: { format: 'r8unorm' }` or `'rg16float'` still requires `out: d.vec4f` and `return d.vec4f(...)`. WebGPU drops the unused channels.

---

## Companion packages

- **`@typegpu/noise`** - real PRNG (`randf`), distributions (uniform, normal, hemisphere, ...), and Perlin noise (`perlin2d`/`perlin3d`) with optional precomputed gradient caches (~10x speedups). **Prefer it over hand-rolled `fract(sin(...))` hashes** - those are fragile and biased. See `references/noise.md`.

- **`@typegpu/sdf`** - 2D/3D signed distance primitives (`sdDisk`, `sdBox2d`, `sdRoundedBox2d`, `sdBezier`, `sdSphere`, `sdBox3d`, `sdCapsule`, `sdPlane`, ...) and operators (`opUnion`, `opSmoothUnion`, `opSmoothDifference`, `opExtrudeX/Y/Z`). All `tgpu.fn` with pinned types, callable directly from `'use gpu'`. For ray marching, UI masking, AA vector drawing. See `references/sdf.md`.

- **[`wgpu-matrix`](https://github.com/greggman/wgpu-matrix)** - TypeGPU's vectors/matrices are indexable exactly the way `wgpu-matrix` expects, so any `d.mat4x4f`/`d.vec3f` can be passed as the `dst` argument to avoid per-call `Float32Array` allocations. Canonical pattern: a small `d.struct` camera uniform (`{ view, proj, viewInv, projInv }`) updated via `.patch()`. In hot paths, prefer pre-allocated `Float32Array`s with subarrays fed to `wgpu-matrix` as `dst`, raw `ArrayBuffer`s, or `common.writeSoA`. See `references/matrices.md`.

---

## Reference files

| File | Contents |
|---|---|
| `references/setup.md` | Install, `unplugin-typegpu`, `tsover` operator overloading |
| `references/shaders.md` | Full shader guide: syntax limits, branch pruning, `std` library, builtins, `console.log` |
| `references/types.md` | **Abstract types, numeric literal rules, when `d.f32()` is redundant vs required, sampler and texture schemas for `tgpu.fn` signatures** — read this whenever writing shader code |
| `references/textures.md` | Creation, views, mipmaps, samplers, storage textures |
| `references/pipelines.md` | Vertex buffers and layouts, attribs wiring, MRT, fullscreen triangle, resolve API |
| `references/noise.md` | `@typegpu/noise` - `randf`, distributions, `perlin2d`/`perlin3d` with static/dynamic caches |
| `references/sdf.md` | `@typegpu/sdf` - 2D/3D primitives, operators, ray marching, AA masks, SDF baking |
| `references/matrices.md` | `wgpu-matrix` integration, column-major layout, camera uniforms, fast-path CPU writes |
| `references/advanced.md` | Buffer reinterpretation, indirect drawing, custom command encoders |
