# Design System — Hello World

> Source of truth: the approved `index.html` (preview: approved design URL in workspace).
> Every value below is extracted from it. Changing a value here without changing the approved design is a defect.

Last updated: 2025-02-14

## 1. Foundations

### 1.1 Color

Semantic tokens. Name by job, never by hue.

| Token | Value | Used for |
|---|---|---|
| `--color-bg` | `#f8fafc` | Page background |
| `--color-surface` | `#ffffff` | Card / panel background |
| `--color-surface-raised` | `#ffffff` | Floating surfaces shown in mockup |
| `--color-border` | `#e2e8f0` | Default border, divider |
| `--color-text` | `#0f172a` | Body text |
| `--color-text-muted` | `#64748b` | Secondary text, captions |
| `--color-primary` | `#4f46e5` | Primary action background |
| `--color-primary-text` | `#ffffff` | Text on primary |
| `--color-accent-secondary` | `#14b8a6` | Secondary accent, status accent |
| `--color-success` | `#16a34a` | Success state |
| `--color-danger` | `#dc2626` | Error state |
| `--color-focus` | `color-mix(in srgb, #4f46e5 70%, white)` | Focus ring |

#### Contrast audit

| Foreground | Background | Ratio | Passes |
|---|---|---|---|
| `--color-text` | `--color-bg` | `16.1:1` | AA |
| `--color-text` | `--color-surface` | `16.1:1` | AA |
| `--color-text-muted` | `--color-surface` | `4.6:1` | AA |
| `--color-primary-text` | `--color-primary` | `5.2:1` | AA |
| `--color-primary` | `--color-surface` | `5.0:1` | AA for UI |
| `--color-success` | `--color-surface` | `3.3:1` | AA Large / UI |
| `--color-danger` | `--color-surface` | `4.8:1` | AA |
| `--color-accent-secondary` | `--color-surface` | `2.5:1` | FAIL for body; used as accent only |

### 1.2 Spacing

Base unit: `4px`. Every margin, padding, and gap in the product uses one of these.

| Token | Value |
|---|---|
| `--space-1` | `4px` |
| `--space-2` | `8px` |
| `--space-3` | `12px` |
| `--space-4` | `16px` |
| `--space-5` | `20px` |
| `--space-6` | `24px` |
| `--space-7` | `28px` |
| `--space-8` | `32px` |
| `--space-9` | `36px` |
| `--space-11` | `44px` |
| `--space-12` | `48px` |
| `--space-14` | `56px` |

### 1.3 Typography

Font families:

- Body: `Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`
- Headings: `Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`
- Mono: browser default monospace fallback only; no explicit mono face in approved design

| Token | Size | Line height | Weight | Used for |
|---|---|---|---|---|
| `--text-xs` | `12px` | `1.4` | `700` | Eyebrow labels |
| `--text-sm` | `14px` | `1.5` | `400/600` | Secondary body, field labels |
| `--text-base` | `16px` | `1.5` | `400` | Body |
| `--text-lg` | `17px` | `1.5` | `400` | Lead paragraph |
| `--text-xl` | `18px` | `1.2` | `600` | Section headings |
| `--text-2xl` | `32px` | `1.05` | `700` | Main headline |
| `--text-3xl` | `64px` | `1.05` | `700` | Large hero scaling via `clamp()` |

Heading levels are used in order and never skipped for visual sizing.

| Token | Value | Used for |
|---|---|---|
| `--font-weight-body` | `400` | Running text |
| `--font-weight-medium` | `600` | Labels, emphasis |
| `--font-weight-heading` | `700` | H1 / strong emphasis |
| `--tracking-tight` | `-0.01em` | Display headings |
| `--tracking-normal` | `0` | Everything else |

### 1.4 Radius, border, shadow, motion

| Token | Value | Used for |
|---|---|---|
| `--radius-xs` | `11px` | Small icon tiles |
| `--radius-sm` | `14px` | Button, input, mini panel |
| `--radius-md` | `16px` | Error box, list item |
| `--radius-lg` | `18px` | Cards, larger panels |
| `--radius-xl` | `22px` | Main panels |
| `--radius-full` | `9999px` | Pills, status chips |
| `--border-width` | `1px` | Default border |
| `--shadow-sm` | `0 10px 25px rgba(79, 70, 229, .25)` | Brand mark only |
| `--shadow-md` | `0 12px 24px rgba(79, 70, 229, .22)` | Hovered primary button |
| `--shadow-lg` | `0 18px 50px rgba(15, 23, 42, .08)` | Resting panel |
| `--duration-fast` | `.18s` | Hover, focus, press |
| `--duration-base` | `.35s` | New list item pop |
| `--easing` | `ease` | All transitions |

Motion respects `prefers-reduced-motion: reduce`: state changes remain, movement is removed.

### 1.5 Layout and breakpoints

| Name | Min width | Container | Columns | Gutter |
|---|---|---|---|---|
| `sm` | `0px` | `100%` | `1` | `16px` |
| `md` | `900px` | `100%` | `2` | `20px` |
| `lg` | `1120px` | `1120px` | `2` | `20px` |
| `xl` | `1280px` | `1120px` | `2` | `20px` |

Z-index scale (only these values are allowed):

| Layer | Value |
|---|---|
| Base | `0` |
| Sticky header | `10` |
| Dropdown | `20` |
| Modal backdrop | `40` |
| Modal | `50` |
| Toast | `60` |

## 2. Components

One subsection per reusable component. Every component lists all states that appear in approved design.

### 2.1 Status pill

**Purpose** — show API health or short state label. Use for top status only; not for long text.

