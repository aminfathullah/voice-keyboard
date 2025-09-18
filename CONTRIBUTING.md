# Contributing to Voice Keyboard

Thanks for your interest in contributing! This file explains how to set up the project locally, coding conventions, and how to open a PR.

## Running locally

Prerequisites:
- Go 1.21+
- PortAudio development headers (`portaudio19-dev` on Debian/Ubuntu)
- `xdotool` installed (Linux only)

Install dependencies:
```bash
go mod tidy
```

Run the app:
```bash
# Option 1: run directly (recommended during development)
go run main.go

# Option 2: build and run
go build -o voice-keyboard
./voice-keyboard
```

Set Google credentials (choose one):
```bash
# Temporary for session
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# Or select via GUI at runtime by clicking "Select Service Account File"
```

## Coding style
- Follow Go idioms and `gofmt` formatting
- Keep UI and business logic separated where possible

## Pull requests
- Fork and create a topic branch
- Keep PRs small and focused
- Include tests for new behavior
- Describe security implications, especially if touching credentials handling

## Issues
Open an issue with steps to reproduce, expected vs actual behavior, and system info (OS, Go version).
