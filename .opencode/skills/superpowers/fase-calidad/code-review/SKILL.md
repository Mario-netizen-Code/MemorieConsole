---
name: code-review
description: Use when completing tasks, implementing features, or before merging to verify work meets requirements. Guides requesting review from subagents AND receiving/responding to feedback.
---

# Code Review

## Overview

This skill covers both sides of the code review process: requesting a structured review of your own work (dispatch to a reviewer subagent) and receiving/responding to review feedback on someone else's code. Both sides share the principle that review is a technical verification, not a social ritual.

## Requesting Review

### When to Request

Request a code review when:
- All plan tasks are implemented and pass local verification (lint, typecheck, tests).
- Before merging any branch.
- When unsure about design decisions or implementation quality.
- After completing a non-trivial feature or refactor.

### How to Dispatch a Reviewer

1. Get the current git SHAs for the branch and base:
   ```
   git rev-parse HEAD
   git merge-base HEAD <base-branch>
   ```
2. Dispatch a reviewer subagent with a structured request. Use this template:

   ```
   Please review changes between <base-sha> and <head-sha>.

   Context: <what this PR/change does>
   Focus areas: <specific concerns, if any>

   Check for:
   - Correctness: does the code do what it claims?
   - Test coverage: are there adequate tests?
   - Edge cases: error handling, empty states, boundary conditions.
   - Project conventions: naming, imports, patterns.
   - Security: no secrets, injection risks, or unsafe patterns.
   ```

3. Collect the reviewer's output.

### Fix Priority

When the reviewer returns feedback, prioritize fixes:

| Priority | Meaning | Action |
|----------|---------|--------|
| **Critical** | Bug, security issue, or correctness problem | Fix immediately. Do not proceed until resolved. |
| **Important** | Design concern, missing coverage, or significant style deviation | Fix before proceeding to next task or merging. |
| **Minor** | Nitpick, naming suggestion, optional refactor | Log for later. Can be addressed in follow-up. |

Never dismiss a Critical or Important finding without addressing it. Minor items can be batched.

## Receiving Review

When you receive a review (from a subagent, peer, or system), follow this process:

### 1. Read Fully

Read the entire review before responding. Do not react to individual comments in isolation.

### 2. Understand

For each comment, make sure you understand the underlying concern. Ask clarifying questions if needed.

### 3. Verify

Before agreeing or disagreeing, verify the claim:
- Is the cited problem real? Run tests, check types, inspect the code.
- Does the reviewer's suggested fix actually solve the problem?
- Does the suggestion introduce new issues?

### 4. Evaluate (YAGNI Check)

Apply YAGNI (You Ain't Gonna Need It):
- Does the suggestion add complexity for a hypothetical future need?
- Is the current code correct and sufficient as-is?
- Would the change make the code harder to read or maintain?

If the suggestion fails YAGNI, do not implement it. Document your reasoning.

### 5. Respond

- If the reviewer is correct: acknowledge and fix.
- If the reviewer is partially correct: explain the nuance and adjust as needed.
- If the reviewer is wrong: push back with technical reasoning (see below).

### 6. Implement

Implement feedback one item at a time. Do not batch all changes into a single blind pass. After each change, re-verify that tests and types still pass.

## Forbidden Responses

**No performative agreement.** Never agree with a reviewer just to be agreeable or move quickly. If you do not understand or do not agree, say so.

**No blind implementation.** Never implement a reviewer suggestion without understanding why it is correct. Verify first, then act.

**No silent dismissal.** If you decide not to address a comment, explain your reasoning explicitly. Do not mark it as resolved without a response.

## When to Push Back

Push back when:

- The reviewer misunderstood the code or its intent.
- The suggestion introduces a bug, regression, or unnecessary complexity.
- The suggestion violates project conventions or established patterns.
- The suggestion is a personal preference rather than an objective improvement.
- The suggestion would require changes outside the scope of the task.

Always push back with reasoning, not defensiveness. Cite code, tests, or documentation to support your position. If the reviewer persists and you remain unconvinced, escalate rather than capitulate.

## Implementation Order

When implementing review feedback:

1. Fix all Critical items first, one at a time.
2. Fix all Important items, one at a time.
3. Address Minor items if time permits, or batch them into a follow-up.
4. After each fix, re-run lint, typecheck, and relevant tests.
5. After all fixes, request re-review if the changes were significant.

Do not skip verification between fixes. A change that fixes one thing may break another.
