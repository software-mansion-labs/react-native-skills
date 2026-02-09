---
name: react-native-executorch
description: Run AI models and LLMs locally on-device in React Native apps using Meta's ExecuTorch. Provides declarative hooks for computer vision, natural language processing, and custom model inference without cloud dependencies.
---

## Overview

This skill provides expertise in using React Native Executorch - a library that enables on-device AI model execution in React Native applications. The library supports various AI tasks including LLMs, computer vision, audio processing, and natural language processing.

## When to Use This Skill

Use this skill when users ask about:

- Running AI models on mobile devices with React Native
- On-device LLMs, image processing, OCR, speech-to-text, text-to-speech
- React Native Executorch library implementation
- Mobile AI without backend/cloud dependencies
- ExecuTorch integration in React Native apps

## Core Capabilities

### 1. Large Language Models (LLMs)

- On-device text generation and chat
- Tool/function calling
- Structured output generation
- Streaming responses
- Multiple model families: Llama 3.2, Qwen 3, Hammer 2.1, SmolLM2, Phi 4

### 2. Computer Vision

- **Image Classification** - Categorize images (EfficientNet V2)
- **Object Detection** - Detect and locate objects (SSDLite 320)
- **Image Segmentation** - Pixel-level classification (DeepLab V3)
- **Style Transfer** - Apply artistic styles (Candy, Mosaic, Udnie, Rain Princess)
- **Text to Image** - Generate images from text (BK-SDM Tiny)
- **Image Embeddings** - Image similarity and search (CLIP)
- **OCR** - Text detection and recognition (horizontal and vertical)

### 3. Audio Processing

- **Speech to Text** - Transcription (Whisper, Moonshine)
- **Text to Speech** - Natural voice synthesis (Kokoro)
- **Voice Activity Detection** - Detect speech segments (FSMN-VAD)

### 4. Natural Language Processing

- **Text Embeddings** - Semantic similarity and search
- **Tokenization** - Text-to-token conversion

## Best Practices

**Quick Links to Detailed Guides:**

- LLMs: [LLM Usage Reference](./references/reference-llms.md)
- Vision: [Computer Vision Part 1](./references/computer-vision-1.md), [Part 2](./computer-vision-2.md)
- Audio: [Audio Models Reference](./audio-models.md)
- Text: [OCR Reference](./ocr-usage.md), [Text Embeddings](./text-embeddings-tokenizer.md)
- Core: [Core Utilities](./core-utils.md), [Models & Loading](./models-loading.md)

### Model Selection

**Consider Device Constraints:**

- Use quantized models for lower-end devices (e.g., `LLAMA3_2_1B_QLORA` instead of `LLAMA3_2_1B`)
- Check memory requirements before recommending models
- Smaller models (135M-1.7B params) for basic tasks
- Larger models (3B-4B params) only for high-end devices

**Model Loading Strategy:**

- Small models (<512MB): Bundle with app using `require()`
- Large models (>512MB): Download from URL on first use
- Always show download progress for remote models
- Implement offline-first with cached models

### Error Handling

**See [Core Utilities Reference](./core-utils.md) for complete error handling guide**

**Always wrap AI operations in try-catch:**

```typescript
import {
  RnExecutorchError,
  RnExecutorchErrorCode,
} from "react-native-executorch";

try {
  await model.forward(input);
} catch (err) {
  if (err instanceof RnExecutorchError) {
    switch (err.code) {
      case RnExecutorchErrorCode.ModuleNotLoaded:
        // Handle not loaded
        break;
      case RnExecutorchErrorCode.MemoryAllocationFailed:
        // Suggest smaller model
        break;
      default:
        console.error(err.message);
    }
  }
}
```

### LLM Implementation Patterns

**See [LLM Usage Reference](./llm-usage.md) for complete details**

**Functional Mode (Stateless) - Best for:**

- Single completions
- Fine-grained control over conversation
- Custom message management

```typescript
const llm = useLLM({ model: LLAMA3_2_1B });
const response = await llm.generate(messages);
```

**Managed Mode (Stateful) - Best for:**

- Chat interfaces
- Automatic conversation history
- Simplified state management

```typescript
useEffect(() => {
  llm.configure({
    chatConfig: { systemPrompt: "...", contextWindowLength: 10 },
    generationConfig: { temperature: 0.7, topp: 0.9 },
  });
}, []);
llm.sendMessage("Hello!");
```

**Always interrupt before unmount:**

```typescript
useEffect(() => {
  return () => {
    if (llm.isGenerating) {
      llm.interrupt();
    }
  };
}, []);
```

