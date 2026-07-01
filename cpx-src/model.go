package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	screenHelp
	screenSnippets
	screenTestEdit
	screenSnippetCreate
	screenMore
)

type dialogKind int

const (
	dialogNone dialogKind = iota
	dialogNew
	dialogNewSingle
	dialogNewRange
	dialogRrun
	dialogTimer
	dialogApiKey
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
type timerTickMsg struct{}

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

	timerRunning bool
	timerEnd     time.Time
	timerDur     time.Duration

	snippetSel int

	testEditFiles []string
	testEditIdx   int
	testEditInput textarea.Model

	// AI Tip
	aiTip        string
	aiTipLoading bool
	ticks        int

	// New inputs for API key and custom snippet
	apiKeyInput       textinput.Model
	snipNameInput     textinput.Model
	snipCodeInput     textarea.Model
	snipCreateFocused int // 0: name, 1: code

	hideFacts bool
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

	loadApiKey(cwd)
	loadAllSnippets(cwd)

	// API key input
	apiInput := textinput.New()
	apiInput.Placeholder = "Paste your free Gemini API key here..."
	apiInput.Width = 50

	// Snippet name input
	snipName := textinput.New()
	snipName.Placeholder = "Snippet Name (e.g. Dijkstra)"
	snipName.Width = 30
	snipName.Focus()

	// Snippet code input
	snipCode := textarea.New()
	snipCode.Placeholder = "Paste snippet code here..."
	snipCode.ShowLineNumbers = true
	snipCode.CharLimit = 0
	snipCode.FocusedStyle.LineNumber = dimStyle
	snipCode.FocusedStyle.CursorLineNumber = activeDimStyle
	snipCode.BlurredStyle.LineNumber = dimStyle
	snipCode.BlurredStyle.CursorLineNumber = dimStyle

	initialTip := "Fetching tip..."
	initialTipLoading := true
	if currentApiKey == "" {
		initialTip = "for funny facts click ctrl + p to add free gemini key"
		initialTipLoading = false
	}

	m := model{
		cwd:               cwd,
		fetchCh:           make(chan fetchEvtMsg, 64),
		customInput:       ta,
		aiTip:             initialTip,
		aiTipLoading:      initialTipLoading,
		apiKeyInput:       apiInput,
		snipNameInput:     snipName,
		snipCodeInput:     snipCode,
		snipCreateFocused: 0,
	}
	m.problems = scanProblems(cwd)
	
	te := textarea.New()
	te.ShowLineNumbers = true
	te.CharLimit = 0
	te.FocusedStyle.LineNumber = dimStyle
	te.FocusedStyle.CursorLineNumber = activeDimStyle
	te.BlurredStyle.LineNumber = dimStyle
	te.BlurredStyle.CursorLineNumber = dimStyle
	m.testEditInput = te
	
	return m
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return timerTickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		launchFetch(m.fetchCh, m.cwd),
		listenFetch(m.fetchCh),
		tickCmd(),
	}
	if currentApiKey != "" {
		cmds = append(cmds, launchAITip())
	}
	return tea.Batch(cmds...)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.customInput.SetWidth(m.width - 4)
		m.customInput.SetHeight((m.height - 10) / 2)
		m.testEditInput.SetWidth(m.width - 4)
		m.testEditInput.SetHeight(m.height - 10)
		m.apiKeyInput.Width = m.width - 10
		m.snipNameInput.Width = m.width - 10
		m.snipCodeInput.SetWidth(m.width - 4)
		m.snipCodeInput.SetHeight(m.height - 12)

	case timerTickMsg:
		m.ticks++
		
		var cmds []tea.Cmd
		cmds = append(cmds, tickCmd()) // unconditional tick
		
		if m.ticks % 300 == 0 && !m.running {
			m.aiTipLoading = true
			cmds = append(cmds, launchAITip())
		}
		
		return m, tea.Batch(cmds...)

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
		if len(m.problems) > 0 {
			m.problemSel = len(m.problems) - 1
		}
		return m, listenFetch(m.fetchCh)
		
	case aiTipMsg:
		m.aiTip = msg.tip
		m.aiTipLoading = false

	case aiAckMsg:
		m.outLines = append(m.outLines, outLine{kind: "ok", text: msg.text})
		// Auto-scroll logs
		availLines := m.height - 7
		if availLines < 5 {
			availLines = 5
		}
		if len(m.outLines) > availLines {
			m.logOffset = len(m.outLines) - availLines
		}

	case refreshMsg:
		m.problems = scanProblems(m.cwd)

	case tea.MouseMsg:
		switch m.screen {
		case screenMain:
			switch msg.Type {
			case tea.MouseWheelUp:
				if m.logOffset > 0 {
					m.logOffset--
				}
			case tea.MouseWheelDown:
				m.logOffset++
			case tea.MouseLeft:
				if msg.Y == 4 {
					if msg.X >= 10 {
						maxProbs := (m.width - 10) / 4
						if maxProbs < 1 { maxProbs = 1 }
						clickedIdx := (msg.X - 10) / 4
						if clickedIdx >= 0 && clickedIdx < maxProbs {
							startP := 0
							if m.problemSel >= maxProbs {
								startP = m.problemSel - maxProbs + 1
							}
							targetP := startP + clickedIdx
							if targetP < len(m.problems) {
								m.problemSel = targetP
							}
						}
					}
				}
			}
		case screenCustomTest:
			return m.updateCustomTest(msg)
		case screenTestEdit:
			return m.updateTestEdit(msg)
		case screenSnippetCreate:
			return m.updateSnippetCreate(msg)
		}

	case tea.KeyMsg:
		if m.screen == screenHelp || m.screen == screenMore {
			m.screen = screenMain
			return m, nil
		}
		if m.screen == screenDialog {
			return m.updateDlg(msg)
		}
		if m.screen == screenCustomTest {
			return m.updateCustomTest(msg)
		}
		if m.screen == screenSnippets {
			return m.updateSnippets(msg)
		}
		if m.screen == screenTestEdit {
			return m.updateTestEdit(msg)
		}
		if m.screen == screenSnippetCreate {
			return m.updateSnippetCreate(msg)
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
		// Locked to problems row (0) for now

	case "down", "j":
		// Locked to problems row (0) for now

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
			m.customInput.SetHeight((m.height - 10) / 2)
			m.customInput.Focus()
			m.screen = screenCustomTest
			return m, textarea.Blink
		}

	case "c":
		m.outLines = nil
		m.logOffset = 0

	case "t":
		if m.timerRunning {
			m.timerRunning = false
		} else {
			m.screen = screenDialog
			inp := mkInput("Timer (e.g. 120m, 2h30m) [Default: 2h]", 20)
			inp.Focus()
			m.dlg = dlgState{kind: dialogTimer, inputs: []textinput.Model{inp}}
			return m, textinput.Blink
		}

	case "s":
		if m.selSolution() != "" && !m.running {
			m.screen = screenSnippets
			m.snippetSel = 0
		}

	case "m":
		if !m.running {
			m.screen = screenMore
		}

	case "ctrl+p":
		m.screen = screenDialog
		m.apiKeyInput.SetValue("")
		m.apiKeyInput.Focus()
		m.dlg = dlgState{kind: dialogApiKey}
		return m, textinput.Blink

	case "ctrl+h":
		m.hideFacts = !m.hideFacts
		return m, nil

	case "x":
		if name := m.selSolution(); name != "" && !m.running {
			inFiles, _ := filepath.Glob(filepath.Join(m.cwd, name+"-*.in"))
			if len(inFiles) > 0 {
				sort.Strings(inFiles)
				m.testEditFiles = inFiles
				m.testEditIdx = 0
				content, _ := os.ReadFile(inFiles[0])
				m.testEditInput.SetValue(string(content))
				m.testEditInput.SetWidth(m.width - 4)
				m.testEditInput.SetHeight(m.height - 10)
				m.testEditInput.Focus()
				m.screen = screenTestEdit
				return m, textarea.Blink
			} else {
				m.outLines = append(m.outLines, outLine{"err", "No .in files to edit"})
			}
		}

	case "?":
		m.screen = screenHelp

	case "ctrl+r":
		m.problems = scanProblems(m.cwd)
	}
	return m, nil
}

