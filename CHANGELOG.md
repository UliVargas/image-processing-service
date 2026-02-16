# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project follows Semantic Versioning.

## [0.2.0-alpha.1] - 2026-02-15

### Added
- Integration test setup isolated with Prisma + SQLite for auth flows.
- Global integration harness under `src/__tests__/setup`.
- Extended auth integration coverage for success and error scenarios.
- Conventional commit enforcement with `husky` + `commitlint`.

### Changed
- Test scripts split to run unit/API checks separately from integration.
- `test:integration` now runs in non-watch mode for CI/local consistency.
- `api` health test updated to use ephemeral port and avoid port collisions.
- `.env.example` aligned with runtime config (`SALT_ROUNDS`).

### Quality
- Auth, users, session and shared service test suites homogenized with RED/GREEN phases.
- Milestone checklist for auth + tests completed.

## [0.1.0-alpha.1] - 2026-02-15

### Added
- Initial project scaffolding with TypeScript + Express architecture.
- Auth and users modules baseline.
- Prisma/PostgreSQL setup and development scripts.
- Testing setup with Vitest and Supertest.
- Repository metadata, README, and MIT license.

### Notes
- Project status: active development (alpha).
- Next milestone: complete auth hardening and test coverage, then start image upload/transformations.
