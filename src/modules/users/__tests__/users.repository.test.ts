import { beforeEach, describe, expect, it, vi } from "vitest";
import type { PrismaClient } from "@/generated/prisma/client";
import { createUsersRepository } from "../users.repository";

describe("Users Repository", () => {
	let prismaMock: PrismaClient;
	let usersRepository: ReturnType<typeof createUsersRepository>;

	beforeEach(() => {
		prismaMock = {
			user: {
				findUnique: vi.fn(),
				create: vi.fn(),
				update: vi.fn(),
				delete: vi.fn(),
			},
		} as unknown as PrismaClient;

		usersRepository = createUsersRepository({
			dbClient: prismaMock,
		});
	});

	it("should create a new user", async () => {
		// Arrange
		const userData = {
			id: "clx123456789",
			username: "johndoe",
			name: "John Doe",
			email: "john.doe@example.com",
			password: "hashedpassword",
		};

		const createdUser = {
			...userData,
			createdAt: new Date(),
			updatedAt: new Date(),
		};

		vi.mocked(prismaMock.user.create).mockResolvedValue(createdUser);

		// Act
		const result = await usersRepository.createUser(userData);

		// Assert
		expect(prismaMock.user.create).toHaveBeenCalledWith({
			data: userData,
		});
		expect(result).toEqual(createdUser);
	});
});
