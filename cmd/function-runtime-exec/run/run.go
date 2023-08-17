/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package run implements a convenience CLI to run and test Composition Functions.
package run

import (
	"context"
	"os"

	"github.com/crossplane/function-runtime-exec/internal/proto/v1beta1"
	"github.com/crossplane/function-runtime-exec/internal/server"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

// Command runs a Composition function.
type Command struct {
	Args               []string `arg:""`
	RunFunctionRequest []byte   `help:"YAML encoded RunFunctionRequest to pass to the function." type:"filecontent"`
}

// Run a Composition container function.
func (c *Command) Run(log logging.Logger) error {
	if len(c.Args) < 1 {
		return errors.Errorf("at least one argument is required")
	}
	runner := server.NewRunner(c.Args[0], c.Args[1:], server.WithLogger(log))

	var req v1beta1.RunFunctionRequest
	err := yaml.Unmarshal(c.RunFunctionRequest, &req)
	if err != nil {
		return errors.Wrap(err, "cannot read RunFunctionRequest")
	}
	resp, err := runner.RunFunction(context.Background(), &req)
	if err != nil {
		return errors.Wrap(err, "cannot run function")
	}

	b, err := yaml.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "cannot marshal response")
	}
	_, err = os.Stdout.Write(b)
	return errors.Wrap(err, "cannot write fio")
}
