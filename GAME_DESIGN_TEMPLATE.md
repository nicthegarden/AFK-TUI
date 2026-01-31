# AFK-TUI - Game Design Document Template

## 🎮 Game Overview

| Field | Value |
|-------|-------|
| **Game Name** | AFK-TUI |
| **Genre** | AFK / Idle / Incremental / TUI |
| **Platform** | Terminal / Command Line |
| **Inspiration** | Melvor Idle |
| **Core Loop** | Progress while away, return to upgrade, set new goals |

---

## 📊 Core Game Loop

```
[Start Activity] → [Wait/AFK] → [Collect Resources/XP] → [Unlock Content] → [Upgrade Skills/Equipment] → Repeat
```

### Key Mechanics

| Mechanic | Description | Implementation Notes |
|----------|-------------|---------------------|
| Offline Progression | Progress continues when game is closed | Save timestamp on exit, calculate delta on startup |
| Tick System | Game updates every X seconds | Configurable tick rate (default: 1s) |
| Resource Generation | Passive income based on active skill | Skill level × equipment modifiers |
| AFK Bonuses | Rewards for longer idle periods | Diminishing returns after threshold |
| Easy modable with module that with a start menu with a logo splashscreen you "AFK" but also manage the saves ,

---

## 🎯 Skills System

### Skill Categories

#### 1. Gathering Skills
| Skill | Description | Primary Resource | Idle Rate |
|-------|-------------|------------------|-----------|
| Woodcutting | Chop trees for logs | Logs | 1 log / 3s at level 1 |
| Mining | Extract ores and gems | Ores | 1 ore / 5s at level 1 |
| Fishing | Catch fish | Fish | 1 fish / 4s at level 1 |
| Farming | Grow crops | Crops | 1 crop / 10s at level 1 |
| Foraging | Gather herbs and materials | Herbs | 1 herb / 6s at level 1 |

#### 2. Processing Skills
| Skill | Description | Input → Output |
|-------|-------------|----------------|
| Smithing | Process ores into bars | Ore → Metal Bars |
| Cooking | Prepare food | Raw Food → Cooked Food |
| Crafting | Create equipment | Materials → Items |
| Fletching | Make ammunition and bows | Wood → Arrows/Bows |
| Alchemy | Create potions | Herbs → Potions |

#### 3. Combat Skills
| Skill | Description | Mechanics |
|-------|-------------|-----------|
| Attack | Melee accuracy | Affects hit chance |
| Strength | Melee damage | Affects damage output |
| Defence | Damage reduction | Reduces incoming damage |
| Hitpoints | Health pool | Determines survivability |
| Ranged | Ranged combat | Distance attacks |
| Magic | Spell casting | Resource-based combat |
| Slayer | Boss hunting | Special monsters/tasks |

#### 4. Utility Skills
| Skill | Description | Unlocks |
|-------|-------------|---------|
| Agility | Movement speed | Shortcut unlocks |
| Thieving | Steal resources | Passive gold generation |
| Construction | Build structures | Storage expansion |
| Summoning | Summon familiars | Combat helpers/gathering boost |

### Skill Progression Formula
```
XP to Next Level = floor(100 × 1.1^(current_level - 1))
Max Level: 120 (virtual levels possible)
```

---

## 🎒 Inventory & Storage System

### Storage Types

| Type | Capacity | Unlocked By |
|------|----------|-------------|
| Inventory | 28 slots | Starting |
| Bank | 500+ slots | Starting |
| Material Storage | Per-skill storage | Skill milestones |
| Equipment Chest | 100 slots | Level 10 Construction |

### Item Categories

| Category | Examples |
|----------|----------|
| Resources | Logs, Ores, Fish, Crops |
| Processed | Bars, Cooked Food, Potions |
| Equipment | Weapons, Armor, Tools |
| Consumables | Food, Potions, Buffs |
| Quest Items | Keys, Special drops |

---

## ⚔️ Combat System

### Combat Mechanics

| Element | Description |
|---------|-------------|
| Combat Triangle | Melee > Ranged > Magic > Melee |
| Auto-Combat | Automatic attacks every 2.4s |
| Food Healing | Auto-eat when HP < threshold |
| Special Attacks | Weapon abilities |
| Drop Tables | Each monster has unique drops |

