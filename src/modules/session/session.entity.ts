// ===============================================================
// Entidad de Dominio (independiente de infraestructura)
// ===============================================================
export interface Session {
	id: string;
	tokenHash: string;
	accessJti: string;
	userId: string;
	deviceInfo: string | null;
	ipAddress: string | null;
	userAgent: string | null;
	expiresAt: Date;
	lastUsedAt: Date | null;
	isValid: boolean;
	createdAt: Date;
}
