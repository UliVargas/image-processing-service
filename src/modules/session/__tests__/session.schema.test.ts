import { faker } from "@faker-js/faker";
import * as v from "valibot";
import { describe, expect, it } from "vitest";
import {
	CreateSessionDataSchema,
	RenewSessionDataSchema,
	SessionSchema,
} from "../session.schema";

describe("session.schema", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should parse session entity shape", () => {
		const input = {
			id: faker.string.uuid(),
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			userId: faker.string.uuid(),
			deviceInfo: "desktop",
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			expiresAt: new Date(),
			lastUsedAt: new Date(),
			isValid: true,
			createdAt: new Date(),
		};

		expect(v.parse(SessionSchema, input)).toEqual(input);
	});

	it("should parse create session payload", () => {
		const input = {
			id: faker.string.uuid(),
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			userId: faker.string.uuid(),
			deviceInfo: "desktop",
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			expiresAt: new Date(),
		};

		expect(v.parse(CreateSessionDataSchema, input)).toEqual(input);
	});

	it("should parse renew session payload", () => {
		const input = {
			tokenHash: faker.string.hexadecimal({ length: 64 }),
			accessJti: faker.string.uuid(),
			expiresAt: new Date(),
			lastUsedAt: new Date(),
			ipAddress: faker.internet.ip(),
			userAgent: faker.internet.userAgent(),
			deviceInfo: "desktop",
		};

		expect(v.parse(RenewSessionDataSchema, input)).toEqual(input);
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should fail parse session entity shape with invalid isValid type", () => {
		expect(() =>
			v.parse(SessionSchema, {
				id: faker.string.uuid(),
				tokenHash: faker.string.hexadecimal({ length: 64 }),
				accessJti: faker.string.uuid(),
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: faker.internet.ip(),
				userAgent: faker.internet.userAgent(),
				expiresAt: new Date(),
				lastUsedAt: new Date(),
				isValid: "true",
				createdAt: new Date(),
			}),
		).toThrow();
	});

	it("should fail parse when create payload has invalid expiresAt", () => {
		expect(() =>
			v.parse(CreateSessionDataSchema, {
				id: faker.string.uuid(),
				tokenHash: faker.string.hexadecimal({ length: 64 }),
				accessJti: faker.string.uuid(),
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: faker.internet.ip(),
				userAgent: faker.internet.userAgent(),
				expiresAt: "not-a-date",
			}),
		).toThrow();
	});

	it("should fail parse when create payload misses required fields", () => {
		expect(() =>
			v.parse(CreateSessionDataSchema, {
				id: faker.string.uuid(),
				tokenHash: faker.string.hexadecimal({ length: 64 }),
				accessJti: faker.string.uuid(),
				userId: faker.string.uuid(),
				deviceInfo: "desktop",
				ipAddress: faker.internet.ip(),
				expiresAt: new Date(),
			}),
		).toThrow();
	});

	it("should fail parse renew session payload with invalid lastUsedAt", () => {
		expect(() =>
			v.parse(RenewSessionDataSchema, {
				tokenHash: faker.string.hexadecimal({ length: 64 }),
				accessJti: faker.string.uuid(),
				expiresAt: new Date(),
				lastUsedAt: "invalid-date",
				ipAddress: faker.internet.ip(),
				userAgent: faker.internet.userAgent(),
				deviceInfo: "desktop",
			}),
		).toThrow();
	});
});
