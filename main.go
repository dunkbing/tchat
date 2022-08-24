package main

import (
	"fmt"
	"github.com/dunkbing/tchat/models/app"
	"github.com/dunkbing/tchat/redis"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//m := chat.New()
	redis.Init()
	m := app.New()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

//package main
//
//import (
//    "fmt"
//    "github.com/charmbracelet/bubbles/textarea"
//    "io"
//    "os"
//
//    "github.com/charmbracelet/bubbles/list"
//    tea "github.com/charmbracelet/bubbletea"
//    "github.com/charmbracelet/lipgloss"
//)
//
//const listHeight = 14
//
//var (
//    titleStyle        = lipgloss.NewStyle().MarginLeft(2)
//    itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
//    selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
//    paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
//    helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
//    quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
//)
//
//type item string
//
//func (i item) FilterValue() string { return string(i) }
//
//type itemDelegate struct{}
//
//func (d itemDelegate) Height() int                             { return 1 }
//func (d itemDelegate) Spacing() int                            { return 0 }
//func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
//func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
//    i, ok := listItem.(item)
//    if !ok {
//        return
//    }
//
//    str := fmt.Sprintf("%d. %s", index+1, i)
//
//    fn := itemStyle.Render
//    if index == m.Index() {
//        fn = func(s string) string {
//            return selectedItemStyle.Render("> " + s)
//        }
//    }
//
//    _, _ = fmt.Fprintf(w, fn(str))
//}
//
//type model struct {
//    list     list.Model
//    items    []item
//    choice   string
//    textarea textarea.Model
//}
//
//func (m model) Init() tea.Cmd {
//    return nil
//}
//
//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//    switch msg := msg.(type) {
//    case tea.WindowSizeMsg:
//        m.list.SetWidth(msg.Width)
//        return m, nil
//
//    case tea.KeyMsg:
//        switch msg.Type {
//        case tea.KeyCtrlC:
//            return m, tea.Quit
//
//        case tea.KeyEnter:
//            i, ok := m.list.SelectedItem().(item)
//            if ok {
//                m.choice = string(i)
//            }
//            return m, tea.Quit
//        }
//    }
//
//    var cmd tea.Cmd
//    m.list, cmd = m.list.Update(msg)
//    return m, cmd
//}
//
//func (m model) View() string {
//    if m.choice != "" {
//        return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
//    }
//    return "\n" + m.list.View() + "\n" + m.textarea.View()
//}
//
//func New() model {
//    items := []list.Item{
//        item("Ramen"),
//        item("Tomato Soup"),
//        item("Hamburgers"),
//        item("Cheeseburgers"),
//        item("Currywurst"),
//        item("Okonomiyaki"),
//        item("Pasta"),
//        item("Fillet Mignon"),
//        item("Caviar"),
//        item("Just Wine"),
//    }
//
//    const defaultWidth = 20
//
//    l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
//    l.Title = "What do you want for dinner?"
//    l.SetShowStatusBar(false)
//    l.SetFilteringEnabled(true)
//    l.Styles.Title = titleStyle
//    l.Styles.PaginationStyle = paginationStyle
//    l.Styles.HelpStyle = helpStyle
//
//    ta := textarea.New()
//    ta.Placeholder = "Find or create a channel..."
//    //ta.Focus()
//    ta.Prompt = "┃ "
//    ta.CharLimit = 280
//    ta.SetWidth(40)
//    ta.SetHeight(3)
//    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
//    ta.ShowLineNumbers = false
//    ta.KeyMap.InsertNewline.SetEnabled(false)
//
//    return model{list: l, textarea: ta}
//}
//
//func main() {
//    m := New()
//
//    if err := tea.NewProgram(m).Start(); err != nil {
//        fmt.Println("Error running program:", err)
//        os.Exit(1)
//    }
//}
