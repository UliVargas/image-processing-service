import dayjs from "dayjs";
import type { UserResponse } from "./users.dto";
import type { User } from "./users.entity";

// ===============================================================
// Mapea una entidad de dominio a un DTO de respuesta
// ===============================================================
export const toUserResponse = (user: User): UserResponse => {
	return {
		id: user.id,
		email: user.email,
		username: user.username,
		name: user.name,
		createdAt: dayjs(user.createdAt).format("YYYY-MM-DDTHH:mm:ss.SSS"),
		updatedAt: dayjs(user.updatedAt).format("YYYY-MM-DDTHH:mm:ss.SSS"),
	};
};

// ===============================================================
// Mapea un array de usuarios a DTOs de respuesta
// ===============================================================
export const toUserListResponse = (users: User[]): UserResponse[] => {
	return users.map(toUserResponse);
};
