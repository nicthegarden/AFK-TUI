package ui

import (
	"afk-tui/internal/engine"
	"afk-tui/internal/models"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	// Colors
	colorPrimary   = lipgloss.Color("#00FF00")
	colorSecondary = lipgloss.Color("#00AAAA")
	colorAccent    = lipgloss.Color("#FFAA00")
	colorWarning   = lipgloss.Color("#FFAA00")
	colorDanger    = lipgloss.Color("#FF5555")
	colorInfo      = lipgloss.Color("#00FFFF")
	colorText      = lipgloss.Color("#FFFFFF")
	colorDim       = lipgloss.Color("#888888")
	colorHighlight = lipgloss.Color("#FFFF00")
	colorDarkBg    = lipgloss.Color("#1a1a2e")
	colorPanelBg   = lipgloss.Color("#16213e")

	// Base styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Background(lipgloss.Color("#0f0f23")).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(1, 2).
			Background(colorDarkBg)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorHighlight).
			Bold(true).
			Background(lipgloss.Color("#333333"))

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorInfo).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(colorText)

	categoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	activityStyle = lipgloss.NewStyle().
			Foreground(colorText)

	lockedStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	// Hotkey styles
	hotkeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHighlight).
			Background(lipgloss.Color("#444400")).
			Padding(0, 1)

	selectedHotkeyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(colorHighlight).
				Padding(0, 1)

	// Animation and log styles
	animationStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	logPanelStyle = lipgloss.NewStyle().
			Background(colorPanelBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(0, 1)

	logEntryStyle = lipgloss.NewStyle().
			Foreground(colorText)

	logTimestampStyle = lipgloss.NewStyle().
				Foreground(colorDim)

	logXPStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff88"))

	logItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffaa00"))

	logLevelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff55ff")).
			Bold(true)

	logPerkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#55ffff")).
			Bold(true)

	// Progress styles
	progressFullStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(lipgloss.Color("#003300"))

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(colorDim)

	progressAnimStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff00")).
				Background(lipgloss.Color("#004400"))

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f0f23")).
			Foreground(colorPrimary).
			Bold(true)

	statusBarInactiveStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1a1a2e")).
				Foreground(colorDim)

	messageStyle = lipgloss.NewStyle().
			Foreground(colorHighlight).
			Background(lipgloss.Color("#444400")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorHighlight)

	// Tier colors
	tier1Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#90EE90"))
	tier2Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	tier3Style     = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDA0DD"))
	legendaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
)

// Animation frames for the activity indicator
var animationFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// View renders the current view
func View(m *engine.Model) string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// If log view is expanded, show full-screen log
	if m.LogViewExpanded {
		return renderExpandedLogView(m)
	}

	var sections []string

	// Header
	sections = append(sections, renderHeader(m))

	// Main content area (reduced height to accommodate bottom panel)
	contentHeight := m.Height - 12 // Reserve space for status bar, log panel, and footer
	if contentHeight < 10 {
		contentHeight = 10
	}

	switch m.State {
	case engine.StateDashboard:
		sections = append(sections, renderDashboard(m, contentHeight))
	case engine.StateSkills:
		sections = append(sections, renderSkillsWithHotkeys(m, contentHeight))
	case engine.StateSkillCategories:
		sections = append(sections, renderCategorySelectionWithHotkeys(m, contentHeight))
	case engine.StateActivitySelection:
		sections = append(sections, renderActivitySelectionWithHotkeys(m, contentHeight))
	case engine.StateInventory:
		sections = append(sections, renderInventory(m, contentHeight))
	case engine.StateEquipment:
		sections = append(sections, renderEquipment(m, contentHeight))
	case engine.StateHelp:
		sections = append(sections, renderHelp(m, contentHeight))
	case engine.StateTraining:
		sections = append(sections, renderTraining(m, contentHeight))
	case engine.StateCharacterSheet:
		sections = append(sections, renderCharacterSheet(m, contentHeight))
	case engine.StateNameEdit:
		sections = append(sections, renderNameEdit(m, contentHeight))
	case engine.StateSlayerTierSelection:
		sections = append(sections, renderSlayerTierSelection(m, contentHeight))
	case engine.StateSlayerMonsterSelection:
		sections = append(sections, renderSlayerMonsterSelection(m, contentHeight))
	case engine.StateCombat:
		sections = append(sections, renderCombat(m, contentHeight))
	default:
		sections = append(sections, renderDashboard(m, contentHeight))
	}

	// Status bar with progress
	sections = append(sections, renderStatusBarWithAnimation(m))

	// Ensure ActivityLog is initialized
	if m.Player.ActivityLog == nil {
		m.Player.ActivityLog = models.NewActivityLog()
	}

	// Log panel (last 3 entries)
	sections = append(sections, renderLogPanel(m))

	// Footer
	sections = append(sections, renderFooter(m))

	// Message overlay
	result := lipgloss.JoinVertical(lipgloss.Left, sections...)

	if m.ShowMessage {
		result = overlayMessage(result, m.CurrentMessage, m.Width, m.Height)
	}

	return result
}

// renderHeader renders the game header
func renderHeader(m *engine.Model) string {
	player := m.Player

	left := headerStyle.Render(" ⚔️ AFK-TUI ")

	center := fmt.Sprintf("💰 %s | ⭐ Total: %d | 👤 %s",
		formatNumber(player.Gold),
		player.GetTotalLevel(),
		player.Name)

	right := headerStyle.Render(" [?]Help ")

	availableWidth := m.Width - lipgloss.Width(left) - lipgloss.Width(right)
	centerPadding := (availableWidth - len(center)) / 2
	if centerPadding < 0 {
		centerPadding = 0
	}

	return left + strings.Repeat(" ", centerPadding) + center + strings.Repeat(" ", availableWidth-lipgloss.Width(center)-centerPadding) + right
}

