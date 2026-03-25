package config

import "os"

const defaultConfig = `version: "0.2"

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
kits:
  docker:               # Docker-in-Docker support (remove to disable)

  java:
    versions:
      - 17
      - 21
      - 25
    default-version: 21

  python:
    # versions:
    #   - 3.14
    # packages:          # Python tools installed via uv
    #   - ansible

  node:
    shadow-node-modules: true
    onboarding: false
    # versions:
    #   - 24
    # packages:          # npm packages installed globally
    #   - turbo

  # apt:                 # System packages installed via apt-get
  #   packages:
  #     - imagemagick
  #     - ffmpeg

  # title:
  #   # Terminal tab title template
  #   # Placeholders: {project}, {agent}, {mode}
  #   tab-title: "🤖 {project}"
  #   allow-agent-terminal-title: false

  # shell:
  #   build:             # Custom commands run at image build time
  #     - "curl -fsSL https://example.com/install.sh | sh"
  #   start:             # Custom commands run at container start
  #     - 'echo "WELCOME TO THE ASYLUM!"'

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

// WriteDefaults writes the default config to the given path if it doesn't exist.
func WriteDefaults(path string) error {
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}
