package firstrun

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDetectCredentials(t *testing.T) {
	home := t.TempDir()

	// No files → nothing detected
	if got := detectCredentials(home); len(got) != 0 {
		t.Fatalf("expected 0 credentials, got %d", len(got))
	}

	// Create Maven settings
	mavenDir := filepath.Join(home, ".m2")
	os.MkdirAll(mavenDir, 0755)
	os.WriteFile(filepath.Join(mavenDir, "settings.xml"), []byte("<settings/>"), 0644)

	got := detectCredentials(home)
	if len(got) != 1 || got[0].Path != ".m2/settings.xml" {
		t.Fatalf("expected maven credential, got %v", got)
	}

}

func TestWriteConfig(t *testing.T) {
	home := t.TempDir()
	asylumDir := filepath.Join(home, ".asylum")

	creds := []credential{
		{Path: ".m2/settings.xml", Label: "Maven settings.xml"},
	}

	if err := writeConfig(asylumDir, creds); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(asylumDir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	var cfg configFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}

	if len(cfg.Volumes) != 1 {
		t.Fatalf("expected 1 volume, got %d", len(cfg.Volumes))
	}
	if cfg.Volumes[0] != "~/.m2/settings.xml:ro" {
		t.Errorf("unexpected volume[0]: %s", cfg.Volumes[0])
	}
}

func TestRunSkipsExistingUser(t *testing.T) {
	home := t.TempDir()
	// Simulate existing user: agents/ directory exists
	os.MkdirAll(filepath.Join(home, ".asylum", "agents"), 0755)
	// Add credentials that would normally trigger prompt
	os.MkdirAll(filepath.Join(home, ".m2"), 0755)
	os.WriteFile(filepath.Join(home, ".m2", "settings.xml"), []byte("<settings/>"), 0644)

	if err := Run(home); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".asylum", "config.yaml")); err == nil {
		t.Fatal("config.yaml should not exist for existing user")
	}
}

func TestRunNoCredentialsSkips(t *testing.T) {
	home := t.TempDir()

	if err := Run(home); err != nil {
		t.Fatal(err)
	}
	// No credentials → no config written, no prompt
	if _, err := os.Stat(filepath.Join(home, ".asylum", "config.yaml")); err == nil {
		t.Fatal("config.yaml should not exist when no credentials found")
	}
}
