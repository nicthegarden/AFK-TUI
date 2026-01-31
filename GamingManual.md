# AFK-TUI Gaming Manual

A comprehensive guide to the AFK-TUI terminal idle game.

## Table of Contents
1. [Getting Started](#getting-started)
2. [Controls & Keybindings](#controls--keybindings)
3. [Core Game Mechanics](#core-game-mechanics)
4. [Skills & Activities](#skills--activities)
5. [Combat System](#combat-system)
6. [Equipment & Tools](#equipment--tools)
7. [Calculations & Formulas](#calculations--formulas)
8. [Save System](#save-system)
9. [Tips & Strategies](#tips--strategies)

---

## Getting Started

### Installation
```bash
# Clone and build
go build -o afk-tui ./cmd/afk-tui
./afk-tui
```

### First Launch
On first launch, you'll start with:
- Level 1 in all skills
- Bronze Axe and Bronze Pickaxe
- 100 HP
- 0 gold

The game automatically saves every minute and when you quit (press `q` or `Ctrl+C`).

---

## Controls & Keybindings

### Global Shortcuts (Work Everywhere)

| Key | Action |
|-----|--------|
| `d` | Go to Dashboard |
| `s` | Go to Skills screen |
| `i` | Go to Inventory |
| `e` | Go to Equipment |
| `c` | Go to Character Sheet |
| `?` or `h` | Toggle Help screen |
| `Space` | Toggle expanded log view (when grinding) |
| `Ctrl+S` | Save game manually |
| `Ctrl+C` | Save and quit |
| `q` | Quit (from dashboard) or go back |

### Navigation Keys (All Screens)

| Key | Action |
|-----|--------|
| `↑` or `k` | Move up |
| `↓` or `j` | Move down |
| `Enter` | Select/Confirm |
| `Esc` or `q` | Go back/Cancel |
| `1-9` | Quick select by number |

### Screen-Specific Controls

#### Skills Screen (`s`)
- `1-9` - Select skill by number
- `Enter` - View skill categories
- `t` - Quick training shortcut (when Combat selected)
- Letter keys - Select skill by first letter (if not conflicting)

#### Category Selection
- `1-9` - Select category by number
- Letter - Select by first letter (skips global keys: s, d, i, e, c, q)

#### Activity Selection
- `1-9` or letter - Start activity
- Letters that conflict with global keys (s, d, i, e, c, q) are skipped

#### Inventory (`i`)
- `v` - Enter sell/vend mode
- `#` (number keys) - Select item by number
- In sell mode: Enter quantity, `max` for all, `Enter` to confirm, `Y/N` to confirm sale

#### Combat/Slayer (`s` → Combat skill)
- `1-5` - Select monster tier
- Letter - Select monster to fight
- `Esc/q` - Flee from combat or go back

---

## Core Game Mechanics

### Tick System
The game runs on a **1-second tick system**:
- Every second, your current activity progresses
- When progress reaches 100%, you complete one action
- Actions grant XP and items instantly
- XP is added immediately; level-ups happen automatically

### Progress Bars
- Activity progress bar shows how close you are to completing an action
- XP bars show progress to next level
- Progress is saved continuously

### Activity Log
The bottom panel shows the last 3 activities. Press `Space` while grinding to expand to full-screen view showing last 1000 entries.

Log types:
- **XP** (green) - Experience gained
- **Items** (yellow) - Items obtained
- **Level Up** (magenta) - Skill level increases
- **Perks** (cyan) - Unlocked skill bonuses
- **System** (white) - Game events, offline progress

### Offline Progress
When you return after being away:
- Maximum 24 hours of offline progress is calculated
- Your character continues the last activity automatically
- XP, items, and level-ups are applied
- A summary appears in the activity log (not a popup)

---

## Skills & Activities

### Skill List (10 Skills, Levels 1-120)

| Skill | Key | Description | Main Activities |
|-------|-----|-------------|-----------------|
| Woodcutting | `w` | Chop trees for logs | Basic, Quality, Exotic trees |
| Mining | `m` | Mine ores and gems | Basic, Intermediate, Advanced, Gems |
| Smithing | `s` | Smelt bars and craft tools | Basic, Alloys, Precious, Elite, Tools |
| Recycling | `r` | Convert logs to fragments | Wood recycling |
| Combat | `c` | Train attributes and fight | Training, Slayer |
| Crafting | `a` | Create items | (Future expansion) |
| Cooking | `k` | Prepare food | (Future expansion) |
| Agility | `g` | Movement skills | (Future expansion) |
| Thieving | `t` | Stealing/pickpocketing | (Future expansion) |

### Activity Structure
```
Skills (s) → Categories → Activities → Actions
```

Example navigation:
- `s` → `1` (Woodcutting) → `1` (Basic) → `1` (Logs) = Chop Logs
- `s` → `5` (Combat) → `t` = Training
- `s` → `2` (Mining) → `4` (Gems) → `1` (Sapphire) = Mine Sapphire

### Activity Types

**Gathering** (Woodcutting, Mining)
- No inputs required
- Grants XP + resources
- Speed affected by tool power

**Production** (Smithing, Recycling)
- Requires input materials
- Automatically consumes materials on action start
- Stops if you run out of materials
- Grants XP + output items

**Training** (Combat - Strength, Dexterity, Defense)
- Trains character attributes
- Also grants half XP to Combat skill
- No items produced

**Slayer** (Combat)
- Fight monsters for Combat XP + Slayer XP
- Automatic combat with ATB system
- Monster drops items on defeat

---

## Combat System

### Character Attributes (Trainable, Levels 1-120)

| Attribute | Affects |
|-----------|---------|
| **Strength** | Melee damage, max hit |
| **Dexterity** | Attack speed (ATB fill), accuracy, ranged damage |
| **Defense** | Damage reduction, HP |
| **Constitution** | Max HP pool (10 HP per level + base 100) |
| **Intelligence** | Magic damage, spell accuracy |

### Combat Stats (Derived from Attributes)

| Stat | Calculation |
|------|-------------|
| **Max HP** | 100 + (Constitution × 10) |
| **Attack** | Dexterity level |
| **Ranged** | Dexterity level |
| **Magic** | Intelligence level |
| **Combat Level** | Based on all stats (formula below) |

### ATB (Active Time Battle) System
Combat uses an ATB gauge (0-100%):
- Player ATB fills based on Dexterity: `Fill Rate = 5 + (Dexterity × 0.5)` per tick
- Monster ATB fills based on monster speed
- When ATB reaches 100%, that side attacks
- ATB resets to 0 after attacking

### Combat Flow
1. Select Combat skill (`s` → `5`)
2. Choose Training (attributes) or Slayer (fighting)
3. For Slayer: Select tier (1-5), then monster
4. Combat runs automatically
5. First to reduce opponent HP to 0 wins

### Monster Tiers

| Tier | Level Range | Examples |
|------|-------------|----------|
| 1 (Novice) | 1-10 | Chicken, Giant Rat, Spider |
| 2 (Intermediate) | 10-30 | Cow, Skeleton, Zombie |
| 3 (Advanced) | 30-60 | Hill Giant, Moss Giant |
| 4 (Expert) | 60-90 | Green Dragon, Blue Dragon |
| 5 (Master) | 90-120 | Red Dragon, Black Dragon |

### Slayer System
- **Slayer Level**: Separate from Combat level
- **Slayer XP**: Gained from killing monsters
- **Slayer Points**: Currency for slayer rewards (future feature)
- Monsters have specific Slayer XP rewards

### Death/Defeat
If defeated in combat:
- HP restored to full
- Combat ends
- Logged in activity log
- No item loss

---

## Equipment & Tools

### Equipment Slots
Your character has 10 equipment slots:
1. Head
2. Body
3. Legs
4. Feet
5. Hands
6. Weapon
7. Off-hand

### Tools
Tools provide **Tool Power** which increases activity speed:
- **Axe** (Weapon slot) - Speeds up Woodcutting
- **Pickaxe** (Weapon slot) - Speeds up Mining

### Tool Power Formula
```
Speed Bonus = Tool Power × 5%
```

Example: Steel Axe (+3 power) = +15% woodcutting speed

### Tool Progression

| Tool | Power | Smithing Level |
|------|-------|----------------|
| Bronze | +1 | 5 |
| Iron | +2 | 20 |
| Steel | +3 | 35 |
| Mithril | +4 | 45 |
| Adamantite | +6 | TBD |
| Runite | +8 | TBD |

### Equipment Stats
Equipment can provide:
- Tool Power (speed)
- Attack bonuses
- Defense bonuses
- HP bonuses

View current stats in Equipment screen (`e`) or Character Sheet (`c`).

---

## Calculations & Formulas

### XP Calculations

**Skill XP Requirements (Levels 1-120)**
```
XP to Next Level = 100 + (Current Level² × 10)
```

Example:
- Level 1 → 2: 100 + (1 × 1 × 10) = 110 XP
- Level 50 → 51: 100 + (50 × 50 × 10) = 25,100 XP
- Level 99 → 100: 100 + (99 × 99 × 10) = 98,110 XP

**Attribute XP Requirements**
```
XP to Next Level = 100 + (Current Level² × 10)
```
(Same formula as skills)

**Combat XP from Training**
- Full XP to trained attribute
- Half XP to Combat skill

### Combat Calculations

**Damage Formula**
```
Base Damage = (Attack × Strength) / 100
Variance = ±20%
Final Damage = Base Damage × (0.8 to 1.2)
```

**ATB Fill Rate**
```
Player Fill = 5 + (Dexterity × 0.5) per tick
Monster Fill = Monster Speed × 5 per tick
```

**HP Calculation**
```
Max HP = 100 + (Constitution × 10)
```

**Combat Level**
```
Combat Level = (Strength + Dexterity + Defense + Constitution + Intelligence) / 5
```

### Tool Speed Calculation
```
Base Ticks per Action = Activity Base (usually 3-10)
Speed Multiplier = 1 - (Tool Power × 0.05)
Final Ticks = Base Ticks × Speed Multiplier (minimum 1)
```

Example: Chop Logs (base 3 ticks) with Steel Axe (+3 power):
```
3 × (1 - 0.15) = 3 × 0.85 = 2.55 → 3 ticks (rounded up)
```

### Drop Rates
```
Drop Chance = 1 / Drop Rate
```
- Common: 1/1 (100%)
- Uncommon: 1/10 (10%)
- Rare: 1/100 (1%)
- Very Rare: 1/1000 (0.1%)

Always-drop items bypass this check.

### Item Value Calculation
Item values are set in templates and used for:
- Selling to gold
- Trading (future)

Formula:
```
Sell Value = Item Base Value × Quantity
```

---

## Save System

### Save Location
```
./afk-tui-save.json
./afk-tui-save.json.backup (previous save)
```

### Auto-Save Triggers
- Every 60 seconds while playing
- When quitting (q or Ctrl+C)
- Manual save (Ctrl+S)

### Save Data Includes
- All skill levels and XP
- Character attributes and XP
- Inventory (unlimited slots)
- Equipment
- Gold
- Combat stats and Slayer progress
- Activity log (last 1000 entries)
- Current activity (if active)
- Perks unlocked

### Offline Progress Calculation
When loading:
1. Calculate time since last save (max 24 hours)
2. Simulate ticks at 1 tick/second
3. Apply XP, items, and level-ups
4. Log summary to activity log
5. Save new state

### Save Compatibility
Old saves are automatically migrated:
- Missing ActivityLog → Created fresh
- Missing CombatStats → Created fresh  
- Missing Attributes → Created fresh
- Missing fields → Initialized to defaults

---

## Tips & Strategies

### Early Game (Levels 1-20)
1. **Start with Woodcutting** - Chop Logs for easy XP
2. **Mine Copper and Tin** - For making Bronze Bars
3. **Craft Bronze Tools** - Smithing gives XP + useful tools
4. **Recycle excess logs** - Turn Logs into Wood Fragments for tool handles

### Tool Priority
1. Bronze Axe → Better woodcutting speed
2. Bronze Pickaxe → Better mining speed  
3. Steel tools (level 35) → Significant speed boost
4. Mithril tools (level 45) → Best early-game tools

### Skill Training Order
1. **Woodcutting** → Easy starter resource
2. **Mining** → Get ores for smithing
3. **Smithing** → Make better tools + bars sell well
4. **Recycling** (optional) → Convert logs to fragments
5. **Combat** → Train attributes when ready to fight

### Combat Tips
1. Train **Constitution** first - More HP = survive longer
2. Train **Defense** - Reduces damage taken
3. Train **Strength** - Deal more damage
4. Train **Dexterity** - Attack faster
5. Start with Tier 1 monsters (Chickens, Rats)
6. Check monster weakness in selection screen

### Inventory Management
- Inventory is **unlimited** (9999 slots)
- No bank needed
- Sell unwanted items with `v` in inventory
- Check item value before selling
- Keep tool materials for crafting

### Money Making
1. Sell excess ore/bars
2. Craft and sell tools
3. Monster drops (combat)
4. Higher-level resources = more gold

### Efficiency Tips
1. Always have the best tool equipped
2. Set current activity before closing game (offline progress)
3. Use number keys (1-9) for faster navigation
4. Check activity log with Spacebar while grinding
5. Save manually before risky actions (Ctrl+S)

### Hotkey Shortcuts
Quick navigation examples:
- `s` → `1` → `1` → `1` = Chop Logs
- `s` → `2` → `1` → `1` = Mine Copper
- `s` → `3` → `1` → `1` = Smelt Bronze
- `s` → `5` → `t` → `1` = Strength Training
- `i` → `v` → `1` → `Enter` → `Y` = Sell item #1

---

## FAQ

**Q: Can I play multiple characters?**
A: Currently no, single save file only.

**Q: What happens if my inventory is full?**
A: Items are dropped and logged. Inventory is unlimited though, so this rarely happens.

**Q: Do I need to keep the game running?**
A: No! Close anytime. Up to 24 hours of offline progress is calculated when you return.

**Q: Can I die permanently?**
A: No. Death restores your HP and ends combat, but no progress is lost.

**Q: What's the max level?**
A: 120 for all skills and attributes.

**Q: How do I get more gold?**
A: Sell items (press `v` in inventory), monster drops, or high-level smithing.

---

## Version History

### Current Version
- 10 skills (5 fully implemented)
- 30+ monsters
- 90+ items
- 80+ activities
- Combat system with ATB
- Equipment and tool power
- Unlimited inventory
- Offline progression (24h max)
- Activity log with 1000 entries
- Save/Load system

---

*Happy grinding! ⚔️*