package selfupdate

import (
	"runtime"
	"testing"
)

func TestAssetName(t *testing.T) {
	want := "asylum-" + runtime.GOOS + "-" + runtime.GOARCH
	if got := AssetName(); got != want {
		t.Errorf("AssetName() = %q, want %q", got, want)
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"", ""},
		{"0.4.0", "v0.4.0"},
		{"v0.4.0", "v0.4.0"},
		{"v1.0.0-rc1", "v1.0.0-rc1"},
	}
	for _, tt := range tests {
		if got := NormalizeVersion(tt.input); got != tt.want {
			t.Errorf("NormalizeVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveChannel(t *testing.T) {
	tests := []struct {
		name          string
		devFlag       bool
		configChannel string
		want          string
	}{
		{"no flag no config", false, "", "stable"},
		{"no flag stable config", false, "stable", "stable"},
		{"no flag dev config", false, "dev", "dev"},
		{"dev flag overrides stable config", true, "stable", "dev"},
		{"dev flag with no config", true, "", "dev"},
		{"dev flag with dev config", true, "dev", "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveChannel(tt.devFlag, tt.configChannel); got != tt.want {
				t.Errorf("ResolveChannel(%v, %q) = %q, want %q", tt.devFlag, tt.configChannel, got, tt.want)
			}
		})
	}
}
