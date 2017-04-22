package main

import (
	"fmt"

	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const Host = "localhost:5001"

func main() {
	// Make an instance of the given environment.
	env, err := gym.Make(Host, "CartPole-v0")
	must(err)

	// Gracefully clean-up when we exit (although this
	// isn't strictly necessary).
	defer env.Close()

	// Start monitoring to "./gym-monitor".
	must(env.Monitor("gym-monitor", true, false))

	// Dump info about the action space.
	actionSpace, err := env.ActionSpace()
	must(err)
	fmt.Printf("Action space: %#v\n", actionSpace)

	// Dump info about the observation space.
	obsSpace, err := env.ObservationSpace()
	must(err)
	fmt.Printf("Observation space: %#v\n", obsSpace)

	// Reset the environment and get initial observation.
	obs, err := env.Reset()
	must(err)
	fmt.Println("Initial obs:", obs)

	// Throw up a GUI window of the environment.
	must(env.Render())

	for {
		// Get a random action.
		var action int
		must(env.SampleAction(&action))

		// Take a step in the environment.
		obs, rew, done, _, err := env.Step(action)
		must(err)

		// Render the updated environment in the GUI.
		must(env.Render())

		// Print output of current step.
		fmt.Printf("Step: rew=%f obs=%v\n", rew, obs)
		if done {
			break
		}
	}

	// This is an example of how you might upload
	// your results to the OpenAI Gym.
	//
	//     env.Close()
	//     must(gym.Upload(Host, "gym-monitor", "", ""))
	//
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
