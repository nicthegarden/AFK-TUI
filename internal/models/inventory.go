package models

import (
	"fmt"
	"sort"
)

// Inventory holds items with limited slots
type Inventory struct {
	Items    []*Item `json:"items"`
	MaxSlots int     `json:"max_slots"`
}

// NewInventory creates a new inventory
func NewInventory(maxSlots int) *Inventory {
	return &Inventory{
		Items:    make([]*Item, 0),
		MaxSlots: maxSlots,
	}
}

// AddItem adds an item to inventory, stacking if possible
func (inv *Inventory) AddItem(item *Item) bool {
	// Try to stack with existing items
	if item.Type == ItemTypeResource || item.Type == ItemTypeMaterial || item.Type == ItemTypeBar {
		for _, existing := range inv.Items {
			if existing.CanStackWith(item) {
				existing.Quantity += item.Quantity
				return true
			}
		}
	}

	// Check if we have space
	if len(inv.Items) >= inv.MaxSlots {
		return false
	}

	// Add as new item
	inv.Items = append(inv.Items, item)
	return true
}

// AddItems adds multiple items
func (inv *Inventory) AddItems(items []*Item) []string {
	var failed []string
	for _, item := range items {
		if !inv.AddItem(item) {
			failed = append(failed, item.Name)
		}
	}
	return failed
}

// RemoveItem removes a quantity of an item
func (inv *Inventory) RemoveItem(itemID string, quantity int) bool {
	for i, item := range inv.Items {
		if item.ID == itemID {
			if item.Quantity >= quantity {
				item.Quantity -= quantity
				if item.Quantity == 0 {
					inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
				}
				return true
			}
			return false
		}
	}
	return false
}

// GetItem finds an item by ID
func (inv *Inventory) GetItem(itemID string) *Item {
	for _, item := range inv.Items {
		if item.ID == itemID {
			return item
		}
	}
	return nil
}

// HasItem checks if inventory has enough of an item
func (inv *Inventory) HasItem(itemID string, quantity int) bool {
	item := inv.GetItem(itemID)
	if item == nil {
		return false
	}
	return item.Quantity >= quantity
}

// GetTotalValue calculates total gold value
func (inv *Inventory) GetTotalValue() int64 {
	total := int64(0)
	for _, item := range inv.Items {
		total += item.Value * int64(item.Quantity)
	}
	return total
}

// SortByType sorts items by type then name
func (inv *Inventory) SortByType() {
	sort.Slice(inv.Items, func(i, j int) bool {
		if inv.Items[i].Type != inv.Items[j].Type {
			return inv.Items[i].Type < inv.Items[j].Type
		}
		return inv.Items[i].Name < inv.Items[j].Name
	})
}

// Count returns number of item stacks
func (inv *Inventory) Count() int {
	return len(inv.Items)
}

// IsFull checks if inventory is full
func (inv *Inventory) IsFull() bool {
	return len(inv.Items) >= inv.MaxSlots
}

// String returns inventory summary
func (inv *Inventory) String() string {
	return fmt.Sprintf("Inventory: %d/%d slots", len(inv.Items), inv.MaxSlots)
}

// Bank is unlimited storage
type Bank struct {
	Items []*Item `json:"items"`
}

// NewBank creates a new bank
func NewBank(initialCapacity int) *Bank {
	return &Bank{
		Items: make([]*Item, 0),
	}
}

// Deposit moves items from inventory to bank
func (b *Bank) Deposit(inv *Inventory, itemID string, quantity int) bool {
	if !inv.HasItem(itemID, quantity) {
		return false
	}

	item := inv.GetItem(itemID)

	// Try to stack in bank
	for _, bankItem := range b.Items {
		if bankItem.ID == itemID {
			bankItem.Quantity += quantity
			inv.RemoveItem(itemID, quantity)
			return true
		}
	}

	// Create new bank item
	newItem := item.Clone()
	newItem.Quantity = quantity
	b.Items = append(b.Items, newItem)
	inv.RemoveItem(itemID, quantity)
	return true
}

// Withdraw moves items from bank to inventory
func (b *Bank) Withdraw(inv *Inventory, itemID string, quantity int) bool {
	for i, item := range b.Items {
		if item.ID == itemID && item.Quantity >= quantity {
			withdrawn := item.Clone()
			withdrawn.Quantity = quantity

			if inv.AddItem(withdrawn) {
				item.Quantity -= quantity
				if item.Quantity == 0 {
					b.Items = append(b.Items[:i], b.Items[i+1:]...)
				}
				return true
			}
			return false
		}
	}
	return false
}

// GetItem finds item in bank
func (b *Bank) GetItem(itemID string) *Item {
	for _, item := range b.Items {
		if item.ID == itemID {
			return item
		}
	}
	return nil
}

// HasItem checks if bank has item
func (b *Bank) HasItem(itemID string, quantity int) bool {
	item := b.GetItem(itemID)
	if item == nil {
		return false
	}
	return item.Quantity >= quantity
}
