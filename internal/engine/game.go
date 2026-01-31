package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"afk-tui/internal/data"
	"afk-tui/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// GameState represents the current state of the game
type GameState int

const (
	StateDashboard GameState = iota
	StateSkills
	StateInventory
	StateEquipment
	StateHelp
	StateSkillCategories
	StateActivitySelection
	StateLogView // New: Full-screen log view
	StateTraining
	StateCombat
	StateSlayerTierSelection
	StateSlayerMonsterSelection
	StateCharacterSheet
	StateNameEdit
)

// ActivityCategory represents a group of activities
type ActivityCategory struct {
	ID          string
	Name        string
	Description string
	Icon        string
	SkillType   models.SkillType
	Activities  []ActivityOption
}

// ActivityOption represents a single selectable activity
type ActivityOption struct {
	ID          string
	Name        string
	Description string
	LevelReq    int
	Input       string
	Output      string
}

// CombatEncounter tracks an active combat session
type CombatEncounter struct {
	Monster          *models.Monster
	PlayerATB        float64 // 0-100
	MonsterATB       float64 // 0-100
	IsPlayerTurn     bool
	CombatTicks      int
	DamageDealt      int
	DamageTaken      int
	LastActionResult string
}

// Model is the main game model for Bubble Tea
type Model struct {
	State            GameState
	Player           *models.Player
	SaveManager      *data.SaveManager
	OfflineProcessor *data.OfflineProcessor

	// Current view state
	SelectedSkill    models.SkillType
	SelectedCategory string
	SelectedActivity string
	CurrentMessage   string
	ShowMessage      bool
	MessageTimer     time.Time
	CursorPosition   int // For menu navigation

	// Log view state
	LogViewExpanded   bool // Spacebar toggles this
	LogScrollPosition int
	LogEntriesPerPage int

	// Combat & Training state
	SelectedTrainingType   string // strength, dexterity, defense
	SelectedSlayerTier     int    // 1-5 for monster tiers
	SelectedMonsterID      string
	CurrentCombatEncounter *CombatEncounter
	NameEditBuffer         string
	NameEditCursor         int

	// Inventory state
	InventoryState InventoryState

	// Ticks
	TickRate  time.Duration
	LastTick  time.Time
	TickCount int // For animation

	// Views
	Width  int
	Height int
}

// NewModel creates a new game model
func NewModel(player *models.Player, saveManager *data.SaveManager) *Model {
	return &Model{
		State:             StateDashboard,
		Player:            player,
		SaveManager:       saveManager,
		OfflineProcessor:  data.NewOfflineProcessor(),
		SelectedSkill:     models.SkillWoodcutting,
		TickRate:          1 * time.Second,
		LastTick:          time.Now(),
		CursorPosition:    0,
		LogViewExpanded:   false,
		LogScrollPosition: 0,
		LogEntriesPerPage: 20,
		TickCount:         0,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Process offline progress
	result := m.OfflineProcessor.CalculateOfflineProgress(m.Player)
	if result.OfflineTime > 0 {
		// Log offline progress to activity log instead of showing popup
		if m.Player.ActivityLog == nil {
			m.Player.ActivityLog = models.NewActivityLog()
		}

		// Format resume-style summary
		hours := int(result.OfflineTime.Hours())
		minutes := int(result.OfflineTime.Minutes()) % 60

		var timeStr string
		if hours > 0 {
			timeStr = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			timeStr = fmt.Sprintf("%dm", minutes)
		}

		// Create concise resume entry
		resumeMsg := fmt.Sprintf("Away for %s: %d actions, %s XP gained",
			timeStr, result.ActionsCompleted, formatNumber(result.XPGained))

		m.Player.ActivityLog.AddEntry(models.LogTypeSystem, resumeMsg, map[string]interface{}{
			"offline_time":   result.OfflineTime.String(),
			"actions":        result.ActionsCompleted,
			"xp_gained":      result.XPGained,
			"activity":       result.ActivityName,
			"skill":          result.SkillName,
			"items_gained":   result.ItemsGained,
			"perks_unlocked": len(result.PerksUnlocked),
		})

		// Also log item summary if any items were gained
		if len(result.ItemsGained) > 0 {
			var itemSummary []string
			for itemID, qty := range result.ItemsGained {
				item := models.GetItemTemplate(itemID)
				name := itemID
				if item != nil {
					name = item.Name
				}
				itemSummary = append(itemSummary, fmt.Sprintf("%d %s", qty, name))
			}
			if len(itemSummary) > 0 {
				itemsMsg := "Items: " + strings.Join(itemSummary, ", ")
				m.Player.ActivityLog.AddEntry(models.LogTypeItem, itemsMsg, nil)
			}
		}

		// Log perks if any were unlocked
		for _, perk := range result.PerksUnlocked {
			m.Player.ActivityLog.AddPerkLog(perk.Name, result.SkillType)
		}
	}

	return tickCmd(m.TickRate)
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle log view scrolling first if in log view mode
		if m.LogViewExpanded {
			return m.handleLogViewInput(msg)
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case TickMsg:
		m.processTick()
		m.TickCount++
		return m, tickCmd(m.TickRate)

	case HideMessageMsg:
		m.ShowMessage = false
		return m, nil
	}

	return m, nil
}

// handleLogViewInput handles input when log view is expanded
func (m *Model) handleLogViewInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "space":
		// Collapse log view
		m.LogViewExpanded = false
		return m, nil

	case "up", "k":
		// Scroll up in log
		if m.LogScrollPosition > 0 {
			m.LogScrollPosition--
		}
		return m, nil

	case "down", "j":
		// Scroll down in log
		if m.Player.ActivityLog == nil {
			m.Player.ActivityLog = models.NewActivityLog()
		}
		maxScroll := m.Player.ActivityLog.GetEntryCount() - m.LogEntriesPerPage
		if maxScroll < 0 {
			maxScroll = 0
		}
		if m.LogScrollPosition < maxScroll {
			m.LogScrollPosition++
		}
		return m, nil

	case "pgup":
		// Page up
		m.LogScrollPosition -= m.LogEntriesPerPage
		if m.LogScrollPosition < 0 {
			m.LogScrollPosition = 0
		}
		return m, nil

	case "pgdown", " ":
		// Page down
		if m.Player.ActivityLog == nil {
			m.Player.ActivityLog = models.NewActivityLog()
		}
		maxScroll := m.Player.ActivityLog.GetEntryCount() - m.LogEntriesPerPage
		if maxScroll < 0 {
			maxScroll = 0
		}
		m.LogScrollPosition += m.LogEntriesPerPage
		if m.LogScrollPosition > maxScroll {
			m.LogScrollPosition = maxScroll
		}
		return m, nil

	case "home":
		// Go to top (oldest)
		m.LogScrollPosition = 0
		return m, nil

	case "end":
		// Go to bottom (newest)
		if m.Player.ActivityLog == nil {
			m.Player.ActivityLog = models.NewActivityLog()
		}
		maxScroll := m.Player.ActivityLog.GetEntryCount() - m.LogEntriesPerPage
		if maxScroll < 0 {
			maxScroll = 0
		}
		m.LogScrollPosition = maxScroll
		return m, nil

	case "esc", "q":
		// Exit log view
		m.LogViewExpanded = false
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (*Model, tea.Cmd) {
	// Global shortcuts first
	switch msg.String() {
	case "ctrl+c":
		m.SaveManager.Save(m.Player)
		return m, tea.Quit

	case "?", "h":
		if m.State == StateHelp {
			m.State = StateDashboard
		} else {
			m.State = StateHelp
		}
		m.CursorPosition = 0
		return m, nil

	case "d":
		m.State = StateDashboard
		m.CursorPosition = 0
		return m, nil

	case "s":
		m.State = StateSkills
		m.CursorPosition = 0
		return m, nil

	case "i":
		m.State = StateInventory
		m.CursorPosition = 0
		return m, nil

	case "e":
		m.State = StateEquipment
		m.CursorPosition = 0
		return m, nil

	case "space":
		// Toggle log view expansion
		if m.Player.CurrentActivity != nil {
			m.LogViewExpanded = !m.LogViewExpanded
			if m.LogViewExpanded {
				// Auto-scroll to bottom when opening
				if m.Player.ActivityLog == nil {
					m.Player.ActivityLog = models.NewActivityLog()
				}
				maxScroll := m.Player.ActivityLog.GetEntryCount() - m.LogEntriesPerPage
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.LogScrollPosition = maxScroll
			}
		} else {
			m.CurrentMessage = "Start an activity first to view logs!"
			m.ShowMessage = true
		}
		return m, nil

	case "ctrl+s":
		if err := m.SaveManager.Save(m.Player); err != nil {
			m.CurrentMessage = fmt.Sprintf("Save failed: %v", err)
		} else {
			m.CurrentMessage = "Game saved!"
		}
		m.ShowMessage = true
		return m, hideMessageCmd(2 * time.Second)

	case "q":
		// Only quit if not in a menu
		if m.State == StateDashboard {
			m.SaveManager.Save(m.Player)
			return m, tea.Quit
		}
		m.State = StateDashboard
		return m, nil
	}

	// State-specific handling
	switch m.State {
	case StateDashboard:
		return m.handleDashboardInput(msg)
	case StateSkills:
		return m.handleSkillsInput(msg)
	case StateSkillCategories:
		return m.handleSkillCategoriesInput(msg)
	case StateActivitySelection:
		return m.handleActivitySelectionInput(msg)
	case StateInventory:
		return m.handleInventoryInput(msg)
	case StateEquipment:
		return m.handleEquipmentInput(msg)
	case StateTraining:
		return m.handleTrainingInput(msg)
	case StateCombat, StateSlayerTierSelection, StateSlayerMonsterSelection:
		return m.handleCombatInput(msg)
	case StateCharacterSheet:
		return m.handleCharacterSheetInput(msg)
	case StateNameEdit:
		return m.handleNameEditInput(msg)
	}

	return m, nil
}

// handleDashboardInput handles dashboard shortcuts
func (m *Model) handleDashboardInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		return m.startActivity("chop_logs")
	case "2":
		return m.startActivity("mine_copper")
	case "3":
		return m.startActivity("smelt_bronze")
	case "4":
		return m.startActivity("recycle_logs")
	case "c":
		m.State = StateCharacterSheet
		return m, nil
	}
	return m, nil
}

