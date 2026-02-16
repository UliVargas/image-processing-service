import type { NextFunction, Request, Response } from "express";
import * as v from "valibot";
import { createError } from "../../shared/errors/app-error";

export type ValidationType = "body" | "query" | "params";

export interface ValidateOptions<T extends v.GenericSchema> {
	schema: T;
	type?: ValidationType;
}

export const validate = <T extends v.GenericSchema>({
	schema,
	type = "body",
}: ValidateOptions<T>) => {
	return (req: Request, _res: Response, next: NextFunction) => {
		const result = v.safeParse(schema, req[type]);

		if (!result.success) {
			const errors: Record<string, string> = {};

			for (const issue of result.issues) {
				const path = issue.path?.map((p) => p.key).join(".") || "general";

				// Personalizar mensaje cuando falta un campo en el objeto
				if (issue.type === "object" && issue.message.includes("Invalid key")) {
					issue.path?.[0]?.key || "campo";
					errors[path] = `Este campo es obligatorio`;
				} else {
					errors[path] = issue.message;
				}
			}

			return next(
				createError({
					message: "Error de validación",
					statusCode: 400,
					code: "VALIDATION_ERROR",
					details: errors,
				}),
			);
		}

		req[type] = result.output;
		next();
	};
};
