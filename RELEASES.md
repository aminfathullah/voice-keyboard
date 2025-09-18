Windows build notes

This repository's GitHub Actions workflow now builds Windows binaries using an MSYS2 environment on the `windows-latest` runner. A few important notes:

- The workflow installs MSYS2 packages including `mingw-w64-x86_64-portaudio` so the binary can link against PortAudio if needed.
- Building GUI apps that use Fyne may require additional native libraries on Windows. The workflow builds a console binary; if you need a bundled app (with icons, .msi or .exe installer), consider using `goreleaser` or an NSIS installer step.
- If you prefer to cross-compile from Linux, the project appears to have cgo/native dependencies; cross-compiling may fail. Using the provided Windows runner avoids cross-compilation issues.

Local Windows build (recommended using MSYS2):

1. Install MSYS2: https://www.msys2.org/
2. Install required packages (run in MSYS2 MINGW64 shell):

```bash
pacman -Syu
pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-portaudio
```

3. Build from the MINGW64 shell:

```bash
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
go build -o voice-keyboard-windows-amd64.exe
```

If your Windows target needs arm64 builds, set `GOARCH=arm64` and an appropriate `CC` (you may need additional toolchains).

Packaging for distribution

- Consider zipping each OS/arch binary and uploading the zip as the release asset.
- For richer installers on Windows, use `goreleaser` with an NSIS/Zip target.
