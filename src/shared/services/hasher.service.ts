import crypto from "node:crypto";

// ===============================================================
// Servicio de hashing para seguridad (contraseñas, tokens, etc.)
// ===============================================================
export const createHasherService = () => ({
	createHash(data: string): string {
		return crypto.createHash("sha256").update(data).digest("hex");
	},
	validarHash(dataOriginal: string, hashAlmacenado: string): boolean {
		const nuevoHash = this.createHash(dataOriginal);
		const bufferNuevo: Buffer = Buffer.from(nuevoHash, "utf8");
		const bufferAlmacenado: Buffer = Buffer.from(hashAlmacenado, "utf8");
		if (bufferNuevo.length !== bufferAlmacenado.length) {
			return false;
		}
		return crypto.timingSafeEqual(bufferNuevo, bufferAlmacenado);
	},
});
