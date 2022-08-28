package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dunkbing/tchat/app"
	"github.com/dunkbing/tchat/redis"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	fmt.Println("Loading")
	redis.Init()
	m := app.New()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
