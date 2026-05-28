---
name: using-superpowers
description: Use when starting any conversation. Establishes how to find and use skills — invoke Skill tool before ANY response including clarifying questions.
---

<SUBAGENT-STOP>
If you were dispatched as a subagent, skip this skill.
</SUBAGENT-STOP>

<EXTREMELY-IMPORTANT>
If you think there is even a 1% chance a skill might apply, you ABSOLUTELY MUST invoke the Skill tool. This is not optional.
</EXTREMELY-IMPORTANT>

## The Rule

Invoke relevant skills BEFORE any response or action. Even a 1% chance → invoke. If the skill turns out wrong, don't use it.

## Instruction Priority

1. **User's explicit instructions** (AGENTS.md, direct requests) — highest
2. **Superpowers skills** — override default system where they conflict
3. **Default system prompt** — lowest

## Red Flags — STOP, You're Rationalizing

| Thought | Reality |
|---------|---------|
| "This is just a simple question" | Questions are tasks. Check for skills. |
| "I need more context first" | Skill check comes BEFORE clarifying questions. |
| "Let me explore the codebase first" | Skills tell you HOW to explore. Check first. |
| "This doesn't need a formal skill" | If a skill exists, use it. |
| "The skill is overkill" | Simple things become complex. Use it. |
| "I know what that means" | Knowing the concept ≠ using the skill. Invoke it. |

## Skill Priority

1. **Process skills first** (brainstorming, debugging) — determine HOW to approach
2. **Implementation skills second** — guide execution

## Skill Types

- **Rigid** (TDD, debugging): Follow exactly. Don't adapt away discipline.
- **Flexible** (patterns): Adapt principles to context. The skill tells you which.

## Platform Access

- **opencode:** Use the `skill` tool
- **Claude Code:** Use the `Skill` tool
- **Copilot CLI:** Uses `skill` tool, auto-discovers from plugins