// handleSkillsInput handles skills screen navigation with number/letter support
func (m *Model) handleSkillsInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	// Number keys for direct skill selection (1-9)
	skillOrder := []models.SkillType{
		models.SkillWoodcutting, models.SkillMining, models.SkillSmithing,
		models.SkillRecycling, models.SkillCombat, models.SkillCrafting,
		models.SkillCooking, models.SkillAgility, models.SkillThieving,
	}

	switch msg.String() {
	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		if m.CursorPosition < len(skillOrder)-1 {
			m.CursorPosition++
			m.SelectedSkill = skillOrder[m.CursorPosition]
		}
		return m, nil

	case "enter":
		if m.CursorPosition < len(skillOrder) {
			m.SelectedSkill = skillOrder[m.CursorPosition]
			// If Combat selected, go to categories (which has Training/Slayer)
			m.State = StateSkillCategories
			m.CursorPosition = 0
		}
		return m, nil

	case "t":
		// Quick training shortcut - only works if combat is selected
		if m.CursorPosition < len(skillOrder) && skillOrder[m.CursorPosition] == models.SkillCombat {
			m.State = StateTraining
			m.CursorPosition = 0
			return m, nil
		}
		return m, nil

	default:
		// Number keys 1-9 for direct selection
		if len(msg.String()) == 1 {
			char := msg.String()[0]
			if char >= '1' && char <= '9' {
				index := int(char - '1')
				if index < len(skillOrder) {
					m.CursorPosition = index
					m.SelectedSkill = skillOrder[index]
					m.State = StateSkillCategories
					m.CursorPosition = 0
					return m, nil
				}
			}
		}
	}

	return m, nil
}

// handleSkillCategoriesInput handles category selection with number/letter support
func (m *Model) handleSkillCategoriesInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	categories := GetCategoriesForSkill(m.SelectedSkill)

	switch msg.String() {
	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		if m.CursorPosition < len(categories)-1 {
			m.CursorPosition++
		}
		return m, nil

	case "enter":
		if m.CursorPosition < len(categories) {
			m.SelectedCategory = categories[m.CursorPosition].ID
			m.State = StateActivitySelection
			m.CursorPosition = 0
		}
		return m, nil

	case "esc", "q":
		m.State = StateSkills
		m.CursorPosition = 0
		return m, nil

	default:
		// Number keys 1-9 for direct category selection (preferred)
		if len(msg.String()) == 1 {
			keyChar := msg.String()[0]
			if keyChar >= '1' && keyChar <= '9' {
				index := int(keyChar - '1')
				if index < len(categories) {
					m.SelectedCategory = categories[index].ID
					m.State = StateActivitySelection
					m.CursorPosition = 0
					return m, nil
				}
			}
			// Letter selection (skip global keys: s, d, i, e, c, q)
			letterChar := rune(strings.ToLower(msg.String())[0])
			if letterChar == 's' || letterChar == 'd' || letterChar == 'i' || letterChar == 'e' || letterChar == 'c' || letterChar == 'q' {
				return m, nil
			}
			for _, cat := range categories {
				if len(cat.Name) > 0 && rune(strings.ToLower(cat.Name)[0]) == letterChar {
					m.SelectedCategory = cat.ID
					m.State = StateActivitySelection
					m.CursorPosition = 0
					return m, nil
				}
			}
		}
	}

	return m, nil
}

