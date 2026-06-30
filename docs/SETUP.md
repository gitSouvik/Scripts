# Setup Guide

One-time setup to get `cpx` and your C++ environment working.

---

## 1 — Install Prerequisites

### Go (required to build cpx)

```bash
brew install go
```

Verify: `go version` should print `go1.21` or higher.

### g++-15 (C++ compiler)

cpx compiles your solutions with `g++-15`. Install GCC via Homebrew:

```bash
brew install gcc
```

Verify: `g++-15 --version` should print version info.

> If your brew installs a different version (e.g. `g++-14`), update the compile command in `cpx-src/exec.go` — search for `g++-15` and change it.

### Competitive Companion

Install the browser extension:

- [Chrome Web Store](https://chrome.google.com/webstore/detail/competitive-companion/cjnmckjndlpiamhfimnnjmnckgghkjbl)
- [Firefox Add-ons](https://addons.mozilla.org/en-US/firefox/addon/competitive-companion/)

In the extension settings, set the port to **54321**.

---

## 2 — Build cpx

```bash
cd cpx-src
go build -o ../cpx .
cd ..
```

This produces a single `cpx` binary in the project root.

---

## 3 — Add to PATH (optional but recommended)

To run `cpx` from any directory:

```bash
echo 'export PATH="$PATH:/path/to/this/repo"' >> ~/.zshrc
source ~/.zshrc
```

Replace `/path/to/this/repo` with the actual path to this folder.

---

## 4 — IDE Setup (clangd autocomplete)

The `.clangd` config and `.ide_includes/bits/stdc++.h` are provided so that `#include <bits/stdc++.h>` works correctly in VS Code or any clangd-powered editor on macOS (which doesn't ship `bits/stdc++.h` natively).

This is purely for the language server — it does not affect compilation.

No extra steps needed; just ensure clangd is installed via your editor's C++ extension.

---

## 5 — Install Debug Headers (optional)

The `debug++.h` header provides rich `dbg()` output. To make it available globally for your compiler:

```bash
sudo cp debug++.h /usr/local/include/bits/debug.h
```

Or place it wherever your g++ looks for includes. Once installed, `dbg(x)` in your solution will print to stderr when running locally via cpx, and silently does nothing when submitted to a judge.

---

## 6 — Verify

```bash
./cpx
```

You should see the TUI. If the problem row is empty, that is expected — create a `.cpp` file using `n` or fetch one from Competitive Companion.
