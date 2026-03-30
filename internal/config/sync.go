package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// SyncKitToConfig inserts a kit's config snippet into the config file's kits
// block using text-based insertion (no YAML roundtrip, so comments and
// indentation are preserved). If the kit key already exists, no modification
// is made.
func SyncKitToConfig(path string, kitName string, snippet string) error {
	// Parse YAML read-only to check if kit already exists.
	doc, err := parseConfigDoc(path)
	if err != nil {
		return err
	}
	if kitsNode := findKitsMapping(doc); kitsNode != nil {
		if kitExistsInMapping(kitsNode, kitName) {
			return nil
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	// Find the "kits:" line.
	kitsIdx := -1
	kitsLineIndent := 0
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(trimmed, "kits:") {
			kitsIdx = i
			kitsLineIndent = len(line) - len(trimmed)
			break
		}
	}
	if kitsIdx < 0 {
		return fmt.Errorf("no kits: mapping found in %s", path)
	}

	entryIndent := kitsLineIndent + 2

	// Find insertion point: after the last active kit entry's full block,
	// before any commented-out kit entries at the entry indent level.
	insertIdx := kitsIdx + 1
	inKitBlock := false
	for i := kitsIdx + 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		isComment := strings.HasPrefix(trimmed, "#")

		// Past the kits block entirely.
		if indent <= kitsLineIndent && !isComment {
			break
		}

		if indent == entryIndent {
			if isComment {
				break // commented-out kit entry — insert before this
			}
			inKitBlock = true
			insertIdx = i + 1
			continue
		}

		// Deeper than entry level: part of the current kit's config.
		if inKitBlock {
			insertIdx = i + 1
		}
	}

	snippetText := strings.TrimRight(snippet, "\n")
	var insert []string
	if insertIdx > kitsIdx+1 {
		insert = append(insert, "") // blank line separator
	}
	insert = append(insert, snippetText)

	result := make([]string, 0, len(lines)+len(insert))
	result = append(result, lines[:insertIdx]...)
	result = append(result, insert...)
	result = append(result, lines[insertIdx:]...)

	return os.WriteFile(path, []byte(strings.Join(result, "\n")), 0644)
}

// SyncKitCommentToConfig appends a commented-out kit block to the config
// file's kits mapping as a foot comment.
func SyncKitCommentToConfig(path string, comment string) error {
	doc, err := parseConfigDoc(path)
	if err != nil {
		return err
	}

	kitsNode := findOrCreateKitsMapping(doc)

	// Append as foot comment on the kits mapping node
	if kitsNode.FootComment != "" {
		kitsNode.FootComment += "\n\n" + comment
	} else {
		kitsNode.FootComment = comment
	}
	return writeConfigDoc(path, doc)
}

// parseConfigDoc reads a YAML file into a yaml.Node document tree.
func parseConfigDoc(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// writeConfigDoc encodes a yaml.Node document tree back to a file.
func writeConfigDoc(path string, doc *yaml.Node) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// findKitsMapping walks the document to find the "kits" mapping node.
// Returns nil if not found.
func findKitsMapping(doc *yaml.Node) *yaml.Node {
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "kits" && root.Content[i+1].Kind == yaml.MappingNode {
			return root.Content[i+1]
		}
	}
	return nil
}

// findOrCreateKitsMapping walks the document to find the "kits" mapping node.
// If it doesn't exist, it creates one and appends it to the root mapping.
func findOrCreateKitsMapping(doc *yaml.Node) *yaml.Node {
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		// Create minimal document structure
		root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		doc.Kind = yaml.DocumentNode
		doc.Content = []*yaml.Node{root}
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	}

	// Walk mapping pairs looking for "kits" key
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "kits" {
			if root.Content[i+1].Kind == yaml.MappingNode {
				return root.Content[i+1]
			}
			// Key exists but value isn't a mapping — replace it
			root.Content[i+1] = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			return root.Content[i+1]
		}
	}

	// No "kits" key — create it
	kitsKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "kits", Tag: "!!str"}
	kitsVal := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	root.Content = append(root.Content, kitsKey, kitsVal)
	return kitsVal
}

// kitExistsInMapping checks if a key is already present in a mapping node.
func kitExistsInMapping(mapping *yaml.Node, name string) bool {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == name {
			return true
		}
	}
	return false
}
