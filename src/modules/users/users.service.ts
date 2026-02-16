import type { Cradle } from "@/shared/di/types";
import { createError } from "@/shared/errors/app-error";
import type { IUsersService } from "./users.ports";

interface Dependencies {
	usersRepository: Cradle["usersRepository"];
	idGenerator: Cradle["idGenerator"];
	encryptor: Cradle["encryptor"];
	config: Cradle["config"];
}

/**
 * Factory del servicio de usuarios
 *
 * Casos de uso disponibles:
 * - createUser: Crear nuevo usuario
 * - getAllUsers: Obtener listado de usuarios
 * - updateUser: Actualizar datos de usuario
 * - deleteUser: Eliminar usuario
 */
export const createUsersService = ({
	usersRepository,
	idGenerator,
	encryptor,
	config,
}: Dependencies): IUsersService => ({
	// ===============================================================
	// Crear usuario
	// ===============================================================
	createUser: async (data) => {
		const existingUser = await usersRepository.findUserByEmailAndUsername({
			email: data.email,
			username: data.username,
		});
		if (existingUser) {
			throw createError({
				message:
					"No se puede crear el usuario con el correo electrónico o nombre de usuario proporcionados",
				statusCode: 409,
				code: "USER_ALREADY_EXISTS",
			});
		}
		const passwordHashed = await encryptor.hash(
			data.password,
			config.saltRounds,
		);
		return await usersRepository.createUser({
			id: idGenerator(),
			email: data.email,
			username: data.username,
			password: passwordHashed,
			name: data.name,
		});
	},

	// ===============================================================
	// Obtener todos los usuarios
	// ===============================================================
	getAllUsers: async () => {
		return await usersRepository.getAllUsers();
	},

	// ===============================================================
	// Actualizar usuario
	// ===============================================================
	updateUser: async (id, data) => {
		const existingUser = await usersRepository.findUserById(id);
		if (!existingUser) {
			throw createError({
				message: "Usuario no encontrado",
				code: "USER_NOT_FOUND",
				statusCode: 404,
			});
		}
		return await usersRepository.updateUser(id, data);
	},

	// ===============================================================
	// Eliminar usuario
	// ===============================================================
	deleteUser: async (id) => {
		return await usersRepository.deleteUser(id);
	},
});
