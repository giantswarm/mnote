# mnote

mnote is a command-line tool for transcribing and summarizing video files using Whisper and ChatGPT. Written in Go, it provides robust audio processing, transcription, and AI-powered summarization capabilities.

## Features

- Video to text transcription using configurable Whisper models
- Text summarization using ChatGPT with customizable prompts
- Support for multiple languages (English, German, Spanish, French, and auto-detection)
- Language-specific model selection
- Force rebuild option for regenerating transcripts and summaries
- Supports various video formats (.mp4, .mkv, .avi, .mov)

## Prerequisites

- Go 1.21 or later
- FFmpeg (for audio extraction)
- OpenAI API key
- Build dependencies for local transcription:
  - gcc
  - cmake
  - pkg-config
- One of:
  - KubeAI installation for remote transcription service
  - Local Whisper model for local transcription

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/giantswarm/mnote.git
cd mnote
```

2. Install Go dependencies:
```bash
go mod download
```

3. Build the binary:
```bash
go build -o mnote ./cmd/mnote
```

4. Move the binary to your PATH:
```bash
sudo mv mnote /usr/local/bin/
```

### Configuration

The configuration file is automatically created at `~/.config/mnote/config` with default values on first run. You can also manually create or modify it:

```bash
# Default configuration file (~/.config/mnote/config)
# Transcription backend configuration
TRANSCRIPTION_BACKEND=kubeai  # Use 'kubeai' or 'local'
TRANSCRIPTION_API_URL=http://localhost:8000/v1/audio/transcriptions  # Only for kubeai backend

# Local Whisper configuration (only for local backend)
LOCAL_MODEL_PATH=~/.config/mnote/models/ggml-base.en.bin
LOCAL_MODEL_SIZE=base  # tiny, base, small, medium, or large

# Language configuration
DEFAULT_LANGUAGE=auto

# Language-specific model configuration (only for kubeai backend)
WHISPER_MODEL_EN=faster-whisper-medium-en-cpu    # Optimized for English
WHISPER_MODEL_DE=systran-faster-whisper-large-v3 # Universal model for German
WHISPER_MODEL_ES=systran-faster-whisper-large-v3 # Universal model for Spanish
WHISPER_MODEL_FR=systran-faster-whisper-large-v3 # Universal model for French

# ChatGPT configuration
CHATGPT_MODEL=gpt-4o
```

### Local Whisper Models

When using the local transcription backend (`TRANSCRIPTION_BACKEND=local`), you need to download and configure a Whisper model. The models have different sizes with varying accuracy and resource requirements:

| Model  | Disk   | Memory  | Accuracy | Use Case |
|--------|--------|---------|----------|----------|
| tiny   | 75 MB  | ~273 MB | Basic    | Quick transcriptions, low resource usage |
| base   | 142 MB | ~388 MB | Good     | General purpose, balanced performance |
| small  | 466 MB | ~852 MB | Better   | Improved accuracy, moderate resources |
| medium | 1.5 GB | ~2.1 GB | High     | High accuracy, higher resource usage |
| large  | 2.9 GB | ~3.9 GB | Best     | Best accuracy, significant resources |

To set up local transcription:

1. Install build dependencies:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install gcc cmake pkg-config

   # macOS
   brew install cmake pkg-config
   ```

2. Download a Whisper model:
   ```bash
   # Create models directory
   mkdir -p ~/.config/mnote/models
   cd ~/.config/mnote/models

   # Download your chosen model (example: base model)
   wget https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin
   ```

3. Configure mnote:
   ```bash
   # Update ~/.config/mnote/config
   TRANSCRIPTION_BACKEND=local
   LOCAL_MODEL_PATH=~/.config/mnote/models/ggml-base.bin
   LOCAL_MODEL_SIZE=base
   ```

### Prompts

Create custom prompts in `~/.config/mnote/prompts/`. The default summarization prompt is automatically created at `~/.config/mnote/prompts/summarize`:

```bash
# Default summary prompt
Create a detailed summary of the following meeting transcript. Structure the summary according to the main topics discussed and organize the information into logical sections. For each topic, summarize who was involved, what was discussed in detail, what decisions were made, what problems or challenges were identified, and what solutions were proposed or implemented.
```

## Usage

### Basic Command

```bash
mnote <video_directory>
```

- **`<video_directory>`**: Path to the directory containing video files.

### Options

- `--prompt <prompt_name>`: Use a custom prompt file from `~/.config/mnote/prompts`.
- `--language <lang_code>`: Specify the language for transcription (de, es, fr, or auto).
                          Defaults to "auto" for automatic detection.
