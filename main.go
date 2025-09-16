package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/gordonklaus/portaudio"
)

func main() {
	// Initialize Google Cloud Speech client
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create speech client: %v", err)
	}
	defer client.Close()

	// Set up audio input
	err = portaudio.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	// Audio processing channel
	audioChan := make(chan []byte, 10)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, 1024, func(in []int16, out []int16) {
		// Convert to bytes
		audioData := make([]byte, len(in)*2)
		for i, sample := range in {
			audioData[i*2] = byte(sample)
			audioData[i*2+1] = byte(sample >> 8)
		}
		select {
		case audioChan <- audioData:
		default:
		}
	})
	if err != nil {
		log.Fatalf("Failed to open audio stream: %v", err)
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		log.Fatalf("Failed to start audio stream: %v", err)
	}
	defer stream.Stop()

	// Set up streaming request
	req := &speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
				},
				InterimResults: true,
			},
		},
	}

	streamClient, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatalf("Failed to create streaming client: %v", err)
	}

	// Send initial config
	if err := streamClient.Send(req); err != nil {
		log.Fatalf("Failed to send config: %v", err)
	}

	go func() {
		for audioData := range audioChan {
			req := &speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: audioData,
				},
			}

			if err := streamClient.Send(req); err != nil {
				log.Printf("Error sending audio: %v", err)
				return
			}
		}
	}()

	// Graceful shutdown handling
	done := make(chan struct{})

	// Receive responses
	go func() {
		defer close(done)
		for {
			resp, err := streamClient.Recv()
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
						// Type the text (debounced/safeguarded)
						safeTypeText(alt.Transcript)
					}
				}
			}
		}
	}()

	// Wait for interrupt (Ctrl+C) or streaming done
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("Shutting down: interrupt received")
	case <-done:
		log.Println("Shutting down: stream closed")
	}

	// Close streamClient to tell server we're done
	_ = streamClient.CloseSend()
	// Give a moment for resources to wind down
	time.Sleep(200 * time.Millisecond)
}

// safeTypeText escapes the text and avoids typing duplicates in quick succession.
var lastTypedMu sync.Mutex
var lastTyped string

func safeTypeText(text string) {
	if text == "" {
		return
	}
	lastTypedMu.Lock()
	if text == lastTyped {
		lastTypedMu.Unlock()
		return
	}
	lastTyped = text
	lastTypedMu.Unlock()

	// Escape single quotes for xdotool
	esc := strings.ReplaceAll(text, "'", "'\\''")
	cmd := exec.Command("bash", "-lc", fmt.Sprintf("xdotool type '%s'", esc))
	if err := cmd.Run(); err != nil {
		log.Printf("Error typing text: %v", err)
	}
}
