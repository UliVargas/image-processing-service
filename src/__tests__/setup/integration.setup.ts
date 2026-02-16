import { afterAll, beforeAll, beforeEach } from "vitest";
import {
	resetSqlitePrisma,
	setupSqlitePrisma,
	teardownSqlitePrisma,
} from "./sqlite-prisma";

beforeAll(async () => {
	await setupSqlitePrisma();
});

beforeEach(async () => {
	await resetSqlitePrisma();
});

afterAll(async () => {
	await teardownSqlitePrisma();
});
