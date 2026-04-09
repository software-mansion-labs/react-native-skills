# TypeGPU Advanced

## Buffer reinterpretation

Pass an existing `GPUBuffer` as `initialData` to create a TypeGPU buffer aliasing the same GPU memory with a different type. Useful for SSBO-as-vertex-buffer patterns.

```ts
const packedBuffer = root
  .createBuffer(d.arrayOf(d.unorm8x4))
  .$usage('vertex');

// Add STORAGE to the underlying GPUBuffer:
packedBuffer.$addFlags(GPUBufferUsage.STORAGE);

// Alias same memory, typed as u32 storage:
const storageView = root.createBuffer(d.arrayOf(d.u32), packedBuffer.buffer);
```

Pairs well with WGSL pack/unpack builtins (`std.pack4x8unorm`, `std.unpack4x8unorm`, `std.pack2x16float`, etc.) - reinterpret a buffer as `u32` storage and pack/unpack in the shader for compact vertex data, color encoding, or quantized weights.

**Caveats:**
- The original buffer's lifecycle is NOT transferred - keep it alive while the alias is in use.
- `$usage()` and `$addFlags()` cannot be called on the aliased buffer.
- The original must have all needed usage flags before the alias is created.

---

## Indirect drawing and dispatching

Indirect buffers let the GPU determine draw/dispatch counts - foundation of GPU-driven rendering.

### Required buffer contents

| Call | Layout |
|---|---|
| `dispatchWorkgroupsIndirect` | 3x `u32`: x, y, z workgroup counts |
| `drawIndirect` | 4x `u32`: vertexCount, instanceCount, firstVertex, firstInstance |
| `drawIndexedIndirect` | indexCount(`u32`), instanceCount(`u32`), firstIndex(`u32`), baseVertex(`i32`), firstInstance(`u32`) |

All indirect methods have two overloads: `(buffer)` (offset 0) and `(buffer, offsetInfo)`. When the indirect params start at the beginning of the buffer, prefer the no-offset overload - it's cleaner and equally safe.

```ts
// Dedicated indirect buffer - no offset needed:
const IndirectParams = d.struct({
  vertexCount:   d.u32,
  instanceCount: d.u32,
  firstVertex:   d.u32,
  firstInstance: d.u32,
});
const indirectBuf = root.createBuffer(IndirectParams).$usage('storage', 'indirect');
pipeline.drawIndirect(indirectBuf); // offset 0 implied
```

### `d.memoryLayoutOf` - safe offset calculation

When indirect params are embedded in a larger struct, use `d.memoryLayoutOf` instead of hardcoding byte offsets:

```ts
const Schema = d.struct({
  someData:      d.arrayOf(d.vec3f, 10),
  vertexCount:   d.u32,
  instanceCount: d.u32,
  firstVertex:   d.u32,
  firstInstance: d.u32,
});

const MyBuffer = root.createBuffer(Schema).$usage('storage', 'indirect');

const drawOffset = d.memoryLayoutOf(Schema, (s) => s.vertexCount); // compute once
pipeline.drawIndirect(MyBuffer, drawOffset); // reuse every frame
```

### Packing indirect params as a vector

`vec4u` guarantees no padding between the four draw params:

```ts
const Schema = d.struct({
  someData:   d.arrayOf(d.vec3f, 10),
  drawParams: d.vec4u, // [vertexCount, instanceCount, firstVertex, firstInstance]
});

const offset = d.memoryLayoutOf(Schema, (s) => s.drawParams);
pipeline.drawIndirect(MyBuffer, offset);
```

---

## Custom command encoders

TypeGPU supports passing an existing `GPUCommandEncoder` or active `GPURenderPassEncoder`/`GPUComputePassEncoder` via `.with(encoder)` or `.with(pass)`, allowing TypeGPU calls to interleave with raw WebGPU commands in a shared command buffer.

> Full documentation for this pattern is not yet in this skill. Refer to the TypeGPU source or examples.

---

## `$addFlags` - raw usage flags

For flags not covered by `$usage` (e.g. `MAP_READ`, `QUERY_RESOLVE`):

```ts
const mappableBuffer = root.createBuffer(d.vec4f).$addFlags(GPUBufferUsage.MAP_READ);
```

`MAP_READ` and `MAP_WRITE` are mutually exclusive with most other usages - setting either overwrites existing flags and adds `COPY_(DST|SRC)`. Other flags are OR'd. Cannot be used on buffers created from an existing `GPUBuffer`.
