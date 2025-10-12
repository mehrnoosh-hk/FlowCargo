package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"flowcargo/internal/shared/config"
)

type Flags struct {
	Environment config.Environment
	ConfigPath  *string
}

func ParseFlags() (Flags, error) {
	return ParseFlagsFromSet(flag.CommandLine, os.Args[1:])
}

func ParseFlagsFromSet(f *flag.FlagSet, args []string) (Flags, error) {
	env := f.String("env", "", "Environment (required): development, test, or production")
	configPath := f.String("config", "", "Path to config file (optional)")

	f.Parse(args)

	parsedFlags, err := validateFlags(env, configPath)

	if err != nil {
		return Flags{}, err
	}
	return parsedFlags, nil
}

func validateFlags(env, configPath *string) (Flags, error) {
	environment, err := validateEnvironment(env)
	if err != nil {
		return Flags{}, err
	}
	path, err := validatePath(configPath)
	if err != nil {
		return Flags{}, err
	}
	return Flags{
		Environment: environment,
		ConfigPath:  path,
	}, nil
}

func validateEnvironment(env *string) (config.Environment, error) {
	if env == nil || *env == "" {
		return config.Environment(""), errors.New("environment flag is required")
	}
	switch config.Environment(*env) {
	case config.Development, config.Test, config.Production:
		return config.Environment(*env), nil
	default:
		return config.Environment(""), fmt.Errorf("invalid environment '%s'. Valid options: development, test, production", *env)
	}
}

func validatePath(path *string) (*string, error) {
	if path == nil || *path == "" {
		return nil, nil
	}
	// check if the file exists
	if _, err := os.Stat(*path); err != nil {
		return nil, fmt.Errorf("cannot stat config file %s: %w", *path, err)
	}
	return path, nil
}
