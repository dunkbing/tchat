package channel

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dunkbing/tchat/redis"
	"io"
)

const listHeight = 14
const channelPrefix = "channel"

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
	list     list.Model
	textarea textarea.Model
	adding   bool
}

func (m Model) Init() tea.Cmd {
	//return textarea.Blink
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var listCmd tea.Cmd
	var taCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	m.textarea, taCmd = m.textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEsc:
			if m.list.FilterState() == list.Filtering {
				return m, tea.Batch(listCmd, taCmd)
			}
			if m.adding {
				m.adding = false
				m.textarea.Reset()
			}
			return m, nil
		case tea.KeyCtrlA:
			m.adding = true
			m.textarea.Focus()
		}
	}

	return m, tea.Batch(listCmd, taCmd)
}

func (m Model) View() string {
	if !m.adding {
		return m.list.View()
	}
	return fmt.Sprintf("%s\n\n%s", m.list.View(), m.textarea.View())
}

func New() Model {
	ctx := context.Background()
	iter := redis.Client.Scan(ctx, 0, fmt.Sprintf("%s*", channelPrefix), 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("Keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Println("err: ", err.Error())
	}
	//listKeys := newListKeyMap()
	items := []list.Item{
		Item("Ramen"),
		Item("Tomato Soup"),
		Item("Hamburgers"),
		Item("Cheeseburgers"),
		Item("Currywurst"),
		Item("Okonomiyaki"),
		Item("Pasta"),
		Item("Fillet Mignon"),
		Item("Caviar"),
		Item("Just Wine"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	//l.AdditionalFullHelpKeys = func() []key.Binding {
	//    return []key.Binding{
	//        listKeys.toggleSpinner,
	//        listKeys.insertItem,
	//        listKeys.toggleTitleBar,
	//        listKeys.toggleStatusBar,
	//        listKeys.togglePagination,
	//        listKeys.toggleHelpMenu,
	//    }
	//}
	//l.KeyMap = list.KeyMap{
	//    CursorUp: key.NewBinding(
	//        key.WithKeys("up"),
	//        key.WithHelp("↑", "up"),
	//    ),
	//    CursorDown: key.NewBinding(
	//        key.WithKeys("down"),
	//        key.WithHelp("↓", "down"),
	//    ),
	//    PrevPage: key.NewBinding(
	//        key.WithKeys("left"),
	//    ),
	//    NextPage: key.NewBinding(
	//        key.WithKeys("right"),
	//    ),
	//    Quit: key.NewBinding(
	//        key.WithKeys("ctrl+c"),
	//    ),
	//    ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	//}

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
