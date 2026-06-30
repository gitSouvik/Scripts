package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// ── Screen / dialog kind ──────────────────────────────────────────────────────

type screen int

const (
	screenMain       screen = iota
	screenDialog
	screenCustomTest
)

type dialogKind int

const (
	dialogNone dialogKind = iota
	dialogNew
	dialogNewSingle
	dialogNewRange
	dialogRrun
)

// ── Messages ──────────────────────────────────────────────────────────────────

type outLine struct {
	kind string // "log","ok","err","input","output","expected","debug","pass","fail","sep","diff-exp","diff-act"
	text string
}

type outLineMsg struct {
	line outLine
	ch   chan outLine
}

type opDoneMsg   struct{}
type fetchEvtMsg string
type refreshMsg  struct{}

// ── File entry ────────────────────────────────────────────────────────────────

func scanProblems(dir string) []string {
	entries, _ := os.ReadDir(dir)
	var probs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if strings.HasSuffix(n, ".cpp") {
			base := strings.TrimSuffix(n, ".cpp")
			if isProblemName(base) {
				probs = append(probs, base)
			}
		}
	}
	sort.Strings(probs)
	return probs
}

// isProblemName returns true for names like A, B, a, A1, B2, a1, A12, d1, etc.
// Pattern: one letter (upper or lower), followed by zero or more digits.
func isProblemName(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !((s[0] >= 'a' && s[0] <= 'z') || (s[0] >= 'A' && s[0] <= 'Z')) {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// ── Dialog state ──────────────────────────────────────────────────────────────

type dlgState struct {
	kind   dialogKind
	optSel int
	inputs []textinput.Model
	focus  int
}

func mkInput(placeholder string, limit int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = limit
	return ti
}

// ── Model ─────────────────────────────────────────────────────────────────────

type model struct {
	width  int
	height int
	cwd    string

	problems   []string
	problemSel int
	focusedRow int // 0: Problems, 1: Shortcuts
	logOffset  int // vertical scroll for logs

	outLines []outLine
	outCh    chan outLine
	running  bool

	screen      screen
	dlg         dlgState
	customInput textarea.Model // custom test case editor

	fetchRunning bool
	fetchMode    string
	fetchCh      chan fetchEvtMsg
}

func newModel() model {
	cwd, _ := os.Getwd()

	ta := textarea.New()
	ta.Placeholder = "Paste or type your custom test case here..."
	ta.ShowLineNumbers = true
	ta.CharLimit = 0 // unlimited

	// Style line numbers: faded gray with extra right-padding for clear visual separation
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		PaddingRight(2)
	activeDimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		PaddingRight(2)
	ta.FocusedStyle.LineNumber = dimStyle
	ta.FocusedStyle.CursorLineNumber = activeDimStyle
	ta.BlurredStyle.LineNumber = dimStyle
	ta.BlurredStyle.CursorLineNumber = dimStyle

	m := model{
		cwd:         cwd,
		fetchCh:     make(chan fetchEvtMsg, 64),
		customInput: ta,
	}
	m.problems = scanProblems(cwd)
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		launchFetch(m.fetchCh, m.cwd),
		listenFetch(m.fetchCh),
	)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case outLineMsg:
		m.outLines = append(m.outLines, msg.line)
		
		// Auto-scroll logic: if we're near the bottom, keep scrolling down
		availLines := m.height - 7
		if availLines < 5 {
			availLines = 5
		}
		if len(m.outLines) > availLines {
			m.logOffset = len(m.outLines) - availLines
		}
		
		return m, listenOutput(msg.ch)

	case opDoneMsg:
		m.running = false
		m.problems = scanProblems(m.cwd)

	case fetchEvtMsg:
		// Don't show log output, just silently update problem list
		m.problems = scanProblems(m.cwd)
		// Auto select the newly added problem by jumping to the end
		m.problemSel = len(m.problems) - 1
		return m, listenFetch(m.fetchCh)

	case refreshMsg:
		m.problems = scanProblems(m.cwd)

	case tea.MouseMsg:
		if m.screen == screenMain {
			if msg.Type == tea.MouseWheelUp {
				if m.logOffset > 0 {
					m.logOffset--
				}
			} else if msg.Type == tea.MouseWheelDown {
				m.logOffset++
			}
		}

	case tea.KeyMsg:
		if m.screen == screenDialog {
			return m.updateDlg(msg)
		}
		if m.screen == screenCustomTest {
			return m.updateCustomTest(msg)
		}
		return m.updateMain(msg)
	}
	return m, nil
}

