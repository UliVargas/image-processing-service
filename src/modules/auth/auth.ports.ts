import type { NextFunction, Request, Response } from "express";
import type { Auth } from "./auth.entity";
import type { LoginData } from "./auth.schema";

export interface ClientInfo {
	ipAddress: string;
	userAgent: string;
	deviceInfo: string;
}

// ===============================================================
// Puerto del Servicio (casos de uso)
// ===============================================================
export interface IAuthService {
	login: (data: LoginData, clientInfo: ClientInfo) => Promise<Auth>;
	logout: (jti: string) => Promise<void>;
	renewSession: (token: string, clientInfo: ClientInfo) => Promise<Auth>;
}

// ===============================================================
// Puerto del Controlador (interfaz HTTP)
// ===============================================================
export interface IAuthController {
	login: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
	logout: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
	renewSession: (
		req: Request,
		res: Response,
		next: NextFunction,
	) => Promise<Response | undefined>;
}
