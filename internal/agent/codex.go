package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Codex struct{}

func (Codex) Name() string             { return "codex" }
func (Codex) Binary() string           { return "codex" }
func (Codex) NativeConfigDir() string  { return "~/.codex" }
func (Codex) ContainerConfigDir() string { return "/home/claude/.codex" }
func (Codex) AsylumConfigDir() string  { return "~/.asylum/agents/codex" }

func (Codex) EnvVars() map[string]string {
	return map[string]string{
		"CODEX_HOME": "/home/claude/.codex",
	}
}

func (Codex) HasSession(projectPath string) bool {
	configDir, err := expandHome("~/.asylum/agents/codex")
	if err != nil {
		return false
	}
	// Codex stores sessions in a global date-organized tree with no per-project
	// metadata. Use a per-project marker to avoid resuming the wrong project.
	encoded := strings.ReplaceAll(projectPath, "/", "-")
	marker := filepath.Join(configDir, "projects", encoded, ".has_session")
	_, err = os.Stat(marker)
	return err == nil
}

func (Codex) WriteMarker(projectPath string) error {
	configDir, err := expandHome("~/.asylum/agents/codex")
	if err != nil {
		return err
	}
	encoded := strings.ReplaceAll(projectPath, "/", "-")
	dir := filepath.Join(configDir, "projects", encoded)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".has_session"), nil, 0644)
}

func (Codex) Command(resume bool, extraArgs []string) []string {
	if resume {
		if len(extraArgs) == 0 {
			return wrapZsh("codex resume --last --yolo")
		}
		fmt.Fprintln(os.Stderr, "! codex: resume skipped because extra args were provided")
	}
	parts := []string{"codex", "--yolo"}
	parts = append(parts, quoteArgs(extraArgs)...)
	return wrapZsh(strings.Join(parts, " "))
}
