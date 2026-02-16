import { faker } from "@faker-js/faker";
import { describe, expect, it } from "vitest";
import { createHasherService } from "../hasher.service";

describe("hasher.service", () => {
	const hasherService = createHasherService();

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should generate deterministic sha256 hash", () => {
		const value = faker.internet.password();
		const hash1 = hasherService.createHash(value);
		const hash2 = hasherService.createHash(value);

		expect(hash1).toBe(hash2);
		expect(hash1).toHaveLength(64);
	});

	it("should validate matching hash", () => {
		const value = faker.internet.password();
		const hash = hasherService.createHash(value);

		expect(hasherService.validarHash(value, hash)).toBe(true);
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should invalidate non-matching hash", () => {
		const original = faker.internet.password();
		const different = faker.internet.password();
		const hash = hasherService.createHash(different);

		expect(hasherService.validarHash(original, hash)).toBe(false);
	});
});
