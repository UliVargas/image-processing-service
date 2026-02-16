// ===============================================================
// Objetos de Transferencia de Datos (Capa de Presentación)
// ===============================================================
export interface UserResponse {
	id: string;
	email: string;
	username: string;
	name: string | null;
	createdAt: string;
	updatedAt: string;
}

export type UserListResponse = UserResponse[];
