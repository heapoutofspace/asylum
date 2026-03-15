## Context

The Asylum Dockerfile has an existing `npm install -g` step for Node.js global packages. OpenSpec just needs to be added to that list.

## Goals / Non-Goals

**Goals:**
- `openspec` command available in PATH inside the container

**Non-Goals:**
- No configuration or setup beyond the npm install

## Decisions

- Add to the existing npm install layer rather than a separate RUN step, to keep the Dockerfile clean and minimize layers.

## Risks / Trade-offs

None — single package addition.
