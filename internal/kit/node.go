package kit

import "github.com/inventage-ai/asylum/internal/onboarding"

func init() {
	Register(&Kit{
		Name:        "node",
		Description: "Node.js global development packages",
		DockerSnippet: `# Install Node.js global packages
RUN bash -c 'export PATH="$HOME/.local/share/fnm:$PATH" && eval "$(fnm env)" && \
    npm install -g \
        typescript \
        @types/node \
        ts-node \
        eslint \
        prettier \
        nodemon'
`,
		BannerLines: `    echo "Node.js:   $(node --version 2>/dev/null || echo 'not found')"
`,
		SubKits: map[string]*Kit{
			"npm": {
				Name:            "node/npm",
				Description:     "npm with caching and onboarding",
				CacheDirs:       map[string]string{"npm": "/home/claude/.npm"},
				OnboardingTasks: []onboarding.Task{onboarding.NPMTask{}},
			},
			"pnpm": {
				Name:        "node/pnpm",
				Description: "pnpm global install",
				DockerSnippet: `# Install pnpm
RUN bash -c 'export PATH="$HOME/.local/share/fnm:$PATH" && eval "$(fnm env)" && npm install -g pnpm'
`,
			},
			"yarn": {
				Name:        "node/yarn",
				Description: "yarn global install",
				DockerSnippet: `# Install yarn
RUN bash -c 'export PATH="$HOME/.local/share/fnm:$PATH" && eval "$(fnm env)" && npm install -g yarn'
`,
			},
		},
	})
}
