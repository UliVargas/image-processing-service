# Arquitectura de la Aplicación

El servicio sigue una versión adaptada de Clean Architecture en Go. El
código principal reside bajo `internal/` y está organizado en capas con
responsabilidades claras.

## Capas

```
┌───────────────────────────────────┐
│  API Layer (cmd/api)              │
│  - chi router                     │
│  - middlewares globales           │
│  - handlers / controllers         │
└───────────────────────────────────┘
               ↓
┌───────────────────────────────────┐
│  Modules (internal/modules)       │
│  - auth, user, session, image     │
│  - servicios / lógica de negocio  │
│  - repositorios (GORM interfaces) │
└───────────────────────────────────┘
               ↓
┌───────────────────────────────────┐
│  Shared (internal/shared)         │
│  - configuración (env loader)     │
│  - cliente de base de datos       │
│  - utilidades de autenticación    │
│  - helpers, errores comunes       │
└───────────────────────────────────┘
```

Cada módulo exporta interfaces que permiten sustituir implementaciones
en pruebas. La inyección de dependencias se realiza manualmente o con un
container ligero según conveniencia.

### Principios aplicados

- **SRP**: cada paquete tiene una única responsabilidad clara.
- **DIP**: el código de alto nivel depende de abstracciones, no de
  implementaciones concretas.
- **ISP**: las interfaces son pequeñas y específicas.
- **OCP/LSP**: las estructuras pueden ser extendidas sin modificar el código
  existente y las subclases (si hubiese) son intercambiables.

## Tecnologías principales

- Go 1.25+
- chi (routing & middleware)
- GORM + PostgreSQL
- github.com/golang-jwt/jwt/v5
- MinIO / AWS S3
- Redis (opcional)
- Docker / Docker Compose

## Estructura de carpetas

```text
.
├── cmd/
│   └── api/              # entry point (main.go) y configuración HTTP
├── internal/
│   ├── api/              # middleware y router.go
│   ├── modules/          # dominios: auth, user, session, image
│   └── shared/           # config, database, auth, utils
├── compose.yml           # servicios auxiliares
├── go.mod
└── docs/                 # documentación de la aplicación
```

Documentos adicionales (setup, variables de entorno, etc.) se encuentran
en el directorio `docs/`.