---
name: on-device-ai
description: "Software Mansion's best practices for on-device AI in React Native using React Native ExecuTorch. Covers LLMs (chat, tool calling, structured output), computer vision (classification, object detection, OCR, segmentation, style transfer, embeddings, text-to-image), speech processing (speech-to-text, text-to-speech, voice activity detection), VisionCamera real-time frame processing, model loading and resource management. Trigger on: 'react-native-executorch', 'ExecuTorch', 'on-device AI', 'on-device ML', 'local AI', 'offline AI', 'useLLM', 'useClassification', 'useObjectDetection', 'useOCR', 'useVerticalOCR', 'useStyleTransfer', 'useTextToImage', 'useImageEmbeddings', 'useImageSegmentation', 'useSpeechToText', 'useTextToSpeech', 'useVAD', 'useTextEmbeddings', 'useTokenizer', 'useExecutorchModule', 'ResourceFetcher', 'real-time frame processing', 'runOnFrame', 'tool calling', 'structured output', or any request to run AI/ML models locally in a React Native app."
---

# On-Device AI

Software Mansion's production patterns for on-device AI in React Native using [React Native ExecuTorch](https://github.com/software-mansion/react-native-executorch).

Load at most one reference file per question. For hook API signatures, model constants, and configuration options, webfetch the relevant page from the official docs at `https://docs.swmansion.com/react-native-executorch/docs/`.

## Decision Tree

Pick the right hook based on the AI task.

```
What AI task does the feature need?
│
├── Text generation, chatbot, or reasoning?
│   └── useLLM                                    → see llm.md
│       ├── Text-only chat → standard useLLM
│       ├── Tool calling → configure with toolsConfig
│       └── Structured JSON output → getStructuredOutputPrompt
│
├── Understanding images?
│   ├── What's in this image? → useClassification  → see vision.md
│   ├── Where are objects? → useObjectDetection    → see vision.md
│   ├── Read text from image? → useOCR / useVerticalOCR → see vision.md
│   ├── Segment regions? → useImageSegmentation    → see vision.md
│   ├── Apply artistic style? → useStyleTransfer   → see vision.md
│   ├── Generate image from text? → useTextToImage  → see vision.md
│   └── Embed image as vector? → useImageEmbeddings → see vision.md
│
├── Speech or audio processing?
│   ├── Transcribe speech → useSpeechToText        → see speech.md
│   ├── Synthesize speech → useTextToSpeech        → see speech.md
│   └── Detect speech segments → useVAD            → see speech.md
│
├── Text utilities?
│   ├── Convert text to vectors → useTextEmbeddings → see vision.md
│   └── Count tokens → useTokenizer
│
├── Real-time camera processing?
│   └── runOnFrame with VisionCamera v5            → see vision.md
│
└── Custom model (.pte)?
    └── useExecutorchModule                        → see setup.md
```

## Critical Rules

- **Always check `isReady` before calling `forward` or `generate`.** Hooks load models asynchronously. Calling inference methods before the model is ready throws `ModuleNotLoaded`.

- **Interrupt LLM generation before unmounting the component.** Unmounting while `isGenerating` is true causes a crash. Call `llm.interrupt()` and wait for `isGenerating` to become false before navigating away.

- **Use quantized models on mobile.** Full-precision models consume too much memory for most devices. React Native ExecuTorch ships quantized variants for all supported models.

- **Audio for speech-to-text must be 16kHz mono.** Mismatched sample rates produce garbled transcriptions silently.

- **Audio from text-to-speech is 24kHz.** Create the `AudioContext` with `{ sampleRate: 24000 }` for playback.

- **Set `pixelFormat: 'rgb'` for VisionCamera frame processing.** The default `yuv` format produces incorrect results with ExecuTorch vision models.

## References

| File | When to read |
|------|-------------|
| `llm.md` | LLM chat (functional and managed), tool calling, structured output, token batching, context window configuration, model selection, generation config |
| `vision.md` | Image classification, object detection, OCR, image segmentation, style transfer, text-to-image, image/text embeddings, VisionCamera real-time frame processing with `runOnFrame` |
| `speech.md` | Speech-to-text (batch and streaming transcription), text-to-speech (batch and streaming synthesis), voice activity detection, audio format requirements |
| `setup.md` | Installation, model loading strategies (bundled, remote, local), `ResourceFetcher` download management, error handling with `RnExecutorchError`, custom models with `useExecutorchModule`, Metro config for `.pte` files |
