## ADDED Requirements

### Requirement: Dark mode CSS variable mapping
The system SHALL map shadcn/ui CSS variables to visual-style-guide color values in the `.dark` block. The mapping MUST preserve shadcn variable names so existing shadcn components work without modification.

#### Scenario: Dark mode active by default
- **WHEN** application loads
- **THEN** `<html>` element has `.dark` class and all shadcn variables resolve to visual-style-guide dark color values

#### Scenario: Background hierarchy
- **WHEN** dark mode is active
- **THEN** background variables map to the 5-level gray scale: canvas `#0A0B0E`, sidebar `#111318`, card `#1A1D24`, hover `#22262F`, active `#2C313B`

#### Scenario: Border and text colors
- **WHEN** dark mode is active
- **THEN** borders use `rgba(255,255,255,0.06)`, primary text `#E8EBF0`, secondary text `#8B91A0`, tertiary text `#5A6070`

### Requirement: Visual style guide custom CSS variables
The system SHALL define additional CSS variables for visual-style-guide concepts not covered by shadcn's default set: `--bg-canvas`, `--bg-sidebar`, `--bg-hover`, `--bg-active`, `--text-primary`, `--text-secondary`, `--text-tertiary`, `--color-brand`, `--color-success`, `--color-warning`, `--color-error`, `--agent-claude`, `--agent-opencode`, `--agent-orchestrator`.

#### Scenario: Custom variables available
- **WHEN** a component needs to reference agent identity colors or brand color
- **THEN** it can use `var(--agent-claude)`, `var(--color-brand)`, etc. via Tailwind's arbitrary value syntax

### Requirement: Light theme placeholder values
The system SHALL keep reasonable light theme values in the `:root` block so that removing `.dark` class does not break the application.

#### Scenario: Light mode switch
- **WHEN** `.dark` class is removed from `<html>`
- **THEN** all shadcn components render with existing light theme values without visual breakage

### Requirement: Typography adherence
The system SHALL use Geist Sans for UI text and Geist Mono for code text. Font sizes MUST follow the visual-style-guide type scale: 20px page title, 14px block title, 14px body, 12px auxiliary, 13px code, 11px micro label. Letter-spacing for UI text MUST be `-0.01em`, for uppercase labels `+0.05em`.

#### Scenario: Code block typography
- **WHEN** a code block is rendered
- **THEN** it uses Geist Mono at 13px with 1.65 line-height

#### Scenario: Badge and label typography
- **WHEN** a badge or status label is rendered
- **THEN** it uses 11px font at Medium weight with +0.05em letter-spacing

### Requirement: Spacing system
The system SHALL use 4px as the base spacing unit. Component padding MUST follow visual-style-guide specs: buttons 8px 16px (small) / 10px 20px (standard), cards 16px, inputs 10px 14px, list items 10px 12px.

#### Scenario: Consistent spacing
- **WHEN** chat UI components are rendered
- **THEN** spacing uses 4px base unit multiples (4, 8, 12, 16, 24, 32px)

### Requirement: Animation constraints
The system SHALL only animate `transform` and `opacity` properties. Animation durations: 120ms micro-interactions (hover/focus), 200ms state changes, 250ms panel transitions, 300ms page transitions. All use `ease-out`. No backdrop-blur animations, gradient flows, particle effects, 3D rotations, or spring physics.

#### Scenario: Hover animation
- **WHEN** user hovers over an interactive element
- **THEN** background color transitions in 120ms ease-out

#### Scenario: Streaming cursor animation
- **WHEN** agent is streaming a response
- **THEN** a blinking cursor `▌` animates with opacity 0→1→0 in 1s loop

#### Scenario: No forbidden animations
- **WHEN** the application runs
- **THEN** no backdrop-blur animations, gradient flows, particle effects, 3D rotations, or spring physics are present
