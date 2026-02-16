import type bcrypt from "bcryptjs";
import type { RequestHandler } from "express";
import type { SignOptions } from "jsonwebtoken";
import type { ValidateOptions } from "@/api/middlewares/validate.middleware";
import type { PrismaClient } from "@/generated/prisma/client";
import type { IAuthService } from "@/modules/auth/auth.ports";
import type { ISessionRepository } from "@/modules/session/session.ports";
import type {
	IUsersRepository,
	IUsersService,
} from "@/modules/users/users.ports";
import type { createHasherService } from "@/shared/services/hasher.service";
import type { ITokenManager } from "@/shared/services/token-manager.service";

export interface Cradle {
	encryptor: typeof bcrypt;
	tokenManager: ITokenManager;
	dbClient: PrismaClient;
	idGenerator: () => string;
	config: {
		accessTokenSecret: string;
		refreshTokenSecret: string;
		accessTokenExpirationTime: SignOptions["expiresIn"];
		refreshTokenExpirationTime: SignOptions["expiresIn"];
		saltRounds: number;
	};
	usersRepository: IUsersRepository;
	sessionRepository: ISessionRepository;
	usersService: IUsersService;
	authService: IAuthService;
	hasherService: ReturnType<typeof createHasherService>;
	authMiddleware: RequestHandler;
	validateMiddleware: <T extends import("valibot").GenericSchema>(
		options: ValidateOptions<T>,
	) => RequestHandler;
}
