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

## Module name

El módulo Go se llama `image-processing-service`; todas las importaciones
internas ya están actualizadas a este nombre.
