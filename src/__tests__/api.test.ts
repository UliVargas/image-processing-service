import request from "supertest";
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { app } from "@/app";

let server: ReturnType<typeof app.listen>;

describe("API Health", () => {
	beforeEach(() => {
		server = app.listen(0);
	});

	afterEach(() => {
		server.close();
	});

	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should return a message indicating the service is running", async () => {
		const response = await request(server).get("/health");
		expect(response.status).toBe(200);
		expect(response.body).toEqual({
			status: "ok",
			timestamp: expect.any(String),
		});
	});
});