// ── Custom Test screen ────────────────────────────────────────────────────────

func (m model) updateCustomTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.customInput.Blur()
			m.screen = screenMain
			return m, nil
		case "ctrl+r":
			inputText := m.customInput.Value()
			name := m.selSolution()
			if name != "" && !m.running {
				ch := make(chan outLine, 512)
				m.outCh = ch
				m.running = true
				m.outLines = nil
				m.logOffset = 0
				return m, tea.Batch(launchCustomRun(name, inputText, ch), listenOutput(ch))
			}
			return m, nil
		case "ctrl+s":
			return m.saveCustomTest(false), nil
		case "ctrl+a":
			return m.saveCustomTest(true), nil
		}
	case tea.MouseMsg:
		msg.Y -= 3
		var cmd tea.Cmd
		m.customInput, cmd = m.customInput.Update(msg)
		return m, cmd
	}
	
	// Default handler for textinput updates
	var cmd tea.Cmd
	m.customInput, cmd = m.customInput.Update(msg)
	return m, cmd
}

func (m model) saveCustomTest(addAnother bool) model {
	name := m.selSolution()
	if name != "" {
		inFiles, _ := filepath.Glob(filepath.Join(m.cwd, name+"-*.in"))
		maxN := 0
		for _, f := range inFiles {
			base := filepath.Base(f)
			base = strings.TrimPrefix(base, name+"-")
			base = strings.TrimSuffix(base, ".in")
			var n int
			fmt.Sscanf(base, "%d", &n)
			if n > maxN {
				maxN = n
			}
		}
		newNum := maxN + 1
		inFile := filepath.Join(m.cwd, fmt.Sprintf("%s-%d.in", name, newNum))
		outFile := filepath.Join(m.cwd, fmt.Sprintf("%s-%d.out", name, newNum))
		os.WriteFile(inFile, []byte(m.customInput.Value()), 0644)
		os.WriteFile(outFile, []byte(""), 0644)
		m.outLines = append(m.outLines, outLine{"ok", fmt.Sprintf("Saved test case %s-%d.in", name, newNum)})
		
		if addAnother {
			m.customInput.SetValue("")
		} else {
			m.customInput.Blur()
			m.screen = screenMain
		}
	}
	return m
}

