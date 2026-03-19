---
name: on-device-ai
description: "Software Mansion's best practices for on-device AI in React Native using React Native ExecuTorch. Use when running machine learning models locally on mobile devices for LLMs, computer vision, OCR, speech processing, text/image embeddings, or any AI feature without cloud dependencies. Trigger on: 'react-native-executorch', 'ExecuTorch', 'on-device AI', 'on-device ML', 'local AI', 'offline AI', 'useLLM', 'useClassification', 'useObjectDetection', 'useOCR', 'useVerticalOCR', 'useStyleTransfer', 'useTextToImage', 'useImageEmbeddings', 'useImageSegmentation', 'useSemanticSegmentation', 'useSpeechToText', 'useTextToSpeech', 'useVAD', 'useTextEmbeddings', 'useTokenizer', 'useExecutorchModule', 'ResourceFetcher', 'image classification', 'object detection', 'style transfer', 'speech-to-text', 'text-to-speech', 'voice activity detection', 'text embeddings', 'image embeddings', 'mobile LLM', 'on-device chatbot', 'document scanning OCR', or any request to run AI/ML models locally in a React Native app."
---

# On-Device AI

Software Mansion's production patterns for on-device AI in React Native using [React Native ExecuTorch](https://github.com/software-mansion/react-native-executorch).

For hook APIs, model constants, configuration options, and TypeScript module classes, webfetch the relevant page from the official docs. The base URL is `https://docs.swmansion.com/react-native-executorch/docs/`. To find the correct page for a specific hook or module, webfetch the sitemap at `https://docs.swmansion.com/react-native-executorch/sitemap.xml` and look for the latest versioned URL matching your topic.

## References

| File | When to read |
|------|-------------|
| `best-practices.md` | Model loading strategy, device constraints, error handling, audio/image processing tips, performance considerations |

## Documentation Map

When you need API details for a specific feature, webfetch the corresponding page from the latest version of the official docs:

| Feature | What to webfetch |
|---------|-----------------|
| Getting started | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/fundamentals/getting-started` |
| Loading models | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/fundamentals/loading-models` |
| LLMs (chat, tool calling, structured output) | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useLLM` |
| Speech-to-text | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useSpeechToText` |
| Text-to-speech | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useTextToSpeech` (if 404, check next version) |
| Voice activity detection | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useVAD` |
| Text embeddings | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useTextEmbeddings` |
| Tokenizer | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/natural-language-processing/useTokenizer` |
| Image classification | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useClassification` |
| Object detection | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useObjectDetection` |
| Image segmentation | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useImageSegmentation` |
| OCR (horizontal text) | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useOCR` |
| OCR (vertical text) | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useVerticalOCR` |
| Style transfer | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useStyleTransfer` |
| Text-to-image | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useTextToImage` |
| Image embeddings | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/computer-vision/useImageEmbeddings` |
| Custom models | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/hooks/executorch-bindings/useExecutorchModule` |
| Resource fetcher | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/utilities/resource-fetcher` |
| Benchmarks | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/benchmarks/inference-time` |
| FAQ | `https://docs.swmansion.com/react-native-executorch/docs/0.6.x/fundamentals/frequently-asked-questions` |
| Full API reference | `https://docs.swmansion.com/react-native-executorch/docs/api-reference` |

If a URL returns 404, the docs may have been restructured. Webfetch the sitemap at `https://docs.swmansion.com/react-native-executorch/sitemap.xml` to find the current URL.
