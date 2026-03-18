package main

import (
	"reflect"
	"testing"
)

func TestResolveMode(t *testing.T) {
	tests := []struct {
		name        string
		positional  []string
		passthrough []string
		wantMode    runMode
		wantExtra   []string
		wantErr     bool
	}{
		{
			name:      "no positional defaults to agent",
			wantMode:  modeAgent,
			wantExtra: nil,
		},
		{
			name:        "no positional passes through passthrough args",
			passthrough: []string{"--some-flag"},
			wantMode:    modeAgent,
			wantExtra:   []string{"--some-flag"},
		},
		{
			name:       "shell positional",
			positional: []string{"shell"},
			wantMode:   modeShell,
			wantExtra:  nil,
		},
		{
			name:        "shell with --admin in passthrough",
			positional:  []string{"shell"},
			passthrough: []string{"--admin"},
			wantMode:    modeAdminShell,
			wantExtra:   nil,
		},
		{
			name:        "shell with other passthrough flags (no --admin)",
			positional:  []string{"shell"},
			passthrough: []string{"--verbose"},
			wantMode:    modeShell,
			wantExtra:   nil,
		},
		{
			name:       "shell with extra positional returns error",
			positional: []string{"shell", "extra"},
			wantErr:    true,
		},
		{
			name:       "ssh-init positional",
			positional: []string{"ssh-init"},
			wantMode:   modeSSHInit,
			wantExtra:  nil,
		},
		{
			name:       "ssh-init with extra positional returns error",
			positional: []string{"ssh-init", "extra"},
			wantErr:    true,
		},
		{
			name:        "arbitrary command mode",
			positional:  []string{"run"},
			passthrough: []string{"arg1", "arg2"},
			wantMode:    modeCommand,
			wantExtra:   []string{"run", "arg1", "arg2"},
		},
		{
			name:       "arbitrary command no passthrough",
			positional: []string{"ls"},
			wantMode:   modeCommand,
			wantExtra:  []string{"ls"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, extra, err := resolveMode(tt.positional, tt.passthrough)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if mode != tt.wantMode {
				t.Errorf("mode = %d, want %d", mode, tt.wantMode)
			}
			if !reflect.DeepEqual(extra, tt.wantExtra) {
				t.Errorf("extra = %v, want %v", extra, tt.wantExtra)
			}
		})
	}
}

func TestParseArgs_Version(t *testing.T) {
	flags, positional, passthrough := parseArgs([]string{"--version"})

	if !flags.Version {
		t.Error("expected Version flag to be true")
	}
	if len(positional) != 0 {
		t.Errorf("expected no positional args, got %v", positional)
	}
	if len(passthrough) != 0 {
		t.Errorf("expected no passthrough args, got %v", passthrough)
	}
}

func TestParseArgs_AgentShorthand(t *testing.T) {
	flags, _, _ := parseArgs([]string{"-a", "gemini"})
	if flags.Agent != "gemini" {
		t.Errorf("agent = %q, want %q", flags.Agent, "gemini")
	}

	flags2, _, _ := parseArgs([]string{"-acodex"})
	if flags2.Agent != "codex" {
		t.Errorf("agent = %q, want %q", flags2.Agent, "codex")
	}
}

func TestParseArgs_DoubleDashSeparator(t *testing.T) {
	flags, positional, passthrough := parseArgs([]string{"--", "fix", "the", "bug"})
	if flags.Version || flags.New {
		t.Error("no flags expected before --")
	}
	if len(positional) != 0 {
		t.Errorf("expected no positional, got %v", positional)
	}
	if !reflect.DeepEqual(passthrough, []string{"fix", "the", "bug"}) {
		t.Errorf("passthrough = %v, want [fix, the, bug]", passthrough)
	}
}

func TestParseArgs_UnknownFlagPassthrough(t *testing.T) {
	_, _, passthrough := parseArgs([]string{"--unknown-flag", "val"})
	if len(passthrough) == 0 || passthrough[0] != "--unknown-flag" {
		t.Errorf("unknown flag should be in passthrough, got %v", passthrough)
	}
}

func TestParseArgs_PositionalArgRouting(t *testing.T) {
	// "shell" is a known positional — no passthrough collected
	_, positional, passthrough := parseArgs([]string{"shell"})
	if len(positional) != 1 || positional[0] != "shell" {
		t.Errorf("positional = %v, want [shell]", positional)
	}
	if len(passthrough) != 0 {
		t.Errorf("passthrough = %v, want empty", passthrough)
	}

	// Unknown positional routes everything after it to passthrough
	_, positional2, passthrough2 := parseArgs([]string{"run", "arg1", "arg2"})
	if len(positional2) != 1 || positional2[0] != "run" {
		t.Errorf("positional = %v, want [run]", positional2)
	}
	if !reflect.DeepEqual(passthrough2, []string{"arg1", "arg2"}) {
		t.Errorf("passthrough = %v, want [arg1, arg2]", passthrough2)
	}
}
