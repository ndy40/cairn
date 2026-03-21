import {
	List,
	Action,
	ActionPanel,
	open,
	Clipboard,
	showHUD,
	confirmAlert,
	Alert,
	showToast,
	Toast,
} from "@vicinae/api";
import { useState, useEffect, useRef } from "react";
import { bmAvailable, bmList, bmSearch, bmDelete, Bookmark } from "./bm";

function formatDate(iso: string): string {
	return iso.slice(0, 10);
}

function BookmarkListItem({
	bookmark,
	onDelete,
}: {
	bookmark: Bookmark;
	onDelete: (bookmark: Bookmark) => void;
}) {
	const title = bookmark.Title || bookmark.URL;
	const accessories = [];

	accessories.push({ text: formatDate(bookmark.CreatedAt) });

	if (bookmark.Tags && bookmark.Tags.length > 0) {
		accessories.push({ text: bookmark.Tags.map((t) => `#${t}`).join(" ") });
	}

	if (bookmark.IsPermanent) {
		accessories.push({ text: "📌" });
	}

	return (
		<List.Item
			title={title}
			subtitle={bookmark.Domain}
			accessories={accessories}
			actions={
				<ActionPanel>
					<Action title="Open in Browser" onAction={() => open(bookmark.URL)} />
					<Action
						title="Copy URL"
						onAction={async () => {
							await Clipboard.copy(bookmark.URL);
							await showHUD("URL copied");
						}}
					/>
					<Action
						title="Delete Bookmark"
						onAction={async () => {
							const confirmed = await confirmAlert({
								title: "Delete Bookmark?",
								message: `"${bookmark.Title || bookmark.URL}" will be permanently deleted.`,
								primaryAction: {
									title: "Delete",
									style: Alert.ActionStyle.Destructive,
								},
							});
							if (!confirmed) return;

							const result = await bmDelete(bookmark.ID);
							if (result.exitCode === 0) {
								await showToast({
									style: Toast.Style.Success,
									title: "Bookmark deleted",
								});
								onDelete(bookmark);
							} else {
								const msg =
									result.exitCode === 1
										? "Bookmark not found"
										: result.stderr.trim() || "Delete failed";
								await showToast({ style: Toast.Style.Failure, title: msg });
							}
						}}
					/>
				</ActionPanel>
			}
		/>
	);
}

export default function SearchBookmarks() {
	const [query, setQuery] = useState("");
	const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [cliError, setCliError] = useState<string | null>(null);
	const [cliAvailable, setCliAvailable] = useState<boolean | null>(null);
	const requestIdRef = useRef(0);

	useEffect(() => {
		let active = true;
		(async () => {
			const available = await bmAvailable();
			if (!active) return;
			if (!available) {
				setCliError(
					"cairn is not installed. Install from: https://github.com/ndy40/bookmark-manager",
				);
				setCliAvailable(false);
				setIsLoading(false);
				return;
			}
			setCliError(null);
			setCliAvailable(true);
		})();
		return () => {
			active = false;
		};
	}, []);

	useEffect(() => {
		if (cliAvailable !== true) return;
		let cancelled = false;
		const requestId = (requestIdRef.current += 1);
		const timer = setTimeout(() => {
			(async () => {
				setIsLoading(true);
				const results =
					query.length >= 1 ? await bmSearch(query) : await bmList();
				if (cancelled || requestId !== requestIdRef.current) return;
				setBookmarks(results);
				setIsLoading(false);
			})();
		}, 250);
		return () => {
			cancelled = true;
			clearTimeout(timer);
		};
	}, [query, cliAvailable]);

	const handleDelete = async () => {
		if (cliAvailable !== true) return;
		const requestId = (requestIdRef.current += 1);
		setIsLoading(true);
		const results = query.length >= 1 ? await bmSearch(query) : await bmList();
		if (requestId !== requestIdRef.current) return;
		setBookmarks(results);
		setIsLoading(false);
	};

	if (cliError) {
		return (
			<List>
				<List.EmptyView title="cairn CLI not found" description={cliError} />
			</List>
		);
	}

	return (
		<List
			isLoading={isLoading}
			searchText={query}
			onSearchTextChange={setQuery}
			searchBarPlaceholder="Search bookmarks..."
		>
			{bookmarks.length === 0 && !isLoading ? (
				<List.EmptyView title="No bookmarks found" />
			) : (
				bookmarks.map((b) => (
					<BookmarkListItem key={b.ID} bookmark={b} onDelete={handleDelete} />
				))
			)}
		</List>
	);
}
