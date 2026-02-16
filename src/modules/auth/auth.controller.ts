import { successResponse } from "@/api/utils/response";
import type { Cradle } from "@/shared/di/types";
import { toAuthResponse } from "./auth.mapper";
import type { IAuthController } from "./auth.ports";

interface Dependencies {
	authService: Cradle["authService"];
}

/**
 * Factory del controlador de autenticación
 *
 * Métodos disponibles:
 * - login: Iniciar sesión
 * - logout: Cerrar sesión (invalidar)
 * - renewSession: Renovar sesión (refresh tokens)
 */
export const createAuthController = ({
	authService,
}: Dependencies): IAuthController => ({
	// ===============================================================
	// POST /login - Iniciar sesión
	// ===============================================================
	login: async (req, res, next) => {
		try {
			const clientInfo = req.clientInfo || {
				ipAddress: "unknown",
				userAgent: "unknown",
				deviceInfo: "unknown",
			};
			const auth = await authService.login(req.body, clientInfo);
			return successResponse({
				res,
				data: toAuthResponse(auth),
				message: "Se inició sesión exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},

	// ===============================================================
	// POST /logout - Cerrar sesión
	// ===============================================================
	logout: async (req, res, next) => {
		try {
			const jti = req.jti as string;
			await authService.logout(jti);
			return successResponse({
				res,
				message: "Se cerró la sesión exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},

	// ===============================================================
	// POST /renew-session - Renovar sesión
	// ===============================================================
	renewSession: async (req, res, next) => {
		try {
			const { refreshToken } = req.body;
			const clientInfo = req.clientInfo || {
				ipAddress: "unknown",
				userAgent: "unknown",
				deviceInfo: "unknown",
			};
			const auth = await authService.renewSession(refreshToken, clientInfo);
			return successResponse({
				res,
				data: toAuthResponse(auth),
				message: "Se renovó la sesión exitosamente",
			});
		} catch (error) {
			next(error);
			return;
		}
	},
});
