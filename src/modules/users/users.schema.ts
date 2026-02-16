import * as v from "valibot";

// ===============================================================
// User Schema Base (para uso interno)
// ===============================================================
export const UserSchema = v.object({
	id: v.string(),
	email: v.pipe(v.string(), v.email()),
	username: v.pipe(v.string(), v.minLength(1), v.maxLength(50)),
	password: v.string(),
	name: v.nullable(v.string()),
	createdAt: v.date(),
	updatedAt: v.date(),
});

// ===============================================================
// User Input Schemas (para validación de requests)
// ===============================================================
export const CreateUserInputSchema = v.object({
	email: v.pipe(
		v.string("El correo electrónico es obligatorio"),
		v.nonEmpty("El correo electrónico no puede estar vacío"),
		v.email("Debes proporcionar un correo electrónico válido"),
		v.maxLength(255, "No puede exceder 255 caracteres"),
	),
	username: v.pipe(
		v.string("El nombre de usuario es obligatorio"),
		v.nonEmpty("El nombre de usuario no puede estar vacío"),
		v.minLength(3, "Debe tener al menos 3 caracteres"),
		v.maxLength(50, "No puede exceder 50 caracteres"),
		v.regex(
			/^[a-zA-Z0-9_-]+$/,
			"Solo puede contener letras, números, guiones y guiones bajos",
		),
	),
	password: v.pipe(
		v.string("La contraseña es obligatoria"),
		v.nonEmpty("La contraseña no puede estar vacía"),
		v.minLength(8, "Debe tener al menos 8 caracteres"),
		v.maxLength(100, "No puede exceder 100 caracteres"),
	),
	name: v.optional(
		v.pipe(v.string(), v.maxLength(100, "No puede exceder 100 caracteres")),
	),
});

export type CreateUserInput = v.InferOutput<typeof CreateUserInputSchema>;

// ===============================================================
// User Data (para crear en repositorio - incluye ID generado)
// ===============================================================
export const CreateUserDataSchema = v.object({
	id: v.string(),
	email: v.pipe(v.string(), v.email()),
	username: v.string(),
	password: v.string(),
	name: v.optional(v.string()),
});

export type CreateUserData = v.InferOutput<typeof CreateUserDataSchema>;

// ===============================================================
// User Id Schema (Para validar el id)
// ===============================================================
export const UserParamsSchema = v.object({
	id: v.pipe(v.string(), v.cuid2("El ID debe ser válido")),
});

export type UserParamsData = v.InferOutput<typeof UserParamsSchema>;

// ===============================================================
// User Update Schema (para validación de requests)
// ===============================================================
export const UpdateUserInputSchema = v.object({
	name: v.optional(
		v.pipe(
			v.string(),
			v.maxLength(100, "El nombre no puede exceder 100 caracteres"),
		),
	),
});

export type UpdateUserInput = v.InferOutput<typeof UpdateUserInputSchema>;
