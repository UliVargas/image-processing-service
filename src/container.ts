import { createId } from "@paralleldrive/cuid2";
import { ENV_CONFIG } from "@shared/config/env.config";
import { prisma } from "@shared/database/db.client";
import { asFunction, asValue, createContainer, Lifetime } from "awilix";
import bcrypt from "bcryptjs";
import { createAuthMiddleware } from "@/api/middlewares/auth.middleware";
import { validate } from "@/api/middlewares/validate.middleware";
import { createTokenManager } from "@/shared/services/token-manager.service";
import { createAuthService } from "./modules/auth/auth.service";
import { createSessionRepository } from "./modules/session/session.repository";
import { createUsersRepository } from "./modules/users/users.repository";
import { createUsersService } from "./modules/users/users.service";
import { createHasherService } from "./shared/services/hasher.service";

const container = createContainer({
	injectionMode: "PROXY",
});

container.register({
	encryptor: asValue(bcrypt),
	dbClient: asValue(prisma),
	idGenerator: asValue(createId),
	config: asValue({
		accessTokenSecret: ENV_CONFIG.ACCESS_TOKEN_SECRET,
		refreshTokenSecret: ENV_CONFIG.REFRESH_TOKEN_SECRET,
		accessTokenExpirationTime: ENV_CONFIG.ACCESS_TOKEN_EXPIRATION_TIME,
		refreshTokenExpirationTime: ENV_CONFIG.REFRESH_TOKEN_EXPIRATION_TIME,
		saltRounds: ENV_CONFIG.SALT_ROUNDS,
	}),
	// ===============================================================
	// Repositories
	// ===============================================================
	usersRepository: asFunction(createUsersRepository, {
		lifetime: Lifetime.SINGLETON,
	}),
	sessionRepository: asFunction(createSessionRepository, {
		lifetime: Lifetime.SINGLETON,
	}),
	// ===============================================================
	// Services
	// ===============================================================
	usersService: asFunction(createUsersService, {
		lifetime: Lifetime.SINGLETON,
	}),
	authService: asFunction(createAuthService, {
		lifetime: Lifetime.SINGLETON,
	}),
	hasherService: asFunction(createHasherService, {
		lifetime: Lifetime.SINGLETON,
	}),
	tokenManager: asFunction(createTokenManager, {
		lifetime: Lifetime.SINGLETON,
	}),
	// ===============================================================
	// Middlewares
	// ===============================================================
	validateMiddleware: asValue(validate),
	authMiddleware: asFunction(createAuthMiddleware, {
		lifetime: Lifetime.SINGLETON,
	}),
});

export { container };