### Monster Tiers

| Tier | Level Range | Difficulty | Drop Quality |
|------|-------------|------------|--------------|
| 1 | 1-20 | Easy | Basic |
| 2 | 20-50 | Medium | Standard |
| 3 | 50-80 | Hard | Rare |
| 4 | 80-100 | Expert | Epic |
| 5 | 100-120 | Master | Legendary |

### Dungeons
| Dungeon | Combat Level | Boss | Unique Reward |
|---------|--------------|------|---------------|
| Forest Cave | 20 | Giant Spider | Spider Silk Armor |
| Volcanic Depths | 50 | Fire Elemental | Ember Weapons |
| Frozen Citadel | 80 | Ice Queen | Frostborn Set |
| Void Realm | 110 | Void Lord | Abyssal Gear |

---

## 🛠️ Equipment System

### Equipment Slots

| Slot | Description | Effect |
|------|-------------|--------|
| Head | Helmets/Hats | Defence bonus |
| Body | Chest armor | Major defence |
| Legs | Leg armor | Defence bonus |
| Feet | Boots | Speed/defence |
| Hands | Gloves | Accuracy bonus |
| Weapon | Main hand | Attack power |
| Off-hand | Shield/Off-weapon | Defence/offence |
| Cape | Cloaks | Skill bonuses |
| Ring | Jewelry | Various bonuses |
| Amulet | Neckwear | Major bonuses |
| Quiver/Ammo | Ranged ammo | Required for ranged |

### Equipment Tiers

| Tier | Material | Level Requirement |
|------|----------|-------------------|
| 1 | Bronze | 1 |
| 2 | Iron | 10 |
| 3 | Steel | 20 |
| 4 | Mithril | 30 |
| 5 | Adamant | 40 |
| 6 | Rune | 50 |
| 7 | Dragon | 60 |
| 8 | Barrows | 70 |
| 9 | Crystal | 80 |
| 10 | Primal | 90 |

---

## 🏆 Achievement System

### Milestones

| Category | Example Milestones |
|----------|-------------------|
| Skill Mastery | Reach level 99, Reach 200M XP |
| Collection | Obtain every item in a category |
| Combat | Defeat X monsters, Complete dungeon |
| Progression | Unlock all skills, Max all skills |
| Challenges | Hardcore mode, Speed runs |

### Rewards
- Titles (displayed next to name)
- Cosmetic overrides
- Permanent bonuses
- Special unlocks

---

## 💰 Economy System

### Currency Types

| Currency | Earned By | Used For |
|----------|-----------|----------|
| Gold | Combat, Thieving, Selling items | Equipment, Consumables |
| Skill Tokens | Leveling skills | Skill-specific unlocks |
| Achievement Points | Completing milestones | Permanent upgrades |
| Slayer Coins | Slayer tasks | Slayer rewards |

### Market/Shop
- General Store (buy/sell basics)
- Skill Shops (skill-specific items)
- Slayer Shop (combat rewards)
- Special Shop (achievement unlocks)

---

## 🖥️ TUI Interface Design

### Screen Layout (80×24 Terminal)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ AFK-TUI v1.0                         Gold: 123,456                    [?]Help │
├──────────────────────────┬──────────────────────────────────────────────────┤
│                          │                                                  │
│  [NAVIGATION]            │              [MAIN CONTENT AREA]                 │
│                          │                                                  │
│  ┌─ Skills ─────────┐   │                                                  │
│  │ • Woodcutting    │   │                                                  │
│  │ • Mining         │   │                                                  │
│  │ • Fishing        │   │                                                  │
│  │ • Combat         │   │                                                  │
│  │ • ...            │   │                                                  │
│  └──────────────────┘   │                                                  │
│                          │                                                  │
│  ┌─ Character ─────┐    │                                                  │
│  │ Level: 45        │    │                                                  │
│  │ HP: 100/100      │    │                                                  │
│  └──────────────────┘    │                                                  │
│                          │                                                  │
├──────────────────────────┴──────────────────────────────────────────────────┤
│ [ACTIVE SKILL: Woodcutting] Progress: [████████░░░░░░░░░░] 45% | 1.2s/tick   │
│ Currently: Cutting Oak Tree (Level 15) | Logs: 1,234 | XP: 45,678 (+12/tick) │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Key Bindings

