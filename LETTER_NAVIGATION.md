# 🎮 Letter-Based Navigation - Quick Guide

## ✨ New Feature: Press Letter to Select!

Navigate entirely using the **first letter** of items. No more arrow keys needed!

---

## 🎯 How It Works

### 3-Step Navigation
```
[s] Skills → [w] Woodcutting → [l] Logs
     ↓            ↓              ↓
   Screen 1    Screen 2      Start Activity!
```

---

## 📋 Skills Menu (Press `s`)

Press the **first letter** of any skill:

| Letter | Skill | What You Can Do |
|--------|-------|-----------------|
| `w` | **W**oodcutting | Chop trees (Logs, Oak, Willow, etc.) |
| `m` | **M**ining | Mine ores (Copper, Iron, Gold, etc.) |
| `s` | **S**mithing | Smelt bars & craft tools |
| `r` | **R**ecycling | Recycle items into fragments |
| `c` | **C**ombat | Fight monsters (when available) |
| `a` | Cr**a**fting | Craft items (when available) |
| `k` | Coo**k**ing | Cook food (when available) |
| `g` | A**g**ility | Movement (when available) |
| `t` | **T**hieving | Steal (when available) |

**Visual:**
```
🎮 Select Skill (Press Letter)

[w] Woodcutting     Lv.15/120 [██████████░░░░░]
[m] Mining          Lv.42/120 [████████████████░░]
[s] Smithing        Lv.28/120 [████████████░░░░░]
[r] Recycling       Lv.18/120 [█████████░░░░░░░░]
```

---

## 🗂️ Category Selection

After selecting a skill, press the **first letter** of the category:

### Woodcutting Categories:
| Letter | Category | Available Woods |
|--------|----------|-----------------|
| `b` | **B**asic | Logs, Oak, Willow |
| `q` | **Q**uality | Maple, Yew, Magic |
| `e` | **E**xotic | Teak, Mahogany |

### Mining Categories:
| Letter | Category | Available Ores |
|--------|----------|----------------|
| `b` | **B**asic | Copper, Tin, Lead, Zinc, Iron |
| `i` | **I**ntermediate | Coal, Nickel, Silver, Gold, Mithril |
| `a` | **A**dvanced | Adamantite, Runite, Platinum, Obsidian |
| `g` | **G**ems | Sapphire, Emerald, Ruby, Diamond, Dragonstone |

### Smithing Categories:
| Letter | Category | What to Make |
|--------|----------|--------------|
| `b` | **B**asic | Bronze, Iron, Lead, Nickel bars |
| `a` | **A**lloys | Steel, Brass, Electrum bars |
| `p` | **P**recious | Silver, Gold, Platinum bars |
| `e` | **E**lite | Mithril, Adamantite, Runite, Obsidian bars |
| `t` | **T**ools | Axes, Pickaxes, Equipment |

**Visual:**
```
🪓 Woodcutting (Lv.15) - Select Category

[b] 🌲 Basic - Easy trees
[q] 🌳 Quality - Better wood  
[e] 🎋 Exotic - Legendary wood

Press letter to select category
```

---

## 🎪 Activity Selection

After selecting a category, press the **first letter** of the activity:

### Example: Woodcutting → Basic
| Letter | Wood | Level | Output |
|--------|------|-------|--------|
| `l` | **L**ogs | Lv.1 | 1x Logs |
| `o` | **O**ak | Lv.15 | 1x Oak Logs |
| `w` | **W**illow | Lv.30 | 1x Willow Logs |

### Example: Mining → Basic
| Letter | Ore | Level | Output |
|--------|-----|-------|--------|
| `c` | **C**opper | Lv.1 | 1x Copper Ore |
| `t` | **T**in | Lv.1 | 1x Tin Ore |
| `l` | **L**ead | Lv.10 | 1x Lead Ore |
| `z` | **Z**inc | Lv.12 | 1x Zinc Ore |
| `i` | **I**ron | Lv.15 | 1x Iron Ore |

### Example: Smithing → Alloys
| Letter | Bar | Level | Recipe |
|--------|-----|-------|--------|
| `s` | **S**teel | Lv.30 | 1 Iron + 2 Coal |
| `b` | **B**rass | Lv.15 | 1 Copper + 1 Zinc |
| `e` | **E**lectrum | Lv.55 | 1 Gold + 1 Silver |

**Visual:**
```
🪓 Woodcutting > Basic

Press letter to chop:

[l] Logs     - Basic wood
[o] Oak      - Sturdy wood    
[w] Willow   - Flexible wood

Selected: Oak
  Input: (none - just chop!)
  Output: 1x Oak Logs
```

---

## 🚀 Quick Examples

### Start Chopping Oak Logs
```bash
s    # Open Skills
w    # Select Woodcutting
b    # Select Basic category
o    # Select Oak → START GRINDING!
```

