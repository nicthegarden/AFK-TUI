package models

import (
	"math/rand"
)

// RollDrop checks if a drop should occur based on drop rate (0.0 to 1.0)
func RollDrop(dropRate float64) bool {
	return rand.Float64() < dropRate
}

// Monster represents an enemy to fight
type Monster struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Level       int           `json:"level"`
	Hitpoints   int           `json:"hitpoints"`
	MaxHP       int           `json:"max_hp"`
	Attack      int           `json:"attack"`
	Defense     int           `json:"defense"`
	Strength    int           `json:"strength"`
	Speed       float64       `json:"speed"` // ATB speed
	Weakness    CombatStyle   `json:"weakness"`
	Resistance  CombatStyle   `json:"resistance"`
	Drops       []MonsterDrop `json:"drops"`
	SlayerXP    int64         `json:"slayer_xp"`
	CombatXP    int64         `json:"combat_xp"`
	Gold        int64         `json:"gold"`
	IsBoss      bool          `json:"is_boss"`
}

// MonsterDrop represents a possible drop from a monster
type MonsterDrop struct {
	ItemID     string  `json:"item_id"`
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	DropRate   float64 `json:"drop_rate"` // 0.0 to 1.0
	AlwaysDrop bool    `json:"always_drop"`
}

// MonsterDatabase contains all monsters
type MonsterDatabase struct {
	monsters map[string]*Monster
}

// GetMonster retrieves a monster by ID
func (db *MonsterDatabase) GetMonster(id string) *Monster {
	if monster, ok := db.monsters[id]; ok {
		return monster
	}
	return nil
}

// GetMonstersByLevelRange returns monsters within a level range
func (db *MonsterDatabase) GetMonstersByLevelRange(minLevel, maxLevel int) []*Monster {
	var result []*Monster
	for _, m := range db.monsters {
		if m.Level >= minLevel && m.Level <= maxLevel && !m.IsBoss {
			result = append(result, m)
		}
	}
	return result
}

// GetBosses returns all boss monsters
func (db *MonsterDatabase) GetBosses() []*Monster {
	var result []*Monster
	for _, m := range db.monsters {
		if m.IsBoss {
			result = append(result, m)
		}
	}
	return result
}

// GetRandomMonster returns a random monster within level range
func (db *MonsterDatabase) GetRandomMonster(minLevel, maxLevel int) *Monster {
	monsters := db.GetMonstersByLevelRange(minLevel, maxLevel)
	if len(monsters) == 0 {
		return nil
	}
	return monsters[rand.Intn(len(monsters))]
}

