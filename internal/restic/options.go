package restic

import (
	"fmt"
	"log"
	"strings"
)

// Names of restic environment variables rsched needs to keep track of.
const (
	EnvResticRepository      = "RESTIC_REPOSITORY"
	EnvResticRepositoryFile  = "RESTIC_REPOSITORY_FILE"
	EnvResticPassword        = "RESTIC_PASSWORD"
	EnvResticPasswordFile    = "RESTIC_PASSWORD_FILE"
	EnvResticPasswordCommand = "RESTIC_PASSWORD_COMMAND"
)

var requiredEnvVars = [][]string{
	{EnvResticRepository, EnvResticRepositoryFile},
	{EnvResticPassword, EnvResticPasswordFile, EnvResticPasswordCommand},
}

// Option is the type of any option accepted by the restic packages public
// functions.
type Option func(*options)

type options struct {
	Restic string
	Runner CmdRunner
	Env    map[string]string
}

func (o *options) Apply(opts []Option) error {
	// Set defaults. May be overwritten further down the line.
	o.Restic = "restic"
	o.Runner = &osRunner{}

	for _, opt := range opts {
		opt(o)
	}

	return o.Validate()
}

func (o *options) Validate() error {
	for _, alternatives := range requiredEnvVars {
		var ok bool

		for _, alt := range alternatives {
			ok = ok || o.Env[alt] != ""
		}
		if !ok {
			return fmt.Errorf("environment requires one of: %s", strings.Join(alternatives, ", "))
		}
	}
	return nil
}

// WithEnv adds the passed key value pairs to the environment that is used
// to call restic.
func WithEnv(env map[string]string) Option {
	return func(o *options) {
		if o.Env == nil {
			o.Env = make(map[string]string)
		}
		for k, v := range env {
			if _, ok := o.Env[k]; ok {
				log.Printf("Replacing already existing value for %q in restic environment", k)
			}
			o.Env[k] = v
		}
	}
}

// WithPassword adds the RESTIC_PASSWORD environment variable to the restic
// environment.
func WithPassword(pw string) Option {
	return WithEnv(map[string]string{EnvResticPassword: pw})
}

// WithRepository adds the RESTIC_REPOSITORY environment variable to the
// restic environment.
func WithRepository(repo string) Option {
	return WithEnv(map[string]string{EnvResticRepository: repo})
}

// WithCmdRunner allows to use a specialized command runner.
//
// This is intended for testing purposes as it allows to test calls to restic
// without actually calling the executable.
func WithCmdRunner(r CmdRunner) Option {
	return func(opts *options) {
		opts.Runner = r
	}
}
