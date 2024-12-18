# Changelog

## [Unreleased]

### Added
- Automatic model downloads from HuggingFace for both local and KubeAI backends
- KubeAI-compatible catalog-based configuration structure for model management
- Support for both local and KubeAI transcription backends with consistent configuration
- Language-specific model configuration with auto-detection support

### Changed
- Simplified model setup by removing manual download requirements
- Enhanced model configuration with unified catalog format
- Improved language handling with automatic model selection

## [0.1.0] - 2024-01-17

### Added
- Complete rewrite in Go with improved architecture
- Modular package structure:
  - config: Configuration management
  - transcribe: Audio transcription with language support
  - summarize: Text summarization with OpenAI integration
  - models: Whisper model management
  - utils: File and media handling utilities
  - process: Complete workflow management
- Comprehensive test suite with mock clients
- Type-safe configuration management
- Improved error handling and logging
- Better dependency management with go modules

### Changed
- Moved from Bash script to Go binary
- Improved file organization (removed temp directory usage)
- Changed output files from .txt to .md
- Restructured configuration format
- Enhanced language model configuration
- Added force rebuild option

### Fixed
- Improved error handling for API calls
- Better file path handling
- More robust configuration loading
- Enhanced input validation
- Fixed transcription using mock client instead of real API
- Fixed config file not being created with default values
- Fixed ChatGPT model configuration to use gpt-4o
- Fixed file extension for transcripts to use .md instead of .json
- Improved README with detailed configuration documentation

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
