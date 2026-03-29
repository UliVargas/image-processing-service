# Image Processing Service

Este repositorio contiene la implementación en **Go** de un servicio backend
para el procesamiento de imágenes. La aplicación permite gestionar usuarios,
subir imágenes, aplicar transformaciones y servirlas mediante una API RESTful.

La lógica de negocio se organiza en módulos (`auth`, `user`, `session`,
`image`) y el código está desacoplado del transporte mediante interfaces. Se
utilizan PostgreSQL para persistencia y MinIO/AWS S3 para almacenamiento de
objetos.

## Getting started

1. Instalar dependencias:
   ```bash
   go mod download
   ```
2. Copiar variables de entorno y editarlas:
   ```bash
   cp .env.example .env
   ```
3. Levantar servicios auxiliares con Docker:
   ```bash
   docker-compose up -d
   ```
4. Ejecutar el servidor:
   ```bash
   go run ./cmd/api
   ```
5. Ejecutar pruebas con cobertura:
   ```bash
   go test ./... -cover
   ```

(La documentación completa de instalación y arquitectura se encuentra en
`docs/`.)

## Project structure

```text
.
├── cmd/          # ejecutables (api)
├── internal/     # código de la aplicación
├── compose.yml   # servicios auxiliares (db, minio, redis)
├── go.mod        # módulo Go (image-processing-service)
└── docs/         # documentación de la aplicación
```
> Las instrucciones específicas de uso y detalles arquitecturales viven en
> archivos dentro de `docs/`.

## Migraciones con Atlas (recomendado)

Este proyecto incluye integración con Atlas + GORM Provider para manejar
migraciones versionadas. En desarrollo, `AutoMigrate` de GORM queda desactivado
por defecto.

### Archivos clave

- `atlas.hcl`: configuración de Atlas.
- `cmd/atlas/main.go`: loader de modelos GORM para Atlas.
- `migrations/`: directorio de migraciones SQL versionadas.

### Instalación

1. Instalar Atlas CLI:
   ```bash
   curl -sSf https://atlasgo.sh | sh
   atlas version
   ```
2. Instalar dependencias del provider:
   ```bash
   go mod tidy
   ```

### Comandos principales

1. Crear una nueva migración desde modelos GORM:
   ```bash
   atlas migrate diff <nombre_cambio> --env gorm
   ```
2. Validar checksum y SQL de migraciones:
   ```bash
   atlas migrate validate --env gorm
   ```
3. Analizar riesgos en la última migración:
   ```bash
   atlas migrate lint --env gorm --latest 1
   ```
4. Aplicar migraciones sobre una base objetivo:
   ```bash
   atlas migrate apply --env gorm -u "$DATABASE_URL"
   ```
5. Ver estado de ejecución:
   ```bash
   atlas migrate status --env gorm -u "$DATABASE_URL"
   ```

### AutoMigrate (solo si lo necesitas temporalmente)

Para habilitar `AutoMigrate` de GORM durante una transición:

```bash
ENABLE_GORM_AUTOMIGRATE=true go run ./cmd/api
```

La recomendación es mantenerlo desactivado y usar únicamente Atlas.
