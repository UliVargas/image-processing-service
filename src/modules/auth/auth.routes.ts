import { makeInvoker } from "awilix-express";
import { Router } from "express";
import { authMiddleware, validateBody } from "@/api/middlewares/configured";
import { createAuthController } from "./auth.controller";
import { LoginDataSchema, RenewSessionDataSchema } from "./auth.schema";

const router: Router = Router();
const api = makeInvoker(createAuthController);

// ===============================================================
// Rutas de autenticación
// ===============================================================

// POST / - Iniciar sesión
router.post("/login", validateBody(LoginDataSchema), api("login"));

// Post /logout - Cerrar sesión
router.post("/logout", authMiddleware, api("logout"));

// POST /renew-session - Renovar sesión
router.post(
	"/renew-session",
	validateBody(RenewSessionDataSchema),
	api("renewSession"),
);

export default router;
