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

// New returns a new validator plugin.
func New(policy string) validator.Plugin {
	return &plugin{
		policyPath: policy,
	}
}

type plugin struct {
	policyPath string
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	// Drone sends its raw config in the request body
	type DroneConfig struct {
		Kind  string `yaml:"kind" json:"kind"`
		Type  string `yaml:"type" json:"type"`
		Name  string `yaml:"name" json:"name"`
		Steps []struct {
			Name      string   `yaml:"name" json:"name"`
			Image     string   `yaml:"image" json:"image"`
			Commands  []string `yaml:"commands" json:"commands"`
			DependsOn []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
		} `yaml:"steps" json:"steps"`
		Trigger struct {
			Branch []string `yaml:"branch" json:"branch"`
		} `yaml:"trigger" json:"trigger"`
	}

	var droneConfig DroneConfig

	err := yaml.Unmarshal([]byte(req.Config.Data), &droneConfig)
	if err != nil {
		return err
	}

	r := rego.New(
		rego.Query("deny = data.drone.validation.deny; msg = data.drone.validation.out"),
		rego.Load([]string{p.policyPath}, nil))

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return err
	}

	rs, err := query.Eval(ctx, rego.EvalInput(droneConfig))
	if err != nil {
		return err
	}

	if rs[0].Bindings["deny"] == true {
		message := fmt.Sprintf("[validator] %v", rs[0].Bindings["msg"])
	 	return errors.New(message)
	}

	if err != nil {
		return err
	}

	return nil
}
