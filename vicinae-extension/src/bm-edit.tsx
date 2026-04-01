import {
	Form,
	List,
	Action,
	ActionPanel,
	showToast,
	Toast,
	popToRoot,
	useNavigation,
} from "@vicinae/api";
import { useState, useEffect } from "react";
import { bmAvailable, bmEdit, bmList, Bookmark } from "./bm";

function EditBookmarkForm({ bookmark }: { bookmark: Bookmark }) {
	const [urlValue, setUrlValue] = useState(bookmark.URL);
	const [tagsValue, setTagsValue] = useState(
		bookmark.Tags ? bookmark.Tags.join(", ") : "",
	);
	const [urlError, setUrlError] = useState<string | undefined>(undefined);
	const [isSubmitting, setIsSubmitting] = useState(false);

	function validateUrl(value: string): boolean {
		if (!value || value.trim() === "") {
			setUrlError("URL is required");
			return false;
		}
		setUrlError(undefined);
		return true;
	}

	async function handleSubmit(values: { url: string; tags: string }) {
		if (!validateUrl(values.url)) {
			return;
		}

		setIsSubmitting(true);

		const urlChanged = values.url !== bookmark.URL;
		const tagsChanged = values.tags !== (bookmark.Tags ? bookmark.Tags.join(", ") : "");

		if (!urlChanged && !tagsChanged) {
			await showToast({ style: Toast.Style.Success, title: "No changes" });
			popToRoot();
			return;
		}

		const { exitCode, stderr } = await bmEdit(
			bookmark.ID,
			urlChanged ? values.url : undefined,
			tagsChanged ? values.tags : undefined,
		);
		setIsSubmitting(false);

		if (exitCode === 0) {
			await showToast({ style: Toast.Style.Success, title: "Bookmark updated" });
			popToRoot();
		} else if (exitCode === 1) {
			if (stderr.includes("Duplicate URL")) {
				setUrlError("Duplicate URL");
			} else {
				setUrlError("Bookmark not found");
			}
		} else {
			setUrlError(stderr.trim() || "Failed to update bookmark");
		}
	}

	return (
		<Form
			isLoading={isSubmitting}
			actions={
				<ActionPanel>
					<Action.SubmitForm title="Save Changes" onSubmit={handleSubmit} />
				</ActionPanel>
			}
		>
			<Form.TextField
				id="url"
				title="URL"
				placeholder="https://example.com"
				value={urlValue}
				onChange={setUrlValue}
				error={urlError}
				onBlur={(e) => validateUrl(e.target.value)}
			/>
			<Form.TextField
				id="tags"
				title="Tags"
				placeholder="work, go, tools  (comma-separated, max 3)"
				value={tagsValue}
				onChange={setTagsValue}
			/>
		</Form>
	);
}

export { EditBookmarkForm };

function formatDate(iso: string): string {
	return iso.slice(0, 10);
}

export default function EditBookmark() {
	const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [cliError, setCliError] = useState<string | null>(null);
	const { push } = useNavigation();

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

	if (cliError) {
		return (
			<List>
				<List.EmptyView title="cairn CLI not found" description={cliError} />
			</List>
		);
	}

	return (
		<List isLoading={isLoading} searchBarPlaceholder="Select a bookmark to edit...">
			{bookmarks.length === 0 && !isLoading ? (
				<List.EmptyView title="No bookmarks saved yet" />
			) : (
				bookmarks.map((b) => {
					const accessories = [];
					accessories.push({ text: formatDate(b.CreatedAt) });
					if (b.Tags && b.Tags.length > 0) {
						accessories.push({ text: b.Tags.map((t) => `#${t}`).join(" ") });
					}
					return (
						<List.Item
							key={b.ID}
							title={b.Title || b.URL}
							subtitle={b.Domain}
							accessories={accessories}
							actions={
								<ActionPanel>
									<Action
										title="Edit Bookmark"
										onAction={() => push(<EditBookmarkForm bookmark={b} />)}
									/>
								</ActionPanel>
							}
						/>
					);
				})
			)}
		</List>
	);
}