// handleActivitySelectionInput handles activity selection with letter support
func (m *Model) handleActivitySelectionInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	categories := GetCategoriesForSkill(m.SelectedSkill)

	var currentCategory *ActivityCategory
	for i := range categories {
		if categories[i].ID == m.SelectedCategory {
			currentCategory = &categories[i]
			break
		}
	}

	if currentCategory == nil {
		m.State = StateSkillCategories
		return m, nil
	}

	activities := currentCategory.Activities
	skill := m.Player.GetSkill(m.SelectedSkill)

	switch msg.String() {
	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		if m.CursorPosition < len(activities)-1 {
			m.CursorPosition++
		}
		return m, nil

	case "enter":
		if m.CursorPosition < len(activities) {
			activityID := activities[m.CursorPosition].ID
			return m.startActivity(activityID)
		}
		return m, nil

	case "esc", "q":
		m.State = StateSkillCategories
		m.CursorPosition = 0
		return m, nil

	default:
		// Letter-based activity selection (skip global navigation keys)
		if len(msg.String()) == 1 {
			char := rune(strings.ToLower(msg.String())[0])
			// Skip global navigation keys: s, d, i, e, c, q
			if char == 's' || char == 'd' || char == 'i' || char == 'e' || char == 'c' || char == 'q' {
				return m, nil
			}
			for _, activity := range activities {
				if len(activity.Name) > 0 && rune(strings.ToLower(activity.Name)[0]) == char {
					if skill.Level >= activity.LevelReq {
						return m.startActivity(activity.ID)
					}
				}
			}
		}
	}

	return m, nil
}

// handleInventoryInput handles inventory screen with sell mode
func (m *Model) handleInventoryInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	msgStr := msg.String()

	// Handle sell mode input
	if m.InventoryState.IsSellMode {
		message, completed, _ := HandleInventoryInput(msgStr, &m.InventoryState, m.Player)

		if message != "" {
			m.CurrentMessage = message
			m.ShowMessage = true
			if completed {
				return m, hideMessageCmd(2 * time.Second)
			}
		}
		return m, nil
	}

	// Normal inventory mode
	switch msgStr {
	case "esc", "q":
		m.State = StateDashboard
		return m, nil

	case "v":
		// Enter sell mode (v = vend/sell)
		m.InventoryState.IsSellMode = true
		m.InventoryState.SelectedItem = 0
		m.InventoryState.QuantityToSell = 0
		m.CurrentMessage = "Sell mode: Enter item number"
		m.ShowMessage = true
		return m, hideMessageCmd(3 * time.Second)
	}

	// Handle sell mode special keys
	if m.InventoryState.IsSellMode && m.InventoryState.SelectedItem > 0 {
		switch msgStr {
		case "enter":
			// Show confirmation
			if m.InventoryState.QuantityToSell <= 0 {
				m.InventoryState.QuantityToSell = m.InventoryState.MaxQuantity
				item := m.Player.Inventory.Items[m.InventoryState.SelectedItem-1]
				m.InventoryState.GoldValue = item.Value * int64(m.InventoryState.QuantityToSell)
			}
			m.InventoryState.ShowConfirmation = true
			return m, nil

		case "backspace":
			// Clear selection
			if m.InventoryState.QuantityToSell > 0 {
				m.InventoryState.QuantityToSell = 0
			} else {
				m.InventoryState.SelectedItem = 0
				m.InventoryState.ItemName = ""
				m.InventoryState.ItemID = ""
			}
			return m, nil

		case "max":
			// Set to max quantity
			if m.InventoryState.SelectedItem > 0 {
				m.InventoryState.QuantityToSell = m.InventoryState.MaxQuantity
				item := m.Player.Inventory.Items[m.InventoryState.SelectedItem-1]
				m.InventoryState.GoldValue = item.Value * int64(m.InventoryState.QuantityToSell)
			}
			return m, nil
		}
	}

	// Number keys for quick sell selection
	if msgStr >= "0" && msgStr <= "9" {
		// Auto-enter sell mode and select item
		if !m.InventoryState.IsSellMode {
			m.InventoryState.IsSellMode = true
			m.InventoryState.SelectedItem = 0
			m.InventoryState.QuantityToSell = 0
		}

		num, _ := strconv.Atoi(msgStr)
		if m.InventoryState.SelectedItem == 0 {
			// Selecting item by number
			m.InventoryState.SelectedItem = num
			if m.InventoryState.SelectedItem > 0 && m.InventoryState.SelectedItem <= len(m.Player.Inventory.Items) {
				item := m.Player.Inventory.Items[m.InventoryState.SelectedItem-1]
				m.InventoryState.ItemName = item.Name
				m.InventoryState.ItemID = item.ID
				m.InventoryState.MaxQuantity = item.Quantity
				m.InventoryState.QuantityToSell = item.Quantity // Default to max
				m.InventoryState.GoldValue = item.Value * int64(item.Quantity)
			}
		} else {
			// Building quantity number
			m.InventoryState.QuantityToSell = m.InventoryState.QuantityToSell*10 + num
			if m.InventoryState.QuantityToSell > m.InventoryState.MaxQuantity {
				m.InventoryState.QuantityToSell = m.InventoryState.MaxQuantity
			}
			item := m.Player.Inventory.Items[m.InventoryState.SelectedItem-1]
			m.InventoryState.GoldValue = item.Value * int64(m.InventoryState.QuantityToSell)
		}
		return m, nil
	}

	return m, nil
}

// handleEquipmentInput handles equipment screen
func (m *Model) handleEquipmentInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State = StateDashboard
		return m, nil
	}
	return m, nil
}

// handleTrainingInput handles training screen input
func (m *Model) handleTrainingInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.State = StateSkills
		m.CursorPosition = 0
		return m, nil

	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		trainingOptions := getTrainingOptions()
		if m.CursorPosition < len(trainingOptions)-1 {
			m.CursorPosition++
		}
		return m, nil

	case "enter":
		trainingOptions := getTrainingOptions()
		if m.CursorPosition < len(trainingOptions) {
			return m.startTraining(trainingOptions[m.CursorPosition].ID)
		}
		return m, nil

	default:
		// Letter shortcuts for training types (skip global keys)
		if len(msg.String()) == 1 {
			char := strings.ToLower(msg.String())[0]
			// Skip global navigation keys: s, d, i, e, c, q
			if char == 's' || char == 'd' || char == 'i' || char == 'e' || char == 'c' || char == 'q' {
				return m, nil
			}
			trainingOptions := getTrainingOptions()
			for _, opt := range trainingOptions {
				if len(opt.Name) > 0 && strings.ToLower(opt.Name)[0] == char {
					return m.startTraining(opt.ID)
				}
			}
		}
	}

	return m, nil
}

