package models

import "fmt"

// CombatStyle represents different combat styles
type CombatStyle string

const (
	CombatStyleMelee  CombatStyle = "melee"
	CombatStyleRanged CombatStyle = "ranged"
	CombatStyleMagic  CombatStyle = "magic"
)

// CombatStats now includes main character attributes
type CombatStats struct {
	// Primary Attributes (trained separately)
	Strength     int `json:"strength"`     // Melee damage, max hit
	Dexterity    int `json:"dexterity"`    // Attack speed, accuracy, ranged damage
	Defense      int `json:"defense"`      // Damage reduction, HP
	Constitution int `json:"constitution"` // HP pool, regeneration
	Intelligence int `json:"intelligence"` // Magic damage, spell accuracy

	// Derived Combat Stats
	Hitpoints    int `json:"hitpoints"`
	MaxHitpoints int `json:"max_hitpoints"`
	Attack       int `json:"attack"` // Accuracy
	Ranged       int `json:"ranged"`
	Magic        int `json:"magic"`

	// Slayer Info
	SlayerLevel  int   `json:"slayer_level"`
	SlayerXP     int64 `json:"slayer_xp"`
	SlayerPoints int   `json:"slayer_points"`

	// Current Task
	CurrentTask *SlayerTask `json:"current_task,omitempty"`
}

// SlayerTask represents an active monster slaying task
type SlayerTask struct {
	MonsterID    string `json:"monster_id"`
	MonsterName  string `json:"monster_name"`
	Amount       int    `json:"amount"`
	Killed       int    `json:"killed"`
	Difficulty   string `json:"difficulty"`
	RewardPoints int    `json:"reward_points"`
}

// CharacterAttributes holds base attributes that can be trained
type CharacterAttributes struct {
	Strength     Attribute `json:"strength"`
	Dexterity    Attribute `json:"dexterity"`
	Defense      Attribute `json:"defense"`
	Constitution Attribute `json:"constitution"`
	Intelligence Attribute `json:"intelligence"`
}

// Attribute represents a trainable attribute
type Attribute struct {
	Level    int   `json:"level"`
	XP       int64 `json:"xp"`
	XPToNext int64 `json:"xp_to_next"`
}

// NewAttribute creates a new attribute at level 1
func NewAttribute() Attribute {
	return Attribute{
		Level:    1,
		XP:       0,
		XPToNext: 100,
	}
}

// AddXP adds XP to attribute
func (a *Attribute) AddXP(amount int64) bool {
	a.XP += amount
	leveledUp := false

	for a.XP >= a.XPToNext && a.Level < 120 {
		a.XP -= a.XPToNext
		a.Level++
		a.XPToNext = calculateAttributeXPToNext(a.Level)
		leveledUp = true
	}

	return leveledUp
}

// CalculateXPToNext for attributes
func calculateAttributeXPToNext(level int) int64 {
	if level >= 120 {
		return 0
	}
	return int64(100 + (level * level * 10))
}

// NewCharacterAttributes creates default attributes
func NewCharacterAttributes() *CharacterAttributes {
	return &CharacterAttributes{
		Strength:     NewAttribute(),
		Dexterity:    NewAttribute(),
		Defense:      NewAttribute(),
		Constitution: NewAttribute(),
		Intelligence: NewAttribute(),
	}
}

// NewCombatStats creates default combat stats
func NewCombatStats() *CombatStats {
	return &CombatStats{
		Strength:     1,
		Dexterity:    1,
		Defense:      1,
		Constitution: 1,
		Intelligence: 1,
		Hitpoints:    100,
		MaxHitpoints: 100,
		Attack:       1,
		Ranged:       1,
		Magic:        1,
		SlayerLevel:  1,
		SlayerXP:     0,
		SlayerPoints: 0,
	}
}

// CalculateDerivedStats updates derived stats from base attributes
func (cs *CombatStats) CalculateDerivedStats(attrs *CharacterAttributes) {
	// Max HP = Constitution * 10 + base 100
	cs.MaxHitpoints = 100 + (attrs.Constitution.Level * 10)
	if cs.Hitpoints > cs.MaxHitpoints {
		cs.Hitpoints = cs.MaxHitpoints
	}

	// Attack accuracy based on Dexterity
	cs.Attack = attrs.Dexterity.Level

	// Ranged based on Dexterity
	cs.Ranged = attrs.Dexterity.Level

	// Magic based on Intelligence
	cs.Magic = attrs.Intelligence.Level
}

// GetCombatLevel calculates overall combat level
func (cs *CombatStats) GetCombatLevel() int {
	// Classic combat level formula: Def + Str + Att + (HP/10)
	base := cs.Defense + cs.Strength + cs.Attack
	hpBonus := cs.MaxHitpoints / 10
	return base + hpBonus
}

// GetTotalAttributeLevel returns sum of all attribute levels
func (attrs *CharacterAttributes) GetTotalAttributeLevel() int {
	return attrs.Strength.Level + attrs.Dexterity.Level + attrs.Defense.Level +
		attrs.Constitution.Level + attrs.Intelligence.Level
}

// String returns formatted combat stats
func (cs *CombatStats) String() string {
	return fmt.Sprintf("Combat Lv.%d | HP: %d/%d | STR:%d DEX:%d DEF:%d",
		cs.GetCombatLevel(), cs.Hitpoints, cs.MaxHitpoints,
		cs.Strength, cs.Dexterity, cs.Defense)
}
