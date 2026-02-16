import { makeInvoker } from "awilix-express";
import { Router } from "express";
import {
	authMiddleware,
	validateBody,
	validateParams,
} from "@/api/middlewares/configured";
import { createUsersController } from "./users.controller";
import {
	CreateUserInputSchema,
	UpdateUserInputSchema,
	UserParamsSchema,
} from "./users.schema";

const router: Router = Router();
const api = makeInvoker(createUsersController);

// ===============================================================
// Rutas de usuarios
// ===============================================================

// POST / - Crear usuario
router.post("/", validateBody(CreateUserInputSchema), api("createUser"));

// Middleware de autenticación aplicado para rutas protegidas
router.use(authMiddleware);

// GET / - Obtener todos los usuarios
router.get("/", api("getAllUsers"));

// PATCH /:id - Actualizar usuario
router.patch(
	"/:id",
	validateParams(UserParamsSchema),
	validateBody(UpdateUserInputSchema),
	api("updateUser"),
);

// DELETE /:id - Eliminar usuario
router.delete("/:id", validateParams(UserParamsSchema), api("deleteUser"));

export default router;
