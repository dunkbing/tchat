package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dunkbing/tchat/redis"
	"github.com/dunkbing/tchat/utils"
	"os"
	"strings"
	"time"
)

const messageExpire = 5
const messagePrefix = "message:"

type Message struct {
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
	CreatedAt time.Time `json:"createdAt"`
}

type Model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

	channel string
	sender  string

	stopped bool
}

func formatMessage(msg Message, style lipgloss.Style) string {
	return fmt.Sprintf("%s%s\n%s", style.Render(msg.Sender+": "), msg.Message, msg.CreatedAt.Format(time.Kitchen))
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Stop() {
	m.stopped = true
}

func (m *Model) SetChannel(c string) {
	m.channel = c
	m.viewport.SetContent(fmt.Sprintf("Welcome to %s", m.channel))
}

func (m *Model) LoadMessages() {
	ctx := context.Background()
	iter := redis.Client.Scan(ctx, 0, fmt.Sprintf("%s%s*", messagePrefix, m.channel), 0).Iterator()
	var messages []string
	for iter.Next(ctx) {
		rawMsg := iter.Val()
		rawMsg = strings.Replace(rawMsg, fmt.Sprintf("%s%s:", messagePrefix, m.channel), "", 1)
		message := Message{}
		_ = json.Unmarshal([]byte(rawMsg), &message)

		if message.Sender == m.sender {
			message.Sender = "You"
		}

		messages = append(messages, formatMessage(message, m.senderStyle))
	}
	if err := iter.Err(); err != nil {
		fmt.Println("err: ", err.Error())
	}
	m.messages = messages
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
}

func (m *Model) SendMessage() {
	message := m.textarea.Value()
	createdAt := time.Now()
	m.messages = append(
		m.messages,
		formatMessage(Message{
			Message:   message,
			Sender:    "You",
			CreatedAt: createdAt,
		}, m.senderStyle),
	)
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
	ctx := context.Background()
	data := Message{
		Message:   message,
		Sender:    m.sender,
		CreatedAt: createdAt,
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err set key: ", err)
		return
	}
	redis.Client.Set(
		ctx,
		fmt.Sprintf("%s%s:%s", messagePrefix, m.channel, jsonStr),
		"",
		time.Minute*messageExpire,
	)
	redis.Client.Publish(ctx, m.channel, jsonStr)
}

func (m *Model) ReceiveMessages() {
	m.stopped = false
	ctx := context.Background()
	pubSub := redis.Client.Subscribe(ctx, m.channel)
	defer pubSub.Close()

	ch := pubSub.Channel()
	for msg := range ch {
		if m.stopped {
			break
		}
		message := Message{}
		_ = json.Unmarshal([]byte(msg.Payload), &message)
		if message.Sender != m.sender {
			m.messages = append(m.messages, formatMessage(message, m.senderStyle))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
		}
	}
}

func (m *Model) Reset() {
	m.textarea.Reset()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m *Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func New() Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 10)
	vp.SetContent(`Welcome`)

	sender, err := os.Hostname()
	if err != nil {
		sender = utils.RandSeq(8)
	}

	return Model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		sender:      sender,
	}
}