// renderStatusBarWithAnimation renders the status bar with animation
func renderStatusBarWithAnimation(m *engine.Model) string {
	if m.Player.CurrentActivity == nil {
		return statusBarInactiveStyle.
			Width(m.Width).
			Render(" 📊 No Activity - Press [s] to select skill → letter to grind | [Space] for logs ")
	}

	activity := m.Player.CurrentActivity
	progress := activity.Progress

	// Get animation frame
	frame := animationFrames[m.TickCount%len(animationFrames)]

	// Create progress bar with gradient effect
	progressBar := renderAnimatedProgressBar(progress, 35, m.TickCount)

	// Format status line
	status := fmt.Sprintf(" %s %s %s | %s | %.0f%% | XP: +%d/tick | [Space] Logs",
		frame,
		activity.Name,
		getSkillIcon(activity.SkillType),
		progressBar,
		progress*100,
		activity.GetXP())

	return statusBarStyle.
		Width(m.Width).
		Render(status)
}

// renderLogPanel renders the last 3 log entries
func renderLogPanel(m *engine.Model) string {
	// Safety check - initialize if nil
	if m.Player.ActivityLog == nil {
		m.Player.ActivityLog = models.NewActivityLog()
	}

	entries := m.Player.ActivityLog.GetRecentEntries(3)

	if len(entries) == 0 {
		return logPanelStyle.
			Width(m.Width).
			Render(" 📜 Activity Log: Start grinding to see logs here... ")
	}

	var lines []string
	lines = append(lines, " 📜 Recent Activity:")

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		formatted := formatLogEntry(entry)
		lines = append(lines, "   "+formatted)
	}

	return logPanelStyle.
		Width(m.Width).
		Render(strings.Join(lines, "\n"))
}

// renderExpandedLogView renders full-screen scrollable log
func renderExpandedLogView(m *engine.Model) string {
	// Safety check - initialize if nil
	if m.Player.ActivityLog == nil {
		m.Player.ActivityLog = models.NewActivityLog()
	}

	// Header
	header := headerStyle.Render(" 📜 Activity Log History ")

	// Get visible entries based on scroll position
	visibleEntries := m.Player.ActivityLog.GetVisibleEntries(m.LogScrollPosition, m.LogEntriesPerPage)

	var lines []string

	// Show scroll position info
	totalEntries := m.Player.ActivityLog.GetEntryCount()
	scrollInfo := fmt.Sprintf(" Showing %d-%d of %d entries | Scroll: ↑/↓/PgUp/PgDn | [Space] Close ",
		m.LogScrollPosition+1,
		m.LogScrollPosition+len(visibleEntries),
		totalEntries)

	lines = append(lines, lipgloss.NewStyle().Foreground(colorDim).Render(scrollInfo))
	lines = append(lines, "")

	// Display entries
	for i, entry := range visibleEntries {
		entryNum := m.LogScrollPosition + i + 1
		formatted := formatLogEntryDetailed(entry, entryNum)
		lines = append(lines, formatted)
	}

	// Fill remaining space if needed
	remainingLines := m.LogEntriesPerPage - len(visibleEntries)
	for i := 0; i < remainingLines; i++ {
		lines = append(lines, "")
	}

	// Footer with controls
	footer := lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render(" [↑/↓] Scroll | [PgUp/PgDn] Page | [Home] Top | [End] Bottom | [Space] Close ")

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	contentBox := boxStyle.
		Width(m.Width - 4).
		Height(m.Height - 4).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Center, header, contentBox, footer)
}

// formatLogEntry formats a log entry for the small panel
func formatLogEntry(entry models.LogEntry) string {
	timeStr := entry.Timestamp.Format("15:04:05")
	timestamp := logTimestampStyle.Render(timeStr)

	switch entry.Type {
	case models.LogTypeXP:
		return fmt.Sprintf("%s %s", timestamp, logXPStyle.Render(entry.Message))
	case models.LogTypeItem:
		return fmt.Sprintf("%s %s", timestamp, logItemStyle.Render(entry.Message))
	case models.LogTypeLevelUp:
		return fmt.Sprintf("%s %s", timestamp, logLevelStyle.Render(entry.Message))
	case models.LogTypePerk:
		return fmt.Sprintf("%s %s", timestamp, logPerkStyle.Render(entry.Message))
	default:
		return fmt.Sprintf("%s %s", timestamp, logEntryStyle.Render(entry.Message))
	}
}

// formatLogEntryDetailed formats a log entry with number for full view
func formatLogEntryDetailed(entry models.LogEntry, num int) string {
	timeStr := entry.Timestamp.Format("15:04:05")
	numStr := fmt.Sprintf("[%3d]", num)

	var typeIcon string
	var style lipgloss.Style

	switch entry.Type {
	case models.LogTypeXP:
		typeIcon = "📈"
		style = logXPStyle
	case models.LogTypeItem:
		typeIcon = "📦"
		style = logItemStyle
	case models.LogTypeLevelUp:
		typeIcon = "🎉"
		style = logLevelStyle
	case models.LogTypePerk:
		typeIcon = "✨"
		style = logPerkStyle
	case models.LogTypeActivity:
		typeIcon = "⚡"
		style = logEntryStyle
	default:
		typeIcon = "📝"
		style = logEntryStyle
	}

	return fmt.Sprintf("%s %s %s %s",
		logTimestampStyle.Render(numStr),
		logTimestampStyle.Render(timeStr),
		style.Render(typeIcon),
		style.Render(entry.Message))
}

