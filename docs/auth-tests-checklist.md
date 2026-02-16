# Auth + Tests Milestone Checklist

Objetivo del hito: cerrar el módulo de autenticación y la integración de pruebas clave antes de iniciar upload/transformaciones.

## Scope del hito

- Módulo `auth` funcional y consistente en login/logout/renew-session.
- Cobertura de pruebas de auth y users en rutas críticas.
- Manejo de errores y validaciones alineado en API.

## Definition of Done (DoD)

- [ ] Login válido retorna access token y refresh token.
- [ ] Login inválido retorna error controlado (credenciales incorrectas).
- [ ] Logout invalida la sesión activa.
- [ ] Renew session rota token o renueva sesión correctamente.
- [ ] Middleware auth protege rutas privadas y rechaza token inválido/expirado.
- [ ] Validaciones de entrada cubren casos inválidos (body/params).
- [ ] Errores se devuelven con formato consistente de API.
- [ ] Tests de integración de auth pasan en local.
- [ ] Tests unitarios de servicios auth/session pasan en local.
- [ ] `pnpm test` sin fallos críticos del módulo auth/users.

## Pruebas mínimas sugeridas

### Integración (API)
- [ ] `POST /api/auth/login` éxito.
- [ ] `POST /api/auth/login` credenciales inválidas.
- [ ] `POST /api/auth/logout` con token válido.
- [ ] `POST /api/auth/logout` sin token.
- [ ] `POST /api/auth/renew-session` éxito.
- [ ] `POST /api/auth/renew-session` token inválido/expirado.

### Unitarias
- [ ] `auth.service` (login, logout, renewSession).
- [ ] `session.repository` (crear, invalidar, buscar sesión activa).
- [ ] `token-manager.service` (sign/verify).
- [ ] `hasher.service` (hash/compare).

## Criterios de calidad

- [ ] Convenciones de commit: `feat|fix|test|chore(scope): mensaje`.
- [ ] Sin secretos en commits.
- [ ] `.env.example` consistente con variables usadas en código.
- [ ] README no promete funcionalidades no implementadas.

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
