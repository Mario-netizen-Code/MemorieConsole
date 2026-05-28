---
name: same-session
description: Use when executing implementation plans with independent tasks in the current session. Dispatches fresh subagents per task with two-stage review (spec compliance then code quality).
---

# Same Session

## Overview

The same-session workflow executes plan tasks in the current workspace without isolation. It dispatches individual subagents for each independent task, reviews each one in two stages (spec compliance first, code quality second), and handles the full lifecycle of subagent communication. Use this when tasks do not share mutable state and can run in parallel.

## When to Use

Use this flowchart to decide between same-session and session-isolated:

```
Are tasks independent (no shared mutable state)?
├── YES ──► Can they fail independently without affecting others?
│   ├── YES ──► Use same-session (dispatch in parallel)
│   └── NO  ──► Use same-session (dispatch sequentially)
└── NO  ──► Are tasks tightly coupled or modifying shared code?
    ├── YES ──► Use session-isolated (git worktree)
    └── NO  ──► Use same-session (dispatch in parallel)
```

**Independent tasks** — tasks that read from shared data but do not write to the same files, modules, or state. Examples: refactoring two unrelated modules, adding tests for different components, writing documentation for separate features.

**Tightly coupled tasks** — tasks that modify the same files, depend on each other's output, or would cause merge conflicts if done in parallel. These require session-isolated.

## The Process

### 1. Read the Plan

Read the written implementation plan (`writing-plans` skill). Identify each task, its dependencies, and whether tasks can run in parallel or must be sequential.

### 2. Create a TodoWrite

Open a tracking document listing all tasks with checkboxes and status fields. This serves as the single source of truth for progress.

### 3. For Each Task (Sequential or Parallel)

#### Dispatch an Implementer Subagent

For each independent task, dispatch a fresh subagent. Provide:

- The task description from the plan.
- Any relevant context (files to modify, existing patterns, conventions).
- The expected outcome / acceptance criteria.

The subagent works autonomously and returns one of four statuses (see below).

#### Stage 1: Review Spec Compliance

After the subagent returns, review the output against the plan:

- Does the implementation match the specification?
- Are all acceptance criteria met?
- Are there missing edge cases or incomplete paths?
- Does it follow the project's existing patterns?

If spec compliance fails, either:
- Fix directly if trivial.
- Redispatch the subagent with specific failure details.

#### Stage 2: Review Code Quality

After spec compliance passes, review code quality:

- Is the code idiomatic for the project's language and frameworks?
- Are types correct (mypy/pyright passing)?
- Are there anti-patterns or unnecessary complexity?
- Are imports and naming consistent with the codebase?

If code quality fails, either:
- Fix directly if minor.
- Redispatch the subagent with quality-specific feedback.

#### Complete the Task

Mark the task done in the TodoWrite. Do not proceed to the next task until the current one passes both review stages (unless tasks were dispatched in parallel, in which case review each as they complete).

### 4. Final Review

Once all tasks are complete and reviewed, do a final check:

- All TodoWrite items are marked done.
- The workspace has no uncommitted changes (or changes are intentional).
- Lint/typecheck/tests pass on the full project.

Then reference `fase-calidad/code-review` for a structured review request.

## Handling Subagent Status

Subagents return one of four statuses:

| Status | Meaning | Action |
|--------|---------|--------|
| `DONE` | Task completed successfully, output is ready | Proceed to stage 1 review |
| `DONE_WITH_CONCERNS` | Task completed but subagent noted potential issues | Review concerns before proceeding; may need additional dispatch |
| `NEEDS_CONTEXT` | Subagent could not proceed due to missing information | Provide the missing context and redispatch |
| `BLOCKED` | Subagent encountered an unresolvable issue | Do not skip or work around. Document the blocker in the TodoWrite. Escalate for guidance |

Never modify a `BLOCKED` status to `DONE` or attempt to complete the task yourself without explicit instruction. "Stop when blocked" applies at the subagent level too.

## Parallel Investigation

When you encounter 2+ independent failures, bugs, or unknowns, dispatch one subagent per failure in parallel for investigation.

Use this when:
- Multiple test failures need root cause analysis.
- Multiple independent code areas need exploration.
- You need to compare approaches or research different solutions.

Dispatch each subagent with:
- A focused question or hypothesis to investigate.
- Relevant files and context.
- A timebox (e.g., "report findings within 5 minutes of tool usage").

Collect all findings, then synthesize into a unified approach before implementing. Do not dispatch implementation subagents in parallel with investigation subagents — investigate first, then implement.

## Key Principles

- **Two-stage review** — always spec compliance first, then code quality. Never skip or reorder.
- **Fresh subagent per task** — do not reuse subagents across tasks to avoid context contamination.
- **Stop when blocked** — at any level (subagent, task, review), if blocked, document and escalate. Never skip.
- **Investigate before implementing** — parallel investigation is for discovery, not execution.
- **Reference writing-plans** for creating the plan, and code-review for structured review after all tasks are done.
