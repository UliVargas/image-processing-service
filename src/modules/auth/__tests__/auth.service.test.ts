import { faker } from "@faker-js/faker";
import type bcrypt from "bcryptjs";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { Cradle } from "@/shared/di/types";
import { createAuthService } from "../auth.service";

describe("auth.service", () => {
	const usersRepositoryMock = {
		findUserByEmailAndUsername: vi.fn(),
	};
	const sessionRepositoryMock = {
		createSession: vi.fn(),
		deleteSession: vi.fn(),
		findSessionByToken: vi.fn(),
		renewSession: vi.fn(),
	};
	const tokenManagerMock = {
		sign: vi.fn(),
		verify: vi.fn(),
	};
	const hasherServiceMock = {
		createHash: vi.fn(),
		validarHash: vi.fn(),
	};
	const compareMock = vi.fn<(...args: [string, string]) => Promise<boolean>>();
	const encryptorMock = {
		compare: compareMock,
	} as unknown as typeof bcrypt;
	let idSequence: string[];
	const idGeneratorMock = vi.fn(
		() => idSequence.shift() || faker.string.uuid(),
	);
	const configMock = {
		accessTokenSecret: "access-secret",
		refreshTokenSecret: "refresh-secret",
		accessTokenExpirationTime: "15m",
		refreshTokenExpirationTime: "7d",
		saltRounds: 10,
	} satisfies Cradle["config"];

	const clientInfo = {
		ipAddress: faker.internet.ip(),
		userAgent: "Mozilla/5.0 Chrome/122.0.0.0 Safari/537.36",
		deviceInfo: "desktop",
	};

	const createService = () =>
		createAuthService({
			usersRepository: usersRepositoryMock as never,
			sessionRepository: sessionRepositoryMock as never,
			tokenManager: tokenManagerMock,
			config: configMock,
			encryptor: encryptorMock,
			hasherService: hasherServiceMock,
			idGenerator: idGeneratorMock,
		});

	beforeEach(() => {
		vi.clearAllMocks();
		idSequence = ["access-jti", "refresh-token", "session-id"];
		tokenManagerMock.sign.mockReturnValue("access-token-value");
		hasherServiceMock.createHash.mockImplementation((v: string) => `hash-${v}`);
		hasherServiceMock.validarHash.mockReturnValue(true);
		sessionRepositoryMock.createSession.mockImplementation(async (data) => ({
			id: data.id,
			tokenHash: data.tokenHash,
			accessJti: data.accessJti,
			userId: data.userId,
			deviceInfo: data.deviceInfo,
			ipAddress: data.ipAddress,
			userAgent: data.userAgent,
			expiresAt: data.expiresAt,
			lastUsedAt: null,
			isValid: true,
			createdAt: new Date(),
		}));
	});

	describe("login", () => {
		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================
		it("should login successfully", async () => {
			const user = {
				id: faker.string.uuid(),
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: "stored-hash",
				name: faker.person.fullName(),
				createdAt: new Date(),
				updatedAt: new Date(),
			};
			usersRepositoryMock.findUserByEmailAndUsername.mockResolvedValue(user);
			compareMock.mockResolvedValue(true);

			const service = createService();
			const result = await service.login(
				{ email: user.email, password: faker.internet.password() },
				clientInfo,
			);

			expect(usersRepositoryMock.findUserByEmailAndUsername).toHaveBeenCalled();
			expect(tokenManagerMock.sign).toHaveBeenCalledWith(
				{ sub: user.id, jti: "access-jti" },
				configMock.accessTokenSecret,
				expect.objectContaining({ expiresIn: "15m" }),
			);
			expect(sessionRepositoryMock.createSession).toHaveBeenCalledWith(
				expect.objectContaining({
					id: "session-id",
					accessJti: "access-jti",
					tokenHash: "hash-refresh-token",
					userId: user.id,
				}),
			);
			expect(result.tokens).toEqual({
				accessToken: "access-token-value",
				refreshToken: "refresh-token",
			});
		});

		// ===============================================================
		// Fase RED (reglas y errores)
		// ===============================================================
		it("should throw UNAUTHORIZED when user is not found", async () => {
			usersRepositoryMock.findUserByEmailAndUsername.mockResolvedValue(null);
			const service = createService();

			await expect(
				service.login(
					{
						email: faker.internet.email(),
						password: faker.internet.password(),
					},
					clientInfo,
				),
			).rejects.toMatchObject({
				statusCode: 401,
				code: "UNAUTHORIZED",
			});
		});

		it("should throw UNAUTHORIZED when password is invalid", async () => {
			usersRepositoryMock.findUserByEmailAndUsername.mockResolvedValue({
				id: faker.string.uuid(),
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: "stored-hash",
				name: faker.person.fullName(),
				createdAt: new Date(),
				updatedAt: new Date(),
			});
			compareMock.mockResolvedValue(false);
			const service = createService();

			await expect(
				service.login(
					{
						username: faker.internet.username(),
						password: faker.internet.password(),
					},
					clientInfo,
				),
			).rejects.toMatchObject({
				statusCode: 401,
				code: "UNAUTHORIZED",
			});
		});
	});

	describe("logout", () => {
		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================
		it("should delegate deleteSession", async () => {
			const service = createService();
			sessionRepositoryMock.deleteSession.mockResolvedValue(undefined);
			await expect(service.logout("jti-1")).resolves.toBeUndefined();
			expect(sessionRepositoryMock.deleteSession).toHaveBeenCalledWith("jti-1");
		});
	});

	describe("renewSession", () => {
		// ===============================================================
		// Fase RED (reglas y errores)
		// ===============================================================
		it("should throw NOT_FOUND when session does not exist", async () => {
			hasherServiceMock.createHash.mockReturnValue("hashed-token");
			sessionRepositoryMock.findSessionByToken.mockResolvedValue(null);
			const service = createService();

			await expect(
				service.renewSession("refresh-token", clientInfo),
			).rejects.toMatchObject({ statusCode: 404, code: "NOT_FOUND" });
		});

		it("should throw SESSION_INVALID when session is invalid", async () => {
			hasherServiceMock.createHash.mockReturnValue("hashed-token");
			sessionRepositoryMock.findSessionByToken.mockResolvedValue({
				id: faker.string.uuid(),
				tokenHash: "hashed-token",
				accessJti: "old-jti",
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: clientInfo.ipAddress,
				userAgent: clientInfo.userAgent,
				expiresAt: new Date(Date.now() + 10_000),
				lastUsedAt: null,
				isValid: false,
				createdAt: new Date(),
			});
			const service = createService();

			await expect(
				service.renewSession("refresh-token", clientInfo),
			).rejects.toMatchObject({ statusCode: 401, code: "SESSION_INVALID" });
			expect(sessionRepositoryMock.deleteSession).toHaveBeenCalledWith(
				"old-jti",
			);
		});

		it("should throw SESSION_EXPIRED when session is expired", async () => {
			hasherServiceMock.createHash.mockReturnValue("hashed-token");
			sessionRepositoryMock.findSessionByToken.mockResolvedValue({
				id: faker.string.uuid(),
				tokenHash: "hashed-token",
				accessJti: "old-jti",
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: clientInfo.ipAddress,
				userAgent: clientInfo.userAgent,
				expiresAt: new Date(Date.now() - 10_000),
				lastUsedAt: null,
				isValid: true,
				createdAt: new Date(),
			});
			const service = createService();

			await expect(
				service.renewSession("refresh-token", clientInfo),
			).rejects.toMatchObject({ statusCode: 401, code: "SESSION_EXPIRED" });
		});

		it("should throw SESSION_DEVICE_MISMATCH when browser or device changes", async () => {
			hasherServiceMock.createHash.mockReturnValue("hashed-token");
			sessionRepositoryMock.findSessionByToken.mockResolvedValue({
				id: faker.string.uuid(),
				tokenHash: "hashed-token",
				accessJti: "old-jti",
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: clientInfo.ipAddress,
				userAgent: "Mozilla/5.0 Firefox/123.0",
				expiresAt: new Date(Date.now() + 10_000),
				lastUsedAt: null,
				isValid: true,
				createdAt: new Date(),
			});
			const service = createService();

			await expect(
				service.renewSession("refresh-token", clientInfo),
			).rejects.toMatchObject({
				statusCode: 401,
				code: "SESSION_DEVICE_MISMATCH",
			});
		});

		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================
		it("should renew session successfully", async () => {
			idSequence = ["new-access-jti", "new-refresh-token"];
			hasherServiceMock.createHash.mockImplementation(
				(v: string) => `hash-${v}`,
			);
			sessionRepositoryMock.findSessionByToken.mockResolvedValue({
				id: "session-1",
				tokenHash: "hash-refresh-token",
				accessJti: "old-jti",
				userId: "user-1",
				deviceInfo: "desktop",
				ipAddress: clientInfo.ipAddress,
				userAgent: clientInfo.userAgent,
				expiresAt: new Date(Date.now() + 10_000),
				lastUsedAt: null,
				isValid: true,
				createdAt: new Date(),
			});
			sessionRepositoryMock.renewSession.mockImplementation(
				async (id, data) => ({
					id,
					tokenHash: data.tokenHash,
					accessJti: data.accessJti,
					userId: "user-1",
					deviceInfo: data.deviceInfo,
					ipAddress: data.ipAddress,
					userAgent: data.userAgent,
					expiresAt: data.expiresAt,
					lastUsedAt: data.lastUsedAt,
					isValid: true,
					createdAt: new Date(),
				}),
			);
			const service = createService();

			const result = await service.renewSession("refresh-token", clientInfo);

			expect(sessionRepositoryMock.renewSession).toHaveBeenCalledWith(
				"session-1",
				expect.objectContaining({
					accessJti: "new-access-jti",
					tokenHash: "hash-new-refresh-token",
				}),
			);
			expect(result.tokens).toEqual({
				accessToken: "access-token-value",
				refreshToken: "new-refresh-token",
			});
		});
	});
});
