# AFK-TUI - Terminal Idle Game

A grindy, feature-rich idle game inspired by Melvor Idle, built for the terminal with Go and Bubble Tea.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)

## Features

### Core Gameplay
- **AFK Progression**: Game continues while you're away (up to 24 hours)
- **Tick-Based System**: Real-time updates every second
- **Skill System**: 10 skills to level up to 120
- **Perk System**: Unlock permanent bonuses as you level up
- **Equipment System**: Equip tools and armor for bonuses

### Skills

#### Gathering
- **Woodcutting**: Chop trees for logs (Logs → Oak → Willow → Maple → Yew → Magic)
- **Mining**: Extract ores (Copper/Tin → Iron → Coal → Silver → Gold → Mithril → Adamantite → Runite)

#### Processing
- **Smithing**: Smelt ores into bars, craft tools and weapons
- **Recycling**: Break down items into materials for crafting
- **Crafting**: Create items from materials (coming soon)
- **Cooking**: Prepare food (coming soon)

#### Combat (Basic)
- **Combat**: Fight monsters, level up combat stats
- Equipment with attack, strength, and defence bonuses

#### Utility
- **Agility**: Movement and efficiency bonuses (coming soon)
- **Thieving**: Steal gold passively (coming soon)

### Recycling & Crafting System
The unique recycling system lets you:
1. **Gather** raw materials (logs, ores)
2. **Recycle** items into fragments (wood fragments, metal fragments)
3. **Craft** better equipment using bars + fragments
4. **Upgrade** your tools for faster gathering

Example:
- Chop logs → Recycle into wood fragments
- Mine copper + tin → Smelt bronze bars
- Smith: 1 Bronze Bar + 5 Wood Fragments = Bronze Axe (Tool Power +1)

### Perks System
Unlock permanent bonuses at milestone levels:
- **Level 5**: 10% speed boost
- **Level 10**: 15% XP boost
- **Level 20**: Double drop chance (5%)
- **Level 35**: 20% speed boost
- **Level 50**: 25% XP boost
- **Level 80**: Triple drop chance (10%)

### Equipment System
10 equipment slots:
- Head, Body, Legs, Feet, Hands
- Weapon, Off-hand
- Cape, Ring, Amulet
- Ammo

Tools provide **Tool Power** which increases gathering speed by 5% per power.

### Progression
- **XP Curve**: Exponential curve requiring ~100-200 hours for level 99
- **120 Levels**: Max level with virtual levels beyond
- **Offline Calculation**: Precise tracking of time away
- **Auto-Save**: Saves every minute and on exit

## Installation

### From Source
```bash
# Clone the repository
git clone <repository>
cd afk-tui

# Build
go build -o afk-tui ./cmd/afk-tui

# Run
./afk-tui
```

### Requirements
- Go 1.21 or later
- Terminal with UTF-8 support
- Minimum 80x24 terminal size

## Controls

### Global Shortcuts
| Key | Action |
|-----|--------|
| `?` or `h` | Toggle help |
| `d` | Dashboard (main view) |
| `s` | Skills view |
| `i` | Inventory |
| `b` | Bank |
| `e` | Equipment |
| `q` or `Ctrl+C` | Save and quit |
| `Space` | Pause/Resume activity |
| `Ctrl+S` | Manual save |

### Quick Actions
| Key | Action |
|-----|--------|
| `1` | Chop Logs (Woodcutting) |
| `2` | Mine Copper (Mining) |
| `3` | Smelt Bronze (Smithing) |
| `4` | Recycle Logs (Recycling) |

### Navigation
| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate menus |
| `Enter` | Select/Confirm |
| `Esc` | Go back |

## Gameplay Guide

### Getting Started

1. **Start the game**: Run `./afk-tui`
2. **Begin gathering**: Press `1` to start chopping logs
3. **Level up**: Watch your woodcutting level increase as you chop
4. **Unlock perks**: At level 5, you'll get your first perk (+10% speed!)

### Early Game Strategy

1. **Chop logs** until level 15 (unlocks oak trees)
2. **Mine copper and tin** to make bronze bars
3. **Recycle logs** for wood fragments
4. **Smith a bronze axe** (needs 1 bar + 5 fragments)
5. **Equip the axe** for faster woodcutting

### Mid Game (Levels 20-50)

1. **Get iron tools** for better efficiency
2. **Level multiple skills** - they all provide perks!
3. **Use the bank** to store excess materials
4. **Craft steel equipment** at level 35

### Late Game (Levels 50+)

1. **High-tier ores**: Mithril, Adamantite, Runite
2. **Optimize perks**: Stack XP and speed bonuses
3. **Maximize efficiency**: Use best tools and equipment
4. **Reach level 120**: The ultimate grind!

### Tips

- **Multi-skill**: Train multiple skills, they all help each other
- **Recycle everything**: Don't waste old items - recycle them!
- **Check perks**: View perks in the skills menu (`s` then navigate)
- **Offline gains**: The game saves automatically, close and come back later
- **Tool upgrades**: Priority #1 is better tools for faster XP

## Save System

The game automatically saves:
- Every 60 seconds while playing
- When you quit (press `q`)
- On manual save (`Ctrl+S`)

Save file: `afk-tui-save.json` (human-readable JSON)

### Offline Progress
When you return:
- Up to 24 hours of progress is calculated
- XP and items are awarded automatically
- New perks are unlocked immediately
- Activity continues from where you left off

## File Structure

```
afk-tui/
├── cmd/afk-tui/
│   └── main.go              # Entry point
├── internal/
│   ├── engine/
│   │   └── game.go          # Game loop & logic
│   ├── models/
│   │   ├── player.go        # Player data
│   │   ├── skill.go         # Skill system
│   │   ├── item.go          # Item database
│   │   ├── activity.go      # Activities & recipes
│   │   ├── equipment.go     # Equipment system
│   │   ├── inventory.go     # Inventory & bank
│   │   └── perk.go          # Perks system
│   ├── data/
│   │   └── save.go          # Save/load & offline calc
│   └── ui/
│       └── view.go          # TUI rendering
├── afk-tui-save.json        # Your save file (created on first run)
└── README.md
```

## Technical Details

### XP Formula
```
XP to Next Level = floor(100 × 1.15^(level-1))
```

### Speed Calculation
```
Effective Speed = Base Speed × (1 + Tool Power × 0.05) × (1 + Perk Bonuses) × (1 + Level Bonus)
```

### Made With
- **Go 1.21** - Language
- **Bubble Tea** - TUI framework
- **Lipgloss** - Terminal styling
- **Love** - For idle games

## Roadmap

### MVP (Done)
- [x] Core game loop
- [x] 5+ skills with activities
- [x] Equipment system
- [x] Smithing with recycled materials
- [x] Perks system
- [x] Save/load with offline progress
- [x] TUI interface

### Phase 2 (Planned)
- [ ] Combat system with monsters
- [ ] Dungeons and bosses
- [ ] Cooking skill
- [ ] Crafting expansion
- [ ] Farming skill
- [ ] Quest system

### Phase 3 (Future)
- [ ] Slayer tasks
- [ ] Pet system
- [ ] Minigames
- [ ] Construction
- [ ] Guilds/Clans

## License

MIT License - Feel free to use, modify, and distribute!

## Credits

Inspired by [Melvor Idle](https://melvoridle.com/) - An amazing idle RPG!

---

**Happy Grinding!** Remember: The numbers go up, and that's what matters. 
