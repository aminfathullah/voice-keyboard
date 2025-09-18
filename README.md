# Voice Keyboard

A Go application with GUI that converts speech to text and types it into the active text field using Google Cloud Speech-to-Text.

## Prerequisites

- Go 1.21 or later
- Google Cloud Project with Speech-to-Text API enabled
- Service Account Key JSON file
- PortAudio and xdotool installed (see installation)

## Installation

1. Clone or download the project.

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install system packages:
   ```bash
   sudo apt update
   sudo apt install portaudio19-dev xdotool
   ```

4. Set up Google Cloud:
   - Create a Google Cloud Project
   - Enable the Speech-to-Text API
   - Create a Service Account and download the JSON key

## Usage

1. Run the application:
   ```bash
   go run main.go
   ```

2. A GUI window will open with the following controls:
   - **Select Service Account File**: Click to choose your Google Cloud service account JSON file
   - **Start Voice Recognition**: Click to begin speech-to-text recognition
   - **Stop Voice Recognition**: Click to stop recognition
   - **Language Selection**: Choose a language for transcription (e.g., `en-US`, `id-ID`). The GUI defaults to `id-ID`.
   - **Hotkey Assignment**: Enter a hotkey combination (e.g., Ctrl+Space) - currently for display only

3. Select your service account JSON file using the file picker

4. Click "Start Voice Recognition" to begin

5. Focus on a text field (e.g., in a browser, editor, etc.)

6. Start speaking. The application will transcribe your speech and type the text into the focused field.

7. Click "Stop Voice Recognition" when done

## GUI Features

- **Service Account Selection**: File picker to select Google Cloud credentials
- **Start/Stop Controls**: Buttons to control voice recognition streaming
- **Status Display**: Shows current state of the application
- **Hotkey Input**: Field for future hotkey assignment functionality

## How it works


## Notes


## Automated Releases (GitHub Releases)

This repository includes a GitHub Actions workflow that builds release binaries and attaches them to a GitHub Release when you push a tag matching `v*.*.*` (for example `v0.1.0`).

How it works:
- When you push a tag like `v1.2.3`, the workflow builds binaries for Linux, macOS, and Windows (amd64 and arm64) and attaches them to the GitHub Release created for that tag.

Trigger a release locally:

```bash
# create a signed or lightweight tag (example lightweight)
git tag v0.1.0
git push origin v0.1.0
```

After the workflow finishes you can download the compiled binaries from the GitHub Releases page for your repository.

Notes:
- Release artifacts contain compiled binaries only. They do not include your service account JSON (which should remain private).
- If you want to customize targets or build flags, edit `.github/workflows/release-build.yml`.
