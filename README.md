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

- Captures audio from the default microphone using PortAudio
- Streams audio to Google Cloud Speech-to-Text API
- Receives real-time transcription
- Uses xdotool to simulate keyboard typing into the active window

## Notes

- Ensure your microphone is working and permissions are granted
- The application types only final transcripts to avoid typing partial words
- For continuous use, keep the application running and switch focus as needed
- The GUI provides easy control without needing command-line interaction