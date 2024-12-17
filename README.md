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
- KubeAI installation for transcription service

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

1. Create configuration directory:
```bash
mkdir -p ~/.config/mnote/prompts
```

2. Create default configuration file:
```bash
cat > ~/.config/mnote/config << EOF
# Default language setting (auto, en, de, es, fr)
DEFAULT_LANGUAGE=auto

# Language-specific model configuration
# English model (optimized for English content)
WHISPER_MODEL_EN=faster-whisper-medium-en-cpu

# Other language models (using universal model)
WHISPER_MODEL_DE=systran-faster-whisper-large-v3
WHISPER_MODEL_ES=systran-faster-whisper-large-v3
WHISPER_MODEL_FR=systran-faster-whisper-large-v3

# Transcription API URL (update with your KubeAI endpoint)
TRANSCRIPTION_API_URL=http://localhost:8000/v1/audio/transcriptions

# OpenAI API Key
OPENAI_API_KEY=your_openai_api_key

# ChatGPT Model for Summarization
CHATGPT_MODEL=gpt-4
EOF
```

3. Create default summary prompt:
```bash
cat > ~/.config/mnote/prompts/summarize << EOF
Create a detailed summary of the following meeting transcript. Structure the summary according to the main topics discussed and organize the information into logical sections. For each topic, summarize who was involved, what was discussed in detail, what decisions were made, what problems or challenges were identified, and what solutions were proposed or implemented.
EOF
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
   Audio files are sent to a Whisper-based transcription API specified in the
   configuration (`TRANSCRIPTION_API_URL`). The script uses language-specific models:
   - English content uses the faster-whisper-medium-en-cpu model by default
   - Other languages use the Systran/faster-whisper-large-v3 universal model
   - Auto-detection (default) intelligently selects the appropriate model
   - Saves transcription results as `.json` files alongside the source video

3. **Summarization**:
   Transcriptions are processed using the `chatgpt` CLI tool with the
   specified ChatGPT model and prompt. If a summary file already exists
   for a video, the ChatGPT processing step is skipped to avoid
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

  3. Create a values file for the models:
     ```bash
     cat > values.yaml <<EOF
     catalog:
       # Universal Whisper model with multilingual support and auto-detection
       systran-faster-whisper-large-v3:
         enabled: true
         features: ["SpeechToText"]
         owner: "Systran"
         url: "hf://Systran/faster-whisper-large-v3"
         engine: "FasterWhisper"
         resourceProfile: "cpu:2"
         minReplicas: 1

       # Default English model (for backward compatibility)
       faster-whisper-medium-en-cpu:
         enabled: true
         features: ["SpeechToText"]
         owner: "Systran"
         url: "hf://Systran/faster-whisper-medium-en"
         engine: "FasterWhisper"
         resourceProfile: "cpu:1"
         minReplicas: 1
     EOF
     ```

  4. Install the models (version 0.9.0):
     ```bash
     helm install kubeai-models kubeai/models --version 0.9.0 -f values.yaml
     ```

  5. Configure mnote:
     Update the `TRANSCRIPTION_API_URL` in your configuration to point to your KubeAI service endpoint.

- **OpenAI API**: You must have an OpenAI API key for the `chatgpt` CLI tool.
  Register at [OpenAI](https://platform.openai.com/).

## Author

Timo Derstappen

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
