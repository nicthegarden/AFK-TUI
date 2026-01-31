package models

import (
	"fmt"
	"math"
)

// ActivityType represents different activity categories
type ActivityType string

const (
	ActivityGathering ActivityType = "gathering"
	ActivityCrafting  ActivityType = "crafting"
	ActivityCombat    ActivityType = "combat"
	ActivityRecycling ActivityType = "recycling"
)

// Activity represents what the player is currently doing
type Activity struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        ActivityType `json:"type"`
	SkillType   SkillType    `json:"skill_type"`

	// Requirements
	RequiredLevel int            `json:"required_level"`
	RequiredItems map[string]int `json:"required_items,omitempty"`

	// Rates
	BaseTicks   int            `json:"base_ticks"`   // Ticks to complete one action
	BaseXP      int64          `json:"base_xp"`      // XP per action
	OutputItems map[string]int `json:"output_items"` // Item ID -> quantity

	// Modifiers (populated at runtime)
	ToolPowerBonus  int     `json:"-"`
	XPMultiplier    float64 `json:"-"`
	SpeedMultiplier float64 `json:"-"`
	DoubleChance    float64 `json:"-"`

	// Progress tracking
	Progress       float64 `json:"progress"` // 0.0 to 1.0
	TicksRemaining int     `json:"ticks_remaining"`
}

// NewActivity creates an activity from a template
func NewActivity(id string) *Activity {
	if template, ok := ActivityDatabase[id]; ok {
		return &Activity{
			ID:             template.ID,
			Name:           template.Name,
			Description:    template.Description,
			Type:           template.Type,
			SkillType:      template.SkillType,
			RequiredLevel:  template.RequiredLevel,
			RequiredItems:  template.RequiredItems,
			BaseTicks:      template.BaseTicks,
			BaseXP:         template.BaseXP,
			OutputItems:    template.OutputItems,
			Progress:       0,
			TicksRemaining: template.BaseTicks,
		}
	}
	return nil
}

// ApplyModifiers applies player bonuses to the activity
func (a *Activity) ApplyModifiers(player *Player) {
	skill := player.GetSkill(a.SkillType)

	// Base tool power from equipment
	stats := player.Equipment.GetTotalStats()
	a.ToolPowerBonus = stats["tool_power"]

	// XP multiplier from perks
	a.XPMultiplier = 1.0 + player.GetSkillMultiplier(a.SkillType)

	// Speed multiplier from perks and tool
	a.SpeedMultiplier = 1.0 + player.GetSkillMultiplier(a.SkillType)
	a.SpeedMultiplier += float64(a.ToolPowerBonus) * 0.05 // 5% per tool power

	// Double drop chance from perks
	for _, perk := range player.UnlockedPerks {
		if perk.SkillType == a.SkillType && perk.Effect == PerkEffectDoubleDrop {
			a.DoubleChance += perk.Value
		}
	}

	// Level bonus to speed (1% per level after 10)
	if skill.Level > 10 {
		a.SpeedMultiplier += float64(skill.Level-10) * 0.01
	}

	// Calculate effective ticks
	effectiveTicks := float64(a.BaseTicks) / a.SpeedMultiplier
	a.TicksRemaining = int(math.Max(1, effectiveTicks))
}

// Tick processes one tick of the activity
func (a *Activity) Tick() bool {
	if a.TicksRemaining > 0 {
		a.TicksRemaining--
	}

	// Update progress based on remaining ticks vs total effective ticks
	// Total effective ticks = BaseTicks / SpeedMultiplier
	totalEffectiveTicks := float64(a.BaseTicks) / a.SpeedMultiplier
	if totalEffectiveTicks < 1 {
		totalEffectiveTicks = 1
	}

	// Progress = 1.0 - (remaining / total)
	remainingRatio := float64(a.TicksRemaining) / totalEffectiveTicks
	a.Progress = math.Min(1.0, math.Max(0.0, 1.0-remainingRatio))

	return a.TicksRemaining <= 0
}

// GetXP calculates XP with all multipliers
func (a *Activity) GetXP() int64 {
	return int64(float64(a.BaseXP) * a.XPMultiplier)
}

