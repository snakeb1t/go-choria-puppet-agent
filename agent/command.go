package agent

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	PuppetCommand = "/opt/puppetlabs/puppet/bin/puppet"
	PuppetAgentSubcommand = "agent"
	EnvironmentFile = "/etc/environment"
)

func runAgentOnce(req *RunonceRequest) (*RunonceResponse, error) {
	args := buildArgs(req)
	cmd := exec.Command(PuppetCommand, args...)
	var outErrBuf bytes.Buffer
	cmd.Stdout = &outErrBuf
	cmd.Stderr = &outErrBuf
	envVars, err := parseEnvironmentFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get env vars")
	}
	cmd.Env = envVars
	start := time.Now()
	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute puppet agent")
	}
	return &RunonceResponse{
		Summary:     outErrBuf.String(),
		InitiatedAt: start.String(),
	}, nil
}

func buildArgs(req *RunonceRequest) []string {
	args := []string{PuppetAgentSubcommand, "-t"}
	if req.Noop {
		args = append(args, "--noop")
	}
	if req.Environment != "" {
		args = append(args, []string{"--environment", req.Environment}...)
	}
	return args
}

func parseEnvironmentFile() ([]string, error) {
	envVars := []string{}
	fh, err := os.Open(EnvironmentFile)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		if strings.ContainsAny(scanner.Text(), "#") {
			continue
		}
		envVars = append(envVars, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return envVars, nil
}