// handleCombatInput handles combat/slayer screen input
func (m *Model) handleCombatInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch m.State {
	case StateSlayerTierSelection:
		return m.handleSlayerTierInput(msg)
	case StateSlayerMonsterSelection:
		return m.handleSlayerMonsterInput(msg)
	case StateCombat:
		return m.handleActiveCombatInput(msg)
	}
	return m, nil
}

// handleSlayerTierInput handles monster tier selection
func (m *Model) handleSlayerTierInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.State = StateSkills
		m.CursorPosition = 0
		return m, nil

	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		if m.CursorPosition < 4 {
			m.CursorPosition++
		}
		return m, nil

	case "enter":
		m.SelectedSlayerTier = m.CursorPosition + 1
		m.State = StateSlayerMonsterSelection
		m.CursorPosition = 0
		return m, nil

	default:
		// Number keys 1-5 for tier selection
		if len(msg.String()) == 1 {
			char := msg.String()[0]
			if char >= '1' && char <= '5' {
				m.SelectedSlayerTier = int(char - '0')
				m.State = StateSlayerMonsterSelection
				m.CursorPosition = 0
				return m, nil
			}
		}
	}

	return m, nil
}

// handleSlayerMonsterInput handles monster selection within a tier
func (m *Model) handleSlayerMonsterInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	monsters := getMonstersForTier(m.SelectedSlayerTier)

	switch msg.String() {
	case "esc", "q":
		m.State = StateSlayerTierSelection
		m.CursorPosition = m.SelectedSlayerTier - 1
		return m, nil

	case "up", "k":
		if m.CursorPosition > 0 {
			m.CursorPosition--
		}
		return m, nil

	case "down", "j":
		if m.CursorPosition < len(monsters)-1 {
			m.CursorPosition++
		}
		return m, nil

	case "enter":
		if m.CursorPosition < len(monsters) {
			return m.startCombat(monsters[m.CursorPosition].ID)
		}
		return m, nil

	default:
		// Letter shortcuts for monsters (skip global keys)
		if len(msg.String()) == 1 {
			char := strings.ToLower(msg.String())[0]
			// Skip global navigation keys: s, d, i, e, c, q
			if char == 's' || char == 'd' || char == 'i' || char == 'e' || char == 'c' || char == 'q' {
				return m, nil
			}
			for i, monster := range monsters {
				if len(monster.Name) > 0 && strings.ToLower(monster.Name)[0] == char {
					return m.startCombat(monster.ID)
					m.CursorPosition = i
					return m, nil
				}
			}
		}
	}

	return m, nil
}

// handleActiveCombatInput handles input during active combat
func (m *Model) handleActiveCombatInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		// Flee combat
		m.CurrentCombatEncounter = nil
		m.State = StateSlayerMonsterSelection
		m.CurrentMessage = "You fled from combat!"
		m.ShowMessage = true
		return m, hideMessageCmd(2 * time.Second)
	}

	return m, nil
}

// handleCharacterSheetInput handles character sheet input
func (m *Model) handleCharacterSheetInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "d":
		m.State = StateDashboard
		return m, nil

	case "n":
		// Start name editing
		m.State = StateNameEdit
		m.NameEditBuffer = m.Player.Name
		m.NameEditCursor = len(m.NameEditBuffer)
		return m, nil
	}

	return m, nil
}

// handleNameEditInput handles name editing mode
func (m *Model) handleNameEditInput(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.State = StateCharacterSheet
		m.NameEditBuffer = ""
		return m, nil

	case tea.KeyEnter:
		// Save new name
		if len(m.NameEditBuffer) > 0 && len(m.NameEditBuffer) <= 20 {
			m.Player.Name = m.NameEditBuffer
			m.CurrentMessage = fmt.Sprintf("Name changed to: %s", m.Player.Name)
		} else {
			m.CurrentMessage = "Name must be 1-20 characters"
		}
		m.ShowMessage = true
		m.State = StateCharacterSheet
		m.NameEditBuffer = ""
		return m, hideMessageCmd(2 * time.Second)

	case tea.KeyBackspace:
		if m.NameEditCursor > 0 {
			m.NameEditBuffer = m.NameEditBuffer[:m.NameEditCursor-1] + m.NameEditBuffer[m.NameEditCursor:]
			m.NameEditCursor--
		}
		return m, nil

	case tea.KeyLeft:
		if m.NameEditCursor > 0 {
			m.NameEditCursor--
		}
		return m, nil

	case tea.KeyRight:
		if m.NameEditCursor < len(m.NameEditBuffer) {
			m.NameEditCursor++
		}
		return m, nil

	case tea.KeyRunes:
		// Add character
		if len(m.NameEditBuffer) < 20 && len(msg.Runes) > 0 {
			m.NameEditBuffer = m.NameEditBuffer[:m.NameEditCursor] + string(msg.Runes) + m.NameEditBuffer[m.NameEditCursor:]
			m.NameEditCursor += len(msg.Runes)
		}
		return m, nil
	}

	return m, nil
}

// getTrainingOptions returns available training activities
func getTrainingOptions() []ActivityOption {
	return []ActivityOption{
		{ID: "strength_training", Name: "Strength", Description: "Train strength at the training dummy", LevelReq: 1},
		{ID: "dexterity_training", Name: "Dexterity", Description: "Train dexterity on the agility course", LevelReq: 1},
		{ID: "defense_training", Name: "Defense", Description: "Train defense with shield drills", LevelReq: 1},
	}
}

// getMonstersForTier returns monsters for a given tier
func getMonstersForTier(tier int) []*models.Monster {
	var minLevel, maxLevel int
	switch tier {
	case 1:
		minLevel, maxLevel = 1, 10
	case 2:
		minLevel, maxLevel = 10, 30
	case 3:
		minLevel, maxLevel = 30, 60
	case 4:
		minLevel, maxLevel = 60, 90
	case 5:
		minLevel, maxLevel = 90, 120
	default:
		minLevel, maxLevel = 1, 10
	}

	return models.Monsters.GetMonstersByLevelRange(minLevel, maxLevel)
}

