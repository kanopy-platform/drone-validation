// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/drone/drone-go/plugin/validator"
	"github.com/open-policy-agent/opa/rego"
	yaml "gopkg.in/yaml.v2"
)

func New(policy string) validator.Plugin {
	return &plugin{
		policyPath: policy,
	}
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {

	var documents []DroneConfig

	dec := yaml.NewDecoder(strings.NewReader(req.Config.Data))
	for {
		var document DroneConfig

		err := dec.Decode(&document)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		documents = append(documents, document)
	}

	r, err := rego.New(
		rego.Query("deny = data.drone.validation.deny; msg = data.drone.validation.out"),
		rego.Load([]string{p.policyPath}, nil)).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	for _, resource := range documents {
		rs, err := r.Eval(ctx, rego.EvalInput(resource))
		if err != nil {
			return err
		}
		if rs[0].Bindings["deny"] == true {
			message := fmt.Sprintf("pipeline %v", rs[0].Bindings["msg"])
			return errors.New(message)
		}

	}
	return nil
}
