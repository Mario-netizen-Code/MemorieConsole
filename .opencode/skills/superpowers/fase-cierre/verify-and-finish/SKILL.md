---
name: verify-and-finish
description: Use when about to claim work is complete, fixed, or passing, before committing or creating PRs. Guides verification then merge/PR/discard decision.
---

# Verify and Finish

## Overview

Claiming work is complete without verification is dishonesty, not efficiency. Once verified, decide how to integrate: merge locally, push+PR, keep as-is, or discard.

**Core principle:** Evidence before claims. Then execute the chosen finish path.

## The Gate Function

```
BEFORE claiming any status:
1. IDENTIFY: What command proves this claim?
2. RUN: Execute the FULL command (fresh, complete)
3. READ: Full output, check exit code, count failures
4. VERIFY: Does output confirm the claim? If NO: state actual status. If YES: claim WITH evidence.
5. ONLY THEN: Make the claim

Skip any step = lying, not verifying
```

## Common Failures

| Claim | Requires | Not Sufficient |
|-------|----------|----------------|
| Tests pass | Output: 0 failures | Previous run, "should pass" |
| Linter clean | Output: 0 errors | Partial check, extrapolation |
| Build succeeds | Exit 0 | Linter passing |
| Bug fixed | Original symptom passes | Code changed, assumed fixed |
| Regression test works | Red-green verified | Test passes once |
| Agent completed | VCS diff shows changes | Agent reports "success" |

## The Finish Menu

**Announce:** "Running verify-and-finish skill."

**Step 1:** Verify tests pass. If failing, stop — must fix first.

**Step 2:** Detect environment:
- `GIT_DIR == GIT_COMMON` (normal repo) — standard 4 options
- `GIT_DIR != GIT_COMMON`, named branch — standard 4 options, provenance cleanup
- `GIT_DIR != GIT_COMMON`, detached HEAD — reduced 3 options (no merge)

**Step 3:** Determine base branch (`git merge-base HEAD main || git merge-base HEAD master`).

**Step 4:** Present options:

```
1. Merge back to <base> locally
2. Push and create a Pull Request
3. Keep the branch as-is
4. Discard this work
```

## Option Details

**1. Merge locally:** Switch to base, pull, merge. Verify tests on merged result. Clean up worktree, delete branch.

**2. Push + PR:** `git push -u origin <branch>`, `gh pr create`. Do NOT clean up worktree — user needs it for PR iteration.

**3. Keep as-is:** Report branch and path. Preserve worktree.

**4. Discard:** Requires typed "discard" confirmation. Clean up worktree, force-delete branch.

## Cleanup

Only for options 1 and 4. Check provenance:
- Under `.worktrees/`, `worktrees/`, or `~/.config/superpowers/worktrees/` → we own cleanup: `git worktree remove` + `git worktree prune`
- Otherwise → harness owns it. Do NOT remove.
- Always `cd` to main repo root before removal. Delete branch after worktree removal.
