import { ENV_CONFIG } from "@shared/config/env.config";
import request from "supertest";
import { afterEach, beforeEach, expect, it } from "vitest";
import { app } from "@/app";

let server: ReturnType<typeof app.listen>;

beforeEach(() => {
	server = app.listen(ENV_CONFIG.PORT || 3000);
});

afterEach(() => {
	server.close();
});

it("should return a message indicating the service is running", async () => {
	const response = await request(server).get("/health");
	expect(response.status).toBe(200);
	expect(response.body).toEqual({
		status: "ok",
		timestamp: expect.any(String),
	});
});
