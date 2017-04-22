# Protocol

All integers are encoded in little endian. All strings are UTF-8. All floats are encoded according to [IEEE 754](https://en.wikipedia.org/wiki/IEEE_floating_point). Booleans (abbreviated "bool") are bytes; they are 0 for false, 1 for true.

Whenever there's an error field, it can be an empty string to indicate success.

## Initial connection

During this stage, the client initiates a connection and requests an environment. The server attempts to create the environment, or fails with an error (e.g. if the environment does not exist).

As a special case, the environment name may be the empty string. In this case, the client may not run any commands which act on an environment.

|Source   |Type    | Description           |
|---------|--------|-----------------------|
|Client   |uint8   | Flags (all 0)         |
|Client   |uint32  | Length of env name    |
|Client   |string  | Environment name      |
|Server   |uint32  | Error length          |
|Server   |string  | Error message         |

## Command packets

Once the handshake has completed, the client may send commands and receive responses. Only one command can be run at once. All packets take the following form:

|Source   |Type    | Description           |
|---------|--------|-----------------------|
|Client   |uint8   | Packet type           |
|Client   |varies  | Packet data           |

Valid packet types are listed below.

### Packet: Reset

This is packet type 0.

This packet resets the environment and gets the initial observation. It can be used as follows:

|Source   |Type                         | Description           |
|---------|-----------------------------|-----------------------|
|Client   |uint8                        | Packet type (0)       |
|Server   |[observation](#observations) | Initial observation   |

### Packet: Step

This is packet type 1.

This packet takes a step in the environment and gets a lot of information back. It can be used as follows:

|Source   |Type                         | Description           |
|---------|-----------------------------|-----------------------|
|Client   |uint8                        | Packet type (1)       |
|Client   |[action](#actions)           | Action to take        |
|Server   |[observation](#observations) | Next observation      |
|Server   |float64                      | Reward                |
|Server   |bool                         | Done                  |
|Server   |uint32                       | Info length           |
|Server   |string                       | Info JSON             |

### Packet: Get Space

This is packet type 2.

This packet gets information about the observation or action space. It can be used as follows:

|Source   |Type                  | Description           |
|---------|----------------------|-----------------------|
|Client   |uint8                 | Packet type (2)       |
|Client   |uint8                 | Which space?          |
|Server   |[space](#spaces)      | Space data            |

The "Which space?" field is 0 for the action space or 1 for the observation space.

### Packet: Sample Actions

This is packet type 3.

This packet samples an action from the action space. It can be used as follows:

|Source   |Type                  | Description           |
|---------|----------------------|-----------------------|
|Client   |uint8                 | Packet type (3)       |
|Server   |[action](#Actions)    | Action data           |

### Packet: Monitor

This is packet type 4.

This packet tells the server to wrap the current environment in a monitor. It can be used as follows:

|Source   |Type                  | Description           |
|---------|----------------------|-----------------------|
|Client   |uint8                 | Packet type (4)       |
|Client   |bool                  | Resume                |
|Client   |bool                  | Force                 |
|Client   |uint32                | Dir path length       |
|Client   |string                | Dir path              |

### Packet: Render

This is packet type 5.

This packet tells the server to render the current environment.

|Source   |Type                  | Description           |
|---------|----------------------|-----------------------|
|Client   |uint8                 | Packet type (5)       |

### Packet: Upload

This is packet type 6.

This packet tells the server to upload a monitor directory to OpenAI Gym.

|Source   |Type                  | Description           |
|---------|----------------------|-----------------------|
|Client   |uint8                 | Packet type (6)       |
|Client   |uint32                | Dir path length       |
|Client   |string                | Dir path              |
|Client   |uint32                | API key length        |
|Client   |string                | API key               |
|Client   |uint32                | Algorithm ID length   |
|Client   |string                | Algorithm ID          |
|Server   |uint32                | Error length          |
|Server   |string                | Error message         |

## Actions

Actions are encoded in a type-specific manner. They are of the form:

|Type    | Description           |
|--------|-----------------------|
|uint8   | Action type ID        |
|uint32  | Length of data        |
|varies  | Data                  |

The available action types are listed below.

### Action: JSON

This is action type 0.

The data inside the action is a JSON string for the space's `from_jsonable` method.

## Observations

Observations are encoded in a type-specific manner. They are of the form:

|Type    | Description           |
|--------|-----------------------|
|uint8   | Observation type ID   |
|uint32  | Length of data        |
|varies  | Data                  |

The available observation types are listed below.

### Observation: JSON

This is observation type 0.

The data in the packet is a JSON string from the space's `to_jsonable` method.

### Observation: Byte List

This is observation type 1.

The data in the packet is a flattened array of bytes. This observation has the following format:

|Type     | Description           |
|---------|-----------------------|
|uint32   | Num dimensions        |
|uint32[] | Dimensions            |
|uint8[]  | Data                  |

This is for observations in things like Atari environments where the observation is a raw 3D array of bytes. The array of bytes is flattened (in C order) into a 1D list of bytes.

## Spaces

Spaces are encoded using JSON:

|Type    | Description      |
|--------|------------------|
|uint32  | Length of data   |
|string  | JSON Data        |

The following examples demonstrate how each space type should be encoded.

Box spaces:

```json
{
  "type": "Box",
  "shape": [2, 3],
  "low": [-1, -1, -1, -1, -1, -1],
  "high": [1, 1, 1, 1, 1, 1]
}
```

Discrete spaces:

```json
{
  "type": "Discrete",
  "n": 5,
}
```

MultiBinary spaces:

```json
{
  "type": "MultiBinary",
  "n": 5,
}
```

MultiDiscrete spaces:

```json
{
  "type": "MultiDiscrete",
  "low": [0, 0, 0],
  "high": [4, 1, 1],
}
```

Tuple spaces:

```json
{
  "type": "Tuple",
  "subspaces": [
    {
      "type": "Discrete",
      "n": 5,
    },
    ...
  ],
}
```
