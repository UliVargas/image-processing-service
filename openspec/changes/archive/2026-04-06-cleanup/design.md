# Design: Cleanup — Debug Logs, Router Extraction, Pagination, Typo

## Technical Approach

Cuatro cambios ortogonales aplicados en un solo PR. El más estructural es la
extracción del router; el más funcional es la paginación. Ninguno requiere
migración de BD.

## Architecture Decisions

| Decisión | Opción elegida | Alternativas | Rationale |
|----------|---------------|--------------|-----------|
| Lugar de `PaginatedResult` | `internal/shared/utils/response.go` | nuevo archivo `pagination.go` | Consistente con `successResponse`/`errorResponse` ya ahí |
| Signatura `GetAll` | `GetAll(page, limit int) ([]*User, int64, error)` | struct `PaginationParams` | Minimal — solo 2 params, no justifica struct propio |
| Defaults de paginación | `page=1, limit=10` en handler | en service/repo | Handler es la frontera de entrada HTTP; convención del proyecto |
| Validación de params | manual `strconv.Atoi` + check rango | `validator` struct | Params son query strings, no JSON body — sin binding automático |
| Router `NewRouter` params | recibe structs de handlers (Handler interfaces) | recibe `*chi.Mux` ya configurado | Encapsula el setup de chi, fácil de testear en el futuro |

## Data Flow

### Paginación

```
GET /api/v1/users?page=2&limit=5
        │
        ▼
handler.GetAll()
  ├─ parse & validate page/limit (defaults, rango)
  ├─ service.GetAll(page=2, limit=5)
  │       ├─ repo.GetAll(page=2, limit=5)
  │       │       └─ SELECT ... LIMIT 5 OFFSET 5  (+ COUNT(*))
  │       └─ returns ([]*User, total int64, error)
  └─ utils.Success(w, 200, PaginatedResult{Data: users, Meta: {page,limit,total}})
```

### Router Extraction

```
main.go
  ├─ wiring: DB, S3, repos, services, handlers
  └─ router.NewRouter(authMW, userHdl, authHdl, fileHdl) → http.Handler
            └─ chi setup, middlewares, rutas (sin lógica propia)
```

## File Changes

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `internal/shared/utils/response.go` | Modify | Agregar `PaginatedResult[T any]` y `PaginatedMeta` |
| `internal/modules/user/repository.go` | Modify | `GetAll(page, limit int) ([]*User, int64, error)` con LIMIT/OFFSET + COUNT |
| `internal/modules/user/service.go` | Modify | `GetAll(page, limit int) ([]*User, int64, error)` |
| `internal/modules/user/handler.go` | Modify | Parse query params, call `GetAll(page, limit)`, fix typo |
| `internal/modules/auth/handler.go` | Modify | Eliminar `log.Println(err)` y `fmt.Println("authUser", ...)` |
| `internal/api/router.go` | Create | `func NewRouter(...) http.Handler` con todo el setup chi |
| `cmd/api/main.go` | Modify | Eliminar definición de rutas; llamar `router.NewRouter(...)` |
| `internal/modules/user/service_test.go` | Modify | Actualizar mock `GetAllFn` signature y assertions |
| `internal/modules/user/repository_test.go` | Modify | Actualizar test `GetAll` para nueva signature |
| `internal/modules/user/handler_test.go` | Modify | Agregar test cases de paginación |

## Interfaces / Contracts

```go
// internal/shared/utils/response.go
type PaginatedMeta struct {
    Total int64 `json:"total"`
    Page  int   `json:"page"`
    Limit int   `json:"limit"`
}

type PaginatedResult[T any] struct {
    Data []T           `json:"data"`
    Meta PaginatedMeta `json:"meta"`
}

// internal/modules/user/repository.go
type Repository interface {
    // ...existing methods...
    GetAll(page, limit int) ([]*User, int64, error) // CHANGED
}

// internal/modules/user/service.go
type Service interface {
    // ...existing methods...
    GetAll(page, limit int) ([]*User, int64, error) // CHANGED
}

// internal/api/router.go
func NewRouter(
    authMiddleware *middleware.authMiddleware,
    authHandler   auth.Handler,
    userHandler   user.Handler,
    fileHandler   file.Handler,
) http.Handler
```

## Testing Strategy

| Capa | Qué testear | Approach |
|------|-------------|----------|
| Unit (service) | `GetAll` con página fuera de rango retorna slice vacío | mock `GetAllFn` nueva signature |
| Integration (repo) | LIMIT/OFFSET correcto, COUNT total real | SQLite in-memory, insertar N rows |
| Unit (handler) | defaults (sin params), limit>100, params no numéricos, page vacía | `httptest.NewRecorder` — patrón ya existente en `handler_test.go` |

## Migration / Rollout

No migration required. Cambio de firma de `GetAll` es interno — no hay otros
consumidores fuera del módulo `user`.

## Open Questions

Ninguna.
