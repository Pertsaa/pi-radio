package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Player struct {
	cmd *exec.Cmd
}

// NewPlayer creates a new player instance.
func NewPlayer() *Player {
	return &Player{}
}

// Play a given MP3 file with a specified volume.
// Volume is a percentage from 0 to 100.
func (p *Player) Play(filePath string, volume int) error {
	// Map volume percentage to mpg123's scale (0-32768).
	scale := int(float64(volume) / 100.0 * 32768)
	if scale < 0 {
		scale = 0
	} else if scale > 32768 {
		scale = 32768
	}

	// Use -q for quiet output and -f for volume scaling.
	p.cmd = exec.Command("mpg123", "-f", fmt.Sprint(scale), "-q", filePath)
	return p.cmd.Start()
}

// Stop the current playback.
func (p *Player) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	// Use SIGTERM to gracefully terminate the process.
	return p.cmd.Process.Signal(syscall.SIGTERM)
}

func main() {
	// Check if a file path argument is provided.
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run your_program.go <path_to_mp3_file>")
		os.Exit(1)
	}

	// The first argument (os.Args[1]) is the file path.
	mp3File := os.Args[1]

	player := NewPlayer()

	// Play at 75% volume.
	fmt.Println("Playing audio at 75% volume...")
	err := player.Play(mp3File, 75)
	if err != nil {
		fmt.Println("Error playing file:", err)
		return
	}

	// Wait for a few seconds.
	time.Sleep(5 * time.Second)

	// Stop playback.
	fmt.Println("Stopping playback.")
	err = player.Stop()
	if err != nil {
		fmt.Println("Error stopping playback:", err)
	}
}