### Audio Processing Requirements

**📄 See [Audio Models Reference](./audio-models.md) for complete details**

**Critical audio format requirements:**

- Speech-to-text: 16kHz sample rate (use AudioContext with `sampleRate: 16000`)
- Text-to-speech: 24kHz sample rate for output
- Always decode audio with correct sample rate before processing

```typescript
const audioContext = new AudioContext({ sampleRate: 16000 });
const decodedAudio = await audioContext.decodeAudioDataSource(audioUri);
const waveform = decodedAudio.getChannelData(0);
```

### Resource Management

**📄 See [Core Utilities Reference](./core-utils.md) for ResourceFetcher details**

**Download management:**

- Check file size before downloading (`ResourceFetcher.getFilesTotalSize()`)
- Show progress during downloads
- Implement pause/resume for large models
- Clean up old models to save space

```typescript
const totalSize = await ResourceFetcher.getFilesTotalSize(modelUrl);
const sizeMB = (totalSize / 1024 / 1024).toFixed(2);
console.log(`Download size: ${sizeMB} MB`);

const uris = await ResourceFetcher.fetch(
  (progress) => console.log(`${(progress * 100).toFixed(1)}%`),
  modelUrl,
);
```

**Memory management:**

- Unload unused models
- Use quantized models when possible
- Monitor device memory constraints
- Suggest appropriate model sizes based on device capabilities

## Common Patterns

**Detailed examples available in:**

- [LLM Usage Reference](./llm-usage.md) - Chat, tools, structured output
- [Computer Vision References](./computer-vision-1.md) - Image tasks
- [OCR Reference](./ocr-usage.md) - Text recognition
- [Audio Models Reference](./audio-models.md) - Speech processing

### LLM Chat Interface

```typescript
import { useLLM, QWEN3_1_7B_QUANTIZED } from 'react-native-executorch';

function ChatApp() {
  const llm = useLLM({ model: QWEN3_1_7B_QUANTIZED });

  useEffect(() => {
    llm.configure({
      chatConfig: {
        systemPrompt: "You are a helpful assistant",
        contextWindowLength: 10
      },
      generationConfig: {
        temperature: 0.7,
        topp: 0.9
      }
    });
  }, []);

  return (
    <View>
      <ScrollView>
        {llm.messageHistory.map((msg, i) => (
          <Text key={i}>{msg.role}: {msg.content}</Text>
        ))}
      </ScrollView>
      <TextInput onSubmitEditing={(e) => llm.sendMessage(e.nativeEvent.text)} />
      {llm.isGenerating && <Button onPress={llm.interrupt} title="Stop" />}
    </View>
  );
}
```

### Tool Calling with LLMs

```typescript
const TOOLS: LLMTool[] = [
  {
    name: "get_weather",
    description: "Get weather in given location",
    parameters: {
      type: "dict",
      properties: {
        location: { type: "string", description: "Location name" },
      },
      required: ["location"],
    },
  },
];

const executeTool = async (call: ToolCall) => {
  switch (call.toolName) {
    case "get_weather":
      const location = call.arguments.location;
      return await fetchWeather(location);
    default:
      return null;
  }
};

llm.configure({
  toolsConfig: {
    tools: TOOLS,
    executeToolCallback: executeTool,
    displayToolCalls: true,
  },
});
```

### Structured Output

```typescript
import {
  getStructuredOutputPrompt,
  fixAndValidateStructuredOutput,
} from "react-native-executorch";
import * as z from "zod/v4";

const schema = z.object({
  username: z.string(),
  question: z.string().optional(),
  bid: z.number(),
  currency: z.string().optional(),
});

useEffect(() => {
  const instructions = getStructuredOutputPrompt(schema);
  llm.configure({
    chatConfig: {
      systemPrompt: `Parse user messages as JSON. ${instructions}`,
    },
  });
}, []);

// After generation completes
const parsed = fixAndValidateStructuredOutput(
  llm.messageHistory.at(-1).content,
  schema,
);
```

### Image Classification

```typescript
import { useClassification, EFFICIENTNET_V2_S } from "react-native-executorch";

const model = useClassification({ model: EFFICIENTNET_V2_S });

const classifyImage = async (imageUri: string) => {
  const results = await model.forward(imageUri);

  const topThree = Object.entries(results)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 3)
    .map(([label, score]) => ({ label, score }));

  return topThree;
};
```

### OCR (Text Recognition)

