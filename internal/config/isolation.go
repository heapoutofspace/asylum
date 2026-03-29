package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// SetAgentIsolation writes the config isolation level for an agent to the
// given config file. Uses yaml.Node to preserve comments and formatting.
func SetAgentIsolation(path, agentName, level string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	// doc is a Document node; its first child is the mapping
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil
	}

	// Find or create "agents" key
	agentsNode := findMapValue(root, "agents")
	if agentsNode == nil {
		root.Content = append(root.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "agents"},
			&yaml.Node{Kind: yaml.MappingNode},
		)
		agentsNode = root.Content[len(root.Content)-1]
	}

	// Find or create the agent entry
	agentNode := findMapValue(agentsNode, agentName)
	if agentNode == nil {
		agentsNode.Content = append(agentsNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: agentName},
			&yaml.Node{Kind: yaml.MappingNode},
		)
		agentNode = agentsNode.Content[len(agentsNode.Content)-1]
	}

	// If the agent node is null/empty scalar, convert to mapping
	if agentNode.Kind == yaml.ScalarNode {
		agentNode.Kind = yaml.MappingNode
		agentNode.Value = ""
		agentNode.Tag = ""
		agentNode.Content = nil
	}

	// Find or create "config" key inside the agent
	configNode := findMapValue(agentNode, "config")
	if configNode == nil {
		agentNode.Content = append(agentNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "config"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: level},
		)
	} else {
		configNode.Value = level
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

// findMapValue finds the value node for a key in a mapping node.
func findMapValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}
