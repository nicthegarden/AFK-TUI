package models

import "fmt"

// PerkEffect types
type PerkEffect string

const (
	PerkEffectXPBoost     PerkEffect = "xp_boost"
	PerkEffectSpeedBoost  PerkEffect = "speed_boost"
	PerkEffectDoubleDrop  PerkEffect = "double_drop"
	PerkEffectAutoCollect PerkEffect = "auto_collect"
	PerkEffectExtraSlot   PerkEffect = "extra_slot"
	PerkEffectGoldBoost   PerkEffect = "gold_boost"
)

// Perk represents a permanent bonus
type Perk struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SkillType   SkillType  `json:"skill_type"`
	LevelReq    int        `json:"level_req"`
	Effect      PerkEffect `json:"effect"`
	Value       float64    `json:"value"`
}

// PerkDatabase contains all perks
type PerkDatabase struct {
	perks []Perk
}

// AllPerks contains every perk in the game
var AllPerks = []Perk{
	// Woodcutting Perks
	{ID: "wc_speed_1", Name: "Quick Chop", Description: "10% faster woodcutting", SkillType: SkillWoodcutting, LevelReq: 5, Effect: PerkEffectSpeedBoost, Value: 0.10},
	{ID: "wc_xp_1", Name: "Nature's Wisdom", Description: "15% more Woodcutting XP", SkillType: SkillWoodcutting, LevelReq: 10, Effect: PerkEffectXPBoost, Value: 0.15},
	{ID: "wc_double", Name: "Double Logs", Description: "5% chance for double logs", SkillType: SkillWoodcutting, LevelReq: 20, Effect: PerkEffectDoubleDrop, Value: 0.05},
	{ID: "wc_speed_2", Name: "Expert Chopper", Description: "20% faster woodcutting", SkillType: SkillWoodcutting, LevelReq: 35, Effect: PerkEffectSpeedBoost, Value: 0.20},
	{ID: "wc_xp_2", Name: "Forest Mastery", Description: "25% more Woodcutting XP", SkillType: SkillWoodcutting, LevelReq: 50, Effect: PerkEffectXPBoost, Value: 0.25},
	{ID: "wc_triple", Name: "Triple Logs", Description: "10% chance for triple logs", SkillType: SkillWoodcutting, LevelReq: 80, Effect: PerkEffectDoubleDrop, Value: 0.10},

	// Mining Perks
	{ID: "mining_speed_1", Name: "Swift Pick", Description: "10% faster mining", SkillType: SkillMining, LevelReq: 5, Effect: PerkEffectSpeedBoost, Value: 0.10},
	{ID: "mining_xp_1", Name: "Rock Sense", Description: "15% more Mining XP", SkillType: SkillMining, LevelReq: 10, Effect: PerkEffectXPBoost, Value: 0.15},
	{ID: "mining_double", Name: "Double Ore", Description: "5% chance for double ore", SkillType: SkillMining, LevelReq: 20, Effect: PerkEffectDoubleDrop, Value: 0.05},
	{ID: "mining_speed_2", Name: "Expert Miner", Description: "20% faster mining", SkillType: SkillMining, LevelReq: 35, Effect: PerkEffectSpeedBoost, Value: 0.20},
	{ID: "mining_xp_2", Name: "Earth Mastery", Description: "25% more Mining XP", SkillType: SkillMining, LevelReq: 50, Effect: PerkEffectXPBoost, Value: 0.25},
	{ID: "mining_triple", Name: "Triple Ore", Description: "10% chance for triple ore", SkillType: SkillMining, LevelReq: 80, Effect: PerkEffectDoubleDrop, Value: 0.10},

	// Smithing Perks
	{ID: "smith_xp_1", Name: "Apprentice Smith", Description: "15% more Smithing XP", SkillType: SkillSmithing, LevelReq: 10, Effect: PerkEffectXPBoost, Value: 0.15},
	{ID: "smith_double", Name: "Efficient Smith", Description: "10% chance to save bars", SkillType: SkillSmithing, LevelReq: 25, Effect: PerkEffectDoubleDrop, Value: 0.10},
	{ID: "smith_xp_2", Name: "Master Smith", Description: "25% more Smithing XP", SkillType: SkillSmithing, LevelReq: 50, Effect: PerkEffectXPBoost, Value: 0.25},
	{ID: "smith_save", Name: "Bar Conservation", Description: "20% chance to save bars", SkillType: SkillSmithing, LevelReq: 80, Effect: PerkEffectDoubleDrop, Value: 0.20},

	// Recycling Perks
	{ID: "recycle_xp_1", Name: "Scavenger", Description: "15% more Recycling XP", SkillType: SkillRecycling, LevelReq: 10, Effect: PerkEffectXPBoost, Value: 0.15},
	{ID: "recycle_bonus", Name: "Bonus Materials", Description: "25% more materials from recycling", SkillType: SkillRecycling, LevelReq: 25, Effect: PerkEffectDoubleDrop, Value: 0.25},
	{ID: "recycle_xp_2", Name: "Master Recycler", Description: "30% more Recycling XP", SkillType: SkillRecycling, LevelReq: 50, Effect: PerkEffectXPBoost, Value: 0.30},
	{ID: "recycle_super", Name: "Super Recycler", Description: "50% more materials from recycling", SkillType: SkillRecycling, LevelReq: 90, Effect: PerkEffectDoubleDrop, Value: 0.50},

	// Combat Perks
	{ID: "combat_xp_1", Name: "Warrior's Path", Description: "15% more Combat XP", SkillType: SkillCombat, LevelReq: 10, Effect: PerkEffectXPBoost, Value: 0.15},
	{ID: "combat_gold", Name: "Loot Master", Description: "25% more gold from combat", SkillType: SkillCombat, LevelReq: 20, Effect: PerkEffectGoldBoost, Value: 0.25},
	{ID: "combat_xp_2", Name: "Battle Veteran", Description: "25% more Combat XP", SkillType: SkillCombat, LevelReq: 50, Effect: PerkEffectXPBoost, Value: 0.25},

	// Global Perks (any skill)
	{ID: "global_slot", Name: "Extra Storage", Description: "+5 inventory slots", SkillType: "", LevelReq: 30, Effect: PerkEffectExtraSlot, Value: 5},
	{ID: "global_auto", Name: "Auto-Collector", Description: "Automatically collect resources", SkillType: "", LevelReq: 60, Effect: PerkEffectAutoCollect, Value: 1},
}

