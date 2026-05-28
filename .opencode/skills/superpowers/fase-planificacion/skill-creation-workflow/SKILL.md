---
name: skill-creation-workflow
description: Planning-first workflow for creating or modifying skills. Use when creating new skills, editing existing skills, or before modifying any skill file. Must present a plan and get user approval before execution.
---

# Skill Creation Workflow

## Overview

Plan first, then execute. Before touching any skill file, present a complete plan of files to create/modify, ask the user for doubts or changes, and only proceed after explicit confirmation.

**Core principle:** No skill file is created or modified without an approved plan first.

**REQUIRED SUB-SKILL:** Use `writing-skills` for the actual content structure (SKILL.md format, frontmatter, CSO rules). This skill adds the planning-and-approval layer on top.

## When to Use

- Creating a brand new skill from scratch
- Editing an existing skill's content, name, or description
- Deleting a skill
- Adding supporting files to a skill directory
- Any task involving `.opencode/skills/` paths

Do NOT use for: reading a skill, loading a skill via the `skill` tool, or using a skill's instructions during implementation.

## The Process

### Phase 1: Plan

Read the context to understand what's needed, then produce a plan listing every file operation.

**Sub-steps:**

1. **Read existing skills** in the relevant area (e.g., `fase-calidad/`, `python/`) to understand conventions — naming, frontmatter style, body sections, directory structure.
2. **Analyze requirements** — what specific problem does the new skill solve? What existing skills overlap? What should the description say for discoverability?
3. **List all file operations** — every file to create, modify, or delete. For each, note:
   - Path (absolute or relative to `.opencode/skills/`)
   - Operation: `Create`, `Modify`, or `Delete`
   - Brief content summary (3-5 lines max)
4. **Check dependencies** — does any file depend on another? Must files be created in order?

### Phase 2: Present

Present the plan to the user using the template below. Then ask:

> Do you have any doubts or changes before I proceed?

Wait for the user's response. If they have doubts, address them. If they request changes, update the plan and re-present. Do not proceed to Phase 3 until the user explicitly confirms.

### Phase 3: Execute

After user approval, execute each file operation:

1. Create directories first, then files
2. Follow the `writing-skills` skill for SKILL.md structure (frontmatter, body template, CSO rules)
3. After creating/modifying each file, verify it was written correctly (check path, size, frontmatter integrity)
4. For modifications, re-read the file before editing (the edit tool requires it)

If blocked during execution, stop and report. Do not make unplanned changes.

## Plan Template

Present plans in this format:

```
## Plan: <Skill Name>

### Files

| Path | Operation | Summary |
|------|-----------|---------|
| `.opencode/skills/<category>/<name>/SKILL.md` | Create | Frontmatter + body covering X, Y, Z |
| `.opencode/skills/<category>/<name>/example.py` | Create | Supporting example script |

### Dependencies

- None / List any ordering requirements

### Questions / Notes

- Any open questions or design decisions noted here
```

## Common Mistakes

- **Skipping the plan** — the most common failure. The plan must exist and be approved before any file operation.
- **Planning but not listing every file** — if you discover a file during execution that wasn't in the plan, stop and re-present.
- **Vague summaries** — "add content" is not a plan. Summarize the key sections: "Covers input validation, retry policies, and resource cleanup patterns."
- **Ignoring `writing-skills`** — this skill handles the workflow, `writing-skills` handles the content structure. Both are needed.
- **Modifying files not in the plan** — scope creep. If execution reveals a needed change, add it to the plan and re-confirm.
