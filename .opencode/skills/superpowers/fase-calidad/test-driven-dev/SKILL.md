---
name: test-driven-dev
description: Use when implementing any feature or bugfix, before writing implementation code. Write test first, watch it fail, write minimal code to pass.
---

# Test-Driven Development (TDD)

## Overview

Write the test first. Watch it fail. Write minimal code to pass.

**Core principle:** If you didn't watch the test fail, you don't know if it tests the right thing. Violating the letter of the rules is violating the spirit.

## The Iron Law

```
NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST
```

Write code before the test? Delete it. Start over. No exceptions:
- Don't keep it as "reference"
- Don't "adapt" it while writing tests
- Don't look at it
- Delete means delete

**No exceptions.** Every rationalization is the same thing: skipping TDD.

| Rationalization | Reality |
|-----------------|---------|
| "Too simple to test" | Simple code breaks. Test takes 30 seconds. |
| "I'll test after" | Passing immediately proves nothing. |
| "Tests after achieve same goals" | After = "what does this do?" Before = "what should this do?" |
| "Already manually tested" | Ad-hoc ≠ systematic. No record, can't re-run. |
| "Deleting X hours is wasteful" | Sunk cost fallacy. Keeping unverified code is debt. |
| "Keep as reference, write tests first" | You'll adapt. That's testing after. |
| "Need to explore first" | Fine. Throw away exploration. Start with TDD. |
| "TDD will slow me down" | TDD is faster than debugging. |
| "Existing code has no tests" | You're improving it. Add tests now. |

## Red-Green-Refactor

### RED — Write Failing Test

Write one minimal test. Clear name. Tests real behavior (no mocks unless unavoidable).

**Verify RED — MANDATORY:** Run the test. Confirm it fails (not errors). Failure must be from feature missing, not typos. Test passes? You're testing existing behavior. Fix test.

### GREEN — Minimal Code

Write simplest code to pass the test. Don't generalize. Don't add features. Don't refactor. Just enough to pass.

**Verify GREEN — MANDATORY:** Run the test. Confirm it passes. Other tests still pass. Output pristine. Test fails? Fix code, not test.

### REFACTOR — Clean Up

After green only: remove duplication, improve names, extract helpers. Keep tests green. Don't add behavior.

### Repeat

Next failing test for next feature or bug.

## Good vs Bad Tests

| Quality | Good | Bad |
|---------|------|-----|
| **Minimal** | One thing. "and" in name? Split it. | `test('validates email and domain and whitespace')` |
| **Clear** | Name describes behavior | `test('test1')` |
| **Shows intent** | Demonstrates desired API | Obscures what code should do |

## Common Rationalizations

| Excuse | Reality |
|--------|---------|
| "TDD is dogmatic, I'm pragmatic" | TDD IS pragmatic. Finds bugs before commit. Prevents regressions. Enables refactoring. |
| "Tests after = same thing" | No. After answers "what does this do?" Before answers "what should this do?" |
| "Manual test faster" | Manual doesn't prove edge cases. You'll re-test every change. |
| "30 min of tests after = TDD" | No. You get coverage, lose proof tests work. |

## Red Flags — STOP and Start Over

- Code before test / test after implementation / test passes immediately
- Can't explain why test failed / tests added "later"
- Rationalizing "just this once"
- "I already manually tested it"
- "Keep as reference" or "adapt existing code"
- "Already spent X hours, deleting is wasteful"
- "This is different because..."

**All of these mean: Delete code. Start over with TDD.**
