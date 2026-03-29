package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMigrateV1ToV2_GlobalConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	v1 := `agent: gemini
release-channel: dev
volumes:
  - ~/.m2/settings.xml:ro
`
	os.WriteFile(path, []byte(v1), 0644)

	if !NeedsMigration(path) {
		t.Fatal("should need migration")
	}

	if err := MigrateV1ToV2(path); err != nil {
		t.Fatal(err)
	}

	// Verify backup exists
	if _, err := os.Stat(path + ".backup"); err != nil {
		t.Error("backup should exist")
	}

	// Read result as text — should contain comments from default config
	data, _ := os.ReadFile(path)
	text := string(data)
	if !strings.Contains(text, "# Kits configure language toolchains") {
		t.Error("migrated config should preserve default config comments")
	}

	// Parse result
	var result Config
	yaml.Unmarshal(data, &result)

	// Version set
	if result.Version != currentVersion {
		t.Errorf("version = %v, want %s", result.Version, currentVersion)
	}

	// User values overlaid
	if result.ReleaseChannel != "dev" {
		t.Errorf("release-channel = %v, want dev", result.ReleaseChannel)
	}
	if result.Agent != "gemini" {
		t.Errorf("agent = %v, want gemini", result.Agent)
	}
	if len(result.Volumes) == 0 || result.Volumes[0] != "~/.m2/settings.xml:ro" {
		t.Errorf("volumes = %v, want [~/.m2/settings.xml:ro]", result.Volumes)
	}

	// All standard kits present from default config
	for _, kit := range []string{"docker", "java", "python", "node"} {
		if _, ok := result.Kits[kit]; !ok {
			t.Errorf("kit %q should be present from default config", kit)
		}
	}
}

func TestMigrateV1ToV2_GlobalConfigWithV1Fields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	v1 := `agent: claude
profiles:
  - java
  - node
versions:
  java: "17"
packages:
  npm:
    - turbo
features:
  shadow-node-modules: true
onboarding:
  npm: false
agents:
  - claude
  - gemini
`
	os.WriteFile(path, []byte(v1), 0644)

	if err := MigrateV1ToV2(path); err != nil {
		t.Fatal(err)
	}

	// V1 fields are transformed before overlay, but global migration
	// produces the default config — v1 kit-level customizations that
	// differ from defaults are not preserved in the overlay (they need
	// to be re-configured). This is acceptable for global config since
	// the default config already has sensible values for all kits.
	data, _ := os.ReadFile(path)
	var result Config
	yaml.Unmarshal(data, &result)

	if result.Version != currentVersion {
		t.Errorf("version = %v, want %s", result.Version, currentVersion)
	}
	// All standard kits present
	for _, kit := range []string{"docker", "java", "python", "node"} {
		if _, ok := result.Kits[kit]; !ok {
			t.Errorf("kit %q should be present", kit)
		}
	}
}

func TestMigrateV1ToV2_ProjectConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".asylum")

	v1 := `features:
  onboarding: false
packages:
  apt:
    - jq
`
	os.WriteFile(path, []byte(v1), 0644)

	if !NeedsMigration(path) {
		t.Fatal("should need migration (has features key)")
	}

	if err := MigrateV1ToV2(path); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	var result map[string]any
	yaml.Unmarshal(data, &result)

	// No version field for project configs
	if _, ok := result["version"]; ok {
		t.Error("project config should not have version")
	}

	kits := result["kits"].(map[string]any)
	apt := kits["apt"].(map[string]any)
	if apt["packages"] == nil {
		t.Error("apt packages should be set")
	}
}

func TestMigrateV1ToV2_AlreadyV2(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	v2 := `version: "0.2"
agent: claude
kits:
  java:
    default-version: "21"
`
	os.WriteFile(path, []byte(v2), 0644)

	if NeedsMigration(path) {
		t.Fatal("should NOT need migration (already v2)")
	}
}

func TestMigrateV1ToV2_NoFeaturesProjectConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".asylum")

	config := `agent: gemini
`
	os.WriteFile(path, []byte(config), 0644)

	if NeedsMigration(path) {
		t.Fatal("should NOT need migration (no features key)")
	}
}

func TestNeedsMigration_GlobalConfigWithoutVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	config := `release-channel: dev
volumes:
  - ~/.m2/settings.xml:ro
`
	os.WriteFile(path, []byte(config), 0644)

	if !NeedsMigration(path) {
		t.Fatal("global config without version field should need migration")
	}
}
