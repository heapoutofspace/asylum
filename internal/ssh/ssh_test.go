package ssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInit_KeyAlreadyExists verifies Init returns nil and skips keygen when
// the key file already exists.
func TestInit_KeyAlreadyExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	sshDir := filepath.Join(home, ".asylum", "ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(sshDir, "id_ed25519")
	if err := os.WriteFile(keyPath, []byte("existing key"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init returned error when key exists: %v", err)
	}
}

// TestInit_CopiesKnownHosts verifies that an existing ~/.ssh/known_hosts is
// copied into the asylum ssh directory.
func TestInit_CopiesKnownHosts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Create ~/.ssh/known_hosts
	dotSSH := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(dotSSH, 0700); err != nil {
		t.Fatal(err)
	}
	knownHostsContent := []byte("github.com ssh-ed25519 AAAA...")
	if err := os.WriteFile(filepath.Join(dotSSH, "known_hosts"), knownHostsContent, 0600); err != nil {
		t.Fatal(err)
	}

	// Also pre-create the key so Init doesn't invoke ssh-keygen.
	sshDir := filepath.Join(home, ".asylum", "ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "id_ed25519"), []byte("existing key"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	dst := filepath.Join(sshDir, "known_hosts")
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("known_hosts not copied: %v", err)
	}
	if strings.TrimSpace(string(got)) != strings.TrimSpace(string(knownHostsContent)) {
		t.Errorf("known_hosts content = %q, want %q", got, knownHostsContent)
	}
}

// TestInit_MergesKnownHosts verifies that existing container-specific entries
// in ~/.asylum/ssh/known_hosts are preserved when merging host known_hosts.
func TestInit_MergesKnownHosts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Create ~/.ssh/known_hosts with one entry.
	dotSSH := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(dotSSH, 0700); err != nil {
		t.Fatal(err)
	}
	hostEntry := "github.com ssh-ed25519 AAAA..."
	if err := os.WriteFile(filepath.Join(dotSSH, "known_hosts"), []byte(hostEntry+"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Pre-create asylum ssh dir with a custom entry and an existing key.
	sshDir := filepath.Join(home, ".asylum", "ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}
	customEntry := "custom.internal ssh-rsa BBBB..."
	if err := os.WriteFile(filepath.Join(sshDir, "known_hosts"), []byte(customEntry+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "id_ed25519"), []byte("existing key"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(sshDir, "known_hosts"))
	if err != nil {
		t.Fatalf("read known_hosts: %v", err)
	}
	content := string(got)
	if !strings.Contains(content, hostEntry) {
		t.Errorf("merged known_hosts missing host entry %q", hostEntry)
	}
	if !strings.Contains(content, customEntry) {
		t.Errorf("merged known_hosts missing custom entry %q", customEntry)
	}
}

// TestInit_NoKnownHosts verifies Init succeeds when ~/.ssh/known_hosts does
// not exist (no copy attempted).
func TestInit_NoKnownHosts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Pre-create the key so Init doesn't invoke ssh-keygen.
	sshDir := filepath.Join(home, ".asylum", "ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "id_ed25519"), []byte("existing key"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init returned error when no known_hosts: %v", err)
	}

	dst := filepath.Join(sshDir, "known_hosts")
	if _, err := os.Stat(dst); err == nil {
		t.Error("known_hosts should not have been created when source does not exist")
	}
}
