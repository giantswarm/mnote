# Changelog

## [0.0.5] - 2024-03-12

### Added
- Support for automatic language detection using Systran/faster-whisper-large-v3
- Language selection option (--language) for explicit language specification
- Comprehensive KubeAI installation documentation

### Changed
- Added Systran/faster-whisper-large-v3 for non-English languages and auto-detection
- Maintained faster-whisper-medium-en-cpu as default model for English
- Language parameter included in transcription API when language is explicitly specified
- Updated KubeAI model configuration to use "FasterWhisper" engine
- Changed output file format from .txt to .md for better markdown compatibility
- Keep all generated files (mp3, json, md) in same directory as source video for improved usability

### Fixed
- Corrected KubeAI model configuration format and engine settings

## [0.0.4] - 2024-12-12

- Enhanced prompt handling to support multiple summaries of the same video using different prompts
- Modified output filename to include prompt name when using custom prompts

## [0.0.3] - 2024-12-12

- Added file existence check to skip unnecessary ChatGPT processing when summary files already exist.

## [0.0.1] - 2024-12-01

- Initial release of mnote.
- Added transcription and summarization functionality.
- Configurable Whisper and ChatGPT models.
