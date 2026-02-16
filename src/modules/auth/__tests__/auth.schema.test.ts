import { faker } from "@faker-js/faker";
import { createId } from "@paralleldrive/cuid2";
import * as v from "valibot";
import { describe, expect, it } from "vitest";
import { LoginDataSchema, RenewSessionDataSchema } from "../auth.schema";

describe("auth.schema", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should parse login by email", () => {
		const input = {
			email: faker.internet.email(),
			password: faker.internet.password(),
		};

		expect(v.parse(LoginDataSchema, input)).toEqual(input);
	});

	it("should parse login by username", () => {
		const input = {
			username: faker.internet.username(),
			password: faker.internet.password(),
		};

		expect(v.parse(LoginDataSchema, input)).toEqual(input);
	});

	it("should parse renew session payload", () => {
		const payload = {
			refreshToken: createId(),
		};

		expect(v.parse(RenewSessionDataSchema, payload)).toEqual(payload);
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should fail login when neither email nor username is provided", () => {
		expect(() =>
			v.parse(LoginDataSchema, {
				password: faker.internet.password(),
			}),
		).toThrow();
	});

	it("should fail login with invalid email format", () => {
		expect(() =>
			v.parse(LoginDataSchema, {
				email: "invalid-email",
				password: faker.internet.password(),
			}),
		).toThrow();
	});

	it("should fail login when password is missing", () => {
		expect(() =>
			v.parse(LoginDataSchema, {
				email: faker.internet.email(),
			}),
		).toThrow();
	});

	it("should fail renew session payload with invalid refreshToken", () => {
		expect(() =>
			v.parse(RenewSessionDataSchema, {
				refreshToken: "not-a-cuid2",
			}),
		).toThrow();
	});
});
