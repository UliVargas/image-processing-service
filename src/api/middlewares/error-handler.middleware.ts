import { isAppError } from "@shared/errors/app-error";
import type { NextFunction, Request, Response } from "express";
import * as v from "valibot";

interface ErrorResponse {
	success: false;
	error: {
		message: string;
		code: string;
		details?: Record<string, string | Record<string, string>>;
	};
	timestamp: string;
	path: string;
}

export const errorHandler = (
	err: Error,
	req: Request,
	res: Response,
	_next: NextFunction,
): void => {
	// Error de aplicación
	if (isAppError(err)) {
		const response: ErrorResponse = {
			success: false,
			error: {
				message: err.message,
				code: err.code,
				details: err.details,
			},
			timestamp: new Date().toISOString(),
			path: req.path,
		};

		res.status(err.statusCode).json(response);
		return;
	}

	// Error de validación de Valibot
	if (v.isValiError(err)) {
		const errors: Record<string, string> = {};
		for (const issue of err.issues) {
			const path = issue.path?.map((p) => p.key).join(".") || "general";
			errors[path] = issue.message;
		}

		const response: ErrorResponse = {
			success: false,
			error: {
				message: "Error de validación",
				code: "VALIDATION_ERROR",
				details: errors,
			},
			timestamp: new Date().toISOString(),
			path: req.path,
		};

		res.status(400).json(response);
		return;
	}

	// Error desconocido (500)
	console.error("Error no manejado:", err);

	const response: ErrorResponse = {
		success: false,
		error: {
			message:
				process.env.NODE_ENV === "production"
					? "Error interno del servidor"
					: err.message,
			code: "INTERNAL_SERVER_ERROR",
		},
		timestamp: new Date().toISOString(),
		path: req.path,
	};

	res.status(500).json(response);
};
