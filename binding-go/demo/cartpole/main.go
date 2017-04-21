package main

import (
	"time"

	gym "github.com/unixpickle/gym-socket-api/binding-go"
)

const Host = "localhost:5001"

func main() {
	conn, err := gym.Make(Host, "CartPole-v0")
	must(err)

	defer conn.Close()

	// TODO: stuff here.
	time.Sleep(time.Second * 2)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
