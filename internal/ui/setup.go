package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type setupStep int

const (
	stepMode setupStep = iota
	stepDbUrl
	stepDbConnecting
	stepServerUrl
	stepServerConnecting
	stepAccount
	stepAccountCreating
	stepDone
)

// Setup is the first-run overlay shown when config.initialized is false. Its
// first (and currently only) stage collects the database URL.
type Setup struct {
	step            setupStep
	modeCursor      int
	modeChoices     []string
	urlInput        textinput.Model
	email           textinput.Model
	username        textinput.Model
	password        textinput.Model
	confirmPassword textinput.Model
	acctField       int
	err             string
}

func newSetupInput(prompt, placeholder string, masked bool) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.Placeholder = placeholder
	ti.SetWidth(50)
	if masked {
		ti.EchoMode = textinput.EchoPassword
	}
	return ti
}

func NewSetup() Setup {
	return Setup{
		step:            stepMode,
		modeChoices:     []string{"Standalone", "Client"},
		urlInput:        newSetupInput("DB URL> ", "postgres://user:pass@host:5432/db", false),
		email:           newSetupInput("Email> ", "you@example.com", false),
		username:        newSetupInput("Username> ", "username", false),
		password:        newSetupInput("Password> ", "", true),
		confirmPassword: newSetupInput("Confirm> ", "", true),
	}
}

func (s Setup) withUrlStep(step setupStep) Setup {
	s.step = step
	s.urlInput.SetValue("")
	if step == stepDbUrl {
		s.urlInput.Prompt = "DB URL> "
		s.urlInput.Placeholder = "postgres://user:pass@host:5432/db"
	} else if step == stepServerUrl {
		s.urlInput.Prompt = "Server URL> "
		s.urlInput.Placeholder = "https://example.com"
	}
	s.urlInput.Focus()
	return s
}

func (s Setup) focusAccountField() Setup {
	s.email.Blur()
	s.username.Blur()
	s.password.Blur()
	s.confirmPassword.Blur()
	switch s.acctField {
	case 0:
		s.email.Focus()
	case 1:
		s.username.Focus()
	case 2:
		s.password.Focus()
	case 3:
		s.confirmPassword.Focus()
	}
	return s
}

func (s Setup) Update(msg tea.Msg) (Setup, tea.Cmd) {
	if p, ok := msg.(tea.PasteMsg); ok {
		return s.paste(p)
	}
	kp, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return s, nil
	}

	switch s.step {
	case stepMode:
		choiceLen := len(s.modeChoices)
		switch kp.String() {
		case "up", "k":
			s.modeCursor = (s.modeCursor - 1 + choiceLen) % choiceLen
		case "down", "j":
			s.modeCursor = (s.modeCursor + 1) % choiceLen
		}
	case stepDbUrl, stepServerUrl:
		var cmd tea.Cmd
		s.urlInput, cmd = s.urlInput.Update(msg)
		return s, cmd
	case stepAccount:
		switch kp.String() {
		case "tab":
			s.acctField = (s.acctField + 1) % 4
			return s.focusAccountField(), nil
		case "shift+tab":
			s.acctField = (s.acctField - 1 + 4) % 4
			return s.focusAccountField(), nil
		}
		var cmd tea.Cmd
		switch s.acctField {
		case 0:
			s.email, cmd = s.email.Update(msg)
		case 1:
			s.username, cmd = s.username.Update(msg)
		case 2:
			s.password, cmd = s.password.Update(msg)
		case 3:
			s.confirmPassword, cmd = s.confirmPassword.Update(msg)
		}
		return s, cmd
	default:
		return s, nil
	}
	return s, nil
}

func (s Setup) View() string {
	var body string
	switch s.step {
	case stepMode:
		body = s.viewMode()
	case stepDbUrl, stepServerUrl:
		body = s.viewUrl()
	case stepDbConnecting, stepServerConnecting:
		body = s.viewConnecting()
	case stepAccount:
		body = s.viewAccount()
	case stepAccountCreating:
		body = s.viewCreatingAccount()
	default:
		return "This program is about to crash"
	}
	title := lipgloss.NewStyle().
		Bold(true).
		PaddingBottom(1).
		Render("First-run Setup")
	stack := lipgloss.JoinVertical(lipgloss.Left, title, body)
	if s.err != "" {
		errLine := lipgloss.NewStyle().
			Foreground(lipgloss.Red).
			PaddingTop(1).
			Render(s.err)
		stack = lipgloss.JoinVertical(lipgloss.Top, stack, errLine)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(stack)
}

func (s Setup) viewMode() string {
	rows := make([]string, len(s.modeChoices))
	for i, c := range s.modeChoices {
		style := lipgloss.NewStyle()
		if i == s.modeCursor {
			style = style.Bold(true).Underline(true)
		}
		rows[i] = style.Render(c)
	}
	hint := lipgloss.NewStyle().Faint(true).PaddingTop(1).
		Render("↑/↓ or j/k to select, Enter to confirm")
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().PaddingBottom(1).Render("Choose your mode"),
		lipgloss.JoinVertical(lipgloss.Top, rows...),
		hint,
	)
}

func (s Setup) viewUrl() string {
	var label string
	if s.step == stepDbUrl {
		label = "Enter the connection string for your database server"
	} else if s.step == stepServerUrl {
		label = "Enter the URL of the server to connect to"
	} else {
		label = "Fatal error this program is about to crash, after this step"
	}
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().PaddingBottom(1).Render(label),
		s.urlInput.View(),
	)
}

func (s Setup) viewConnecting() string {
	var msg string
	if s.step == stepDbConnecting {
		msg = "Connecting to database…"

	} else if s.step == stepServerConnecting {
		msg = "Connecting to server…"
	} else {
		msg = "Fatal error this program is going to crash after this step"
	}
	return lipgloss.NewStyle().Faint(true).Render(msg)
}

func (s Setup) viewCreatingAccount() string {
	return lipgloss.NewStyle().
		Faint(true).
		Render("Creating a new account...")
}

func (s Setup) viewAccount() string {
	fields := []struct {
		label string
		input textinput.Model
	}{
		{"Email", s.email},
		{"Username", s.username},
		{"Password", s.password},
		{"Confirm", s.confirmPassword},
	}
	rows := make([]string, len(fields))
	for i, f := range fields {
		lbl := lipgloss.NewStyle().Width(10)
		if i == s.acctField {
			lbl = lbl.Bold(true).Foreground(lipgloss.Yellow)
		}
		rows[i] = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lbl.Render(f.label),
			f.input.View(),
		)
	}
	hint := lipgloss.NewStyle().
		Faint(true).
		Render("Tab to move, Enter to advance / submit on the last field")
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().
			PaddingBottom(1).
			Render("Create your account"),
		lipgloss.JoinVertical(
			lipgloss.Top,
			rows...,
		),
		hint,
	)
}

func (s Setup) paste(p tea.PasteMsg) (Setup, tea.Cmd) {
	switch s.step {
	case stepDbUrl, stepServerUrl:
		var cmd tea.Cmd
		s.urlInput, cmd = s.urlInput.Update(p)
		return s, cmd
	case stepAccount:
		var cmd tea.Cmd
		switch s.acctField {
		case 0:
			s.email, cmd = s.email.Update(p)
		case 1:
			s.username, cmd = s.username.Update(p)
		case 2:
			s.password, cmd = s.password.Update(p)
		case 3:
			s.confirmPassword, cmd = s.confirmPassword.Update(p)
		}
		return s, cmd
	}
	return s, nil
}
