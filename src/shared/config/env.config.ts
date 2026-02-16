import "dotenv/config";

// ===============================================================
// Configuración de variables de entorno
// ===============================================================
const ENV_CONFIG = {
	PORT: process.env.PORT || "3000",
	NODE_ENV: process.env.NODE_ENV || "development",
	DATABASE_URL: process.env.DATABASE_URL || "",
	ACCESS_TOKEN_SECRET:
		process.env.ACCESS_TOKEN_SECRET || "your-default-jwt-secret",
	REFRESH_TOKEN_SECRET:
		process.env.REFRESH_TOKEN_SECRET || "your-default-refresh-jwt-secret",
	ACCESS_TOKEN_EXPIRATION_TIME:
		process.env.ACCESS_TOKEN_EXPIRATION_TIME || "15m",
	REFRESH_TOKEN_EXPIRATION_TIME:
		process.env.REFRESH_TOKEN_EXPIRATION_TIME || "7d",
	LOG_LEVEL: process.env.LOG_LEVEL || "info",
	SALT_ROUNDS: parseInt(process.env.SALT_ROUNDS || "10", 10),
};

export { ENV_CONFIG };
