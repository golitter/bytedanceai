## 1. taskctl CLI (Go)

- [x] 1.1 Create `src/skills/builtin/taskctl/go.mod` and `main.go` skeleton
- [x] 1.2 Implement path self-bootstrapping: parse exe path to resolve task_id, agent_name, shared_dir
- [x] 1.3 Implement `help` command: print all commands with descriptions
- [x] 1.4 Implement `ls` command: recursively list `shared/.agent/` file structure
- [x] 1.5 Implement `summary` command: read config.yaml + plans/ files and output concatenated content
- [x] 1.6 Implement `common-memory` command: read all files in `memory/common/` sorted by name
- [x] 1.7 Implement `sub-memory` command: read files in `memory/{agent_name}/` only, enforcing agent isolation
- [x] 1.8 Compile Go binary as `exe` and commit to `src/skills/builtin/taskctl/exe`

## 2. skill.md Instruction File

- [x] 2.1 Create `src/skills/builtin/taskctl/skill.md` with usage instructions for all taskctl commands (help, ls, summary, common-memory, sub-memory), referencing exe path relative to agent skill directory

## 3. Skill Provisioner (Python)

- [x] 3.1 Create `src/skills/__init__.py` and `src/skills/provisioner.py` with `SkillProvisioner` class
- [x] 3.2 Implement `provision(worktree_path, task_id, agent_name)`: determine target skill dir (`.claude/skills/` or `.opencode/skills/`), copy builtin skill files (exe + skill.md) to target
- [x] 3.3 Implement `init_shared_dirs(worktrees_root, task_id, agent_name)`: create `shared/.agent/memory/common/` and `shared/.agent/memory/{agent_name}/` if not exist

## 4. Workspace Manager Integration

- [x] 4.1 Modify `WorkspaceManager.create()` to call `SkillProvisioner.provision()` after worktree creation
- [x] 4.2 Modify `WorkspaceManager.create()` to call `SkillProvisioner.init_shared_dirs()` after worktree creation
- [x] 4.3 Ensure provisioner is not called if worktree creation fails

## 5. Testing

- [x] 5.1 Add unit tests for taskctl CLI (Go tests): path resolution, each command, agent isolation
- [x] 5.2 Add unit tests for SkillProvisioner: file copying, shared directory creation, idempotency on second agent
- [x] 5.3 Add integration test: full workspace create → verify skill files exist in agent dir → verify shared dirs exist → run exe commands
