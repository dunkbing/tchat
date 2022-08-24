package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dunkbing/tchat/models/channel"
	"github.com/dunkbing/tchat/models/chat"
)

type Model struct {
	channel  channel.Model
	chat     chat.Model
	chatting bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	channelModel, chCmd := m.channel.Update(msg)
	chatModel, cCmd := m.chat.Update(msg)
	m.channel = channelModel.(channel.Model)
	m.chat = chatModel.(chat.Model)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.chatting {
				m.chatting = false
			}
		case tea.KeyEnter:
			if !m.chatting {
				m.chatting = true
			}
		}

	case error:
		return m, nil
	}

	if m.chatting {
		return m, tea.Batch(chCmd)
	}
	return m, tea.Batch(cCmd)
}

func (m Model) View() string {
	if m.chatting {
		return m.chat.View()
	}
	return m.channel.View()
}

func New() tea.Model {
	channelModel := channel.New()
	chatModel := chat.New()
	model := Model{
		channel:  channelModel,
		chat:     chatModel,
		chatting: false,
	}
	return model
}
