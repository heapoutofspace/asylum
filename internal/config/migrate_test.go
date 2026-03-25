package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMigrateV1ToV2_GlobalConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	v1 := `agent: claude
profiles:
  - java
  - node
versions:
  java: "17"
packages:
  apt:
    - ffmpeg
  npm:
    - turbo
  pip:
    - ansible
  run:
    - "curl https://example.com | sh"
features:
  shadow-node-modules: true
  allow-agent-terminal-title: false
onboarding:
  npm: false
tab-title: "🤖 {project}"
agents:
  - claude
  - gemini
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

	// Parse result
	data, _ := os.ReadFile(path)
	var result map[string]any
	yaml.Unmarshal(data, &result)

	// Version set
	if result["version"] != currentVersion {
		t.Errorf("version = %v, want %s", result["version"], currentVersion)
	}

	// Old keys removed
	for _, key := range []string{"profiles", "versions", "packages", "features", "onboarding", "tab-title"} {
		if _, ok := result[key]; ok {
			t.Errorf("old key %q should be removed", key)
		}
	}

	// Kits created
	kits, ok := result["kits"].(map[string]any)
	if !ok {
		t.Fatal("kits should be a map")
	}

	// Docker kit added (was always on in v1)
	if _, ok := kits["docker"]; !ok {
		t.Error("docker kit should be added during migration")
	}

	// Java kit
	java, ok := kits["java"].(map[string]any)
	if !ok {
		t.Fatal("java kit should exist")
	}
	if java["default-version"] != "17" {
		t.Errorf("java default-version = %v, want 17", java["default-version"])
	}

	// Node kit
	node, ok := kits["node"].(map[string]any)
	if !ok {
		t.Fatal("node kit should exist")
	}
	if node["shadow-node-modules"] != true {
		t.Error("node shadow-node-modules should be true")
	}
	if node["onboarding"] != false {
		t.Error("node onboarding should be false")
	}

	// Title kit
	title, ok := kits["title"].(map[string]any)
	if !ok {
		t.Fatal("title kit should exist")
	}
	if title["tab-title"] != "🤖 {project}" {
		t.Errorf("title tab-title = %v", title["tab-title"])
	}

	// Agents converted from list to map
	agents, ok := result["agents"].(map[string]any)
	if !ok {
		t.Fatal("agents should be a map")
	}
	if _, ok := agents["claude"]; !ok {
		t.Error("claude should be in agents")
	}
	if _, ok := agents["gemini"]; !ok {
		t.Error("gemini should be in agents")
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
