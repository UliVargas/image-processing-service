import type { Cradle } from "@/shared/di/types";
import { createError } from "@/shared/errors/app-error";
import type { ClientInfo, IAuthService } from "./auth.ports";
import { getBrowser, getDeviceType } from "./auth.utils";

interface Dependencies {
	usersRepository: Cradle["usersRepository"];
	sessionRepository: Cradle["sessionRepository"];
	config: Cradle["config"];
	tokenManager: Cradle["tokenManager"];
	encryptor: Cradle["encryptor"];
	hasherService: Cradle["hasherService"];
	idGenerator: Cradle["idGenerator"];
}

/**
 * Factory del servicio de autenticación
 *
 * Casos de uso disponibles:
 * - login: Crear nuevo usuario
 * - logout: Cerrar sesión (invalidar)
 * - renewSession: Renovar sesión (refresh tokens)
 */
export const createAuthService = ({
	usersRepository,
	sessionRepository,
	tokenManager,
	config,
	encryptor,
	hasherService,
	idGenerator,
}: Dependencies): IAuthService => ({
	// ===============================================================
	// Iniciar de sesión
	// ===============================================================
	login: async (data, clientInfo: ClientInfo) => {
		const existingUser = await usersRepository.findUserByEmailAndUsername({
			email: "email" in data ? data.email : undefined,
			username: "username" in data ? data.username : undefined,
		});
		if (!existingUser) {
			throw createError({
				statusCode: 401,
				code: "UNAUTHORIZED",
				message: "Credenciales inválidas",
			});
		}
		const isValid = await encryptor.compare(
			data.password,
			existingUser.password,
		);
		if (!isValid) {
			throw createError({
				statusCode: 401,
				code: "UNAUTHORIZED",
				message: "Credenciales inválidas",
			});
		}

		const accessTokenJti = idGenerator();
		const refreshToken = idGenerator();

		const accessToken = tokenManager.sign(
			{
				sub: existingUser.id,
				jti: accessTokenJti,
			},
			config.accessTokenSecret,
			{
				expiresIn: config.accessTokenExpirationTime,
			},
		);

		const expiresAt = new Date();
		expiresAt.setDate(expiresAt.getDate() + 7);

		const session = await sessionRepository.createSession({
			id: idGenerator(),
			tokenHash: hasherService.createHash(refreshToken),
			accessJti: accessTokenJti,
			userId: existingUser.id,
			ipAddress: clientInfo.ipAddress,
			userAgent: clientInfo.userAgent,
			deviceInfo: clientInfo.deviceInfo,
			expiresAt,
		});

		return {
			session,
			tokens: {
				accessToken,
				refreshToken,
			},
		};
	},

	// ===============================================================
	// Cerrar sesión
	// ===============================================================
	logout: async (jti) => {
		await sessionRepository.deleteSession(jti);
	},

	// ===============================================================
	// Renovar  sesión
	// ===============================================================
	renewSession: async (token, clientInfo) => {
		const tokenHash = hasherService.createHash(token);
		const session = await sessionRepository.findSessionByToken(tokenHash);
		if (!session) {
			throw createError({
				statusCode: 404,
				code: "NOT_FOUND",
				message: "Sesión no encontrada",
			});
		}

		if (!session.isValid) {
			await sessionRepository.deleteSession(session.accessJti);
			throw createError({
				statusCode: 401,
				code: "SESSION_INVALID",
				message: "Sesión inválida",
			});
		}

		if (session.expiresAt.getTime() < Date.now()) {
			await sessionRepository.deleteSession(session.accessJti);
			throw createError({
				statusCode: 401,
				code: "SESSION_EXPIRED",
				message: "Sesión expirada",
			});
		}

		// Detectar cambios significativos de dispositivo o navegador
		const currentDeviceType = getDeviceType(session.userAgent || "");
		const newDeviceType = getDeviceType(clientInfo.userAgent);
		const currentBrowser = getBrowser(session.userAgent || "");
		const newBrowser = getBrowser(clientInfo.userAgent);

		const deviceTypeChanged = currentDeviceType !== newDeviceType;
		const browserChanged = currentBrowser !== newBrowser;
		const ipChanged = session.ipAddress !== clientInfo.ipAddress;

		// Invalidar si cambia el tipo de dispositivo o navegador
		if (deviceTypeChanged || browserChanged) {
			await sessionRepository.deleteSession(session.accessJti);
			throw createError({
				statusCode: 401,
				code: "SESSION_DEVICE_MISMATCH",
				message:
					"Sesión invalidada por cambio de dispositivo o navegador. Cada dispositivo/navegador debe tener su propia sesión.",
			});
		}

		// Solo logging si cambia la IP (permitido, pero monitoreado)
		if (ipChanged) {
			console.warn(
				`[SEGURIDAD] Cambio de IP detectado al renovar sesión:\n` +
					`  Usuario: ${session.userId}\n` +
					`  Dispositivo: ${currentDeviceType} (${currentBrowser})\n` +
					`  IP: ${session.ipAddress} → ${clientInfo.ipAddress}`,
			);
		}

		const accessTokenJti = idGenerator();
		const refreshToken = idGenerator();

		const accessToken = tokenManager.sign(
			{
				sub: session.userId,
				jti: accessTokenJti,
			},
			config.accessTokenSecret,
			{
				expiresIn: config.accessTokenExpirationTime,
			},
		);

		const expiresAt = new Date();
		expiresAt.setDate(expiresAt.getDate() + 7);

		const updatedSession = await sessionRepository.renewSession(session.id, {
			tokenHash: hasherService.createHash(refreshToken),
			accessJti: accessTokenJti,
			expiresAt,
			lastUsedAt: new Date(),
			ipAddress: clientInfo.ipAddress,
			userAgent: clientInfo.userAgent,
			deviceInfo: clientInfo.deviceInfo,
		});

		return {
			session: updatedSession,
			tokens: {
				accessToken,
				refreshToken,
			},
		};
	},
});
