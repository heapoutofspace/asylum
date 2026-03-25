## MODIFIED Requirements

### Requirement: Dockerfile decomposition
The monolithic Dockerfile SHALL be split into an embedded core fragment (OS, shell, build tools, Docker CLI, language managers) and an embedded tail fragment (oh-my-zsh, shell config, entrypoint COPY, final USER/WORKDIR). The Docker engine installation is NOT part of the core fragment — it is assembled from the docker kit's DockerSnippet.

#### Scenario: Core fragment content
- **WHEN** the core Dockerfile is examined
- **THEN** it contains OS package installation, Docker CLI installation, GitHub/GitLab CLIs, user creation, mise/fnm/uv manager installation — but no Docker engine installation, no agent CLI installations, and no language-specific installations
