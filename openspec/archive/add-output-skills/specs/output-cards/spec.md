## ADDED Requirements

### Requirement: HtmlCard renders HTML in sandboxed iframe
HtmlCard SHALL render HTML content inside an `<iframe>` with `sandbox` attribute. The HTML content SHALL come from the HtmlBlock's content field, injected via `srcdoc`.

#### Scenario: HtmlCard renders styled HTML
- **WHEN** HtmlBlock contains `<div style="background:red;padding:20px">Hello</div>`
- **THEN** HtmlCard SHALL render an iframe showing a red box with "Hello"

#### Scenario: HtmlCard sandbox restricts scripts
- **WHEN** HtmlBlock contains `<script>alert(1)</script>`
- **THEN** HtmlCard iframe SHALL NOT execute the script (sandbox attribute)

### Requirement: ImageCard displays workspace image
ImageCard SHALL display an image from the workspace via the file API proxy. Image path SHALL be relative to workspace root.

#### Scenario: ImageCard shows image
- **WHEN** ImageBlock contains `path: "chart.png"`
- **THEN** ImageCard SHALL render `<img src="/api/workspace/{id}/files/chart.png" />`

#### Scenario: ImageCard shows fallback for missing image
- **WHEN** image API returns 404
- **THEN** ImageCard SHALL show a "图片加载失败" placeholder

### Requirement: AttachmentCard provides download link
AttachmentCard SHALL display file name, size, and a download button. Download SHALL go through the Go proxy file API.

#### Scenario: AttachmentCard shows download link
- **WHEN** AttachmentBlock contains `path: "report.pdf"`
- **THEN** AttachmentCard SHALL render file name "report.pdf" and a download link pointing to `/api/workspace/{id}/files/report.pdf`

### Requirement: DiffCard shows workspace diff with edit capability
DiffCard SHALL fetch diff from workspace API on mount, render with react-diff-viewer, and provide "编辑文件" / "接受变更" / "拒绝变更" actions.

#### Scenario: DiffCard fetches and renders diff
- **WHEN** DiffBlock is rendered for a workspace
- **THEN** DiffCard SHALL call `GET /api/workspace/{id}/diff` and render the unified diff

#### Scenario: DiffCard edit opens CodeMirror
- **WHEN** user clicks "编辑文件" on a file in the diff
- **THEN** DiffCard SHALL fetch the full file content via file API and open a CodeMirror editor

#### Scenario: DiffCard accept commits changes
- **WHEN** user clicks "接受变更"
- **THEN** DiffCard SHALL call `POST /api/workspace/{id}/commit`

#### Scenario: DiffCard reject reverts changes
- **WHEN** user clicks "拒绝变更"
- **THEN** DiffCard SHALL call `POST /api/workspace/{id}/revert`

### Requirement: PreviewCard renders iframe with agent preview URL
PreviewCard SHALL render an `<iframe>` pointing to the URL from PreviewBlock. The URL is a local HTTP server started by AgentEnd.

#### Scenario: PreviewCard shows preview
- **WHEN** PreviewBlock contains `url: "http://localhost:3928/index.html"`
- **THEN** PreviewCard SHALL render `<iframe src="http://localhost:3928/index.html" />`

#### Scenario: PreviewCard open in new tab
- **WHEN** user clicks "在新标签页打开"
- **THEN** SHALL open the preview URL in a new browser tab

### Requirement: MessageBubble renders blocks array
MessageBubble SHALL iterate over `message.blocks` and render each block with the corresponding card component based on `block.type`.

#### Scenario: MessageBubble renders mixed blocks
- **WHEN** a message has blocks: `[TextBlock, HtmlBlock, TextBlock, DiffBlock]`
- **THEN** MessageBubble SHALL render MarkdownRenderer, HtmlCard, MarkdownRenderer, DiffCard in order
