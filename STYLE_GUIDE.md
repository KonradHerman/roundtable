# Style Guide

## Color Palette (Gruvbox Dark Theme)

This project uses the Gruvbox Dark color scheme with specific accent colors.

### Base Colors

| Color Name | Hex | HSL | Usage |
|------------|-----|-----|-------|
| Background | `#282828` | `30 9% 15.7%` | Main background color |
| Foreground | `#ebdbb2` | `42 19% 85.3%` | Main text color |
| Card | `#3c3836` | `33 10% 20%` | Card backgrounds |
| Card Foreground | `#ebdbb2` | `42 19% 85.3%` | Text on cards |

### Accent Colors

| Color Name | Hex | HSL | Usage |
|------------|-----|-----|-------|
| Primary | `#d79921` | `43 100% 49%` | Primary actions, focus states, brand color (orange/gold) |
| Primary Foreground | `#282828` | `30 9% 15.7%` | Text on primary buttons |
| Secondary | `#458588` | `192 28% 40%` | Secondary actions (blue/teal) |
| Secondary Foreground | `#ebdbb2` | `42 19% 85.3%` | Text on secondary buttons |
| Accent | `#b16286` | `327 25% 54%` | Highlights, special elements (purple) |
| Accent Foreground | `#ebdbb2` | `42 19% 85.3%` | Text on accent elements |

### Utility Colors

| Color Name | Hex | HSL | Usage |
|------------|-----|-----|-------|
| Muted | `#504945` | `32 10% 27%` | Muted backgrounds |
| Muted Foreground | `#928374` | `40 11% 56%` | Secondary text, placeholders |
| Destructive | `#cc241d` | `0 62% 45%` | Errors, destructive actions |
| Destructive Foreground | `#ebdbb2` | `42 19% 85.3%` | Text on error states |
| Border | `#504945` | `32 10% 27%` | Default borders |
| Input Border | `#504945` | `32 10% 27%` | Input field borders |
| Ring | `#d79921` | `43 100% 49%` | Focus ring color |

## Input Fields

### Standard Input
```css
Background: #f9f5d7 (Light cream - less yellow than standard Gruvbox)
Text Color: #1a1a1a (nearly black for high contrast)
Placeholder: #888888 (medium gray)
Border: #504945 (muted)
Focus Border: #d79921 (primary)
```

**Rationale**: Light background provides clear contrast with the dark theme while maintaining the Gruvbox aesthetic. Dark text ensures readability against the light input background.

### Usage Example
```html
<input class="input" placeholder="Enter text" />
```

## Color Combinations

### Text on Dark Backgrounds
- **Primary text**: `#ebdbb2` on `#282828` or `#3c3836`
- **Secondary text**: `#928374` on `#282828` or `#3c3836`
- **Muted text**: `#928374` with reduced opacity

### Text on Light Backgrounds (Inputs)
- **Input text**: `#1a1a1a` on `#f9f5d7`
- **Placeholder**: `#888888` on `#f9f5d7`

### Buttons

#### Primary Button
```
Background: #d79921 (primary gold)
Text: #282828 (dark)
Hover: #d79921 with 90% opacity
```

#### Secondary Button
```
Background: #458588 (teal/blue)
Text: #ebdbb2 (light)
Hover: #458588 with 80% opacity
```

#### Destructive Button
```
Background: #cc241d (red)
Text: #ebdbb2 (light)
Hover: #cc241d with 90% opacity
```

## Typography

### Font Families
- **Body**: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif`
- **Monospace** (for codes): `font-mono` (Tailwind's default monospace stack)

### Font Sizes
- **Large heading**: `text-5xl` (48px)
- **Section heading**: `text-2xl` (24px)
- **Body text**: `text-base` (16px) or `text-lg` (18px)
- **Small text**: `text-sm` (14px)
- **Room codes**: `text-2xl` (24px) with `font-mono` and `font-bold`

## Components

### Cards
```css
Background: #3c3836
Text: #ebdbb2
Border Radius: 1rem (rounded-2xl)
Padding: 1.5rem (p-6)
Shadow: shadow-lg
```

### Borders
- **Default**: 2px solid `#504945`
- **Focus**: 2px solid `#d79921`

### Border Radius
- **Default**: `0.75rem` (rounded-xl)
- **Cards**: `1rem` (rounded-2xl)

## Spacing

### Touch Targets
- **Minimum height**: 48px for all interactive elements
- **Buttons**: 48px minimum height and width
- **Inputs**: 48px minimum height

### Padding
- **Cards**: 1.5rem (24px)
- **Buttons**: 1rem horizontal, 1rem vertical
- **Inputs**: 1rem horizontal, 1rem vertical

## Accessibility

### Contrast Ratios
- Text on dark backgrounds (`#ebdbb2` on `#282828`): ~11:1 ✓
- Input text on light backgrounds (`#1a1a1a` on `#fbf1c7`): ~14:1 ✓
- Primary button (`#282828` on `#d79921`): ~7:1 ✓

All combinations meet WCAG AAA standards for normal text.

### Focus States
- All interactive elements must have visible focus states
- Focus ring color: `#d79921` (primary)
- Focus ring width: 2px

### Touch Targets
- All interactive elements must be at least 48x48px for mobile accessibility

## Design Principles

1. **Contrast First**: Ensure high contrast between text and backgrounds, especially in input fields
2. **Gruvbox Consistency**: Stick to the Gruvbox palette for all colors
3. **Progressive Disclosure**: Show UI elements progressively as needed (e.g., reveal name input after room code)
4. **Touch-Friendly**: All interactive elements meet minimum touch target sizes
5. **No Unnecessary Clicks**: Streamline workflows to minimize required interactions

## Quick Reference

### When to Use Which Color

- **Primary (`#d79921`)**: Main CTAs, brand elements, "Host a Game" button
- **Secondary (`#458588`)**: Secondary actions, "Join Game" button
- **Muted (`#928374`)**: Labels, secondary text, disabled states
- **Destructive (`#cc241d`)**: Errors, delete actions, warnings
- **Input background (`#f9f5d7`)**: All text input fields
- **Input text (`#1a1a1a`)**: Text inside input fields (must be dark for contrast)

