package main

import (
	"fmt"
	"math/rand"

	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const Host = "localhost:5001"

func main() {
	conn, err := gym.Make(Host, "CartPole-v0")
	must(err)

	defer conn.Close()

	obs, err := conn.Reset()
	must(err)
	fmt.Println("Initial obs:", obs)

	for {
		action := rand.Intn(2)
		obs, rew, done, _, err := conn.Step(action)
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
