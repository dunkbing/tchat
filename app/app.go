package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dunkbing/tchat/channel"
	"github.com/dunkbing/tchat/chat"
)

const (
	BrowsingChannel  = "browsing-channel"
	CreatingChannel  = "creating-channel"
	FilteringChannel = "filtering-channel"
	Chatting         = "chatting"
)

type Model struct {
	channel *channel.Model
	chat    *chat.Model
	status  string
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	channelModel, channelCmd := m.channel.Update(msg)
	chatModel, chatCmd := m.chat.Update(msg)
	m.channel = channelModel.(*channel.Model)
	m.chat = chatModel.(*chat.Model)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.chat.Stop()
			return m, tea.Quit
		case tea.KeyCtrlA:
			m.channel.PreAddChannel()
			m.status = CreatingChannel
		case tea.KeyEsc:
			if m.status == Chatting {
				m.status = BrowsingChannel
				m.chat.Stop()
			}
			m.channel.PostAddChannel()
			m.chat.Reset()
			return m, nil
		case tea.KeyEnter:
			if m.status == CreatingChannel {
				m.status = BrowsingChannel
				m.channel.AddChannel()
				return m, channelCmd
			}
			if m.status == BrowsingChannel {
				m.status = Chatting
				m.chat.Reset()
				m.chat.SetChannel(m.channel.SelectedChannel())
				m.chat.LoadMessages()
				go m.chat.ReceiveMessages()
				return m, chatCmd
			}
			if m.status == Chatting {
				m.chat.SendMessage()
				return m, chatCmd
			}
		}

	case error:
		return m, nil
	}

	if m.status == Chatting {
		return m, tea.Batch(chatCmd)
	}
	return m, tea.Batch(channelCmd)
}

func (m *Model) View() string {
	if m.status == Chatting {
		return m.chat.View()
	}
	return m.channel.View()
}

func (m *Model) StartReceivingMsg() {
	m.chat.ReceiveMessages()
}

func New() tea.Model {
	channelModel := channel.New()

	chatModel := chat.New()
	model := &Model{
		channel: &channelModel,
		chat:    &chatModel,
		status:  BrowsingChannel,
	}
	channelModel.OnAddChannel = func(s string) {
		model.status = BrowsingChannel
	}
	channelModel.OnChooseChannel = func(s string) {
		if model.status == BrowsingChannel {
			model.status = Chatting
		}
	}
	//chatModel.ReceiveMessages()
	return model
}