// GetOutput calculates output items with double chance
func (a *Activity) GetOutput() map[string]int {
	output := make(map[string]int)

	for itemID, quantity := range a.OutputItems {
		dropMult := 1

		// Check double/triple drop
		if a.DoubleChance > 0 {
			// Simple implementation - guaranteed extra at certain thresholds
			if a.DoubleChance >= 1.0 {
				dropMult = 2
			}
		}

		output[itemID] = quantity * dropMult
	}

	return output
}

// Reset resets progress for next action
func (a *Activity) Reset() {
	a.Progress = 0
	a.TicksRemaining = int(float64(a.BaseTicks) / a.SpeedMultiplier)
}

// CanDo checks if player can perform this activity
func (a *Activity) CanDo(player *Player) error {
	skill := player.GetSkill(a.SkillType)

	if skill.Level < a.RequiredLevel {
		return fmt.Errorf("requires level %d %s (you have %d)",
			a.RequiredLevel, SkillNames[a.SkillType], skill.Level)
	}

	// Check required items
	for itemID, quantity := range a.RequiredItems {
		if !player.Inventory.HasItem(itemID, quantity) {
			item := GetItemTemplate(itemID)
			itemName := itemID
			if item != nil {
				itemName = item.Name
			}
			return fmt.Errorf("need %d %s", quantity, itemName)
		}
	}

	return nil
}

// ActivityDatabase contains all activities
type ActivityTemplate struct {
	ID            string
	Name          string
	Description   string
	Type          ActivityType
	SkillType     SkillType
	RequiredLevel int
	RequiredItems map[string]int
	BaseTicks     int
	BaseXP        int64
	OutputItems   map[string]int
}

