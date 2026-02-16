import dayjs from "dayjs";
import type { AuthResponse } from "./auth.dto";
import type { Auth } from "./auth.entity";

// ===============================================================
// Mapea una entidad de dominio a un DTO de respuesta
// ===============================================================
export const toAuthResponse = (auth: Auth): AuthResponse => {
	return {
		session: {
			sessionId: auth.session.id,
			expiresAt: dayjs(auth.session.expiresAt).format(
				"YYYY-MM-DDTHH:mm:ss.SSS",
			),
		},
		tokens: {
			accessToken: auth.tokens.accessToken,
			refreshToken: auth.tokens.refreshToken,
		},
	};
};
