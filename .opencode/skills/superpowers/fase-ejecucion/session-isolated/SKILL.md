---
name: session-isolated
description: Use when you have a written implementation plan and need to execute it in an isolated workspace with review checkpoints. Creates or detects isolated workspace via git worktrees, then executes plan with verified tasks.
---

# Session Isolated

## Overview

The session-isolated workflow handles the full lifecycle of executing an implementation plan in a workspace that is isolated from your main development branch. It checks for existing isolation, creates a workspace if needed, executes plan tasks with verification checkpoints, and hands off to code review and verification skills when done.

## Step 0: Detect Existing Isolation

Before creating anything, check if already isolated. Use `infraestructura/git-worktrees` Step 0 for the full detection logic. In summary:

1. **Compare `git rev-parse --git-dir` vs `--git-common-dir`** — different paths = linked worktree.
2. **Submodule guard** — `git rev-parse --show-superproject-working-tree` — if returns a path, it's a submodule (treat as normal repo).
3. **CI check** — if `CI`, `GITHUB_ACTIONS`, or similar env vars are set → skip worktree creation entirely (already in CI workspace).

If already isolated (linked worktree, not submodule): skip to Step 2.
If CI environment: execute plan tasks in-place; no worktree needed.
Otherwise: proceed to Step 1 for worktree creation.

## Step 1: Create Isolated Workspace

Delegate to `infraestructura/git-worktrees` — follow its instructions for native tools or `git worktree` fallback. After the worktree is created and verified, return here.

## Step 2: Execute Plan Tasks

1. **Load the plan** — read the written plan file. Understand each task, its prerequisites, and expected outcome.
2. **Create a TodoWrite** — open a tracking document listing all tasks with checkboxes.
3. **For each task, execute in order:**
   a. Complete the implementation.
   b. Run verification commands (lint, typecheck, tests).
   c. If the task is blocked, follow the "stop when blocked" principle (see below).
   d. Mark the task complete in the TodoWrite.
   e. Optionally commit with a descriptive message: `git add .` then `git commit -m "<scope>: <description>"`
4. **Stop when blocked** — if a task cannot proceed due to unresolved issues, do not attempt workarounds. Document the blocker, report the status, and wait for guidance. Do not skip blocked tasks and continue.

## Step 3: Complete Development

Once all tasks are executed, do not declare completion. Instead:

1. Reference `fase-calidad/code-review` — request a code review of the work.
2. Reference `fase-cierre/verify-and-finish` — verify everything and finalize.

Follow the instructions in those skills to close out the work. Do not skip directly to commit or merge.

## Key Principles

- **Always detect first** — never assume you need to create a workspace from scratch.
- **Delegate to git-worktrees** — use `infraestructura/git-worktrees` for workspace creation; don't duplicate those steps.
- **Verify after every task** — always run the project's lint/typecheck/test commands after completing implementation work.
- **Stop when blocked** — do not invent workarounds, skip tasks, or proceed with uncertainty. Document and escalate.
- **Hand off, don't finish** — after executing all tasks, pass to code-review and verify-and-finish rather than declaring done.
