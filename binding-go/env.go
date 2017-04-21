package gym

import (
	"bufio"
	"encoding/binary"
	"net"

	"github.com/unixpickle/essentials"
)

var byteOrder = binary.LittleEndian

// Env is a handle on a Gym environment.
type Env interface {
	// Reset resets the environment.
	Reset() (observation interface{}, err error)

	// Step takes an action.
	Step(action interface{}) (obs interface{}, reward float64,
		done bool, info interface{}, err error)

	// ActionSpace gets the action space.
	ActionSpace() (*Space, error)

	// ObservationSpace gets the observation space.
	ObservationSpace() (*Space, error)

	// SampleAction samples from the action space.
	SampleAction() (interface{}, error)

	// Monitor sets the environment up to save results
	// to the given directory.
	Monitor(dir string, force, resume bool) error

	// Close stops and cleans up the environment.
	Close() error
}

type connEnv struct {
	Buf  *bufio.ReadWriter
	Conn net.Conn
}

// Make creates an Env by connecting to an API server and
// requesting the given environment.
func Make(host, envName string) (env Env, err error) {
	defer essentials.AddCtxTo("make environment", &err)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	if err := handshake(rw, envName); err != nil {
		conn.Close()
		return nil, err
	}

	return &connEnv{Buf: rw, Conn: conn}, nil
}

func (c *connEnv) Reset() (observation interface{}, err error) {
	panic("nyi")
}

func (c *connEnv) Step(action interface{}) (obs interface{}, reward float64,
	done bool, info interface{}, err error) {
	panic("nyi")
}

func (c *connEnv) ActionSpace() (*Space, error) {
	panic("nyi")
}

func (c *connEnv) ObservationSpace() (*Space, error) {
	panic("nyi")
}

func (c *connEnv) SampleAction() (interface{}, error) {
	panic("nyi")
}

func (c *connEnv) Monitor(dir string, force, resume bool) error {
	panic("nyi")
}

func (c *connEnv) Close() error {
	return c.Conn.Close()
}
