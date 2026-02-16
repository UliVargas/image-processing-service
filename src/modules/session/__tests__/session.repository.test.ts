import { faker } from "@faker-js/faker";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { PrismaClient } from "@/generated/prisma/client";
import { createSessionRepository } from "../session.repository";

describe("session.repository", () => {
	let prismaMock: PrismaClient;
	let sessionRepository: ReturnType<typeof createSessionRepository>;

	beforeEach(() => {
		prismaMock = {
			session: {
				create: vi.fn(),
				update: vi.fn(),
				findFirst: vi.fn(),
				deleteMany: vi.fn(),
			},
		} as unknown as PrismaClient;

		sessionRepository = createSessionRepository({
			dbClient: prismaMock,
		});
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should create session", async () => {
		const sessionData = {
			id: faker.string.uuid(),
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			userId: faker.string.uuid(),
			deviceInfo: "desktop",
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			expiresAt: new Date(Date.now() + 60_000),
		};

		const prismaSession = {
			...sessionData,
			lastUsedAt: null,
			isValid: true,
			createdAt: new Date(),
		};

		vi.mocked(prismaMock.session.create).mockResolvedValue(prismaSession);

		const result = await sessionRepository.createSession(sessionData);

		expect(prismaMock.session.create).toHaveBeenCalledWith({
			data: sessionData,
		});
		expect(result).toEqual(prismaSession);
	});

	it("should renew session", async () => {
		const sessionId = faker.string.uuid();
		const renewData = {
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			expiresAt: new Date(Date.now() + 120_000),
			lastUsedAt: new Date(),
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			deviceInfo: "desktop",
		};

		const prismaSession = {
			id: sessionId,
			...renewData,
			userId: faker.string.uuid(),
			isValid: true,
			createdAt: new Date(),
		};

		vi.mocked(prismaMock.session.update).mockResolvedValue(prismaSession);

		const result = await sessionRepository.renewSession(sessionId, renewData);

		expect(prismaMock.session.update).toHaveBeenCalledWith({
			where: { id: sessionId },
			data: {
				expiresAt: renewData.expiresAt,
				accessJti: renewData.accessJti,
				tokenHash: renewData.tokenHash,
				lastUsedAt: renewData.lastUsedAt,
				ipAddress: renewData.ipAddress,
				userAgent: renewData.userAgent,
				deviceInfo: renewData.deviceInfo,
			},
		});
		expect(result).toEqual(prismaSession);
	});

	it("should find session by token", async () => {
		const tokenHash = faker.string.hexadecimal({ length: 64 });
		const prismaSession = {
			id: faker.string.uuid(),
			tokenHash,
			accessJti: faker.string.uuid(),
			userId: faker.string.uuid(),
			deviceInfo: "desktop",
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			expiresAt: new Date(Date.now() + 120_000),
			lastUsedAt: null,
			isValid: true,
			createdAt: new Date(),
		};

		vi.mocked(prismaMock.session.findFirst).mockResolvedValue(prismaSession);

		const result = await sessionRepository.findSessionByToken(tokenHash);

		expect(prismaMock.session.findFirst).toHaveBeenCalledWith({
			where: { tokenHash },
		});
		expect(result).toEqual(prismaSession);
	});

	it("should delete session by jti", async () => {
		const jti = faker.string.uuid();
		vi.mocked(prismaMock.session.deleteMany).mockResolvedValue({ count: 1 });

		await sessionRepository.deleteSession(jti);

		expect(prismaMock.session.deleteMany).toHaveBeenCalledWith({
			where: { accessJti: jti },
		});
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should return null when session by token does not exist", async () => {
		const tokenHash = faker.string.hexadecimal({ length: 64 });
		vi.mocked(prismaMock.session.findFirst).mockResolvedValue(null);

		const result = await sessionRepository.findSessionByToken(tokenHash);

		expect(result).toBeNull();
	});

	it("should propagate create session errors", async () => {
		const sessionData = {
			id: faker.string.uuid(),
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			userId: faker.string.uuid(),
			deviceInfo: "desktop",
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			expiresAt: new Date(Date.now() + 60_000),
		};
		const dbError = new Error("db-create-session-error");
		vi.mocked(prismaMock.session.create).mockRejectedValue(dbError);

		await expect(sessionRepository.createSession(sessionData)).rejects.toThrow(
			dbError,
		);
	});

	it("should propagate renew session errors", async () => {
		const sessionId = faker.string.uuid();
		const renewData = {
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			expiresAt: new Date(Date.now() + 120_000),
			lastUsedAt: new Date(),
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			deviceInfo: "desktop",
		};
		const dbError = new Error("db-renew-session-error");
		vi.mocked(prismaMock.session.update).mockRejectedValue(dbError);

		await expect(
			sessionRepository.renewSession(sessionId, renewData),
		).rejects.toThrow(dbError);
	});

	it("should propagate delete session errors", async () => {
		const jti = faker.string.uuid();
		const dbError = new Error("db-delete-session-error");
		vi.mocked(prismaMock.session.deleteMany).mockRejectedValue(dbError);

		await expect(sessionRepository.deleteSession(jti)).rejects.toThrow(dbError);
	});
});