var ActivityDatabase = map[string]*ActivityTemplate{
	// Woodcutting
	"chop_logs": {
		ID: "chop_logs", Name: "Chop Logs", Description: "Chop basic trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 1, BaseTicks: 4, BaseXP: 10,
		OutputItems: map[string]int{"logs": 1},
	},
	"chop_oak": {
		ID: "chop_oak", Name: "Chop Oak", Description: "Chop oak trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 15, BaseTicks: 6, BaseXP: 20,
		OutputItems: map[string]int{"oak_logs": 1},
	},
	"chop_willow": {
		ID: "chop_willow", Name: "Chop Willow", Description: "Chop willow trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 30, BaseTicks: 8, BaseXP: 35,
		OutputItems: map[string]int{"willow_logs": 1},
	},
	"chop_maple": {
		ID: "chop_maple", Name: "Chop Maple", Description: "Chop maple trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 45, BaseTicks: 10, BaseXP: 55,
		OutputItems: map[string]int{"maple_logs": 1},
	},
	"chop_yew": {
		ID: "chop_yew", Name: "Chop Yew", Description: "Chop yew trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 60, BaseTicks: 14, BaseXP: 85,
		OutputItems: map[string]int{"yew_logs": 1},
	},
	"chop_magic": {
		ID: "chop_magic", Name: "Chop Magic", Description: "Chop magic trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 75, BaseTicks: 20, BaseXP: 125,
		OutputItems: map[string]int{"magic_logs": 1},
	},

	// Mining
	"mine_copper": {
		ID: "mine_copper", Name: "Mine Copper", Description: "Mine copper ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 1, BaseTicks: 5, BaseXP: 12,
		OutputItems: map[string]int{"copper_ore": 1},
	},
	"mine_tin": {
		ID: "mine_tin", Name: "Mine Tin", Description: "Mine tin ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 1, BaseTicks: 5, BaseXP: 12,
		OutputItems: map[string]int{"tin_ore": 1},
	},
	"mine_iron": {
		ID: "mine_iron", Name: "Mine Iron", Description: "Mine iron ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 15, BaseTicks: 7, BaseXP: 25,
		OutputItems: map[string]int{"iron_ore": 1},
	},
	"mine_coal": {
		ID: "mine_coal", Name: "Mine Coal", Description: "Mine coal",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 30, BaseTicks: 8, BaseXP: 35,
		OutputItems: map[string]int{"coal": 1},
	},
	"mine_silver": {
		ID: "mine_silver", Name: "Mine Silver", Description: "Mine silver ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 40, BaseTicks: 12, BaseXP: 50,
		OutputItems: map[string]int{"silver_ore": 1},
	},
	"mine_gold": {
		ID: "mine_gold", Name: "Mine Gold", Description: "Mine gold ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 50, BaseTicks: 15, BaseXP: 65,
		OutputItems: map[string]int{"gold_ore": 1},
	},
	"mine_mithril": {
		ID: "mine_mithril", Name: "Mine Mithril", Description: "Mine mithril ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 65, BaseTicks: 20, BaseXP: 90,
		OutputItems: map[string]int{"mithril_ore": 1},
	},
	"mine_adamantite": {
		ID: "mine_adamantite", Name: "Mine Adamantite", Description: "Mine adamantite ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 80, BaseTicks: 30, BaseXP: 120,
		OutputItems: map[string]int{"adamantite_ore": 1},
	},
	"mine_runite": {
		ID: "mine_runite", Name: "Mine Runite", Description: "Mine runite ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 95, BaseTicks: 45, BaseXP: 160,
		OutputItems: map[string]int{"runite_ore": 1},
	},

	// Smithing
	"smelt_bronze": {
		ID: "smelt_bronze", Name: "Smelt Bronze", Description: "Smelt bronze bar (1 copper + 1 tin)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 1, BaseTicks: 6, BaseXP: 15,
		RequiredItems: map[string]int{"copper_ore": 1, "tin_ore": 1},
		OutputItems:   map[string]int{"bronze_bar": 1},
	},
	"smelt_iron": {
		ID: "smelt_iron", Name: "Smelt Iron", Description: "Smelt iron bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 15, BaseTicks: 8, BaseXP: 30,
		RequiredItems: map[string]int{"iron_ore": 1},
		OutputItems:   map[string]int{"iron_bar": 1},
	},
	"smelt_steel": {
		ID: "smelt_steel", Name: "Smelt Steel", Description: "Smelt steel bar (1 iron + 2 coal)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 30, BaseTicks: 10, BaseXP: 45,
		RequiredItems: map[string]int{"iron_ore": 1, "coal": 2},
		OutputItems:   map[string]int{"steel_bar": 1},
	},
	"smelt_silver": {
		ID: "smelt_silver", Name: "Smelt Silver", Description: "Smelt silver bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 40, BaseTicks: 12, BaseXP: 55,
		RequiredItems: map[string]int{"silver_ore": 1},
		OutputItems:   map[string]int{"silver_bar": 1},
	},
	"smelt_gold": {
		ID: "smelt_gold", Name: "Smelt Gold", Description: "Smelt gold bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 50, BaseTicks: 15, BaseXP: 70,
		RequiredItems: map[string]int{"gold_ore": 1},
		OutputItems:   map[string]int{"gold_bar": 1},
	},
	"smelt_mithril": {
		ID: "smelt_mithril", Name: "Smelt Mithril", Description: "Smelt mithril bar (1 mithril + 4 coal)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 65, BaseTicks: 20, BaseXP: 95,
		RequiredItems: map[string]int{"mithril_ore": 1, "coal": 4},
		OutputItems:   map[string]int{"mithril_bar": 1},
	},
	"smelt_adamantite": {
		ID: "smelt_adamantite", Name: "Smelt Adamantite", Description: "Smelt adamantite bar (1 adamantite + 6 coal)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 80, BaseTicks: 30, BaseXP: 125,
		RequiredItems: map[string]int{"adamantite_ore": 1, "coal": 6},
		OutputItems:   map[string]int{"adamantite_bar": 1},
	},
	"smelt_runite": {
		ID: "smelt_runite", Name: "Smelt Runite", Description: "Smelt runite bar (1 runite + 8 coal)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 95, BaseTicks: 45, BaseXP: 165,
		RequiredItems: map[string]int{"runite_ore": 1, "coal": 8},
		OutputItems:   map[string]int{"runite_bar": 1},
	},

	// Crafting with recycled materials
	"smith_bronze_axe": {
		ID: "smith_bronze_axe", Name: "Smith Bronze Axe", Description: "Craft bronze axe from fragments",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 5, BaseTicks: 10, BaseXP: 25,
		RequiredItems: map[string]int{"bronze_bar": 1, "wood_fragments": 5},
		OutputItems:   map[string]int{"bronze_axe": 1},
	},
	"smith_bronze_pickaxe": {
		ID: "smith_bronze_pickaxe", Name: "Smith Bronze Pickaxe", Description: "Craft bronze pickaxe from fragments",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 5, BaseTicks: 10, BaseXP: 25,
		RequiredItems: map[string]int{"bronze_bar": 1, "wood_fragments": 5},
		OutputItems:   map[string]int{"bronze_pickaxe": 1},
	},
	"smith_iron_axe": {
		ID: "smith_iron_axe", Name: "Smith Iron Axe", Description: "Craft iron axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 20, BaseTicks: 12, BaseXP: 45,
		RequiredItems: map[string]int{"iron_bar": 2, "wood_fragments": 10},
		OutputItems:   map[string]int{"iron_axe": 1},
	},
	"smith_iron_pickaxe": {
		ID: "smith_iron_pickaxe", Name: "Smith Iron Pickaxe", Description: "Craft iron pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 20, BaseTicks: 12, BaseXP: 45,
		RequiredItems: map[string]int{"iron_bar": 2, "wood_fragments": 10},
		OutputItems:   map[string]int{"iron_pickaxe": 1},
	},
	"smith_steel_axe": {
		ID: "smith_steel_axe", Name: "Smith Steel Axe", Description: "Craft steel axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 35, BaseTicks: 15, BaseXP: 70,
		RequiredItems: map[string]int{"steel_bar": 2, "oak_logs": 2},
		OutputItems:   map[string]int{"steel_axe": 1},
	},
	"smith_steel_pickaxe": {
		ID: "smith_steel_pickaxe", Name: "Smith Steel Pickaxe", Description: "Craft steel pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 35, BaseTicks: 15, BaseXP: 70,
		RequiredItems: map[string]int{"steel_bar": 2, "oak_logs": 2},
		OutputItems:   map[string]int{"steel_pickaxe": 1},
	},

	// Recycling
	"recycle_logs": {
		ID: "recycle_logs", Name: "Recycle Logs", Description: "Recycle logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 1, BaseTicks: 3, BaseXP: 8,
		RequiredItems: map[string]int{"logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 1},
	},
	"recycle_oak_logs": {
		ID: "recycle_oak_logs", Name: "Recycle Oak Logs", Description: "Recycle oak logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 15, BaseTicks: 4, BaseXP: 15,
		RequiredItems: map[string]int{"oak_logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 2},
	},
	"recycle_bronze_items": {
		ID: "recycle_bronze_items", Name: "Recycle Bronze Items", Description: "Recycle bronze equipment",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 10, BaseTicks: 5, BaseXP: 20,
		RequiredItems: map[string]int{"bronze_sword": 1},
		OutputItems:   map[string]int{"metal_fragments": 3, "copper_fragments": 2},
	},

	// Extended Mining - Additional Ores
	"mine_lead": {
		ID: "mine_lead", Name: "Mine Lead", Description: "Mine lead ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 10, BaseTicks: 6, BaseXP: 18,
		OutputItems: map[string]int{"lead_ore": 1},
	},
	"mine_zinc": {
		ID: "mine_zinc", Name: "Mine Zinc", Description: "Mine zinc ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 12, BaseTicks: 6, BaseXP: 20,
		OutputItems: map[string]int{"zinc_ore": 1},
	},
	"mine_nickel": {
		ID: "mine_nickel", Name: "Mine Nickel", Description: "Mine nickel ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 25, BaseTicks: 9, BaseXP: 32,
		OutputItems: map[string]int{"nickel_ore": 1},
	},
	"mine_platinum": {
		ID: "mine_platinum", Name: "Mine Platinum", Description: "Mine platinum ore (rare!)",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 70, BaseTicks: 25, BaseXP: 110,
		OutputItems: map[string]int{"platinum_ore": 1},
	},
	"mine_obsidian": {
		ID: "mine_obsidian", Name: "Mine Obsidian", Description: "Mine obsidian ore",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 90, BaseTicks: 40, BaseXP: 150,
		OutputItems: map[string]int{"obsidian_ore": 1},
	},

	// Gem Mining
	"mine_sapphire": {
		ID: "mine_sapphire", Name: "Mine Sapphire", Description: "Mine uncut sapphire",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 20, BaseTicks: 10, BaseXP: 30,
		OutputItems: map[string]int{"uncut_sapphire": 1},
	},
	"mine_emerald": {
		ID: "mine_emerald", Name: "Mine Emerald", Description: "Mine uncut emerald",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 35, BaseTicks: 14, BaseXP: 50,
		OutputItems: map[string]int{"uncut_emerald": 1},
	},
	"mine_ruby": {
		ID: "mine_ruby", Name: "Mine Ruby", Description: "Mine uncut ruby",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 55, BaseTicks: 20, BaseXP: 75,
		OutputItems: map[string]int{"uncut_ruby": 1},
	},
	"mine_diamond": {
		ID: "mine_diamond", Name: "Mine Diamond", Description: "Mine uncut diamond",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 75, BaseTicks: 30, BaseXP: 110,
		OutputItems: map[string]int{"uncut_diamond": 1},
	},
	"mine_dragonstone": {
		ID: "mine_dragonstone", Name: "Mine Dragonstone", Description: "Mine uncut dragonstone (legendary!)",
		Type: ActivityGathering, SkillType: SkillMining,
		RequiredLevel: 100, BaseTicks: 60, BaseXP: 200,
		OutputItems: map[string]int{"uncut_dragonstone": 1},
	},

	// Additional Smelting - Alloys
	"smelt_lead": {
		ID: "smelt_lead", Name: "Smelt Lead", Description: "Smelt lead bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 10, BaseTicks: 7, BaseXP: 22,
		RequiredItems: map[string]int{"lead_ore": 1},
		OutputItems:   map[string]int{"lead_bar": 1},
	},
	"smelt_brass": {
		ID: "smelt_brass", Name: "Smelt Brass", Description: "Smelt brass bar (1 copper + 1 zinc)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 15, BaseTicks: 8, BaseXP: 28,
		RequiredItems: map[string]int{"copper_ore": 1, "zinc_ore": 1},
		OutputItems:   map[string]int{"brass_bar": 1},
	},
	"smelt_electrum": {
		ID: "smelt_electrum", Name: "Smelt Electrum", Description: "Smelt electrum bar (1 gold + 1 silver)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 55, BaseTicks: 18, BaseXP: 85,
		RequiredItems: map[string]int{"gold_ore": 1, "silver_ore": 1},
		OutputItems:   map[string]int{"electrum_bar": 1},
	},
	"smelt_nickel": {
		ID: "smelt_nickel", Name: "Smelt Nickel", Description: "Smelt nickel bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 30, BaseTicks: 11, BaseXP: 40,
		RequiredItems: map[string]int{"nickel_ore": 1},
		OutputItems:   map[string]int{"nickel_bar": 1},
	},
	"smelt_platinum": {
		ID: "smelt_platinum", Name: "Smelt Platinum", Description: "Smelt platinum bar",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 75, BaseTicks: 28, BaseXP: 130,
		RequiredItems: map[string]int{"platinum_ore": 1, "coal": 4},
		OutputItems:   map[string]int{"platinum_bar": 1},
	},
	"smelt_obsidian": {
		ID: "smelt_obsidian", Name: "Smelt Obsidian", Description: "Smelt obsidian bar (2 obsidian + 2 coal)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 95, BaseTicks: 50, BaseXP: 180,
		RequiredItems: map[string]int{"obsidian_ore": 2, "coal": 2},
		OutputItems:   map[string]int{"obsidian_bar": 1},
	},

	// Gem Cutting (requires Crafting skill - using Smithing for now)
	"cut_sapphire": {
		ID: "cut_sapphire", Name: "Cut Sapphire", Description: "Cut sapphire gem",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 20, BaseTicks: 8, BaseXP: 35,
		RequiredItems: map[string]int{"uncut_sapphire": 1},
		OutputItems:   map[string]int{"sapphire": 1},
	},
	"cut_emerald": {
		ID: "cut_emerald", Name: "Cut Emerald", Description: "Cut emerald gem",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 35, BaseTicks: 12, BaseXP: 60,
		RequiredItems: map[string]int{"uncut_emerald": 1},
		OutputItems:   map[string]int{"emerald": 1},
	},
	"cut_ruby": {
		ID: "cut_ruby", Name: "Cut Ruby", Description: "Cut ruby gem",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 55, BaseTicks: 18, BaseXP: 90,
		RequiredItems: map[string]int{"uncut_ruby": 1},
		OutputItems:   map[string]int{"ruby": 1},
	},
	"cut_diamond": {
		ID: "cut_diamond", Name: "Cut Diamond", Description: "Cut diamond gem",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 75, BaseTicks: 28, BaseXP: 130,
		RequiredItems: map[string]int{"uncut_diamond": 1},
		OutputItems:   map[string]int{"diamond": 1},
	},
	"cut_dragonstone": {
		ID: "cut_dragonstone", Name: "Cut Dragonstone", Description: "Cut dragonstone gem",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 100, BaseTicks: 45, BaseXP: 220,
		RequiredItems: map[string]int{"uncut_dragonstone": 1},
		OutputItems:   map[string]int{"dragonstone": 1},
	},

	// Pottery (Crafting)
	"gather_clay": {
		ID: "gather_clay", Name: "Gather Clay", Description: "Gather clay from riverbeds",
		Type: ActivityGathering, SkillType: SkillCrafting,
		RequiredLevel: 1, BaseTicks: 4, BaseXP: 8,
		OutputItems: map[string]int{"clay": 1},
	},
	"soften_clay": {
		ID: "soften_clay", Name: "Soften Clay", Description: "Process clay for molding",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 5, BaseTicks: 5, BaseXP: 12,
		RequiredItems: map[string]int{"clay": 1},
		OutputItems:   map[string]int{"soft_clay": 1},
	},
	"make_pottery": {
		ID: "make_pottery", Name: "Make Pottery", Description: "Craft basic pottery",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 10, BaseTicks: 8, BaseXP: 20,
		RequiredItems: map[string]int{"soft_clay": 2},
		OutputItems:   map[string]int{"pottery": 1},
	},
	"make_bowl": {
		ID: "make_bowl", Name: "Make Bowl", Description: "Craft clay bowl",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 15, BaseTicks: 10, BaseXP: 28,
		RequiredItems: map[string]int{"soft_clay": 3},
		OutputItems:   map[string]int{"bowl": 1},
	},
	"make_vase": {
		ID: "make_vase", Name: "Make Vase", Description: "Craft decorative vase",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 25, BaseTicks: 14, BaseXP: 45,
		RequiredItems: map[string]int{"soft_clay": 4},
		OutputItems:   map[string]int{"vase": 1},
	},

	// Tanning (Crafting)
	"tan_leather": {
		ID: "tan_leather", Name: "Tan Leather", Description: "Process cow hide into leather",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 5, BaseTicks: 6, BaseXP: 15,
		RequiredItems: map[string]int{"cow_hide": 1},
		OutputItems:   map[string]int{"leather": 1},
	},
	"make_hard_leather": {
		ID: "make_hard_leather", Name: "Make Hard Leather", Description: "Reinforce leather material",
		Type: ActivityCrafting, SkillType: SkillCrafting,
		RequiredLevel: 30, BaseTicks: 12, BaseXP: 40,
		RequiredItems: map[string]int{"leather": 2, "nickel_bar": 1},
		OutputItems:   map[string]int{"hard_leather": 1},
	},

	// Advanced Tool Crafting
	"smith_mithril_axe": {
		ID: "smith_mithril_axe", Name: "Smith Mithril Axe", Description: "Craft mithril axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 45, BaseTicks: 18, BaseXP: 90,
		RequiredItems: map[string]int{"mithril_bar": 2, "willow_logs": 2},
		OutputItems:   map[string]int{"mithril_axe": 1},
	},
	"smith_mithril_pickaxe": {
		ID: "smith_mithril_pickaxe", Name: "Smith Mithril Pickaxe", Description: "Craft mithril pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 45, BaseTicks: 18, BaseXP: 90,
		RequiredItems: map[string]int{"mithril_bar": 2, "willow_logs": 2},
		OutputItems:   map[string]int{"mithril_pickaxe": 1},
	},
	"smith_adamantite_axe": {
		ID: "smith_adamantite_axe", Name: "Smith Adamantite Axe", Description: "Craft adamantite axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 65, BaseTicks: 25, BaseXP: 120,
		RequiredItems: map[string]int{"adamantite_bar": 2, "maple_logs": 2},
		OutputItems:   map[string]int{"adamantite_axe": 1},
	},
	"smith_adamantite_pickaxe": {
		ID: "smith_adamantite_pickaxe", Name: "Smith Adamantite Pickaxe", Description: "Craft adamantite pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 65, BaseTicks: 25, BaseXP: 120,
		RequiredItems: map[string]int{"adamantite_bar": 2, "maple_logs": 2},
		OutputItems:   map[string]int{"adamantite_pickaxe": 1},
	},
	"smith_runite_axe": {
		ID: "smith_runite_axe", Name: "Smith Runite Axe", Description: "Craft runite axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 90, BaseTicks: 40, BaseXP: 160,
		RequiredItems: map[string]int{"runite_bar": 2, "yew_logs": 2},
		OutputItems:   map[string]int{"runite_axe": 1},
	},
	"smith_runite_pickaxe": {
		ID: "smith_runite_pickaxe", Name: "Smith Runite Pickaxe", Description: "Craft runite pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 90, BaseTicks: 40, BaseXP: 160,
		RequiredItems: map[string]int{"runite_bar": 2, "yew_logs": 2},
		OutputItems:   map[string]int{"runite_pickaxe": 1},
	},
	"smith_dragon_axe": {
		ID: "smith_dragon_axe", Name: "Smith Dragon Axe", Description: "Craft legendary dragon axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 105, BaseTicks: 60, BaseXP: 250,
		RequiredItems: map[string]int{"dragon_bar": 2, "magic_logs": 2, "dragonstone": 1},
		OutputItems:   map[string]int{"dragon_axe": 1},
	},
	"smith_dragon_pickaxe": {
		ID: "smith_dragon_pickaxe", Name: "Smith Dragon Pickaxe", Description: "Craft legendary dragon pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 105, BaseTicks: 60, BaseXP: 250,
		RequiredItems: map[string]int{"dragon_bar": 2, "magic_logs": 2, "dragonstone": 1},
		OutputItems:   map[string]int{"dragon_pickaxe": 1},
	},

	// Training Activities
	"strength_training": {
		ID: "strength_training", Name: "Strength Training", Description: "Train strength at the training dummy",
		Type: ActivityCombat, SkillType: SkillCombat,
		RequiredLevel: 1, BaseTicks: 8, BaseXP: 15,
		OutputItems: map[string]int{}, // No items, trains attributes directly
	},
	"dexterity_training": {
		ID: "dexterity_training", Name: "Dexterity Training", Description: "Train dexterity on the agility course",
		Type: ActivityCombat, SkillType: SkillCombat,
		RequiredLevel: 1, BaseTicks: 8, BaseXP: 15,
		OutputItems: map[string]int{},
	},
	"defense_training": {
		ID: "defense_training", Name: "Defense Training", Description: "Train defense with shield drills",
		Type: ActivityCombat, SkillType: SkillCombat,
		RequiredLevel: 1, BaseTicks: 8, BaseXP: 15,
		OutputItems: map[string]int{},
	},

	// Gem-Studded Tools
	"smith_sapphire_axe": {
		ID: "smith_sapphire_axe", Name: "Smith Sapphire Axe", Description: "Craft enchanted sapphire axe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 55, BaseTicks: 22, BaseXP: 100,
		RequiredItems: map[string]int{"steel_bar": 2, "sapphire": 1, "oak_logs": 2},
		OutputItems:   map[string]int{"sapphire_axe": 1},
	},
	"smith_emerald_pickaxe": {
		ID: "smith_emerald_pickaxe", Name: "Smith Emerald Pickaxe", Description: "Craft enchanted emerald pickaxe",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 65, BaseTicks: 25, BaseXP: 115,
		RequiredItems: map[string]int{"mithril_bar": 2, "emerald": 1, "willow_logs": 2},
		OutputItems:   map[string]int{"emerald_pickaxe": 1},
	},

	// Extended Woodcutting
	"chop_teak": {
		ID: "chop_teak", Name: "Chop Teak", Description: "Chop teak trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 90, BaseTicks: 25, BaseXP: 145,
		OutputItems: map[string]int{"teak_logs": 1},
	},
	"chop_mahogany": {
		ID: "chop_mahogany", Name: "Chop Mahogany", Description: "Chop mahogany trees",
		Type: ActivityGathering, SkillType: SkillWoodcutting,
		RequiredLevel: 105, BaseTicks: 35, BaseXP: 200,
		OutputItems: map[string]int{"mahogany_logs": 1},
	},

	// Extended Recycling
	"recycle_willow_logs": {
		ID: "recycle_willow_logs", Name: "Recycle Willow Logs", Description: "Recycle willow logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 30, BaseTicks: 5, BaseXP: 22,
		RequiredItems: map[string]int{"willow_logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 3},
	},
	"recycle_maple_logs": {
		ID: "recycle_maple_logs", Name: "Recycle Maple Logs", Description: "Recycle maple logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 45, BaseTicks: 6, BaseXP: 30,
		RequiredItems: map[string]int{"maple_logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 4},
	},
	"recycle_yew_logs": {
		ID: "recycle_yew_logs", Name: "Recycle Yew Logs", Description: "Recycle yew logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 60, BaseTicks: 8, BaseXP: 42,
		RequiredItems: map[string]int{"yew_logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 5},
	},
	"recycle_magic_logs": {
		ID: "recycle_magic_logs", Name: "Recycle Magic Logs", Description: "Recycle magic logs into fragments",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 75, BaseTicks: 12, BaseXP: 60,
		RequiredItems: map[string]int{"magic_logs": 1},
		OutputItems:   map[string]int{"wood_fragments": 8, "magic_essence": 1},
	},
	"recycle_iron_items": {
		ID: "recycle_iron_items", Name: "Recycle Iron Items", Description: "Recycle iron equipment",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 20, BaseTicks: 6, BaseXP: 35,
		RequiredItems: map[string]int{"iron_sword": 1},
		OutputItems:   map[string]int{"metal_fragments": 5, "iron_fragments": 2},
	},
	"recycle_steel_items": {
		ID: "recycle_steel_items", Name: "Recycle Steel Items", Description: "Recycle steel equipment",
		Type: ActivityRecycling, SkillType: SkillRecycling,
		RequiredLevel: 35, BaseTicks: 8, BaseXP: 55,
		RequiredItems: map[string]int{"steel_sword": 1},
		OutputItems:   map[string]int{"metal_fragments": 7, "steel_fragments": 3},
	},

	// Advanced Smelting - Legendary
	"smelt_dragon": {
		ID: "smelt_dragon", Name: "Smelt Dragon Bar", Description: "Smelt legendary dragon bar (2 mithril + 1 runite + 4 coal + 1 dragonstone)",
		Type: ActivityCrafting, SkillType: SkillSmithing,
		RequiredLevel: 110, BaseTicks: 80, BaseXP: 300,
		RequiredItems: map[string]int{"mithril_bar": 2, "runite_bar": 1, "coal": 4, "dragonstone": 1},
		OutputItems:   map[string]int{"dragon_bar": 1},
	},
}

// GetActivitiesForSkill returns all activities for a skill
func GetActivitiesForSkill(skillType SkillType) []*Activity {
	var activities []*Activity
	for _, template := range ActivityDatabase {
		if template.SkillType == skillType {
			activities = append(activities, NewActivity(template.ID))
		}
	}
	return activities
}
