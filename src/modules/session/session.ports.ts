import type { Session } from "./session.entity";
import type { CreateSessionData, RenewSessionData } from "./session.schema";

// ===============================================================
// Puerto del Repositorio (abstracción de persistencia)
// ===============================================================
export interface ISessionRepository {
	createSession: (data: CreateSessionData) => Promise<Session>;
	deleteSession: (id: string) => Promise<void>;
	renewSession: (id: string, data: RenewSessionData) => Promise<Session>;
	findSessionByToken: (token: string) => Promise<Session | null>;
}