// ── Main screen ───────────────────────────────────────────────────────────────

func (m model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.focusedRow > 0 {
			m.focusedRow--
		}

	case "down", "j":
		if m.focusedRow < 1 {
			m.focusedRow++
		}

	case "left", "h":
		if m.focusedRow == 0 && m.problemSel > 0 {
			m.problemSel--
		}

	case "right", "l":
		if m.focusedRow == 0 && m.problemSel < len(m.problems)-1 {
			m.problemSel++
		}

	case "r":
		if name := m.selSolution(); name != "" && !m.running {
			ch := make(chan outLine, 512)
			m.outCh = ch
			m.running = true
			m.outLines = nil
			m.logOffset = 0
			return m, tea.Batch(launchRun(name, ch), listenOutput(ch))
		}

	case "d":
		if m.selSolution() != "" && !m.running {
			inp := mkInput("blank=keyboard · 0=all · N=case N", 8)
			inp.Focus()
			m.screen = screenDialog
			m.dlg = dlgState{kind: dialogRrun, inputs: []textinput.Model{inp}}
			return m, textinput.Blink
		}

	case "g":
		m.outLines = append(m.outLines, outLine{"err", "[gen] not ready yet — write gen.cpp first"})
		m.logOffset = 0

	case "n":
		m.screen = screenDialog
		m.dlg = dlgState{kind: dialogNew, optSel: 0}

	case "e":
		if m.selSolution() != "" && !m.running {
			m.customInput.SetWidth(m.width - 4)
			m.customInput.SetHeight(m.height - 10)
			m.customInput.Focus()
			m.screen = screenCustomTest
			return m, textarea.Blink
		}

	case "c":
		m.outLines = nil
		m.logOffset = 0

	case "ctrl+r":
		m.problems = scanProblems(m.cwd)
	}
	return m, nil
}

// ── Custom Test screen ────────────────────────────────────────────────────────

func (m model) updateCustomTest(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.customInput.Blur()
		m.screen = screenMain
		return m, nil
	case "ctrl+r":
		inputText := m.customInput.Value()
		name := m.selSolution()
		m.customInput.Blur()
		m.screen = screenMain
		if name != "" && !m.running {
			ch := make(chan outLine, 512)
			m.outCh = ch
			m.running = true
			m.outLines = nil
			m.logOffset = 0
			return m, tea.Batch(launchCustomRun(name, inputText, ch), listenOutput(ch))
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.customInput, cmd = m.customInput.Update(msg)
		return m, cmd
	}
}

// ── Dialog screen ─────────────────────────────────────────────────────────────

