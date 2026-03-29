## MODIFIED Requirements

### Requirement: GitHub Pages deployment
A GitHub Actions workflow at `.github/workflows/docs.yml` SHALL deploy the `dev` version of the docs site to GitHub Pages using mike on push to main. The release workflow at `.github/workflows/release.yml` SHALL deploy stable versions using mike when a version tag is pushed.

#### Scenario: Dev docs deploy on push
- **WHEN** a commit is pushed to main that changes files in `docs/` or `mkdocs.yml`
- **THEN** the workflow deploys the `dev` version via `mike deploy dev --push`

#### Scenario: Stable docs deploy on release
- **WHEN** a version tag (`v*`) is pushed
- **THEN** the release workflow deploys the docs as a stable version via `mike deploy <version> latest --update-aliases --push`

#### Scenario: No deploy on unrelated changes
- **WHEN** a commit is pushed to main that does not change `docs/` or `mkdocs.yml`
- **THEN** the docs workflow does not run
