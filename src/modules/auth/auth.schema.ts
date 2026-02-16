import * as v from "valibot";

// ===============================================================
// Login Data (para validación de requests)
// ===============================================================

/* Se usa union para que se pueda recibir email o username */
export const LoginDataSchema = v.union([
	v.object({
		email: v.pipe(v.string(), v.email()),
		password: v.string(),
	}),
	v.object({
		username: v.string(),
		password: v.string(),
	}),
]);

export type LoginData = v.InferOutput<typeof LoginDataSchema>;

// ===============================================================
// Renew Session Data (para validación de requests)
// ===============================================================
export const RenewSessionDataSchema = v.object({
	refreshToken: v.pipe(v.string(), v.cuid2("El refreshToken debe ser válido")),
});

export type RenewSessionData = v.InferOutput<typeof RenewSessionDataSchema>;
