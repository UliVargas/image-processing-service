# 🖼️ Image Processing Service

Servicio backend construido con Node.js + TypeScript para autenticación, gestión de usuarios y base de un flujo de procesamiento de imágenes.

## 📌 Estado del proyecto

Actualmente el proyecto cubre autenticación y usuarios en producción local, con base técnica preparada para extender el dominio de imágenes y transformaciones.

## ✅ Características implementadas

- Autenticación con JWT (login, logout y renovación de sesión)
- Gestión de usuarios con validación de entrada
- API REST con Express 5 y arquitectura modular
- Persistencia con PostgreSQL + Prisma
- Inyección de dependencias con Awilix
- Middlewares de seguridad (`helmet`, `cors`, `compression`)
- Testing con Vitest y Supertest
- Infraestructura local con Docker Compose (PostgreSQL, MinIO y Redis)

## 🧱 Stack tecnológico

### Backend
- Node.js
- TypeScript
- Express

### Datos
- PostgreSQL
- Prisma ORM

### Seguridad y validación
- JWT
- bcryptjs
- Valibot

### Calidad y pruebas
- Vitest
- Supertest

### Infra local
- Docker Compose
- MinIO
- Redis

## 📁 Estructura principal

```text
src/
├── api/                  # Rutas y middlewares HTTP
├── modules/
│   ├── auth/             # Casos de uso de autenticación
│   ├── users/            # Casos de uso de usuarios
│   ├── images/           # Dominio de imágenes (en progreso)
│   └── transformations/  # Dominio de transformaciones (en progreso)
├── shared/               # Config, errores, DB y servicios transversales
├── app.ts                # Setup de Express
├── container.ts          # Registro de dependencias
└── index.ts              # Entrada de la aplicación
```

## 🚀 Inicio rápido

### Requisitos

- Node.js 20+
- pnpm
- Docker + Docker Compose

### 1) Instalar dependencias

```bash
pnpm install
```

### 2) Configurar entorno

```bash
cp .env.example .env
```

### 3) Levantar infraestructura

```bash
docker compose up -d
```

### 4) Ejecutar migraciones

```bash
pnpm prisma:migrate
```

### 5) Iniciar el servidor

```bash
pnpm dev
```

Health check: `GET /health`

## 🧪 Scripts

```bash
pnpm dev
pnpm build
pnpm start
pnpm test
pnpm test:watch
pnpm test:coverage
```

## 🗺️ Roadmap corto

- Finalizar módulo de imágenes
- Implementar transformaciones con `sharp`
- Agregar cobertura de integración para nuevos flujos
- Consolidar documentación técnica pública

## 📄 Licencia

MIT
