# AFK-TUI - Complete Controls Guide

## Quick Reference Card

```
╔════════════════════════════════════════════════════════════════╗
║  GLOBAL SHORTCUTS                                               ║
╠════════════════════════════════════════════════════════════════╣
║  [d] Dashboard    [s] Skills      [i] Inventory  [b] Bank       ║
║  [e] Equipment    [?] Help        [Space] Pause  [q] Quit       ║
╚════════════════════════════════════════════════════════════════╝

╔════════════════════════════════════════════════════════════════╗
║  QUICK START KEYS (From Anywhere)                              ║
╠════════════════════════════════════════════════════════════════╣
║  [1] Chop Logs       - Woodcutting Lv.1                        ║
║  [2] Mine Copper     - Mining Lv.1                             ║
║  [3] Smelt Bronze    - Smithing Lv.1 (needs copper+tin)        ║
║  [4] Recycle Logs    - Recycling Lv.1                          ║
╚════════════════════════════════════════════════════════════════╝
```

## Detailed Controls

### Navigation

| Key | Action | Context |
|-----|--------|---------|
| `↑` / `k` | Move up | Menus, lists |
| `↓` / `j` | Move down | Menus, lists |
| `Enter` | Select / Confirm | All screens |
| `Esc` | Go back / Cancel | All screens |
| `q` | Save and quit | Global |
| `Ctrl+C` | Force quit | Emergency |

### Screen Shortcuts

| Key | Screen | Description |
|-----|--------|-------------|
| `d` | **Dashboard** | Main overview, current activity, skill summary |
| `s` | **Skills** | All skills, levels, XP, perks |
| `i` | **Inventory** | Items, recycle mode, equipment |
| `b` | **Bank** | Unlimited storage |
| `e` | **Equipment** | View and manage gear |
| `?` or `h` | **Help** | This help screen |

### Activity Controls

| Key | Action | Notes |
|-----|--------|-------|
| `Space` | Pause/Resume | Stops current activity |
| `1-4` | Quick start | Start common activities instantly |
| `Ctrl+S` | Manual save | Saves game state |

### Skill-Specific Controls

#### Skills Screen (`s`)
```
Navigate: ↑/↓ or j/k
View activities: a (when skill selected)
Back: Esc or d
```

#### Inventory Screen (`i`)
```
Recycle mode: r
Deposit to bank: b (on selected item)
Back: Esc or d
```

#### Smithing Screen (from Skills → Smithing)
```
[1] Smelt Bronze      - 1 Copper + 1 Tin
[2] Smelt Iron        - 1 Iron Ore
[3] Smelt Steel       - 1 Iron + 2 Coal
[4] Smelt Silver      - 1 Silver Ore (Lv.40)
[5] Smelt Gold        - 1 Gold Ore (Lv.50)
[6] Smith Bronze Axe  - 1 Bar + 5 Wood Fragments
[7] Smith Iron Axe    - 2 Bars + 10 Fragments
[8] Smith Steel Axe   - 2 Bars + 2 Oak Logs
Back: Esc
```

#### Recycling Screen (from Inventory)
```
[1] Recycle Logs      → 1 Wood Fragments
[2] Recycle Oak Logs  → 2 Wood Fragments (Lv.15)
Back: Esc
```

## Activity Quick Start

### Gathering Activities

| Activity | Req. | XP/tick | Output |
|----------|------|---------|--------|
| Chop Logs | Lv.1 | 10 | 1 Logs |
| Chop Oak | Lv.15 | 20 | 1 Oak Logs |
| Chop Willow | Lv.30 | 35 | 1 Willow Logs |
| Mine Copper | Lv.1 | 12 | 1 Copper Ore |
| Mine Tin | Lv.1 | 12 | 1 Tin Ore |
| Mine Iron | Lv.15 | 25 | 1 Iron Ore |
| Mine Coal | Lv.30 | 35 | 1 Coal |

### Processing Activities

| Activity | Req. | XP/tick | Input → Output |
|----------|------|---------|----------------|
| Smelt Bronze | Lv.1 | 15 | 1Cu+1Ti → 1 Bar |
| Smelt Iron | Lv.15 | 30 | 1 Ore → 1 Bar |
| Smelt Steel | Lv.30 | 45 | 1Fe+2Co → 1 Bar |
| Recycle Logs | Lv.1 | 8 | 1 Logs → 1 Fragment |
| Smith Bronze Axe | Lv.5 | 25 | 1Bar+5Frag → Axe |

## Pro Tips

### Efficiency Tips
1. **Start with [1]** - Chop logs immediately
2. **Get a better axe ASAP** - Smith bronze axe at Lv.5
3. **Multi-skill** - Switch between gathering and processing
4. **Recycle everything** - Don't waste old items
5. **Check perks often** - Press `s` to see what's coming

### Keyboard Shortcuts for Speed
```
Start chopping:     1
Check progress:     d
View skills:        s
Make bronze bar:    3
Recycle items:      i → r → 1
Smith new axe:      s → ↓ → a → 6
```

### AFK Strategy
```
1. Start activity (e.g., press 1 to chop)
2. Press q to quit (saves automatically)
3. Come back later!
4. Game calculates up to 24 hours offline progress
```

## Command Combinations

### Early Game Loop (Levels 1-15)
```
1 (chop) → wait → d (check) → 2 (mine) → 3 (smelt) → 4 (recycle) → repeat
```

### Mid Game Loop (Levels 15-35)
```
s (skills) → check oak unlock → 1 (chop oak) → mine iron → 
smelt iron → smith iron axe → equip → faster chopping!
```

### Equipment Check
```
e (equipment) - see current gear and bonuses
```

## Status Bar Icons

When playing, you'll see:

```
AFK-TUI v1.0                    Gold: 1,234    [?]Help
```

**Dashboard shows:**
```
[ACTIVE SKILL: Woodcutting] Progress: [████████░░░░] 67% | 0.8s/tick
Currently: Chopping Oak Tree | Logs: 156 | XP: 12,345 (+18/tick)
```

## Troubleshooting

### "Need X item" error
- You don't have required materials
- Switch to gathering that item first

### "Requires level X" error
- Activity needs higher skill level
- Check skills screen (`s`) to see current level

### Activity paused
- Press `Space` to resume
- Or press a number to start new activity

### Full inventory
- Go to bank (`b`) to deposit
- Or recycle items (`i` → `r`)

## Game Save Location

Save file: `afk-tui-save.json` (in game directory)

The save is plain text JSON - you can even edit it if you're careful!

---

**Remember**: The game saves automatically every minute and when you press `q` or `Ctrl+S`
