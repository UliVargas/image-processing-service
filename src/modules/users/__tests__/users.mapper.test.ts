import { faker } from "@faker-js/faker";
import dayjs from "dayjs";
import { describe, expect, it } from "vitest";
import { toUserListResponse, toUserResponse } from "../users.mapper";

describe("users.mapper", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should map one user to response dto", () => {
		const user = {
			id: faker.string.uuid(),
			email: faker.internet.email(),
			username: faker.internet.username(),
			name: faker.person.fullName(),
			password: faker.internet.password(),
			createdAt: new Date("2026-02-15T10:20:30.123Z"),
			updatedAt: new Date("2026-02-15T11:20:30.123Z"),
		};

		const result = toUserResponse(user);

		expect(result).toEqual({
			id: user.id,
			email: user.email,
			username: user.username,
			name: user.name,
			createdAt: dayjs(user.createdAt).format("YYYY-MM-DDTHH:mm:ss.SSS"),
			updatedAt: dayjs(user.updatedAt).format("YYYY-MM-DDTHH:mm:ss.SSS"),
		});
	});

	it("should map user list to response dto list", () => {
		const users = [
			{
				id: faker.string.uuid(),
				email: faker.internet.email(),
				username: faker.internet.username(),
				name: faker.person.fullName(),
				password: faker.internet.password(),
				createdAt: new Date("2026-02-15T10:20:30.123Z"),
				updatedAt: new Date("2026-02-15T11:20:30.123Z"),
			},
		];

		const result = toUserListResponse(users);
		expect(result).toHaveLength(1);
		expect(result[0]).toMatchObject({
			id: users[0].id,
			email: users[0].email,
			username: users[0].username,
			name: users[0].name,
		});
	});
});
