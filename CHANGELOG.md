# Changelog

## [0.0.5] - 2024-03-12

### Added
- Support for automatic language detection using Systran/faster-whisper-large-v3
- Language selection option (--language) for backward compatibility
- Comprehensive KubeAI installation documentation

### Changed
- Switched to Systran/faster-whisper-large-v3 as the universal model for all languages
- Removed language parameter from transcription API for automatic language detection
- Updated KubeAI model configuration to use "FasterWhisper" engine
- Simplified model configuration by using a single universal model
- Maintained faster-whisper-medium-en-cpu as optional model for backward compatibility

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