// GetPerksForLevel returns perks unlocked at a specific level
func GetPerksForLevel(skillType SkillType, level int) []Perk {
	var unlocked []Perk
	for _, perk := range AllPerks {
		if perk.SkillType == skillType && perk.LevelReq == level {
			unlocked = append(unlocked, perk)
		}
	}
	return unlocked
}

// GetAllPerksForSkill returns all perks for a skill
func GetAllPerksForSkill(skillType SkillType) []Perk {
	var perks []Perk
	for _, perk := range AllPerks {
		if perk.SkillType == skillType {
			perks = append(perks, perk)
		}
	}
	return perks
}

// PerkManager tracks unlocked perks
type PerkManager struct {
	Unlocked map[string]bool
	Perks    []Perk
}

// NewPerkManager creates a new perk manager
func NewPerkManager() *PerkManager {
	return &PerkManager{
		Unlocked: make(map[string]bool),
		Perks:    []Perk{},
	}
}

// UnlockPerk adds a perk if not already unlocked
func (pm *PerkManager) UnlockPerk(perk Perk) bool {
	if pm.Unlocked[perk.ID] {
		return false
	}
	pm.Unlocked[perk.ID] = true
	pm.Perks = append(pm.Perks, perk)
	return true
}

// GetBonus returns the total bonus for an effect type
func (pm *PerkManager) GetBonus(skillType SkillType, effect PerkEffect) float64 {
	bonus := 0.0
	for _, perk := range pm.Perks {
		if perk.Effect == effect && (perk.SkillType == skillType || perk.SkillType == "") {
			bonus += perk.Value
		}
	}
	return bonus
}

// String representation
func (p Perk) String() string {
	return fmt.Sprintf("%s (%s +%.0f%%)", p.Name, p.Effect, p.Value*100)
}
