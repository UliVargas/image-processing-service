import { container } from "@/container";

const authMiddleware = container.resolve("authMiddleware");
const validateMiddleware = container.resolve("validateMiddleware");

export { authMiddleware, validateMiddleware };
export {
	validateBody,
	validateParams,
	validateQuery,
} from "./validate.helpers";
