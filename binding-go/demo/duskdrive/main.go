package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/unixpickle/essentials"
	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const (
	Host   = "localhost:5001"
	Frames = 50
)

func main() {
	// Make an instance of the given environment.
	env, err := gym.Make(Host, "flashgames.DuskDrive-v0")
	if err != nil {
		essentials.Die(err.Error() +
			" (install Universe and run server with --universe)")
	}

	// Gracefully clean-up when we exit (although this
	// isn't strictly necessary).
	defer env.Close()

	// Limit observations to game screen.
	must(env.UniverseWrap("CropObservations", nil))
	must(env.UniverseWrap("Vision", nil))

	// Tell the environment how to run.
	must(env.UniverseConfigure(map[string]interface{}{"remotes": 1}))

	// Reset the environment and get initial observation.
	log.Println("Resetting...")
	lastObs, err := env.Reset()
	must(err)

	// Play for a bit and record the FPS.
	log.Println("Running...")
	startTime := time.Now().UnixNano()
	for i := 0; i < Frames; i++ {
		action := []interface{}{[]interface{}{"KeyEvent", "ArrowUp", true}}
		lastObs, _, _, _, err = env.Step(action)
		must(err)
	}
	seconds := float64(time.Now().UnixNano()-startTime) / 1e9
	fmt.Println("Played at", Frames/seconds, "fps")

	// Save a screenshot to screenshot.png.
	var rawFrame [][][]float64
	must(lastObs.Unmarshal(&rawFrame))
	width, height := len(rawFrame[0]), len(rawFrame)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i, pixel := range lastObs.(gym.Uint8Obs).Uint8Obs() {
		img.Pix[4*(i/3)+(i%3)] = pixel

		// Set alpha to 1.
		img.Pix[4*(i/3)+3] = 0xff
	}
	out, err := os.Create("screenshot.png")
	must(err)
	must(png.Encode(out, img))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
