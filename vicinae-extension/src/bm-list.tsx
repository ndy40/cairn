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
import { useState, useEffect } from "react";
import { bmAvailable, bmList, bmDelete, bmPin, Bookmark } from "./bm";

function formatDate(iso: string): string {
	return iso.slice(0, 10);
}

function BookmarkListItem({
	bookmark,
	onDelete,
	onPin,
}: {
	bookmark: Bookmark;
	onDelete: (bookmark: Bookmark) => void;
	onPin: () => void;
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
						title={bookmark.IsPermanent ? "Unpin Bookmark" : "Pin Bookmark"}
						onAction={async () => {
							const result = await bmPin(bookmark.ID);
							if (result.exitCode === 0) {
								await showToast({
									style: Toast.Style.Success,
									title: bookmark.IsPermanent ? "Bookmark unpinned" : "Bookmark pinned",
								});
								onPin();
							} else {
								const msg =
									result.exitCode === 1
										? "Bookmark not found"
										: result.stderr.trim() || "Pin failed";
								await showToast({ style: Toast.Style.Failure, title: msg });
							}
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

export default function ListBookmarks() {
	const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [cliError, setCliError] = useState<string | null>(null);

	useEffect(() => {
		let active = true;
		(async () => {
			const available = await bmAvailable();
			if (!active) return;
			if (!available) {
				setCliError(
					"cairn is not installed. Install from: https://github.com/ndy40/bookmark-manager",
				);
				setIsLoading(false);
				return;
			}
			const results = await bmList();
			if (!active) return;
			setBookmarks(results);
			setIsLoading(false);
		})();
		return () => {
			active = false;
		};
	}, []);

	const handleDelete = async () => {
		setIsLoading(true);
		const results = await bmList();
		setBookmarks(results);
		setIsLoading(false);
	};

	const handlePin = async () => {
		setIsLoading(true);
		const results = await bmList();
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
		<List isLoading={isLoading} searchBarPlaceholder="Filter bookmarks...">
			{bookmarks.length === 0 && !isLoading ? (
				<List.EmptyView title="No bookmarks saved yet" />
			) : (
				bookmarks.map((b) => (
					<BookmarkListItem key={b.ID} bookmark={b} onDelete={handleDelete} onPin={handlePin} />
				))
			)}
		</List>
	);
}
