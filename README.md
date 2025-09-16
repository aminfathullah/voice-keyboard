# Voice Keyboard

A Go application that converts speech to text and types it into the active text field using Google Cloud Speech-to-Text.

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
   - Set the environment variable:
     ```bash
     export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
     ```

## Usage

1. Run the application:
   ```bash
   go run main.go
   ```

2. Focus on a text field (e.g., in a browser, editor, etc.)

3. Start speaking. The application will transcribe your speech and type the text into the focused field.

4. The program handles pauses automatically through streaming recognition.

## How it works

- Captures audio from the default microphone using PortAudio
- Streams audio to Google Cloud Speech-to-Text API
- Receives real-time transcription
- Uses xdotool to simulate keyboard typing into the active window

## Notes

- Ensure your microphone is working and permissions are granted
- The application types only final transcripts to avoid typing partial words
- For continuous use, keep the application running and switch focus as needed