# AEOlyzer Design System

## Overview
AEOlyzer is a modern AI chatbot interface inspired by Claude, featuring a sophisticated dark theme with carefully curated colors, typography, and interactions. The design prioritizes clarity, accessibility, and a premium user experience.

---

## Color Palette

### Primary Colors (System/Dark Mode)
- **Background (Primary)**: `#2b2a27` - Deep charcoal base for main content areas
- **Background (Secondary)**: `#393836` - Slightly lighter for input fields and cards
- **Background (Tertiary)**: `#343330` - For badges and subtle containers
- **Background (Dark Hover)**: `#252422` - Active/hover state for all interactive elements

### Primary Colors (Light Mode)
- **Background (Primary)**: `#faf9f5` - Warm off-white base for main content areas
- **Background (Secondary)**: `#f0efeb` - Slightly darker for input fields and cards
- **Background (Tertiary)**: `#e6e5e1` - For badges and subtle containers
- **Background (Hover)**: `#e6e5e1` - Active/hover state for all interactive elements

### Brand Color
- **Orange (Brand)**: `#e07b53` - Primary accent used for:
  - Logo and sunburst animation
  - "Upgrade" links and badges
  - Toggle switches when enabled
  - Focus rings and selection states
  - Checkmarks and selection indicators

### Text Colors (System/Dark Mode)
- **Text (Primary)**: `#ececec` - Bright white for main content and headings
- **Text (Secondary)**: `#a3a29e` - Muted grey for sidebar items, icons, placeholders
- **Text (Tertiary)**: `#7a7974` - Lighter grey for subtle text and hints
- **Text (Warm)**: `#b8977e` - Warm beige for greeting text and brand messaging
- **Text (Sidebar)**: `#d4d4d0` - Slightly greyish white for sidebar brand name

### Text Colors (Light Mode)
- **Text (Primary)**: `#2b2a27` - Charcoal for main content and headings
- **Text (Secondary)**: `#7a7974` - Muted grey for sidebar items, icons, placeholders
- **Text (Tertiary)**: `#a3a29e` - Lighter grey for subtle text and hints
- **Text (Warm)**: `#b8977e` - Warm beige for greeting text and brand messaging
- **Text (Sidebar)**: `#2b2a27` - Charcoal for sidebar brand name

### Border & Divider
- **Borders**: `#4a4945` - Subtle grey for input borders and dividers
- **Scrollbar**: `#4a4945` - Matches border color for consistency

---

## Typography

### Font Families
- **Display/Headings**: Rokkitt (weight: 100-300) - Thin serif font for greeting text and prominent headings
- **Body/UI**: Geist (default) - Modern sans-serif for all UI text, buttons, and body content
- **Monospace**: Geist Mono (for code blocks if needed)

### Sizing & Usage
- **Greeting Text**: `3rem` (48px), font-light, line-height 1.1
- **Heading 1**: `2.75rem` - Large titles
- **Body**: `15px` (0.9375rem), line-height 1.75
- **Small**: `13px` (0.8125rem) - Secondary information
- **Tiny**: `11px` (0.6875rem) - Badges, labels

### Font Weights
- Light (`300`) - Greeting text, welcoming content
- Regular (`400`) - Standard body text
- Medium (`500`) - Emphasized text, some button labels

---

## Layout & Spacing

### Container Widths
- **Chat Input**: `max-w-2xl` (42rem) - Narrower, focused input area
- **Main Content**: Full width with padding
- **Sidebar**: Fixed width with collapsible toggle

### Spacing Scale
All spacing follows Tailwind's spacing scale:
- Padding/Margin: `px-4`, `py-2`, `mb-4`, `gap-4` (consistent 4px units)
- Between welcome text and textarea: `mb-4` (16px)
- Between greeting and free plan badge: `mb-6` (24px)
- Between logo and greeting: `gap-4` (16px)

### Corner Radius
- **Large containers**: `rounded-xl` (12px) - Input field, modals
- **Medium elements**: `rounded-lg` (8px) - Buttons, quick actions
- **Small elements**: `rounded-md` (6px) - Icon buttons, subtle interactive elements
- **Notes**: Corners are less rounded than typical UI, creating a more modern, sharp aesthetic

---

## Components

### Input Area
- **Height**: Min 32px, max 200px with auto-resize
- **Padding**: `pt-5 pb-4 px-4` - Slightly generous padding
- **Default Border**: 1px subtle border (`border-muted-foreground/20` in light mode, `border-border` in dark mode)
- **Focus State**: Very thin, sharp 1px border matching the brand color (`border-accent` / `#e07b53`). No fuzzy outset rings (`ring-0`) are used to maintain a crisp look.
- **Caret Color**: Matches text color for visibility (`caret-foreground`)
- **Placeholder**: Light grey (`#7a7974`), font-light, letter-spacing 0.01em

