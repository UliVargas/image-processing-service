# user-list-pagination Specification

## Purpose

Define el comportamiento de paginación del endpoint `GET /api/v1/users`.
Garantiza que el sistema no cargue todos los usuarios en memoria y que
el cliente tenga información suficiente para navegar entre páginas.

## Requirements

### Requirement: Paginación por query params

El sistema MUST aceptar los parámetros `page` y `limit` como query params en
`GET /api/v1/users`. Si no se proporcionan, MUST aplicar valores por defecto:
`page=1`, `limit=10`. El valor de `limit` MUST estar restringido al rango [1, 100].

#### Scenario: Petición sin params — defaults aplicados

- GIVEN un usuario autenticado
- WHEN hace `GET /api/v1/users` sin query params
- THEN la respuesta tiene status 200
- AND devuelve como máximo 10 usuarios
- AND `meta.page = 1`, `meta.limit = 10`

#### Scenario: Petición con page y limit explícitos

- GIVEN un usuario autenticado y al menos 15 usuarios en BD
- WHEN hace `GET /api/v1/users?page=2&limit=5`
- THEN la respuesta tiene status 200
- AND devuelve los usuarios 6-10 del total
- AND `meta.page = 2`, `meta.limit = 5`, `meta.total >= 15`

#### Scenario: page fuera de rango (mayor al total de páginas)

- GIVEN un usuario autenticado y 3 usuarios en BD
- WHEN hace `GET /api/v1/users?page=99&limit=10`
- THEN la respuesta tiene status 200
- AND `data` es un array vacío
- AND `meta.total = 3`

#### Scenario: limit = 0 — rechazado

- GIVEN un usuario autenticado
- WHEN hace `GET /api/v1/users?limit=0`
- THEN la respuesta tiene status 422
- AND `error.code = "VALIDATION_FAILED"`

#### Scenario: limit > 100 — rechazado

- GIVEN un usuario autenticado
- WHEN hace `GET /api/v1/users?limit=200`
- THEN la respuesta tiene status 422
- AND `error.code = "VALIDATION_FAILED"`

#### Scenario: params no numéricos — rechazados

- GIVEN un usuario autenticado
- WHEN hace `GET /api/v1/users?page=abc&limit=xyz`
- THEN la respuesta tiene status 422
- AND `error.code = "VALIDATION_FAILED"`

### Requirement: Estructura de respuesta paginada

El sistema MUST envolver la lista de usuarios en un objeto con dos campos:
`data` (array de usuarios) y `meta` (metadatos de paginación). La respuesta
MUST incluir `meta.total` (total de registros), `meta.page` y `meta.limit`.

#### Scenario: Respuesta incluye metadata completa

- GIVEN una petición paginada exitosa
- WHEN el sistema responde
- THEN el cuerpo contiene `{ success: true, data: { data: [...], meta: { total, page, limit } } }`

### Requirement: Require autenticación

El sistema MUST rechazar peticiones sin token JWT válido.

#### Scenario: Sin token — 401

- GIVEN ningún token en la cabecera Authorization
- WHEN hace `GET /api/v1/users`
- THEN la respuesta tiene status 401
- AND `error.code = "UNAUTHORIZED"`
