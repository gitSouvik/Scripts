package main

import (
	"fmt"
	"strings"
)

// ── View ──────────────────────────────────────────────────────────────────────

func (m model) View() string {
	var b strings.Builder

	// Top Banner & Status
	b.WriteString("\n  \033[1;36mcpx\033[0m — Competitive Programming Tools\n")
	b.WriteString(fmt.Sprintf("  \033[1;34m%s\033[0m\n\n", m.cwd))

	// ROW 1: Problems
	if m.focusedRow == 0 && m.screen == screenMain {
		b.WriteString("\033[1;36m│\033[0m ")
	} else {
		b.WriteString("  ")
	}

	if len(m.problems) == 0 {
		b.WriteString("\033[2;37m(No single-letter .cpp files found)\033[0m\n")
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
		b.WriteString("\033[1;36m│\033[0m ")
	} else {
		b.WriteString("  ")
	}
	
	if m.running {
		b.WriteString("\033[2;37mType :  [r] run   [e] custom   [d] debug   [c] clear   [n] new   [g] gen   [q] quit\033[0m\n")
	} else {
		b.WriteString("\033[2;37mType :\033[0m  \033[1;32m[r]\033[0m run   \033[1;32m[e]\033[0m custom   \033[1;32m[d]\033[0m debug   \033[1;32m[c]\033[0m clear   \033[1;32m[n]\033[0m new   \033[1;31m[g]\033[0m gen   \033[1;32m[q]\033[0m quit\n")
	}
	lineWidth := m.width - 4
	if lineWidth < 0 {
		lineWidth = 0
	}
	b.WriteString(fmt.Sprintf("  \033[2;37m%s\033[0m\n", strings.Repeat("─", lineWidth)))

	// BODY: Custom Test Editor
	if m.screen == screenCustomTest {
		name := m.selSolution()
		b.WriteString(fmt.Sprintf("\n  \033[1;36m─── Custom Test\033[0m \033[2;37m(%s.cpp)\033[0m\n\n", name))
		b.WriteString(m.customInput.View())
		b.WriteString("\n\n  \033[2;37mctrl+r\033[0m: run with this input   \033[2;37mesc\033[0m: cancel\n")
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
		}
		b.WriteString("\n")
	} else {
		// Output lines with scroll
		availLines := m.height - 10
		if availLines < 5 {
			availLines = 5
		}
		
		start := m.logOffset
		if start < 0 {
			start = 0
		}
		if start > len(m.outLines) - availLines {
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
			b.WriteString(fmt.Sprintf("  %s%s%s\033[0m\n", color, prefix, l.text))
		}
	}

	return b.String()
}
