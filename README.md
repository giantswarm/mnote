# mnote

**mnote** is a CLI tool for summarizing meetings using AI. It transcribes audio from video files, processes the transcription using a Whisper-based API, and generates detailed summaries with ChatGPT.

## Features

- Extracts audio from video files.
- Transcribes audio using a configurable Whisper-based transcription service.
- Summarizes transcripts using ChatGPT with customizable prompts.
- Configurable transcription API, Whisper model, and ChatGPT model.

---

## Installation

### Prerequisites
Ensure the following tools are installed on your system and available in the `PATH`:
- `ffmpeg` (for audio extraction)
- `curl` (for API requests)
- `jq` (for processing JSON output)
- `chatgpt` ([chatgpt-cli](https://github.com/kardolus/chatgpt-cli) for summarization)

### Environment Variable

Set your OpenAI API key for the `chatgpt` tool:

```bash
export OPENAI_API_KEY="your_openai_api_key"
```

### Clone and Setup

Clone the repository:

```bash
git clone https://github.com/teemow/mnote.git
cd mnote
```

:> [!WARNING]
>
Make the script executable:

```bash
chmod +x mnote.sh
```

(Optional) Add it to your `PATH`:

```bash
sudo mv mnote.sh /usr/local/bin/mnote
```

---

## Configuration

### Default Configuration

Upon first run, **mnote** will create a configuration directory at
`~/.config/mnote` with the following structure:

```
~/.config/mnote/
├── config
└── prompts/
    └── summarize
```

### Configuration File (`~/.config/mnote/config`)

The `config` file contains the following default values:

```ini
# Transcription API URL
TRANSCRIPTION_API_URL=https://example.com/openai/v1/audio/transcriptions

# Whisper Model for Transcription
WHISPER_MODEL=faster-whisper-medium-en-cpu

# ChatGPT Model for Summarization
CHATGPT_MODEL=gpt-4o-2024-05-13
```

You can edit these values to customize the transcription API, Whisper model,
and ChatGPT model.

### Prompts

Prompts are stored in `~/.config/mnote/prompts`. The default prompt
(`summarize`) is created automatically:

```plaintext
Create a detailed summary of the following meeting transcript. Structure the summary according to the main topics discussed and organize the information into logical sections. For each topic, summarize who was involved, what was discussed in detail, what decisions were made, what problems or challenges were identified, and what solutions were proposed or implemented. If specific names are included in the transcript, use them to accurately attribute the statements. Also document all important feedback and planned actions. Pay attention to details on time frames, responsibilities, open questions and any next steps. Conclude the summary with a brief overview of the key findings and next steps.
```

To add a custom prompt, create a new file in the `prompts` Directory
(e.g., `meeting`) and reference it using the `--prompt` option.

---

## Usage

### Basic Command

```bash
mnote <video_directory>
```

- **`<video_directory>`**: Path to the directory containing video files.

### Options

- `--prompt <prompt_name>`: Use a custom prompt file from `~/.config/mnote/prompts`.
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

---

## How It Works

1. **Audio Extraction**:
   The tool uses `ffmpeg` to extract audio from video files.

2. **Transcription**:
   Audio files are sent to a Whisper-based transcription API specified in the
   configuration (`TRANSCRIPTION_API_URL`).

3. **Summarization**:
   Transcriptions are processed using the `chatgpt` CLI tool with the
   specified ChatGPT model and prompt. If a summary file already exists
   for a video, the ChatGPT processing step is skipped to avoid
   unnecessary API calls.

4. **Output**:
   Summarized meeting notes are saved as `.txt` files in the same directory
   as the input videos. When using custom prompts, the prompt name is included
   in the output filename (e.g., `video_meeting.txt` for the "meeting" prompt).
   The default "summarize" prompt maintains the original filename format
   (e.g., `video.txt`).

---

## Supported File Formats

- `.mp4`
- `.mkv`
- `.avi`
- `.mov`

---

## Dependencies

Ensure the following tools are installed:

- `ffmpeg`: [Installation Guide](https://ffmpeg.org/download.html)
- `curl`: [Installation Guide](https://curl.se/)
- `jq`: [Installation Guide](https://stedolan.github.io/jq/download/)
- `chatgpt`: Install from [chatgpt-cli](https://github.com/kardolus/chatgpt-cli)

---

## Notes

- **Transcription Service**: The transcription service is based on KubeAI,
  deployed via a Helm chart in a Kubernetes cluster. Ensure the service is
  properly configured and accessible via the `TRANSCRIPTION_API_URL` in the
  configuration file.
- **OpenAI API**: You must have an OpenAI API key for the `chatgpt` CLI tool.
  Register at [OpenAI](https://platform.openai.com/).

---

## Author

Timo Derstappen

---

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
