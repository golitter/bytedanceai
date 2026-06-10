## ADDED Requirements

### Requirement: AgentEnd starts preview HTTP server
AgentEnd SHALL be able to start a local HTTP server serving workspace files for preview. The server SHALL serve static files from the workspace worktree directory.

#### Scenario: Start preview server
- **WHEN** Agent wants to preview a workspace
- **THEN** AgentEnd SHALL start an HTTP server on a dynamically allocated port serving the worktree directory

#### Scenario: Preview server serves files
- **WHEN** preview server receives `GET /index.html`
- **THEN** SHALL return the file from the worktree with correct Content-Type

#### Scenario: Preview server serves sub-resources
- **WHEN** preview server receives `GET /styles/main.css`
- **THEN** SHALL return the CSS file with Content-Type `text/css`

### Requirement: Preview URL communicated via aka_yhy block
Agent SHALL include the preview server URL in the `aka_yhy preview` block so the frontend knows where to point the iframe.

#### Scenario: Agent outputs preview URL
- **WHEN** Agent creates a web page and starts preview
- **THEN** SHALL output ` ```aka_yhy\ntype: preview\nurl: http://localhost:{port}/index.html\n``` `

### Requirement: Preview server lifecycle
Preview server SHALL remain running as long as the workspace is active. SHALL be cleaned up when workspace is cleaned up.

#### Scenario: Preview server cleanup
- **WHEN** workspace is cleaned up or removed
- **THEN** associated preview server SHALL be stopped
