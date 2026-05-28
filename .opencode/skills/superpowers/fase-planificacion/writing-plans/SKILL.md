---
name: writing-plans
description: Use when you have a spec or requirements for a multi-step task, before touching code. Creates detailed implementation plans with bite-sized tasks.
---

# Writing Plans

## Overview

Creates detailed, executable implementation plans from an approved design spec or set of requirements. Each plan assumes the engineer has zero context for the codebase — every action, file path, and command must be explicitly documented.

The output is a plan document that can be handed off to a subagent or executed inline.

## Plan Header Format

Every plan document starts with a header:

```markdown
# Plan: <Feature or Fix Name>

**Inspired by:** `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`

**Sub-skills required:** [sub-skill-1], [sub-skill-2]

## Tasks
```

The `Sub-skills required` field is **REQUIRED**. List any skill names needed during execution (e.g., `python-testing-patterns`, `systematic-debugging`). If none, write `None`.

## Task Structure

Each task follows this format:

```markdown
### Task N: <short verb phrase>

**Files:**
- `Create/Modify: path/to/file.py` — what changes
- `Test: path/to/test_file.py` — what to assert (if applicable)

**Steps:**

1. Run `pytest path/to/test_file.py -k "test_name"` — expect failure (red)
2. Edit `path/to/file.py` to add implementation — 3-5 lines of actual code shown here
3. Run `pytest path/to/test_file.py -k "test_name"` — expect pass (green)
4. Run `.\venv\Scripts\python -m ruff check .` — must pass
5. Run `.\venv\Scripts\python -m mypy .` — must pass
6. Commit with message: "feat: implement <thing>"
```

### Bite-Sized Tasks

Each task must be **2-5 minutes** of work maximum. If a step feels longer, split it:

- **Too big**: "Implement the full parser"
- **Good**: "Add tokenizer for numeric literals" / "Add tokenizer for string literals" / "Wire tokenizer into parser"

Each task should have exactly one logical change: write failing test → make it pass → verify.

### Files Section

Every task must list every file it touches with one of these prefixes:

- `Create:` — new file
- `Modify:` — existing file
- `Delete:` — file to remove
- `Test:` — test file (always include the test command)

### Code in Tasks

Each task must contain the **exact code** to write, not a description:

```
**Steps:**

1. Edit `models.py` — add after the `Chunk` class:
   ```python
   @dataclass
   class Document:
       chunks: list[Chunk]
       metadata: dict[str, str]
   ```
```

## No Placeholders

**"TBD", "TODO", "implement later", and empty sections are plan failures.** Every task must contain actual, complete code or content. If you cannot write the code because you lack information, ask the user before writing the plan.

- ❌ `# TODO: add validation later`
- ✅ `if not value: raise ValueError("value must not be empty")`
- ❌ `# implement the sorting function`
- ✅ A complete sorting function with correct logic, types, and tests

## Self-Review

Before presenting the plan, review:

| Check | Description |
|---|---|
| **Spec coverage** | Does every requirement in the spec have a corresponding task? |
| **Placeholder scan** | Search for "TODO", "TBD", "FIXME", "implement", "later" in the plan |
| **Type consistency** | Do types in earlier tasks match types in later tasks? |
| **Completeness** | Are lint, typecheck, and test commands in every applicable task? |
| **Granularity** | Is every task 2-5 minutes? If not, split. |

## Execution Handoff

After writing the plan, offer the user two options:

> This plan has [N] tasks. Would you like me to:
>
> **A) Subagent-Driven Execution (recommended)** — dispatch tasks to parallel subagents for faster completion
> **B) Inline Execution** — execute tasks one by one in this session

If the user chooses A, invoke **same-session** (fase-ejecucion). If B, invoke **session-isolated** (fase-ejecucion).
