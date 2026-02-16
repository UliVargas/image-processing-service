// ===============================================================
// Entidad de Dominio (independiente de infraestructura)
// ===============================================================
export interface User {
	id: string;
	email: string;
	username: string;
	name: string | null;
	password: string;
	createdAt: Date;
	updatedAt: Date;
}
