# Proposal: Cleanup — Debug logs, Router extraction, Pagination

## Intent

Eliminar deuda técnica acumulada durante el desarrollo inicial:
1. **Debug logs** en handlers de producción (`log.Println(err)` en `auth/handler.go`, `fmt.Println("authUser")`) exponen información interna en logs de producción.
2. **Router embebido en `main.go`** dificulta testing de rutas y hace el entry-point difícil de leer.
3. **`GetAll` sin paginación** carga todos los usuarios en memoria — peligroso a escala.
4. **Typo en respuesta**: `"Usuaurio eliminado"` → `"Usuario eliminado"`.

## Scope

### In Scope
- Eliminar `log.Println(err)` y `fmt.Println(...)` de handlers de auth
- Extraer rutas de `main.go` → `internal/api/router.go`
- Agregar paginación (`page`, `limit`) a `GET /api/v1/users`
- Corregir typo en mensaje de `Delete` user
- Actualizar tests existentes para reflejar los cambios

### Out of Scope
- Logging estructurado (slog/zap) — cambio separado
- Paginación en listado de archivos
- Cualquier nueva feature

## Capabilities

### New Capabilities
- `user-list-pagination`: Paginación en el endpoint `GET /api/v1/users` con query params `page` y `limit`

### Modified Capabilities
None

## Approach

- **Router**: crear `internal/api/router.go` con `func NewRouter(...) http.Handler`; `main.go` solo construye dependencias y llama `router.NewRouter(...)`.
- **Debug logs**: eliminar directamente; auth handler no necesita logging — los errores llegan a `utils.HandleError`.
- **Pagination**: agregar `PaginatedResult[T]` en `utils`, método `GetAll(page, limit int)` en repo/service, query params parseados en handler.
- **Typo**: fix en string literal.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `cmd/api/main.go` | Modified | Solo wiring; delega rutas a `NewRouter` |
| `internal/api/router.go` | New | Toda la definición de rutas |
| `internal/modules/auth/handler.go` | Modified | Eliminar debug logs |
| `internal/modules/user/handler.go` | Modified | Pagination params + fix typo |
| `internal/modules/user/service.go` | Modified | `GetAll(page, limit)` |
| `internal/modules/user/repository.go` | Modified | `GetAll(page, limit)` con LIMIT/OFFSET |
| `internal/shared/utils/response.go` | Modified | Agregar `PaginatedResult[T]` |
| `internal/modules/user/*_test.go` | Modified | Actualizar mocks/assertions |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Breaking change en `GetAll` signature (afecta tests) | Med | Actualizar tests en el mismo PR |
| Router refactor rompe middleware chain | Low | Verificar orden de middlewares |

## Rollback Plan

`git revert` del commit de cleanup. Ningún cambio de esquema de BD involucrado — rollback seguro y sin migraciones.

## Dependencies

Ninguna externa.

## Success Criteria

- [ ] `go test ./... -cover` pasa sin errores
- [ ] `go vet ./...` sin warnings
- [ ] No aparece `log.Println` ni `fmt.Println` en handlers
- [ ] `GET /api/v1/users?page=1&limit=10` responde con metadata de paginación
- [ ] `main.go` no contiene definiciones de rutas
- [ ] Typo "Usuaurio" corregido
