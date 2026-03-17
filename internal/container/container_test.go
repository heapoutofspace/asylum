package container

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyDir(t *testing.T) {
	t.Run("copies files and nested directories", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()

		os.MkdirAll(filepath.Join(src, "sub"), 0755)
		os.WriteFile(filepath.Join(src, "file.txt"), []byte("hello"), 0644)
		os.WriteFile(filepath.Join(src, "sub", "nested.txt"), []byte("world"), 0644)

		if err := copyDir(src, dst); err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile(filepath.Join(dst, "file.txt"))
		if err != nil || string(data) != "hello" {
			t.Errorf("file.txt: got %q, err %v", data, err)
		}
		data, err = os.ReadFile(filepath.Join(dst, "sub", "nested.txt"))
		if err != nil || string(data) != "world" {
			t.Errorf("sub/nested.txt: got %q, err %v", data, err)
		}
	})

	t.Run("preserves file permissions", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()

		os.WriteFile(filepath.Join(src, "exec.sh"), []byte("#!/bin/sh"), 0755)

		if err := copyDir(src, dst); err != nil {
			t.Fatal(err)
		}

		info, err := os.Stat(filepath.Join(dst, "exec.sh"))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0755 {
			t.Errorf("permissions = %o, want 0755", info.Mode().Perm())
		}
	})

	t.Run("recreates symlinks", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()

		os.WriteFile(filepath.Join(src, "target.txt"), []byte("data"), 0644)
		os.Symlink("target.txt", filepath.Join(src, "link.txt"))

		if err := copyDir(src, dst); err != nil {
			t.Fatal(err)
		}

		linkTarget, err := os.Readlink(filepath.Join(dst, "link.txt"))
		if err != nil {
			t.Fatalf("Readlink: %v", err)
		}
		if linkTarget != "target.txt" {
			t.Errorf("symlink target = %q, want %q", linkTarget, "target.txt")
		}
	})

	t.Run("propagates error on unreadable source file", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("root ignores permission bits")
		}
		src := t.TempDir()
		dst := t.TempDir()

		path := filepath.Join(src, "unreadable.txt")
		os.WriteFile(path, []byte("data"), 0000)
		defer os.Chmod(path, 0644)

		if err := copyDir(src, dst); err == nil {
			t.Error("expected error reading unreadable file")
		}
	})
}

func TestSafeHostname(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want string
	}{
		{
			name: "simple name",
			dir:  "/home/user/myproject",
			want: "asylum-myproject",
		},
		{
			name: "underscores become dashes",
			dir:  "/home/user/my_project",
			want: "asylum-my-project",
		},
		{
			name: "uppercase lowercased",
			dir:  "/home/user/MyProject",
			want: "asylum-myproject",
		},
		{
			name: "leading dash stripped",
			dir:  "/home/user/_project",
			want: "asylum-project",
		},
		{
			name: "trailing dash stripped after truncation",
			// base name: 56 a's + hyphen + more: truncation at 56 lands on hyphen
			dir:  "/home/user/" + strings.Repeat("a", 55) + "-extra",
			want: "asylum-" + strings.Repeat("a", 55),
		},
		{
			name: "exact 56-char input not truncated",
			dir:  "/home/user/" + strings.Repeat("a", 56),
			want: "asylum-" + strings.Repeat("a", 56),
		},
		{
			name: "all non-alphanumeric becomes dashes then empty -> project",
			dir:  "/home/user/___",
			want: "asylum-project",
		},
		{
			name: "empty base falls back to project",
			dir:  "/",
			want: "asylum-project",
		},
		{
			name: "result within Docker 63-char limit",
			dir:  "/home/user/" + strings.Repeat("b", 63),
			want: "asylum-" + strings.Repeat("b", 56),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeHostname(tt.dir)
			if got != tt.want {
				t.Errorf("safeHostname(%q) = %q, want %q", tt.dir, got, tt.want)
			}
			if len(got) > 63 {
				t.Errorf("hostname too long: %d chars", len(got))
			}
			if strings.HasSuffix(got, "-") {
				t.Errorf("hostname has trailing dash: %q", got)
			}
		})
	}
}
