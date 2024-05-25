package component

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle                = lipgloss.NewStyle().Bold(true)
	DescriptionStyle          = lipgloss.NewStyle()
	DropdownItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedDropdownItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)
