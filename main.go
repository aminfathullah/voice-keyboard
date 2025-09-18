package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/gordonklaus/portaudio"
)

type VoiceKeyboard struct {
	client             *speech.Client
	stream             *portaudio.Stream
	audioChan          chan []byte
	streamClient       speechpb.Speech_StreamingRecognizeClient
	ctx                context.Context
	cancel             context.CancelFunc
	running            bool
	serviceAccountFile string
	language           string
	lastTyped          string
	lastTypedMu        sync.Mutex
}

func NewVoiceKeyboard() *VoiceKeyboard {
	return &VoiceKeyboard{
		audioChan: make(chan []byte, 10),
		running:   false,
		language:  "id-ID",
	}
}

func (vk *VoiceKeyboard) SetServiceAccountFile(filePath string) {
	vk.serviceAccountFile = filePath
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", filePath)
}

func (vk *VoiceKeyboard) SetLanguage(lang string) {
	if lang == "" {
		return
	}
	vk.language = lang
}

func (vk *VoiceKeyboard) Initialize() error {
	vk.ctx, vk.cancel = context.WithCancel(context.Background())

	var err error
	vk.client, err = speech.NewClient(vk.ctx)
	if err != nil {
		return fmt.Errorf("failed to create speech client: %v", err)
	}

	err = portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %v", err)
	}

	return nil
}

func (vk *VoiceKeyboard) Start() error {
	if vk.running {
		return fmt.Errorf("already running")
	}

	var err error
	vk.stream, err = portaudio.OpenDefaultStream(1, 0, 16000, 64, func(in []int16, out []int16) {
		audioData := make([]byte, len(in)*2)
		for i, sample := range in {
			audioData[i*2] = byte(sample)
			audioData[i*2+1] = byte(sample >> 8)
		}
		select {
		case vk.audioChan <- audioData:
		default:
		}
	})
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}

	err = vk.stream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}

	// Set up streaming request
	// Use selected language
	req := &speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    vk.language,
				},
				InterimResults: true,
			},
		},
	}

	vk.streamClient, err = vk.client.StreamingRecognize(vk.ctx)
	if err != nil {
		return fmt.Errorf("failed to create streaming client: %v", err)
	}

	// Send initial config
	if err := vk.streamClient.Send(req); err != nil {
		return fmt.Errorf("failed to send config: %v", err)
	}

	// Start audio processing goroutine
	go vk.processAudio()

	// Start response processing goroutine
	go vk.processResponses()

	vk.running = true
	return nil
}

func (vk *VoiceKeyboard) Stop() error {
	if !vk.running {
		return fmt.Errorf("not running")
	}

	vk.running = false
	vk.cancel()

	if vk.stream != nil {
		vk.stream.Stop()
		vk.stream.Close()
	}

	if vk.streamClient != nil {
		vk.streamClient.CloseSend()
	}

	if vk.client != nil {
		vk.client.Close()
	}

	portaudio.Terminate()
	close(vk.audioChan)

	return nil
}

func (vk *VoiceKeyboard) processAudio() {
	for audioData := range vk.audioChan {
		req := &speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: audioData,
			},
		}

		if err := vk.streamClient.Send(req); err != nil {
			log.Printf("Error sending audio: %v", err)
			return
		}
	}
}

func (vk *VoiceKeyboard) processResponses() {
	for {
		select {
		case <-vk.ctx.Done():
			return
		default:
			resp, err := vk.streamClient.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Error receiving response: %v", err)
				return
			}

			for _, result := range resp.Results {
				if result.IsFinal {
					for _, alt := range result.Alternatives {
						fmt.Printf("Transcript: %s\n", alt.Transcript)
						vk.safeTypeText(alt.Transcript)
					}
				} else if len(result.Alternatives) > 0 {
					alt := result.Alternatives[0]
					if alt.Confidence > 0.5 {
						fmt.Printf("Interim: %s\n", alt.Transcript)
						vk.safeTypeText(alt.Transcript)
					}
				}
			}
		}
	}
}

func (vk *VoiceKeyboard) safeTypeText(text string) {
	if text == "" {
		return
	}
	vk.lastTypedMu.Lock()
	if text == vk.lastTyped {
		vk.lastTypedMu.Unlock()
		return
	}
	vk.lastTyped = text
	vk.lastTypedMu.Unlock()

	esc := strings.ReplaceAll(text, "'", "'\\''")
	cmd := exec.Command("bash", "-lc", fmt.Sprintf("xdotool type '%s'", esc))
	if err := cmd.Run(); err != nil {
		log.Printf("Error typing text: %v", err)
	}
}

func main() {
	vk := NewVoiceKeyboard()

	// Create Fyne app
	a := app.New()
	w := a.NewWindow("Voice Keyboard")

	// Service account file path label and button
	serviceAccountLabel := widget.NewLabel("No service account file selected")
	serviceAccountButton := widget.NewButton("Select Service Account File", func() {
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Printf("Error selecting file: %v", err)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			filePath := reader.URI().Path()
			vk.SetServiceAccountFile(filePath)
			serviceAccountLabel.SetText("Service Account: " + filePath)
		}, w)
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dialog.Show()
	})

	// Start/Stop buttons
	startButton := widget.NewButton("Start Voice Recognition", nil)
	stopButton := widget.NewButton("Stop Voice Recognition", nil)
	stopButton.Disable()

	// Set up button callbacks
	startButton.OnTapped = func() {
		if err := vk.Initialize(); err != nil {
			log.Printf("Failed to initialize: %v", err)
			return
		}
		if err := vk.Start(); err != nil {
			log.Printf("Failed to start: %v", err)
			return
		}
		startButton.Disable()
		stopButton.Enable()
		log.Println("Voice recognition started")
	}

	stopButton.OnTapped = func() {
		if err := vk.Stop(); err != nil {
			log.Printf("Failed to stop: %v", err)
			return
		}
		startButton.Enable()
		stopButton.Disable()
		log.Println("Voice recognition stopped")
	}

	// Status label
	statusLabel := widget.NewLabel("Status: Stopped")

	// Language selector
	languages := []string{"en-US", "id-ID", "es-ES", "fr-FR", "de-DE", "ja-JP", "zh-CN"}
	langSelect := widget.NewSelect(languages, func(s string) {
		vk.SetLanguage(s)
		statusLabel.SetText("Language: " + s)
	})
	langSelect.SetSelected(vk.language)

	// Hotkey input field
	hotkeyEntry := widget.NewEntry()
	hotkeyEntry.SetPlaceHolder("Enter hotkey (e.g., Ctrl+Space)")

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Voice Keyboard Control Panel"),
		serviceAccountLabel,
		serviceAccountButton,
		container.NewHBox(startButton, stopButton),
		widget.NewLabel("Language:"),
		langSelect,
		widget.NewLabel("Hotkey Assignment:"),
		hotkeyEntry,
		statusLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
