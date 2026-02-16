import type { NextFunction, Request, Response } from "express";

// ===============================================================
// Obtiene la IP real del cliente considerando proxies
// ===============================================================
const getClientIp = (req: Request): string => {
	const xRealIp = req.headers["x-real-ip"];
	if (xRealIp && typeof xRealIp === "string") {
		return xRealIp;
	}

	const xForwardedFor = req.headers["x-forwarded-for"];
	if (xForwardedFor) {
		const ips = (
			typeof xForwardedFor === "string" ? xForwardedFor : xForwardedFor[0]
		).split(",");
		return ips[0].trim();
	}

	const cfConnectingIp = req.headers["cf-connecting-ip"];
	if (cfConnectingIp && typeof cfConnectingIp === "string") {
		return cfConnectingIp;
	}

	return req.ip || "unknown";
};

// ===============================================================
// Obtiene el User-Agent del cliente
// ===============================================================
const getUserAgent = (req: Request): string => {
	return req.headers["user-agent"] || "unknown";
};

// ===============================================================
// Extrae información básica del dispositivo desde el User-Agent
// ===============================================================
const getDeviceInfo = (req: Request): string => {
	const userAgent = req.headers["user-agent"] || "";

	let browser = "Unknown";
	if (userAgent.includes("Chrome")) browser = "Chrome";
	else if (userAgent.includes("Firefox")) browser = "Firefox";
	else if (userAgent.includes("Safari")) browser = "Safari";
	else if (userAgent.includes("Edge")) browser = "Edge";
	else if (userAgent.includes("Opera")) browser = "Opera";

	let os = "Unknown";
	if (userAgent.includes("Windows")) os = "Windows";
	else if (userAgent.includes("Mac OS")) os = "macOS";
	else if (userAgent.includes("Linux")) os = "Linux";
	else if (userAgent.includes("Android")) os = "Android";
	else if (
		userAgent.includes("iOS") ||
		userAgent.includes("iPhone") ||
		userAgent.includes("iPad")
	)
		os = "iOS";

	let device = "Desktop";
	if (userAgent.includes("Mobile")) device = "Mobile";
	else if (userAgent.includes("Tablet")) device = "Tablet";

	return `${browser} on ${os} (${device})`;
};

// ===============================================================
// Middleware que agrega información del cliente al request
// ===============================================================
export const clientInfoMiddleware = (
	req: Request,
	_res: Response,
	next: NextFunction,
): void => {
	req.clientInfo = {
		ipAddress: getClientIp(req),
		userAgent: getUserAgent(req),
		deviceInfo: getDeviceInfo(req),
	};
	next();
};
