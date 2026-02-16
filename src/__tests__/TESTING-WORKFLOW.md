# Testing Workflow (proyecto)

Flujo recomendado al desarrollar una feature:

1. Escribir primero un test del comportamiento de negocio clave que actualmente falla.
2. Implementar el mínimo código para hacerlo pasar.
3. Refactorizar sin cambiar comportamiento.

## Naming recomendado para tests

- `should <resultado esperado> when <condición>`
- Nombres orientados al comportamiento, no a detalles de implementación.

## Orden sugerido por capa

1. Unit tests de servicios (reglas de negocio)
2. Tests de controladores/repositorios (contratos e interacción)
3. Tests de integración (flujo end-to-end)

## Ejemplo aplicado en este repo

- `src/modules/users/__tests__/users.service.test.ts`

El archivo está estructurado por caso de uso (`createUser`, `updateUser`, etc.) y mantiene una secuencia práctica de implementación guiada por tests.
