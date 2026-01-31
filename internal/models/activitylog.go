package models

import (
	"fmt"
	"time"
)

// LogEntry represents a single log entry for the activity log
type LogEntry struct {
	Timestamp    time.Time              `json:"timestamp"`
	Type         LogType                `json:"type"`
	Message      string                 `json:"message"`
	SkillType    SkillType              `json:"skill_type,omitempty"`    // For XP logs
	XPAmount     int64                  `json:"xp_amount,omitempty"`     // XP gained
	SkillLevel   int                    `json:"skill_level,omitempty"`   // Current level after gain
	ItemName     string                 `json:"item_name,omitempty"`     // For item logs
	ItemQuantity int                    `json:"item_quantity,omitempty"` // Items gained
	ActivityName string                 `json:"activity_name,omitempty"` // Activity name
	IsCombined   bool                   `json:"is_combined,omitempty"`   // XP + Item combined
	Details      map[string]interface{} `json:"details,omitempty"`
}

// LogType categorizes log entries
type LogType string

const (
	LogTypeXP       LogType = "XP"
	LogTypeItem     LogType = "ITEM"
	LogTypeLevelUp  LogType = "LEVEL"
	LogTypePerk     LogType = "PERK"
	LogTypeActivity LogType = "ACTIVITY"
	LogTypeSystem   LogType = "SYSTEM"
	LogTypeSell     LogType = "SELL"
)

// ActivityLog tracks all game events
type ActivityLog struct {
	Entries        []LogEntry `json:"entries"`
	MaxEntries     int        `json:"max_entries"`
	ScrollPos      int        `json:"scroll_pos"`
	PendingXPEntry *LogEntry  `json:"-"` // For combining XP + Item
}

// NewActivityLog creates a new activity log
func NewActivityLog() *ActivityLog {
	return &ActivityLog{
		Entries:    make([]LogEntry, 0),
		MaxEntries: 1000,
		ScrollPos:  0,
	}
}

// AddEntry adds a new log entry
func (al *ActivityLog) AddEntry(entryType LogType, message string, details map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Type:      entryType,
		Message:   message,
		Details:   details,
	}

	al.Entries = append(al.Entries, entry)

	// Trim if exceeds max
	if len(al.Entries) > al.MaxEntries {
		al.Entries = al.Entries[len(al.Entries)-al.MaxEntries:]
	}

	// Reset scroll to bottom
	al.ScrollPos = len(al.Entries) - 1
	if al.ScrollPos < 0 {
		al.ScrollPos = 0
	}
}

// StartXPEntry starts a pending XP entry (to be combined with items)
func (al *ActivityLog) StartXPEntry(skill SkillType, amount int64, level int) {
	al.PendingXPEntry = &LogEntry{
		Timestamp:  time.Now(),
		Type:       LogTypeXP,
		SkillType:  skill,
		XPAmount:   amount,
		SkillLevel: level,
		IsCombined: false,
	}
}

// AddItemToPending adds an item to the pending XP entry and finalizes it
func (al *ActivityLog) AddItemToPending(itemName string, quantity int) {
	if al.PendingXPEntry == nil {
		// No pending XP, just log the item
		al.AddItemLog(itemName, quantity, "")
		return
	}

	// Combine XP + Item into single entry
	al.PendingXPEntry.ItemName = itemName
	al.PendingXPEntry.ItemQuantity = quantity
	al.PendingXPEntry.IsCombined = true

	// Create combined message
	al.PendingXPEntry.Message = fmt.Sprintf("%s +%d%% (+%d %s)",
		SkillNames[al.PendingXPEntry.SkillType],
		al.PendingXPEntry.XPAmount,
		quantity,
		itemName)

	al.Entries = append(al.Entries, *al.PendingXPEntry)
	al.PendingXPEntry = nil

	// Trim and scroll
	if len(al.Entries) > al.MaxEntries {
		al.Entries = al.Entries[len(al.Entries)-al.MaxEntries:]
	}
	al.ScrollPos = len(al.Entries) - 1
}

// FinalizePendingXP adds the pending XP entry without items
func (al *ActivityLog) FinalizePendingXP() {
	if al.PendingXPEntry == nil {
		return
	}

	// Just XP, no items
	al.PendingXPEntry.Message = fmt.Sprintf("%s +%d%%",
		SkillNames[al.PendingXPEntry.SkillType],
		al.PendingXPEntry.XPAmount)

	al.Entries = append(al.Entries, *al.PendingXPEntry)
	al.PendingXPEntry = nil

	if len(al.Entries) > al.MaxEntries {
		al.Entries = al.Entries[len(al.Entries)-al.MaxEntries:]
	}
	al.ScrollPos = len(al.Entries) - 1
}