```typescript
import { useOCR, OCR_ENGLISH } from "react-native-executorch";

const model = useOCR({ model: OCR_ENGLISH });

const extractText = async (imageUrl: string) => {
  const detections = await model.forward(imageUrl);

  for (const detection of detections) {
    console.log("Text:", detection.text);
    console.log("Confidence:", detection.score);
    console.log("Position:", detection.bbox);
  }
};
```

### Speech to Text with Streaming

```typescript
import { useSpeechToText, WHISPER_TINY_EN } from "react-native-executorch";
import { AudioRecorder } from "react-native-audio-api";

const stt = useSpeechToText({ model: WHISPER_TINY_EN });

const recorder = new AudioRecorder({
  sampleRate: 16000,
  bufferLengthInSamples: 1600,
});

recorder.onAudioReady(({ buffer }) => {
  stt.streamInsert(buffer.getChannelData(0));
});

recorder.start();
await stt.stream();

// Access transcription
console.log(stt.committedTranscription);
console.log(stt.nonCommittedTranscription);

// Cleanup
recorder.stop();
stt.streamStop();
```

### Text to Speech

```typescript
import {
  useTextToSpeech,
  KOKORO_MEDIUM,
  KOKORO_VOICE_AF_HEART,
} from "react-native-executorch";
import { AudioContext } from "react-native-audio-api";

const tts = useTextToSpeech({
  model: KOKORO_MEDIUM,
  voice: KOKORO_VOICE_AF_HEART,
});

const audioContext = new AudioContext({ sampleRate: 24000 });

const speak = async (text: string) => {
  const waveform = await tts.forward(text, 1.0);

  const buffer = audioContext.createBuffer(1, waveform.length, 24000);
  buffer.getChannelData(0).set(waveform);

  const source = audioContext.createBufferSource();
  source.buffer = buffer;
  source.connect(audioContext.destination);
  source.start();
};
```

## Model Recommendations by Use Case

**📄 See [Available Models Reference](./models-loading.md) for complete model catalog**

### Chatbots & Assistants

- **Entry-level:** `SMOLLM2_1_360M_QUANTIZED` - Fast, minimal memory
- **Mid-range:** `QWEN3_1_7B_QUANTIZED` - Good balance
- **High-end:** `LLAMA3_2_3B_QLORA` - Best quality

### Function Calling

- **Recommended:** `HAMMER2_1_1_5B_QUANTIZED` - Optimized for tools
- **Alternative:** `QWEN3_4B_QUANTIZED` - Strong reasoning

### Structured Output

- **Recommended:** `QWEN3_4B_QUANTIZED` - Excellent JSON generation
- **Alternative:** `PHI_4_MINI_4B_QUANTIZED` - Good format adherence

### Image Tasks

- **Classification:** `EFFICIENTNET_V2_S`
- **Object Detection:** `SSDLITE_320_MOBILENET_V3_LARGE`
- **Segmentation:** `DEEPLAB_V3_RESNET50`
- **OCR:** `OCR_ENGLISH` or language-specific variants

### Speech Tasks

- **Transcription (English):** `WHISPER_TINY_EN_QUANTIZED` - Fast
- **Transcription (Multilingual):** `WHISPER_BASE` - Balanced
- **Text-to-Speech:** `KOKORO_MEDIUM` - Natural voices

## Important Constraints

### Technical Limitations

- **Audio formats:** Must be 16kHz for STT, 24kHz output for TTS
- **Image sizes:** Vary by model (typically 224x224 or specific to model)
- **Context windows:** Limited by model (check token limits)
- **Token limits:** Text embeddings have max token counts (check model specs)

### Platform Considerations

- Models run on-device - no internet required after download
- Memory constraints vary by device
- Processing time varies by device capability
- Some models iOS/Android optimized (check model details)

### OCR Specifics

**📄 See [OCR Usage Reference](./ocr-usage.md) for complete guide**

- `useOCR`: For horizontal text only
- `useVerticalOCR`: For vertical text (experimental)
- Use `independentCharacters: true` for CJK languages
- Recognizer must match target language alphabet

### LLM Reasoning Mode

- Some models (Qwen 3) have reasoning mode
- Disable with `/no_think` suffix in prompt
- Increases token usage if enabled

## Troubleshooting Guide

### Model Not Loading

- Check model source is valid URL or file path
- Verify sufficient device storage
- Ensure network connectivity for remote models
- Check error code for specific issue

### Out of Memory

- Switch to quantized model variant
- Use smaller model size
- Clear app cache
- Restart app to free memory

### Poor LLM Quality

- Adjust temperature (lower = more focused)
- Modify top-p sampling
- Improve system prompt
- Try larger model if device supports it

### Audio Processing Issues