### Quick Action Buttons
- **Layout**: Flex wrap, centered, gap-2
- **Styling**: Dark background (`#2b2a27`), 1px border (`#4a4945`)
- **Padding**: `px-4 py-2`
- **Hover State**: Dark background (`#252422`)
- **Icon Size**: 16px, stroke-width 1.5 for thin elegant appearance

### Model Selector Dropdown
- **Selected Item**: Ring-2 with orange (`#e07b53`)
- **Hover**: Dark background (`#252422`)
- **Upgrade Badge**: Orange background with white text

### Toggle/Switch Controls
- **Size**: `w-11 h-6` - Standard toggle dimensions
- **Off State**: Grey background (`#4a4945`)
- **On State**: Orange background (`#e07b53`)
- **Animation**: Smooth transition
- **Thumb**: Circular, white, with proper padding

### Selection/Theme Buttons
- **Ring**: `ring-2` with orange (`#e07b53`) when selected
- **Border**: `1px solid #4a4945` when unselected
- **Hover**: Dark background (`#252422`)

---

## Interactions & States

### Hover States
- **All Interactive Elements**: Change background to dark hover color (`#252422`)
- **Consistency**: Never use bright/light hovers; always go darker
- **Transition**: Smooth `transition-colors` for all state changes

### Focus States
- **Form Inputs & Textareas**: A crisp 1px solid border using the brand orange (`border-accent`). Do not use outset fuzzy shadows or `ring-` utilities for input focus to ensure it stays sharp and elegant.
- **Buttons**: Generally no ring, rely on background change (`hover:bg-muted` or similar)

### Active/Selected States
- **Theme Selection**: `border-accent` + dark background
- **Tab Navigation**: Dark background (`#252422`)
- **Checkmarks**: Orange color (`#e07b53`) for visibility
- **Model Selected**: Orange checkmark in dropdown

### Loading/Thinking State
- **Logo Animation**: Spinning outer ring with pulsing center
- **Color**: Orange (`#e07b53`)
- **"Thinking..." text**: Muted grey (`#a3a29e`)

---

## Sidebar

### Spacing
- **Navigation Items**: `py-1.5` - Compact vertical spacing
- **Recent Chats**: `py-1` - Even tighter spacing
- **Gap Between Sections**: Clear visual separation with muted text headers

### Text Colors
- **Active Item**: Bright white (`#ececec`) with dark hover background
- **Inactive Item**: Muted grey (`#a3a29e`) that brightens on hover to (`#d4d4d0`)
- **Brand Name**: Soft grey (`#d4d4d0`)
- **Section Headers**: Very muted (`#a3a29e`)

### Icons
- **Size**: 20px for navigation, 16px for quick actions
- **Stroke Width**: 1.5 for thin, elegant appearance
- **Color**: Inherits from text color (muted grey by default)

---

## Accessibility

### Color Contrast
- Text on dark backgrounds meets WCAG AA standards
- Orange accent color (#e07b53) chosen for visibility against dark backgrounds
- Muted greys reserved for secondary/tertiary information

### Focus Indicators
- All interactive elements have visible focus states
- Focus rings use the brand orange for clarity
- Keyboard navigation fully supported

### Typography
- Minimum font size: 11px only for badges/labels
- Standard body text: 15px for comfortable reading
- Generous line height: 1.75 for body, 1.1 for headings

---

## Theme Modes (System & Light Mode)

The application uses the dark theme as the default **System Mode**, alongside a carefully calibrated **Light Mode**:
- **System Mode (Dark)**: Optimized for low-light viewing and reduced eye strain through strategic use of warm accents.
- **Light Mode**: Inverted color hierarchy maintaining the warm beige and orange accents, optimized for bright environments.
- Transitions between themes are immediate but respect user OS preferences via the `next-themes` provider.

---

## Design Principles

1. **Strictly Tailwind CSS**: Tailwind is the sole source of truth for styling. Absolutely no inline `style={{...}}` objects or hardcoded hex colors inside React components.
2. **Semantic Design Tokens**: All colors use CSS variables mapped through Tailwind v4's `@theme inline` block in `globals.css` (e.g., `text-sidebar-text`, `bg-sidebar-hover`). Do not use direct hex codes.
3. **Minimalist Elegance**: Clean, uncluttered interface with purposeful use of whitespace.
4. **Consistent Interaction**: All hover and focus states follow the same pattern (native Tailwind `hover:`, `group-hover:`, etc.). Avoid using React `onMouseEnter/onMouseLeave` state just for visual changes.
5. **Brand Consistency**: Orange accent (`#e07b53` / `border-accent`) used sparingly but consistently for active states and input focus.
6. **Typography Hierarchy**: Serif display font for welcoming content, sans-serif for functional UI.
7. **Accessibility First**: High contrast text (especially in Light Mode), clear focus states, semantic sizing.
8. **Performance**: Native CSS transitions, no layout jumping (e.g., icons that appear on hover should use `opacity-0 group-hover:opacity-100` rather than conditional React rendering).

---

## Future Considerations

- Gradient overlays for visual depth (if needed)
- Additional accent colors for error/warning states
- Animation refinement for loading states and transitions
