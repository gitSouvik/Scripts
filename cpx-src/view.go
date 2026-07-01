package main

	import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// ── View ──────────────────────────────────────────────────────────────────────

func (m model) View() string {
	var b strings.Builder

	// Top Banner & Status
	timerStr := ""
	timerLen := 0
	if m.timerRunning {
		rem := time.Until(m.timerEnd)
		if rem < 0 {
			timerStr = "\033[5;31m TIME \033[0m" // flashing red
			timerLen = 6
		} else {
			timerText := fmt.Sprintf("%02d:%02d:%02d", int(rem.Hours()), int(rem.Minutes())%60, int(rem.Seconds())%60)
			timerLen = len(timerText)
			if rem < 5*time.Minute {
				timerStr = fmt.Sprintf("\033[1;31m%s\033[0m", timerText) // red
			} else {
				timerStr = fmt.Sprintf("\033[1;32m%s\033[0m", timerText) // green
			}
		}
	}
	
	titleLen := 25 // "cpx · code · test · debug"
	pad := (m.width - titleLen) / 2
	if pad < 0 { pad = 0 }
	
	banner := "\033[1;36mcpx\033[0m\033[2;37m · code · test · debug\033[0m"

	if timerStr != "" {
		leftPad := 2
		spacesAfterTimer := pad - leftPad - timerLen
		if spacesAfterTimer < 1 {
			spacesAfterTimer = 1
		}
		b.WriteString(fmt.Sprintf("\n%s%s%s%s\n\n", strings.Repeat(" ", leftPad), timerStr, strings.Repeat(" ", spacesAfterTimer), banner))
	} else {
		b.WriteString(fmt.Sprintf("\n%s%s\n\n", strings.Repeat(" ", pad), banner))
	}

	// ROW 1: Problems
	if m.focusedRow == 0 && m.screen == screenMain {
		b.WriteString("  \033[2;37mprob : \033[0m ")
	} else {
		b.WriteString("  \033[2;37mprob : \033[0m ")
	}

	if len(m.problems) == 0 {
		b.WriteString("\033[2;37m (No single-letter .cpp files found)\033[0m\n")
	} else {
		maxProbs := (m.width - 10) / 4
		if maxProbs < 1 {
			maxProbs = 1
		}

		startP := 0
		if m.problemSel >= maxProbs {
			startP = m.problemSel - maxProbs + 1
		}

		endP := startP + maxProbs
		if endP > len(m.problems) {
			endP = len(m.problems)
		}

		if startP > 0 {
			b.WriteString("\033[2;37m... \033[0m")
		}
		for i := startP; i < endP; i++ {
			p := m.problems[i]
			if i == m.problemSel {
				b.WriteString(fmt.Sprintf("\033[1;32m[%s]\033[0m ", p))
			} else {
				b.WriteString(fmt.Sprintf(" %s  ", p))
			}
		}
		if endP < len(m.problems) {
			b.WriteString("\033[2;37m...\033[0m")
		}
		b.WriteString("\n")
	}

	// ROW 2: Shortcuts
	if m.focusedRow == 1 && m.screen == screenMain {
		b.WriteString("  \033[2;37mcmds :\033[0m  ")
	} else {
		b.WriteString("  \033[2;37mcmds :\033[0m  ")
	}

	if m.running {
		b.WriteString("\033[2;37m[?] help   [r] run   [R] +run   [x] tests   [X] +tests   [m] more   [c] clear\033[0m\n")
	} else {
		b.WriteString("\033[1;32m[?]\033[0m help   \033[1;32m[r]\033[0m run   \033[1;32m[R]\033[0m +run   \033[1;32m[x]\033[0m tests   \033[1;32m[X]\033[0m +tests   \033[1;32m[m]\033[0m more   \033[1;32m[c]\033[0m clear\n")
	}
	lineWidth := m.width - 4
	if lineWidth < 0 {
		lineWidth = 0
	}
	b.WriteString(fmt.Sprintf("  \033[2;37m%s\033[0m\n", strings.Repeat("─", lineWidth)))
	
	// AI Tip
	if m.aiTip != "" && !m.hideFacts {
		words := strings.Fields(m.aiTip)
		maxW := m.width - 15
		if maxW < 10 {
			maxW = 10
		}
		var lines []string
		currentLine := ""
		for _, w := range words {
			if len(currentLine)+len(w)+1 > maxW {
				lines = append(lines, currentLine)
				currentLine = w
			} else {
				if currentLine == "" {
					currentLine = w
				} else {
					currentLine += " " + w
				}
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}
		
		for i, line := range lines {
			if i == 0 {
				b.WriteString(fmt.Sprintf("  \033[1;33mFacts!\033[0m \033[2;37m%s\033[0m\n", line))
			} else {
				b.WriteString(fmt.Sprintf("         \033[2;37m%s\033[0m\n", line))
			}
		}
	}

	// BODY: Help Screen
	if m.screen == screenHelp {
		b.WriteString("\n  \033[1;36m╭── CPX Documentation ───────────────────────────────────────────\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m  \033[2;37mPress any key to close this view\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m\n")
		
		// 1. Navigation & UI
		b.WriteString("  \033[1;36m│\033[0m  \033[1;37m✦ Navigation & UI\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32m←  /  →\033[0m     Move left/right between problems\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mh  /  l\033[0m     (Vim-style) Move left/right between problems\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32m↑  /  ↓\033[0m     Switch focus between problem row and commands row\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mscroll\033[0m      Scroll the output log up and down\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mctrl+h\033[0m      Toggle visibility of AI facts\n")
		b.WriteString("  \033[1;36m│\033[0m\n")
		
		// 2. Core Execution
		b.WriteString("  \033[1;36m│\033[0m  \033[1;37m✦ Execution Commands\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mr\033[0m           Compile and run against all standard `.in` test cases\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mR\033[0m           Interactive Run: Compiles, then awaits live custom input\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32md\033[0m           Debug Run: Select a specific test case number to run\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mc\033[0m           Clear the output console log\n")
		b.WriteString("  \033[1;36m│\033[0m\n")
		
		// 3. Testing & Generators
		b.WriteString("  \033[1;36m│\033[0m  \033[1;37m✦ Test Management\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mx\033[0m           Edit standard tests — view and modify `.in` files\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mX\033[0m           Add Custom Test — quick input editor for isolated testing\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;31mg\033[0m           Gen — Stress test generator (requires a `gen.cpp` file)\n")
		b.WriteString("  \033[1;36m│\033[0m\n")
		
		// 4. Tools & Features
		b.WriteString("  \033[1;36m│\033[0m  \033[1;37m✦ Extended Features\033[0m\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mt\033[0m           Timer — set a countdown timer for virtual contests\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32ms\033[0m           Snippets — open the code snippet injection menu\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mn\033[0m           New — dialog to generate a new problem file or range\n")
		b.WriteString("  \033[1;36m│\033[0m    \033[1;32mq\033[0m           Quit CPX environment\n")
		b.WriteString("  \033[1;36m╰────────────────────────────────────────────────────────────────\033[0m\n")
		return b.String()
	}

	// BODY: Test Case Editor
	if m.screen == screenTestEdit {
		fname := ""
		if len(m.testEditFiles) > 0 {
			fname = filepath.Base(m.testEditFiles[m.testEditIdx])
		}
		b.WriteString(fmt.Sprintf("\n  \033[1;36m─── Edit Tests\033[0m \033[2;37m— %s (%d/%d)\033[0m\n\n", fname, m.testEditIdx+1, len(m.testEditFiles)))
		b.WriteString(m.testEditInput.View())
		b.WriteString("\n\n  \033[2;37mctrl+s\033[0m: save   \033[2;37mtab\033[0m: next test   \033[2;37mshift+tab\033[0m: prev test   \033[2;37mesc\033[0m: done\n")
		return b.String()
	}

	// BODY: More Guide
	if m.screen == screenMore {
		b.WriteString("\n  \033[1;35m╭── Advanced Guide & Setup ──────────────────────────────────────\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m  \033[2;37mPress any key to close this view\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m\n")
		
		b.WriteString("  \033[1;35m│\033[0m  \033[1;37m✦ Competitive Companion Fetching\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m    CPX runs a background server listening on port \033[1;32m54321\033[0m.\n")
		b.WriteString("  \033[1;35m│\033[0m    Click the \033[1;32m+\033[0m icon in your browser extension to push problems.\n")
		b.WriteString("  \033[1;35m│\033[0m    Files are auto-generated: \033[2;37mA.cpp, A-1.in, A-1.out\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m\n")
		
		b.WriteString("  \033[1;35m│\033[0m  \033[1;37m✦ Debug Macros & Compilation\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m    Include \033[1;32mdebug.h\033[0m or \033[1;32mdebug++.h\033[0m in your files.\n")
		b.WriteString("  \033[1;35m│\033[0m    Use \033[1;36mdbg(var)\033[0m to print directly to the stderr console log.\n")
		b.WriteString("  \033[1;35m│\033[0m    Macros are seamlessly disabled when pushed to the judge.\n")
		b.WriteString("  \033[1;35m│\033[0m\n")
		
		b.WriteString("  \033[1;35m│\033[0m  \033[1;37m✦ Precompiled Headers (PCH)\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m    CPX takes advantage of \033[2;37mstdc++.h.gch\033[0m for ultra-fast compilation.\n")
		b.WriteString("  \033[1;35m│\033[0m    Compilation flags: \033[2;37m-O2 -std=c++20 -DWOOF_\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m\n")
		
		b.WriteString("  \033[1;35m│\033[0m  \033[1;37m✦ Custom Configuration\033[0m\n")
		b.WriteString("  \033[1;35m│\033[0m    \033[1;32mctrl+p\033[0m      Set your free Gemini API key for smart facts.\n")
		b.WriteString("  \033[1;35m╰────────────────────────────────────────────────────────────────\033[0m\n")
		return b.String()
	}

	// BODY: Snippets
	if m.screen == screenSnippets {
		b.WriteString("\n  \033[1;36m─── Snippets\033[0m \033[2;37m— ↑/↓ select • enter inject • n new • esc cancel\033[0m\n\n")
		for i, s := range snippets {
			if i == m.snippetSel {
				b.WriteString(fmt.Sprintf("  \033[1;36m❯\033[0m %-30s \033[2;37m%s\033[0m\n", s.name, strings.Split(s.code, "\n")[0]))
			} else {
				b.WriteString(fmt.Sprintf("    %-30s \033[2;37m%s\033[0m\n", s.name, strings.Split(s.code, "\n")[0]))
			}
		}
		return b.String()
	}

	// BODY: Snippet Creator
	if m.screen == screenSnippetCreate {
		b.WriteString("\n  \033[1;36m─── Create Custom Snippet\033[0m \033[2;37m— tab: switch • ctrl+s: save • esc: cancel\033[0m\n\n")
		nameFocused := "  "
		codeFocused := "  "
		if m.snipCreateFocused == 0 {
			nameFocused = "\033[1;36m❯\033[0m "
		} else {
			codeFocused = "\033[1;36m❯\033[0m "
		}
		b.WriteString(fmt.Sprintf("  %sName: %s\n\n", nameFocused, m.snipNameInput.View()))
		b.WriteString(fmt.Sprintf("  %sCode:\n", codeFocused))
		b.WriteString(m.snipCodeInput.View())
		b.WriteString("\n")
		return b.String()
	}

	// BODY: Custom Test Editor
	if m.screen == screenCustomTest {
		name := m.selSolution()
		b.WriteString(fmt.Sprintf("\n  \033[1;36m─── +tests\033[0m \033[2;37m(%s.cpp)\033[0m\n\n", name))
		b.WriteString(m.customInput.View())
		
		availLines := (m.height - 10) / 2
		if availLines < 5 { availLines = 5 }
		
		// Split pane logs
		fmt.Fprintf(&b, "\n  \033[2;37m%s\033[0m\n", strings.Repeat("─", m.width-4))
		b.WriteString(m.renderLogs(availLines))
		
		b.WriteString("\n  \033[2;37mctrl+r\033[0m: run   \033[2;37mctrl+s\033[0m: add as test   \033[2;37mctrl+a\033[0m: save & add another   \033[2;37mesc\033[0m: cancel\n")
		return b.String()
	}

	// BODY: Logs or Dialogs
	if m.screen == screenDialog {
		d := m.dlg
		b.WriteString("\n\033[1;36m  ─── \033[0m")
		switch d.kind {
		case dialogNew:
			b.WriteString("New File\033[1;36m ──────────────────────────────────\033[0m\n")
			opt1 := "  Single File (e.g. E.cpp)"
			opt2 := "  Range of Files (e.g. a.cpp to e.cpp)"
			if d.optSel == 0 {
				b.WriteString(fmt.Sprintf("\n  \033[1;36m❯\033[0m %s\n%s\n", opt1, opt2))
			} else {
				b.WriteString(fmt.Sprintf("\n%s\n  \033[1;36m❯\033[0m %s\n", opt1, opt2))
			}
			b.WriteString("\n  \033[2;37m↑/↓: select • enter: confirm • esc: cancel\033[0m")

		case dialogNewSingle:
			b.WriteString("New Single File\033[1;36m ───────────────────────────\033[0m\n\n")
			b.WriteString(fmt.Sprintf("  Name: %s\n", d.inputs[0].View()))
			b.WriteString(fmt.Sprintf("  Ext:  %s\n", d.inputs[1].View()))
			b.WriteString("\n  \033[2;37mtab: next • enter: create • esc: cancel\033[0m")

		case dialogNewRange:
			b.WriteString("New Range of Files\033[1;36m ────────────────────────\033[0m\n\n")
			b.WriteString(fmt.Sprintf("  Up to Char (e.g. 'e'): %s\n", d.inputs[0].View()))
			b.WriteString(fmt.Sprintf("  Ext:                   %s\n", d.inputs[1].View()))
			b.WriteString("\n  \033[2;37mtab: next • enter: create • esc: cancel\033[0m")

		case dialogRrun:
			b.WriteString("Debug Run\033[1;36m ─────────────────────────────────\033[0m\n\n")
			b.WriteString(fmt.Sprintf("  Test N: %s\n", d.inputs[0].View()))
			b.WriteString("\n  \033[2;37menter: run • esc: cancel\033[0m")
			
		case dialogTimer:
			b.WriteString("Contest Timer\033[1;36m ─────────────────────────────\033[0m\n\n")
			b.WriteString(fmt.Sprintf("  Duration: %s\n", d.inputs[0].View()))
			b.WriteString("\n  \033[2;37menter: start • esc: cancel\033[0m")
			
		case dialogApiKey:
			b.WriteString("Enter Gemini API Key\033[1;36m ───────────────────────\033[0m\n\n")
			b.WriteString(fmt.Sprintf("  Key: %s\n", m.apiKeyInput.View()))
			b.WriteString("\n  \033[2;37menter: save • esc: cancel\033[0m")
		}
		b.WriteString("\n")
	} else {
		availLines := m.height - 10
		if m.inlineInputActive {
			availLines -= (m.customInput.Height() + 2)
		}
		if availLines < 5 {
			availLines = 5
		}
		b.WriteString(m.renderLogs(availLines))
		if m.inlineInputActive {
			b.WriteString("\n  \033[1;36mType custom input:\033[0m \033[2;37m(enter or ctrl+r to run, esc to cancel)\033[0m\n")
			b.WriteString(m.customInput.View())
		}
	}

	return b.String()
}

func (m model) renderLogs(availLines int) string {
	var b strings.Builder
	
	start := m.logOffset
	if start < 0 {
		start = 0
	}
	if start > len(m.outLines)-availLines {
		start = len(m.outLines) - availLines
	}
	if start < 0 {
		start = 0
	}

	end := start + availLines
	if end > len(m.outLines) {
		end = len(m.outLines)
	}

	for _, l := range m.outLines[start:end] {
		prefix := ""
		color := "\033[0m"
		switch l.kind {
		case "log":
			color = "\033[2;37m"
		case "ok":
			color = "\033[1;32m"
			prefix = "✔ "
		case "err":
			color = "\033[1;31m"
			prefix = "✘ "
		case "input":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "output":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "output-err":
			color = "\033[1;31m"
			prefix = "\033[36m│\033[0m "
		case "diff-exp":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "diff-act":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "expected":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "debug":
			color = ""
			prefix = "\033[36m│\033[0m "
		case "pass":
			color = "\033[1;32m"
			prefix = "✔ "
		case "fail":
			color = "\033[1;31m"
			prefix = "✘ "
		case "sep":
			color = "\033[2;37m"
		}
		
		if l.kind == "diff-act" || l.kind == "diff-exp" {
			b.WriteString(fmt.Sprintf("  %s%s\033[0m\n", prefix, l.text))
		} else {
			b.WriteString(fmt.Sprintf("  %s%s%s\033[0m\n", color, prefix, l.text))
		}
	}
	return b.String()
}