// AddXPLog logs XP gain (legacy, use StartXPEntry + FinalizePendingXP or AddItemToPending)
func (al *ActivityLog) AddXPLog(skill SkillType, amount int64, newTotal int64, level int) {
	al.AddEntry(LogTypeXP, fmt.Sprintf("%s +%d%% (Lv.%d)", SkillNames[skill], amount, level), map[string]interface{}{
		"skill":  skill,
		"amount": amount,
		"total":  newTotal,
		"level":  level,
	})
}

// AddItemLog logs item gain (legacy)
func (al *ActivityLog) AddItemLog(itemName string, quantity int, itemID string) {
	al.AddEntry(LogTypeItem, fmt.Sprintf("+%d %s", quantity, itemName), map[string]interface{}{
		"item":     itemID,
		"quantity": quantity,
		"name":     itemName,
	})
}

// AddSellLog logs item selling
func (al *ActivityLog) AddSellLog(itemName string, quantity int, gold int64) {
	al.AddEntry(LogTypeSell, fmt.Sprintf("Sold %dx %s for %d gold", quantity, itemName, gold), map[string]interface{}{
		"item":     itemName,
		"quantity": quantity,
		"gold":     gold,
	})
}

// AddLevelUpLog logs level up
func (al *ActivityLog) AddLevelUpLog(skill SkillType, newLevel int) {
	al.AddEntry(LogTypeLevelUp, fmt.Sprintf("🎉 %s Level %d!", SkillNames[skill], newLevel), map[string]interface{}{
		"skill":     skill,
		"new_level": newLevel,
	})
}

// AddPerkLog logs perk unlock
func (al *ActivityLog) AddPerkLog(perkName string, skill SkillType) {
	al.AddEntry(LogTypePerk, fmt.Sprintf("✨ Perk: %s", perkName), map[string]interface{}{
		"perk":  perkName,
		"skill": skill,
	})
}

// AddActivityLog logs activity start/stop
func (al *ActivityLog) AddActivityLog(activityName string, started bool) {
	action := "Stopped"
	if started {
		action = "Started"
	}
	al.AddEntry(LogTypeActivity, fmt.Sprintf("%s %s", action, activityName), map[string]interface{}{
		"activity": activityName,
		"started":  started,
	})
}

// GetRecentEntries returns the most recent N entries
func (al *ActivityLog) GetRecentEntries(n int) []LogEntry {
	if len(al.Entries) == 0 {
		return []LogEntry{}
	}

	start := len(al.Entries) - n
	if start < 0 {
		start = 0
	}

	return al.Entries[start:]
}

// GetVisibleEntries returns entries for the scrollable view
func (al *ActivityLog) GetVisibleEntries(start, count int) []LogEntry {
	if len(al.Entries) == 0 {
		return []LogEntry{}
	}

	if start < 0 {
		start = 0
	}

	end := start + count
	if end > len(al.Entries) {
		end = len(al.Entries)
	}

	if start >= end {
		return []LogEntry{}
	}

	return al.Entries[start:end]
}

// ScrollUp moves the view up
func (al *ActivityLog) ScrollUp(amount int) {
	al.ScrollPos -= amount
	if al.ScrollPos < 0 {
		al.ScrollPos = 0
	}
}

// ScrollDown moves the view down
func (al *ActivityLog) ScrollDown(amount int) {
	al.ScrollPos += amount
	if al.ScrollPos >= len(al.Entries) {
		al.ScrollPos = len(al.Entries) - 1
	}
	if al.ScrollPos < 0 {
		al.ScrollPos = 0
	}
}

// ScrollToBottom jumps to most recent
func (al *ActivityLog) ScrollToBottom() {
	al.ScrollPos = len(al.Entries) - 1
	if al.ScrollPos < 0 {
		al.ScrollPos = 0
	}
}

// ScrollToTop jumps to oldest
func (al *ActivityLog) ScrollToTop() {
	al.ScrollPos = 0
}

// GetEntryCount returns total entries
func (al *ActivityLog) GetEntryCount() int {
	return len(al.Entries)
}
