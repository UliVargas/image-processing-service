import { execSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { PrismaBetterSqlite3 } from "@prisma/adapter-better-sqlite3";

const testDbDir = path.join(process.cwd(), "src", "__tests__", "prisma");
const testDbPath = path.join(testDbDir, "test.sqlite");
const testDbUrl = `file:${testDbPath}`;

process.env.TEST_DATABASE_URL = process.env.TEST_DATABASE_URL || testDbUrl;
process.env.DATABASE_URL = process.env.TEST_DATABASE_URL;

let prismaClient: unknown;

const run = (command: string) => {
	execSync(command, {
		cwd: process.cwd(),
		stdio: "pipe",
		env: {
			...process.env,
			TEST_DATABASE_URL: process.env.TEST_DATABASE_URL,
			DATABASE_URL: process.env.TEST_DATABASE_URL,
		},
	});
};

export const setupSqlitePrisma = async () => {
	if (!fs.existsSync(testDbDir)) {
		fs.mkdirSync(testDbDir, { recursive: true });
	}

	if (fs.existsSync(testDbPath)) {
		fs.rmSync(testDbPath);
	}

	run("pnpm prisma generate --schema src/__tests__/prisma/schema.prisma");
	run("pnpm prisma db push --schema src/__tests__/prisma/schema.prisma");

	const prismaModulePath = path.join(
		process.cwd(),
		"src/generated/prisma-test/client.ts",
	);
	const prismaModule = await import(prismaModulePath);
	const PrismaClient = prismaModule.PrismaClient as new (args: {
		adapter: PrismaBetterSqlite3;
	}) => {
		$disconnect: () => Promise<void>;
		session: { deleteMany: () => Promise<unknown> };
		user: { deleteMany: () => Promise<unknown> };
	};

	const adapter = new PrismaBetterSqlite3({
		url: process.env.TEST_DATABASE_URL as string,
	});

	prismaClient = new PrismaClient({ adapter });
};

export const getSqlitePrisma = () => {
	if (!prismaClient) {
		throw new Error("SQLite Prisma client is not initialized");
	}
	return prismaClient;
};

export const resetSqlitePrisma = async () => {
	const client = getSqlitePrisma() as {
		session: { deleteMany: () => Promise<unknown> };
		user: { deleteMany: () => Promise<unknown> };
	};

	await client.session.deleteMany();
	await client.user.deleteMany();
};

export const teardownSqlitePrisma = async () => {
	if (!prismaClient) return;

	const client = prismaClient as {
		$disconnect: () => Promise<void>;
	};
	await client.$disconnect();
	prismaClient = undefined;

	if (fs.existsSync(testDbPath)) {
		fs.rmSync(testDbPath);
	}
};
