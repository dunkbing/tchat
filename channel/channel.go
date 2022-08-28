package channel

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dunkbing/tchat/redis"
	"io"
	"strings"
	"time"
)

const listHeight = 14
const channelPrefix = "channel:"
const defaultChannel = "Default"
const channelExpire = 60

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type Item string

func (i Item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	_, _ = fmt.Fprintf(w, fn(str))
}

type Model struct {
	list            list.Model
	textarea        textarea.Model
	adding          bool
	OnAddChannel    func(string)
	OnChooseChannel func(string)
}

func (m *Model) PreAddChannel() {
	m.adding = true
	m.textarea.Focus()
}

func (m *Model) AddChannel() {
	newChan := m.textarea.Value()
	if len(newChan) > 0 {
		m.list.InsertItem(0, Item(newChan))
		ctx := context.Background()
		redis.Client.Set(ctx, fmt.Sprintf("%s%s", channelPrefix, newChan), newChan, time.Minute*channelExpire)
	}
	m.textarea.Reset()
	if m.OnAddChannel != nil {
		m.OnAddChannel(newChan)
	}
	m.PostAddChannel()
}

func (m *Model) PostAddChannel() {
	m.textarea.Reset()
	m.adding = false
}

func (m *Model) SelectedChannel() string {
	return m.list.SelectedItem().FilterValue()
}

func (m *Model) Init() tea.Cmd {
	//return textarea.Blink
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmd tea.Cmd
	var taCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	m.textarea, taCmd = m.textarea.Update(msg)

	return m, tea.Batch(listCmd, taCmd)
}

func (m *Model) View() string {
	if !m.adding {
		return m.list.View()
	}
	return fmt.Sprintf("%s\n\n%s", m.list.View(), m.textarea.View())
}

func New() Model {
	items := []list.Item{
		Item(defaultChannel),
	}
	ctx := context.Background()
	iter := redis.Client.Scan(ctx, 0, fmt.Sprintf("%s*", channelPrefix), 0).Iterator()
	for iter.Next(ctx) {
		channel := iter.Val()
		channel = strings.Replace(channel, channelPrefix, "", 1)
		items = append(items, Item(channel))
	}
	if err := iter.Err(); err != nil {
		fmt.Println("err: ", err.Error())
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "All chats"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	l.KeyMap = list.KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}

	ta := textarea.New()
	ta.Placeholder = "Create a channel..."
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(40)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := Model{list: l, textarea: ta, adding: false}

	return m
}
