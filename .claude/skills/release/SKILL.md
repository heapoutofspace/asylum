---
name: release
description: Create a new release — updates CHANGELOG.md, commits, and tags. Use when the user wants to cut a release.
disable-model-invocation: true
argument-hint: "[version]"
---

Create a new release for this project.

**Version**: `$ARGUMENTS`
If a version is provided (e.g., `/release 0.3.0`), use it. If no version is provided, determine the latest git tag (e.g., `v0.2.1`), increment the patch number by one (e.g., `0.2.2`), and use that.

**Steps**

1. **Verify branch**

   Run `git branch --show-current`. If the current branch is not `main`, abort with:
   "Releases must be created from the `main` branch. Current branch: `<branch>`."

2. **Check for uncommitted changes**

   Run `git status` to check for uncommitted changes (staged or unstaged, excluding untracked files).

   If there are uncommitted changes:
   - Show the user what's changed (summary, not full diff)
   - Ask: "There are uncommitted changes. Commit them before releasing, or ignore and release from the current HEAD?"
   - If they want to commit: create an appropriate commit first, then continue
   - If they want to ignore: continue without committing (the changes will remain uncommitted after the release)

3. **Ensure CHANGELOG.md is up to date**

   Read `CHANGELOG.md` and check the `## Unreleased` section.

   Run `git log <last-tag>..HEAD --oneline --no-merges` to see all commits since the last release.

   Compare the commits against the unreleased changelog entries. If significant changes are missing from the changelog:
   - Add the missing entries under the appropriate categories (Added/Changed/Fixed/Removed)
   - Show the user what was added

   If the Unreleased section is empty or has no meaningful entries:
   - Warn the user: "No unreleased changes in CHANGELOG.md. Are you sure you want to release?"
   - If confirmed, proceed. Otherwise, abort.

4. **Clean up and organize changelog entries**

   Before moving entries to the versioned section, process the unreleased entries:

   - **Remove internal-only entries**: Drop fixes for bugs introduced since the last release (not present in the previous release) and changes to features that were also added since the last release. Users upgrading from the last release never saw these bugs or intermediate states, so they're noise. Check git tags and the previous release's changelog to determine what was already shipped.
   - **Merge related entries**: If multiple entries describe changes to the same feature (e.g., an "Added" entry and a later "Changed" entry for the same thing), merge them into a single entry that describes the final state.
   - **Order by importance**: Within each category (Added/Changed/Fixed), order entries by user impact — breaking changes and major features first, minor improvements last.

5. **Move entries to a versioned section and add summary**

   In `CHANGELOG.md`:
   - Replace `## Unreleased` content with an empty Unreleased section
   - Insert a new version section below it: `## <version> — YYYY-MM-DD` (today's date)
   - Move the cleaned-up entries into this new section
   - Write a 2-3 sentence executive summary and place it between the version header and the first `### Added` section. The summary should highlight the most important user-facing changes in plain language — think "what would someone care about when deciding to upgrade?" Don't just list categories.

   The result should look like:
   ```
   ## Unreleased

   ## 0.3.0 — 2026-03-18

   Summary text goes here. Two to three sentences covering the highlights.

   ### Added
   - ...
   ```

6. **Commit and tag**

   - Stage only `CHANGELOG.md` (and any files from step 1 if the user chose to commit)
   - Create a commit with message: `Release v<version>`
   - Create an annotated tag: `git tag v<version> -m "v<version>"`

7. **Summary**

   Show:
   - The new version number
   - The changelog entries included in this release
   - Remind: `git push origin main --tags` to trigger the release workflow

   Do NOT push automatically — let the user decide when to push.
