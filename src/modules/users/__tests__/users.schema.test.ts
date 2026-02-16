import { faker } from "@faker-js/faker";
import { createId } from "@paralleldrive/cuid2";
import * as v from "valibot";
import { describe, expect, it } from "vitest";
import {
	CreateUserInputSchema,
	UpdateUserInputSchema,
	UserParamsSchema,
} from "../users.schema";

describe("users.schema", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should validate create user input", () => {
		const input = {
			email: faker.internet.email(),
			username: faker.internet.username().replace(/[^a-zA-Z0-9_-]/g, "abc"),
			password: faker.internet.password({ length: 10 }),
			name: faker.person.fullName(),
		};

		const parsed = v.parse(CreateUserInputSchema, input);
		expect(parsed).toMatchObject(input);
	});

	it("should validate update user input", () => {
		const parsed = v.parse(UpdateUserInputSchema, {
			name: faker.person.fullName(),
		});
		expect(parsed).toEqual({ name: expect.any(String) });
	});

	it("should validate user params id", () => {
		const parsed = v.parse(UserParamsSchema, {
			id: createId(),
		});
		expect(parsed.id).toBeTruthy();
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should fail create user input with invalid email", () => {
		expect(() =>
			v.parse(CreateUserInputSchema, {
				email: "not-an-email",
				username: "valid_user",
				password: "password123",
			}),
		).toThrow();
	});

	it("should fail create user input with short username", () => {
		expect(() =>
			v.parse(CreateUserInputSchema, {
				email: faker.internet.email(),
				username: "ab",
				password: "password123",
			}),
		).toThrow();
	});

	it("should fail create user input with invalid username chars", () => {
		expect(() =>
			v.parse(CreateUserInputSchema, {
				email: faker.internet.email(),
				username: "invalid user",
				password: "password123",
			}),
		).toThrow();
	});

	it("should fail create user input with short password", () => {
		expect(() =>
			v.parse(CreateUserInputSchema, {
				email: faker.internet.email(),
				username: "valid_user",
				password: "1234567",
			}),
		).toThrow();
	});

	it("should fail update user input when name is too long", () => {
		expect(() =>
			v.parse(UpdateUserInputSchema, {
				name: "a".repeat(101),
			}),
		).toThrow();
	});

	it("should fail user params id when id is invalid", () => {
		expect(() =>
			v.parse(UserParamsSchema, {
				id: "invalid-id",
			}),
		).toThrow();
	});
});
