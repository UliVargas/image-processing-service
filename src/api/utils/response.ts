import type { Response } from "express";

interface SuccessResponseOptions<T> {
	res: Response;
	data?: T;
	message?: string;
	statusCode?: number;
}

export const successResponse = <T>({
	res,
	data,
	message,
	statusCode = 200,
}: SuccessResponseOptions<T>) => {
	return res.status(statusCode).json({
		success: true,
		...(message && { message }),
		...(data && { data }),
		timestamp: new Date().toISOString(),
	});
};
