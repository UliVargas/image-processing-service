# Tasks: Cleanup — Debug Logs, Router Extraction, Pagination, Typo

## Phase 1: Foundation

- [x] 1.1 `internal/shared/utils/response.go` — agregar `PaginatedMeta` struct (Total int64, Page int, Limit int)
- [x] 1.2 `internal/shared/utils/response.go` — agregar `PaginatedResult[T any]` struct (Data []T, Meta PaginatedMeta)
- [x] 1.3 `internal/modules/user/repository.go` — actualizar interfaz `Repository`: `GetAll(page, limit int) ([]*User, int64, error)`
- [x] 1.4 `internal/modules/user/service.go` — actualizar interfaz `Service`: `GetAll(page, limit int) ([]*User, int64, error)`

## Phase 2: Core Implementation

- [x] 2.1 `internal/modules/user/repository.go` — implementar `GetAll` con `LIMIT/OFFSET` y `COUNT(*)` en una transacción
- [x] 2.2 `internal/modules/user/service.go` — implementar `GetAll(page, limit)` delegando a repo y retornando `([]*User, int64, error)`
- [x] 2.3 `internal/modules/user/handler.go` — parsear `page`/`limit` con `strconv.Atoi`, aplicar defaults (1/10), validar rango limit [1,100]
- [x] 2.4 `internal/modules/user/handler.go` — responder con `PaginatedResult[*User]` usando `utils.Success`
- [x] 2.5 `internal/modules/user/handler.go` — corregir typo `"Usuaurio eliminado"` → `"Usuario eliminado"`
- [x] 2.6 `internal/modules/auth/handler.go` — eliminar `log.Println(err)` del método `SignUp`
- [x] 2.7 `internal/modules/auth/handler.go` — eliminar `fmt.Println("authUser", authUser)` del método `SignOut`

## Phase 3: Wiring

- [x] 3.1 Crear `internal/api/router.go` — función `NewRouter(authMW, authHdl, userHdl, fileHdl) http.Handler` con todo el setup de chi (middlewares + rutas)
- [x] 3.2 `cmd/api/main.go` — eliminar definición de rutas; reemplazar por llamada a `router.NewRouter(...)`

## Phase 4: Testing

- [x] 4.1 [RED] `user/service_test.go` — actualizar `GetAllFn` del mock a nueva signature `(page, limit int) ([]*User, int64, error)`
- [x] 4.2 [RED] `user/repository_test.go` — actualizar `TestRepository_GetAll` para nueva signature; agregar caso LIMIT/OFFSET correcto
- [x] 4.3 [RED] `user/handler_test.go` — agregar caso: sin params → defaults (page=1, limit=10), verifica `meta` en respuesta
- [x] 4.4 [RED] `user/handler_test.go` — agregar caso: `?page=2&limit=5` → segundo bloque de resultados
- [x] 4.5 [RED] `user/handler_test.go` — agregar caso: `?limit=0` → 422 `VALIDATION_FAILED`
- [x] 4.6 [RED] `user/handler_test.go` — agregar caso: `?limit=200` → 422 `VALIDATION_FAILED`
- [x] 4.7 [RED] `user/handler_test.go` — agregar caso: `?page=abc&limit=xyz` → 422 `VALIDATION_FAILED`
- [x] 4.8 [GREEN] Hacer pasar todos los tests rojos (`go test ./... -cover`)
- [x] 4.9 Ejecutar `go vet ./...` — cero warnings