### Start Mining Iron
```bash
s    # Open Skills
m    # Select Mining
b    # Select Basic category
i    # Select Iron → START GRINDING!
```

### Smelt Steel Bar
```bash
s    # Open Skills
s    # Select Smithing
a    # Select Alloys category
s    # Select Steel → START GRINDING!
```

### Recycle Oak Logs
```bash
s    # Open Skills
r    # Select Recycling
w    # Select Wood category
o    # Select Oak → START GRINDING!
```

### Mine Sapphire
```bash
s    # Open Skills
m    # Select Mining
g    # Select Gems category
s    # Select Sapphire → START GRINDING!
```

---

## 🎨 Visual Indicators

### Hotkey Display
Every option shows its letter hotkey:
```
[w] Woodcutting     ← Press 'w'
[m] Mining          ← Press 'm'
[s] Smithing        ← Press 's'
```

### Selection Highlight
Selected items are highlighted:
```
  [l] Logs
► [o] Oak          ← Currently selected (cursor)
  [w] Willow
```

### Locked Activities
Activities you can't do yet (level too low) are dimmed:
```
[l] Logs        ← Available
[m] Mahogany    ← Dimmer (requires Lv.105)
```

---

## 🎮 Complete Navigation Reference

### Global Keys (Work Everywhere)
| Key | Action |
|-----|--------|
| `s` | Open Skills menu |
| `d` | Dashboard |
| `i` | Inventory |
| `b` | Bank |
| `e` | Equipment |
| `Space` | Pause/Resume grind |
| `Ctrl+S` | Save game |
| `q` or `Ctrl+C` | Save & quit |
| `?` or `h` | Help |

### Navigation Keys
| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Move cursor up/down |
| `Enter` | Confirm selection |
| `Esc` | Go back |
| **Letter** | Quick select by first letter |

### Menu-Specific
| Menu | Press | Result |
|------|-------|--------|
| Skills | `w` | Go to Woodcutting categories |
| Categories | `b` | Select Basic category |
| Activities | `o` | Start Oak/Obsidian/etc. |

---

## 💡 Pro Tips

### Speed Run Navigation
Once you learn the letters, you can grind super fast:
```bash
s → w → b → l    # Chop Logs in 4 keystrokes!
s → m → b → c    # Mine Copper in 4 keystrokes!
s → s → a → s    # Smelt Steel in 4 keystrokes!
```

### Muscle Memory
Common grinds become automatic:
- **Logs**: `s → w → b → l`
- **Copper**: `s → m → b → c`
- **Bronze**: `s → s → b → b`
- **Oak**: `s → w → b → o`

### Arrow Keys Still Work
Don't like letters? Use arrow keys:
```
s          # Skills
↓ ↓ ↓      # Navigate to Mining
Enter      # Select
↓ ↓        # Navigate to Intermediate
Enter      # Select
↓ ↓ ↓      # Navigate to Gold
Enter      # Start!
```

---

## 🎯 Quick Reference Card

```
┌─────────────────────────────────────┐
│  LETTER NAVIGATION                  │
├─────────────────────────────────────┤
│                                     │
│  [s] → [w] → [b] → [l] = Chop Logs │
│  [s] → [m] → [b] → [c] = Mine Cu   │
│  [s] → [s] → [a] → [s] = Smelt St  │
│  [s] → [r] → [w] → [o] = Recycle   │
│                                     │
├─────────────────────────────────────┤
│  SKILLS        CATEGORIES           │
├─────────────────────────────────────┤
│  [w]oodcutting [b]asic              │
│  [m]ining      [i]ntermediate       │
│  [s]mithing    [a]dvanced           │
│  [r]ecycling   [e]lite              │
│  [c]ombat      [g]ems               │
│  cr[a]fting    [t]ools              │
│  coo[k]ing                          │
│  a[g]ility                          │
│  [t]hieving                         │
│                                     │
└─────────────────────────────────────┘
```

---

## 🌟 Status Bar (Always Visible)

Even while navigating, the status bar shows your grind:
```
 ⛏️ Mine Iron [████████████████████░░░░░░░░░░░░░░░░] 45% | XP: +25/tick
```

This updates every second so you can see progress!

---

## 🎊 Summary

**Before:**
- Navigate with arrows only
- Slow to find activities
- Easy to get lost in menus

**After:**
- ✅ Press letter to jump directly
- ✅ Lightning-fast navigation
- ✅ No more arrow key spam!
- ✅ Works with arrow keys too (your choice)

**Example:**
```
OLD: s → ↓↓↓↓ Enter ↓ Enter ↓↓↓ Enter (8 keys)
NEW: s → w → b → o                    (4 keys)
     
2x FASTER! 🚀
```

---

## 🚀 Try It Now!

```bash
./afk-tui-letter-nav
```

Then try:
```
s → w → b → l    # Chop Logs
```

Enjoy the speed! ⚡
