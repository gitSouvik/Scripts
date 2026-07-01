<div align="center">

```
  ██████╗██████╗ ██╗  ██╗
 ██╔════╝██╔══██╗╚██╗██╔╝
 ██║     ██████╔╝ ╚███╔╝ 
 ██║     ██╔═══╝  ██╔██╗ 
 ╚██████╗██║     ██╔╝ ██╗
  ╚═════╝╚═╝     ╚═╝  ╚═╝
```

**A fast, native Terminal UI for competitive programming.**

*Built with Go · Powered by Bubble Tea · No dependencies to install*

[![Go](https://img.shields.io/badge/go-1.21%2B-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![C++](https://img.shields.io/badge/C%2B%2B-20-00599C?style=flat-square&logo=cplusplus&logoColor=white)](https://isocpp.org/)
[![Platform](https://img.shields.io/badge/platform-macOS%20%2F%20Linux-lightgrey?style=flat-square)]()
[![Contributing](https://img.shields.io/badge/docs-CONTRIBUTING.md-blue?style=flat-square)](./CONTRIBUTING.md)

</div>

---

## What is this?

`cpx` is a single binary that replaces a pile of bash scripts and python fetchers with a clean, keyboard-driven terminal UI. It handles the full competitive programming workflow end-to-end:

- **Fetch** problems from Competitive Companion (auto-starts on launch)
- **Run** your solution against all sample test cases with pass/fail output
- **Debug** with custom hand-typed or pasted test input
- **Create** single files or ranges (`a.cpp` → `e.cpp`) in one keystroke

---

## Quickstart

### Prerequisites

| Tool | Notes |
|---|---|
| **Go 1.21+** | `brew install go` |
| **g++-15** | `brew install gcc` or use your distro's package |
| **Competitive Companion** | [Chrome](https://chrome.google.com/webstore/detail/competitive-companion/cjnmckjndlpiamhfimnnjmnckgghkjbl) / [Firefox](https://addons.mozilla.org/en-US/firefox/addon/competitive-companion/) |

### Build

```bash
git clone https://github.com/<you>/cpx.git
cd cpx
cd cpx-src && go build -o ../cpx . && cd ..
```

### Run

Navigate to your contest directory and launch:

```bash
cd ~/contest
cpx
```

> **Tip**: Add it to your PATH for global access:
> ```bash
> echo 'export PATH="$PATH:/path/to/cpx-dir"' >> ~/.zshrc && source ~/.zshrc
> ```

---

## UI Layout

```
                        00:00:00        cpx · code · test · debug          ← Timer (extreme left) & Title (centered)
  prob :  [A]  B   C   D                                                   ← Problem row
  cmds :  [?] help   [r] run   [x] tests   [e] +tests   [s] snip   [m] more   [c] clear   [q] quit

  Facts! There are 10 types of people in the world...                       ← Dynamic Gemini facts
  ─────────────────────────────────────────────────────────────────────

  [run] compiling A.cpp...
  ✔ compiled successfully (1200ms)
  ─────────────────────────────────────
  ✔ Test 1 Passed (4ms)
  ─────────────────────────────────────
  ✘ Test 2 Failed (3ms)
  │ Output:
  │ 5
  │ Expected:
  │ 4
```

---

## Keyboard Shortcuts

### Navigation (always available)
| Key | Action |
|---|---|
| `←` / `→` | Move between problems |
| Mouse Click | Click on problem letter to jump focus directly to it |

### Commands (Shortcut row focused or anywhere on main screen)
| Key | Action |
|---|---|
| `r` | Run selected problem against all `.in`/`.out` test cases (includes side-by-side token diff) |
| `x` | Open interactive test cases editor |
| `e` | Open custom test editor (`+tests`) — split-pane (input on top, logs on bottom), `ctrl+r` to run, `ctrl+s`/`ctrl+a` to save |
| `s` | Open code snippet injection menu (press `n` here to open Snippet Creator) |
| `m` | Open "More Options Guide" for hidden options |
| `c` | Clear output log |
| `q` | Quit |

### Hidden Options (Shown in `[m] more` guide)
| Key | Action |
|---|---|
| `t` | Start or stop the contest countdown timer |
| `n` | New file dialog — create a single file or a range (`a`→`e`) |
| `ctrl+p` | Open Gemini API key input dialog |
| `ctrl+h` | Toggle facts visibility |

---

## Workflow

### 1. Fetch a Problem
Start `cpx`. The fetch server is already running on port `54321`. Click the **+** icon in your Competitive Companion browser extension. The problem files are created silently and appear in the problem row instantly.

Files created:
- `A.cpp` — pre-filled from `template.cpp`
- `A-1.in`, `A-1.out`, `A-2.in`, `A-2.out`, … — all sample test cases

### 2. Write Your Solution & Inject Snippets
Open `A.cpp` in your editor and write your solution. 
- Need a boilerplate algorithm? Press `s` in `cpx` to open snippets, select one, and hit `enter` to inject it.
- Want to create a new snippet? Press `s` -> `n` to open the **Snippet Creator**, type the name, paste the code, and hit `ctrl+s`.

### 3. Test It
Press `r` in `cpx`. Your solution is compiled and run against every sample test. Passing tests show green `✔`, failing tests show red `✘` with output and expected output. A token-level **Smart Diff** highlights exact mismatched tokens in red.

### 4. Custom Tests (`+tests`)
Press `e` to open the split-pane custom test editor. Type or paste your input on the top half. Press `ctrl+r` to run it—the output appears immediately on the bottom half without closing the editor. You can tweak and rerun instantly. Press `ctrl+s` to save it as a permanent test case, or `esc` to close.

### 5. Debug Run
Press `d` to select a test case number, or `0` to run all test cases without comparing to expected output. Useful for viewing debug statements.

### 6. Contest Timer & Facts
- Start a contest countdown timer by pressing `t` (from the `m` menu) and setting a duration (e.g. `2h` or `120m`).
- A new programming fact is fetched every 5 minutes from Gemini. Press `ctrl+h` to hide/show the facts bar.


---

## Project Structure

```
cpx/
├── cpx                     ← Compiled binary (run this)
├── cpx-src/                ← Go source code
│   ├── main.go             ← Entry point
│   ├── model.go            ← State, key handling, Bubble Tea model
│   ├── view.go             ← Terminal UI rendering
│   ├── exec.go             ← Compile, run, test logic
│   ├── fetch.go            ← Competitive Companion HTTP server
│   ├── go.mod
│   └── go.sum
│
├── template.cpp            ← C++20 solution template (loaded for each problem)
├── debug.h                 ← Basic dbg() macro header
├── debug++.h               ← Extended dbg() for tuples, queues, 2D arrays, etc.
│
├── .clangd                 ← clangd config for IDE autocomplete
├── .ide_includes/bits/     ← Dummy stdc++.h for clangd on macOS
│
└── docs/
    ├── SETUP.md            ← Full setup guide (compiler, clangd, PATH)
    └── WORKFLOW.md         ← End-to-end contest workflow
```

---

## Debug Macros

The `debug.h` / `debug++.h` headers provide a `dbg(...)` macro that pretty-prints any variable to `stderr`. It is **active only when compiled with `-DWOOF_`** (set by cpx automatically), so it is silently disabled on the judge.

```cpp
vector<int> v = {1, 2, 3};
dbg(v); // → [v] = [1, 2, 3]

pair<int,int> p = {4, 5};
dbg(p); // → [p] = (4, 5)
```

To install on macOS (so the header is picked up system-wide for your compiler):

```bash
sudo cp debug++.h /usr/local/include/bits/debug.h
```

---

## License

Personal toolchain — fork and adapt freely.