| Key | Action |
|-----|--------|
| `↑↓` | Navigate menus |
| `Enter` | Select / Confirm |
| `Esc/q` | Back / Quit |
| `Space` | Pause/Resume current activity |
| `i` | Open Inventory |
| `b` | Open Bank |
| `s` | Open Skills |
| `c` | Open Combat |
| `1-9` | Quick actions |
| `?` | Help menu |
| `Ctrl+S` | Save game |
| `Ctrl+L` | Load game |

### Views/Windows

1. **Dashboard** - Overview of all skills, current activity, notifications
2. **Skill View** - Detail for selected skill, actions available, XP progress
3. **Inventory** - Items, equipment, management
4. **Bank** - Storage, sorting, tabs
5. **Combat** - Monster selection, combat stats, loot log
6. **Equipment** - Gear management, comparisons
7. **Shop** - Purchase/sell interface
8. **Settings** - Game options, save/load

---

## 💾 Save System

### Save Data Structure (Go Structs)

```go
package models

import (
	"time"
)

// Player represents the main player data
type Player struct {
	Name           string                 `json:"name"`
	CreatedAt      time.Time              `json:"created_at"`
	LastOnline     time.Time              `json:"last_online"`
	TotalPlaytime  time.Duration          `json:"total_playtime"`
	Skills         map[SkillType]*Skill   `json:"skills"`
	Inventory      *Inventory             `json:"inventory"`
	Bank           *Bank                  `json:"bank"`
	Equipment      Equipment              `json:"equipment"`
	Settings       Settings               `json:"settings"`
	Unlocks        []string               `json:"unlocks"`
	CurrentActivity *ActivityState        `json:"current_activity,omitempty"`
}

// Skill represents a single skill
type Skill struct {
	Type      SkillType `json:"type"`
	Level     int       `json:"level"`
	XP        int64     `json:"xp"`
	XPToNext  int64     `json:"xp_to_next"`
}

// Item represents an item in inventory
type Item struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Equipped  bool   `json:"equipped,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Inventory holds items
type Inventory struct {
	Items     []*Item `json:"items"`
	MaxSlots  int     `json:"max_slots"`
}

// Equipment slots
type Equipment struct {
	Head     string `json:"head,omitempty"`
	Body     string `json:"body,omitempty"`
	Legs     string `json:"legs,omitempty"`
	Feet     string `json:"feet,omitempty"`
	Hands    string `json:"hands,omitempty"`
	Weapon   string `json:"weapon,omitempty"`
	OffHand  string `json:"off_hand,omitempty"`
	Cape     string `json:"cape,omitempty"`
	Ring     string `json:"ring,omitempty"`
	Amulet   string `json:"amulet,omitempty"`
	Ammo     string `json:"ammo,omitempty"`
}

// Settings for the player
type Settings struct {
	TickRate          time.Duration `json:"tick_rate"`
	AutoSave          bool          `json:"auto_save"`
	AutoSaveInterval  time.Duration `json:"auto_save_interval"`
	AutoEatThreshold  float64       `json:"auto_eat_threshold"`
	MaxOfflineTime    time.Duration `json:"max_offline_time"`
}

// ActivityState tracks current activity for offline calc
type ActivityState struct {
	ActivityID string    `json:"activity_id"`
	SkillType  SkillType `json:"skill_type"`
	StartedAt  time.Time `json:"started_at"`
}

// SkillType enum
type SkillType string

