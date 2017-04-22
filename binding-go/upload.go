package gym

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/unixpickle/essentials"
)

// Upload connects to the API host and uses it to upload a
// monitor directory to the Gym website.
//
// If the API key is "", the OPENAI_GYM_API_KEY
// environment variable is used.
//
// If the directory is a relative path, it should be
// relative to the current working directory.
func Upload(apiHost, dir, apiKey, algorithmID string) (err error) {
	essentials.AddCtxTo("upload monitor", &err)
	env, err := Make(apiHost, "")
	if err != nil {
		return err
	}
	c := env.(*connEnv)
	defer c.Close()

	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_GYM_API_KEY")
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if err := writePacketType(c.Buf, packetUpload); err != nil {
		return err
	}
	for _, str := range []string{absDir, apiKey, algorithmID} {
		if err := writeByteField(c.Buf, []byte(str)); err != nil {
			return err
		}
	}
	if err := c.Buf.Flush(); err != nil {
		return err
	}

	if errMsg, err := readByteField(c.Buf); err != nil {
		return err
	} else if len(errMsg) > 0 {
		return errors.New(string(errMsg))
	}

	return nil
}
