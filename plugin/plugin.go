// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/drone/drone-go/plugin/validator"
	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v2"
)

func New(policy string) validator.Plugin {
	return &plugin{
		policyPath: policy,
	}
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {

	var droneConfig DroneConfig

	err := yaml.Unmarshal([]byte(req.Config.Data), &droneConfig)
	if err != nil {
		return err
	}

	r := rego.New(
		rego.Query("deny = data.drone.validation.deny; msg = data.drone.validation.out"),
		rego.Load([]string{p.policyPath}, nil),
		rego.Input(droneConfig))

	rs, err := r.Eval(ctx)
	if err != nil {
		return err
	}
	if rs[0].Bindings["deny"] == true {
		message := fmt.Sprintf("pipeline %v", rs[0].Bindings["msg"])
		return errors.New(message)
	}

	return nil
}
