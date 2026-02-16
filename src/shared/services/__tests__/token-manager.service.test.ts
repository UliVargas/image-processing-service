import { faker } from "@faker-js/faker";
import { describe, expect, it } from "vitest";
import { createTokenManager } from "../token-manager.service";

describe("token-manager.service", () => {
	const tokenManager = createTokenManager();

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should sign and verify token payload", () => {
		const payload = {
			sub: faker.string.uuid(),
			jti: faker.string.uuid(),
		};
		const secret = "test-secret-1234567890";

		const token = tokenManager.sign(payload, secret, { expiresIn: "15m" });
		const decoded = tokenManager.verify(token, secret) as {
			sub: string;
			jti: string;
		};

		expect(token).toEqual(expect.any(String));
		expect(decoded.sub).toBe(payload.sub);
		expect(decoded.jti).toBe(payload.jti);
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should throw when verifying with wrong secret", () => {
		const token = tokenManager.sign(
			{ sub: faker.string.uuid() },
			"correct-secret",
			{ expiresIn: "15m" },
		);

		expect(() => tokenManager.verify(token, "wrong-secret")).toThrow();
	});
});
