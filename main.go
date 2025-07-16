package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/models"
	"github.com/zigzter/chatterm/utils"
)

func main() {
	configPath := utils.SetupPath()
	f, err := tea.LogToFile(configPath+"/debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer f.Close()
	sql := db.OpenDB(configPath)
	defer sql.Close()
	db.CreateTables(sql)

	m := models.InitialRootModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Received shutdown signal, initiating graceful shutdown...")
		
		// Safe WebSocket connection cleanup
		if m.IsChatInitialized && m.Chat.WsClient != nil {
			log.Println("Closing WebSocket connection...")
			m.Chat.WsClient.SafeClose()
		} else {
			log.Println("No WebSocket client to close")
		}
		
		// Ensure all logs are flushed
		log.Println("Shutdown complete")
		os.Exit(0)
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
