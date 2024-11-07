package table

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	c "github.com/nabinthapaa/vantage-cli/constants"
	"github.com/nabinthapaa/vantage-cli/modules/conservation_mode"
	"github.com/nabinthapaa/vantage-cli/modules/fn_lock"
	"github.com/nabinthapaa/vantage-cli/modules/usb_charging"
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

func (m model) updateTableValues() {
	var err error
	m.c_mode.value, err = conservation_mode.GetCurrentValue()
	m.fn_lock.value, err = fn_lock.GetCurrentValue()
	m.usb_charging.value, err = usb_charging.GetCurrentValue()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func initializeModel() *model {
	var err error
	c_mode, err := conservation_mode.GetCurrentValue()
	fn_lock, err := fn_lock.GetCurrentValue()
	usb_charging, err := usb_charging.GetCurrentValue()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	c_mode_struct := &File{
		c.CONSERVATION_MODE, c_mode,
	}
	fn_struct := &File{
		c.FN_LOCK, fn_lock,
	}
	usb_charging_struct := &File{
		c.USB_CHARGING, usb_charging,
	}

	columns := []table.Column{
		{Title: "Modes", Width: 30},
		{Title: "Status", Width: 15},
		{Title: "", Width: 0},
	}
	rows := []table.Row{
		{"Conservation Mode", c_mode_struct.value, c_mode_struct.filename},
		{"Fn Lock", fn_struct.value, fn_struct.filename},
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
		{"Conservation Mode", m.c_mode.value, m.c_mode.filename},
		{"Fn Lock", m.fn_lock.value, m.fn_lock.filename},
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
			file, value := m.t.SelectedRow()[2], m.t.SelectedRow()[1]
			switch file {
			case c.CONSERVATION_MODE:
				if err := conservation_mode.UpdateCurrentValue(value); err != nil {
					log.Fatalf("Error: ", err)
				}
			case c.FN_LOCK:
				if err := fn_lock.UpdateCurrentValue(value); err != nil {
					log.Fatalf("Error: ", err)
				}
			case c.USB_CHARGING:
				if err := usb_charging.UpdateCurrentValue(value); err != nil {
					log.Fatalf("Error: ", err)
				}
			}
			m.updateTableValues()
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
