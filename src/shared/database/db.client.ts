import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import { ENV_CONFIG } from "@shared/config/env.config";
import { PrismaClient } from "@/generated/prisma/client";

// ===============================================================
// Configuración del cliente de base de datos (Prisma)
// ===============================================================
const connectionString = `${ENV_CONFIG.DATABASE_URL}`;

const adapter = new PrismaPg({ connectionString });
const prisma = new PrismaClient({ adapter });

export { prisma };