const (
	SkillWoodcutting SkillType = "woodcutting"
	SkillMining      SkillType = "mining"
	SkillFishing     SkillType = "fishing"
	SkillCombat      SkillType = "combat"
	// ... more skills
)
```

### JSON Save Format Example

```json
{
  "name": "PlayerOne",
  "created_at": "2026-01-30T10:00:00Z",
  "last_online": "2026-01-31T08:00:00Z",
  "total_playtime": 86400000000000,
  "skills": {
    "woodcutting": {
      "type": "woodcutting",
      "level": 45,
      "xp": 123456,
      "xp_to_next": 789
    },
    "mining": {
      "type": "mining",
      "level": 30,
      "xp": 65432,
      "xp_to_next": 123
    }
  },
  "inventory": {
    "items": [
      {"id": "oak_logs", "name": "Oak Logs", "quantity": 1234},
      {"id": "bronze_axe", "name": "Bronze Axe", "quantity": 1, "equipped": true}
    ],
    "max_slots": 28
  },
  "current_activity": {
    "activity_id": "cut_oak",
    "skill_type": "woodcutting",
    "started_at": "2026-01-31T07:00:00Z"
  },
  "unlocks": [
    "woodcutting_unlocked",
    "mining_unlocked",
    "oak_tree_unlocked"
  ]
}
```

### Offline Calculation

```
offline_time = current_time - last_online
max_offline = 24 hours (configurable)

for each active skill:
    ticks = min(offline_time / tick_rate, max_offline / tick_rate)
    resources = ticks × rate
    xp = ticks × xp_per_tick
    
apply_limits(resources, storage_capacity)
```

---

## 🎨 Visual Design

### Color Palette

| Element | Color Code | Description |
|---------|-----------|-------------|
| Background | `black` / `#000000` | Terminal background |
| Primary Text | `white` / `#ffffff` | Standard text |
| Secondary | `gray` / `#808080` | Dimmed text |
| Success | `green` / `#00ff00` | Gains, positive |
| Warning | `yellow` / `#ffff00` | Caution, medium level |
| Danger | `red` / `#ff0000` | Damage, error |
| Info | `cyan` / `#00ffff` | Tips, highlights |
| Rare | `magenta` / `#ff00ff` | Rare items |
| Epic | `bright_yellow` | Epic items |
| Legendary | `bright_cyan` | Legendary items |

### Progress Bars

| Type | Characters |
|------|-----------|
| XP Bar | `█` filled, `░` empty |
| Health | `♥` filled, `♡` empty (colored) |
| Activity | `▓` filled, `▒` empty |

### Icons/Unicode

| Meaning | Unicode | Display |
|---------|---------|---------|
| Level Up | ↗ | ↗ |
| Achievement | ★ | ★ |
| Locked | 🔒 | 🔒 or [L] |
| Equipped | ⚔ | ⚔ or [E] |
| New | ✨ | ✨ or [!] |
| Skill | ⭐ | ⭐ |
| Combat | ⚔ | ⚔ |
| Resource | 📦 | 📦 or [R] |

---

## 🔧 Technical Architecture (Go Implementation)

### Project Structure

```
afk-tui/
├── cmd/
│   └── afk-tui/
│       └── main.go
├── internal/
│   ├── engine/
│   │   ├── game.go          # Main game loop
│   │   └── tick.go          # Tick management
│   ├── models/
│   │   ├── player.go        # Player struct & methods
│   │   ├── skill.go         # Skill definitions
│   │   ├── item.go          # Item definitions
│   │   ├── activity.go      # Activity system
│   │   └── combat.go        # Combat mechanics
│   ├── data/
│   │   ├── save.go          # Save/load system
│   │   └── offline.go       # Offline calculation
│   ├── ui/
│   │   ├── app.go           # Bubble Tea app
│   │   ├── views/
│   │   │   ├── dashboard.go
│   │   │   ├── skills.go
│   │   │   ├── inventory.go
│   │   │   └── combat.go
│   │   └── components/
│   │       ├── progressbar.go
│   │       └── list.go
│   └── utils/
│       └── math.go          # XP calculations, etc.
├── assets/
│   └── data/
│       ├── items.json
│       ├── skills.json
│       └── monsters.json
├── go.mod
└── README.md
```

### Game Loop (Go)

