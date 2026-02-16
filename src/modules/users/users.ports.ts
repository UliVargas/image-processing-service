import type { NextFunction, Request, Response } from "express";
import type { User } from "./users.entity";
import type {
	CreateUserData,
	CreateUserInput,
	UpdateUserInput,
} from "./users.schema";

// ===============================================================
// Puerto del Repositorio (abstracción de persistencia)
// ===============================================================
export interface IUsersRepository {
	createUser: (data: CreateUserData) => Promise<User>;
	findUserByEmailAndUsername: (data: {
		email?: string;
		username?: string;
	}) => Promise<User | null>;
	findUserById: (id: string) => Promise<User | null>;
	getAllUsers: () => Promise<User[]>;
	updateUser: (id: string, data: UpdateUserInput) => Promise<User>;
	deleteUser: (id: string) => Promise<void>;
}

// ===============================================================
// Puerto del Servicio (casos de uso)
// ===============================================================
export interface IUsersService {
	createUser: (data: CreateUserInput) => Promise<User>;
	getAllUsers: () => Promise<User[]>;
	updateUser: (id: string, data: UpdateUserInput) => Promise<User>;
	deleteUser: (id: string) => Promise<void>;
}

// ===============================================================
// Puerto del Controlador (interfaz HTTP)
// ===============================================================
export interface IUsersController {
	createUser: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
	getAllUsers: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
	updateUser: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
	deleteUser: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
}
