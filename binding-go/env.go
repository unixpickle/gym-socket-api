package gym

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"

	"github.com/unixpickle/essentials"
)

// Env is a handle on a Gym environment.
//
// The methods on an Env are thread-safe.
type Env interface {
	// Reset resets the environment.
	Reset() (obs Obs, err error)

	// Step takes an action.
	Step(action interface{}) (obs Obs, reward float64,
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

	CmdLock sync.Mutex
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

func (c *connEnv) Reset() (obs Obs, err error) {
	defer essentials.AddCtxTo("reset environment", &err)
	c.CmdLock.Lock()
	defer c.CmdLock.Unlock()
	if err := writePacketType(c.Buf, packetReset); err != nil {
		return nil, err
	}
	if err := c.Buf.Flush(); err != nil {
		return nil, err
	}
	return readObservation(c.Buf)
}

func (c *connEnv) Step(action interface{}) (obs Obs, reward float64,
	done bool, info interface{}, err error) {
	defer essentials.AddCtxTo("step environment", &err)
	c.CmdLock.Lock()
	defer c.CmdLock.Unlock()
	err = writePacketType(c.Buf, packetStep)
	if err != nil {
		return
	}
	err = writeAction(c.Buf, action)
	if err != nil {
		return
	}
	err = c.Buf.Flush()
	if err != nil {
		return
	}
	obs, err = readObservation(c.Buf)
	if err != nil {
		return
	}
	reward, err = readReward(c.Buf)
	if err != nil {
		return
	}
	done, err = readBool(c.Buf)
	if err != nil {
		return
	}
	infoData, err := readByteField(c.Buf)
	if err != nil {
		return
	}
	err = json.Unmarshal(infoData, &info)
	return
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
