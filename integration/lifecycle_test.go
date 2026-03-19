//go:build integration

package integration_test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/inventage-ai/asylum/internal/docker"
)

// startDetachedContainer starts an asylum container in detached mode with
// sleep infinity, mimicking the new lifecycle. Returns the container name.
func startDetachedContainer(t *testing.T, name string) {
	t.Helper()
	cmd := exec.Command("docker", "run", "-d", "--init",
		"--name", name,
		"asylum:latest", "sleep", "infinity")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker run -d failed: %v\noutput: %s", err, out)
	}
}

func TestDetachedContainerLifecycle(t *testing.T) {
	ensureBaseImage(t)
	name := "asylum-lifecycle-test"
	t.Cleanup(func() {
		exec.Command("docker", "rm", "-f", name).Run()
	})

	// Start detached container
	startDetachedContainer(t, name)

	if !docker.IsRunning(name) {
		t.Fatal("container should be running after detached start")
	}

	// Exec a command into it
	cmd := exec.Command("docker", "exec", name, "echo", "hello")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker exec failed: %v\noutput: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "hello" {
		t.Errorf("exec output = %q, want \"hello\"", strings.TrimSpace(string(out)))
	}

	// Container still running after exec exits
	if !docker.IsRunning(name) {
		t.Fatal("container should still be running after exec exits")
	}

	// Remove container
	docker.RemoveContainer(name)
	if docker.IsRunning(name) {
		t.Fatal("container should not be running after removal")
	}
}

func TestMultipleExecSessions(t *testing.T) {
	ensureBaseImage(t)
	name := "asylum-multi-session-test"
	t.Cleanup(func() {
		exec.Command("docker", "rm", "-f", name).Run()
	})

	startDetachedContainer(t, name)

	// Start a long-running exec session in the background
	session1 := exec.Command("docker", "exec", name, "sleep", "30")
	if err := session1.Start(); err != nil {
		t.Fatalf("start session 1: %v", err)
	}

	// Give it a moment to start
	time.Sleep(500 * time.Millisecond)

	// Run a second exec session that exits quickly
	cmd := exec.Command("docker", "exec", name, "echo", "session2")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("session 2 exec failed: %v\noutput: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "session2" {
		t.Errorf("session 2 output = %q, want \"session2\"", strings.TrimSpace(string(out)))
	}

	// Container still running (session 1 is still active)
	if !docker.IsRunning(name) {
		t.Fatal("container should still be running while session 1 is active")
	}

	// Kill session 1
	session1.Process.Kill()
	session1.Wait()

	// Container still running (sleep infinity keeps it alive)
	if !docker.IsRunning(name) {
		t.Fatal("container should still be running (sleep infinity is PID 1's child)")
	}
}

func TestEntrypointRunsInDetachedMode(t *testing.T) {
	ensureBaseImage(t)
	name := "asylum-entrypoint-detached-test"
	t.Cleanup(func() {
		exec.Command("docker", "rm", "-f", name).Run()
	})

	startDetachedContainer(t, name)

	// Wait for entrypoint to finish setup (sleep becomes the main process)
	for i := 0; i < 30; i++ {
		out, _ := exec.Command("docker", "exec", name, "pgrep", "-x", "sleep").CombinedOutput()
		if strings.TrimSpace(string(out)) != "" {
			break
		}
		time.Sleep(time.Second)
	}

	// Verify entrypoint ran: git safe.directory should be set
	out, err := exec.Command("docker", "exec", name, "git", "config", "--global", "--get-all", "safe.directory").CombinedOutput()
	if err != nil {
		t.Fatalf("git config failed: %v\noutput: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "*" {
		t.Errorf("safe.directory = %q, want \"*\"", strings.TrimSpace(string(out)))
	}
}