// startTraining starts a training activity
func (m *Model) startTraining(trainingID string) (*Model, tea.Cmd) {
	// Check if training activity exists
	if _, ok := models.ActivityDatabase[trainingID]; !ok {
		m.CurrentMessage = fmt.Sprintf("Unknown training: %s", trainingID)
		m.ShowMessage = true
		return m, hideMessageCmd(2 * time.Second)
	}

	return m.startActivity(trainingID)
}

// startCombat starts combat with a monster
func (m *Model) startCombat(monsterID string) (*Model, tea.Cmd) {
	monster := models.Monsters.GetMonster(monsterID)
	if monster == nil {
		m.CurrentMessage = fmt.Sprintf("Unknown monster: %s", monsterID)
		m.ShowMessage = true
		return m, hideMessageCmd(2 * time.Second)
	}

	// Initialize combat encounter
	m.CurrentCombatEncounter = &CombatEncounter{
		Monster:      monster,
		PlayerATB:    0,
		MonsterATB:   0,
		IsPlayerTurn: true,
	}
	m.SelectedMonsterID = monsterID
	m.State = StateCombat
	m.CurrentMessage = fmt.Sprintf("Combat started: %s!", monster.Name)
	m.ShowMessage = true

	return m, hideMessageCmd(2 * time.Second)
}

// startActivity starts a new activity
func (m *Model) startActivity(activityID string) (*Model, tea.Cmd) {
	activity := models.NewActivity(activityID)
	if activity == nil {
		m.CurrentMessage = fmt.Sprintf("Unknown activity: %s", activityID)
		m.ShowMessage = true
		return m, hideMessageCmd(2 * time.Second)
	}

	// Check requirements
	if err := activity.CanDo(m.Player); err != nil {
		m.CurrentMessage = err.Error()
		m.ShowMessage = true
		return m, hideMessageCmd(3 * time.Second)
	}

	// Apply modifiers
	activity.ApplyModifiers(m.Player)

	// Set as current activity
	m.Player.CurrentActivity = activity
	m.SelectedActivity = activityID

	// Log activity start
	if m.Player.ActivityLog == nil {
		m.Player.ActivityLog = models.NewActivityLog()
	}
	m.Player.ActivityLog.AddActivityLog(activity.Name, true)

	m.CurrentMessage = fmt.Sprintf("Started: %s", activity.Name)
	m.ShowMessage = true

	// Return to dashboard to see progress
	m.State = StateDashboard

	return m, hideMessageCmd(2 * time.Second)
}

// processTick handles game tick logic
func (m *Model) processTick() {
	// Process active combat first
	if m.State == StateCombat && m.CurrentCombatEncounter != nil {
		m.processCombatTick()
		return
	}

	if m.Player.CurrentActivity == nil {
		return
	}

	activity := m.Player.CurrentActivity
	activity.ApplyModifiers(m.Player)

	// Handle training activities specially
	if activity.ID == "strength_training" || activity.ID == "dexterity_training" || activity.ID == "defense_training" {
		m.processTrainingTick(activity)
		return
	}

	// Consume required items for crafting
	if activity.Type == models.ActivityCrafting || activity.Type == models.ActivityRecycling {
		if activity.Progress == 0 {
			for itemID, qty := range activity.RequiredItems {
				if !m.Player.Inventory.HasItem(itemID, qty) {
					m.Player.CurrentActivity = nil
					if m.Player.ActivityLog == nil {
						m.Player.ActivityLog = models.NewActivityLog()
					}
					m.Player.ActivityLog.AddActivityLog(activity.Name, false)
					m.CurrentMessage = fmt.Sprintf("Ran out of %s", itemID)
					m.ShowMessage = true
					return
				}
			}

			// Consume items
			for itemID, qty := range activity.RequiredItems {
				m.Player.Inventory.RemoveItem(itemID, qty)
			}
		}
	}

	// Process tick
	completed := activity.Tick()

	if completed {
		// Action completed!
		xpGained := activity.GetXP()
		outputItems := activity.GetOutput()

		skill := m.Player.GetSkill(activity.SkillType)
		oldLevel := skill.Level

		// Add XP
		perks := m.Player.AddXP(activity.SkillType, xpGained)

		// Start combined XP log
		if m.Player.ActivityLog == nil {
			m.Player.ActivityLog = models.NewActivityLog()
		}
		m.Player.ActivityLog.StartXPEntry(activity.SkillType, xpGained, skill.Level)

		// Add items to inventory and combine with XP log
		for itemID, qty := range outputItems {
			item := models.NewItem(itemID, "", qty)
			if m.Player.Inventory.AddItem(item) {
				itemTemplate := models.GetItemTemplate(itemID)
				if itemTemplate != nil {
					m.Player.ActivityLog.AddItemToPending(itemTemplate.Name, qty)
				}
			} else {
				m.CurrentMessage = fmt.Sprintf("Inventory full! Dropped %s", item.Name)
				m.ShowMessage = true
			}
		}

		// Finalize XP log (if no items were added, this just logs XP)
		m.Player.ActivityLog.FinalizePendingXP()

		// Check for level up
		if skill.Level > oldLevel {
			if m.Player.ActivityLog == nil {
				m.Player.ActivityLog = models.NewActivityLog()
			}
			m.Player.ActivityLog.AddLevelUpLog(activity.SkillType, skill.Level)
		}

		// Check for new perks
		if len(perks) > 0 {
			for _, perk := range perks {
				if m.Player.ActivityLog == nil {
					m.Player.ActivityLog = models.NewActivityLog()
				}
				m.Player.ActivityLog.AddPerkLog(perk.Name, activity.SkillType)
				m.CurrentMessage = fmt.Sprintf("Perk Unlocked: %s!", perk.Name)
				m.ShowMessage = true
			}
		}

		// Reset for next action
		activity.Reset()
	}

	m.LastTick = time.Now()
}

