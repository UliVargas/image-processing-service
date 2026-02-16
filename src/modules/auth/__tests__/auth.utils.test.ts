import { describe, expect, it } from "vitest";
import { getBrowser, getDeviceType } from "../auth.utils";

describe("auth.utils", () => {
	// ===============================================================
	// Fase GREEN (flujo esperado)
	// ===============================================================
	it("should detect device types", () => {
		expect(getDeviceType("Mozilla/5.0 (Linux; Android 10; Mobile)")).toBe(
			"mobile",
		);
		expect(
			getDeviceType("Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)"),
		).toBe("ios");
		expect(getDeviceType("Mozilla/5.0 (Tablet; rv:109.0)")).toBe("tablet");
		expect(getDeviceType("Mozilla/5.0 (Windows NT 10.0; Win64; x64)")).toBe(
			"desktop",
		);
	});

	it("should detect browsers", () => {
		expect(getBrowser("Mozilla/5.0 Edg/122.0.0.0")).toBe("edge");
		expect(getBrowser("Mozilla/5.0 Chrome/122.0.0.0 Safari/537.36")).toBe(
			"chrome",
		);
		expect(getBrowser("Mozilla/5.0 Firefox/123.0")).toBe("firefox");
		expect(getBrowser("Mozilla/5.0 Version/17.0 Safari/605.1.15")).toBe(
			"safari",
		);
		expect(getBrowser("Mozilla/5.0 OPR/105.0.0.0")).toBe("opera");
		expect(getBrowser("unknown-agent")).toBe("other");
	});
});
