import { faker } from "@faker-js/faker";
import dayjs from "dayjs";
import { describe, expect, it } from "vitest";
import { toAuthResponse } from "../auth.mapper";

describe("auth.mapper", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should map auth entity to response dto", () => {
		const auth = {
			session: {
				id: faker.string.uuid(),
				tokenHash: faker.string.hexadecimal({ length: 64 }),
				accessJti: faker.string.uuid(),
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: faker.internet.ip(),
				userAgent: faker.internet.userAgent(),
				expiresAt: new Date("2026-02-15T12:30:45.123Z"),
				lastUsedAt: null,
				isValid: true,
				createdAt: new Date(),
			},
			tokens: {
				accessToken: faker.string.alphanumeric(40),
				refreshToken: faker.string.alphanumeric(40),
			},
		};

		const result = toAuthResponse(auth);

		expect(result).toEqual({
			session: {
				sessionId: auth.session.id,
				expiresAt: dayjs(auth.session.expiresAt).format(
					"YYYY-MM-DDTHH:mm:ss.SSS",
				),
			},
			tokens: {
				accessToken: auth.tokens.accessToken,
				refreshToken: auth.tokens.refreshToken,
			},
		});
	});
});
