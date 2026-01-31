package engine

import (
	"afk-tui/internal/models"
	"fmt"
	"strconv"
	"strings"
)

// InventoryState tracks inventory view state
type InventoryState struct {
	IsSellMode       bool
	SelectedItem     int
	QuantityToSell   int
	ShowConfirmation bool
	ItemName         string
	ItemID           string
	MaxQuantity      int
	GoldValue        int64
}

// HandleInventoryInput processes inventory key presses
func HandleInventoryInput(msg string, state *InventoryState, player *models.Player) (string, bool, int64) {
	// If in sell mode with confirmation showing
	if state.IsSellMode && state.ShowConfirmation {
		switch msg {
		case "y", "Y":
			// Confirm sell
			if state.ItemID != "" {
				// Remove items
				player.Inventory.RemoveItem(state.ItemID, state.QuantityToSell)
				// Add gold
				player.Gold += state.GoldValue
				// Log it
				player.ActivityLog.AddSellLog(state.ItemName, state.QuantityToSell, state.GoldValue)

				message := fmt.Sprintf("Sold %dx %s for %d gold", state.QuantityToSell, state.ItemName, state.GoldValue)
				ResetInventoryState(state)
				return message, true, state.GoldValue
			}
		case "n", "N", "esc":
			// Cancel
			ResetInventoryState(state)
			return "Sell cancelled", false, 0
		}
		return "", false, 0
	}

	// If in sell mode selecting item
	if state.IsSellMode && !state.ShowConfirmation {
		switch msg {
		case "esc", "q":
			ResetInventoryState(state)
			return "Sell mode cancelled", false, 0

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Building quantity number
			num, _ := strconv.Atoi(msg)
			if state.SelectedItem > 0 && state.SelectedItem <= len(player.Inventory.Items) {
				// Adding to quantity
				state.QuantityToSell = state.QuantityToSell*10 + num
				if state.QuantityToSell > state.MaxQuantity {
					state.QuantityToSell = state.MaxQuantity
				}
			} else {
				// Selecting item by number
				state.SelectedItem = num
				if state.SelectedItem > 0 && state.SelectedItem <= len(player.Inventory.Items) {
					item := player.Inventory.Items[state.SelectedItem-1]
					state.ItemName = item.Name
					state.ItemID = item.ID
					state.MaxQuantity = item.Quantity
					state.QuantityToSell = item.Quantity // Default to max
					state.GoldValue = item.Value * int64(item.Quantity)
				}
			}
			return "", false, 0

		case "enter":
			if state.SelectedItem > 0 && state.SelectedItem <= len(player.Inventory.Items) {
				state.ShowConfirmation = true
			}
			return "", false, 0

		case "max":
			// Set to max quantity
			if state.SelectedItem > 0 && state.SelectedItem <= len(player.Inventory.Items) {
				state.QuantityToSell = state.MaxQuantity
				item := player.Inventory.Items[state.SelectedItem-1]
				state.GoldValue = item.Value * int64(state.QuantityToSell)
			}
			return "", false, 0

		case "backspace":
			// Clear selection
			state.SelectedItem = 0
			state.QuantityToSell = 0
			state.ItemName = ""
			state.ItemID = ""
			return "", false, 0
		}
	}

	// Normal inventory mode
	switch msg {
	case "v":
		// Enter sell mode (v = vend/sell)
		state.IsSellMode = true
		state.SelectedItem = 0
		state.QuantityToSell = 0
		return "Sell mode: Enter item number, then quantity (or 'max')", false, 0

	case "esc", "q":
		return "", false, 0 // Return to dashboard
	}

	return "", false, 0
}

// ResetInventoryState resets the inventory state
func ResetInventoryState(state *InventoryState) {
	state.IsSellMode = false
	state.SelectedItem = 0
	state.QuantityToSell = 0
	state.ShowConfirmation = false
	state.ItemName = ""
	state.ItemID = ""
	state.MaxQuantity = 0
	state.GoldValue = 0
}

// GetSellConfirmationText returns the confirmation dialog text
func GetSellConfirmationText(state *InventoryState) string {
	if !state.ShowConfirmation {
		return ""
	}

	return fmt.Sprintf("\n╔════════════════════════════════════════╗\n"+
		"║  CONFIRM SELL                         ║\n"+
		"║                                       ║\n"+
		"║  Sell: %dx %-25s ║\n"+
		"║  For:  %d gold                        ║\n"+
		"║                                       ║\n"+
		"║  [Y] Yes    [N] No                    ║\n"+
		"╚════════════════════════════════════════╝",
		state.QuantityToSell,
		state.ItemName,
		state.GoldValue)
}

// ParseSellInput parses sell mode input
func ParseSellInput(input string) (itemNum int, quantity int, isMax bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "max" {
		return 0, 0, true
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return 0, 0, false
	}

	// First number is item number
	itemNum, _ = strconv.Atoi(parts[0])

	// Second number (if exists) is quantity
	if len(parts) > 1 {
		quantity, _ = strconv.Atoi(parts[1])
	}

	return itemNum, quantity, false
}
