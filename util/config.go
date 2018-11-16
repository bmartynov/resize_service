package util

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/bmartynov/resizer_service/config"
	"github.com/bmartynov/resizer_service/crawler"
	"github.com/bmartynov/resizer_service/resizer"
)

func NewResizerConfig() (
	creawlerConf *crawler.Config,
	resizeConf *resizer.Config,
	transport *config.HttpTransportConfig,
	err error,
) {
	c := config.NewResizerConfig()

	err = envconfig.Process("RESIZER_SERVICE", c)
	if err != nil {
		return
	}

	err = c.Validate()
	if err != nil {
		return
	}

	transport = &c.Http
	creawlerConf = &c.Crawler
	resizeConf = &c.Resizer

	return
}
