declare global {
	namespace Express {
		interface Request {
			userId?: string;
			jti?: string;
			clientInfo?: {
				ipAddress: string;
				userAgent: string;
				deviceInfo: string;
			};
		}
	}
}

export {};
