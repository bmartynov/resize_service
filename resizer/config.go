package resizer

import (
	"fmt"
	"github.com/disintegration/imaging"
)

var availableFilters = map[string]imaging.ResampleFilter{
	"nearest_neighbor":   imaging.NearestNeighbor,
	"box":                imaging.Box,
	"linear":             imaging.Linear,
	"hermite":            imaging.Hermite,
	"mitchell_netravali": imaging.MitchellNetravali,
	"catmull_rom":        imaging.CatmullRom,
	"b_spline":           imaging.BSpline,
	"gaussian":           imaging.Gaussian,
	"bartlett":           imaging.Bartlett,
	"lanczos":            imaging.Lanczos,
	"hann":               imaging.Hann,
	"hamming":            imaging.Hamming,
	"blackman":           imaging.Blackman,
	"welch":              imaging.Welch,
	"cosine":             imaging.Cosine,
}

type Config struct {
	Filter string `envconfig:"FILTER" required:"true"`
}

func (c *Config) getFilter() imaging.ResampleFilter {
	return availableFilters[c.Filter]
}

func (c *Config) Validate() error {
	if _, ok := availableFilters[c.Filter]; !ok {
		return fmt.Errorf("filter %s unknown", c.Filter)
	}
	return nil
}
