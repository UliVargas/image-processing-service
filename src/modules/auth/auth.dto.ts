// ===============================================================
// Objetos de Transferencia de Datos (Capa de Presentación)
// ===============================================================
export interface AuthResponse {
	session: {
		sessionId: string;
		expiresAt: string;
	};
	tokens: {
		accessToken: string;
		refreshToken: string;
	};
}
