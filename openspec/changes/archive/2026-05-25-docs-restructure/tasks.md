## 1. Root docs/ migration

- [x] 1.1 Create new directories: `docs/design/`, `docs/reference/`, `docs/guides/`, `docs/testing/`
- [x] 1.2 Move `docs/architecture/three-tier-design.md` → `docs/design/three-tier-design.md`
- [x] 1.3 Move `docs/common/skills.md` → `docs/reference/skills.md`
- [x] 1.4 Move `docs/common/contract-layer.md` → `docs/guides/contract-layer.md`
- [x] 1.5 Move `docs/common/git-conventions.md` → `docs/guides/git-conventions.md`
- [x] 1.6 Move `docs/common/makefile-guide.md` → `docs/guides/makefile-guide.md`
- [x] 1.7 Move `docs/setup.md` → `docs/guides/setup.md`
- [x] 1.8 Move `docs/monorepo-setup.md` → `docs/guides/monorepo-setup.md`
- [x] 1.9 Move `docs/test_guide/inactive-cleanup.md` → `docs/testing/inactive-cleanup.md`
- [x] 1.10 Remove empty old directories: `docs/architecture/`, `docs/common/`, `docs/test_guide/`

## 2. Frontend docs/ migration

- [x] 2.1 Create new directories: `frontend/docs/design/`, `frontend/docs/reference/`
- [x] 2.2 Move `frontend/docs/impl/01-architecture.md` → `frontend/docs/design/01-architecture.md`
- [x] 2.3 Move `frontend/docs/impl/02-components.md` → `frontend/docs/design/02-components.md`
- [x] 2.4 Move `frontend/docs/impl/03-data-flow.md` → `frontend/docs/design/03-data-flow.md`
- [x] 2.5 Move `frontend/docs/impl/04-theme.md` → `frontend/docs/design/04-theme.md`
- [x] 2.6 Move `frontend/docs/development-strategy.md` → `frontend/docs/design/development-strategy.md`
- [x] 2.7 Move `frontend/docs/common/tech-stack.md` → `frontend/docs/reference/tech-stack.md`
- [x] 2.8 Move `frontend/docs/visual-style-guide.md` → `frontend/docs/reference/visual-style-guide.md`
- [x] 2.9 Remove empty old directories: `frontend/docs/impl/`, `frontend/docs/common/`

## 3. Backend docs/ migration

- [x] 3.1 Create new directories: `backend/docs/design/`, `backend/docs/reference/`
- [x] 3.2 Move `backend/docs/impl/phase1-go-glue.md` → `backend/docs/design/phase1-go-glue.md`
- [x] 3.3 Move `backend/docs/common/tech-stack.md` → `backend/docs/reference/tech-stack.md`
- [x] 3.4 Remove empty old directories: `backend/docs/impl/`, `backend/docs/common/`

## 4. Agentend docs/ migration

- [x] 4.1 Create new directories: `agentend/docs/design/`, `agentend/docs/reference/`, `agentend/docs/guides/`, `agentend/docs/testing/`, `agentend/docs/backlog/`
- [x] 4.2 Move `agentend/docs/impl/architecture.md` → `agentend/docs/design/architecture.md`
- [x] 4.3 Move `agentend/docs/impl/01-schemas.md` → `agentend/docs/design/01-schemas.md`
- [x] 4.4 Move `agentend/docs/impl/02-adapters.md` → `agentend/docs/design/02-adapters.md`
- [x] 4.5 Move `agentend/docs/impl/03-session.md` → `agentend/docs/design/03-session.md`
- [x] 4.6 Move `agentend/docs/impl/04-rules.md` → `agentend/docs/design/04-rules.md`
- [x] 4.7 Move `agentend/docs/impl/05-api.md` → `agentend/docs/design/05-api.md`
- [x] 4.8 Move `agentend/docs/impl/06-app-wiring.md` → `agentend/docs/design/06-app-wiring.md`
- [x] 4.9 Move `agentend/docs/impl/07-session-mapping.md` → `agentend/docs/design/07-session-mapping.md`
- [x] 4.10 Move `agentend/docs/impl/08-workspace.md` → `agentend/docs/design/08-workspace.md`
- [x] 4.11 Move `agentend/docs/impl/09-cli-session-id-writeback.md` → `agentend/docs/design/09-cli-session-id-writeback.md`
- [x] 4.12 Move `agentend/docs/impl/10-workspace-branch-cleanup.md` → `agentend/docs/design/10-workspace-branch-cleanup.md`
- [x] 4.13 Move `agentend/docs/impl/11-orchestrator-planning.md` → `agentend/docs/design/11-orchestrator-planning.md`
- [x] 4.14 Move `agentend/docs/impl/skills/taskctl.md` → `agentend/docs/design/skills/taskctl.md`
- [x] 4.15 Move `agentend/docs/common/adapter-diff.md` → `agentend/docs/reference/adapter-diff.md`
- [x] 4.16 Move `agentend/docs/common/details.md` → `agentend/docs/reference/details.md`
- [x] 4.17 Move `agentend/docs/playbooks/01-session-id-writeback.md` → `agentend/docs/testing/01-session-id-writeback.md`
- [x] 4.18 Move `agentend/docs/playbooks/02-workspace-git-ops.md` → `agentend/docs/testing/02-workspace-git-ops.md`
- [x] 4.19 Move `agentend/docs/playbooks/03-taskctl-merge.md` → `agentend/docs/testing/03-taskctl-merge.md`
- [x] 4.20 Move `agentend/docs/playbooks/04-orchestrator-planning.md` → `agentend/docs/testing/04-orchestrator-planning.md`
- [x] 4.21 Move `agentend/docs/todos/orchestrator-drawbacks.md` → `agentend/docs/backlog/orchestrator-drawbacks.md`
- [x] 4.22 Move `agentend/docs/todos/session-persistence.md` → `agentend/docs/backlog/session-persistence.md`
- [x] 4.23 Move `agentend/docs/todos/system-architecture-and-frontend-cards.md` → `agentend/docs/backlog/system-architecture-and-frontend-cards.md`
- [x] 4.24 Remove empty old directories: `agentend/docs/impl/`, `agentend/docs/common/`, `agentend/docs/playbooks/`, `agentend/docs/todos/`

## 5. Update AGENTS.md links

- [x] 5.1 Update root `AGENTS.md` — fix docs paths (e.g. `docs/common/makefile-guide.md` → `docs/guides/makefile-guide.md`)
- [x] 5.2 Update `frontend/AGENTS.md` — fix docs paths (e.g. `docs/common/tech-stack.md` → `docs/reference/tech-stack.md`)
- [x] 5.3 Update `backend/AGENTS.md` — fix docs paths
- [x] 5.4 Update `agentend/AGENTS.md` — fix docs paths (e.g. `docs/common/details.md` → `docs/reference/details.md`)
- [x] 5.5 Verify `contracts/AGENTS.md` has no docs paths to update

## 6. Verification

- [x] 6.1 Global grep for old paths (`docs/common/`, `docs/impl/`, `docs/playbooks/`, `docs/todos/`, `docs/test_guide/`, `docs/architecture/`) — confirm no stale references except in preserved directories
- [x] 6.2 Verify all new directories have at least one file
- [x] 6.3 Verify no empty old directories remain
- [x] 6.4 Verify all AGENTS.md doc links resolve to existing files
