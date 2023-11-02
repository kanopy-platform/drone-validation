// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/validator"
)

var noContext = context.Background()

var req = &validator.Request{
	Build: drone.Build{
		After:  "3d21ec53a331a6f037a91c368710b99387d012c1",
		Parent: 1,
		Deploy: "production",
		Sender: "test",
		Event:  "promote",
	},
	Repo: drone.Repo{
		Slug:   "repo/test",
		Config: ".drone.yml",
		Branch: "main",
	},
	Config: drone.Config{},
}

func getSamplePipeline(sample string) (string, error) {
	path := filepath.Join("testdata", sample)
	sampleData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(sampleData), err
}

func checkOutput(plugin validator.Plugin, sampleFile, expected string) func(*testing.T) {
	return func(t *testing.T) {
		pipeline, err := getSamplePipeline(sampleFile)
		if err != nil {
			t.Error(err)
			return
		}
		req.Config.Data = pipeline

		validate := plugin.Validate(noContext, req)
		if validate == nil && expected == "" {
			return
		}
		if validate.Error() != expected {
			t.Errorf("Invalid evaluation output, returned: %v, expected: %v.", validate, expected)
		}
	}
}

func TestPlugin(t *testing.T) {
	plugin := New("../policy")
	pipeline, err := getSamplePipeline("authorized-type.yml")
	if err != nil {
		t.Error(err)
		return
	}
	req.Config.Data = pipeline

	err = plugin.Validate(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("authorized pipeline type", checkOutput(plugin, "authorized-type.yml", ""))
	t.Run("unauthorized pipeline type", checkOutput(plugin, "unauthorized-type.yml", "type 'k8s' is not supported"))
	t.Run("empty pipeline type", checkOutput(plugin, "default-type.yml", "type 'docker' is not supported"))
	t.Run("broken pipeline config file", checkOutput(plugin, "broken-config.yml", "yaml: line 2: mapping values are not allowed in this context"))
}

func TestValidateInvalidPolicy(t *testing.T) {
	plugin := New("testdata/empty.rego")
	pipeline, err := getSamplePipeline("authorized-type.yml")
	if err != nil {
		t.Error(err)
		return
	}
	req.Config.Data = pipeline
	err = plugin.Validate(noContext, req)
	if err == nil {
		t.Error(err)
		return
	}
}

func TestValidateDefaultPolicy(t *testing.T) {
	want := &plugin{
		policyPath: "./policy",
	}
	got := New("")

	if !reflect.DeepEqual(got, want) {
		t.Error("invalid default policyPath")
		return

	}
}
