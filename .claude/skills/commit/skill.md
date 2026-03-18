---
name: commit
description: Create a commit from current changes with an appropriate message. Handles OpenSpec change lifecycle (verify, archive) and CHANGELOG updates automatically.
disable-model-invocation: true
---

Create a commit from the current working tree changes.

**Steps**

1. **Assess current changes**

   Run `git status` (excluding untracked files) and `git diff --stat` to understand what's changed.

   If there are no staged or unstaged changes, abort with: "Nothing to commit."

2. **Check for active OpenSpec changes**

   Run `openspec list --json` to check for active (non-archived) changes.

   For each active change:
   - Read its `tasks.md` to check task completion status
   - Determine if the current changes relate to this change (check if modified files overlap with what the change's tasks describe)

   **If changes relate to an active OpenSpec change:**

   a. **Check task completion**:
      - Parse `tasks.md` for `- [ ]` (incomplete) vs `- [x]` (complete)
      - If there are incomplete tasks: ask the user "Change '<name>' has X incomplete tasks. Commit anyway?"
        - If user declines: abort
        - If user confirms: proceed without archiving
      - If all tasks are complete: continue to verification

   b. **Run verification if not already done this session**:
      - Use the Skill tool to invoke `opsx:verify` for the change
      - If CRITICAL issues are found: alert the user and ask whether to proceed
      - If no critical issues (or verify was already run): continue to archive

   c. **Archive the change**:
      - Use the Skill tool to invoke `opsx:archive` for the change

3. **Update CHANGELOG.md**

   Read `CHANGELOG.md` and the `## Unreleased` section.

   Review the changes being committed. If they represent a user-facing change (new feature, bug fix, behavior change), check whether the Unreleased section already mentions it.

   If a significant change is missing:
   - Add a concise entry under the appropriate category (Added/Changed/Fixed/Removed)
   - Show the user what was added
   - Stage `CHANGELOG.md`

   Skip CHANGELOG updates for:
   - Pure refactoring with no behavior change
   - Test-only changes
   - Documentation/comment updates
   - OpenSpec artifact files

4. **Stage and commit**

   - Stage all relevant changed files (modified + new files related to the work)
   - Do NOT stage files that look like secrets (`.env`, credentials, etc.)
   - Do NOT stage unrelated changes that happen to be in the working tree — ask the user if unclear
   - Draft a concise commit message (1-2 sentences) that focuses on "why" not "what"
   - Create the commit

5. **Show summary**

   Display:
   - Commit hash and message
   - Files included
   - Whether a change was archived
   - Whether CHANGELOG was updated
   - `git status` after commit to confirm clean state

**Commit Message Guidelines**

- First line: imperative mood, under 72 chars (e.g., "Add worktree volume mounting")
- Focus on the user-facing change, not implementation details
- For bug fixes: describe what was broken
- For features: describe the capability added
- For refactors: "refactor:" prefix, describe what was simplified

**Guardrails**
- Never commit files that likely contain secrets
- Never amend existing commits — always create new ones
- If the working tree has changes to many unrelated files, ask the user which to include
- Always read the diff before writing the commit message
- Use HEREDOC syntax for commit messages to preserve formatting
- Do NOT push after committing
