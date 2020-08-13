// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"

	"github.com/drone/drone-go/plugin/validator"
)

// New returns a new validator plugin.
func New(min, max int) validator.Plugin {
	return &plugin{
		min: min,
		max: max,
	}
}

type plugin struct {
	// TODO replace or remove these configuration
	// parameters. They are for demo purposes only.
	min int
	max int
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	// TODO replace or remove these checks
	// They are for demo purposes only.
	// fmt.Println(req.Config.Data)
	switch {
	case len(req.Config.Data) > p.max:
		return errors.New("validator: configuration exceeds the maximum length")
	case len(req.Config.Data) < p.min:
		return errors.New("validator: configuration does not meet minimum length requirement")
	}

	// a nil error indicates the configuration is valid.
	return nil
}