**Anatomy** — `[dot] [label]`

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-surface`, `--color-border`, `--color-text` | Neutral status chip |
| Success | `--color-success` | Connected / healthy state |
| Danger | `--color-danger` | Unreachable / error state |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | `38px` | `10px 14px` | `--text-sm` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Neutral border, white fill | `--color-surface`, `--color-border` |
| Hover | No special hover shown | None |
| Focus (keyboard) | No interactive focus in mockup | None |
| Active / pressed | None shown | None |
| Disabled | Not shown | None |
| Loading | Not shown | None |
| Error | Red dot / label for unreachable API | `--color-danger` |
| Empty | Not shown | None |

**Accessibility** — status text readable without color; no keyboard interaction shown.

### 2.2 Button

**Purpose** — primary form and refresh action. Use for direct actions only.

**Anatomy** — `[label]`

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Primary | `--color-primary`, `--color-primary-text` | Main action |
| Secondary | `#e0e7ff`, `#1e1b4b` | Support action |
| Ghost | `#ffffff`, `--color-border`, `--color-text` | Low-emphasis action |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | `44px` | `12px 16px` | `--text-sm` / `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Flat fill, rounded rectangle | variant tokens |
| Hover | Lift `translateY(-1px)` and shadow on primary | `--shadow-md`, `--duration-fast` |
| Focus (keyboard) | Visible outline ring | `--color-focus` |
| Active / pressed | Not explicitly shown beyond hover transition | None |
| Disabled | Not shown | None |
| Loading | Not shown | None |
| Error | Not shown | None |
| Empty | Not shown | None |

**Accessibility** — minimum hit target met by padding; keyboard focus visible.

### 2.3 Text input

**Purpose** — accept name and message. Use for single-line text field.

**Anatomy** — `[label] [input]`

**Variants**

| Variant | Tokens | When to use |
|---|---|---|
| Default | `--color-surface`, `--color-border`, `--color-text` | Normal input |

**Sizes**

| Size | Height | Padding | Text token |
|---|---|---|---|
| Default | `44px` min | `12px 14px` | `--text-base` |

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White fill, 1px border | `--color-surface`, `--color-border` |
| Hover | None shown | None |
| Focus (keyboard) | Visible focus ring | `--color-focus` |
| Active / pressed | Typing state only | None |
| Disabled | Not shown | None |
| Loading | Not shown | None |
| Error | Not shown as field border; page-level validation message used instead | `--color-danger` |
| Empty | Placeholder text visible | `--color-text-muted` |

**Accessibility** — label associated by `for` / `id`.

### 2.4 Textarea

**Purpose** — accept greeting message.

**Anatomy** — `[label] [textarea]`

**States** — same as Text input.

### 2.5 Card / panel

**Purpose** — group sections on page.

**Anatomy** — `[content]`

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White surface, border, subtle shadow | `--color-surface`, `--color-border`, `--shadow-lg`, `--radius-xl` |
| Hover | None on main panels | None |
| Focus (keyboard) | Not interactive | None |
| Active / pressed | None | None |
| Disabled | None | None |
| Loading | Not shown | None |
| Error | Error variant uses pink/red surface and border | `#fff1f2`, `#fecaca`, `#9f1239` |
| Empty | Empty state shown inside panel, not as a separate card style | None |

### 2.6 Greeting list item

**Purpose** — show stored greetings newest first.

**Anatomy** — `[name] [created_at] [message]`

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | White card with border | `--color-surface`, `--color-border`, `--radius-lg` |
| Hover | Slight lift on pointer hover | `--duration-fast` |
| Focus (keyboard) | Not interactive | None |
| Active / pressed | None | None |
| Disabled | None | None |
| Loading | Not shown | None |
| Error | None | None |
| Empty | Not rendered; list wrapper shows empty message instead | None |

### 2.7 Empty state

**Purpose** — explain missing greetings and next action.

**Anatomy** — `[message]`

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Dashed border, muted copy | `--color-border`, `--color-text-muted`, `--color-bg` |
| Hover | None | None |
| Focus (keyboard) | None | None |
| Active / pressed | None | None |
| Disabled | None | None |
| Loading | Not shown | None |
| Error | Not shown | None |
| Empty | This is its only shown state | None |

### 2.8 Error box

**Purpose** — show API unreachable and validation messages.

**Anatomy** — `[message]`

**States**

| State | Visual change | Tokens |
|---|---|---|
| Default | Solid light-red box or filled pale-red alert | `#fef2f2`, `#fee2e2`, `#fecaca`, `#991b1b`, `#9f1239` |
| Hover | None | None |
| Focus (keyboard) | None | None |
| Active / pressed | None | None |
| Disabled | None | None |
| Loading | None | None |
| Error | This is the error state | same |
| Empty | Hidden when not needed | None |

## 3. Content and formatting

- Voice: friendly, plain, short.
- Dates: `YYYY-MM-DD HH:MM` in mockup; app may format timezone consistently, but must keep newest-first ordering.
- Numbers: plain decimal, no grouping needed in this demo.
- Buttons: sentence case.
- Headings: sentence case with short verbs.
- Labels: title case.
- Empty-state and error-message wording pattern: state first, next action second.

## 4. Known deviations

| Where | Deviation | Why it stands | Follow-up |
|---|---|---|---|
| Header status area | Mockup shows `Health: ok` and `API unreachable state` side by side, though runtime can show only one live chip | Stakeholder explicitly called it out as a runtime note | App must render one real status chip only |
| Live hello panel | Mockup shows `Connected` chip even though API may be down | Static preview artifact, not runtime rule | Use real API result in app |
| Accent use | Secondary teal accent has low contrast on white if used as body text | It is used only as icon / accent fill | Keep it non-text |

## 5. Change log

| Date | Change | Design PR |
|---|---|---|
| 2025-02-14 | Initial design system extracted from approved mockup | TBD |
