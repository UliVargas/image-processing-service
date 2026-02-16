import type { Session } from "../session/session.entity";

// ===============================================================
// Entidad de Dominio (independiente de infraestructura)
// ===============================================================
export interface Auth {
	session: Session;
	tokens: {
		accessToken: string;
		refreshToken: string;
	};
}
