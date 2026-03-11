import { spawnSync } from "child_process";

export interface Bookmark {
	id: number;
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

export function bmAvailable(): boolean {
	const result = spawnSync("which", ["cairn"], { encoding: "utf8" });
	return result.status === 0;
}

export function bmList(): Bookmark[] {
	const result = spawnSync("cairn", ["list", "--json"], { encoding: "utf8" });
	if (result.status !== 0) {
		return [];
	}
	try {
		return JSON.parse(result.stdout) as Bookmark[];
	} catch {
		return [];
	}
}

export function bmSearch(query: string): Bookmark[] {
	const result = spawnSync(
		"cairn",
		["search", query, "--json", "--limit", "20"],
		{
			encoding: "utf8",
		},
	);
	if (result.status !== 0) {
		return [];
	}
	try {
		return JSON.parse(result.stdout) as Bookmark[];
	} catch {
		return [];
	}
}

export function bmDelete(id: number): { exitCode: number; stderr: string } {
	const result = spawnSync("cairn", ["delete", String(id)], {
		encoding: "utf8",
	});
	return {
		exitCode: result.status ?? 3,
		stderr: result.stderr ?? "",
	};
}

export function bmAdd(
	url: string,
	tags?: string,
): { exitCode: number; stderr: string } {
	const args = ["add", url];
	if (tags && tags.trim() !== "") {
		args.push("--tags", tags);
	}
	const result = spawnSync("cairn", args, { encoding: "utf8" });
	return {
		exitCode: result.status ?? 3,
		stderr: result.stderr ?? "",
	};
}
