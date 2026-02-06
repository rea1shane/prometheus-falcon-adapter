package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Global   Global   `yaml:"global"`
	Falcon   Falcon   `yaml:"falcon"`
	Exporter Exporter `yaml:"exporter"`
}

type Global struct {
	ScrapeInterval time.Duration `yaml:"scrape_interval"`
	ScrapeTimeout  time.Duration `yaml:"scrape_timeout"`
}

type Falcon struct {
	Target string `yaml:"target"`
}

type Exporter struct {
	Targets []string `yaml:"targets"`
}

func (c Config) Content() []any {
	return []any{
		"global.scrape_interval",
		c.Global.ScrapeInterval.String(),
		"global.scrape_timeout",
		c.Global.ScrapeTimeout.String(),
		"falcon.target",
		c.Falcon.Target,
		"exporter.targets",
		strings.Join(c.Exporter.Targets, ", "),
	}
}

func New(file string) (Config, error) {
	var cfg Config

	// Read file
	data, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal YAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// Verify global configuration
	if cfg.Global.ScrapeInterval <= 0 {
		return Config{}, fmt.Errorf("non-positive global.scrape_interval")
	}
	if cfg.Global.ScrapeTimeout <= 0 {
		return Config{}, fmt.Errorf("non-positive global.scrape_timeout")
	}
	if cfg.Global.ScrapeTimeout >= cfg.Global.ScrapeInterval {
		return Config{}, fmt.Errorf("global.scrape_timeout not less than global.scrape_interval")
	}

	// Verify targets
	if err := verifyTarget(cfg.Falcon.Target); err != nil {
		return Config{}, fmt.Errorf("invalid falcon.target: %w", err)
	}

	for _, target := range cfg.Exporter.Targets {
		if err := verifyTarget(target); err != nil {
			return Config{}, fmt.Errorf("invalid exporter.targets %s: %w", target, err)
		}
	}

	return cfg, nil
}

func verifyTarget(in string) error {
	if strings.TrimSpace(in) == "" {
		return fmt.Errorf("empty target")
	}

	if !strings.HasPrefix(in, "http://") && !strings.HasPrefix(in, "https://") {
		return fmt.Errorf("target not start with http:// or https://")
	}

	if _, err := ExtractHostname(in); err != nil {
		return fmt.Errorf("failed to extract hostname: %w", err)
	}

	return nil
}

func ExtractHostname(in string) (string, error) {
	parsed, err := url.Parse(in)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("empty hostname")
	}

	return hostname, nil
}
