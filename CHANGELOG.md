# Changelog

## [0.0.6] - 2024-03-12

### Added
- Added force rebuild flag (--force) to regenerate transcription and summary files
- Enhanced script output to display selected language and model information
- Switched to .md file extension for better markdown compatibility

### Changed
- Refactored language-specific model configuration
  - English now uses consistent configuration format (WHISPER_MODEL_EN)
  - Maintained special faster-whisper-medium-en-cpu model for English
  - Other languages use systran-faster-whisper-large-v3 universal model
- Improved file organization by keeping all files in source directory
- Updated configuration format to support language-specific models

### Fixed
- Fixed script name references (removed .sh suffix)
- Improved help text to reflect new configuration options

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