// renderAnimatedProgressBar creates an animated progress bar
func renderAnimatedProgressBar(progress float64, width int, tickCount int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}
	empty := width - filled

	// Create animated effect in the filled portion
	fullChar := "█"
	animChars := []string{"▓", "▒", "░"}
	animIndex := tickCount % len(animChars)

	var bar string
	if filled > 0 {
		// Most of the bar is solid
		solidFilled := filled - 1
		if solidFilled < 0 {
			solidFilled = 0
		}
		bar = strings.Repeat(fullChar, solidFilled)
		// Last character animates
		if filled > 0 {
			bar += animChars[animIndex]
		}
	}

	emptyChar := "░"
	bar += strings.Repeat(emptyChar, empty)

	// Color based on progress
	if progress >= 1 {
		return lipgloss.NewStyle().Foreground(colorPrimary).Render(bar)
	} else if progress > 0.7 {
		return progressAnimStyle.Render(bar)
	}

	return lipgloss.NewStyle().Foreground(colorSecondary).Render(bar)
}

// renderFooter renders the footer bar
func renderFooter(m *engine.Model) string {
	hotkeys := []string{
		"[d]Dashboard",
		"[c]Char",
		"[s]Skills",
		"[i]Inv",
		"[e]Equip",
		"[Space]Logs",
		"[Ctrl+S]Save",
		"[q]Quit",
	}

	footer := strings.Join(hotkeys, " ")

	if len(footer) > m.Width {
		footer = footer[:m.Width-3] + "..."
	}

	return lipgloss.NewStyle().
		Foreground(colorDim).
		Background(lipgloss.Color("#222222")).
		Width(m.Width).
		Render(footer)
}

// renderDashboard renders the main dashboard
func renderDashboard(m *engine.Model, height int) string {
	player := m.Player

	var skillLines []string
	skillLines = append(skillLines, labelStyle.Render("📊 Skills"))
	skillLines = append(skillLines, strings.Repeat("─", 20))

	skillOrder := []models.SkillType{
		models.SkillWoodcutting,
		models.SkillMining,
		models.SkillSmithing,
		models.SkillRecycling,
		models.SkillCrafting,
	}

	for _, skillType := range skillOrder {
		skill := player.GetSkill(skillType)
		progress := float64(skill.XP) / float64(skill.XPToNext)
		progressBar := renderProgressBar(progress, 12)

		line := fmt.Sprintf("%-12s Lv.%3d %s",
			models.SkillNames[skillType],
			skill.Level,
			progressBar)
		skillLines = append(skillLines, line)
	}

	skillsBox := boxStyle.
		Height(height - 2).
		Width(35).
		Render(lipgloss.JoinVertical(lipgloss.Left, skillLines...))

	var infoLines []string

	infoLines = append(infoLines, labelStyle.Render("⚡ Quick Navigation"))
	infoLines = append(infoLines, strings.Repeat("─", 30))
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "Press letter to grind:")
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "[s] → [w] → [l] = Chop Logs")
	infoLines = append(infoLines, "[s] → [m] → [c] = Mine Copper")
	infoLines = append(infoLines, "[s] → [m] → [i] = Smelt Iron")
	infoLines = append(infoLines, "[s] → [s] → [b] = Smelt Bronze")
	infoLines = append(infoLines, "[s] → [r] → [l] = Recycle Logs")
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, labelStyle.Render("🛡️ Equipment"))
	infoLines = append(infoLines, player.Equipment.String())

	toolPower := player.Equipment.GetToolPower()
	if toolPower > 0 {
		infoLines = append(infoLines, fmt.Sprintf("Tool Power: +%d (+%.0f%% speed)", toolPower, float64(toolPower)*5))
	}

	infoLines = append(infoLines, "")
	infoLines = append(infoLines, labelStyle.Render("📦 Storage (Unlimited)"))
	infoLines = append(infoLines, fmt.Sprintf("Items: %d", player.Inventory.Count()))

	infoBox := boxStyle.
		Height(height - 2).
		Width(m.Width - 41).
		Render(lipgloss.JoinVertical(lipgloss.Left, infoLines...))

	return lipgloss.JoinHorizontal(lipgloss.Top, skillsBox, infoBox)
}

