package models

import "fmt"

// Equipment holds all equipped items
type Equipment struct {
	Head    *Item `json:"head,omitempty"`
	Body    *Item `json:"body,omitempty"`
	Legs    *Item `json:"legs,omitempty"`
	Feet    *Item `json:"feet,omitempty"`
	Hands   *Item `json:"hands,omitempty"`
	Weapon  *Item `json:"weapon,omitempty"`
	Offhand *Item `json:"offhand,omitempty"`
	Cape    *Item `json:"cape,omitempty"`
	Ring    *Item `json:"ring,omitempty"`
	Amulet  *Item `json:"amulet,omitempty"`
	Ammo    *Item `json:"ammo,omitempty"`
}

// NewEquipment creates empty equipment
func NewEquipment() *Equipment {
	return &Equipment{}
}

// GetSlot gets item in a specific slot
func (e *Equipment) GetSlot(slot EquipmentSlot) *Item {
	switch slot {
	case SlotHead:
		return e.Head
	case SlotBody:
		return e.Body
	case SlotLegs:
		return e.Legs
	case SlotFeet:
		return e.Feet
	case SlotHands:
		return e.Hands
	case SlotWeapon:
		return e.Weapon
	case SlotOffhand:
		return e.Offhand
	case SlotCape:
		return e.Cape
	case SlotRing:
		return e.Ring
	case SlotAmulet:
		return e.Amulet
	case SlotAmmo:
		return e.Ammo
	}
	return nil
}

// SetSlot sets item in a specific slot
func (e *Equipment) SetSlot(slot EquipmentSlot, item *Item) *Item {
	var old *Item
	switch slot {
	case SlotHead:
		old = e.Head
		e.Head = item
	case SlotBody:
		old = e.Body
		e.Body = item
	case SlotLegs:
		old = e.Legs
		e.Legs = item
	case SlotFeet:
		old = e.Feet
		e.Feet = item
	case SlotHands:
		old = e.Hands
		e.Hands = item
	case SlotWeapon:
		old = e.Weapon
		e.Weapon = item
	case SlotOffhand:
		old = e.Offhand
		e.Offhand = item
	case SlotCape:
		old = e.Cape
		e.Cape = item
	case SlotRing:
		old = e.Ring
		e.Ring = item
	case SlotAmulet:
		old = e.Amulet
		e.Amulet = item
	case SlotAmmo:
		old = e.Ammo
		e.Ammo = item
	}
	return old
}

// Equip equips an item from inventory
func (e *Equipment) Equip(inv *Inventory, item *Item) (*Item, error) {
	if !item.IsEquipable() {
		return nil, fmt.Errorf("item cannot be equipped")
	}

	if !inv.HasItem(item.ID, 1) {
		return nil, fmt.Errorf("item not in inventory")
	}

	// Remove from inventory
	inv.RemoveItem(item.ID, 1)

	// Unequip current item if any
	oldItem := e.SetSlot(item.Slot, item)

	// Put old item back in inventory if exists
	if oldItem != nil {
		inv.AddItem(oldItem)
	}

	return oldItem, nil
}

// Unequip removes item from slot and returns to inventory
func (e *Equipment) Unequip(slot EquipmentSlot, inv *Inventory) (*Item, error) {
	item := e.GetSlot(slot)
	if item == nil {
		return nil, fmt.Errorf("nothing equipped in that slot")
	}

	if !inv.AddItem(item) {
		return nil, fmt.Errorf("inventory full")
	}

	e.SetSlot(slot, nil)
	return item, nil
}

// GetTotalStats calculates total equipment stats
func (e *Equipment) GetTotalStats() map[string]int {
	stats := map[string]int{
		"attack":     0,
		"strength":   0,
		"defence":    0,
		"tool_power": 0,
	}

	slots := []*Item{e.Head, e.Body, e.Legs, e.Feet, e.Hands, e.Weapon, e.Offhand, e.Cape, e.Ring, e.Amulet}

	for _, item := range slots {
		if item != nil && item.Stats != nil {
			for stat, value := range item.Stats {
				if current, ok := stats[stat]; ok {
					stats[stat] = current + value
				}
			}
		}
		if item != nil {
			stats["tool_power"] += item.ToolPower
		}
	}

	return stats
}

// GetToolPower returns total tool power bonus
func (e *Equipment) GetToolPower() int {
	return e.GetTotalStats()["tool_power"]
}

// String returns equipment summary
func (e *Equipment) String() string {
	stats := e.GetTotalStats()
	return fmt.Sprintf("ATK:%d STR:%d DEF:%d PWR:%d",
		stats["attack"], stats["strength"], stats["defence"], stats["tool_power"])
}
