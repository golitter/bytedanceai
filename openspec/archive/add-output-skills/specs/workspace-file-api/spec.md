## ADDED Requirements

### Requirement: Workspace file read API
AgentEnd SHALL provide `GET /v1/workspace/{id}/files/{path}` to read files from the workspace worktree. Files SHALL be returned using streaming FileResponse, not read-all-into-memory.

#### Scenario: Read image file
- **WHEN** `GET /v1/workspace/{id}/files/chart.png` is called
- **THEN** SHALL return the file with Content-Type `image/png` using FileResponse

#### Scenario: Read non-existent file
- **WHEN** `GET /v1/workspace/{id}/files/nonexistent.txt` is called
- **THEN** SHALL return HTTP 404

#### Scenario: Path traversal prevention
- **WHEN** `GET /v1/workspace/{id}/files/../../etc/passwd` is called
- **THEN** SHALL return HTTP 403

### Requirement: Workspace file write API
AgentEnd SHALL provide `PUT /v1/workspace/{id}/files/{path}` to write files to the workspace worktree. This enables DiffCard edit functionality.

#### Scenario: Write edited file
- **WHEN** `PUT /v1/workspace/{id}/files/src/app.tsx` with file content body
- **THEN** SHALL write the content to the worktree file and return HTTP 200

### Requirement: Workspace diff API
AgentEnd SHALL provide `GET /v1/workspace/{id}/diff` to get the full workspace diff against HEAD. SHALL execute `git diff HEAD` in the worktree directory.

#### Scenario: Get workspace diff
- **WHEN** `GET /v1/workspace/{id}/diff` is called
- **THEN** SHALL return the unified diff text of all uncommitted changes

#### Scenario: No changes
- **WHEN** workspace has no uncommitted changes
- **THEN** SHALL return HTTP 200 with empty string body

### Requirement: Workspace commit API
AgentEnd SHALL provide `POST /v1/workspace/{id}/commit` to commit current changes in the worktree.

#### Scenario: Commit changes
- **WHEN** `POST /v1/workspace/{id}/commit` is called
- **THEN** SHALL execute `git add -A && git commit` in the worktree

### Requirement: Workspace revert API
AgentEnd SHALL provide `POST /v1/workspace/{id}/revert` to discard uncommitted changes in the worktree.

#### Scenario: Revert all changes
- **WHEN** `POST /v1/workspace/{id}/revert` is called
- **THEN** SHALL execute `git checkout HEAD -- .` in the worktree

### Requirement: Go backend proxies all workspace APIs
Go backend SHALL proxy all workspace file/diff/commit/revert API calls to AgentEnd without storing any data. Go SHALL use `io.Copy` for file transfer, not `io.ReadAll`.

#### Scenario: Go proxies file read
- **WHEN** `GET /api/workspace/{id}/files/chart.png` hits Go backend
- **THEN** Go SHALL proxy to AgentEnd `GET /v1/workspace/{id}/files/chart.png` using streaming response

#### Scenario: Go proxies diff request
- **WHEN** `GET /api/workspace/{id}/diff` hits Go backend
- **THEN** Go SHALL proxy to AgentEnd and return the diff text
