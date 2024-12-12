# Changelog

## [0.0.5] - 2024-03-12

### Added
- Support for multiple languages (German, Spanish, French)
- Language-specific Whisper models for optimized transcription
- Automatic language detection using openai/whisper-large-v3
- Language selection option (--language) with support for de, es, fr, and auto
- Comprehensive KubeAI installation documentation

### Changed
- Restored faster-whisper-medium-en-cpu as default model for English transcription
- Enhanced transcription API to support language-specific parameters

## [0.0.4] - 2024-12-12

- Enhanced prompt handling to support multiple summaries of the same video using different prompts
- Modified output filename to include prompt name when using custom prompts

## [0.0.3] - 2024-12-12

- Added file existence check to skip unnecessary ChatGPT processing when summary files already exist.

## [0.0.1] - 2024-12-01

- Initial release of mnote.
- Added transcription and summarization functionality.
- Configurable Whisper and ChatGPT models.
