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

// Package start implements the reference Composition Function runner.
// It exposes a gRPC API that may be used to run Composition Functions.
package start

import (
	"github.com/crossplane/function-runtime-exec/internal/server"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

// Command starts a gRPC API to run Composition Functions.
type Command struct {
	Args []string `arg:""`

	Network           string `help:"Network on which to listen for gRPC connections." default:"tcp"`
	Address           string `help:"Address at which to listen for gRPC connections." default:"0.0.0.0:1234"`
	TLSServerCertsDir string `help:"Folder containing server certs (tls.key, tls.crt) and the CA used to verify client certificates (ca.crt)" env:"TLS_SERVER_CERTS_DIR"`
}

// Run a Composition Function gRPC API.
func (c *Command) Run(log logging.Logger) error {
	// TODO(negz): Expose a healthz endpoint and otel metrics.
	if len(c.Args) < 1 {
		return errors.Errorf("at least one argument is required")
	}
	f := server.NewRunner(c.Args[0], c.Args[1:], server.WithLogger(log), server.WithServerTLSCertPath(c.TLSServerCertsDir))
	return errors.Wrap(f.ListenAndServe(c.Network, c.Address), "cannot listen for and serve gRPC API")
}
