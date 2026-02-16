import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

export default defineConfig({
	plugins: [tsconfigPaths({ projects: ["./tsconfig.test.json"] })],
	test: {
		environment: "node",
		include: ["src/**/*.integration.test.ts"],
		setupFiles: ["src/__tests__/setup/integration.setup.ts"],
	},
});
