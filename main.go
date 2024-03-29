// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"os"

	"github.com/drone/drone-go/plugin/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kanopy-platform/drone-validation/pkg/plugin"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind       string `envconfig:"DRONE_BIND"`
	Debug      bool   `envconfig:"DRONE_DEBUG"`
	Secret     string `envconfig:"DRONE_VALIDATE_PLUGIN_SECRET"`
	PolicyPath string `envconfig:"DRONE_POLICY_PATH"`
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	handler := validator.Handler(
		spec.Secret,
		plugin.New(plugin.WithPolicyPath(spec.PolicyPath)),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	http.HandleFunc("/healthz", healthz)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
