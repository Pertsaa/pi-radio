package radio

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type Radio struct {
	AudioDir   string
	AudioFiles []AudioFile
	Volume     float64
	Paused     bool
	StreamFile *AudioFile
	Streamer   beep.StreamSeekCloser
	PauseCtrl  *beep.Ctrl
	VolumeCtrl *effects.Volume
}

func New(audioDir string) *Radio {
	return &Radio{AudioDir: audioDir}
}

type AudioFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r *Radio) ScanAudioFiles() error {
	entries, err := os.ReadDir(r.AudioDir)
	if err != nil {
		return err
	}

	r.AudioFiles = []AudioFile{}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mp3") {
			r.AudioFiles = append(r.AudioFiles, AudioFile{ID: uuid.New().String(), Name: entry.Name()})
		}
	}

	return nil
}

func (r *Radio) Play(audioFileID string) error {
	err := r.Stop()
	if err != nil {
		return err
	}

	var audioFile AudioFile
	for _, file := range r.AudioFiles {
		if file.ID == audioFileID {
			audioFile = file
		}
	}
	if audioFile.ID == "" {
		return fmt.Errorf("audio file not found: %s", audioFileID)
	}

	f, err := os.Open(fmt.Sprintf("%s/%s", r.AudioDir, audioFile.Name))
	if err != nil {
		return err
	}

	fmt.Printf("playing: %s", audioFile.Name)

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}

	r.Streamer = streamer

	// err check here would fail on replay (speaker can only be initialized once)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	pauseCtrl := &beep.Ctrl{Streamer: streamer, Paused: false}

	r.PauseCtrl = pauseCtrl

	volumeCtrl := &effects.Volume{
		Streamer: pauseCtrl,
		Base:     2,
		Volume:   r.Volume,
		Silent:   false,
	}

	r.VolumeCtrl = volumeCtrl

	r.Paused = false
	r.StreamFile = &audioFile

	speaker.Play(volumeCtrl)

	return nil
}

func (r *Radio) SetPaused(paused bool) {
	if r.PauseCtrl == nil {
		return
	}

	speaker.Lock()
	r.PauseCtrl.Paused = paused
	r.Paused = paused
	speaker.Unlock()
}

func (r *Radio) SetVolume(volume float64) {
	if r.VolumeCtrl == nil {
		return
	}

	speaker.Lock()
	r.VolumeCtrl.Volume += volume
	speaker.Unlock()
}

func (r *Radio) Stop() error {
	if r.Streamer != nil {
		err := r.Streamer.Close()
		if err != nil {
			return err
		}
	}

	r.Streamer = nil
	r.VolumeCtrl = nil
	r.PauseCtrl = nil
	r.StreamFile = nil
	r.Paused = true

	return nil
}
