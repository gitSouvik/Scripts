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
  cpx — Competitive Programming Tools
  /path/to/contest

│ [A]  B   C   D                          ← Problem row (focused)
  Type :  [r] run  [e] custom  [d] debug  [c] clear  [n] new  [g] gen  [q] quit

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

The `│` bar on the left shows which row is focused. Use `↑`/`↓` to switch rows.

---

## Keyboard Shortcuts

### Navigation (always available)
| Key | Action |
|---|---|
| `↑` / `↓` | Switch focus between Problem row and Shortcut row |
| `←` / `→` | Move between problems (when Problem row is focused) |

### Commands (Shortcut row focused or anywhere on main screen)
| Key | Action |
|---|---|
| `r` | Run selected problem against all `.in`/`.out` test cases |
| `e` | Open custom test editor (paste/type, then `ctrl+r` to run) |
| `d` | Debug run — opens dialog to pick a specific test case number |
| `c` | Clear output log |
| `n` | New file dialog — create a single file or a range (`a`→`e`) |
| `q` | Quit |

---

## Workflow

### 1. Fetch a Problem
Start `cpx`. The fetch server is already running on port `54321`. Click the **+** icon in your Competitive Companion browser extension and the problem files are created silently and appear in the problem row instantly.

Files created:
- `A.cpp` — pre-filled from `template.cpp`
- `A-1.in`, `A-1.out`, `A-2.in`, `A-2.out`, … — all sample test cases

### 2. Write Your Solution
Open `A.cpp` in your editor and write your solution.

### 3. Test It
Press `r` in cpx. Your solution is compiled and run against every sample test. Passing tests show green `✔`, failing tests show red `✘` with the output and expected output side by side.

### 4. Custom Tests
Press `e` to open the custom test editor. Paste or type your own input. Press `ctrl+r` to run it. Press `esc` to cancel.

### 5. Debug Run
Press `d`, then enter a test case number (`1`, `2`, etc.) or `0` to run all test cases without comparison. Useful for inspecting debug output from `dbg()`.

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