```go
package engine

import (
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

type GameEngine struct {
	TickRate        time.Duration
	Running         bool
	CurrentActivity *models.Activity
	Player          *models.Player
	SaveManager     *data.SaveManager
	LastTick        time.Time
}

func NewGameEngine(player *models.Player, saveMgr *data.SaveManager) *GameEngine {
	return &GameEngine{
		TickRate:    1 * time.Second,
		Running:     true,
		Player:      player,
		SaveManager: saveMgr,
		LastTick:    time.Now(),
	}
}

func (e *GameEngine) Init() tea.Cmd {
	// Process offline progress on startup
	offlineTime := time.Since(e.Player.LastOnline)
	if offlineTime > 0 {
		e.processOfflineProgress(offlineTime)
	}
	
	return e.tickCmd()
}

func (e *GameEngine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		// Process current activity
		if e.CurrentActivity != nil {
			e.processTick(e.CurrentActivity)
		}
		e.LastTick = time.Now()
		return e, e.tickCmd()
		
	case tea.KeyMsg:
		return e.handleInput(msg)
	}
	return e, nil
}

func (e *GameEngine) tickCmd() tea.Cmd {
	return tea.Tick(e.TickRate, func(t time.Time) tea.Msg {
		return TickMsg{Time: t}
	})
}

func (e *GameEngine) processTick(activity *models.Activity) {
	// Calculate gains based on rates
	resources := activity.CalculateYield(e.Player)
	xp := activity.CalculateXP(e.Player)
	
	// Apply to player
	e.Player.AddXP(activity.Skill, xp)
	e.Player.Inventory.AddItems(resources)
	
	// Check triggers (drops, unlocks, etc.)
	e.checkTriggers(activity)
	
	// Auto-save periodically
	if e.Player.ShouldAutoSave() {
		e.SaveManager.Save(e.Player)
	}
}

func (e *GameEngine) processOfflineProgress(offlineTime time.Duration) {
	maxOffline := 24 * time.Hour // Configurable
	if offlineTime > maxOffline {
		offlineTime = maxOffline
	}
	
	ticks := int(offlineTime / e.TickRate)
	
	if e.CurrentActivity != nil {
		activity := e.CurrentActivity
		for i := 0; i < ticks; i++ {
			resources := activity.CalculateYield(e.Player)
			xp := activity.CalculateXP(e.Player)
			e.Player.AddXP(activity.Skill, xp)
			e.Player.Inventory.AddItems(resources)
		}
	}
	
	// Apply storage limits
	e.Player.Inventory.ApplyLimits()
}

type TickMsg struct {
	Time time.Time
}
```

### Core Components (Go)

| Package | Type | Responsibility |
|---------|------|---------------|
| `engine` | `GameEngine` | Main loop, tick management |
| `models` | `Player` | Stats, inventory, progress |
| `models` | `SkillManager` | Skill definitions, XP curves |
| `models` | `Activity` | Skill actions, rates, requirements |
| `models` | `ItemManager` | Item definitions, properties |
| `models` | `CombatEngine` | Combat calculations |
| `data` | `SaveManager` | JSON serialization, offline calc |
| `ui` | `Bubble Tea Model` | Terminal UI with Bubble Tea |
| `ui` | `KeyMap` | Keyboard shortcuts |

---

## 📋 Content Checklist

### MVP (Minimum Viable Product)

- [ ] 5 Gathering Skills (Woodcutting, Mining, Fishing, Farming, Foraging)
- [ ] 3 Processing Skills (Smithing, Cooking, Crafting)
- [ ] Basic Combat (Attack/Defence/Hitpoints)
- [ ] 10+ Monsters across 2 dungeons
- [ ] 50+ Items
- [ ] Basic TUI (skills, inventory, combat)
- [ ] Save/Load system
- [ ] Offline progression

### Phase 2

- [ ] All Melvor-inspired skills
- [ ] Advanced combat (Ranged, Magic)
- [ ] 5+ Dungeons with bosses
- [ ] 100+ Items
- [ ] Equipment system
- [ ] Achievement system
- [ ] Shop/Economy

### Phase 3

- [ ] Slayer system
- [ ] Farming expansion
- [ ] Construction
- [ ] Pet system
- [ ] Leaderboards (optional)
- [ ] Mod support (optional)

---

## 🎯 Balancing Guidelines

### XP Curve
- Early levels: Fast (1-10 minutes per level)
- Mid levels: Moderate (30 min - 2 hours per level)
- Late levels: Slow (4+ hours per level)
- 99 should take ~100-200 hours per skill

