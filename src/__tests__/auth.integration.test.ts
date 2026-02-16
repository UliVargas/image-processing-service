import { createId } from "@paralleldrive/cuid2";
import bcrypt from "bcryptjs";
import jwt from "jsonwebtoken";
import request from "supertest";
import { afterAll, beforeAll, beforeEach, describe, expect, it } from "vitest";
import { setupIntegrationContainer } from "@/__tests__/setup/integration-container";
import { getSqlitePrisma } from "@/__tests__/setup/sqlite-prisma";
import { app } from "@/app";
import { container } from "@/container";
import type { PrismaClient } from "@/generated/prisma-test/client";

describe("auth.integration", () => {
	let restoreContainer: (() => void) | undefined;
	let dbClient: PrismaClient;

	const testUser = {
		id: createId(),
		email: "integration@example.com",
		username: "integration_user",
		password: "P@ssw0rd123",
		name: "Integration User",
	};

	beforeAll(() => {
		restoreContainer = setupIntegrationContainer();
		dbClient = getSqlitePrisma() as unknown as PrismaClient;
	});

	beforeEach(async () => {
		await dbClient.user.create({
			data: {
				id: testUser.id,
				email: testUser.email,
				username: testUser.username,
				password: bcrypt.hashSync(testUser.password, 10),
				name: testUser.name,
			},
		});
	});

	afterAll(() => {
		restoreContainer?.();
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	describe("success cases", () => {
		it("POST /api/auth/login should succeed with valid credentials", async () => {
			const response = await request(app).post("/api/auth/login").send({
				email: testUser.email,
				password: testUser.password,
			});

			expect(response.status).toBe(200);
			expect(response.body.success).toBe(true);
			expect(response.body.data.tokens.accessToken).toEqual(expect.any(String));
			expect(response.body.data.tokens.refreshToken).toEqual(
				expect.any(String),
			);
		});

		it("POST /api/auth/logout should succeed with valid token", async () => {
			const loginResponse = await request(app).post("/api/auth/login").send({
				email: testUser.email,
				password: testUser.password,
			});

			const accessToken = loginResponse.body.data.tokens.accessToken as string;
			const logoutResponse = await request(app)
				.post("/api/auth/logout")
				.set("Authorization", `Bearer ${accessToken}`)
				.send();

			expect(logoutResponse.status).toBe(200);
			expect(logoutResponse.body.success).toBe(true);
		});

		it("POST /api/auth/renew-session should succeed with valid refresh token", async () => {
			const loginResponse = await request(app).post("/api/auth/login").send({
				email: testUser.email,
				password: testUser.password,
			});
			const refreshToken = loginResponse.body.data.tokens
				.refreshToken as string;

			const renewResponse = await request(app)
				.post("/api/auth/renew-session")
				.send({ refreshToken });

			expect(renewResponse.status).toBe(200);
			expect(renewResponse.body.success).toBe(true);
			expect(renewResponse.body.data.tokens.accessToken).toEqual(
				expect.any(String),
			);
			expect(renewResponse.body.data.tokens.refreshToken).toEqual(
				expect.any(String),
			);
		});
	});

	// ===============================================================
	// Fase RED (reglas y errores)
	// ===============================================================
	describe("error cases", () => {
		it("POST /api/auth/login should fail with invalid credentials", async () => {
			const response = await request(app).post("/api/auth/login").send({
				email: testUser.email,
				password: "invalid-password",
			});

			expect(response.status).toBe(401);
			expect(response.body.success).toBe(false);
			expect(response.body.error.code).toBe("UNAUTHORIZED");
		});

		it("POST /api/auth/login should fail with invalid payload", async () => {
			const response = await request(app).post("/api/auth/login").send({
				email: testUser.email,
			});

			expect(response.status).toBe(400);
			expect(response.body.success).toBe(false);
			expect(response.body.error.code).toBe("VALIDATION_ERROR");
			expect(response.body.error.details).toEqual(expect.any(Object));
		});

		it("POST /api/auth/logout should fail without token", async () => {
			const response = await request(app).post("/api/auth/logout").send();

			expect(response.status).toBe(401);
			expect(response.body.success).toBe(false);
			expect(response.body.error.code).toBe("TOKEN_NOT_PROVIDED");
		});

		it("POST /api/auth/logout should fail with invalid token", async () => {
			const response = await request(app)
				.post("/api/auth/logout")
				.set("Authorization", "Bearer this-is-not-a-valid-jwt")
				.send();

			expect(response.status).toBe(401);
			expect(response.body.success).toBe(false);
			expect(response.body.error.code).toBe("INVALID_TOKEN");
		});

		it("POST /api/auth/logout should fail with expired token", async () => {
			const config = container.resolve("config");
			const expiredToken = jwt.sign(
				{ sub: testUser.id, jti: createId() },
				config.accessTokenSecret,
				{ expiresIn: -1 },
			);

			const response = await request(app)
				.post("/api/auth/logout")
				.set("Authorization", `Bearer ${expiredToken}`)
				.send();

			expect(response.status).toBe(401);
			expect(response.body.success).toBe(false);
			expect(response.body.error.code).toBe("TOKEN_EXPIRED");
		});

		it("POST /api/auth/renew-session should fail with invalid refresh token", async () => {
			const renewResponse = await request(app)
				.post("/api/auth/renew-session")
				.send({ refreshToken: createId() });

			expect(renewResponse.status).toBe(404);
			expect(renewResponse.body.success).toBe(false);
			expect(renewResponse.body.error.code).toBe("NOT_FOUND");
		});

		it("POST /api/auth/renew-session should fail with invalid payload", async () => {
			const renewResponse = await request(app)
				.post("/api/auth/renew-session")
				.send({ refreshToken: "invalid-token" });

			expect(renewResponse.status).toBe(400);
			expect(renewResponse.body.success).toBe(false);
			expect(renewResponse.body.error.code).toBe("VALIDATION_ERROR");
			expect(renewResponse.body.error.details).toEqual(expect.any(Object));
		});

		it("POST /api/auth/renew-session should fail with device mismatch", async () => {
			const loginResponse = await request(app)
				.post("/api/auth/login")
				.set("User-Agent", "Mozilla/5.0 Chrome/122.0.0.0 Safari/537.36")
				.send({
					email: testUser.email,
					password: testUser.password,
				});
			const refreshToken = loginResponse.body.data.tokens
				.refreshToken as string;

			const renewResponse = await request(app)
				.post("/api/auth/renew-session")
				.set("User-Agent", "Mozilla/5.0 Firefox/123.0")
				.send({ refreshToken });

			expect(renewResponse.status).toBe(401);
			expect(renewResponse.body.success).toBe(false);
			expect(renewResponse.body.error.code).toBe("SESSION_DEVICE_MISMATCH");
		});

		it("auth errors should keep consistent response format", async () => {
			const response = await request(app).post("/api/auth/login").send({
				email: testUser.email,
				password: "wrong-password",
			});

			expect(response.status).toBe(401);
			expect(response.body).toEqual(
				expect.objectContaining({
					success: false,
					error: expect.objectContaining({
						message: expect.any(String),
						code: expect.any(String),
					}),
					timestamp: expect.any(String),
					path: expect.any(String),
				}),
			);
		});
	});
});
