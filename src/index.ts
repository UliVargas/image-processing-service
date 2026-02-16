import { ENV_CONFIG } from "@shared/config/env.config";
import { app } from "./app";

app.listen(ENV_CONFIG.PORT, () => {
	console.log(`Server is running on http://localhost:${ENV_CONFIG.PORT}`);
	console.log(`🏥 Health: http://localhost:${ENV_CONFIG.PORT}/health`);
});
