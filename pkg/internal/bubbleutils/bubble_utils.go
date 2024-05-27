package bubbleutils

import tea "github.com/charmbracelet/bubbletea"

func RenderAll[T tea.Model, S ~[]T](mods S) []string {
	out := make([]string, 0, len(mods))
	for _, v := range mods {
		out = append(out, v.View())
	}
	return out
}

func UpdateAll(msg tea.Msg, mods []tea.Model) tea.Cmd {
	cmds := make([]tea.Cmd, len(mods))
	for i, v := range mods {
		mods[i], cmds[i] = v.Update(msg)
	}
	return tea.Batch(cmds...)
}

func InitAll[T tea.Model, S ~[]T](mods S) tea.Cmd {
	cmds := make([]tea.Cmd, len(mods))
	for i, v := range mods {
		cmds[i] = v.Init()
	}
	return tea.Batch(cmds...)
}
