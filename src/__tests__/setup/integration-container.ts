import { asValue } from "awilix";
import { container } from "@/container";
import { createAuthService } from "@/modules/auth/auth.service";
import { createSessionRepository } from "@/modules/session/session.repository";
import { createUsersRepository } from "@/modules/users/users.repository";
import { createUsersService } from "@/modules/users/users.service";
import type { Cradle } from "@/shared/di/types";
import { getSqlitePrisma } from "./sqlite-prisma";

type RestoreState = {
	dbClient: Cradle["dbClient"];
	usersRepository: Cradle["usersRepository"];
	sessionRepository: Cradle["sessionRepository"];
	usersService: Cradle["usersService"];
	authService: Cradle["authService"];
};

export const setupIntegrationContainer = (): (() => void) => {
	const previous: RestoreState = {
		dbClient: container.resolve("dbClient"),
		usersRepository: container.resolve("usersRepository"),
		sessionRepository: container.resolve("sessionRepository"),
		usersService: container.resolve("usersService"),
		authService: container.resolve("authService"),
	};

	const testDbClient = getSqlitePrisma() as unknown as Cradle["dbClient"];
	const usersRepository = createUsersRepository({ dbClient: testDbClient });
	const sessionRepository = createSessionRepository({ dbClient: testDbClient });
	const usersService = createUsersService({
		usersRepository,
		idGenerator: container.resolve("idGenerator"),
		encryptor: container.resolve("encryptor"),
		config: container.resolve("config"),
	});
	const authService = createAuthService({
		usersRepository,
		sessionRepository,
		config: container.resolve("config"),
		tokenManager: container.resolve("tokenManager"),
		encryptor: container.resolve("encryptor"),
		hasherService: container.resolve("hasherService"),
		idGenerator: container.resolve("idGenerator"),
	});

	container.register({
		dbClient: asValue(testDbClient),
		usersRepository: asValue(usersRepository),
		sessionRepository: asValue(sessionRepository),
		usersService: asValue(usersService),
		authService: asValue(authService),
	});

	return () => {
		container.register({
			dbClient: asValue(previous.dbClient),
			usersRepository: asValue(previous.usersRepository),
			sessionRepository: asValue(previous.sessionRepository),
			usersService: asValue(previous.usersService),
			authService: asValue(previous.authService),
		});
	};
};