### Resource Economy
- Lower tier = Abundant, fast gathering
- Higher tier = Rare, slow but valuable
- Processing adds value multiplier (2-5x)
- Consumables should be sustainable

### Combat Difficulty
- Monsters should be defeatable at equal combat level
- Food should heal 10-20% of max HP
- Higher tier monsters = better drops but more food consumption
- Bosses require preparation and gear

---

## 📝 Notes & Ideas

### Potential Features
- Pet system (rare drops that provide bonuses)
- Farming mechanics (time-based crop growth)
- Construction (build rooms for bonuses)
- Agility courses (shortcuts, passive benefits)
- Guilds/Clans (optional multiplayer aspect)
- Minigames (fishing trawler, etc.)

### Technical Decisions (Go Stack)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Language** | Go | Fast compilation, excellent concurrency, single binary deployment |
| **TUI Framework** | Bubble Tea | Elegant Elm-style architecture, Charm ecosystem (Lipgloss, Bubbles) |
| **Styling** | Lipgloss | CSS-like styling for terminal apps |
| **Save Format** | JSON (pretty-printed) | Human-readable, easy debugging, version control friendly |
| **Config** | YAML/TOML | Player preferences and game settings |
| **Tick System** | Real-time with Bubble Tea | `tea.Tick()` for precise intervals |
| **State Management** | Immutable messages | Follows Bubble Tea architecture |
| **Concurrency** | Goroutines | For background saves, offline calculations |

### Go Dependencies (go.mod)

```go
module github.com/yourusername/afk-tui

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/charmbracelet/bubbles v0.18.0
    github.com/fsnotify/fsnotify v1.7.0  // For config hot-reload
    gopkg.in/yaml.v3 v3.1.0               // Config files
)
```

---

## 🚀 Getting Started (Go)

### Quick Start Checklist

#### 1. Project Setup
```bash
# Create project
go mod init github.com/yourusername/afk-tui

# Install dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbles
go get gopkg.in/yaml.v3

# Create directory structure
mkdir -p cmd/afk-tui internal/{engine,models,data,ui/views,ui/components,utils} assets/data
```

#### 2. Implementation Steps

1. [ ] **Create data models** (`internal/models/`)
   - Define `Player`, `Skill`, `Item`, `Activity` structs
   - Add JSON tags for serialization

2. [ ] **Implement save system** (`internal/data/`)
   - JSON marshal/unmarshal functions
   - File I/O with error handling
   - Offline time calculation

3. [ ] **Build game engine** (`internal/engine/`)
   - Initialize Bubble Tea program
   - Implement tick loop with `tea.Tick`
   - Handle offline progress on startup

4. [ ] **Create TUI views** (`internal/ui/`)
   - Dashboard view (overview)
   - Skill view (list + detail)
   - Inventory view (grid/list)
   - Combat view (monster selection)

5. [ ] **Add content**
   - Define 3-5 skills in JSON
   - Create 10+ items
   - Set up 2 dungeons with monsters

6. [ ] **Implement game mechanics**
   - XP calculation formula
   - Resource generation rates
   - Combat calculations

7. [ ] **Polish & test**
   - Add colors with Lipgloss
   - Test offline progression
   - Balance rates

### Main Entry Point Example

```go
// cmd/afk-tui/main.go
package main

import (
	"fmt"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/afk-tui/internal/engine"
	"github.com/yourusername/afk-tui/internal/data"
)

func main() {
	// Load or create player
	saveManager := data.NewSaveManager("./save.json")
	player, err := saveManager.Load()
	if err != nil {
		player = data.NewPlayer("Adventurer")
	}
	
	// Create game engine
	game := engine.NewGameEngine(player, saveManager)
	
	// Start Bubble Tea program
	p := tea.NewProgram(game, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	
	// Save on exit
	saveManager.Save(player)
}
```

### Useful Go Resources

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Effective Go](https://go.dev/doc/effective_go)

---

*Template Version: 1.1 (Go Edition)*
*Last Updated: 2026-01-30*
*For: AFK-TUI Game Development*
