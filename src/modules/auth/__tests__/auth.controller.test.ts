import { faker } from "@faker-js/faker";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { createAuthController } from "../auth.controller";

const createResMock = () => {
	const json = vi.fn();
	const status = vi.fn().mockReturnValue({ json });
	return { status, json };
};

describe("auth.controller", () => {
	const authServiceMock = {
		login: vi.fn(),
		logout: vi.fn(),
		renewSession: vi.fn(),
	};
	const next = vi.fn();
	let controller: ReturnType<typeof createAuthController>;

	beforeEach(() => {
		vi.clearAllMocks();
		controller = createAuthController({
			authService: authServiceMock,
		});
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	describe("success cases", () => {
		it("should login and return success response", async () => {
			const req = {
				body: {
					email: faker.internet.email(),
					password: faker.internet.password(),
				},
				clientInfo: {
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					deviceInfo: "desktop",
				},
			};
			const res = createResMock();
			const authResponse = {
				session: {
					id: faker.string.uuid(),
					expiresAt: new Date(),
					tokenHash: faker.string.hexadecimal({ length: 64 }),
					accessJti: faker.string.uuid(),
					userId: faker.string.uuid(),
					deviceInfo: "desktop",
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					lastUsedAt: null,
					isValid: true,
					createdAt: new Date(),
				},
				tokens: {
					accessToken: faker.string.alphanumeric(30),
					refreshToken: faker.string.alphanumeric(30),
				},
			};
			authServiceMock.login.mockResolvedValue(authResponse);

			await controller.login(req as never, res as never, next);

			expect(authServiceMock.login).toHaveBeenCalledWith(
				req.body,
				req.clientInfo,
			);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Se inició sesión exitosamente",
				}),
			);
		});

		it("should use unknown clientInfo when absent on login", async () => {
			const req = {
				body: {
					email: faker.internet.email(),
					password: faker.internet.password(),
				},
			};
			const res = createResMock();
			authServiceMock.login.mockResolvedValue({
				session: {
					id: faker.string.uuid(),
					expiresAt: new Date(),
					tokenHash: faker.string.hexadecimal({ length: 64 }),
					accessJti: faker.string.uuid(),
					userId: faker.string.uuid(),
					deviceInfo: "desktop",
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					lastUsedAt: null,
					isValid: true,
					createdAt: new Date(),
				},
				tokens: {
					accessToken: faker.string.alphanumeric(30),
					refreshToken: faker.string.alphanumeric(30),
				},
			});

			await controller.login(req as never, res as never, next);

			expect(authServiceMock.login).toHaveBeenCalledWith(req.body, {
				ipAddress: "unknown",
				userAgent: "unknown",
				deviceInfo: "unknown",
			});
		});

		it("should logout and return success response", async () => {
			const req = { jti: faker.string.uuid() };
			const res = createResMock();
			authServiceMock.logout.mockResolvedValue(undefined);

			await controller.logout(req as never, res as never, next);

			expect(authServiceMock.logout).toHaveBeenCalledWith(req.jti);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Se cerró la sesión exitosamente",
				}),
			);
		});

		it("should renew session and return success response", async () => {
			const req = {
				body: {
					refreshToken: faker.string.alphanumeric(24),
				},
				clientInfo: {
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					deviceInfo: "desktop",
				},
			};
			const res = createResMock();
			authServiceMock.renewSession.mockResolvedValue({
				session: {
					id: faker.string.uuid(),
					expiresAt: new Date(),
					tokenHash: faker.string.hexadecimal({ length: 64 }),
					accessJti: faker.string.uuid(),
					userId: faker.string.uuid(),
					deviceInfo: "desktop",
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					lastUsedAt: null,
					isValid: true,
					createdAt: new Date(),
				},
				tokens: {
					accessToken: faker.string.alphanumeric(30),
					refreshToken: faker.string.alphanumeric(30),
				},
			});

			await controller.renewSession(req as never, res as never, next);

			expect(authServiceMock.renewSession).toHaveBeenCalledWith(
				req.body.refreshToken,
				req.clientInfo,
			);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Se renovó la sesión exitosamente",
				}),
			);
		});
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	describe("error cases", () => {
		it("should call next with error", async () => {
			const req = { body: {} };
			const res = createResMock();
			const error = new Error("boom");
			authServiceMock.login.mockRejectedValue(error);

			await controller.login(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});

		it("should call next with error on logout", async () => {
			const req = { jti: faker.string.uuid() };
			const res = createResMock();
			const error = new Error("logout-error");
			authServiceMock.logout.mockRejectedValue(error);

			await controller.logout(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});

		it("should call next with error on renewSession", async () => {
			const req = {
				body: { refreshToken: faker.string.alphanumeric(24) },
				clientInfo: {
					ipAddress: faker.internet.ip(),
					userAgent: faker.internet.userAgent(),
					deviceInfo: "desktop",
				},
			};
			const res = createResMock();
			const error = new Error("renew-error");
			authServiceMock.renewSession.mockRejectedValue(error);

			await controller.renewSession(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});
	});
});
