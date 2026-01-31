package models

// ItemType categorizes items
type ItemType string

const (
	ItemTypeResource   ItemType = "resource"
	ItemTypeTool       ItemType = "tool"
	ItemTypeWeapon     ItemType = "weapon"
	ItemTypeArmor      ItemType = "armor"
	ItemTypeConsumable ItemType = "consumable"
	ItemTypeMaterial   ItemType = "material"
	ItemTypeBar        ItemType = "bar"
)

// Item represents any item in the game
type Item struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         ItemType               `json:"type"`
	Quantity     int                    `json:"quantity"`
	Value        int64                  `json:"value"` // Gold value
	Requirements map[SkillType]int      `json:"requirements,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`

	// For equipment
	Slot  EquipmentSlot  `json:"slot,omitempty"`
	Stats map[string]int `json:"stats,omitempty"` // attack, defence, etc.

	// For tools
	ToolPower int `json:"tool_power,omitempty"` // Gathering speed bonus

	// For recycled materials
	RecycleValue map[string]int `json:"recycle_value,omitempty"` // What you get when recycling
}

// NewItem creates a new item with defaults from ItemDatabase
func NewItem(id string, name string, quantity int) *Item {
	// Check if exists in database
	if template, ok := ItemDatabase[id]; ok {
		item := &Item{
			ID:           template.ID,
			Name:         template.Name,
			Description:  template.Description,
			Type:         template.Type,
			Quantity:     quantity,
			Value:        template.Value,
			Requirements: template.Requirements,
			Slot:         template.Slot,
			Stats:        template.Stats,
			ToolPower:    template.ToolPower,
			RecycleValue: template.RecycleValue,
			Metadata:     make(map[string]interface{}),
		}
		return item
	}

	// Fallback for unknown items
	return &Item{
		ID:       id,
		Name:     name,
		Quantity: quantity,
		Type:     ItemTypeResource,
		Metadata: make(map[string]interface{}),
	}
}

// Clone creates a copy of the item
func (i *Item) Clone() *Item {
	return NewItem(i.ID, i.Name, i.Quantity)
}

// CanStackWith checks if two items can stack
func (i *Item) CanStackWith(other *Item) bool {
	return i.ID == other.ID && i.Type != ItemTypeTool && i.Type != ItemTypeWeapon && i.Type != ItemTypeArmor
}

// IsEquipable returns true if item can be equipped
func (i *Item) IsEquipable() bool {
	return i.Type == ItemTypeTool || i.Type == ItemTypeWeapon || i.Type == ItemTypeArmor
}

// IsRecyclable returns true if item can be recycled
func (i *Item) IsRecyclable() bool {
	return i.RecycleValue != nil && len(i.RecycleValue) > 0
}

// GetRecycleYield calculates what you get from recycling
func (i *Item) GetRecycleYield() map[string]int {
	if i.RecycleValue == nil {
		return nil
	}

	yield := make(map[string]int)
	for material, amount := range i.RecycleValue {
		yield[material] = amount * i.Quantity
	}
	return yield
}

// EquipmentSlot represents equipment slots
type EquipmentSlot string

const (
	SlotHead    EquipmentSlot = "head"
	SlotBody    EquipmentSlot = "body"
	SlotLegs    EquipmentSlot = "legs"
	SlotFeet    EquipmentSlot = "feet"
	SlotHands   EquipmentSlot = "hands"
	SlotWeapon  EquipmentSlot = "weapon"
	SlotOffhand EquipmentSlot = "offhand"
	SlotCape    EquipmentSlot = "cape"
	SlotRing    EquipmentSlot = "ring"
	SlotAmulet  EquipmentSlot = "amulet"
	SlotAmmo    EquipmentSlot = "ammo"
)

