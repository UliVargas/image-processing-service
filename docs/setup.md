# Guía de Configuración y Uso

Instrucciones básicas para ejecutar la aplicación localmente.

## Requisitos previos

- Go 1.25 o superior
- Docker y Docker Compose (para servicios auxiliares)
- PostgreSQL y MinIO (pueden levantarse con Compose)

## Variables de entorno

Copia el archivo de ejemplo y actualiza según tu entorno:

```bash
cp .env.example .env
```

Contiene valores para la base de datos, JWT, MinIO y opcionalmente Redis.

## Levantar dependencias

```bash
docker-compose up -d
```

Servirá:

- PostgreSQL en `localhost:5432`
- MinIO en `http://localhost:9000` (bucket `images`)
- Redis en `localhost:6379` (si se usa)

## Ejecutar la aplicación

```
go run ./cmd/api
```

El servidor escucha en el puerto definido en `PORT` (por defecto 3000).

## Ejecución de pruebas

```bash
go test ./... -coverprofile=coverage.out
```

Puedes especificar paquetes concretos o usar `-run` para filtrar pruebas.

## Construcción Docker

```bash
docker build -t image-processing-service .
```

y luego

```bash
docker run --env-file .env -p 3000:3000 image-processing-service
```

## Limpieza

```bash
docker-compose down
```

---

Para más detalles sobre la arquitectura y el diseño de los módulos consulta
`docs/architecture.md`.