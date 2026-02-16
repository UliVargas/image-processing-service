// ===============================================================
// Utilidades para detección de dispositivo y navegador
// ===============================================================

export const getDeviceType = (userAgent: string): string => {
	const ua = userAgent.toLowerCase();
	if (ua.includes("mobile") || ua.includes("android")) return "mobile";
	if (ua.includes("iphone") || ua.includes("ipad")) return "ios";
	if (ua.includes("tablet")) return "tablet";
	return "desktop";
};

export const getBrowser = (userAgent: string): string => {
	const ua = userAgent.toLowerCase();
	if (ua.includes("edg")) return "edge";
	if (ua.includes("chrome")) return "chrome";
	if (ua.includes("firefox")) return "firefox";
	if (ua.includes("safari") && !ua.includes("chrome")) return "safari";
	if (ua.includes("opera") || ua.includes("opr")) return "opera";
	return "other";
};
