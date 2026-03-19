# On-Device AI Best Practices

For the full hook API, model constants, and configuration options, webfetch the relevant page from the official docs. See `SKILL.md` for the documentation map.

---

## Model Loading Strategy

Models can be loaded from three sources. Choose based on model size and UX requirements:

| Source | Best for | Trade-off |
|--------|----------|-----------|
| Bundled with app (assets folder) | Small models (< 512MB) | Available immediately, increases app size |
| Remote URL (downloaded on first use) | Large models (> 512MB) | Keeps app small, requires internet on first use |
| Local file system | User-managed models | Maximum flexibility, requires custom download/file management UI |

**Guidelines:**
- Bundle small models for instant availability
- Download large models on first use with progress tracking via `ResourceFetcher`
- Prefer quantized model variants to save memory and storage
- Show download progress UI for remote models

---

## Device Constraints and Model Selection

**Memory:** Low-end devices handle smaller models (135M-1.7B parameters) and quantized variants. High-end devices can run larger models (3B-4B parameters).

**Processing power:** Lower-end devices have longer inference times. Check the official benchmarks before choosing a model for your target devices.

**Storage:** Large models require significant disk space. Implement cleanup mechanisms via `ResourceFetcher` to remove unused models. Monitor total downloaded model size.

**Guidelines:**
- Always check model memory requirements against target device specs
- Prefer quantized model variants on lower-end devices
- Test on the lowest-spec device you plan to support
- Consider providing a cloud API fallback for devices that cannot run the model

---

## Error Handling

Use `RnExecutorchError` and its error codes for robust error handling. Common failure modes:

| Failure | Likely cause | Recovery |
|---------|-------------|----------|
| Model not loading | Invalid source URL/path, insufficient storage | Verify model source, check available disk space |
| Out of memory | Model too large for device | Switch to a smaller or quantized variant |
| Poor LLM quality | Suboptimal generation config | Adjust temperature/top-p, improve system prompt |
| Download failures | Network issues | Implement retry logic, use `ResourceFetcher` pause/resume |

Always provide user-visible loading states and error messages. Model loading and inference are long-running operations.

---

## Audio Processing

Audio must be in the correct sample rate for processing:

- **Speech-to-text and VAD input:** 16kHz sample rate
- **Text-to-speech output:** 24kHz sample rate

Always decode/resample audio to the correct rate before processing. Mismatched sample rates produce garbled results silently.

---

## Image Processing

Images can be provided as:
- Remote URLs (http/https), automatically cached
- Local file URIs (file://)
- Base64-encoded strings

Image preprocessing (resizing, normalization) is handled automatically by the vision hooks. You do not need to manually resize or normalize images before passing them.

---

## Text Token Limits

Text embeddings and LLMs have maximum token limits. Text exceeding these limits is truncated. Use `useTokenizer` to count tokens before processing when working with variable-length input.

---

## Implementation Checklist

**Planning:**
- Identify the AI task (chat, vision, audio, search)
- Consider device memory constraints and target devices
- Choose an appropriate model from available options
- Decide on model loading strategy (bundled vs. remote)

**Development:**
- Select the correct hook for your task
- Configure model loading source
- Implement error handling with user-visible feedback
- Add loading states for model initialization and inference
- Set up `ResourceFetcher` for large model downloads

**Testing:**
- Test on the target minimum-spec device
- Verify offline functionality after initial model download
- Check memory usage stays within device limits
- Test error handling paths (network failure, insufficient memory, invalid inputs)
- Measure inference time for acceptable UX

**Deployment:**
- Finalize model bundling strategy (size vs. download trade-off)
- Implement download progress UI for remote models
- Plan for model version updates
