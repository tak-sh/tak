package keyregistry

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Submit key.Binding
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

var DefaultKeys = &KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
}

var DebugKeys = &DebugKeyMap{
	KeyMap: DefaultKeys,
	Reload: key.NewBinding(
		key.WithKeys(tea.KeyF8.String()),
		key.WithHelp("F8", "rerun step"),
	),
	Previous: key.NewBinding(
		key.WithKeys(tea.KeyF6.String()),
		key.WithHelp("F6", "previous step"),
	),
	Next: key.NewBinding(
		key.WithKeys(tea.KeyF7.String()),
		key.WithHelp("F7", "next step"),
	),
}

type DebugKeyMap struct {
	*KeyMap
	Reload   key.Binding
	Previous key.Binding
	Next     key.Binding
}

func (d *DebugKeyMap) ShortHelp() []key.Binding {
	return append([]key.Binding{d.Next, d.Previous, d.Reload}, d.KeyMap.ShortHelp()...)
}

func (d *DebugKeyMap) FullHelp() [][]key.Binding {
	return append(d.KeyMap.FullHelp(),
		[]key.Binding{
			d.Next, d.Previous, d.Reload,
		},
	)
}
