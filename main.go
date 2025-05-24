package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func loadMp3(filepath string) (*mp3.Decoder, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	fileReader := bytes.NewReader(file)

	decodedMp3, err := mp3.NewDecoder(fileReader)
	if err != nil {
		return nil, err
	}

	return decodedMp3, nil
}

func main() {
	// Need to load the track first since we need the sample rate
	track, err := loadMp3("[insert path here]")
	if err != nil {
		log.Fatalln(err)
	}

	contextOptions := &oto.NewContextOptions{}
	contextOptions.SampleRate = track.SampleRate()
	contextOptions.ChannelCount = 2
	contextOptions.Format = oto.FormatSignedInt16LE

	context, readyChan, err := oto.NewContext(contextOptions)
	if err != nil {
		log.Fatalln(err)
	}

	// Wait for the hardware audio devices to be ready
	<-readyChan

	player := context.NewPlayer(track)

	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
