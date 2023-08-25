/*
Copyright 2019 The Crossplane Authors.

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

// Package server is the server for the function runtime exec.
package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"os/exec"
	"path/filepath"

	"github.com/crossplane/function-runtime-exec/internal/proto/v1beta1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/crossplane/crossplane-runtime/pkg/certificates"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
)

// Runner runs a Composition Function as a subprocess.
type Runner struct {
	// UnimplementedFunctionRunnerServiceServer is an unimplemented server API
	// for FunctionRunnerService service.
	v1beta1.UnimplementedFunctionRunnerServiceServer
	command string
	args    []string

	log       logging.Logger
	certsPath string
}

// RunnerOpts are options for a Runner.
type RunnerOpts func(*Runner)

// WithLogger configures the supplied Runner to use the supplied logger.
func WithLogger(l logging.Logger) RunnerOpts {
	return func(r *Runner) {
		r.log = l
	}
}

// WithServerTLSCertPath configures the supplied Runner to use the supplied
// server TLS certificate path.
func WithServerTLSCertPath(p string) RunnerOpts {
	return func(r *Runner) {
		r.certsPath = p
	}
}

// NewRunner returns a new Runner executing the supplied command with the
// supplied arguments and configured by the supplied options.
func NewRunner(command string, args []string, opts ...RunnerOpts) *Runner {
	r := &Runner{
		command: command,
		args:    args,
		log:     logging.NewNopLogger(),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// RunFunction implements the FunctionRunnerServiceServer interface.
func (r *Runner) RunFunction(ctx context.Context, req *v1beta1.RunFunctionRequest) (*v1beta1.RunFunctionResponse, error) {
	r.log.Debug("Running", "command", r.command, "args", r.args)
	cmd := exec.CommandContext(ctx, r.command, r.args...) //nolint:gosec // We want to run arbitrary commands.
	b, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}
	cmd.Stdin = bytes.NewReader(b)
	stdOut, err := cmd.StdoutPipe()
	defer func() {
		_ = stdOut.Close()
	}()
	if err != nil {
		return nil, err
	}
	stdErr, err := cmd.StderrPipe()
	defer func() {
		_ = stdErr.Close()
	}()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	stdOutBytes, err := io.ReadAll(stdOut)
	if err != nil {
		return nil, err
	}
	stdErrBytes, err := io.ReadAll(stdErr)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			r.log.Debug("Ran", "command", r.command, "args", r.args, "stdout", string(stdOutBytes), "stderr", string(stdErrBytes), "exitCode", exitErr.ExitCode())
		}
		return nil, err
	}
	r.log.Debug("Ran", "command", r.command, "args", r.args, "stdout", string(stdOutBytes), "stderr", string(stdErrBytes))
	res := &v1beta1.RunFunctionResponse{}
	if err := json.Unmarshal(stdOutBytes, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListenAndServe gRPC connections at the supplied address.
func (r *Runner) ListenAndServe(network, address string) error {
	r.log.Debug("Listening", "network", network, "address", address, "command", r.command, "args", r.args, "certsPath", r.certsPath)
	lis, err := net.Listen(network, address)
	if err != nil {
		return errors.Wrapf(err, "while trying to listen on network: %s, address: %s", network, address)
	}

	var opts []grpc.ServerOption
	if r.certsPath != "" {
		tlsConfig, err := certificates.LoadMTLSConfig(filepath.Join(r.certsPath, "ca.crt"), filepath.Join(r.certsPath, "tls.crt"), filepath.Join(r.certsPath, "tls.key"), true)
		if err != nil {
			return errors.Wrap(err, "while loading mTLS config")
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}
	// TODO(negz): Limit concurrent function runs?
	srv := grpc.NewServer(opts...)
	reflection.Register(srv)
	v1beta1.RegisterFunctionRunnerServiceServer(srv, r)
	return errors.Wrap(srv.Serve(lis), "while running grpc server")
}
