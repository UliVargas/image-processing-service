import { errorHandler } from "@api/middlewares/error-handler.middleware";
import { scopePerRequest } from "awilix-express";
import compression from "compression";
import cors from "cors";
import express, { type Express } from "express";
import helmet from "helmet";
import morgan from "morgan";
import { container } from "@/container";
import { clientInfoMiddleware } from "./api/middlewares/client-info.middleware";
import { router } from "./api/routes";

const app: Express = express();

app.use(helmet());
app.use(cors());
app.use(compression());
app.use(morgan("dev"));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(scopePerRequest(container));
app.use(clientInfoMiddleware);

app.get("/health", (_req, res) => {
	res.json({ status: "ok", timestamp: new Date().toISOString() });
});

app.use("/api", router);

app.use(errorHandler);

export { app };
