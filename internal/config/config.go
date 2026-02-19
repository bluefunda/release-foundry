package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RepoEntry defines a single repository in a batch config.
type RepoEntry struct {
	Owner         string   `yaml:"owner"`
	Repo          string   `yaml:"repo"`
	Edition       string   `yaml:"edition"`
	BaseBranch    string   `yaml:"baseBranch"`
	IncludeLabels []string `yaml:"includeLabels"`
	ExcludeLabels []string `yaml:"excludeLabels"`
}

// ReposConfig is the top-level YAML structure for multi-repo batch mode.
type ReposConfig struct {
	Defaults struct {
		Owner      string `yaml:"owner"`
		BaseBranch string `yaml:"baseBranch"`
	} `yaml:"defaults"`
	Repos []RepoEntry `yaml:"repos"`
}

// LoadReposConfig reads and parses a YAML config file, applying defaults.
func LoadReposConfig(path string) (*ReposConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg ReposConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if len(cfg.Repos) == 0 {
		return nil, fmt.Errorf("config %s: no repos defined", path)
	}

	// Apply defaults to each repo entry.
	for i := range cfg.Repos {
		if cfg.Repos[i].Owner == "" {
			cfg.Repos[i].Owner = cfg.Defaults.Owner
		}
		if cfg.Repos[i].BaseBranch == "" {
			cfg.Repos[i].BaseBranch = cfg.Defaults.BaseBranch
		}
		if cfg.Repos[i].BaseBranch == "" {
			cfg.Repos[i].BaseBranch = "main"
		}
		if cfg.Repos[i].Owner == "" {
			return nil, fmt.Errorf("config %s: repo[%d] (%s) has no owner", path, i, cfg.Repos[i].Repo)
		}
		if cfg.Repos[i].Repo == "" {
			return nil, fmt.Errorf("config %s: repo[%d] has no repo name", path, i)
		}
	}

	return &cfg, nil
}
