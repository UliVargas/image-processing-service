import { Router } from "express";
import authRoutes from "../modules/auth/auth.routes";
import userRoutes from "../modules/users/users.routes";

export const router: Router = Router();

router.use("/auth", authRoutes);
router.use("/users", userRoutes);
