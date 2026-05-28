---
name: git-worktrees
description: Use when starting feature work that needs isolation from current workspace. Ensures isolated workspace exists via native tools or git worktree fallback.
---

# Git Worktrees

## Overview

Ensure work happens in an isolated workspace. Prefer native worktree tools. Fall back to manual git worktrees only when no native tool is available.

**Core principle:** Detect existing isolation first. Use native tools. Fall back to git. Never fight the harness.

## Step 0: Detect Existing Isolation

Run `git rev-parse --git-dir` — if the output path differs from the main repo's `.git` directory, you are inside a linked worktree. To detect programmatically:

**PowerShell:**
```powershell
$gitDir = git rev-parse --git-dir 2>$null
$gitCommon = git rev-parse --git-common-dir 2>$null
$isSubmodule = git rev-parse --show-superproject-working-tree 2>$null
```

**Linux/macOS:**
```bash
GIT_DIR=$(git rev-parse --git-dir 2>/dev/null)
GIT_COMMON=$(git rev-parse --git-common-dir 2>/dev/null)
IS_SUBMODULE=$(git rev-parse --show-superproject-working-tree 2>/dev/null)
```

**Logic:**
- If `$gitDir -ne $gitCommon` AND `$isSubmodule` is empty → you are in a linked worktree. Skip to Project Setup.
- If `$isSubmodule` is non-empty → treat as normal repo (submodule guard).
- Otherwise → normal repo. Ask for consent to create a worktree. Honor declared preference without asking. If declined, work in place.

## Step 1: Create Worktree

**1a. Native tools (preferred):** If a worktree tool exists (e.g., `EnterWorktree`, `WorktreeCreate`), use it. Skip to Project Setup. Native tools handle placement, branches, and cleanup automatically.

**1b. Git fallback (only if no native tool):**

**Directory priority:**
1. Explicit user preference from instructions
2. Existing `.worktrees/` (hidden, preferred) or `worktrees/` directory
3. Global `~/.config/superpowers/worktrees/<project>/` (backward compat)
4. Default to `.worktrees/` at project root

## Safety Verification

For project-local directories, MUST verify the directory is gitignored:

```bash
git check-ignore -q .worktrees || git check-ignore -q worktrees
```

If NOT ignored: add to `.gitignore`, commit, then proceed. Global directories need no verification.

Create: `git worktree add "<path>" -b "<branch>"` and `cd` into it.

**Sandbox fallback:** If permission error, work in place instead.

## Project Setup

Auto-detect and run:
- `package.json` → `npm install`
- `Cargo.toml` → `cargo build`
- `requirements.txt` → `pip install -r requirements.txt`
- `pyproject.toml` → `poetry install`
- `go.mod` → `go mod download`

## Verify Baseline

Run project-appropriate test command. If tests fail, report and ask whether to proceed. If passing, report ready.

```
Worktree ready at <path>
Tests passing (<N> tests, 0 failures)
Ready to implement <feature>
```