func (m model) updateDlg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	d := m.dlg

	switch d.kind {

	case dialogNew:
		switch msg.String() {
		case "esc", "q":
			m.screen = screenMain
			m.dlg = dlgState{}
		case "up", "k":
			if d.optSel > 0 {
				d.optSel--
			}
		case "down", "j":
			if d.optSel < 1 {
				d.optSel++
			}
		case "1":
			d.optSel = 0
			return m.openSubDlg(0)
		case "m", "2":
			d.optSel = 1
			return m.openSubDlg(1)
		case "enter":
			return m.openSubDlg(d.optSel)
		}
		m.dlg = d

	case dialogNewSingle:
		switch msg.String() {
		case "esc":
			m.screen = screenDialog
			m.dlg = dlgState{kind: dialogNew, optSel: 0}
		case "tab", "down":
			d = fwdFocus(d)
		case "shift+tab", "up":
			d = bwdFocus(d)
		case "enter":
			if d.focus < len(d.inputs)-1 {
				d = fwdFocus(d)
			} else {
				name := strings.TrimSpace(d.inputs[0].Value())
				ext := strings.TrimSpace(d.inputs[1].Value())
				if name != "" && ext != "" {
					doCreate(m.cwd, []string{name}, ext)
					m.problems = scanProblems(m.cwd)
					m.outLines = append(m.outLines, outLine{kind: "ok", text: fmt.Sprintf("created  %s.%s", name, ext)})
				}
				m.screen = screenMain
				m.dlg = dlgState{}
				return m, nil
			}
		default:
			var cmd tea.Cmd
			d.inputs[d.focus], cmd = d.inputs[d.focus].Update(msg)
			m.dlg = d
			return m, cmd
		}
		m.dlg = d

	case dialogNewRange:
		switch msg.String() {
		case "esc":
			m.screen = screenDialog
			m.dlg = dlgState{kind: dialogNew, optSel: 1}
		case "tab", "down":
			d = fwdFocus(d)
		case "shift+tab", "up":
			d = bwdFocus(d)
		case "enter":
			if d.focus < len(d.inputs)-1 {
				d = fwdFocus(d)
			} else {
				endChar := strings.TrimSpace(d.inputs[0].Value())
				ext := strings.TrimSpace(d.inputs[1].Value())
				if len(endChar) == 1 && ext != "" {
					end := rune(endChar[0])
					if end >= 'a' && end <= 'z' {
						var names []string
						for c := rune('a'); c <= end; c++ {
							names = append(names, string(c))
						}
						doCreate(m.cwd, names, ext)
						m.problems = scanProblems(m.cwd)
						m.outLines = append(m.outLines, outLine{
							kind: "ok",
							text: fmt.Sprintf("created  a.%s  →  %s.%s  (%d files)", ext, string(end), ext, int(end-'a'+1)),
						})
					}
				}
				m.screen = screenMain
				m.dlg = dlgState{}
				return m, nil
			}
		default:
			var cmd tea.Cmd
			d.inputs[d.focus], cmd = d.inputs[d.focus].Update(msg)
			m.dlg = d
			return m, cmd
		}
		m.dlg = d

	case dialogRrun:
		switch msg.String() {
		case "esc":
			m.screen = screenMain
			m.dlg = dlgState{}
		case "enter":
			testNum := strings.TrimSpace(d.inputs[0].Value())
			name := m.selSolution()
			m.screen = screenMain
			m.dlg = dlgState{}
			if name != "" {
				ch := make(chan outLine, 512)
				m.outCh = ch
				m.running = true
				m.outLines = nil
				m.logOffset = 0
				return m, tea.Batch(launchRrun(name, testNum, ch), listenOutput(ch))
			}
		default:
			var cmd tea.Cmd
			d.inputs[0], cmd = d.inputs[0].Update(msg)
			m.dlg = d
			return m, cmd
		}
		m.dlg = d
	}

	return m, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (m model) selSolution() string {
	if m.problemSel < 0 || m.problemSel >= len(m.problems) {
		return ""
	}
	return m.problems[m.problemSel]
}

func (m model) openSubDlg(opt int) (model, tea.Cmd) {
	if opt == 0 {
		i0 := mkInput("E", 50)
		i0.Focus()
		i1 := mkInput("cpp", 20)
		m.dlg = dlgState{kind: dialogNewSingle, inputs: []textinput.Model{i0, i1}}
	} else {
		i0 := mkInput("e", 1)
		i0.CharLimit = 1
		i0.Focus()
		i1 := mkInput("cpp", 20)
		m.dlg = dlgState{kind: dialogNewRange, inputs: []textinput.Model{i0, i1}}
	}
	m.screen = screenDialog
	return m, textinput.Blink
}

func fwdFocus(d dlgState) dlgState {
	if d.focus < len(d.inputs)-1 {
		d.inputs[d.focus].Blur()
		d.focus++
		d.inputs[d.focus].Focus()
	}
	return d
}

func bwdFocus(d dlgState) dlgState {
	if d.focus > 0 {
		d.inputs[d.focus].Blur()
		d.focus--
		d.inputs[d.focus].Focus()
	}
	return d
}

func listenOutput(ch chan outLine) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return opDoneMsg{}
		}
		return outLineMsg{line: msg, ch: ch}
	}
}

func listenFetch(ch chan fetchEvtMsg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

// doCreate creates name.ext files in cwd
func doCreate(cwd string, names []string, ext string) {
	for _, name := range names {
		p := filepath.Join(cwd, name+"."+ext)
		f, err := os.Create(p)
		if err == nil {
			f.Close()
		}
	}
}