// renderSkillsWithHotkeys renders skills with letter hotkeys
func renderSkillsWithHotkeys(m *engine.Model, height int) string {
	player := m.Player

	var lines []string
	lines = append(lines, headerStyle.Render(" 🎮 Select Skill (Press Letter) "))
	lines = append(lines, "")

	skills := []struct {
		key       rune
		name      string
		skillType models.SkillType
	}{
		{'w', "Woodcutting", models.SkillWoodcutting},
		{'m', "Mining", models.SkillMining},
		{'s', "Smithing", models.SkillSmithing},
		{'r', "Recycling", models.SkillRecycling},
		{'c', "Combat", models.SkillCombat},
		{'a', "Crafting", models.SkillCrafting},
		{'k', "Cooking", models.SkillCooking},
		{'g', "Agility", models.SkillAgility},
		{'t', "Thieving", models.SkillThieving},
	}

	for i, s := range skills {
		skill := player.GetSkill(s.skillType)
		isSelected := i == m.CursorPosition

		progress := float64(skill.XP) / float64(skill.XPToNext)
		progressBar := renderProgressBar(progress, 12)

		var hotkey string
		var line string

		if isSelected {
			hotkey = selectedHotkeyStyle.Render(string(s.key))
			line = selectedStyle.Render(fmt.Sprintf("%s %-12s Lv.%3d/120 %s",
				hotkey, s.name, skill.Level, progressBar))
		} else {
			hotkey = hotkeyStyle.Render(string(s.key))
			line = fmt.Sprintf("%s %-12s Lv.%3d/120 %s",
				hotkey, s.name, skill.Level, progressBar)
		}

		lines = append(lines, line)

		if isSelected {
			lines = append(lines, fmt.Sprintf("     XP: %s / %s",
				formatNumber(skill.XP), formatNumber(skill.XPToNext)))

			perks := models.GetAllPerksForSkill(s.skillType)
			unlockedCount := 0
			for _, perk := range perks {
				if skill.Level >= perk.LevelReq {
					unlockedCount++
				}
			}
			lines = append(lines, fmt.Sprintf("     Perks: %d/%d", unlockedCount, len(perks)))
		}
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [1-9] Select  [↑/↓] Navigate  [Enter] View  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderCategorySelectionWithHotkeys renders categories with letter hotkeys
func renderCategorySelectionWithHotkeys(m *engine.Model, height int) string {
	player := m.Player
	skill := player.GetSkill(m.SelectedSkill)

	categories := engine.GetCategoriesForSkill(m.SelectedSkill)

	var lines []string
	lines = append(lines, headerStyle.Render(fmt.Sprintf(" %s %s (Lv.%d) - Select Category ",
		getSkillIcon(m.SelectedSkill),
		models.SkillNames[m.SelectedSkill],
		skill.Level)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Press letter to select category:"))
	lines = append(lines, "")

	for i, cat := range categories {
		isSelected := i == m.CursorPosition

		hotkeyLetter := getFirstLetter(cat.Name)

		var hotkey string
		var line string

		if isSelected {
			hotkey = selectedHotkeyStyle.Render(string(hotkeyLetter))
			line = selectedStyle.Render(fmt.Sprintf("%s %s %s - %s",
				hotkey, cat.Icon, cat.Name, cat.Description))
		} else {
			hotkey = hotkeyStyle.Render(string(hotkeyLetter))
			line = fmt.Sprintf("%s %s %s %s - %s",
				hotkey, cat.Icon, cat.Name,
				dimStyle.Render(fmt.Sprintf("(%d activities)", len(cat.Activities))),
				cat.Description)
		}

		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [1-9] Select  [↑/↓] Navigate  [Enter] Confirm  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderActivitySelectionWithHotkeys renders activities with letter hotkeys
func renderActivitySelectionWithHotkeys(m *engine.Model, height int) string {
	player := m.Player
	skill := player.GetSkill(m.SelectedSkill)

	categories := engine.GetCategoriesForSkill(m.SelectedSkill)

	var currentCategory *engine.ActivityCategory
	for i := range categories {
		if categories[i].ID == m.SelectedCategory {
			currentCategory = &categories[i]
			break
		}
	}

	if currentCategory == nil {
		return boxStyle.Render("Error: Category not found")
	}

	var lines []string
	lines = append(lines, headerStyle.Render(fmt.Sprintf(" %s %s > %s ",
		getSkillIcon(m.SelectedSkill),
		models.SkillNames[m.SelectedSkill],
		currentCategory.Name)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(fmt.Sprintf("Press letter to %s:", getActionVerb(m.SelectedSkill))))
	lines = append(lines, "")

	for i, activity := range currentCategory.Activities {
		isSelected := i == m.CursorPosition
		canDo := skill.Level >= activity.LevelReq

		hotkeyLetter := getFirstLetter(activity.Name)

		var hotkey string
		var line string

		if !canDo {
			hotkey = dimStyle.Render(string(hotkeyLetter))
			line = lockedStyle.Render(fmt.Sprintf("%s %s (Lv.%d)",
				hotkey, activity.Name, activity.LevelReq))
		} else if isSelected {
			hotkey = selectedHotkeyStyle.Render(string(hotkeyLetter))
			line = selectedStyle.Render(fmt.Sprintf("%s %s", hotkey, activity.Name))
			lines = append(lines, line)
			lines = append(lines, fmt.Sprintf("       %s", activity.Description))
			if activity.Input != "" {
				lines = append(lines, fmt.Sprintf("       Input: %s", activity.Input))
			}
			lines = append(lines, fmt.Sprintf("       Output: %s", activityStyle.Render(activity.Output)))
			lines = append(lines, "")
			continue
		} else {
			hotkey = hotkeyStyle.Render(string(hotkeyLetter))
			line = fmt.Sprintf("%s %s", hotkey, activity.Name)
		}

		lines = append(lines, line)
		lines = append(lines, "")
	}

	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [1-9/letter] Start  [↑/↓] Navigate  [Enter] Start  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderInventory renders the inventory
func renderInventory(m *engine.Model, height int) string {
	player := m.Player

	var lines []string

	// Show sell mode header if active
	if m.InventoryState.IsSellMode {
		lines = append(lines, headerStyle.Render(fmt.Sprintf(" 💰 SELL MODE (%d items) ",
			player.Inventory.Count())))
		lines = append(lines, "")
		lines = append(lines, labelStyle.Render("Enter item # then quantity (or 'max')"))
		lines = append(lines, "")
	} else {
		lines = append(lines, headerStyle.Render(fmt.Sprintf(" 🎒 Inventory (%d items - Unlimited) ",
			player.Inventory.Count())))
		lines = append(lines, "")
	}

	if len(player.Inventory.Items) == 0 {
		lines = append(lines, dimStyle.Render("Inventory is empty"))
	} else {
		itemsByType := make(map[models.ItemType][]*models.Item)
		for _, item := range player.Inventory.Items {
			itemsByType[item.Type] = append(itemsByType[item.Type], item)
		}

		categories := []models.ItemType{
			models.ItemTypeResource,
			models.ItemTypeBar,
			models.ItemTypeMaterial,
			models.ItemTypeTool,
			models.ItemTypeWeapon,
		}

		itemNum := 1
		for _, itemType := range categories {
			items := itemsByType[itemType]
			if len(items) > 0 {
				lines = append(lines, categoryStyle.Render(string(itemType)))
				for _, item := range items {
					var line string
					if m.InventoryState.IsSellMode {
						// Numbered list in sell mode
						isSelected := m.InventoryState.SelectedItem == itemNum
						if isSelected {
							line = selectedStyle.Render(fmt.Sprintf("  [%2d] %-22s x%d (%s gold)",
								itemNum, item.Name, item.Quantity, formatNumber(item.Value)))
						} else {
							line = fmt.Sprintf("  [%2d] %-22s x%d (%s gold)",
								itemNum, item.Name, item.Quantity, formatNumber(item.Value))
						}
					} else {
						// Normal view
						line = fmt.Sprintf("  %-25s x%d", item.Name, item.Quantity)
						if item.IsRecyclable() {
							line += tier1Style.Render(" [R]")
						}
					}
					lines = append(lines, line)
					itemNum++
				}
				lines = append(lines, "")
			}
		}
	}

	// Show sell confirmation dialog
	if m.InventoryState.IsSellMode && m.InventoryState.ShowConfirmation {
		lines = append(lines, "")
		confirmBox := lipgloss.NewStyle().
			Background(lipgloss.Color("#333300")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorHighlight).
			Padding(1, 2).
			Render(fmt.Sprintf("CONFIRM SELL\n\nSell: %dx %s\nFor: %s gold\n\n[Y] Yes  [N] No",
				m.InventoryState.QuantityToSell,
				m.InventoryState.ItemName,
				formatNumber(m.InventoryState.GoldValue)))
		lines = append(lines, confirmBox)
	}

	// Show sell mode status
	if m.InventoryState.IsSellMode && !m.InventoryState.ShowConfirmation {
		if m.InventoryState.SelectedItem > 0 {
			lines = append(lines, "")
			lines = append(lines, labelStyle.Render(fmt.Sprintf("Selected: %s | Quantity: %d | Value: %s gold",
				m.InventoryState.ItemName,
				m.InventoryState.QuantityToSell,
				formatNumber(m.InventoryState.GoldValue))))
			lines = append(lines, dimStyle.Render("Enter: Confirm | Backspace: Clear | max: All | Esc: Cancel"))
		}
	}

	if !m.InventoryState.IsSellMode {
		lines = append(lines, dimStyle.Render(fmt.Sprintf("Total Value: %s gold",
			formatNumber(player.Inventory.GetTotalValue()))))
	}

	// Footer
	if m.InventoryState.IsSellMode {
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(colorInfo).
			Render("  [#] Select item  [qty] Enter amount  [max] All  [Enter] Confirm  [Esc] Cancel  "))
	} else {
		lines = append(lines, lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(colorInfo).
			Render("  [#] Quick sell  [v] Sell/Vend  [Esc/q] Back  "))
	}

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderEquipment renders equipment screen
func renderEquipment(m *engine.Model, height int) string {
	player := m.Player

	var lines []string
	lines = append(lines, headerStyle.Render(" 🛡️ Equipment "))
	lines = append(lines, "")

	slots := []struct {
		name string
		item *models.Item
	}{
		{"Head", player.Equipment.Head},
		{"Body", player.Equipment.Body},
		{"Legs", player.Equipment.Legs},
		{"Feet", player.Equipment.Feet},
		{"Hands", player.Equipment.Hands},
		{"Weapon", player.Equipment.Weapon},
		{"Off-hand", player.Equipment.Offhand},
	}

	for _, slot := range slots {
		if slot.item != nil {
			lines = append(lines, fmt.Sprintf("%-10s: %s", slot.name, slot.item.Name))
			if slot.item.ToolPower > 0 {
				lines = append(lines, fmt.Sprintf("           Power: +%d", slot.item.ToolPower))
			}
		} else {
			lines = append(lines, fmt.Sprintf("%-10s: %s", slot.name, dimStyle.Render("Empty")))
		}
	}

	lines = append(lines, "")
	stats := player.Equipment.GetTotalStats()
	lines = append(lines, fmt.Sprintf("Tool Power: %d (+%.0f%% speed)",
		stats["tool_power"], float64(stats["tool_power"])*5))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderHelp renders help screen
func renderHelp(m *engine.Model, height int) string {
	var lines []string
	lines = append(lines, headerStyle.Render(" ❓ Controls Help "))
	lines = append(lines, "")

	helpItems := []struct {
		key  string
		desc string
	}{
		{"Navigation", ""},
		{"s → w → l", "Skills → Woodcutting → Logs (example)"},
		{"letter keys", "Quick select by first letter"},
		{"↑/↓ or j/k", "Navigate menus"},
		{"Enter", "Confirm selection"},
		{"", ""},
		{"Activity", ""},
		{"Space", "Toggle expanded log view (while grinding)"},
		{"↑/↓/PgUp/PgDn", "Scroll in log view"},
		{"Home/End", "Jump to top/bottom of log"},
		{"", ""},
		{"Global", ""},
		{"d", "Dashboard"},
		{"c", "Character Sheet"},
		{"i", "Inventory"},
		{"e", "Equipment"},
		{"?/h", "This help"},
		{"Ctrl+S", "Save game"},
		{"q", "Save & quit"},
		{"", ""},
		{"Combat", ""},
		{"s → c → t", "Skills → Combat → Training"},
		{"s → c → s", "Skills → Combat → Slayer"},
		{"In Combat", "ATB system - auto-attacks when ready"},
		{"Esc/q", "Flee from combat"},
	}

	for _, item := range helpItems {
		if item.key == "" {
			lines = append(lines, "")
		} else if item.desc == "" {
			lines = append(lines, labelStyle.Render(item.key))
		} else {
			lines = append(lines, fmt.Sprintf("  %-15s %s",
				hotkeyStyle.Render(item.key),
				item.desc))
		}
	}

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderTraining renders the training screen
func renderTraining(m *engine.Model, height int) string {
	player := m.Player

	var lines []string
	lines = append(lines, headerStyle.Render(" 💪 Training - Select Attribute "))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Current Attributes:"))
	lines = append(lines, fmt.Sprintf("  Strength:     Lv.%d", player.Attributes.Strength.Level))
	lines = append(lines, fmt.Sprintf("  Dexterity:    Lv.%d", player.Attributes.Dexterity.Level))
	lines = append(lines, fmt.Sprintf("  Defense:      Lv.%d", player.Attributes.Defense.Level))
	lines = append(lines, fmt.Sprintf("  Constitution: Lv.%d", player.Attributes.Constitution.Level))
	lines = append(lines, fmt.Sprintf("  Intelligence: Lv.%d", player.Attributes.Intelligence.Level))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Press letter to train:"))
	lines = append(lines, "")

	trainingOptions := []struct {
		key         rune
		name        string
		description string
	}{
		{'s', "Strength", "Train at the training dummy"},
		{'d', "Dexterity", "Run the agility course"},
		{'e', "Defense", "Practice shield drills"},
	}

	for i, opt := range trainingOptions {
		isSelected := i == m.CursorPosition

		var hotkey string
		var line string

		if isSelected {
			hotkey = selectedHotkeyStyle.Render(string(opt.key))
			line = selectedStyle.Render(fmt.Sprintf("%s %s - %s", hotkey, opt.name, opt.description))
		} else {
			hotkey = hotkeyStyle.Render(string(opt.key))
			line = fmt.Sprintf("%s %s - %s", hotkey, opt.name, dimStyle.Render(opt.description))
		}

		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [letter] Start  [↑/↓] Navigate  [Enter] Start  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderCharacterSheet renders the character sheet screen
func renderCharacterSheet(m *engine.Model, height int) string {
	player := m.Player

	var lines []string
	lines = append(lines, headerStyle.Render(fmt.Sprintf(" 👤 Character Sheet - %s ", player.Name)))
	lines = append(lines, "")

	// Combat Level
	combatLevel := player.CombatStats.GetCombatLevel()
	lines = append(lines, categoryStyle.Render("⚔️ Combat Info"))
	lines = append(lines, fmt.Sprintf("  Combat Level: %d", combatLevel))
	lines = append(lines, fmt.Sprintf("  Total Level:  %d", player.GetTotalLevel()))
	lines = append(lines, "")

	// Attributes
	lines = append(lines, categoryStyle.Render("💪 Attributes"))
	attrs := player.Attributes
	lines = append(lines, fmt.Sprintf("  Strength:     Lv.%d (XP: %s)", attrs.Strength.Level, formatNumber(attrs.Strength.XP)))
	lines = append(lines, fmt.Sprintf("  Dexterity:    Lv.%d (XP: %s)", attrs.Dexterity.Level, formatNumber(attrs.Dexterity.XP)))
	lines = append(lines, fmt.Sprintf("  Defense:      Lv.%d (XP: %s)", attrs.Defense.Level, formatNumber(attrs.Defense.XP)))
	lines = append(lines, fmt.Sprintf("  Constitution: Lv.%d (XP: %s)", attrs.Constitution.Level, formatNumber(attrs.Constitution.XP)))
	lines = append(lines, fmt.Sprintf("  Intelligence: Lv.%d (XP: %s)", attrs.Intelligence.Level, formatNumber(attrs.Intelligence.XP)))
	lines = append(lines, "")

	// Derived Stats
	lines = append(lines, categoryStyle.Render("🛡️ Combat Stats"))
	lines = append(lines, fmt.Sprintf("  HP:     %d/%d", player.CombatStats.Hitpoints, player.CombatStats.MaxHitpoints))
	lines = append(lines, fmt.Sprintf("  Attack: %d", player.CombatStats.Attack))
	lines = append(lines, fmt.Sprintf("  Ranged: %d", player.CombatStats.Ranged))
	lines = append(lines, fmt.Sprintf("  Magic:  %d", player.CombatStats.Magic))
	lines = append(lines, "")

	// Slayer Info
	lines = append(lines, categoryStyle.Render("🗡️ Slayer"))
	lines = append(lines, fmt.Sprintf("  Slayer Level: %d", player.CombatStats.SlayerLevel))
	lines = append(lines, fmt.Sprintf("  Slayer XP:    %s", formatNumber(player.CombatStats.SlayerXP)))
	lines = append(lines, fmt.Sprintf("  Slayer Points: %d", player.CombatStats.SlayerPoints))
	lines = append(lines, "")

	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [n] Change Name  [Esc/d] Dashboard  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderNameEdit renders the name editing screen
func renderNameEdit(m *engine.Model, height int) string {
	var lines []string
	lines = append(lines, headerStyle.Render(" ✏️ Change Name "))
	lines = append(lines, "")
	lines = append(lines, "Enter new name:")
	lines = append(lines, "")

	// Show current buffer with cursor
	displayName := m.NameEditBuffer
	if m.TickCount%2 == 0 {
		// Add cursor indicator
		displayName = m.NameEditBuffer[:m.NameEditCursor] + "▌" + m.NameEditBuffer[m.NameEditCursor:]
	}

	nameBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorHighlight).
		Padding(1, 2).
		Width(30).
		Render(displayName)

	lines = append(lines, "  "+nameBox)
	lines = append(lines, "")
	lines = append(lines, dimStyle.Render(fmt.Sprintf("  Length: %d/20", len(m.NameEditBuffer))))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [Enter] Confirm  [Esc] Cancel  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderSlayerTierSelection renders the monster tier selection screen
func renderSlayerTierSelection(m *engine.Model, height int) string {
	player := m.Player

	var lines []string
	lines = append(lines, headerStyle.Render(" 🗡️ Slayer - Select Tier "))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Press number or navigate to select monster tier:"))
	lines = append(lines, "")

	tierNames := []struct {
		num  int
		name string
		min  int
		max  int
		desc string
	}{
		{1, "Novice", 1, 10, "Chickens, rats, spiders, goblins"},
		{2, "Intermediate", 10, 30, "Cows, skeletons, zombies, barbarians"},
		{3, "Advanced", 30, 60, "Hill giants, moss giants, ice warriors"},
		{4, "Expert", 60, 90, "Green dragons, blue dragons, abyssal demons"},
		{5, "Master", 90, 120, "Red dragons, black dragons, legendary foes"},
	}

	for i, tier := range tierNames {
		isSelected := i == m.CursorPosition
		canAccess := player.Attributes.GetTotalAttributeLevel() >= tier.min

		tierStr := fmt.Sprintf("[%d] %s", tier.num, tier.name)
		rangeStr := fmt.Sprintf("Lv.%d-%d", tier.min, tier.max)

		var line string
		if !canAccess {
			line = lockedStyle.Render(fmt.Sprintf("  %s %s - %s", tierStr, rangeStr, tier.desc))
		} else if isSelected {
			hotkey := selectedHotkeyStyle.Render(fmt.Sprintf("%d", tier.num))
			line = selectedStyle.Render(fmt.Sprintf("%s %s - %s", hotkey, tier.name, tier.desc))
		} else {
			hotkey := hotkeyStyle.Render(fmt.Sprintf("%d", tier.num))
			line = fmt.Sprintf("%s %s %s", hotkey, tier.name, dimStyle.Render(tier.desc))
		}

		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [1-5] Select  [↑/↓] Navigate  [Enter] Confirm  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderSlayerMonsterSelection renders monster selection within a tier
func renderSlayerMonsterSelection(m *engine.Model, height int) string {
	monsters := getMonstersForTierUI(m.SelectedSlayerTier)

	var lines []string
	tierNames := []string{"Novice", "Intermediate", "Advanced", "Expert", "Master"}
	tierName := "Unknown"
	if m.SelectedSlayerTier >= 1 && m.SelectedSlayerTier <= 5 {
		tierName = tierNames[m.SelectedSlayerTier-1]
	}

	lines = append(lines, headerStyle.Render(fmt.Sprintf(" 🗡️ %s Monsters ", tierName)))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Press letter to fight monster:"))
	lines = append(lines, "")

	for i, monster := range monsters {
		isSelected := i == m.CursorPosition

		hotkeyLetter := getFirstLetter(monster.Name)

		var hotkey string
		var line string

		if isSelected {
			hotkey = selectedHotkeyStyle.Render(string(hotkeyLetter))
			line = selectedStyle.Render(fmt.Sprintf("%s %s (Lv.%d) - HP:%d", hotkey, monster.Name, monster.Level, monster.MaxHP))
			lines = append(lines, line)
			lines = append(lines, fmt.Sprintf("       %s", monster.Description))
			lines = append(lines, fmt.Sprintf("       Attack:%d Defense:%d Strength:%d", monster.Attack, monster.Defense, monster.Strength))
			lines = append(lines, fmt.Sprintf("       XP: %s combat, %s slayer", formatNumber(monster.CombatXP), formatNumber(monster.SlayerXP)))
			lines = append(lines, "")
		} else {
			hotkey = hotkeyStyle.Render(string(hotkeyLetter))
			line = fmt.Sprintf("%s %s (Lv.%d) - HP:%d", hotkey, monster.Name, monster.Level, monster.MaxHP)
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [letter] Fight  [↑/↓] Navigate  [Enter] Fight  [Esc/q] Back  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// renderCombat renders the active combat screen
func renderCombat(m *engine.Model, height int) string {
	if m.CurrentCombatEncounter == nil {
		return boxStyle.Render("Error: No active combat")
	}

	encounter := m.CurrentCombatEncounter
	monster := encounter.Monster

	var lines []string
	lines = append(lines, headerStyle.Render(fmt.Sprintf(" ⚔️ Combat vs %s ", monster.Name)))
	lines = append(lines, "")

	// Monster info
	lines = append(lines, categoryStyle.Render("👹 Enemy"))
	lines = append(lines, fmt.Sprintf("  %s (Lv.%d)", monster.Name, monster.Level))
	lines = append(lines, fmt.Sprintf("  HP: %d/%d", monster.Hitpoints, monster.MaxHP))

	// Monster HP bar
	monsterHPPercent := float64(monster.Hitpoints) / float64(monster.MaxHP)
	monsterHPBar := renderProgressBar(monsterHPPercent, 25)
	if monsterHPPercent > 0.5 {
		monsterHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(monsterHPBar)
	} else if monsterHPPercent > 0.25 {
		monsterHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaa00")).Render(monsterHPBar)
	} else {
		monsterHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(monsterHPBar)
	}
	lines = append(lines, fmt.Sprintf("  %s", monsterHPBar))
	lines = append(lines, "")

	// Player info
	lines = append(lines, categoryStyle.Render("🛡️ You"))
	lines = append(lines, fmt.Sprintf("  HP: %d/%d", m.Player.CombatStats.Hitpoints, m.Player.CombatStats.MaxHitpoints))

	// Player HP bar
	playerHPPercent := float64(m.Player.CombatStats.Hitpoints) / float64(m.Player.CombatStats.MaxHitpoints)
	playerHPBar := renderProgressBar(playerHPPercent, 25)
	if playerHPPercent > 0.5 {
		playerHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(playerHPBar)
	} else if playerHPPercent > 0.25 {
		playerHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaa00")).Render(playerHPBar)
	} else {
		playerHPBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(playerHPBar)
	}
	lines = append(lines, fmt.Sprintf("  %s", playerHPBar))

	// ATB Bars
	lines = append(lines, "")
	lines = append(lines, categoryStyle.Render("⚡ ATB Gauge"))

	playerATBBar := renderProgressBar(encounter.PlayerATB/100.0, 25)
	playerATBBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff")).Render(playerATBBar)
	lines = append(lines, fmt.Sprintf("  You:    %s %.0f%%", playerATBBar, encounter.PlayerATB))

	monsterATBBar := renderProgressBar(encounter.MonsterATB/100.0, 25)
	monsterATBBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff00ff")).Render(monsterATBBar)
	lines = append(lines, fmt.Sprintf("  Enemy:  %s %.0f%%", monsterATBBar, encounter.MonsterATB))
	lines = append(lines, "")

	// Combat log
	if encounter.LastActionResult != "" {
		lines = append(lines, categoryStyle.Render("📜 Last Action"))
		lines = append(lines, "  "+encounter.LastActionResult)
		lines = append(lines, "")
	}

	// Combat stats
	lines = append(lines, dimStyle.Render(fmt.Sprintf("Combat Ticks: %d | Damage Dealt: %d | Damage Taken: %d",
		encounter.CombatTicks, encounter.DamageDealt, encounter.DamageTaken)))

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Background(lipgloss.Color("#333333")).
		Foreground(colorInfo).
		Render("  [Esc/q] Flee Combat  "))

	return boxStyle.
		Height(height).
		Width(m.Width - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// Helper function to get monsters for tier
func getMonstersForTierUI(tier int) []*models.Monster {
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

// Helper functions

func getSkillIcon(skill models.SkillType) string {
	switch skill {
	case models.SkillWoodcutting:
		return "🪓"
	case models.SkillMining:
		return "⛏️"
	case models.SkillSmithing:
		return "🔥"
	case models.SkillRecycling:
		return "♻️"
	case models.SkillCrafting:
		return "🛠️"
	default:
		return "⭐"
	}
}

func getActionVerb(skill models.SkillType) string {
	switch skill {
	case models.SkillWoodcutting:
		return "chop"
	case models.SkillMining:
		return "mine"
	case models.SkillSmithing:
		return "craft"
	case models.SkillRecycling:
		return "recycle"
	default:
		return "do"
	}
}

func getFirstLetter(s string) rune {
	if len(s) == 0 {
		return '?'
	}
	return rune(strings.ToLower(s)[0])
}

func renderProgressBar(progress float64, width int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}
	empty := width - filled

	fullChar := "█"
	emptyChar := "░"

	bar := strings.Repeat(fullChar, filled) + strings.Repeat(emptyChar, empty)

	if progress >= 1 {
		return lipgloss.NewStyle().Foreground(colorPrimary).Render(bar)
	}

	return lipgloss.NewStyle().Foreground(colorSecondary).Render(bar)
}

func overlayMessage(content string, message string, width int, height int) string {
	msg := messageStyle.Width(width - 10).Render(message)

	msgHeight := lipgloss.Height(msg)
	msgWidth := lipgloss.Width(msg)

	topPadding := (height - msgHeight) / 2
	leftPadding := (width - msgWidth) / 2

	lines := strings.Split(content, "\n")
	if len(lines) > topPadding {
		var result []string
		result = append(result, lines[:topPadding]...)

		msgLines := strings.Split(msg, "\n")
		for _, line := range msgLines {
			result = append(result, strings.Repeat(" ", leftPadding)+line)
		}

		result = append(result, lines[topPadding+msgHeight:]...)
		return strings.Join(result, "\n")
	}

	return content + "\n\n" + msg
}

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
