import { successResponse } from "@/api/utils/response";
import type { Cradle } from "@/shared/di/types";
import { toUserListResponse, toUserResponse } from "./users.mapper";
import type { IUsersController } from "./users.ports";

interface Dependencies {
	usersService: Cradle["usersService"];
}

/**
 * Factory del controlador de usuarios
 *
 * Métodos disponibles:
 * - createUser: POST /users
 * - getAllUsers: GET /users
 * - updateUser: PATCH /users/:id
 */
export const createUsersController = ({
	usersService,
}: Dependencies): IUsersController => ({
	// ===============================================================
	// POST /users - Crear usuario
	// ===============================================================
	createUser: async (req, res, next) => {
		try {
			const user = await usersService.createUser(req.body);
			return successResponse({
				res,
				data: toUserResponse(user),
				message: "Usuario creado exitosamente",
				statusCode: 201,
			});
		} catch (error) {
			next(error);
			return;
		}
	},

	// ===============================================================
	// GET /users - Obtener todos los usuarios
	// ===============================================================
	getAllUsers: async (_req, res, next) => {
		try {
			const users = await usersService.getAllUsers();
			return successResponse({
				res,
				data: toUserListResponse(users),
				message: "Usuarios obtenidos exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},

	// ===============================================================
	// PATCH /users/:id - Actualizar usuario
	// ===============================================================
	updateUser: async (req, res, next) => {
		try {
			const updatedUser = await usersService.updateUser(
				req.params.id as string,
				req.body,
			);
			return successResponse({
				res,
				data: toUserResponse(updatedUser),
				message: "Usuario actualizado exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},

	// ===============================================================
	// DELETE /users/:id - Eliminar usuario
	// ===============================================================
	deleteUser: async (req, res, next) => {
		try {
			await usersService.deleteUser(req.params.id as string);
			return successResponse({
				res,
				message: "Usuario eliminado exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},
});
