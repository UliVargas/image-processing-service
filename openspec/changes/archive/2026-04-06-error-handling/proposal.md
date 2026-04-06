# Proposal: Error Handling Centralization

## Intent

Los errores HTTP genéricos (`ErrInvalidJSON`, `ErrInvalidIDFormat`) están duplicados en 3 paquetes distintos con código idéntico. El patrón para producir errores de validación con `details` se repite 7+ veces. `ErrAlreadyExists` tiene el mismo código y mensaje en dos servicios. Esto hace que un cambio trivial (e.g., el mensaje de INVALID_JSON) requiera editar múltiples archivos. El objetivo es centralizar sin sobre-ingeniería: mover solo lo que es verdaderamente compartido, dejar los errores de dominio donde están.

## Scope

### In Scope
- Mover `ErrInvalidJSON`, `ErrInvalidIDFormat`, `ErrAlreadyExists` a `utils/errors.go`
- Añadir helper `utils.ValidationError(details interface{}) *AppError`
- Eliminar las definiciones duplicadas en `auth/handler.go`, `user/handler.go`, `api/middleware/auth.go`, `auth/service.go`
- Reemplazar los 7+ copy-pastes del patrón `NewError(ErrValidation.StatusCode, ...)` por `utils.ValidationError(...)`

### Out of Scope
- Errores de dominio específicos (`ErrNotFound`, `ErrInvalidPassword`, `ErrInvalidCredentials`, `ErrInvalidSession`) — quedan en sus servicios
- Wrapping de errores con `fmt.Errorf("%w")` — no aporta valor aquí
- Catálogo centralizado de todos los errores (sobre-ingeniería)
- Tests nuevos — es refactor puro, los tests existentes validan el comportamiento

## Capabilities

### New Capabilities
None

### Modified Capabilities
None — refactor puro. Sin cambios de comportamiento observable.

## Approach

Añadir 4 sentinelas + 1 helper en `utils/errors.go`. Luego borrar las definiciones duplicadas en cada handler/middleware/service y sustituir los call sites. Los tests existentes (40) sirven de red de seguridad.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/shared/utils/errors.go` | Modified | +`ErrInvalidJSON`, `ErrInvalidIDFormat`, `ErrAlreadyExists`, `ValidationError()` |
| `internal/modules/auth/handler.go` | Modified | Eliminar 3 vars duplicadas; usar `utils.*` |
| `internal/modules/auth/service.go` | Modified | Eliminar `ErrAlreadyExists`; usar `utils.ErrAlreadyExists` |
| `internal/modules/user/handler.go` | Modified | Eliminar 3 vars duplicadas; reemplazar patrón ValidationError |
| `internal/modules/file/handler.go` | Modified | Eliminar `ErrValidation`; reemplazar patrón ValidationError |
| `internal/api/middleware/auth.go` | Modified | Eliminar `ErrInvalidIDFormat`; usar `utils.ErrInvalidIDFormat` |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Romper import cycle (`utils` importa algo de modules) | Low | `utils` no importa ningún módulo — dependencia unidireccional |
| Cambio de comportamiento silencioso | Low | 40 tests existentes validan respuestas HTTP exactas |

## Rollback Plan

`git revert` del commit. Sin migraciones ni cambios de schema.

## Dependencies

Ninguna externa.

## Success Criteria

- [ ] `go build ./...` limpio
- [ ] `go vet ./...` limpio
- [ ] `go test ./...` — todos los tests existentes pasan sin modificación
- [ ] Zero instancias de `ErrInvalidJSON`, `ErrInvalidIDFormat`, `ErrAlreadyExists` definidas fuera de `utils/errors.go`
- [ ] Zero instancias del patrón `NewError(ErrValidation.StatusCode, ErrValidation.Code, ErrValidation.Message, ...)` en el código