// ── Snippets screen ───────────────────────────────────────────────────────────

func (m model) updateSnippets(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.screen = screenMain
	case "n":
		m.screen = screenSnippetCreate
		m.snipNameInput.SetValue("")
		m.snipCodeInput.SetValue("")
		m.snipNameInput.Focus()
		m.snipCodeInput.Blur()
		m.snipCreateFocused = 0
		return m, nil
	case "up", "k":
		if m.snippetSel > 0 {
			m.snippetSel--
		}
	case "down", "j":
		if m.snippetSel < len(snippets)-1 {
			m.snippetSel++
		}
	case "enter":
		name := m.selSolution()
		if name != "" {
			cppFile := filepath.Join(m.cwd, name+".cpp")
			content, err := os.ReadFile(cppFile)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				insertIdx := 0
				for i := len(lines) - 1; i >= 0; i-- {
					if strings.HasPrefix(strings.TrimSpace(lines[i]), "#include") || strings.HasPrefix(strings.TrimSpace(lines[i]), "using namespace") {
						insertIdx = i + 1
						break
					}
				}
				snippetCode := "\n// --- " + snippets[m.snippetSel].name + " ---\n" + snippets[m.snippetSel].code + "\n"
				newLines := append(lines[:insertIdx], append(strings.Split(snippetCode, "\n"), lines[insertIdx:]...)...)
				os.WriteFile(cppFile, []byte(strings.Join(newLines, "\n")), 0644)
				m.outLines = append(m.outLines, outLine{"ok", fmt.Sprintf("Injected snippet: %s", snippets[m.snippetSel].name)})
			}
		}
		m.screen = screenMain
	}
	return m, nil
}