- `--help`: Display the help message.

### Examples

#### Summarize a Directory of Videos

```bash
mnote /path/to/videos
```

Uses the default prompt (`summarize`) to process all supported video files in
the directory.

#### Use a Custom Prompt

```bash
mnote --prompt meeting /path/to/videos
```

Uses the custom prompt file `~/.config/mnote/prompts/meeting` for summarization.

#### Specify Language for Transcription

```bash
mnote --language de /path/to/videos     # German
mnote --language es /path/to/videos     # Spanish
mnote --language fr /path/to/videos     # French
mnote --language auto /path/to/videos   # Auto-detect language
```

## How It Works

1. **Audio Extraction**:
   The tool uses `ffmpeg` to extract audio from video files, saving the `.mp3` file
   in the same directory as the source video.

2. **Transcription**:
   Audio files are processed using one of two backends:

   a. **KubeAI Backend** (default):
      - Audio files are sent to a Whisper-based transcription API
      - Uses language-specific models:
        - English content uses the faster-whisper-medium-en-cpu model
        - Other languages use the Systran-faster-whisper-large-v3 universal model
        - Auto-detection intelligently selects the appropriate model

   b. **Local Backend**:
      - Uses whisper.cpp for local transcription
      - Supports multiple model sizes (tiny to large)
      - Provides offline transcription capability
      - Auto-detection uses the configured model size

   Transcriptions are saved as `.md` files alongside the source video.

3. **Summarization**:
   Transcriptions are processed using the OpenAI API with the configured
   ChatGPT model (gpt-4o by default) and specified prompt. If a summary file already exists
   for a video and --force is not used, the summarization step is skipped to avoid
   unnecessary API calls.

4. **Output**:
   Summarized meeting notes are saved as `.md` files in the same directory
   as the input videos. When using custom prompts, the prompt name is included
   in the output filename (e.g., `video_meeting.md` for the "meeting" prompt).
   The default "summarize" prompt maintains the original filename format
   (e.g., `video.md`).

## Supported File Formats

- `.mp4`
- `.mkv`
- `.avi`
- `.mov`

## Dependencies

Ensure the following tools are installed:

- `ffmpeg`: [Installation Guide](https://ffmpeg.org/download.html)
- `curl`: [Installation Guide](https://curl.se/)
- `jq`: [Installation Guide](https://stedolan.github.io/jq/download/)
- `chatgpt`: Install from [chatgpt-cli](https://github.com/kardolus/chatgpt-cli)

## Notes

- **KubeAI Installation**: The transcription service requires KubeAI with appropriate Whisper models.
  Follow these steps to set up the service:

  1. Add the KubeAI Helm repository:
     ```bash
     helm repo add kubeai https://charts.kubeai.com
     helm repo update
     ```

  2. Install KubeAI (version 0.9.0):
     ```bash
     helm install kubeai kubeai/kubeai --version 0.9.0
     ```

  3. Create a values file (`kubeai-models.yaml`):
     ```yaml
     catalog:
       # Default English model (optimized for English content)
       faster-whisper-medium-en-cpu:
         enabled: true
         features: ["SpeechToText"]
         owner: "Systran"
         url: "hf://Systran/faster-whisper-medium-en"
         engine: "FasterWhisper"
         resourceProfile: "cpu:1"
         minReplicas: 1

       # Universal model for multilingual support and auto-detection
       systran-faster-whisper-large-v3:
         enabled: true
         features: ["SpeechToText"]
         owner: "Systran"
         url: "hf://Systran/faster-whisper-large-v3"
         engine: "FasterWhisper"
         resourceProfile: "cpu:2"
         minReplicas: 1
     ```

  4. Install the models (version 0.9.0):
     ```bash
     helm install kubeai-models kubeai/models --version 0.9.0 -f kubeai-models.yaml
     ```

  5. Configure mnote:
     Update the `TRANSCRIPTION_API_URL` in your configuration to point to your KubeAI service endpoint.
     ```bash
     # Example configuration
     TRANSCRIPTION_BACKEND=kubeai
     TRANSCRIPTION_API_URL=http://kubeai.your-domain/v1/audio/transcriptions
     ```

  6. Port forward the service (for local development):
     ```bash
     kubectl port-forward svc/kubeai 8000:80 -n kubeai
     ```
     Then use `http://localhost:8000/v1/audio/transcriptions` as your `TRANSCRIPTION_API_URL`.

- **OpenAI API**: You must have an OpenAI API key for the `chatgpt` CLI tool.
  Register at [OpenAI](https://platform.openai.com/).

Timo Derstappen

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
