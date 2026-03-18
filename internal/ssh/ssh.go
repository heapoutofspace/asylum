package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/inventage-ai/asylum/internal/log"
)

// mergeKnownHosts combines host entries with any existing entries in dst,
// deduplicating by full line content. Host entries come first, then any
// extra lines from the existing file that weren't in the host data.
func mergeKnownHosts(hostData []byte, dstPath string) []byte {
	hostLines := nonEmptyLines(string(hostData))

	seen := make(map[string]bool, len(hostLines))
	for _, l := range hostLines {
		seen[l] = true
	}

	existing, err := os.ReadFile(dstPath)
	if err == nil {
		for _, l := range nonEmptyLines(string(existing)) {
			if !seen[l] {
				hostLines = append(hostLines, l)
				seen[l] = true
			}
		}
	}

	if len(hostLines) == 0 {
		return hostData
	}
	return []byte(strings.Join(hostLines, "\n") + "\n")
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home dir: %w", err)
	}

	sshDir := filepath.Join(home, ".asylum", "ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("create ssh dir: %w", err)
	}

	knownHosts := filepath.Join(home, ".ssh", "known_hosts")
	if info, err := os.Stat(knownHosts); err == nil && !info.IsDir() {
		hostData, err := os.ReadFile(knownHosts)
		if err != nil {
			return fmt.Errorf("read known_hosts: %w", err)
		}
		dst := filepath.Join(sshDir, "known_hosts")
		merged := mergeKnownHosts(hostData, dst)
		if err := os.WriteFile(dst, merged, 0600); err != nil {
			return fmt.Errorf("write known_hosts: %w", err)
		}
		log.Success("merged known_hosts")
	}

	keyPath := filepath.Join(sshDir, "id_ed25519")
	if _, err := os.Stat(keyPath); err == nil {
		log.Info("SSH key already exists at %s", keyPath)
		log.Info("replace with your own keys if needed")
		return nil
	}

	hostname, _ := os.Hostname()
	comment := fmt.Sprintf("asylum@%s", hostname)

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-C", comment, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ssh-keygen: %w", err)
	}

	pubKey, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}

	log.Success("SSH key generated")
	fmt.Printf("\nPublic key:\n%s\n", pubKey)
	log.Info("add this key to your Git hosting provider")
	log.Info("or replace with your own keys at %s", sshDir)

	return nil
}
