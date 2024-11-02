package table

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nabinthapaa/vantage-cli/constants"
)

type File struct {
	filename string
	value    string
}

type model struct {
	t            table.Model
	c_mode       *File
	fn_lock      *File
	usb_charging *File
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func GetValues(s string, on int, off int) string {
	file, err := os.OpenFile(s, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open conservation_mode file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		log.Fatal("Failed to read from conservation_mode file")
	}
	fileContent := scanner.Text()

	mode, err := strconv.Atoi(fileContent)
	if err != nil {
		log.Fatalf("Error parsing mode: %v", err)
	}

	if mode == on {
		return "On"
	} else if mode == off {
		return "Off"
	} else {
		return "Invalid value"
	}
}

func handleOnOff(s string, val string) error {
	var mode string
	file, err := os.OpenFile(s, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open conservation_mode file: %v", err)
	}
	defer file.Close()

	if val == "On" {
		mode = "0"
	} else if val == "Off" {
		mode = "1"
	}

	if _, err := file.WriteString(mode); err != nil {
		return fmt.Errorf("Failed to write to conservation_mode file: %v", err)
	}

	return nil
}

func (m model) updateValue() {
	c_mode := GetValues(constants.CONSERVATION_MODE_FILE, 1, 0)
	fn_lock := GetValues(constants.FN_LOCK, 1, 0)
	usb_charging := GetValues(constants.USB_CHARGING, 1, 0)

	m.c_mode.value = c_mode
	m.fn_lock.value = fn_lock
	m.usb_charging.value = usb_charging
}

func initializeModel() *model {
	c_mode := GetValues(constants.CONSERVATION_MODE_FILE, 1, 0)
	fn_lock := GetValues(constants.FN_LOCK, 1, 0)
	usb_charging := GetValues(constants.USB_CHARGING, 1, 0)

	c_mode_struct := &File{
		constants.CONSERVATION_MODE_FILE, c_mode,
	}
	fn_struct := &File{
		constants.FN_LOCK, fn_lock,
	}
	usb_charging_struct := &File{
		constants.CONSERVATION_MODE_FILE, usb_charging,
	}

	columns := []table.Column{
		{Title: "Modes", Width: 30},
		{Title: "Status", Width: 15},
		{Title: "", Width: 0},
	}
	rows := []table.Row{
		{"Fn Lock", fn_struct.value, fn_struct.filename},
		{"Conservation Mode", c_mode_struct.value, c_mode_struct.filename},
		{"Usb Charging", usb_charging_struct.value, usb_charging_struct.filename},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(14),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &model{
		t,
		c_mode_struct,
		fn_struct,
		usb_charging_struct,
	}
}

func (m model) updateTable() []table.Row {
	rows := []table.Row{
		{"Fn Lock", m.fn_lock.value, m.fn_lock.filename},
		{"Conservation Mode", m.c_mode.value, m.c_mode.filename},
		{"Usb Charging", m.usb_charging.value, m.usb_charging.filename},
	}

	return rows
}

func Run() error {
	model := initializeModel()
	_, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			handleOnOff(m.t.SelectedRow()[2], m.t.SelectedRow()[1])
			m.updateValue()
			m.t.SetRows(m.updateTable())
			return m, nil
		}
	}
	m.t, cmd = m.t.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return baseStyle.Render(m.t.View()) + "\n  " + m.t.HelpView() + "\n"
}
