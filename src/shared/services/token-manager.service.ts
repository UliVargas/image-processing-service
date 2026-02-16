import jwt, { type SignOptions } from "jsonwebtoken";

// ===============================================================
// Puerto del Servicio de gestión de tokens (JWT)
// ===============================================================
export interface ITokenManager {
	sign(
		payload: string | object | Buffer,
		secret: string,
		options?: SignOptions,
	): string;
	verify(token: string, secret: string): string | object;
}

export const createTokenManager = (): ITokenManager => ({
	sign: (payload, secret, options) => {
		return jwt.sign(payload, secret, options);
	},
	verify: (token, secret) => {
		return jwt.verify(token, secret);
	},
});
