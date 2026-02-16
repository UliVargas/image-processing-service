# Auth + Tests Milestone Checklist

Objetivo del hito: cerrar el módulo de autenticación y la integración de pruebas clave antes de iniciar upload/transformaciones.

## Scope del hito

- Módulo `auth` funcional y consistente en login/logout/renew-session.
- Cobertura de pruebas de auth y users en rutas críticas.
- Manejo de errores y validaciones alineado en API.

## Definition of Done (DoD)

- [x] Login válido retorna access token y refresh token.
- [x] Login inválido retorna error controlado (credenciales incorrectas).
- [x] Logout invalida la sesión activa.
- [x] Renew session rota token o renueva sesión correctamente.
- [x] Middleware auth protege rutas privadas y rechaza token inválido/expirado.
- [x] Validaciones de entrada cubren casos inválidos (body/params).
- [x] Errores se devuelven con formato consistente de API.
- [x] Tests de integración de auth pasan en local.
- [x] Tests unitarios de servicios auth/session pasan en local.
- [x] `pnpm test` sin fallos críticos del módulo auth/users.

## Pruebas mínimas sugeridas

### Integración (API)
- [x] `POST /api/auth/login` éxito.
- [x] `POST /api/auth/login` credenciales inválidas.
- [x] `POST /api/auth/logout` con token válido.
- [x] `POST /api/auth/logout` sin token.
- [x] `POST /api/auth/renew-session` éxito.
- [x] `POST /api/auth/renew-session` token inválido/expirado.

### Unitarias
- [x] `auth.service` (login, logout, renewSession).
- [x] `session.repository` (crear, invalidar, buscar sesión activa).
- [x] `token-manager.service` (sign/verify).
- [x] `hasher.service` (hash/compare).

## Criterios de calidad

- [x] Convenciones de commit: `feat|fix|test|chore(scope): mensaje`.
- [x] Sin secretos en commits.
- [x] `.env.example` consistente con variables usadas en código.
- [x] README no promete funcionalidades no implementadas.

## Cierre de versión para este hito

Cuando todos los checks estén en verde:

1. Actualizar `CHANGELOG.md` con cambios del hito.
2. Subir versión en `package.json` a `0.2.0-alpha.1`.
3. Commit de release:
   - `chore(release): prepare v0.2.0-alpha.1`
4. Crear tag anotado:
   - `git tag -a v0.2.0-alpha.1 -m "auth module completed + core tests integrated"`
5. Publicar:
   - `git push origin feat/auth-tests-completion`
   - Merge a `main`
   - `git push origin main --follow-tags`
