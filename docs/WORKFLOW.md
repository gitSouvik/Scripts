# Contest Workflow

A step-by-step guide for using `cpx` during a competitive programming contest.

---

## Before the Contest

Navigate to your working directory and start cpx:

```bash
cd ~/contest    # wherever you keep your solutions
cpx
```

The fetch server starts automatically on port `54321`. You are ready.

---

## During the Contest

### Step 1 — Fetch the Problem

Open the problem page in your browser. Click the green **+** in the Competitive Companion extension.

cpx receives the problem silently in the background and:
- Creates `A.cpp` (or the appropriate letter) pre-filled from `template.cpp`
- Creates all sample test cases as `A-1.in` / `A-1.out`, `A-2.in` / `A-2.out`, …
- The problem appears immediately in the top row of the TUI

### Step 2 — Write Your Solution

Open `A.cpp` in your editor and code your solution.

Use `dbg(x)` anywhere to inspect variables during local testing — it prints to stderr and is compiled out on the judge automatically.

### Step 3 — Run the Samples

In cpx, make sure `A` is selected (use `←`/`→`), then press **`r`**.

cpx will:
1. Compile with `g++-15 -O2 -std=c++20 -DWOOF_`
2. Run against each `A-N.in` file
3. Compare output to `A-N.out`
4. Show `✔ Test N Passed` or `✘ Test N Failed` with output/expected

### Step 4 — Debug a Failing Test

**Option A — Debug a specific test case:**
Press **`d`**, enter the test number (e.g. `2`), press Enter. cpx runs only that test and shows your full output including any `dbg()` lines.

**Option B — Custom input:**
Press **`e`** to open the custom test editor. Type or paste your own test case. Press `ctrl+r` to run. Press `esc` to cancel.

### Step 5 — Submit

When all samples pass, copy your solution and submit to the judge.

---

## Creating Files Manually

If you want to create solution files without fetching from Competitive Companion:

Press **`n`** in cpx to open the New File dialog. Two options:

1. **Single File** — Enter a name (e.g. `E`) and extension (`cpp`) → creates `E.cpp`
2. **Range** — Enter an end character (e.g. `e`) and extension (`cpp`) → creates `a.cpp`, `b.cpp`, `c.cpp`, `d.cpp`, `e.cpp`

---

## Keyboard Reference

| Key | Action |
|---|---|
| `↑` / `↓` | Switch focus between Problem row and Shortcut row |
| `←` / `→` | Move between problems |
| `r` | Run against all sample test cases |
| `e` | Open custom test input editor |
| `d` | Debug run (pick a test case number) |
| `c` | Clear output log |
| `n` | Create new file(s) |
| `q` | Quit |
| Mouse wheel | Scroll through output logs |
