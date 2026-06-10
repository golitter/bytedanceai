## ADDED Requirements

### Requirement: Color tokens SHALL replace all hardcoded color values

All `.tsx` component files SHALL use CSS variable-based Tailwind semantic classes (e.g., `text-destructive`, `bg-success/10`, `text-foreground`) instead of:
- Hex literals (`#fff`, `#000`, `#FFFFFF`)
- `rgba()` / `rgb()` / `hsl()` literals
- Tailwind direct color classes (`text-red-500`, `bg-green-500`, `text-emerald-300`)

Inline `style={{ color: '...' }}` SHALL only be used for dynamically computed values. Static color values MUST use Tailwind classes.

#### Scenario: Hardcoded hex in RightSidebar
- **WHEN** a component uses `color: '#fff'` in inline style
- **THEN** it SHALL be replaced with `text-foreground` Tailwind class

#### Scenario: Tailwind direct color in error states
- **WHEN** a component uses `text-red-500` for error text
- **THEN** it SHALL be replaced with `text-destructive` Tailwind class

#### Scenario: Diff color indicators
- **WHEN** diff components use `bg-green-500/15` or `bg-red-500/15`
- **THEN** they SHALL be replaced with `bg-success/15` or `bg-destructive/15`

#### Scenario: Admin overlay background
- **WHEN** a component uses `background: 'rgba(0,0,0,0.5)'` in inline style
- **THEN** it SHALL be replaced with `bg-black/50` Tailwind class

### Requirement: Transitions SHALL only animate transform and opacity

All CSS transitions in component files SHALL use `transition-[transform,opacity]` with duration `120ms-300ms` and `ease-out` easing.

The following transition types are FORBIDDEN:
- `transition-colors` (animates color/background-color)
- `transition-all` (animates all properties)
- Any transition on properties other than `transform` or `opacity`

#### Scenario: Hover state transition
- **WHEN** a component uses `transition-colors duration-120` for hover effect
- **THEN** it SHALL be replaced with `transition-[transform,opacity] duration-120 ease-out`

#### Scenario: Button hover with background change
- **WHEN** a button changes background on hover without animation
- **THEN** the background change SHALL be instant (no transition), matching VSG philosophy

#### Scenario: Admin page transition-all
- **WHEN** an admin page uses `transition-all`
- **THEN** it SHALL be replaced with `transition-[transform,opacity]` or removed entirely

### Requirement: Inline styles SHALL be replaced with Tailwind classes for static values

Components SHALL NOT use `style={{ ... }}` for static CSS values (colors, borders, padding, background). These SHALL be expressed as Tailwind utility classes.

Inline styles are ONLY acceptable for:
- Dynamic computed values (e.g., `height` from measurement)
- Dynamic agent colors from constants
- Dynamic width from resize

#### Scenario: Admin card border
- **WHEN** an admin page uses `style={{ border: '1px solid var(--border)' }}`
- **THEN** it SHALL be replaced with `border border-border` Tailwind class

#### Scenario: Admin card background
- **WHEN** an admin page uses `style={{ background: 'var(--bg-card)' }}`
- **THEN** it SHALL be replaced with `bg-card` Tailwind class

#### Scenario: Dynamic sidebar width
- **WHEN** a component uses `style={{ width: dynamicWidth }}`
- **THEN** inline style is ACCEPTABLE because the value is dynamically computed

### Requirement: Font sizes SHALL conform to the defined scale

Components SHALL only use the following font sizes: 11px, 12px, 13px, 14px, 20px. The value `10px` is NOT in the allowed scale and SHALL be replaced with `11px` (the closest valid size for micro-labels).

#### Scenario: text-[10px] in AnnouncementsSection
- **WHEN** a component uses `text-[10px]` for a micro-label
- **THEN** it SHALL be replaced with `text-[11px]`

#### Scenario: text-[10px] in DiffFileTabs
- **WHEN** a diff tab uses `text-[10px]` for a count badge
- **THEN** it SHALL be replaced with `text-[11px]`

### Requirement: Magic strings SHALL be consolidated to constants

All hardcoded string literals for agent types, message roles, chat statuses, and user identity SHALL be replaced with constants from `lib/constants.ts`.

#### Scenario: Agent type string literal
- **WHEN** a component uses `'claude-code'` or `'orchestrator'` as a string literal
- **THEN** it SHALL reference `AGENT_TYPES.CLAUDE_CODE` or `AGENT_TYPES.ORCHESTRATOR` from constants

#### Scenario: Message role string literal
- **WHEN** a component uses `msg.role === 'user'`
- **THEN** it SHALL reference `MESSAGE_ROLES.USER` from constants

#### Scenario: Hardcoded username
- **WHEN** a component uses `'田乐檬'` as a hardcoded string
- **THEN** it SHALL reference a `CURRENT_USER_NAME` constant or derive from user state

#### Scenario: Status string comparison
- **WHEN** a component uses `session.status === 'streaming'`
- **THEN** it SHALL reference `CHAT_STATUSES.STREAMING` from constants

### Requirement: Admin pages SHALL use TanStack Query for server data

All 7 admin pages (StatisticsPage, DashboardPage, WorkspacePage, ServiceHealthPage, UserManagementPage, SessionCleanupPage, AgentOverviewPage) SHALL replace `useState + useEffect + fetch` patterns with `useQuery` from TanStack React Query.

#### Scenario: Admin page data loading
- **WHEN** an admin page uses `const [data, setData] = useState(null)` with `useEffect(() => { fetch(...).then(setData) })`
- **THEN** it SHALL be replaced with `const { data } = useQuery({ queryKey: ['admin-xxx'], queryFn: fetchAdminXxx })`

#### Scenario: Admin page loading state
- **WHEN** an admin page manually manages `isLoading` state
- **THEN** it SHALL use TanStack Query's built-in `isLoading` return value

#### Scenario: Admin page error state
- **WHEN** an admin page manually manages error state
- **THEN** it SHALL use TanStack Query's built-in `error` return value

### Requirement: Chat store SHALL be split by domain

The `chat.ts` store (974 lines) SHALL be split into domain-specific stores, each under 300 lines:
- `navigation-store.ts` — navigation state (active tab, sidebar visibility)
- `session-store.ts` — session list, active session, session CRUD
- `message-store.ts` — messages, streaming state, runtime blocks

#### Scenario: Store file size
- **WHEN** a Zustand store file exceeds 300 lines
- **THEN** it SHALL be split along domain boundaries

#### Scenario: Component import after split
- **WHEN** a component previously imported from `stores/chat`
- **THEN** it SHALL import from the appropriate split store file without functional changes

### Requirement: Shadows SHALL only appear on popup elements

`shadow-*` Tailwind classes SHALL only be used on popup menus, dropdowns, dialogs, and tooltips. Buttons, cards, and static panels SHALL NOT use shadows.

#### Scenario: Shadow on toggle button
- **WHEN** a toggle button uses `shadow-sm`
- **THEN** the shadow SHALL be removed, using background color difference for depth instead
