// ===============================================================
// Objetos de Transferencia de Datos (Capa de Presentación)
// ===============================================================
export interface SessionResponse {
	id: string;
	tokenHash: string;
	userId: string;
	deviceInfo: string | null;
	ipAddress: string | null;
	userAgent: string | null;
	expiresAt: Date;
	lastUsedAt: Date | null;
	isValid: boolean;
	createdAt: Date;
}

export type UserListResponse = SessionResponse[];
