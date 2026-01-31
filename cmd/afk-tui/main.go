package main

import (
	"afk-tui/internal/data"
	"afk-tui/internal/engine"
	"afk-tui/internal/models"
	"afk-tui/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// GameWrapper wraps the engine model to provide the View method
type GameWrapper struct {
	model *engine.Model
}

// NewGameWrapper creates a new wrapper
func NewGameWrapper(player *models.Player, saveManager *data.SaveManager) *GameWrapper {
	return &GameWrapper{
		model: engine.NewModel(player, saveManager),
	}
}

// Init implements tea.Model
func (w *GameWrapper) Init() tea.Cmd {
	return w.model.Init()
}

// Update implements tea.Model
func (w *GameWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := w.model.Update(msg)
	w.model = newModel
	return w, cmd
}

// View implements tea.Model
func (w *GameWrapper) View() string {
	return ui.View(w.model)
}

func main() {
	// Initialize save manager
	saveManager := data.NewSaveManager("")

	// Try to load existing save
	var player *models.Player
	var err error

	if saveManager.Exists() {
		player, err = saveManager.Load()
		if err != nil {
			fmt.Printf("Warning: Failed to load save: %v\n", err)
			fmt.Println("Starting new game...")
			player = models.NewPlayer("Adventurer")
		} else {
			fmt.Println("Loaded existing save!")
		}
	} else {
		fmt.Println("Welcome to AFK-TUI!")
		fmt.Println("Starting new game...")
		player = models.NewPlayer("Adventurer")
	}

	// Create game wrapper
	game := NewGameWrapper(player, saveManager)

	// Configure Bubble Tea program
	p := tea.NewProgram(
		game,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}

	// Save on exit
	if err := saveManager.Save(player); err != nil {
		fmt.Printf("Error saving game: %v\n", err)
	} else {
		fmt.Println("\nGame saved successfully!")
	}
}
