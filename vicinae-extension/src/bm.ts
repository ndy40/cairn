import { execFile } from "child_process";
import { promisify } from "util";

const execFileAsync = promisify(execFile);
const CAIRN_TIMEOUT_MS = 3000;
const CAIRN_MAX_BUFFER = 5 * 1024 * 1024;
const LIST_CACHE_TTL_MS = 2000;
let listCache: { ts: number; data: Bookmark[] } | null = null;

export interface Bookmark {
	ID: number;
	URL: string;
	Domain: string;
	Title: string;
	Description: string;
	CreatedAt: string;
	Tags: string[];
	LastVisitedAt: string | null;
	IsPermanent: boolean;
	IsArchived: boolean;
	ArchivedAt: string | null;
}

export async function bmAvailable(): Promise<boolean> {
	try {
		await execFileAsync("which", ["cairn"], {
			encoding: "utf8",
			timeout: 1000,
		});
		return true;
	} catch {
		return false;
	}
}

async function runCairn(args: string[]): Promise<{
	stdout: string;
	stderr: string;
	exitCode: number;
}> {
	try {
		const { stdout, stderr } = await execFileAsync("cairn", args, {
			encoding: "utf8",
			timeout: CAIRN_TIMEOUT_MS,
			maxBuffer: CAIRN_MAX_BUFFER,
		});
		return { stdout, stderr, exitCode: 0 };
	} catch (err) {
		const error = err as NodeJS.ErrnoException & {
			stdout?: string;
			stderr?: string;
			code?: number | string;
		};
		const exitCode = typeof error.code === "number" ? error.code : 3;
		return {
			stdout: typeof error.stdout === "string" ? error.stdout : "",
			stderr: typeof error.stderr === "string" ? error.stderr : "",
			exitCode,
		};
	}
}

export async function bmList(): Promise<Bookmark[]> {
	if (listCache && Date.now() - listCache.ts < LIST_CACHE_TTL_MS) {
		return listCache.data;
	}
	const result = await runCairn(["list", "--json"]);
	if (result.exitCode !== 0) {
		return [];
	}
	try {
		const parsed = JSON.parse(result.stdout);
		const data: Bookmark[] = Array.isArray(parsed) ? parsed : [];
		listCache = { ts: Date.now(), data };
		return data;
	} catch {
		return [];
	}
}

export async function bmSearch(query: string): Promise<Bookmark[]> {
	const result = await runCairn(["search", query, "--json", "--limit", "20"]);
	if (result.exitCode !== 0) {
		return [];
	}
	try {
		const parsed = JSON.parse(result.stdout);
		return Array.isArray(parsed) ? parsed : [];
	} catch {
		return [];
	}
}

export async function bmDelete(
	id: number,
): Promise<{ exitCode: number; stderr: string }> {
	const result = await runCairn(["delete", String(id)]);
	if (result.exitCode === 0) {
		listCache = null;
	}
	return {
		exitCode: result.exitCode,
		stderr: result.stderr,
	};
}

export async function bmPin(
	id: number,
): Promise<{ exitCode: number; stderr: string }> {
	const result = await runCairn(["pin", String(id)]);
	if (result.exitCode === 0) {
		listCache = null;
	}
	return { exitCode: result.exitCode, stderr: result.stderr };
}

export async function bmEdit(
	id: number,
	url?: string,
	tags?: string,
	title?: string,
): Promise<{ exitCode: number; stderr: string }> {
	const args = ["edit", String(id)];
	if (url !== undefined && url.trim() !== "") {
		args.push("--url", url);
	}
	if (tags !== undefined) {
		args.push("--tags", tags);
	}
	if (title !== undefined && title.trim() !== "") {
		args.push("--title", title.trim());
	}
	const result = await runCairn(args);
	if (result.exitCode === 0) {
		listCache = null;
	}
	return { exitCode: result.exitCode, stderr: result.stderr };
}

export async function bmAdd(
	url: string,
	tags?: string,
): Promise<{ exitCode: number; stderr: string }> {
	const args = ["add", url];
	if (tags && tags.trim() !== "") {
		args.push("--tags", tags);
	}
	const result = await runCairn(args);
	if (result.exitCode === 0 || result.exitCode === 2) {
		listCache = null;
	}
	return {
		exitCode: result.exitCode,
		stderr: result.stderr,
	};
}
