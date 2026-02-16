import type { Session as PrismaSession } from "@/generated/prisma/client";
import type { Cradle } from "@/shared/di/types";
import type { Session } from "./session.entity";
import type { ISessionRepository } from "./session.ports";

interface Dependencies {
	dbClient: Cradle["dbClient"];
}

// ===============================================================
// Mapper: Prisma → Dominio
// ===============================================================
const toDomainSession = (prismasession: PrismaSession): Session => ({
	id: prismasession.id,
	tokenHash: prismasession.tokenHash,
	accessJti: prismasession.accessJti,
	userId: prismasession.userId,
	deviceInfo: prismasession.deviceInfo,
	ipAddress: prismasession.ipAddress,
	userAgent: prismasession.userAgent,
	expiresAt: prismasession.expiresAt,
	lastUsedAt: prismasession.lastUsedAt,
	isValid: prismasession.isValid,
	createdAt: prismasession.createdAt,
});

/**
 * Factory del repositorio de sesiones
 *
 * Métodos disponibles:
 * - createSession: Insertar sesión en BD
 * - renewSession: Renovar sesión en BD
 */
export const createSessionRepository = ({
	dbClient,
}: Dependencies): ISessionRepository => ({
	// ===============================================================
	// Insertar sesión en BD
	// ===============================================================
	createSession: async (data) => {
		const prismaSession = await dbClient.session.create({
			data,
		});
		return toDomainSession(prismaSession);
	},
	renewSession: async (id, data) => {
		const prismaSession = await dbClient.session.update({
			where: { id },
			data: {
				expiresAt: data.expiresAt,
				accessJti: data.accessJti,
				tokenHash: data.tokenHash,
				lastUsedAt: data.lastUsedAt,
				ipAddress: data.ipAddress,
				userAgent: data.userAgent,
				deviceInfo: data.deviceInfo,
			},
		});
		return toDomainSession(prismaSession);
	},
	findSessionByToken: async (token) => {
		const prismaSession = await dbClient.session.findFirst({
			where: { tokenHash: token },
		});
		return prismaSession ? toDomainSession(prismaSession) : null;
	},
	deleteSession: async (jti) => {
		await dbClient.session.deleteMany({
			where: { accessJti: jti },
		});
	},
});
