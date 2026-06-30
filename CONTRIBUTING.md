# Contributing & Modifying cpx

This document explains the architecture of cpx and how to make changes, whether you are fixing a bug, adding a command, tweaking the UI, or adapting it for your own workflow.

---

## Architecture Overview

cpx is a single Go module in `cpx-src/`. It uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework, which follows the Elm architecture: **Model → Update → View**.

```
cpx-src/
├── main.go      Entry point. Initialises the Bubble Tea program.
├── model.go     All application state (problems list, screen, dialogs, etc.)
│                and all keyboard / mouse input handling.
├── view.go      Pure rendering — takes model state, returns a string for the terminal.
├── exec.go      Compile and run logic: run, debug, custom test, gen.
└── fetch.go     HTTP server that receives problems from Competitive Companion.
```

---

## Build

```bash
cd cpx-src
go build -o ../cpx .
```

Or run in-place without producing a binary (useful during development):

```bash
cd cpx-src
go run .
```

---

## How to Make Changes

### Change a keyboard shortcut

Open `cpx-src/model.go` and find `func (m model) updateMain`. Each `case "x":` block handles that keypress. Add, remove, or change cases freely.

```go
case "x":
    // your new action here
```

### Add a new screen / dialog

1. Add a constant to the `screen` or `dialogKind` enum at the top of `model.go`.
2. Add a handler function like `updateCustomTest` — called from the `tea.KeyMsg` routing block in `Update`.
3. Add a rendering branch in `view.go` inside `View()`.

### Change compile flags

Open `cpx-src/exec.go` and find `compileFile`. The compile command is:

```go
cmd := exec.Command("g++-15", "-O2", "-std=c++20", "-DWOOF_", "-o", name, cppFile)
```

Change `g++-15` to match your installed compiler version, or add/remove flags as needed.

### Change the fetch server port

Open `cpx-src/fetch.go`. The port is hardcoded in `launchFetch`:

```go
fetchServer = &http.Server{Addr: ":54321", Handler: mux}
```

Change `54321` to any free port. Make sure to match it in the Competitive Companion extension settings.

### Modify the template loaded for each fetched problem

Edit `template.cpp` in the repo root. This file is read by `fetch.go` and written into each new problem's `.cpp` file. The fetch code looks for `template.cpp` relative to the directory where cpx is launched.

### Change the UI layout

All rendering is in `cpx-src/view.go` inside `func (m model) View() string`. It builds a `strings.Builder` from top to bottom. ANSI color codes are used directly (e.g. `\033[1;32m` = bold green, `\033[2;37m` = dim white).

To change colors, adjust the ANSI codes in the `switch l.kind` block in `view.go`.

---

## About `.ide_includes` and `.clangd`

> **These files are only needed on macOS with VS Code (or any clangd-powered editor).** You can delete them on Linux — it already has `bits/stdc++.h` natively.

### Why they exist

`#include <bits/stdc++.h>` is a GCC convenience header that includes everything. It exists in Homebrew GCC (`/usr/local/include/bits/stdc++.h`) and works fine at **compile time** on macOS.

However, **clangd** (the C++ language server that powers autocomplete in VS Code) uses Apple's clang include paths, which do **not** have `bits/stdc++.h`. This causes clangd to show a false "file not found" error and break autocomplete — even though the code compiles perfectly.

The fix:
- `.ide_includes/bits/stdc++.h` — a dummy header that includes all standard headers individually, satisfying clangd's include resolution.
- `.clangd` — tells clangd to add `.ide_includes` to its include search path, and to compile with `-std=c++20 -DWOOF_` so macros resolve correctly.

### If you are on Linux

Delete both:

```bash
rm -rf .ide_includes .clangd
```

### If you use a different editor (not clangd-based)

Same — these files are irrelevant. Delete them or ignore them.

---

## Debug Headers (`debug.h` / `debug++.h`)

The `dbg(...)` macro prints any variable to `stderr` with its name and value. It is active only when compiled with `-DWOOF_` (which cpx sets automatically). It is a no-op on the judge.

### Install system-wide (macOS / Linux)

```bash
# find where your g++ looks for includes
g++-15 -v 2>&1 | grep "include"

# copy to that path
sudo cp debug++.h /usr/local/include/bits/debug.h
```

### Usage in code

```cpp
#ifdef WOOF_
#include <bits/debug.h>
#endif

// then anywhere in your solution:
dbg(x);                  // → [x] = 42
dbg(v);                  // → [v] = [1, 2, 3]
dbg(a, b, c);            // → [a, b, c] = 1, 2, 3
dbg(make_pair(1, 2));    // → [(1, 2)] = (1, 2)
```

`debug++.h` supports: `int`, `string`, `vector`, `pair`, `tuple`, `map`, `set`, `stack`, `queue`, `priority_queue`, 2D vectors, C-arrays, and more.

---

## Dependency Management

cpx uses Go modules. Dependencies are pinned in `cpx-src/go.mod` and `cpx-src/go.sum`.

To update all dependencies:

```bash
cd cpx-src
go get -u ./...
go mod tidy
```

Current dependencies:
- [`github.com/charmbracelet/bubbletea`](https://github.com/charmbracelet/bubbletea) — TUI framework
- [`github.com/charmbracelet/bubbles`](https://github.com/charmbracelet/bubbles) — textinput and textarea components
- [`github.com/charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) — terminal styling
