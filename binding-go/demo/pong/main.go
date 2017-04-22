package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const (
	Host   = "localhost:5001"
	Frames = 50
)

func main() {
	// Make an instance of the given environment.
	env, err := gym.Make(Host, "Pong-v0")
	must(err)

	// Gracefully clean-up when we exit (although this
	// isn't strictly necessary).
	defer env.Close()

	// Get the screen dimensions.
	obsSpace, err := env.ObservationSpace()
	must(err)
	height, width := obsSpace.Shape[0], obsSpace.Shape[1]

	// Reset the environment and get initial observation.
	lastObs, err := env.Reset()
	must(err)

	// Play randomly for a bit and record the FPS.
	startTime := time.Now().UnixNano()
	for i := 0; i < Frames; i++ {
		var action int
		must(env.SampleAction(&action))
		lastObs, _, _, _, err = env.Step(action)
		must(err)
	}
	seconds := float64(time.Now().UnixNano()-startTime) / 1e9
	fmt.Println("Played at", Frames/seconds, "fps")

	// Save a screenshot to pong.png.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	data := lastObs.(gym.Uint8Obs).Uint8Obs()
	for i, pixel := range data {
		img.Pix[4*(i/3)+(i%3)] = pixel

		// Set alpha to 1.
		img.Pix[4*(i/3)+3] = 0xff
	}

	out, err := os.Create("pong.png")
	must(err)
	must(png.Encode(out, img))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
