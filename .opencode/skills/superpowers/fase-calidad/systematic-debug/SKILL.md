---
name: systematic-debug
description: Use when encountering any bug, test failure, or unexpected behavior, before proposing fixes. Finds root cause before applying any fix.
---

# Systematic Debugging

## Overview

Random fixes waste time and create new bugs. Quick patches mask underlying issues.

**Core principle:** ALWAYS find root cause before attempting fixes. Symptom fixes are failure. Violating the letter of this process is violating the spirit.

## The Iron Law

```
NO FIXES WITHOUT ROOT CAUSE INVESTIGATION FIRST
```

If you haven't completed Phase 1, you cannot propose fixes.

## The Four Phases

You MUST complete each phase before proceeding to the next.

### Phase 1: Root Cause Investigation

1. **Read Error Messages** — Don't skip past errors. Read stack traces completely. Note line numbers, file paths, error codes.
2. **Reproduce Consistently** — Exact steps? Happens every time? Not reproducible? Gather more data, don't guess.
3. **Check Recent Changes** — Git diff, recent commits, new dependencies, config changes, environmental differences.
4. **Gather Evidence in Multi-Component Systems** — For each component boundary: log what enters/exits, verify environment propagation, check state at each layer. Run once to gather evidence showing WHERE it breaks, then analyze.
5. **Trace Data Flow** — Where does bad value originate? What called this with bad value? Keep tracing up until you find the source. Fix at source, not at symptom.

### Phase 2: Pattern Analysis

1. **Find Working Examples** — Locate similar working code in same codebase.
2. **Compare Against References** — Read reference implementation COMPLETELY. Don't skim.
3. **Identify Differences** — List every difference between working and broken. Don't assume "that can't matter."
4. **Understand Dependencies** — Other components needed? Settings, config, environment? Assumptions?

### Phase 3: Hypothesis and Testing

1. **Form Single Hypothesis** — "I think X is root cause because Y." Be specific.
2. **Test Minimally** — Smallest possible change. One variable at a time. Don't fix multiple things at once.
3. **Verify Before Continuing** — Worked? Phase 4. Didn't work? New hypothesis. DON'T stack fixes.
4. **When You Don't Know** — Say it. Don't pretend. Ask for help. Research more.

### Phase 4: Implementation

1. **Create Failing Test** — Simplest reproduction. Automated if possible. MUST have before fixing.
2. **Implement Single Fix** — Address root cause. ONE change. No "while I'm here" improvements.
3. **Verify Fix** — Test passes? No other tests broken? Issue actually resolved?
4. **If Fix Doesn't Work** — STOP. Count attempts. If < 3: return to Phase 1. If ≥ 3: do NOT attempt fix #4 — question architecture.

## Red Flags — Stop Signals

- "Quick fix for now, investigate later"
- "Just try changing X and see if it works"
- "Add multiple changes, run tests"
- "Skip the test, I'll manually verify"
- "I don't fully understand but this might work"
- **"One more fix attempt" (when already tried 2+)**
- **Each fix reveals new problem in different place**

**All of these mean: STOP. Return to Phase 1.**

## Common Rationalizations

| Excuse | Reality |
|--------|---------|
| "Issue is simple, don't need process" | Process is fast for simple bugs too. |
| "Emergency, no time for process" | Systematic is FASTER than guess-and-check. |
| "Just try this first, then investigate" | First fix sets the pattern. Do it right. |
| "Multiple fixes at once saves time" | Can't isolate what worked. Causes new bugs. |
| "I see the problem, let me fix it" | Seeing symptoms ≠ understanding root cause. |

## When to Question Architecture

**3+ fixes failed?** Stop and ask:
- Is this pattern fundamentally sound?
- Are we "sticking with it through inertia"?
- Should we refactor architecture vs. continue fixing symptoms?

This is NOT a failed hypothesis — this is a wrong architecture. Discuss with your human partner before attempting more fixes.
