package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"afk-tui/internal/models"
)

// SaveManager handles saving and loading game state
type SaveManager struct {
	SavePath string
}

// NewSaveManager creates a new save manager
func NewSaveManager(saveDir string) *SaveManager {
	if saveDir == "" {
		saveDir = "."
	}
	return &SaveManager{
		SavePath: filepath.Join(saveDir, "afk-tui-save.json"),
	}
}

// Save saves the player to disk
func (sm *SaveManager) Save(player *models.Player) error {
	player.UpdateLastOnline()

	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	// Create backup if save exists
	if _, err := os.Stat(sm.SavePath); err == nil {
		backupPath := sm.SavePath + ".backup"
		os.Rename(sm.SavePath, backupPath)
	}

	// Write new save
	if err := os.WriteFile(sm.SavePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save: %w", err)
	}

	return nil
}

// Load loads the player from disk
func (sm *SaveManager) Load() (*models.Player, error) {
	data, err := os.ReadFile(sm.SavePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save file not found")
		}
		return nil, fmt.Errorf("failed to read save: %w", err)
	}

	var player models.Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, fmt.Errorf("failed to unmarshal save: %w", err)
	}

	// Initialize ActivityLog if nil (for old saves)
	if player.ActivityLog == nil {
		player.ActivityLog = models.NewActivityLog()
		player.ActivityLog.AddEntry(models.LogTypeSystem, "Save loaded - Activity log initialized", nil)
	}

	// Initialize CombatStats if nil (for old saves)
	if player.CombatStats == nil {
		player.CombatStats = models.NewCombatStats()
	}

	// Initialize Attributes if nil (for old saves)
	if player.Attributes == nil {
		player.Attributes = models.NewCharacterAttributes()
	}

	// Restore current activity if present
	if player.CurrentActivity != nil {
		// Re-link activity to template
		player.CurrentActivity = models.NewActivity(player.CurrentActivity.ID)
	}

	return &player, nil
}

// Exists checks if save file exists
func (sm *SaveManager) Exists() bool {
	_, err := os.Stat(sm.SavePath)
	return err == nil
}

// Delete removes save file
func (sm *SaveManager) Delete() error {
	return os.Remove(sm.SavePath)
}

// OfflineProcessor handles offline progress calculation
type OfflineProcessor struct {
	MaxOfflineTime time.Duration
}

// NewOfflineProcessor creates processor with default 24h max
func NewOfflineProcessor() *OfflineProcessor {
	return &OfflineProcessor{
		MaxOfflineTime: 24 * time.Hour,
	}
}

// CalculateOfflineProgress calculates what happened while away
func (op *OfflineProcessor) CalculateOfflineProgress(player *models.Player) *OfflineResult {
	offlineDuration := time.Since(player.LastOnline)
	if offlineDuration > op.MaxOfflineTime {
		offlineDuration = op.MaxOfflineTime
	}

	if player.CurrentActivity == nil || offlineDuration < time.Second {
		return &OfflineResult{
			OfflineTime:    0,
			TicksProcessed: 0,
		}
	}

	activity := player.CurrentActivity
	activity.ApplyModifiers(player)

	// Calculate how many ticks fit in offline time
	// Assume tick rate of 1 second for simplicity
	tickRate := time.Second
	totalTicks := int(offlineDuration / tickRate)

	// Calculate actions completed
	ticksPerAction := activity.TicksRemaining
	if ticksPerAction <= 0 {
		ticksPerAction = 1
	}
	actionsCompleted := totalTicks / ticksPerAction

	// Calculate XP gained
	xpPerAction := activity.GetXP()
	totalXP := xpPerAction * int64(actionsCompleted)

	// Calculate items gained
	outputPerAction := activity.GetOutput()
	totalItems := make(map[string]int)
	for itemID, qty := range outputPerAction {
		totalItems[itemID] = qty * actionsCompleted
	}

	// Apply XP and items to player
	perks := player.AddXP(activity.SkillType, totalXP)

	var failedItems []string
	for itemID, qty := range totalItems {
		item := models.NewItem(itemID, "", qty)
		if !player.Inventory.AddItem(item) {
			failedItems = append(failedItems, itemID)
		}
	}

	return &OfflineResult{
		OfflineTime:      offlineDuration,
		TicksProcessed:   totalTicks,
		ActionsCompleted: actionsCompleted,
		XPGained:         totalXP,
		ItemsGained:      totalItems,
		PerksUnlocked:    perks,
		FailedItems:      failedItems,
		ActivityName:     activity.Name,
		SkillName:        models.SkillNames[activity.SkillType],
		SkillType:        activity.SkillType,
	}
}

// OfflineResult contains offline calculation results
type OfflineResult struct {
	OfflineTime      time.Duration
	TicksProcessed   int
	ActionsCompleted int
	XPGained         int64
	ItemsGained      map[string]int
	PerksUnlocked    []models.Perk
	FailedItems      []string
	ActivityName     string
	SkillName        string
	SkillType        models.SkillType
}

// String returns formatted offline summary
func (or *OfflineResult) String() string {
	if or.OfflineTime == 0 {
		return "Welcome back!"
	}

	hours := int(or.OfflineTime.Hours())
	minutes := int(or.OfflineTime.Minutes()) % 60

	summary := fmt.Sprintf("Welcome back! You were away for ")
	if hours > 0 {
		summary += fmt.Sprintf("%dh ", hours)
	}
	summary += fmt.Sprintf("%dm\n\n", minutes)

	summary += fmt.Sprintf("While you were away:\n")
	summary += fmt.Sprintf("  Activity: %s (%s)\n", or.ActivityName, or.SkillName)
	summary += fmt.Sprintf("  Actions: %d\n", or.ActionsCompleted)
	summary += fmt.Sprintf("  XP Gained: %d\n", or.XPGained)

	if len(or.ItemsGained) > 0 {
		summary += "  Items Gained:\n"
		for itemID, qty := range or.ItemsGained {
			item := models.GetItemTemplate(itemID)
			name := itemID
			if item != nil {
				name = item.Name
			}
			summary += fmt.Sprintf("    - %d %s\n", qty, name)
		}
	}

	if len(or.PerksUnlocked) > 0 {
		summary += "  Perks Unlocked:\n"
		for _, perk := range or.PerksUnlocked {
			summary += fmt.Sprintf("    - %s\n", perk.Name)
		}
	}

	if len(or.FailedItems) > 0 {
		summary += "  (Inventory was full for some items)\n"
	}

	return summary
}
