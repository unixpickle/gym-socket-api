# gym-socket-api

This is an API for accessing [OpenAI Gym](https://gym.openai.com) from other programming languages. It is intended to scale up to difficult, high-dimensional environments without being slow.

**Disclaimer:** I don't know Python.

# Installation

To download the code, run:

```
git clone https://github.com/unixpickle/gym-socket-api
```

To install the dependencies, run:

```
pip install -r requirements.txt
```

# Usage

## Server

You can spin up a new server like so:

```
python .
```

By default, this will listen on port 5001. You can change the port using the `--port` flag:

```
python . --port 1337
```

## Client

**Go client:** See [binding-go/demo/cartpole](binding-go/demo/cartpole) to jump right into the Go bindings. If you'd like a more comprehensive guide, see the [Godoc](https://godoc.org/github.com/unixpickle/gym-socket-api/binding-go).

# Why not openai/gym-http-api?

There are already official language bindings for OpenAI Gym in [openai/gym-http-api](https://github.com/openai/gym-http-api). Here are some reasons why gym-socket-api is still necessary:

**Performance:** games like Atari Pong generate video frames many times a second. It should be possible to play games like this faster than real-time. Since gym-http-api is committed to JSON, serializing and deserializing video frames is extremely slow. In contrast, gym-socket-api can achieve over 300 FPS on games like Atari Pong.

**Session management:** when a client process exits ungracefully (e.g. from Ctrl+C), the language bindings should clean up any resources that the process allocated. This means that, for example, any GUI windows related to the environment should be closed. Since gym-http-api does not use persistent sockets or timeouts, it cannot do this.

**Stability:** I've had the gym-http-api server crash on me before. With gym-socket-api, each client gets its own server process. This way, even a segmentation fault can't bring down the server.