func (m model) updateSnippetCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.screen = screenSnippets
			return m, nil
		case "tab", "down":
			m.snipCreateFocused = (m.snipCreateFocused + 1) % 2
			if m.snipCreateFocused == 0 {
				m.snipNameInput.Focus()
				m.snipCodeInput.Blur()
			} else {
				m.snipNameInput.Blur()
				m.snipCodeInput.Focus()
			}
		case "shift+tab", "up":
			m.snipCreateFocused = (m.snipCreateFocused - 1 + 2) % 2
			if m.snipCreateFocused == 0 {
				m.snipNameInput.Focus()
				m.snipCodeInput.Blur()
			} else {
				m.snipNameInput.Blur()
				m.snipCodeInput.Focus()
			}
		case "ctrl+s":
			name := strings.TrimSpace(m.snipNameInput.Value())
			code := m.snipCodeInput.Value()
			if name != "" && code != "" {
				saveCustomSnippet(m.cwd, name, code)
				m.screen = screenMain
				m.outLines = append(m.outLines, outLine{"ok", fmt.Sprintf("Saved snippet: %s", name)})
				return m, nil
			}
		default:
			var cmd tea.Cmd
			if m.snipCreateFocused == 0 {
				m.snipNameInput, cmd = m.snipNameInput.Update(msg)
			} else {
				m.snipCodeInput, cmd = m.snipCodeInput.Update(msg)
			}
			return m, cmd
		}
	case tea.MouseMsg:
		if msg.Y == 3 || msg.Y == 4 {
			m.snipCreateFocused = 0
			m.snipNameInput.Focus()
			m.snipCodeInput.Blur()
			var cmd tea.Cmd
			m.snipNameInput, cmd = m.snipNameInput.Update(msg)
			return m, cmd
		} else if msg.Y >= 6 {
			m.snipCreateFocused = 1
			m.snipNameInput.Blur()
			m.snipCodeInput.Focus()
			msg.Y -= 6
			var cmd tea.Cmd
			m.snipCodeInput, cmd = m.snipCodeInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// ── Test Edit screen ──────────────────────────────────────────────────────────

func (m model) updateTestEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			os.WriteFile(m.testEditFiles[m.testEditIdx], []byte(m.testEditInput.Value()), 0644)
			m.testEditInput.Blur()
			m.screen = screenMain
			return m, nil
		case "ctrl+s", "ctrl+w":
			os.WriteFile(m.testEditFiles[m.testEditIdx], []byte(m.testEditInput.Value()), 0644)
			m.outLines = append(m.outLines, outLine{"log", fmt.Sprintf("Saved %s", filepath.Base(m.testEditFiles[m.testEditIdx]))})
		case "tab":
			os.WriteFile(m.testEditFiles[m.testEditIdx], []byte(m.testEditInput.Value()), 0644)
			m.testEditIdx = (m.testEditIdx + 1) % len(m.testEditFiles)
			content, _ := os.ReadFile(m.testEditFiles[m.testEditIdx])
			m.testEditInput.SetValue(string(content))
		case "shift+tab":
			os.WriteFile(m.testEditFiles[m.testEditIdx], []byte(m.testEditInput.Value()), 0644)
			m.testEditIdx = (m.testEditIdx - 1 + len(m.testEditFiles)) % len(m.testEditFiles)
			content, _ := os.ReadFile(m.testEditFiles[m.testEditIdx])
			m.testEditInput.SetValue(string(content))
		}
	case tea.MouseMsg:
		msg.Y -= 4
		var cmd tea.Cmd
		m.testEditInput, cmd = m.testEditInput.Update(msg)
		return m, cmd
	}
	
	// Default handler for textinput updates
	var cmd tea.Cmd
	m.testEditInput, cmd = m.testEditInput.Update(msg)
	return m, cmd
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
		case dialogTimer:
		switch msg.String() {
		case "esc":
			m.screen = screenMain
			m.dlg = dlgState{}
		case "enter":
			val := strings.TrimSpace(d.inputs[0].Value())
			if val == "" {
				val = "2h"
			}
			dur, err := time.ParseDuration(val)
			if err == nil {
				m.timerDur = dur
				m.timerEnd = time.Now().Add(dur)
				m.timerRunning = true
			}
			m.screen = screenMain
			m.dlg = dlgState{}
			return m, tickCmd()
		default:
			var cmd tea.Cmd
			d.inputs[0], cmd = d.inputs[0].Update(msg)
			m.dlg = d
			return m, cmd
		}

	case dialogApiKey:
		switch msg.String() {
		case "esc":
			m.screen = screenMain
			m.dlg = dlgState{}
		case "enter":
			key := strings.TrimSpace(m.apiKeyInput.Value())
			saveApiKey(m.cwd, key)
			m.screen = screenMain
			m.dlg = dlgState{}
			m.aiTip = "Fetching tip..."
			m.aiTipLoading = true
			return m, launchAITip()
		default:
			var cmd tea.Cmd
			m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
			return m, cmd
		}
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
