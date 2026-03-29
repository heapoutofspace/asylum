package config

import (
	"os"

	"github.com/inventage-ai/asylum/internal/kit"
)

const configHeader = `version: "0.2"

# Release channel for self-update (stable, dev)
release-channel: stable

# Agent to start by default (claude, gemini, codex, opencode)
agent: claude

# Agent CLIs to install in the container image.
# Remove or comment out agents you don't use to speed up image builds.
agents:
  claude:
  # gemini:
  # codex:
  # opencode:

# Kits configure language toolchains and tools installed in the container.
# A kit is active when its key is present (even with no options).
# Comment out or remove a kit to disable it entirely.
kits:`

const configFooter = `
# Port forwarding (host:container or just port for same on both sides)
# ports:
#   - "3000"
#   - "8080:80"

# Additional volume mounts
# Supports: /path, /host:/container, /host:/container:ro, ~/path
# volumes:
#   - ~/shared-data:/data
#   - ~/.aws

# Environment variables
# env:
#   GITHUB_TOKEN: ghp_xxx
#   NODE_ENV: development
`

// DefaultConfig returns the full default config assembled from the header,
// kit ConfigSnippets, and footer.
func DefaultConfig() string {
	return configHeader + kit.AssembleConfigSnippets() + configFooter
}

// WriteDefaults writes the default config to the given path if it doesn't
// already exist. It uses O_CREATE|O_EXCL to avoid a TOCTOU race.
func WriteDefaults(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	_, err = f.WriteString(DefaultConfig())
	return err
}
