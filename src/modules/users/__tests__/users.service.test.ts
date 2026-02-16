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

	describe("createUser", () => {
		it("should create a new user successfully", async () => {
			// Arrange
			const userData = {
				email: "test@example.com",
				username: "testuser",
				password: "password123",
				name: "Test User",
			};
			const createdUser = {
				id: "clx123456789",
				...userData,
				password: "hashed-password",
				createdAt: new Date(),
				updatedAt: new Date(),
			};

			vi.mocked(
				usersRepositoryMock.findUserByEmailAndUsername,
			).mockResolvedValue(null);
			vi.mocked(usersRepositoryMock.createUser).mockResolvedValue(createdUser);

			// Act
			const result = await usersService.createUser(userData);

			// Assert
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
				id: "clx123456789",
				email: userData.email,
				username: userData.username,
				name: userData.name,
				createdAt: expect.any(String),
				updatedAt: expect.any(String),
			});
		});
	});
});
