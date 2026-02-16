import { faker } from "@faker-js/faker";
import type bcrypt from "bcryptjs";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { IUsersRepository } from "../users.ports";
import { createUsersService } from "../users.service";

describe("Users Service", () => {
	let usersService: ReturnType<typeof createUsersService>;
	let usersRepositoryMock: IUsersRepository;
	let idGeneratorMock: () => string;
	let encryptorMock: typeof bcrypt;
	let configMock: {
		saltRounds: number;
		accessTokenSecret: string;
		refreshTokenSecret: string;
		accessTokenExpirationTime: "15m";
		refreshTokenExpirationTime: "7d";
	};

	beforeEach(() => {
		usersRepositoryMock = {
			findUserByEmailAndUsername: vi.fn(),
			createUser: vi.fn(),
			deleteUser: vi.fn(),
			findUserById: vi.fn(),
			updateUser: vi.fn(),
			getAllUsers: vi.fn(),
		};
		idGeneratorMock = vi.fn().mockReturnValue("clx123456789");
		encryptorMock = {
			hash: vi.fn().mockResolvedValue("hashed-password"),
		} as unknown as typeof bcrypt;

		configMock = {
			saltRounds: 10,
			accessTokenSecret: "test-access-secret",
			refreshTokenSecret: "test-refresh-secret",
			accessTokenExpirationTime: "15m",
			refreshTokenExpirationTime: "7d",
		};

		usersService = createUsersService({
			usersRepository: usersRepositoryMock,
			idGenerator: idGeneratorMock,
			encryptor: encryptorMock,
			config: configMock,
		});
	});

	// ===============================================================
	// createUser
	// ===============================================================
	describe("createUser", () => {
		// ===============================================================
		// Fase RED (reglas y errores)
		// ===============================================================
		it("should throw USER_ALREADY_EXISTS when user already exists", async () => {
			const userData = {
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: faker.internet.password({ length: 12 }),
				name: faker.person.fullName(),
			};

			vi.mocked(
				usersRepositoryMock.findUserByEmailAndUsername,
			).mockResolvedValue({
				id: faker.string.uuid(),
				email: userData.email,
				username: userData.username,
				password: faker.internet.password(),
				name: userData.name,
				createdAt: new Date(),
				updatedAt: new Date(),
			});

			await expect(usersService.createUser(userData)).rejects.toMatchObject({
				statusCode: 409,
				code: "USER_ALREADY_EXISTS",
			});
			expect(usersRepositoryMock.createUser).not.toHaveBeenCalled();
		});

		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================

		it("should create a new user successfully", async () => {
			const userData = {
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: faker.internet.password({ length: 12 }),
				name: faker.person.fullName(),
			};
			const createdUser = {
				id: faker.string.alphanumeric(10),
				...userData,
				password: "hashed-password",
				createdAt: new Date(),
				updatedAt: new Date(),
			};

			vi.mocked(
				usersRepositoryMock.findUserByEmailAndUsername,
			).mockResolvedValue(null);
			vi.mocked(usersRepositoryMock.createUser).mockResolvedValue(createdUser);

			const result = await usersService.createUser(userData);

			expect(
				usersRepositoryMock.findUserByEmailAndUsername,
			).toHaveBeenCalledWith({
				email: userData.email,
				username: userData.username,
			});
			expect(encryptorMock.hash).toHaveBeenCalledWith(
				userData.password,
				configMock.saltRounds,
			);
			expect(usersRepositoryMock.createUser).toHaveBeenCalledWith({
				id: "clx123456789",
				...userData,
				password: "hashed-password",
			});
			expect(result).toEqual({
				id: createdUser.id,
				email: userData.email,
				username: userData.username,
				name: userData.name,
				password: "hashed-password",
				createdAt: expect.any(Date),
				updatedAt: expect.any(Date),
			});
		});
	});

	// ===============================================================
	// getAllUsers
	// ===============================================================
	describe("getAllUsers", () => {
		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================
		it("should return all users", async () => {
			const users = [
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

			vi.mocked(usersRepositoryMock.getAllUsers).mockResolvedValue(users);

			await expect(usersService.getAllUsers()).resolves.toEqual(users);
		});
	});

	// ===============================================================
	// updateUser
	// ===============================================================
	describe("updateUser", () => {
		// ===============================================================
		// Fase RED (reglas y errores)
		// ===============================================================
		it("should throw USER_NOT_FOUND when user does not exist", async () => {
			const userId = faker.string.uuid();
			vi.mocked(usersRepositoryMock.findUserById).mockResolvedValue(null);

			await expect(
				usersService.updateUser(userId, { name: faker.person.fullName() }),
			).rejects.toMatchObject({
				statusCode: 404,
				code: "USER_NOT_FOUND",
			});
		});

		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================

		it("should update user when it exists", async () => {
			const userId = faker.string.uuid();
			const updatePayload = { name: faker.person.fullName() };
			const existingUser = {
				id: userId,
				email: faker.internet.email(),
				username: faker.internet.username(),
				password: faker.internet.password(),
				name: faker.person.fullName(),
				createdAt: new Date(),
				updatedAt: new Date(),
			};
			const updatedUser = { ...existingUser, ...updatePayload };

			vi.mocked(usersRepositoryMock.findUserById).mockResolvedValue(
				existingUser,
			);
			vi.mocked(usersRepositoryMock.updateUser).mockResolvedValue(updatedUser);

			await expect(
				usersService.updateUser(userId, updatePayload),
			).resolves.toEqual(updatedUser);
			expect(usersRepositoryMock.updateUser).toHaveBeenCalledWith(
				userId,
				updatePayload,
			);
		});
	});

	// ===============================================================
	// deleteUser
	// ===============================================================
	describe("deleteUser", () => {
		// ===============================================================
		// Fase GREEN (flujo esperado)
		// ===============================================================
		it("should delegate deletion to repository", async () => {
			const userId = faker.string.uuid();
			vi.mocked(usersRepositoryMock.deleteUser).mockResolvedValue();

			await expect(usersService.deleteUser(userId)).resolves.toBeUndefined();
			expect(usersRepositoryMock.deleteUser).toHaveBeenCalledWith(userId);
		});
	});
});
