package main

import (
	"fmt"

	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const Host = "localhost:5001"

func main() {
	conn, err := gym.Make(Host, "CartPole-v0")
	must(err)

	defer conn.Close()

	must(conn.Monitor("gym-monitor", true, false))

	actionSpace, err := conn.ActionSpace()
	must(err)
	fmt.Printf("Action space: %#v\n", actionSpace)

	obsSpace, err := conn.ObservationSpace()
	must(err)
	fmt.Printf("Observation space: %#v\n", obsSpace)

	obs, err := conn.Reset()
	must(err)
	fmt.Println("Initial obs:", obs)

	must(conn.Render())

	for {
		var action int
		must(conn.SampleAction(&action))
		obs, rew, done, _, err := conn.Step(action)
		must(conn.Render())
		must(err)
		fmt.Printf("Step: rew=%f obs=%v\n", rew, obs)
		if done {
			break
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