// processCombatTick handles combat logic per tick
func (m *Model) processCombatTick() {
	if m.CurrentCombatEncounter == nil {
		return
	}

	encounter := m.CurrentCombatEncounter
	player := m.Player

	// Fill ATB bars based on speed
	playerSpeed := models.GetATBFill(1, player.Attributes.Dexterity.Level)
	monsterSpeed := encounter.Monster.Speed * 5 // Convert speed to ATB fill rate

	encounter.PlayerATB += playerSpeed
	encounter.MonsterATB += monsterSpeed

	// Process player turn
	if encounter.PlayerATB >= 100 {
		encounter.PlayerATB = 0
		encounter.IsPlayerTurn = true

		// Calculate damage
		damage := models.CalculateDamage(
			player.CombatStats.Attack,
			player.Attributes.Strength.Level,
			encounter.Monster.Defense,
			models.CombatStyleMelee,
			encounter.Monster.Weakness,
			encounter.Monster.Resistance,
		)

		if damage > 0 {
			encounter.Monster.Hitpoints -= damage
			encounter.DamageDealt += damage
			encounter.LastActionResult = fmt.Sprintf("You hit %s for %d damage!", encounter.Monster.Name, damage)

			// Check if monster defeated
			if encounter.Monster.Hitpoints <= 0 {
				// Monster defeated!
				encounter.Monster.Hitpoints = 0
				m.completeCombat(encounter)
				return
			}
		} else {
			encounter.LastActionResult = fmt.Sprintf("You missed %s!", encounter.Monster.Name)
		}
	}

	// Process monster turn
	if encounter.MonsterATB >= 100 {
		encounter.MonsterATB = 0
		encounter.IsPlayerTurn = false

		// Calculate monster damage
		damage := models.CalculateDamage(
			encounter.Monster.Attack,
			encounter.Monster.Strength,
			player.Attributes.Defense.Level,
			models.CombatStyleMelee,
			"",
			"",
		)

		if damage > 0 {
			player.CombatStats.Hitpoints -= damage
			encounter.DamageTaken += damage
			encounter.LastActionResult = fmt.Sprintf("%s hits you for %d damage!", encounter.Monster.Name, damage)

			// Check if player defeated
			if player.CombatStats.Hitpoints <= 0 {
				player.CombatStats.Hitpoints = 0
				m.handlePlayerDefeat()
				return
			}
		} else {
			encounter.LastActionResult = fmt.Sprintf("%s missed you!", encounter.Monster.Name)
		}
	}

	encounter.CombatTicks++
	m.LastTick = time.Now()
}

// completeCombat handles monster defeat rewards
func (m *Model) completeCombat(encounter *CombatEncounter) {
	monster := encounter.Monster
	player := m.Player

	// Award XP
	combatXP := monster.CombatXP
	slayerXP := monster.SlayerXP

	// Add Combat skill XP
	m.Player.AddXP(models.SkillCombat, combatXP)

	// Add Slayer XP (directly to combat stats)
	player.CombatStats.SlayerXP += slayerXP
	// Check for slayer level up
	for player.CombatStats.SlayerXP >= models.CalculateXPToNext(player.CombatStats.SlayerLevel) && player.CombatStats.SlayerLevel < 120 {
		player.CombatStats.SlayerXP -= models.CalculateXPToNext(player.CombatStats.SlayerLevel)
		player.CombatStats.SlayerLevel++
	}

	// Award gold
	player.Gold += monster.Gold

	// Process drops
	dropMessages := []string{}
	for _, drop := range monster.Drops {
		// Roll for drop
		if drop.AlwaysDrop || (models.RollDrop(drop.DropRate)) {
			item := models.NewItem(drop.ItemID, drop.ItemName, drop.Quantity)
			if player.Inventory.AddItem(item) {
				dropMessages = append(dropMessages, fmt.Sprintf("%dx %s", drop.Quantity, drop.ItemName))
			}
		}
	}

	// Log the kill
	if player.ActivityLog == nil {
		player.ActivityLog = models.NewActivityLog()
	}
	player.ActivityLog.AddActivityLog(fmt.Sprintf("Defeated %s", monster.Name), true)

	// Show completion message
	dropStr := ""
	if len(dropMessages) > 0 {
		dropStr = fmt.Sprintf(" Drops: %s", dropMessages)
	}
	m.CurrentMessage = fmt.Sprintf("Victory! +%s XP, +%s gold%s",
		formatNumber(combatXP), formatNumber(monster.Gold), dropStr)
	m.ShowMessage = true

	// Reset monster HP for next fight
	monster.Hitpoints = monster.MaxHP

	// Return to monster selection
	m.CurrentCombatEncounter = nil
	m.State = StateSlayerMonsterSelection
}

// handlePlayerDefeat handles when player is defeated
func (m *Model) handlePlayerDefeat() {
	player := m.Player

	// Reset player HP
	player.CombatStats.Hitpoints = player.CombatStats.MaxHitpoints

	// Clear current combat
	m.CurrentCombatEncounter = nil

	// Log the defeat
	if player.ActivityLog == nil {
		player.ActivityLog = models.NewActivityLog()
	}
	player.ActivityLog.AddActivityLog("Defeated in combat", false)

	m.CurrentMessage = "You were defeated! HP restored."
	m.ShowMessage = true
	m.State = StateSlayerMonsterSelection
}

// processTrainingTick handles training activity progression
func (m *Model) processTrainingTick(activity *models.Activity) {
	player := m.Player

	// Process the tick
	completed := activity.Tick()

	if completed {
		xpGained := activity.GetXP()

		// Add XP to appropriate attribute
		switch activity.ID {
		case "strength_training":
			if player.Attributes.Strength.AddXP(xpGained) {
				m.CurrentMessage = "Strength Level Up!"
				m.ShowMessage = true
			}
		case "dexterity_training":
			if player.Attributes.Dexterity.AddXP(xpGained) {
				m.CurrentMessage = "Dexterity Level Up!"
				m.ShowMessage = true
			}
		case "defense_training":
			if player.Attributes.Defense.AddXP(xpGained) {
				m.CurrentMessage = "Defense Level Up!"
				m.ShowMessage = true
			}
		}

		// Also add Combat skill XP
		skill := player.GetSkill(models.SkillCombat)
		oldLevel := skill.Level
		player.AddXP(models.SkillCombat, xpGained/2) // Half XP to combat skill

		if skill.Level > oldLevel {
			if player.ActivityLog == nil {
				player.ActivityLog = models.NewActivityLog()
			}
			player.ActivityLog.AddLevelUpLog(models.SkillCombat, skill.Level)
		}

		// Update derived stats
		player.CombatStats.CalculateDerivedStats(player.Attributes)

		// Log activity
		if player.ActivityLog == nil {
			player.ActivityLog = models.NewActivityLog()
		}
		player.ActivityLog.AddXPLog(models.SkillCombat, xpGained, skill.XP, skill.Level)

		// Reset activity
		activity.Reset()
	}

	m.LastTick = time.Now()
}

// formatNumber helper for combat messages
func formatNumber(n int64) string {
	if n >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(n)/1000000000)
	}
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