- Verify correct sample rate (16kHz for STT)
- Check audio format compatibility
- Ensure AudioContext configured correctly
- Validate waveform data structure

### Download Failures

- Implement retry logic with exponential backoff
- Check network connectivity
- Verify URL accessibility
- Use pause/resume for large files

## Code Generation Guidelines

When generating code for React Native Executorch:

1. **Always import from 'react-native-executorch'**
2. **Use TypeScript for type safety**
3. **Include proper error handling with RnExecutorchError**
4. **Show loading states during model operations**
5. **Implement cleanup in useEffect returns**
6. **Validate inputs before passing to models**
7. **Use appropriate model constants, not hardcoded URLs**
8. **Include progress tracking for downloads**
9. **Handle async operations properly**
10. **Add comments for complex preprocessing/postprocessing**

## Example: Complete LLM Chat App

```typescript
import React, { useEffect, useState } from 'react';
import { View, Text, TextInput, Button, ScrollView, ActivityIndicator } from 'react-native';
import { useLLM, QWEN3_1_7B_QUANTIZED, RnExecutorchError } from 'react-native-executorch';

export default function ChatApp() {
  const llm = useLLM({ model: QWEN3_1_7B_QUANTIZED });
  const [input, setInput] = useState('');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    try {
      llm.configure({
        chatConfig: {
          systemPrompt: 'You are a helpful AI assistant.',
          contextWindowLength: 10
        },
        generationConfig: {
          temperature: 0.7,
          topp: 0.9
        }
      });
    } catch (err) {
      if (err instanceof RnExecutorchError) {
        setError(err.message);
      }
    }
  }, []);

  useEffect(() => {
    return () => {
      if (llm.isGenerating) {
        llm.interrupt();
      }
    };
  }, []);

  const handleSend = () => {
    if (!input.trim()) return;

    try {
      llm.sendMessage(input);
      setInput('');
      setError(null);
    } catch (err) {
      if (err instanceof RnExecutorchError) {
        setError(err.message);
      }
    }
  };

  return (
    <View style={{ flex: 1, padding: 20 }}>
      <ScrollView style={{ flex: 1 }}>
        {llm.messageHistory.map((msg, i) => (
          <View key={i} style={{ marginVertical: 5 }}>
            <Text style={{ fontWeight: 'bold' }}>{msg.role}:</Text>
            <Text>{msg.content}</Text>
          </View>
        ))}
        {llm.isGenerating && <ActivityIndicator />}
      </ScrollView>

      {error && <Text style={{ color: 'red' }}>{error}</Text>}

      <View style={{ flexDirection: 'row' }}>
        <TextInput
          value={input}
          onChangeText={setInput}
          placeholder="Type a message..."
          style={{ flex: 1, borderWidth: 1, padding: 10 }}
        />
        <Button onPress={handleSend} title="Send" disabled={llm.isGenerating} />
        {llm.isGenerating && (
          <Button onPress={() => llm.interrupt()} title="Stop" />
        )}
      </View>
    </View>
  );
}
```

## Internal Reference Files

When working with React Native Executorch, consult these reference files for detailed implementation guidance:

- **[OCR Usage](./ocr-usage.md)** - Complete guide for useOCR and useVerticalOCR with language support
- **[Text Embeddings & Tokenizer](./text-embeddings-tokenizer.md)** - useTextEmbeddings and useTokenizer reference
- **[Available Models & Loading](./models-loading.md)** - Complete model catalog and loading strategies
- **[LLM Usage](./llm-usage.md)** - Comprehensive LLM guide with functional/managed modes, tool calling, structured output
- **[Computer Vision (Part 1)](./computer-vision-1.md)** - Classification, Segmentation, Object Detection
- **[Computer Vision (Part 2)](./computer-vision-2.md)** - Style Transfer, Text-to-Image, Image Embeddings
- **[Audio Models](./audio-models.md)** - Speech-to-Text, Text-to-Speech, Voice Activity Detection
- **[Core Utilities](./core-utils.md)** - useExecutorchModule, ResourceFetcher, Error Handling

**Always check these references before implementing** to ensure you have the most accurate and detailed information for the specific feature you're working with.

## External References

- **Documentation:** https://docs.swmansion.com/react-native-executorch
- **HuggingFace Models:** https://huggingface.co/software-mansion/collections
- **GitHub:** https://github.com/software-mansion/react-native-executorch
- **API Reference:** https://docs.swmansion.com/react-native-executorch/docs/api-reference

## Version Notes

This skill is based on React Native Executorch library documentation. Always verify the latest API changes and model availability at the official documentation before implementation.
