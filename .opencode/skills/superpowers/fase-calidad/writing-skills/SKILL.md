---
name: writing-skills
description: Use when creating new skills, editing existing skills, or verifying skills work before deployment. Applies TDD to documentation.
---

# Writing Skills

## Overview

**Writing skills IS TDD for documentation.** Write pressure scenarios (tests), watch the agent fail without the skill (RED), write the skill (GREEN), refactor to close loopholes.

**Core principle:** If you didn't watch an agent fail without the skill, you don't know if the skill teaches the right thing.

**REQUIRED BACKGROUND:** You MUST use `test-driven-dev` before using this skill. This adapts TDD to documentation.

## SKILL.md Structure

### Frontmatter (YAML)

```yaml
---
name: my-skill-name
description: Use when [triggering conditions]. Never summarize workflow — only describe when to use.
---
```

- `name`: letters, numbers, hyphens only. Max 64 chars.
- `description`: Start with "Use when...". Max 500 chars. Third person. **NEVER** summarize the skill's workflow.

### Body Template

```markdown
# Skill Name

## Overview
Core principle in 1-2 sentences.

## When to Use
Bullet list with symptoms and use cases. When NOT to use.

## Core Pattern
Before/after comparison or code examples.

## Quick Reference
Table for common operations.

## Common Mistakes
What goes wrong + fixes.
```

## Claude Search Optimization (CSO)

The description is what Claude reads to decide whether to load the skill. **It must answer "Should I read this right now?" using only triggering conditions.**

- ✅ `description: Use when tests have race conditions or pass/fail inconsistently`
- ❌ `description: Use for async testing with condition-based-waiting pattern that replaces timeouts`

The second form causes Claude to follow the description instead of reading the full skill. A summary of workflow in the description creates a shortcut Claude will take.

### Keyword Coverage

Use words Claude would search for: error messages, symptoms, synonyms, tool names.

### Cross-Referencing Other Skills

```markdown
**REQUIRED SUB-SKILL:** Use `test-driven-dev`
```

Never use `@` links — they force-load files and burn context.

## File Organization

```
skills/
  skill-name/
    SKILL.md              # Main reference (required)
    supporting-file.*     # Only if needed
```

Keep references inline unless they exceed 100 lines. Use separate files only for heavy API docs or reusable scripts.

## Common Mistakes

- **Narrative storytelling** — don't describe how you solved a problem once. State the pattern generically.
- **Multi-language examples** — one excellent example beats mediocre implementations in 5 languages.
- **Code in flowcharts** — can't copy-paste. Use markdown code blocks instead.
- **No rationalization table** — discipline skills need explicit counters for common excuses.