// GetCategoriesForSkill returns categories for a skill
func GetCategoriesForSkill(skill models.SkillType) []ActivityCategory {
	switch skill {
	case models.SkillWoodcutting:
		return []ActivityCategory{
			{
				ID: "basic", Name: "Basic", Icon: "🌲", SkillType: skill,
				Description: "Easy trees",
				Activities: []ActivityOption{
					{ID: "chop_logs", Name: "Logs", Description: "Basic wood", LevelReq: 1, Output: "1x Logs"},
					{ID: "chop_oak", Name: "Oak", Description: "Sturdy wood", LevelReq: 15, Output: "1x Oak Logs"},
					{ID: "chop_willow", Name: "Willow", Description: "Flexible wood", LevelReq: 30, Output: "1x Willow Logs"},
				},
			},
			{
				ID: "quality", Name: "Quality", Icon: "🌳", SkillType: skill,
				Description: "Better wood",
				Activities: []ActivityOption{
					{ID: "chop_maple", Name: "Maple", Description: "Quality wood", LevelReq: 45, Output: "1x Maple Logs"},
					{ID: "chop_yew", Name: "Yew", Description: "Rare wood", LevelReq: 60, Output: "1x Yew Logs"},
					{ID: "chop_magic", Name: "Magic", Description: "Enchanted wood", LevelReq: 75, Output: "1x Magic Logs"},
				},
			},
			{
				ID: "exotic", Name: "Exotic", Icon: "🎋", SkillType: skill,
				Description: "Legendary wood",
				Activities: []ActivityOption{
					{ID: "chop_teak", Name: "Teak", Description: "Tropical hardwood", LevelReq: 90, Output: "1x Teak Logs"},
					{ID: "chop_mahogany", Name: "Mahogany", Description: "Premium wood", LevelReq: 105, Output: "1x Mahogany Logs"},
				},
			},
		}

	case models.SkillMining:
		return []ActivityCategory{
			{
				ID: "basic", Name: "Basic", Icon: "⚒️", SkillType: skill,
				Description: "Beginner ores",
				Activities: []ActivityOption{
					{ID: "mine_copper", Name: "Copper", Description: "Soft, easy ore", LevelReq: 1, Output: "1x Copper Ore"},
					{ID: "mine_tin", Name: "Tin", Description: "For bronze alloy", LevelReq: 1, Output: "1x Tin Ore"},
					{ID: "mine_lead", Name: "Lead", Description: "Heavy soft metal", LevelReq: 10, Output: "1x Lead Ore"},
					{ID: "mine_zinc", Name: "Zinc", Description: "For brass making", LevelReq: 12, Output: "1x Zinc Ore"},
					{ID: "mine_iron", Name: "Iron", Description: "Strong base metal", LevelReq: 15, Output: "1x Iron Ore"},
				},
			},
			{
				ID: "intermediate", Name: "Intermediate", Icon: "⛏️", SkillType: skill,
				Description: "Mid-level ores",
				Activities: []ActivityOption{
					{ID: "mine_coal", Name: "Coal", Description: "For smelting", LevelReq: 30, Output: "1x Coal"},
					{ID: "mine_nickel", Name: "Nickel", Description: "Corrosion resistant", LevelReq: 25, Output: "1x Nickel Ore"},
					{ID: "mine_silver", Name: "Silver", Description: "Precious metal", LevelReq: 40, Output: "1x Silver Ore"},
					{ID: "mine_gold", Name: "Gold", Description: "Valuable ore", LevelReq: 50, Output: "1x Gold Ore"},
					{ID: "mine_mithril", Name: "Mithril", Description: "Lightweight", LevelReq: 65, Output: "1x Mithril Ore"},
				},
			},
			{
				ID: "advanced", Name: "Advanced", Icon: "💎", SkillType: skill,
				Description: "High-level ores",
				Activities: []ActivityOption{
					{ID: "mine_adamantite", Name: "Adamantite", Description: "Extremely tough", LevelReq: 80, Output: "1x Adamantite Ore"},
					{ID: "mine_runite", Name: "Runite", Description: "Mystical metal", LevelReq: 95, Output: "1x Runite Ore"},
					{ID: "mine_platinum", Name: "Platinum", Description: "Ultra rare!", LevelReq: 70, Output: "1x Platinum Ore"},
					{ID: "mine_obsidian", Name: "Obsidian", Description: "Volcanic glass", LevelReq: 90, Output: "1x Obsidian Ore"},
				},
			},
			{
				ID: "gems", Name: "Gems", Icon: "💍", SkillType: skill,
				Description: "Precious gems",
				Activities: []ActivityOption{
					{ID: "mine_sapphire", Name: "Sapphire", Description: "Blue gemstone", LevelReq: 20, Output: "1x Uncut Sapphire"},
					{ID: "mine_emerald", Name: "Emerald", Description: "Green gemstone", LevelReq: 35, Output: "1x Uncut Emerald"},
					{ID: "mine_ruby", Name: "Ruby", Description: "Red gemstone", LevelReq: 55, Output: "1x Uncut Ruby"},
					{ID: "mine_diamond", Name: "Diamond", Description: "Clear gemstone", LevelReq: 75, Output: "1x Uncut Diamond"},
					{ID: "mine_dragonstone", Name: "Dragonstone", Description: "MYTHICAL!", LevelReq: 100, Output: "1x Uncut Dragonstone"},
				},
			},
		}

	case models.SkillSmithing:
		return []ActivityCategory{
			{
				ID: "basic", Name: "Basic", Icon: "🔥", SkillType: skill,
				Description: "Essential bars",
				Activities: []ActivityOption{
					{ID: "smelt_bronze", Name: "Bronze", Description: "1 Cu + 1 Ti", LevelReq: 1, Input: "1 Copper + 1 Tin", Output: "1x Bronze Bar"},
					{ID: "smelt_iron", Name: "Iron", Description: "Basic iron", LevelReq: 15, Input: "1 Iron Ore", Output: "1x Iron Bar"},
					{ID: "smelt_lead", Name: "Lead", Description: "Soft metal", LevelReq: 10, Input: "1 Lead Ore", Output: "1x Lead Bar"},
					{ID: "smelt_nickel", Name: "Nickel", Description: "Resistant", LevelReq: 30, Input: "1 Nickel Ore", Output: "1x Nickel Bar"},
				},
			},
			{
				ID: "alloys", Name: "Alloys", Icon: "⚙️", SkillType: skill,
				Description: "Advanced mixtures",
				Activities: []ActivityOption{
					{ID: "smelt_steel", Name: "Steel", Description: "1 Fe + 2 Coal", LevelReq: 30, Input: "1 Iron + 2 Coal", Output: "1x Steel Bar"},
					{ID: "smelt_brass", Name: "Brass", Description: "1 Cu + 1 Zn", LevelReq: 15, Input: "1 Copper + 1 Zinc", Output: "1x Brass Bar"},
					{ID: "smelt_electrum", Name: "Electrum", Description: "1 Au + 1 Ag", LevelReq: 55, Input: "1 Gold + 1 Silver", Output: "1x Electrum Bar"},
				},
			},
			{
				ID: "precious", Name: "Precious", Icon: "💰", SkillType: skill,
				Description: "Valuable bars",
				Activities: []ActivityOption{
					{ID: "smelt_silver", Name: "Silver", Description: "Pure silver", LevelReq: 40, Input: "1 Silver Ore", Output: "1x Silver Bar"},
					{ID: "smelt_gold", Name: "Gold", Description: "Pure gold", LevelReq: 50, Input: "1 Gold Ore", Output: "1x Gold Bar"},
					{ID: "smelt_platinum", Name: "Platinum", Description: "Ultra valuable", LevelReq: 75, Input: "1 Platinum + 4 Coal", Output: "1x Platinum Bar"},
				},
			},
			{
				ID: "elite", Name: "Elite", Icon: "🗡️", SkillType: skill,
				Description: "Legendary materials",
				Activities: []ActivityOption{
					{ID: "smelt_mithril", Name: "Mithril", Description: "1 Mithril + 4 Coal", LevelReq: 65, Input: "1 Mithril Ore + 4 Coal", Output: "1x Mithril Bar"},
					{ID: "smelt_adamantite", Name: "Adamantite", Description: "1 Adamantite + 6 Coal", LevelReq: 80, Input: "1 Adamantite + 6 Coal", Output: "1x Adamantite Bar"},
					{ID: "smelt_runite", Name: "Runite", Description: "1 Runite + 8 Coal", LevelReq: 95, Input: "1 Runite Ore + 8 Coal", Output: "1x Runite Bar"},
					{ID: "smelt_obsidian", Name: "Obsidian", Description: "2 Obsidian + 2 Coal", LevelReq: 95, Input: "2 Obsidian + 2 Coal", Output: "1x Obsidian Bar"},
				},
			},
			{
				ID: "tools", Name: "Tools", Icon: "🛠️", SkillType: skill,
				Description: "Create equipment",
				Activities: []ActivityOption{
					{ID: "smith_bronze_axe", Name: "Bronze Axe", Description: "Woodcutting tool", LevelReq: 5, Input: "1 Bronze + 5 Wood Frag", Output: "Bronze Axe"},
					{ID: "smith_iron_axe", Name: "Iron Axe", Description: "Better axe", LevelReq: 20, Input: "2 Iron + 10 Wood Frag", Output: "Iron Axe"},
					{ID: "smith_steel_axe", Name: "Steel Axe", Description: "Quality axe", LevelReq: 35, Input: "2 Steel + 2 Oak Logs", Output: "Steel Axe"},
					{ID: "smith_mithril_axe", Name: "Mithril Axe", Description: "Superior axe", LevelReq: 45, Input: "2 Mithril + 2 Willow", Output: "Mithril Axe"},
					{ID: "smith_bronze_pickaxe", Name: "Bronze Pick", Description: "Mining tool", LevelReq: 5, Input: "1 Bronze + 5 Wood Frag", Output: "Bronze Pickaxe"},
					{ID: "smith_iron_pickaxe", Name: "Iron Pick", Description: "Better pick", LevelReq: 20, Input: "2 Iron + 10 Wood Frag", Output: "Iron Pickaxe"},
					{ID: "smith_steel_pickaxe", Name: "Steel Pick", Description: "Quality pick", LevelReq: 35, Input: "2 Steel + 2 Oak Logs", Output: "Steel Pickaxe"},
					{ID: "smith_mithril_pickaxe", Name: "Mithril Pick", Description: "Superior pick", LevelReq: 45, Input: "2 Mithril + 2 Willow", Output: "Mithril Pickaxe"},
				},
			},
		}

	case models.SkillRecycling:
		return []ActivityCategory{
			{
				ID: "wood", Name: "Wood", Icon: "🪵", SkillType: skill,
				Description: "Turn logs into fragments",
				Activities: []ActivityOption{
					{ID: "recycle_logs", Name: "Logs", Description: "Basic wood", LevelReq: 1, Input: "1 Logs", Output: "1x Wood Fragments"},
					{ID: "recycle_oak_logs", Name: "Oak", Description: "Oak wood", LevelReq: 15, Input: "1 Oak Logs", Output: "2x Wood Fragments"},
					{ID: "recycle_willow_logs", Name: "Willow", Description: "Willow wood", LevelReq: 30, Input: "1 Willow Logs", Output: "3x Wood Fragments"},
					{ID: "recycle_maple_logs", Name: "Maple", Description: "Maple wood", LevelReq: 45, Input: "1 Maple Logs", Output: "4x Wood Fragments"},
					{ID: "recycle_yew_logs", Name: "Yew", Description: "Yew wood", LevelReq: 60, Input: "1 Yew Logs", Output: "5x Wood Fragments"},
					{ID: "recycle_magic_logs", Name: "Magic", Description: "Magic wood", LevelReq: 75, Input: "1 Magic Logs", Output: "8x Wood Frag + Magic Essence"},
				},
			},
		}

	case models.SkillCombat:
		return []ActivityCategory{
			{
				ID: "training", Name: "Training", Icon: "💪", SkillType: skill,
				Description: "Train combat attributes",
				Activities: []ActivityOption{
					{ID: "strength_training", Name: "Strength", Description: "Train strength at dummy", LevelReq: 1},
					{ID: "dexterity_training", Name: "Dexterity", Description: "Train dexterity course", LevelReq: 1},
					{ID: "defense_training", Name: "Defense", Description: "Shield drills", LevelReq: 1},
				},
			},
			{
				ID: "slayer", Name: "Slayer", Icon: "⚔️", SkillType: skill,
				Description: "Fight monsters for XP and loot",
				Activities: []ActivityOption{
					{ID: "tier1", Name: "Tier 1", Description: "Level 1-10 monsters", LevelReq: 1},
					{ID: "tier2", Name: "Tier 2", Description: "Level 10-30 monsters", LevelReq: 10},
					{ID: "tier3", Name: "Tier 3", Description: "Level 30-60 monsters", LevelReq: 30},
					{ID: "tier4", Name: "Tier 4", Description: "Level 60-90 monsters", LevelReq: 60},
					{ID: "tier5", Name: "Tier 5", Description: "Level 90+ monsters", LevelReq: 90},
				},
			},
		}

	default:
		return []ActivityCategory{}
	}
}

// TickMsg is sent on each tick
type TickMsg struct {
	Time time.Time
}

// tickCmd returns a command that ticks every duration
func tickCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return TickMsg{Time: t}
	})
}

// HideMessageMsg hides the message
type HideMessageMsg struct{}

// hideMessageCmd returns a command to hide message after duration
func hideMessageCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return HideMessageMsg{}
	})
}
