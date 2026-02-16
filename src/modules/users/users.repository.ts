import type { User as PrismaUser } from "@/generated/prisma/client";
import type { Cradle } from "@/shared/di/types";
import type { User } from "./users.entity";
import type { IUsersRepository } from "./users.ports";

interface Dependencies {
	dbClient: Cradle["dbClient"];
}

// ===============================================================
// Mapper: Prisma → Dominio
// ===============================================================
const toDomainUser = (prismaUser: PrismaUser): User => ({
	id: prismaUser.id,
	email: prismaUser.email,
	username: prismaUser.username,
	name: prismaUser.name,
	password: prismaUser.password,
	createdAt: prismaUser.createdAt,
	updatedAt: prismaUser.updatedAt,
});

/**
 * Factory del repositorio de usuarios
 *
 * Métodos disponibles:
 * - createUser: Insertar usuario en BD
 * - findUserByEmailAndUsername: Buscar por email o username
 * - getAllUsers: Obtener todos los usuarios
 * - updateUser: Actualizar usuario por ID
 * - deleteUser: Eliminar usuario por ID
 * - findUserById: Buscar usuario por ID
 */
export const createUsersRepository = ({
	dbClient,
}: Dependencies): IUsersRepository => ({
	// ===============================================================
	// Insertar usuario en BD
	// ===============================================================
	createUser: async (data) => {
		const prismaUser = await dbClient.user.create({
			data,
		});
		return toDomainUser(prismaUser);
	},

	// ===============================================================
	// Buscar usuario por email o username
	// ===============================================================
	findUserByEmailAndUsername: async ({ email, username }) => {
		const conditions = [email && { email }, username && { username }].filter(
			(c): c is { email: string } | { username: string } => Boolean(c),
		);

		if (!conditions.length) return null;

		const prismaUser = await dbClient.user.findFirst({
			where: {
				OR: conditions,
			},
		});
		return prismaUser ? toDomainUser(prismaUser) : null;
	},

	// ===============================================================
	// Buscar usuario por ID
	// ===============================================================
	findUserById: async (id) => {
		const prismaUser = await dbClient.user.findUnique({
			where: {
				id,
			},
		});
		return prismaUser ? toDomainUser(prismaUser) : null;
	},

	// ===============================================================
	// Obtener todos los usuarios
	// ===============================================================
	getAllUsers: async () => {
		const prismaUsers = await dbClient.user.findMany();
		return prismaUsers.map(toDomainUser);
	},

	// ===============================================================
	// Actualizar usuario por ID
	// ===============================================================
	updateUser: async (id, data) => {
		const prismaUser = await dbClient.user.update({
			where: { id },
			data,
		});
		return toDomainUser(prismaUser);
	},

	// ===============================================================
	// Eliminar usuario por ID
	// ===============================================================
	deleteUser: async (id) => {
		await dbClient.user.delete({
			where: { id },
		});
	},
});
