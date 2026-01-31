package models

import (
	"fmt"
	"math"
	"time"
)

// SkillType represents different skills
type SkillType string

const (
	SkillWoodcutting SkillType = "woodcutting"
	SkillMining      SkillType = "mining"
	SkillFishing     SkillType = "fishing"
	SkillSmithing    SkillType = "smithing"
	SkillRecycling   SkillType = "recycling"
	SkillCombat      SkillType = "combat"
	SkillCrafting    SkillType = "crafting"
	SkillCooking     SkillType = "cooking"
	SkillAgility     SkillType = "agility"
	SkillThieving    SkillType = "thieving"
)

// SkillNames maps types to display names
var SkillNames = map[SkillType]string{
	SkillWoodcutting: "Woodcutting",
	SkillMining:      "Mining",
	SkillFishing:     "Fishing",
	SkillSmithing:    "Smithing",
	SkillRecycling:   "Recycling",
	SkillCombat:      "Combat",
	SkillCrafting:    "Crafting",
	SkillCooking:     "Cooking",
	SkillAgility:     "Agility",
	SkillThieving:    "Thieving",
}

// Skill represents a single skill with progression
type Skill struct {
	Type           SkillType `json:"type"`
	Level          int       `json:"level"`
	XP             int64     `json:"xp"`
	XPToNext       int64     `json:"xp_to_next"`
	TotalProcessed int64     `json:"total_processed"` // Total items processed/gathered
}

// NewSkill creates a new skill at level 1
func NewSkill(skillType SkillType) *Skill {
	s := &Skill{
		Type:     skillType,
		Level:    1,
		XP:       0,
		XPToNext: 100,
	}
	return s
}

// AddXP adds experience and handles level-ups
func (s *Skill) AddXP(amount int64) []Perk {
	var unlockedPerks []Perk
	s.XP += amount
	s.TotalProcessed += 1

	// Check for level ups
	for s.XP >= s.XPToNext && s.Level < 120 {
		s.XP -= s.XPToNext
		s.Level++
		s.XPToNext = CalculateXPToNext(s.Level)

		// Check for perks at this level
		if perks := GetPerksForLevel(s.Type, s.Level); len(perks) > 0 {
			unlockedPerks = append(unlockedPerks, perks...)
		}
	}

	return unlockedPerks
}

// CalculateXPToNext uses exponential curve
func CalculateXPToNext(level int) int64 {
	if level >= 120 {
		return 0
	}
	base := 100.0
	multiplier := math.Pow(1.15, float64(level-1))
	return int64(base * multiplier)
}

// GetXPForLevel calculates total XP needed to reach a level
func GetXPForLevel(level int) int64 {
	if level <= 1 {
		return 0
	}
	total := int64(0)
	for i := 1; i < level; i++ {
		total += CalculateXPToNext(i)
	}
	return total
}

// Player represents the game state
type Player struct {
	Name            string               `json:"name"`
	CreatedAt       time.Time            `json:"created_at"`
	LastOnline      time.Time            `json:"last_online"`
	TotalPlaytime   time.Duration        `json:"total_playtime"`
	Skills          map[SkillType]*Skill `json:"skills"`
	Inventory       *Inventory           `json:"inventory"`
	Equipment       *Equipment           `json:"equipment"`
	UnlockedPerks   []Perk               `json:"unlocked_perks"`
	CurrentActivity *Activity            `json:"current_activity,omitempty"`
	Gold            int64                `json:"gold"`
	CombatStats     *CombatStats         `json:"combat_stats"`
	Attributes      *CharacterAttributes `json:"attributes"` // Trainable stats
	ActivityLog     *ActivityLog         `json:"activity_log"`

	// Session tracking (not saved)
	SessionStart time.Time `json:"-"`
}

// NewPlayer creates a new player
func NewPlayer(name string) *Player {
	now := time.Now()
	p := &Player{
		Name:          name,
		CreatedAt:     now,
		LastOnline:    now,
		Skills:        make(map[SkillType]*Skill),
		Inventory:     NewInventory(9999), // Unlimited inventory (9999 slots)
		Equipment:     NewEquipment(),
		UnlockedPerks: []Perk{},
		Gold:          0,
		CombatStats:   NewCombatStats(),
		Attributes:    NewCharacterAttributes(),
		ActivityLog:   NewActivityLog(),
		SessionStart:  now,
	}

	// Initialize all skills
	for _, skillType := range []SkillType{
		SkillWoodcutting, SkillMining, SkillFishing,
		SkillSmithing, SkillRecycling, SkillCombat,
		SkillCrafting, SkillCooking, SkillAgility, SkillThieving,
	} {
		p.Skills[skillType] = NewSkill(skillType)
	}

	// Give starting items
	p.Inventory.AddItem(NewItem("bronze_axe", "Bronze Axe", 1))
	p.Inventory.AddItem(NewItem("bronze_pickaxe", "Bronze Pickaxe", 1))

	return p
}

// GetSkill safely retrieves a skill
func (p *Player) GetSkill(skillType SkillType) *Skill {
	if skill, ok := p.Skills[skillType]; ok {
		return skill
	}
	return NewSkill(skillType)
}

// AddXP adds XP to a skill and returns unlocked perks
func (p *Player) AddXP(skillType SkillType, amount int64) []Perk {
	skill := p.GetSkill(skillType)
	unlockedPerks := skill.AddXP(amount)
	p.UnlockedPerks = append(p.UnlockedPerks, unlockedPerks...)
	return unlockedPerks
}

// GetOfflineTime returns time since last online
func (p *Player) GetOfflineTime() time.Duration {
	return time.Since(p.LastOnline)
}

// UpdateLastOnline updates the last online time
func (p *Player) UpdateLastOnline() {
	p.LastOnline = time.Now()
}

// GetTotalLevel returns sum of all skill levels + attribute levels
func (p *Player) GetTotalLevel() int {
	total := 0
	for _, skill := range p.Skills {
		total += skill.Level
	}
	// Add attribute levels
	if p.Attributes != nil {
		total += p.Attributes.GetTotalAttributeLevel()
	}
	return total
}

// GetTotalXP returns sum of all XP
func (p *Player) GetTotalXP() int64 {
	total := int64(0)
	for _, skill := range p.Skills {
		total += skill.XP
	}
	return total
}

// GetSkillMultiplier returns multiplier from perks
func (p *Player) GetSkillMultiplier(skillType SkillType) float64 {
	multiplier := 1.0
	for _, perk := range p.UnlockedPerks {
		if perk.SkillType == skillType {
			switch perk.Effect {
			case "xp_boost":
				multiplier += perk.Value
			case "speed_boost":
				multiplier += perk.Value
			case "double_drop":
				// Handled separately
				multiplier += perk.Value
			}
		}
	}
	return multiplier
}

// CanEquip checks if player meets requirements for equipment
func (p *Player) CanEquip(item *Item) bool {
	for skillType, level := range item.Requirements {
		skill := p.GetSkill(skillType)
		if skill.Level < level {
			return false
		}
	}
	return true
}

// String returns player summary
func (p *Player) String() string {
	return fmt.Sprintf("%s (Total: %d)", p.Name, p.GetTotalLevel())
}
