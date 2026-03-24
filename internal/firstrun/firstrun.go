package firstrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/inventage-ai/asylum/internal/log"
	"gopkg.in/yaml.v3"
)

// credential is a host file that can be mounted into the container.
type credential struct {
	Path  string // relative to home, e.g. ".m2/settings.xml"
	Label string // display name
}

var credentials = []credential{
	{Path: ".m2/settings.xml", Label: "Maven settings.xml"},
}

// Run detects a first-run condition and prompts the user to mount
// package manager credentials. It writes ~/.asylum/config.yaml if
// the user accepts. Uses ~/.asylum/agents/ as the signal that asylum
// has been used before (created by EnsureAgentConfig, not the installer).
func Run(homeDir string) error {
	agentsDir := filepath.Join(homeDir, ".asylum", "agents")
	if _, err := os.Stat(agentsDir); err == nil {
		return nil // existing user
	}

	found := detectCredentials(homeDir)
	if len(found) == 0 {
		return nil
	}

	if accepted := prompt(found); accepted {
		asylumDir := filepath.Join(homeDir, ".asylum")
		if err := writeConfig(asylumDir, found); err != nil {
			return fmt.Errorf("write config: %w", err)
		}
		log.Success("wrote %s", filepath.Join(asylumDir, "config.yaml"))
	}
	return nil
}

func detectCredentials(homeDir string) []credential {
	var found []credential
	for _, c := range credentials {
		path := filepath.Join(homeDir, c.Path)
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			found = append(found, c)
		}
	}
	return found
}

func prompt(found []credential) bool {
	fmt.Println()
	log.Info("Package manager credentials found:")
	for _, c := range found {
		fmt.Printf("  - %s (~/%s)\n", c.Label, c.Path)
	}
	fmt.Print("Mount these into the sandbox (read-only)? [Y/n] ")
	var answer string
	fmt.Scanln(&answer)
	return !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "n")
}

type configFile struct {
	Volumes []string `yaml:"volumes"`
}

func writeConfig(asylumDir string, creds []credential) error {
	if err := os.MkdirAll(asylumDir, 0755); err != nil {
		return err
	}
	var volumes []string
	for _, c := range creds {
		volumes = append(volumes, "~/"+c.Path+":ro")
	}
	cfg := configFile{Volumes: volumes}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(asylumDir, "config.yaml"), data, 0644)
}
