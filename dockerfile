FROM node:24-alpine AS base

# Instalar dependencias necesarias para Sharp
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    cairo-dev \
    jpeg-dev \
    pango-dev \
    giflib-dev

WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./
COPY pnpm-lock.yaml ./

# Instalar pnpm
RUN npm install -g pnpm

FROM base AS dependencies

# Instalar dependencias de producción
RUN pnpm install --frozen-lockfile --prod

FROM base AS build

# Instalar todas las dependencias (incluidas dev)
RUN pnpm install --frozen-lockfile

# Copiar código fuente
COPY . .

# Generar cliente de Prisma
RUN pnpm prisma generate

# Compilar TypeScript
RUN pnpm build

FROM base AS production

# Copiar dependencias de producción
COPY --from=dependencies /app/node_modules ./node_modules

# Copiar código compilado
COPY --from=build /app/dist ./dist
COPY --from=build /app/prisma ./prisma

# Copiar archivos necesarios
COPY package*.json ./

# Exponer puerto
EXPOSE 3000

# Comando para ejecutar
CMD ["node", "dist/index.js"]