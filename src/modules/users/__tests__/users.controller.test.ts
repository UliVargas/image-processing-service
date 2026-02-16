import { faker } from "@faker-js/faker";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { createUsersController } from "../users.controller";

const createResMock = () => {
	const json = vi.fn();
	const status = vi.fn().mockReturnValue({ json });
	return { status, json };
};

describe("users.controller", () => {
	const usersServiceMock = {
		createUser: vi.fn(),
		getAllUsers: vi.fn(),
		updateUser: vi.fn(),
		deleteUser: vi.fn(),
	};

	const next = vi.fn();
	let controller: ReturnType<typeof createUsersController>;

	beforeEach(() => {
		vi.clearAllMocks();
		controller = createUsersController({
			usersService: usersServiceMock,
		});
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	describe("success cases", () => {
		it("should create user and return 201 response", async () => {
			const req = {
				body: {
					email: faker.internet.email(),
					username: faker.internet.username(),
					password: faker.internet.password(),
					name: faker.person.fullName(),
				},
			};
			const user = {
				id: faker.string.uuid(),
				...req.body,
				createdAt: new Date(),
				updatedAt: new Date(),
			};
			const res = createResMock();
			usersServiceMock.createUser.mockResolvedValue(user);

			await controller.createUser(req as never, res as never, next);

			expect(usersServiceMock.createUser).toHaveBeenCalledWith(req.body);
			expect(res.status).toHaveBeenCalledWith(201);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Usuario creado exitosamente",
				}),
			);
		});

		it("should return users list", async () => {
			const req = {};
			const res = createResMock();
			const users = [
				{
					id: faker.string.uuid(),
					email: faker.internet.email(),
					username: faker.internet.username(),
					name: faker.person.fullName(),
					password: faker.internet.password(),
					createdAt: new Date(),
					updatedAt: new Date(),
				},
			];
			usersServiceMock.getAllUsers.mockResolvedValue(users);

			await controller.getAllUsers(req as never, res as never, next);

			expect(usersServiceMock.getAllUsers).toHaveBeenCalled();
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Usuarios obtenidos exitosamente",
				}),
			);
		});

		it("should update user", async () => {
			const req = {
				params: { id: faker.string.uuid() },
				body: { name: faker.person.fullName() },
			};
			const res = createResMock();
			const updatedUser = {
				id: req.params.id,
				email: faker.internet.email(),
				username: faker.internet.username(),
				name: req.body.name,
				password: faker.internet.password(),
				createdAt: new Date(),
				updatedAt: new Date(),
			};
			usersServiceMock.updateUser.mockResolvedValue(updatedUser);

			await controller.updateUser(req as never, res as never, next);

			expect(usersServiceMock.updateUser).toHaveBeenCalledWith(
				req.params.id,
				req.body,
			);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Usuario actualizado exitosamente",
				}),
			);
		});

		it("should delete user", async () => {
			const req = {
				params: { id: faker.string.uuid() },
			};
			const res = createResMock();
			usersServiceMock.deleteUser.mockResolvedValue(undefined);

			await controller.deleteUser(req as never, res as never, next);

			expect(usersServiceMock.deleteUser).toHaveBeenCalledWith(req.params.id);
			expect(res.json).toHaveBeenCalledWith(
				expect.objectContaining({
					success: true,
					message: "Usuario eliminado exitosamente",
				}),
			);
		});
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	describe("error cases", () => {
		it("should forward errors to next on createUser", async () => {
			const req = { body: {} };
			const res = createResMock();
			const error = new Error("boom");
			usersServiceMock.createUser.mockRejectedValue(error);

			await controller.createUser(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});

		it("should forward errors to next on getAllUsers", async () => {
			const req = {};
			const res = createResMock();
			const error = new Error("getAllUsers-error");
			usersServiceMock.getAllUsers.mockRejectedValue(error);

			await controller.getAllUsers(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});

		it("should forward errors to next on updateUser", async () => {
			const req = {
				params: { id: faker.string.uuid() },
				body: { name: faker.person.fullName() },
			};
			const res = createResMock();
			const error = new Error("updateUser-error");
			usersServiceMock.updateUser.mockRejectedValue(error);

			await controller.updateUser(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});

		it("should forward errors to next on deleteUser", async () => {
			const req = {
				params: { id: faker.string.uuid() },
			};
			const res = createResMock();
			const error = new Error("deleteUser-error");
			usersServiceMock.deleteUser.mockRejectedValue(error);

			await controller.deleteUser(req as never, res as never, next);

			expect(next).toHaveBeenCalledWith(error);
		});
	});
});
