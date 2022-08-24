package main

import (
	"fmt"
	"github.com/dunkbing/tchat/models/app"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//m := chat.New()
	m := app.New()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
