# Protocol

All integers are encoded in little endian. All strings are UTF-8. Whenever there's an error field, it can be an empty string to indicate success.

## Initial connection

During this stage, the client initiates a connection and requests an environment. The server attempts to create the environment, or fails with an error (e.g. if the environment does not exist).

|Source   |Type    | Description           |
|---------|--------|-----------------------|
|Client   |uint8   | Flags (all 0)         |
|Client   |uint32  | Length of env name    |
|Client   |string  | Environment name      |
|Server   |uint32  | Error length          |
|Server   |string  | Error message         |
