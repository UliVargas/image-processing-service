import { faker } from "@faker-js/faker";
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
				findFirst: vi.fn(),
				findMany: vi.fn(),
				create: vi.fn(),
				update: vi.fn(),
				delete: vi.fn(),
			},
		} as unknown as PrismaClient;

		usersRepository = createUsersRepository({
			dbClient: prismaMock,
		});
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should create a new user", async () => {
		// Arrange
		const userData = {
			id: faker.string.uuid(),
			username: faker.internet.username(),
			name: faker.person.fullName(),
			email: faker.internet.email(),
			password: faker.internet.password(),
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

	it("should find user by email or username", async () => {
		const prismaUser = {
			id: faker.string.uuid(),
			email: faker.internet.email(),
			username: faker.internet.username(),
			password: faker.internet.password(),
			name: faker.person.fullName(),
			createdAt: new Date(),
			updatedAt: new Date(),
		};

		vi.mocked(prismaMock.user.findFirst).mockResolvedValue(prismaUser);

		const result = await usersRepository.findUserByEmailAndUsername({
			email: prismaUser.email,
			username: prismaUser.username,
		});

		expect(prismaMock.user.findFirst).toHaveBeenCalledWith({
			where: {
				OR: [{ email: prismaUser.email }, { username: prismaUser.username }],
			},
		});
		expect(result).toEqual(prismaUser);
	});

	it("should find user by id", async () => {
		const prismaUser = {
			id: faker.string.uuid(),
			email: faker.internet.email(),
			username: faker.internet.username(),
			password: faker.internet.password(),
			name: faker.person.fullName(),
			createdAt: new Date(),
			updatedAt: new Date(),
		};

		vi.mocked(prismaMock.user.findUnique).mockResolvedValue(prismaUser);

		const result = await usersRepository.findUserById(prismaUser.id);

		expect(prismaMock.user.findUnique).toHaveBeenCalledWith({
			where: { id: prismaUser.id },
		});
		expect(result).toEqual(prismaUser);
	});

	it("should get all users", async () => {
		const prismaUsers = [
			{
				id: faker.string.uuid(),
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: faker.internet.password(),
				name: faker.person.fullName(),
				createdAt: new Date(),
				updatedAt: new Date(),
			},
		];

		vi.mocked(prismaMock.user.findMany).mockResolvedValue(prismaUsers);

		const result = await usersRepository.getAllUsers();

		expect(prismaMock.user.findMany).toHaveBeenCalled();
		expect(result).toEqual(prismaUsers);
	});

	it("should update user", async () => {
		const userId = faker.string.uuid();
		const updateData = { name: faker.person.fullName() };
		const updatedUser = {
			id: userId,
			email: faker.internet.email(),
			username: faker.internet.username(),
			password: faker.internet.password(),
			name: updateData.name,
			createdAt: new Date(),
			updatedAt: new Date(),
		};

		vi.mocked(prismaMock.user.update).mockResolvedValue(updatedUser);

		const result = await usersRepository.updateUser(userId, updateData);

		expect(prismaMock.user.update).toHaveBeenCalledWith({
			where: { id: userId },
			data: updateData,
		});
		expect(result).toEqual(updatedUser);
	});

	it("should delete user", async () => {
		const userId = faker.string.uuid();
		vi.mocked(prismaMock.user.delete).mockResolvedValue({} as never);

		await usersRepository.deleteUser(userId);

		expect(prismaMock.user.delete).toHaveBeenCalledWith({
			where: { id: userId },
		});
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	it("should return null when no email/username is provided", async () => {
		const result = await usersRepository.findUserByEmailAndUsername({});
		expect(prismaMock.user.findFirst).not.toHaveBeenCalled();
		expect(result).toBeNull();
	});

	it("should return null when user by id does not exist", async () => {
		const userId = faker.string.uuid();
		vi.mocked(prismaMock.user.findUnique).mockResolvedValue(null);

		const result = await usersRepository.findUserById(userId);

		expect(result).toBeNull();
	});

	it("should propagate create user errors", async () => {
		const userData = {
			id: faker.string.uuid(),
			username: faker.internet.username(),
			name: faker.person.fullName(),
			email: faker.internet.email(),
			password: faker.internet.password(),
		};
		const dbError = new Error("db-create-error");
		vi.mocked(prismaMock.user.create).mockRejectedValue(dbError);

		await expect(usersRepository.createUser(userData)).rejects.toThrow(dbError);
	});

	it("should propagate update user errors", async () => {
		const userId = faker.string.uuid();
		const updateData = { name: faker.person.fullName() };
		const dbError = new Error("db-update-error");
		vi.mocked(prismaMock.user.update).mockRejectedValue(dbError);

		await expect(
			usersRepository.updateUser(userId, updateData),
		).rejects.toThrow(dbError);
	});

	it("should propagate delete user errors", async () => {
		const userId = faker.string.uuid();
		const dbError = new Error("db-delete-error");
		vi.mocked(prismaMock.user.delete).mockRejectedValue(dbError);

		await expect(usersRepository.deleteUser(userId)).rejects.toThrow(dbError);
	});
});
