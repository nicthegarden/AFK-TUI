# 🎮 Skills Menu Controls

## How to Open
Press **`s`** from anywhere in the game to open the Skills menu.

---

## Skills Menu Layout

```
┌─────────────────────────────────────────────────────────────┐
│  🎮 Skills                                                  │
│                                                             │
│  Navigate: ↑/↓ or j/k | [c] Categories | [Enter] View       │
│                                                             │
│  ► Woodcutting     Lv.15/120 [██████████░░░░░]              │  ← Selected (► cursor)
│     XP: 12,450 / 15,230                                     │
│     Perks: 2/6 unlocked                                     │
│                                                             │
│    Mining          Lv.42/120 [████████████████░░]           │
│                                                             │
│    Smithing        Lv.28/120 [████████████░░░░░]            │
│                                                             │
│    Recycling       Lv.18/120 [█████████░░░░░░░░]            │
│                                                             │
│    Combat          Lv. 5/120 [██░░░░░░░░░░░░░░░]            │
│                                                             │
│    Crafting        Lv.12/120 [██████░░░░░░░░░░░]            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Controls Reference

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `↑` (Up Arrow) | Move cursor up | Navigate to skill above |
| `↓` (Down Arrow) | Move cursor down | Navigate to skill below |
| `k` | Move cursor up | Alternative: vim-style |
| `j` | Move cursor down | Alternative: vim-style |

### Actions

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` or `a` | **View Categories** | Open category selection for highlighted skill |
| `c` | **Categories** | Jump directly to Categories menu |
| `d` | Dashboard | Return to main dashboard |
| `i` | Inventory | Open inventory |
| `b` | Bank | Open bank |
| `e` | Equipment | View equipment |
| `?` or `h` | Help | Show help screen |

### Global

| Key | Action | Description |
|-----|--------|-------------|
| `Space` | Pause/Resume | Toggle current activity |
| `Ctrl+S` | Save | Manual save |
| `q` or `Ctrl+C` | Quit | Save and exit game |

---

## Skill Details (When Selected)

When you navigate to a skill with `↑`/`↓`, you'll see:

### Expanded Information
```
► Woodcutting     Lv.15/120 [██████████░░░░░]  ← Selected skill
   XP: 12,450 / 15,230                           Current XP / XP needed
   Perks: 2/6 unlocked                           How many perks you have
```

### Perk Preview
```
► Mining          Lv.42/120 [████████████████░░]
   XP: 45,230 / 52,100
   Perks: 4/6 unlocked
   
   Perks:                                        ← List of all perks
     ✓ Lv.5: Quick Chop (+10% speed)            ✓ = Unlocked
     ✓ Lv.10: Nature's Wisdom (+15% XP)         ○ = Locked
     ○ Lv.20: Double Logs
     ○ Lv.35: Expert Chopper
     ✓ Lv.50: Forest Mastery
     ○ Lv.80: Triple Logs
```

---

## Quick Workflow Examples

### Check Skill Progress
```bash
s           # Open Skills menu
↓           # Navigate to Mining
↓           # Navigate to Smithing
...         # Check levels and perks
Esc or d    # Return to dashboard
```

### Start Grinding a Skill
```bash
s           # Open Skills menu
↓           # Navigate to desired skill (e.g., Mining)
Enter or a  # Open Categories for Mining
1           # Select "Basic Ores" category
2           # Select "Mine Copper" activity
# You are now mining copper! Check status bar at bottom.
```

### Shortcut: Direct to Categories
```bash
s           # Open Skills menu
↓           # Navigate to Smithing
c           # Jump to Categories (same as Enter)
3           # Select "Precious Metals"
1           # Select "Smelt Silver"
# Now smelting silver!
```

---

## Skill Icons & Colors

| Skill | Icon | Color (Level-Based) |
|-------|------|---------------------|
| Woodcutting | 🪓 | 🟢 Lv.1-29 / 🔵 Lv.30-69 / 🟣 Lv.70-99 / 🟡 Lv.100+ |
| Mining | ⛏️ | Same color coding |
| Smithing | 🔥 | Same color coding |
| Recycling | ♻️ | Same color coding |
| Combat | ⚔️ | Same color coding |
| Crafting | 🛠️ | Same color coding |
| Cooking | 🍳 | Same color coding |
| Agility | 🏃 | Same color coding |
| Thieving | 🎭 | Same color coding |

**Color Guide:**
- 🟢 **Green** (Lv.1-29) - Beginner
- 🔵 **Blue** (Lv.30-69) - Intermediate  
- 🟣 **Purple** (Lv.70-99) - Advanced
- 🟡 **Gold** (Lv.100+) - Legendary

---

## Tips & Tricks

### Efficient Navigation
```bash
# Use j/k for quick navigation (vim-style)
s           # Open Skills
j j j       # Down 3 times to Smithing
k           # Up 1 to Recycling
Enter       # View categories
```

### Check Multiple Skills
```bash
s           # Open Skills
↓           # Woodcutting
↓           # Mining  
↓           # Smithing
# Compare levels quickly
```

### Perk Planning
```bash
s           # Open Skills
↓           # Navigate to skill you want to train
# Read the locked perks (○) to see what's coming
# Plan your grinding to reach those levels!
```

---

## Status Bar (Always Visible)

Even while browsing the Skills menu, the **status bar** at the bottom shows:

```
 ⛏️ Mine Copper [████████████████████░░░░░░░░░░░░░░░░] 45% | XP: +12/tick
```

This updates every second so you can see your active grind progressing!

---

## Related Menus

From Skills menu, you can navigate to:

- **`[c]` Categories** → Browse activities by category
- **`[d]` Dashboard** → Return to main view
- **`[i]` Inventory** → Check your items
- **`[b]` Bank** → Access storage
- **`[e]` Equipment** → View equipped gear
- **`[?]` Help** → Show all controls

---

## Summary

**Quick Reference Card:**
```
┌─────────────────────────────────────┐
│  SKILLS MENU CONTROLS               │
├─────────────────────────────────────┤
│  Navigate:  ↑ ↓  or  j k            │
│  Action:    Enter  or  a            │
│  Categories:      c                 │
│  Global:    d i b e ?  Space        │
│  Save:      Ctrl+S                  │
│  Quit:      q  or  Ctrl+C           │
└─────────────────────────────────────┘
```

**Most Important:**
- **`↑/↓`** or **`j/k`** = Move cursor
- **`Enter`** or **`a`** = View skill categories
- **`c`** = Jump to Categories
- **`Space`** = Pause/Resume grinding
- **`q`** = Save and quit

Happy grinding! 🎮
