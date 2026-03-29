package kit

func init() {
	Register(&Kit{
		Name:        "title",
		Description: "Terminal tab title and agent title configuration",
		DefaultOn:   true,
		ConfigSnippet: `  # title:              # Terminal tab title configuration
  #   # Placeholders: {project}, {agent}, {mode}
  #   tab-title: "🤖 {project}"
  #   allow-agent-terminal-title: false
`,
	})
}