// Global monster database instance
var Monsters = &MonsterDatabase{
	monsters: map[string]*Monster{
		// Tier 1: Beginner (Level 1-10)
		"chicken": {
			ID: "chicken", Name: "Chicken", Description: "A harmless farm chicken",
			Level: 1, Hitpoints: 10, MaxHP: 10, Attack: 1, Defense: 1, Strength: 1,
			Speed: 1.0, SlayerXP: 10, CombatXP: 10, Gold: 5,
			Drops: []MonsterDrop{
				{ItemID: "raw_chicken", ItemName: "Raw Chicken", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "feather", ItemName: "Feather", Quantity: 5, DropRate: 0.5},
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
			},
		},
		"rat": {
			ID: "rat", Name: "Giant Rat", Description: "An oversized sewer rat",
			Level: 2, Hitpoints: 15, MaxHP: 15, Attack: 2, Defense: 1, Strength: 2,
			Speed: 1.2, SlayerXP: 15, CombatXP: 15, Gold: 8,
			Drops: []MonsterDrop{
				{ItemID: "raw_rat_meat", ItemName: "Raw Rat Meat", Quantity: 1, DropRate: 0.8},
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
			},
		},
		"spider": {
			ID: "spider", Name: "Giant Spider", Description: "A venomous arachnid",
			Level: 4, Hitpoints: 25, MaxHP: 25, Attack: 3, Defense: 2, Strength: 3,
			Speed: 1.5, Weakness: CombatStyleMelee, Resistance: CombatStyleRanged,
			SlayerXP: 25, CombatXP: 20, Gold: 15,
			Drops: []MonsterDrop{
				{ItemID: "spider_silk", ItemName: "Spider Silk", Quantity: 1, DropRate: 0.4},
				{ItemID: "eye_of_newt", ItemName: "Eye of Newt", Quantity: 1, DropRate: 0.2},
			},
		},
		"goblin": {
			ID: "goblin", Name: "Goblin", Description: "A weak but aggressive goblin",
			Level: 5, Hitpoints: 30, MaxHP: 30, Attack: 4, Defense: 3, Strength: 4,
			Speed: 1.0, SlayerXP: 30, CombatXP: 25, Gold: 20,
			Drops: []MonsterDrop{
				{ItemID: "bronze_dagger", ItemName: "Bronze Dagger", Quantity: 1, DropRate: 0.05},
				{ItemID: "coins", ItemName: "Coins", Quantity: 10, DropRate: 0.8},
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
			},
		},

		// Tier 2: Low Level (Level 10-30)
		"cow": {
			ID: "cow", Name: "Cow", Description: "A docile dairy cow",
			Level: 10, Hitpoints: 60, MaxHP: 60, Attack: 5, Defense: 5, Strength: 6,
			Speed: 0.8, SlayerXP: 50, CombatXP: 40, Gold: 30,
			Drops: []MonsterDrop{
				{ItemID: "cow_hide", ItemName: "Cow Hide", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "raw_beef", ItemName: "Raw Beef", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
			},
		},
		"skeleton": {
			ID: "skeleton", Name: "Skeleton", Description: "An undead warrior",
			Level: 15, Hitpoints: 80, MaxHP: 80, Attack: 10, Defense: 8, Strength: 10,
			Speed: 1.0, Weakness: CombatStyleMagic, Resistance: CombatStyleRanged,
			SlayerXP: 70, CombatXP: 60, Gold: 50,
			Drops: []MonsterDrop{
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "bronze_scimitar", ItemName: "Bronze Scimitar", Quantity: 1, DropRate: 0.03},
			},
		},
		"zombie": {
			ID: "zombie", Name: "Zombie", Description: "A shambling corpse",
			Level: 20, Hitpoints: 120, MaxHP: 120, Attack: 12, Defense: 10, Strength: 12,
			Speed: 0.7, Weakness: CombatStyleMagic,
			SlayerXP: 100, CombatXP: 80, Gold: 70,
			Drops: []MonsterDrop{
				{ItemID: "bones", ItemName: "Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "iron_dagger", ItemName: "Iron Dagger", Quantity: 1, DropRate: 0.02},
			},
		},
		"barbarian": {
			ID: "barbarian", Name: "Barbarian", Description: "A fierce tribal warrior",
			Level: 25, Hitpoints: 150, MaxHP: 150, Attack: 15, Defense: 12, Strength: 18,
			Speed: 1.1, SlayerXP: 120, CombatXP: 100, Gold: 100,
			Drops: []MonsterDrop{
				{ItemID: "coins", ItemName: "Coins", Quantity: 50, DropRate: 0.9},
				{ItemID: "bronze_axe", ItemName: "Bronze Axe", Quantity: 1, DropRate: 0.05},
			},
		},

		// Tier 3: Mid Level (Level 30-60)
		"hill_giant": {
			ID: "hill_giant", Name: "Hill Giant", Description: "A towering brute",
			Level: 35, Hitpoints: 250, MaxHP: 250, Attack: 20, Defense: 18, Strength: 25,
			Speed: 0.6, Weakness: CombatStyleRanged,
			SlayerXP: 200, CombatXP: 150, Gold: 200,
			Drops: []MonsterDrop{
				{ItemID: "giant_bones", ItemName: "Giant Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "iron_full_helm", ItemName: "Iron Full Helm", Quantity: 1, DropRate: 0.05},
				{ItemID: "big_bones", ItemName: "Big Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
			},
		},
		"moss_giant": {
			ID: "moss_giant", Name: "Moss Giant", Description: "A giant covered in moss",
			Level: 45, Hitpoints: 300, MaxHP: 300, Attack: 25, Defense: 22, Strength: 30,
			Speed: 0.7, Weakness: CombatStyleMagic,
			SlayerXP: 250, CombatXP: 200, Gold: 300,
			Drops: []MonsterDrop{
				{ItemID: "big_bones", ItemName: "Big Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "steel_sword", ItemName: "Steel Sword", Quantity: 1, DropRate: 0.03},
			},
		},
		"ice_warrior": {
			ID: "ice_warrior", Name: "Ice Warrior", Description: "A frozen knight",
			Level: 55, Hitpoints: 350, MaxHP: 350, Attack: 30, Defense: 30, Strength: 28,
			Speed: 1.0, Weakness: CombatStyleMelee, Resistance: CombatStyleMagic,
			SlayerXP: 300, CombatXP: 250, Gold: 400,
			Drops: []MonsterDrop{
				{ItemID: "mithril_dagger", ItemName: "Mithril Dagger", Quantity: 1, DropRate: 0.02},
				{ItemID: "ice_shard", ItemName: "Ice Shard", Quantity: 1, DropRate: 0.3},
			},
		},

		// Tier 4: High Level (Level 60-90)
		"green_dragon": {
			ID: "green_dragon", Name: "Green Dragon", Description: "A fierce young dragon",
			Level: 65, Hitpoints: 500, MaxHP: 500, Attack: 45, Defense: 40, Strength: 50,
			Speed: 1.2, Weakness: CombatStyleRanged,
			SlayerXP: 500, CombatXP: 400, Gold: 1000,
			Drops: []MonsterDrop{
				{ItemID: "dragon_bones", ItemName: "Dragon Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "green_dragonhide", ItemName: "Green Dragonhide", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "dragon_dagger", ItemName: "Dragon Dagger", Quantity: 1, DropRate: 0.01},
			},
		},
		"blue_dragon": {
			ID: "blue_dragon", Name: "Blue Dragon", Description: "A powerful adult dragon",
			Level: 75, Hitpoints: 650, MaxHP: 650, Attack: 55, Defense: 50, Strength: 60,
			Speed: 1.3, Weakness: CombatStyleMelee, Resistance: CombatStyleMagic,
			SlayerXP: 700, CombatXP: 550, Gold: 1500,
			Drops: []MonsterDrop{
				{ItemID: "dragon_bones", ItemName: "Dragon Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "blue_dragonhide", ItemName: "Blue Dragonhide", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "rune_dagger", ItemName: "Rune Dagger", Quantity: 1, DropRate: 0.005},
			},
		},
		"abyssal_demon": {
			ID: "abyssal_demon", Name: "Abyssal Demon", Description: "A creature from the abyss",
			Level: 85, Hitpoints: 800, MaxHP: 800, Attack: 65, Defense: 60, Strength: 70,
			Speed: 1.5, Weakness: CombatStyleMagic,
			SlayerXP: 1000, CombatXP: 750, Gold: 2000,
			Drops: []MonsterDrop{
				{ItemID: "abyssal_whip", ItemName: "Abyssal Whip", Quantity: 1, DropRate: 0.002},
				{ItemID: "abyssal_head", ItemName: "Abyssal Head", Quantity: 1, DropRate: 0.1},
			},
		},

		// Tier 5: Elite (Level 90+)
		"red_dragon": {
			ID: "red_dragon", Name: "Red Dragon", Description: "An ancient and powerful dragon",
			Level: 95, Hitpoints: 1000, MaxHP: 1000, Attack: 80, Defense: 75, Strength: 85,
			Speed: 1.4, Weakness: CombatStyleRanged,
			SlayerXP: 1500, CombatXP: 1000, Gold: 3000,
			Drops: []MonsterDrop{
				{ItemID: "dragon_bones", ItemName: "Dragon Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "red_dragonhide", ItemName: "Red Dragonhide", Quantity: 2, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "dragonfire_shield", ItemName: "Dragonfire Shield", Quantity: 1, DropRate: 0.001},
			},
		},
		"black_dragon": {
			ID: "black_dragon", Name: "Black Dragon", Description: "The most fearsome dragon",
			Level: 110, Hitpoints: 1200, MaxHP: 1200, Attack: 95, Defense: 90, Strength: 100,
			Speed:    1.5,
			SlayerXP: 2000, CombatXP: 1400, Gold: 5000,
			Drops: []MonsterDrop{
				{ItemID: "dragon_bones", ItemName: "Dragon Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "black_dragonhide", ItemName: "Black Dragonhide", Quantity: 2, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "draconic_visage", ItemName: "Draconic Visage", Quantity: 1, DropRate: 0.0005},
			},
		},

		// Bosses (Raids)
		"giant_mole": {
			ID: "giant_mole", Name: "Giant Mole", Description: "An oversized burrowing creature",
			Level: 40, Hitpoints: 600, MaxHP: 600, Attack: 35, Defense: 30, Strength: 40,
			Speed: 1.0, IsBoss: true,
			SlayerXP: 400, CombatXP: 300, Gold: 1500,
			Drops: []MonsterDrop{
				{ItemID: "mole_claw", ItemName: "Mole Claw", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "mole_skin", ItemName: "Mole Skin", Quantity: 3, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "amulet_of_strength", ItemName: "Amulet of Strength", Quantity: 1, DropRate: 0.05},
			},
		},
		"dagannoth_king": {
			ID: "dagannoth_king", Name: "Dagannoth Rex", Description: "Ruler of the Dagannoth",
			Level: 70, Hitpoints: 1000, MaxHP: 1000, Attack: 60, Defense: 50, Strength: 70,
			Speed: 1.5, IsBoss: true, Weakness: CombatStyleMagic,
			SlayerXP: 800, CombatXP: 600, Gold: 3000,
			Drops: []MonsterDrop{
				{ItemID: "berserker_ring", ItemName: "Berserker Ring", Quantity: 1, DropRate: 0.01},
				{ItemID: "archers_ring", ItemName: "Archers Ring", Quantity: 1, DropRate: 0.01},
			},
		},
		"kalphite_queen": {
			ID: "kalphite_queen", Name: "Kalphite Queen", Description: "Matriarch of the Kalphite hive",
			Level: 80, Hitpoints: 1500, MaxHP: 1500, Attack: 75, Defense: 70, Strength: 80,
			Speed: 1.8, IsBoss: true,
			SlayerXP: 1200, CombatXP: 900, Gold: 5000,
			Drops: []MonsterDrop{
				{ItemID: "kalphite_head", ItemName: "Kalphite Head", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "dragon_chainbody", ItemName: "Dragon Chainbody", Quantity: 1, DropRate: 0.02},
			},
		},
		"godwars_boss": {
			ID: "godwars_boss", Name: "General Graardor", Description: "Leader of the Bandos army",
			Level: 100, Hitpoints: 2000, MaxHP: 2000, Attack: 90, Defense: 85, Strength: 100,
			Speed: 2.0, IsBoss: true,
			SlayerXP: 2000, CombatXP: 1500, Gold: 10000,
			Drops: []MonsterDrop{
				{ItemID: "bandos_tassets", ItemName: "Bandos Tassets", Quantity: 1, DropRate: 0.01},
				{ItemID: "bandos_chestplate", ItemName: "Bandos Chestplate", Quantity: 1, DropRate: 0.01},
				{ItemID: "godsword_shard", ItemName: "Godsword Shard", Quantity: 1, DropRate: 0.05},
			},
		},
		"corporal_beast": {
			ID: "corporal_beast", Name: "Corporeal Beast", Description: "A massive ethereal creature",
			Level: 120, Hitpoints: 3000, MaxHP: 3000, Attack: 110, Defense: 100, Strength: 120,
			Speed: 2.2, IsBoss: true,
			SlayerXP: 3000, CombatXP: 2500, Gold: 20000,
			Drops: []MonsterDrop{
				{ItemID: "corp_bones", ItemName: "Corporeal Bones", Quantity: 1, DropRate: 1.0, AlwaysDrop: true},
				{ItemID: "spectral_sigil", ItemName: "Spectral Sigil", Quantity: 1, DropRate: 0.005},
				{ItemID: "arcane_sigil", ItemName: "Arcane Sigil", Quantity: 1, DropRate: 0.005},
				{ItemID: "elysian_sigil", ItemName: "Elysian Sigil", Quantity: 1, DropRate: 0.002},
			},
		},
	},
}

// CalculateDamage calculates damage based on attacker and defender stats
func CalculateDamage(attackerAttack, attackerStrength int, defenderDefense int, style CombatStyle, weakness, resistance CombatStyle) int {
	// Base damage calculation
	accuracy := float64(attackerAttack) * 2
	evasion := float64(defenderDefense) * 1.5

	// Hit chance
	hitChance := accuracy / (accuracy + evasion)
	if hitChance > 0.95 {
		hitChance = 0.95
	}
	if hitChance < 0.05 {
		hitChance = 0.05
	}

	// Check if hit lands
	if rand.Float64() > hitChance {
		return 0 // Miss
	}

	// Damage calculation
	maxHit := attackerStrength
	minHit := 1
	if attackerStrength > 10 {
		minHit = attackerStrength / 10
	}

	damage := minHit + rand.Intn(maxHit-minHit+1)

	// Apply style modifiers
	if style == weakness {
		damage = int(float64(damage) * 1.5) // 50% bonus for weakness
	} else if style == resistance {
		damage = int(float64(damage) * 0.5) // 50% reduction for resistance
	}

	return damage
}

// GetATBFill calculates how much ATB bar fills per tick
func GetATBFill(speed, dexterity int) float64 {
	// Base speed + dexterity bonus
	totalSpeed := float64(speed) + (float64(dexterity) * 0.02)

	// ATB fills per tick (max 100)
	fillAmount := totalSpeed * 5

	return fillAmount
}

// MonsterTiers for organizing monsters
type MonsterTier struct {
	Name     string
	MinLevel int
	MaxLevel int
	GoldMult float64
	XPBoost  float64
}

// GetMonsterTier returns tier info for a monster level
func GetMonsterTier(level int) MonsterTier {
	switch {
	case level < 10:
		return MonsterTier{Name: "Novice", MinLevel: 1, MaxLevel: 10, GoldMult: 1.0, XPBoost: 1.0}
	case level < 30:
		return MonsterTier{Name: "Intermediate", MinLevel: 10, MaxLevel: 30, GoldMult: 1.5, XPBoost: 1.2}
	case level < 60:
		return MonsterTier{Name: "Advanced", MinLevel: 30, MaxLevel: 60, GoldMult: 2.0, XPBoost: 1.5}
	case level < 90:
		return MonsterTier{Name: "Expert", MinLevel: 60, MaxLevel: 90, GoldMult: 3.0, XPBoost: 2.0}
	case level < 120:
		return MonsterTier{Name: "Master", MinLevel: 90, MaxLevel: 120, GoldMult: 5.0, XPBoost: 3.0}
	default:
		return MonsterTier{Name: "Legendary", MinLevel: 120, MaxLevel: 999, GoldMult: 10.0, XPBoost: 5.0}
	}
}
