---
name: brainstorming
description: Use for creative work before implementation. Explores user intent, requirements and design through guided discovery and structured specs.
---

# Brainstorming

## Overview

Transforms vague ideas or feature requests into concrete, approved designs before any code is written. Uses structured dialogue to explore intent, surface constraints, produce a written spec, and hand off to the writing-plans skill for implementation planning.

## When to Use

- Any creative work: new features, components, functionality, or behavioral changes
- Requirements are ambiguous, incomplete, or conflicting
- Before touching any implementation code
- When you need user approval on a design direction before proceeding

Do NOT use for bugs, test-only tasks, or purely mechanical refactors where the design is already clear.

## The Process

### 1. Explore — gather context

Investigate existing code, docs, recent commits, and any related files to understand the current state. Look for patterns, conventions, existing similar features, and architectural constraints.

Do NOT propose solutions yet. Only collect information.

### 2. Ask — clarify intent, one question at a time

Ask the user clarifying questions **one at a time**. Wait for an answer before asking the next.

Prefer multiple-choice or yes/no formats:

```
Would you prefer:
A) A CLI flag to control behavior
B) Automatic detection based on file type
C) Both
```

Keep asking until you have enough context to propose approaches. Typically 2-4 rounds.

### 3. Propose — 2-3 approaches with trade-offs

Present 2-3 distinct approaches. For each include:

- **How it works** (1-2 sentences)
- **Pros** (short)
- **Cons** (short)
- **Effort estimate** (Small / Medium / Large)

End with a recommendation and a clear ask:

```
I recommend Approach B because [reason]. Shall I proceed with this?
```

### 4. Present Design — structured sections for approval

Once the approach is chosen, present a structured design covering:

- **Goal**: one-sentence summary
- **Changes needed**: files to touch, public API surface
- **Behavior**: how it works from the user's perspective
- **Edge cases**: what we will handle and what we won't
- **Out of scope**: explicitly stated

Ask the user to approve or adjust before writing the spec.

### 5. Write Spec — design document to disk

Write a design document to `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`.

The spec contains:

```markdown
# Design: <Title>

## Goal

## Approach

## Changes

## API / Behavior

## Edge Cases

## Out of Scope

## Open Questions
```

Create the directory if it does not exist.

### 6. Self-Review — before presenting to user

Review your own spec for:

| Check | Description |
|---|---|
| **Placeholders** | Any "TBD", "TODO", or empty sections? |
| **Contradictions** | Do any sections conflict with each other? |
| **Ambiguity** | Could a reader misinterpret any statement? |
| **Scope** | Are boundaries clearly stated (in scope vs out)? |
| **Completeness** | Are all decisions from the proposal reflected? |

Fix any issues found before presenting to the user.

### 7. User Review — get sign-off

Present the written spec for user review. Ask explicitly:

> Please review the design spec at `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md`. Do you approve this design?

Wait for explicit approval before proceeding.

### 8. Transition to Plan — invoke writing-plans

Once approved, invoke the writing-plans skill to create an implementation plan:

```
Proceed with the writing-plans skill to create an implementation plan from this spec.
```

## Visual Companion

When useful, supplement the design with visual aids in the spec:

```
[ASCII diagram]

Table: File A | File B | Interaction
```

For complex architectures, offer to describe a diagram in enough detail that it could be rendered as a Mermaid or ASCII flowchart. The goal is to clarify relationships, not pixel-perfect mockups.

## Key Principles

- **One question at a time**: never batch questions — each answer informs the next
- **Multiple choice preferred**: guides the user without overwhelming them
- **No implementation without approval**: the spec is a contract; once approved it must be followed
- **HARD-GATE**: Do NOT write any implementation code until the design is formally approved by the user
- **Specs are committed**: the design doc becomes part of the repository for future reference
- **Terminal state**: success means the writing-plans skill is invoked; failure means the user rejected the design
