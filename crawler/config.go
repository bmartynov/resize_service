package crawler

import (
	"errors"
	"time"
)

type Config struct {
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func (c *Config) Validate() error {
	if c.Timeout.Seconds() < 1 {
		return errors.New("crawler timeout too small")
	}
	return nil
}
