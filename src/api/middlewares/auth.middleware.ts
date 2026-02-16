import type { NextFunction, Request, Response } from "express";
import {
	JsonWebTokenError,
	type JwtPayload,
	TokenExpiredError,
} from "jsonwebtoken";
import type { Cradle } from "@/shared/di/types";
import { createError } from "@/shared/errors/app-error";

interface Dependencies {
	tokenManager: Cradle["tokenManager"];
	config: Cradle["config"];
}

export const createAuthMiddleware =
	({ tokenManager, config }: Dependencies) =>
	(req: Request, _res: Response, next: NextFunction): void => {
		try {
			const token = req.headers.authorization?.split(" ")[1];
			if (!token) {
				throw createError({
					message: "Token no proporcionado",
					statusCode: 401,
					code: "TOKEN_NOT_PROVIDED",
				});
			}
			const decoded = tokenManager.verify(token, config.accessTokenSecret);
			req.userId = (decoded as JwtPayload).sub as string;
			req.jti = (decoded as JwtPayload).jti as string;
			next();
		} catch (error) {
			if (error instanceof TokenExpiredError) {
				next(
					createError({
						message:
							"El token ha expirado. Por favor, inicia sesión nuevamente",
						statusCode: 401,
						code: "TOKEN_EXPIRED",
					}),
				);
				return;
			}
			if (error instanceof JsonWebTokenError) {
				next(
					createError({
						message: "Token inválido o mal formado",
						statusCode: 401,
						code: "INVALID_TOKEN",
					}),
				);
				return;
			}
			next(error);
		}
	};
