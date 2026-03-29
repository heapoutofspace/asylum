## MODIFIED Requirements

### Requirement: Update to latest stable release
The `self-update` subcommand SHALL query the GitHub Releases API for the latest non-prerelease and download the matching binary for the current OS and architecture, replacing the running binary atomically. When a version argument is provided, the specified release tag SHALL be fetched instead.

#### Scenario: Successful stable update
- **WHEN** `asylum self-update` is run and a newer stable release exists
- **THEN** the binary is downloaded, replaces the current executable, and the new version is printed

#### Scenario: Already up to date
- **WHEN** `asylum self-update` is run and the running version matches the latest stable release
- **THEN** a message is printed indicating the binary is already up to date, and no download occurs

#### Scenario: Targeted version update
- **WHEN** `asylum self-update 0.4.0` is run
- **THEN** the release tagged `v0.4.0` is fetched and installed

#### Scenario: Already at targeted version
- **WHEN** `asylum self-update 0.4.0` is run and the current version is `0.4.0`
- **THEN** a message is printed indicating the binary is already at the requested version
