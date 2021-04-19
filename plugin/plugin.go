// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

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
	f := bytes.NewBufferString(req.Config.Data)
	resources, err := Parse(f)
    if err != nil {
		return err
	}
	for _, resource := range resources {


	r := rego.New(
		rego.Query("deny = data.drone.validation.deny; msg = data.drone.validation.out"),
		rego.Load([]string{p.policyPath}, nil),
		rego.Input(resource))

	rs, err := r.Eval(ctx)
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

func Parse(r io.Reader) ([]*DroneConfig, error) {
	const newline = '\n'
	var resources []*DroneConfig
	var resource *DroneConfig

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isSeparator(line) {
			resource = nil
		}
		if resource == nil {
			resource = &DroneConfig{}
			resources = append(resources, resource)
		}
		if isSeparator(line) {
			continue
		}
		if isTerminator(line) {
			break
		}
		if scanner.Err() == io.EOF {
			break
		}
		resource.Data = append(
			resource.Data,
			line...,
		)
		resource.Data = append(
			resource.Data,
			newline,
		)
	}
	for _, resource := range resources {
		err := yaml.Unmarshal(resource.Data, resource)
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

func isSeparator(s string) bool {
	return strings.HasPrefix(s, "---")
}

func isTerminator(s string) bool {
	return strings.HasPrefix(s, "...")
}
