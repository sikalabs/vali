package vali

import (
	"encoding/json"
	"fmt"
	"os"
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

func ReadConfig() (ValiConfig, error) {
	for _, path := range []string{
		"vali.local.json",
		"vali.json",
		"/vali.json",
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var config ValiConfig
		if err := json.Unmarshal(data, &config); err != nil {
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
