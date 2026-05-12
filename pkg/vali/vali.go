package vali

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ValiConfig struct {
	Meta struct {
		Version int
	}
	Env   []string
	Files []string
}

type ValiOutput struct {
	OK           bool
	MissingEnv   []string
	MissingFiles []string
}

func ReadConfigFromFile(path string) (ValiConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ValiConfig{}, fmt.Errorf("cannot read config file %s: %w", path, err)
	}
	var config ValiConfig
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return ValiConfig{}, fmt.Errorf("cannot parse config file %s: %w", path, err)
		}
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			return ValiConfig{}, fmt.Errorf("cannot parse config file %s: %w", path, err)
		}
	}
	return config, nil
}

func ReadConfig() (ValiConfig, error) {
	for _, path := range []string{
		"vali.local.yaml",
		"vali.local.yml",
		"vali.local.json",
		"vali.yaml",
		"vali.yml",
		"vali.json",
		"/vali.yaml",
		"/vali.yml",
		"/vali.json",
	} {
		config, err := ReadConfigFromFile(path)
		if err != nil {
			continue
		}
		return config, nil
	}
	return ValiConfig{}, fmt.Errorf("no vali config file found")
}

func Validate(config ValiConfig) ValiOutput {
	out := ValiOutput{OK: true}

	// Evironment Variable Exists Validation
	for _, name := range config.Env {
		if os.Getenv(name) == "" {
			out.MissingEnv = append(out.MissingEnv, name)
			out.OK = false
		}
	}

	// File Exists Validation
	for _, path := range config.Files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			out.MissingFiles = append(out.MissingFiles, path)
			out.OK = false
		}
	}

	return out
}

func PrintValidateOutput(out ValiOutput) {
	if !out.OK {
		fmt.Fprintln(os.Stderr, "Validation failed")
		for _, name := range out.MissingEnv {
			fmt.Fprintf(os.Stderr, "  missing env:  %s\n", name)
		}
		for _, path := range out.MissingFiles {
			fmt.Fprintf(os.Stderr, "  missing file: %s\n", path)
		}
		return
	}
	fmt.Println("OK")
}
