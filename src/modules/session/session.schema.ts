import * as v from "valibot";

// ===============================================================
// Session Schema Base (para uso interno)
// ===============================================================
export const SessionSchema = v.object({
	id: v.string(),
	tokenHash: v.string(),
	accessJti: v.string(),
	userId: v.string(),
	deviceInfo: v.string(),
	ipAddress: v.string(),
	userAgent: v.string(),
	expiresAt: v.date(),
	lastUsedAt: v.date(),
	isValid: v.boolean(),
	createdAt: v.date(),
});

// ===============================================================
// Session Data (para crear una sesión en repositorio - incluye ID generado)
// ===============================================================
export const CreateSessionDataSchema = v.object({
	id: v.string(),
	tokenHash: v.string(),
	accessJti: v.string(),
	userId: v.string(),
	deviceInfo: v.string(),
	ipAddress: v.string(),
	userAgent: v.string(),
	expiresAt: v.date(),
});

export type CreateSessionData = v.InferOutput<typeof CreateSessionDataSchema>;

// ===============================================================
// Session Data (para renovar una sesión en repositorio)
// ===============================================================
export const RenewSessionDataSchema = v.object({
	tokenHash: v.string(),
	accessJti: v.string(),
	expiresAt: v.date(),
	lastUsedAt: v.date(),
	ipAddress: v.string(),
	userAgent: v.string(),
	deviceInfo: v.string(),
});

export type RenewSessionData = v.InferOutput<typeof RenewSessionDataSchema>;