// ItemDatabase contains all item definitions
var ItemDatabase = map[string]*Item{
	// Resources - Wood
	"logs": {
		ID: "logs", Name: "Logs", Type: ItemTypeResource, Value: 5,
		Description:  "Basic wood from any tree",
		RecycleValue: map[string]int{"wood_fragments": 1},
	},
	"oak_logs": {
		ID: "oak_logs", Name: "Oak Logs", Type: ItemTypeResource, Value: 15,
		Description:  "Sturdy oak wood",
		RecycleValue: map[string]int{"wood_fragments": 2},
	},
	"willow_logs": {
		ID: "willow_logs", Name: "Willow Logs", Type: ItemTypeResource, Value: 30,
		Description:  "Flexible willow wood",
		RecycleValue: map[string]int{"wood_fragments": 3},
	},
	"maple_logs": {
		ID: "maple_logs", Name: "Maple Logs", Type: ItemTypeResource, Value: 60,
		Description:  "Quality maple wood",
		RecycleValue: map[string]int{"wood_fragments": 4},
	},
	"yew_logs": {
		ID: "yew_logs", Name: "Yew Logs", Type: ItemTypeResource, Value: 120,
		Description:  "Rare yew wood",
		RecycleValue: map[string]int{"wood_fragments": 5},
	},
	"magic_logs": {
		ID: "magic_logs", Name: "Magic Logs", Type: ItemTypeResource, Value: 250,
		Description:  "Enchanted wood with magical properties",
		RecycleValue: map[string]int{"wood_fragments": 8, "magic_essence": 1},
	},

	// Resources - Ores
	"copper_ore": {
		ID: "copper_ore", Name: "Copper Ore", Type: ItemTypeResource, Value: 8,
		Description:  "Common copper ore",
		RecycleValue: map[string]int{"metal_fragments": 1},
	},
	"tin_ore": {
		ID: "tin_ore", Name: "Tin Ore", Type: ItemTypeResource, Value: 8,
		Description:  "Common tin ore",
		RecycleValue: map[string]int{"metal_fragments": 1},
	},
	"iron_ore": {
		ID: "iron_ore", Name: "Iron Ore", Type: ItemTypeResource, Value: 20,
		Description:  "Strong iron ore",
		RecycleValue: map[string]int{"metal_fragments": 2},
	},
	"coal": {
		ID: "coal", Name: "Coal", Type: ItemTypeResource, Value: 25,
		Description:  "Used for smelting",
		RecycleValue: map[string]int{"coal_fragments": 1},
	},
	"silver_ore": {
		ID: "silver_ore", Name: "Silver Ore", Type: ItemTypeResource, Value: 50,
		Description:  "Precious silver ore",
		RecycleValue: map[string]int{"metal_fragments": 3, "silver_fragments": 1},
	},
	"gold_ore": {
		ID: "gold_ore", Name: "Gold Ore", Type: ItemTypeResource, Value: 100,
		Description:  "Valuable gold ore",
		RecycleValue: map[string]int{"metal_fragments": 4, "gold_fragments": 1},
	},
	"mithril_ore": {
		ID: "mithril_ore", Name: "Mithril Ore", Type: ItemTypeResource, Value: 200,
		Description:  "Lightweight mithril ore",
		RecycleValue: map[string]int{"metal_fragments": 5, "mithril_fragments": 1},
	},
	"adamantite_ore": {
		ID: "adamantite_ore", Name: "Adamantite Ore", Type: ItemTypeResource, Value: 400,
		Description:  "Heavy adamantite ore",
		RecycleValue: map[string]int{"metal_fragments": 6, "adamantite_fragments": 1},
	},
	"runite_ore": {
		ID: "runite_ore", Name: "Runite Ore", Type: ItemTypeResource, Value: 1000,
		Description:  "Mystical runite ore",
		RecycleValue: map[string]int{"metal_fragments": 8, "rune_fragments": 1},
	},

	// Bars
	"bronze_bar": {
		ID: "bronze_bar", Name: "Bronze Bar", Type: ItemTypeBar, Value: 20,
		Description:  "Copper and tin alloy",
		RecycleValue: map[string]int{"metal_fragments": 2, "copper_fragments": 1},
	},
	"iron_bar": {
		ID: "iron_bar", Name: "Iron Bar", Type: ItemTypeBar, Value: 50,
		Description:  "Pure iron bar",
		RecycleValue: map[string]int{"metal_fragments": 3},
	},
	"steel_bar": {
		ID: "steel_bar", Name: "Steel Bar", Type: ItemTypeBar, Value: 120,
		Description:  "Strong steel alloy",
		RecycleValue: map[string]int{"metal_fragments": 4, "steel_fragments": 1},
	},
	"silver_bar": {
		ID: "silver_bar", Name: "Silver Bar", Type: ItemTypeBar, Value: 150,
		Description:  "Pure silver bar",
		RecycleValue: map[string]int{"metal_fragments": 3, "silver_fragments": 2},
	},
	"gold_bar": {
		ID: "gold_bar", Name: "Gold Bar", Type: ItemTypeBar, Value: 250,
		Description:  "Pure gold bar",
		RecycleValue: map[string]int{"metal_fragments": 4, "gold_fragments": 2},
	},
	"mithril_bar": {
		ID: "mithril_bar", Name: "Mithril Bar", Type: ItemTypeBar, Value: 500,
		Description:  "Lightweight mithril bar",
		RecycleValue: map[string]int{"metal_fragments": 5, "mithril_fragments": 2},
	},
	"adamantite_bar": {
		ID: "adamantite_bar", Name: "Adamantite Bar", Type: ItemTypeBar, Value: 1000,
		Description:  "Heavy adamantite bar",
		RecycleValue: map[string]int{"metal_fragments": 6, "adamantite_fragments": 2},
	},
	"runite_bar": {
		ID: "runite_bar", Name: "Runite Bar", Type: ItemTypeBar, Value: 2500,
		Description:  "Powerful runite bar",
		RecycleValue: map[string]int{"metal_fragments": 8, "rune_fragments": 2},
	},

	// Recycled Materials
	"wood_fragments": {
		ID: "wood_fragments", Name: "Wood Fragments", Type: ItemTypeMaterial, Value: 1,
		Description: "Tiny pieces of wood from recycling",
	},
	"metal_fragments": {
		ID: "metal_fragments", Name: "Metal Fragments", Type: ItemTypeMaterial, Value: 2,
		Description: "Scrap metal from recycling",
	},
	"copper_fragments": {
		ID: "copper_fragments", Name: "Copper Fragments", Type: ItemTypeMaterial, Value: 3,
		Description: "Pure copper fragments",
	},
	"silver_fragments": {
		ID: "silver_fragments", Name: "Silver Fragments", Type: ItemTypeMaterial, Value: 10,
		Description: "Pure silver fragments",
	},
	"gold_fragments": {
		ID: "gold_fragments", Name: "Gold Fragments", Type: ItemTypeMaterial, Value: 25,
		Description: "Pure gold fragments",
	},
	"mithril_fragments": {
		ID: "mithril_fragments", Name: "Mithril Fragments", Type: ItemTypeMaterial, Value: 50,
		Description: "Lightweight mithril fragments",
	},
	"adamantite_fragments": {
		ID: "adamantite_fragments", Name: "Adamantite Fragments", Type: ItemTypeMaterial, Value: 100,
		Description: "Heavy adamantite fragments",
	},
	"rune_fragments": {
		ID: "rune_fragments", Name: "Rune Fragments", Type: ItemTypeMaterial, Value: 250,
		Description: "Mystical runite fragments",
	},
	"steel_fragments": {
		ID: "steel_fragments", Name: "Steel Fragments", Type: ItemTypeMaterial, Value: 5,
		Description: "High-quality steel fragments",
	},
	"coal_fragments": {
		ID: "coal_fragments", Name: "Coal Fragments", Type: ItemTypeMaterial, Value: 3,
		Description: "Compressed carbon fragments",
	},
	"magic_essence": {
		ID: "magic_essence", Name: "Magic Essence", Type: ItemTypeMaterial, Value: 100,
		Description: "Concentrated magical energy",
	},

	// Tools
	"bronze_axe": {
		ID: "bronze_axe", Name: "Bronze Axe", Type: ItemTypeTool, Value: 50,
		Description: "Basic woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 1,
		Requirements: map[SkillType]int{SkillWoodcutting: 1},
		RecycleValue: map[string]int{"metal_fragments": 3, "copper_fragments": 2},
	},
	"iron_axe": {
		ID: "iron_axe", Name: "Iron Axe", Type: ItemTypeTool, Value: 150,
		Description: "Better woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 3,
		Requirements: map[SkillType]int{SkillWoodcutting: 10},
		RecycleValue: map[string]int{"metal_fragments": 5},
	},
	"steel_axe": {
		ID: "steel_axe", Name: "Steel Axe", Type: ItemTypeTool, Value: 400,
		Description: "Efficient woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 5,
		Requirements: map[SkillType]int{SkillWoodcutting: 20},
		RecycleValue: map[string]int{"metal_fragments": 6, "steel_fragments": 2},
	},
	"bronze_pickaxe": {
		ID: "bronze_pickaxe", Name: "Bronze Pickaxe", Type: ItemTypeTool, Value: 50,
		Description: "Basic mining tool",
		Slot:        SlotWeapon, ToolPower: 1,
		Requirements: map[SkillType]int{SkillMining: 1},
		RecycleValue: map[string]int{"metal_fragments": 3, "copper_fragments": 2},
	},
	"iron_pickaxe": {
		ID: "iron_pickaxe", Name: "Iron Pickaxe", Type: ItemTypeTool, Value: 150,
		Description: "Better mining tool",
		Slot:        SlotWeapon, ToolPower: 3,
		Requirements: map[SkillType]int{SkillMining: 10},
		RecycleValue: map[string]int{"metal_fragments": 5},
	},
	"steel_pickaxe": {
		ID: "steel_pickaxe", Name: "Steel Pickaxe", Type: ItemTypeTool, Value: 400,
		Description: "Efficient mining tool",
		Slot:        SlotWeapon, ToolPower: 5,
		Requirements: map[SkillType]int{SkillMining: 20},
		RecycleValue: map[string]int{"metal_fragments": 6, "steel_fragments": 2},
	},

	// Weapons
	"bronze_sword": {
		ID: "bronze_sword", Name: "Bronze Sword", Type: ItemTypeWeapon, Value: 100,
		Description:  "Basic melee weapon",
		Slot:         SlotWeapon,
		Stats:        map[string]int{"attack": 5, "strength": 3},
		Requirements: map[SkillType]int{SkillCombat: 1},
		RecycleValue: map[string]int{"metal_fragments": 4, "copper_fragments": 2},
	},
	"iron_sword": {
		ID: "iron_sword", Name: "Iron Sword", Type: ItemTypeWeapon, Value: 250,
		Description:  "Stronger melee weapon",
		Slot:         SlotWeapon,
		Stats:        map[string]int{"attack": 10, "strength": 7},
		Requirements: map[SkillType]int{SkillCombat: 10},
		RecycleValue: map[string]int{"metal_fragments": 6},
	},
	"steel_sword": {
		ID: "steel_sword", Name: "Steel Sword", Type: ItemTypeWeapon, Value: 600,
		Description:  "Reliable melee weapon",
		Slot:         SlotWeapon,
		Stats:        map[string]int{"attack": 15, "strength": 12},
		Requirements: map[SkillType]int{SkillCombat: 20},
		RecycleValue: map[string]int{"metal_fragments": 8, "steel_fragments": 3},
	},

	// Armor
	"bronze_helmet": {
		ID: "bronze_helmet", Name: "Bronze Helmet", Type: ItemTypeArmor, Value: 80,
		Description:  "Basic head protection",
		Slot:         SlotHead,
		Stats:        map[string]int{"defence": 3},
		Requirements: map[SkillType]int{SkillCombat: 1},
		RecycleValue: map[string]int{"metal_fragments": 3, "copper_fragments": 2},
	},
	"bronze_body": {
		ID: "bronze_body", Name: "Bronze Body", Type: ItemTypeArmor, Value: 150,
		Description:  "Basic body protection",
		Slot:         SlotBody,
		Stats:        map[string]int{"defence": 5},
		Requirements: map[SkillType]int{SkillCombat: 1},
		RecycleValue: map[string]int{"metal_fragments": 5, "copper_fragments": 3},
	},
	"bronze_legs": {
		ID: "bronze_legs", Name: "Bronze Legs", Type: ItemTypeArmor, Value: 120,
		Description:  "Basic leg protection",
		Slot:         SlotLegs,
		Stats:        map[string]int{"defence": 4},
		Requirements: map[SkillType]int{SkillCombat: 1},
		RecycleValue: map[string]int{"metal_fragments": 4, "copper_fragments": 2},
	},

	// Additional Ores & Minerals
	"lead_ore": {
		ID: "lead_ore", Name: "Lead Ore", Type: ItemTypeResource, Value: 15,
		Description:  "Soft heavy metal ore",
		RecycleValue: map[string]int{"metal_fragments": 1, "lead_fragments": 1},
	},
	"zinc_ore": {
		ID: "zinc_ore", Name: "Zinc Ore", Type: ItemTypeResource, Value: 18,
		Description:  "Used for brass making",
		RecycleValue: map[string]int{"metal_fragments": 1, "zinc_fragments": 1},
	},
	"nickel_ore": {
		ID: "nickel_ore", Name: "Nickel Ore", Type: ItemTypeResource, Value: 35,
		Description:  "Corrosion-resistant metal",
		RecycleValue: map[string]int{"metal_fragments": 2, "nickel_fragments": 1},
	},
	"platinum_ore": {
		ID: "platinum_ore", Name: "Platinum Ore", Type: ItemTypeResource, Value: 500,
		Description:  "Extremely rare precious metal",
		RecycleValue: map[string]int{"metal_fragments": 5, "platinum_fragments": 1},
	},
	"obsidian_ore": {
		ID: "obsidian_ore", Name: "Obsidian Ore", Type: ItemTypeResource, Value: 800,
		Description:  "Volcanic glass ore",
		RecycleValue: map[string]int{"obsidian_fragments": 1},
	},

	// Gems - Uncut
	"uncut_sapphire": {
		ID: "uncut_sapphire", Name: "Uncut Sapphire", Type: ItemTypeResource, Value: 150,
		Description:  "A rough blue gemstone",
		RecycleValue: map[string]int{"gem_fragments": 2},
	},
	"uncut_emerald": {
		ID: "uncut_emerald", Name: "Uncut Emerald", Type: ItemTypeResource, Value: 300,
		Description:  "A rough green gemstone",
		RecycleValue: map[string]int{"gem_fragments": 3},
	},
	"uncut_ruby": {
		ID: "uncut_ruby", Name: "Uncut Ruby", Type: ItemTypeResource, Value: 600,
		Description:  "A rough red gemstone",
		RecycleValue: map[string]int{"gem_fragments": 4},
	},
	"uncut_diamond": {
		ID: "uncut_diamond", Name: "Uncut Diamond", Type: ItemTypeResource, Value: 1500,
		Description:  "A rough clear gemstone",
		RecycleValue: map[string]int{"gem_fragments": 5},
	},
	"uncut_dragonstone": {
		ID: "uncut_dragonstone", Name: "Uncut Dragonstone", Type: ItemTypeResource, Value: 5000,
		Description:  "A rare mystical gemstone",
		RecycleValue: map[string]int{"gem_fragments": 10, "magic_essence": 1},
	},

	// Gems - Cut
	"sapphire": {
		ID: "sapphire", Name: "Sapphire", Type: ItemTypeMaterial, Value: 300,
		Description:  "A cut blue gemstone",
		RecycleValue: map[string]int{"gem_fragments": 3},
	},
	"emerald": {
		ID: "emerald", Name: "Emerald", Type: ItemTypeMaterial, Value: 600,
		Description:  "A cut green gemstone",
		RecycleValue: map[string]int{"gem_fragments": 5},
	},
	"ruby": {
		ID: "ruby", Name: "Ruby", Type: ItemTypeMaterial, Value: 1200,
		Description:  "A cut red gemstone",
		RecycleValue: map[string]int{"gem_fragments": 7},
	},
	"diamond": {
		ID: "diamond", Name: "Diamond", Type: ItemTypeMaterial, Value: 3000,
		Description:  "A cut clear gemstone",
		RecycleValue: map[string]int{"gem_fragments": 10},
	},
	"dragonstone": {
		ID: "dragonstone", Name: "Dragonstone", Type: ItemTypeMaterial, Value: 10000,
		Description:  "A mystical cut gemstone",
		RecycleValue: map[string]int{"gem_fragments": 20, "magic_essence": 2},
	},

	// Additional Bars & Alloys
	"lead_bar": {
		ID: "lead_bar", Name: "Lead Bar", Type: ItemTypeBar, Value: 30,
		Description:  "Soft heavy metal bar",
		RecycleValue: map[string]int{"metal_fragments": 2, "lead_fragments": 1},
	},
	"brass_bar": {
		ID: "brass_bar", Name: "Brass Bar", Type: ItemTypeBar, Value: 40,
		Description:  "Copper-zinc alloy",
		RecycleValue: map[string]int{"metal_fragments": 2, "copper_fragments": 1, "zinc_fragments": 1},
	},
	"electrum_bar": {
		ID: "electrum_bar", Name: "Electrum Bar", Type: ItemTypeBar, Value: 200,
		Description:  "Gold-silver alloy",
		RecycleValue: map[string]int{"metal_fragments": 3, "gold_fragments": 1, "silver_fragments": 1},
	},
	"nickel_bar": {
		ID: "nickel_bar", Name: "Nickel Bar", Type: ItemTypeBar, Value: 80,
		Description:  "Corrosion-resistant bar",
		RecycleValue: map[string]int{"metal_fragments": 3, "nickel_fragments": 1},
	},
	"platinum_bar": {
		ID: "platinum_bar", Name: "Platinum Bar", Type: ItemTypeBar, Value: 1500,
		Description:  "Extremely valuable bar",
		RecycleValue: map[string]int{"metal_fragments": 5, "platinum_fragments": 2},
	},
	"obsidian_bar": {
		ID: "obsidian_bar", Name: "Obsidian Bar", Type: ItemTypeBar, Value: 2500,
		Description:  "Hardened volcanic glass",
		RecycleValue: map[string]int{"obsidian_fragments": 2},
	},
	"dragon_bar": {
		ID: "dragon_bar", Name: "Dragon Bar", Type: ItemTypeBar, Value: 10000,
		Description:  "Legendary metal alloy",
		RecycleValue: map[string]int{"metal_fragments": 10, "dragon_fragments": 1},
	},

	// Crafting Materials
	"clay": {
		ID: "clay", Name: "Clay", Type: ItemTypeResource, Value: 3,
		Description:  "Soft earth material",
		RecycleValue: map[string]int{"earth_fragments": 1},
	},
	"soft_clay": {
		ID: "soft_clay", Name: "Soft Clay", Type: ItemTypeMaterial, Value: 5,
		Description:  "Processed clay ready for molding",
		RecycleValue: map[string]int{"earth_fragments": 1},
	},
	"pottery": {
		ID: "pottery", Name: "Pottery", Type: ItemTypeMaterial, Value: 15,
		Description:  "Basic clay vessel",
		RecycleValue: map[string]int{"earth_fragments": 2},
	},
	"bowl": {
		ID: "bowl", Name: "Bowl", Type: ItemTypeMaterial, Value: 25,
		Description:  "Clay bowl for mixing",
		RecycleValue: map[string]int{"earth_fragments": 2},
	},
	"vase": {
		ID: "vase", Name: "Vase", Type: ItemTypeMaterial, Value: 50,
		Description:  "Decorative clay vase",
		RecycleValue: map[string]int{"earth_fragments": 3},
	},

	// Leather & Hides
	"cow_hide": {
		ID: "cow_hide", Name: "Cow Hide", Type: ItemTypeResource, Value: 10,
		Description:  "Basic animal hide",
		RecycleValue: map[string]int{"leather_fragments": 1},
	},
	"leather": {
		ID: "leather", Name: "Leather", Type: ItemTypeMaterial, Value: 30,
		Description:  "Treated leather material",
		RecycleValue: map[string]int{"leather_fragments": 2},
	},
	"hard_leather": {
		ID: "hard_leather", Name: "Hard Leather", Type: ItemTypeMaterial, Value: 80,
		Description:  "Reinforced leather",
		RecycleValue: map[string]int{"leather_fragments": 3},
	},
	"dragon_hide": {
		ID: "dragon_hide", Name: "Dragon Hide", Type: ItemTypeResource, Value: 2000,
		Description:  "Rare dragon scales",
		RecycleValue: map[string]int{"leather_fragments": 5, "dragon_fragments": 1},
	},
	"dragon_leather": {
		ID: "dragon_leather", Name: "Dragon Leather", Type: ItemTypeMaterial, Value: 5000,
		Description:  "Legendary dragon material",
		RecycleValue: map[string]int{"leather_fragments": 8, "dragon_fragments": 2},
	},

	// Additional Fragment Types
	"lead_fragments": {
		ID: "lead_fragments", Name: "Lead Fragments", Type: ItemTypeMaterial, Value: 4,
		Description: "Heavy soft metal fragments",
	},
	"zinc_fragments": {
		ID: "zinc_fragments", Name: "Zinc Fragments", Type: ItemTypeMaterial, Value: 5,
		Description: "Alloy-making fragments",
	},
	"nickel_fragments": {
		ID: "nickel_fragments", Name: "Nickel Fragments", Type: ItemTypeMaterial, Value: 10,
		Description: "Corrosion-resistant fragments",
	},
	"platinum_fragments": {
		ID: "platinum_fragments", Name: "Platinum Fragments", Type: ItemTypeMaterial, Value: 100,
		Description: "Precious platinum fragments",
	},
	"obsidian_fragments": {
		ID: "obsidian_fragments", Name: "Obsidian Fragments", Type: ItemTypeMaterial, Value: 150,
		Description: "Volcanic glass fragments",
	},
	"gem_fragments": {
		ID: "gem_fragments", Name: "Gem Fragments", Type: ItemTypeMaterial, Value: 50,
		Description: "Crushed gemstone pieces",
	},
	"earth_fragments": {
		ID: "earth_fragments", Name: "Earth Fragments", Type: ItemTypeMaterial, Value: 1,
		Description: "Basic earth material",
	},
	"leather_fragments": {
		ID: "leather_fragments", Name: "Leather Fragments", Type: ItemTypeMaterial, Value: 3,
		Description: "Scraps of leather",
	},
	"dragon_fragments": {
		ID: "dragon_fragments", Name: "Dragon Fragments", Type: ItemTypeMaterial, Value: 500,
		Description: "Rare dragon essence",
	},

	// Advanced Tools
	"mithril_axe": {
		ID: "mithril_axe", Name: "Mithril Axe", Type: ItemTypeTool, Value: 1200,
		Description: "Superior woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 8,
		Requirements: map[SkillType]int{SkillWoodcutting: 40},
		RecycleValue: map[string]int{"metal_fragments": 8, "mithril_fragments": 3},
	},
	"adamantite_axe": {
		ID: "adamantite_axe", Name: "Adamantite Axe", Type: ItemTypeTool, Value: 3000,
		Description: "Elite woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 12,
		Requirements: map[SkillType]int{SkillWoodcutting: 60},
		RecycleValue: map[string]int{"metal_fragments": 10, "adamantite_fragments": 3},
	},
	"runite_axe": {
		ID: "runite_axe", Name: "Runite Axe", Type: ItemTypeTool, Value: 8000,
		Description: "Ultimate woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 18,
		Requirements: map[SkillType]int{SkillWoodcutting: 85},
		RecycleValue: map[string]int{"metal_fragments": 12, "rune_fragments": 3},
	},
	"mithril_pickaxe": {
		ID: "mithril_pickaxe", Name: "Mithril Pickaxe", Type: ItemTypeTool, Value: 1200,
		Description: "Superior mining tool",
		Slot:        SlotWeapon, ToolPower: 8,
		Requirements: map[SkillType]int{SkillMining: 40},
		RecycleValue: map[string]int{"metal_fragments": 8, "mithril_fragments": 3},
	},
	"adamantite_pickaxe": {
		ID: "adamantite_pickaxe", Name: "Adamantite Pickaxe", Type: ItemTypeTool, Value: 3000,
		Description: "Elite mining tool",
		Slot:        SlotWeapon, ToolPower: 12,
		Requirements: map[SkillType]int{SkillMining: 60},
		RecycleValue: map[string]int{"metal_fragments": 10, "adamantite_fragments": 3},
	},
	"runite_pickaxe": {
		ID: "runite_pickaxe", Name: "Runite Pickaxe", Type: ItemTypeTool, Value: 8000,
		Description: "Ultimate mining tool",
		Slot:        SlotWeapon, ToolPower: 18,
		Requirements: map[SkillType]int{SkillMining: 85},
		RecycleValue: map[string]int{"metal_fragments": 12, "rune_fragments": 3},
	},
	"dragon_axe": {
		ID: "dragon_axe", Name: "Dragon Axe", Type: ItemTypeTool, Value: 25000,
		Description: "Legendary woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 25,
		Requirements: map[SkillType]int{SkillWoodcutting: 100},
		RecycleValue: map[string]int{"metal_fragments": 15, "dragon_fragments": 2},
	},
	"dragon_pickaxe": {
		ID: "dragon_pickaxe", Name: "Dragon Pickaxe", Type: ItemTypeTool, Value: 25000,
		Description: "Legendary mining tool",
		Slot:        SlotWeapon, ToolPower: 25,
		Requirements: map[SkillType]int{SkillMining: 100},
		RecycleValue: map[string]int{"metal_fragments": 15, "dragon_fragments": 2},
	},

	// Gem-Studded Tools (Special)
	"sapphire_axe": {
		ID: "sapphire_axe", Name: "Sapphire Axe", Type: ItemTypeTool, Value: 2500,
		Description: "Enchanted woodcutting tool",
		Slot:        SlotWeapon, ToolPower: 10,
		Requirements: map[SkillType]int{SkillWoodcutting: 50},
		RecycleValue: map[string]int{"metal_fragments": 6, "steel_fragments": 2, "gem_fragments": 5},
	},
	"emerald_pickaxe": {
		ID: "emerald_pickaxe", Name: "Emerald Pickaxe", Type: ItemTypeTool, Value: 4500,
		Description: "Enchanted mining tool",
		Slot:        SlotWeapon, ToolPower: 12,
		Requirements: map[SkillType]int{SkillMining: 60},
		RecycleValue: map[string]int{"metal_fragments": 7, "mithril_fragments": 2, "gem_fragments": 8},
	},
}

// GetItemTemplate retrieves an item template
func GetItemTemplate(id string) *Item {
	if template, ok := ItemDatabase[id]; ok {
		return template
	}
	return nil
}
