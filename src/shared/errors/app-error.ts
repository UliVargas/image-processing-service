export interface AppError extends Error {
	statusCode: number;
	code: string;
	details?: Record<string, string | Record<string, string>>;
}

interface CreateErrorOptions {
	message: string;
	statusCode: number;
	code: string;
	details?: Record<string, string | Record<string, string>>;
}

// ===============================================================
// Función para crear errores personalizados de la aplicación
// ===============================================================
export const createError = ({
	message,
	statusCode,
	code,
	details,
}: CreateErrorOptions): AppError => {
	const error = new Error(message) as AppError;
	error.statusCode = statusCode;
	error.code = code;
	error.details = details;
	return error;
};

export const isAppError = (error: unknown): error is AppError => {
	return error instanceof Error && "statusCode" in error && "code" in error;
};
