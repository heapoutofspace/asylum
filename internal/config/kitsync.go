package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/inventage-ai/asylum/internal/kit"
	"github.com/inventage-ai/asylum/internal/log"
)

// SyncNewKits detects kits not yet in state.json, prompts for activation
// (if interactive), updates the config file, and saves the new state.
// Returns true if any new kits were processed.
func SyncNewKits(asylumDir string, interactive bool) (bool, error) {
	state, err := LoadState(asylumDir)
	if err != nil {
		return false, fmt.Errorf("load state: %w", err)
	}

	newKits := NewKits(kit.All(), state)
	if len(newKits) == 0 {
		return false, nil
	}

	configPath := filepath.Join(asylumDir, "config.yaml")

	for _, name := range newKits {
		k := kit.Get(name)
		if k == nil {
			continue
		}

		switch k.Tier {
		case kit.TierAlwaysOn:
			log.Info("new kit: %s (always active)", name)

		case kit.TierDefault:
			activate := true
			if interactive {
				activate = promptActivateKit(name, k.Description)
			}
			if activate && interactive {
				if k.ConfigNodes != nil {
					if err := SyncKitToConfig(configPath, name, k.ConfigNodes); err != nil {
						log.Error("sync kit %s: %v", name, err)
					}
				}
			} else {
				// Non-interactive or declined: add as comment
				if k.ConfigComment != "" {
					if err := SyncKitCommentToConfig(configPath, k.ConfigComment); err != nil {
						log.Error("sync kit %s: %v", name, err)
					}
				}
				log.Info("kit %s added as comment — uncomment in config.yaml to enable", name)
			}

		case kit.TierOptIn:
			log.Info("new kit available: %s — uncomment in config.yaml to enable", name)
			if k.ConfigComment != "" {
				if err := SyncKitCommentToConfig(configPath, k.ConfigComment); err != nil {
					log.Error("sync kit %s: %v", name, err)
				}
			}
		}
	}

	// Update state with all currently registered kits
	state.KnownKits = kit.All()
	if err := SaveState(asylumDir, state); err != nil {
		return true, fmt.Errorf("save state: %w", err)
	}

	return true, nil
}

func promptActivateKit(name, description string) bool {
	label := name
	if description != "" {
		label = name + " (" + description + ")"
	}
	fmt.Printf("  New kit: %s — activate? [Y/n] ", label)
	var answer string
	fmt.Scanln(&answer)
	return !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "n")
}
