package config

import (
	"net"
	"strconv"

	"github.com/pkg/errors"

	"github.com/bmartynov/resizer_service/crawler"
	"github.com/bmartynov/resizer_service/resizer"
)

func validateAddress(address string) error {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("port invalid")
	}

	if p <= 0 || p > 65536 {
		return errors.New("port invalid")
	}
	return nil
}

type HttpTransportConfig struct {
	Address string `envconfig:"ADDRESS" required:"true"`
}

func (c *HttpTransportConfig) Validate() error {
	if err := validateAddress(c.Address); err != nil {
		return err
	}
	return nil
}

type ResizerConfig struct {
	Http    HttpTransportConfig `envconfig:"HTTP" required:"true"`
	Crawler crawler.Config      `envconfig:"CRAWLER" required:"true"`
	Resizer resizer.Config      `envconfig:"RESIZER" required:"true"`
}

func (c *ResizerConfig) Validate() error {
	if err := c.Http.Validate(); err != nil {
		return errors.Wrap(err, "http.validate")
	}
	if err := c.Crawler.Validate(); err != nil {
		return errors.Wrap(err, "crawler.validate")
	}
	if err := c.Resizer.Validate(); err != nil {
		return errors.Wrap(err, "resizer.validate")
	}
	return nil
}

func NewResizerConfig() *ResizerConfig {
	return &ResizerConfig{
		Http:    HttpTransportConfig{},
		Crawler: crawler.Config{},
		Resizer: resizer.Config{},
	}
}
