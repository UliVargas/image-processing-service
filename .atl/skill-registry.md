# Skill Registry — image-processing-service
Generated: 2026-04-06

## Project Conventions

| File | Description |
|------|-------------|
| AGENTS.md / CLAUDE.md | Not found |
| openspec/config.yaml | SDD project config (strict_tdd: true, openspec mode) |

## Available User Skills

| Skill | Trigger |
|-------|---------|
| sdd-init | "sdd init", "iniciar sdd", "openspec init" |
| sdd-explore | Before committing to a change — explore/investigate |
| sdd-propose | Create a change proposal |
| sdd-spec | Write specifications with scenarios |
| sdd-design | Create technical design document |
| sdd-tasks | Break down a change into tasks |
| sdd-apply | Implement tasks from the change |
| sdd-verify | Validate implementation against specs |
| sdd-archive | Sync delta specs and archive a completed change |
| branch-pr | Create pull request (issue-first workflow) |
| issue-creation | Create GitHub issue |
| go-testing | Go testing patterns (Bubbletea TUI, teatest) |
| judgment-day | Adversarial dual review protocol |
| skill-creator | Create new AI agent skills |
| skill-registry | Update this registry |

## SDD Workflow

```
sdd-explore → sdd-propose → sdd-spec → sdd-design → sdd-tasks → sdd-apply → sdd-verify → sdd-archive
```

## Project Notes

- **Module**: image-processing-service (Go 1.25+)
- **Router**: chi/v5
- **ORM**: GORM + PostgreSQL (prod), SQLite (tests)
- **Storage**: S3-compatible (MinIO / AWS S3)
- **Auth**: JWT (golang-jwt/jwt/v5)
- **Migrations**: Atlas + GORM provider
- **Test style**: GIVEN/WHEN/THEN, mock structs with Fn fields, SQLite in-memory for repos
- **Error pattern**: `utils.NewError(status, "CODE", "message", details)` → `*AppError`
- **ID generation**: `utils.GenerateID()` (cuid2)
