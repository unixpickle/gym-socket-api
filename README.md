# gym-socket-api

This will be a fast and stable API for accessing [OpenAI Gym](https://gym.openai.com) from other programming languages. My main goal is to create bindings that can scale up to difficult, high-dimensional environments without being slow.

# Why not use the official bindings?

There are already official language bindings for OpenAI Gym in [openai/gym-http-api](https://github.com/openai/gym-http-api). Here are some reasons why gym-socket-api is still necessary:

 * Performance &mdash; games like Atari Pong generate video frames many times a second. It should be possible to play games like this faster than real-time. Since gym-http-api is committed to JSON, serializing and deserializing video frames is extremely slow. This makes training times **significantly worse** than they could be.
 * Session management &mdash; when a client process exits ungracefully (e.g. from Ctrl+C), the language bindings should clean up any resources that the process allocated. This means that, for example, any GUI windows related to the environment should be closed. Since gym-http-api does not use persistent sockets or timeouts, it cannot do this.
 * Stability &mdash; I've had the gym-http-api server crash on me before. It would be ideal if this did not happen.
