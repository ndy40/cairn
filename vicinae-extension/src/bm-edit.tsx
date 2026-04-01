import {
	Form,
	Action,
	ActionPanel,
	showToast,
	Toast,
	popToRoot,
} from "@vicinae/api";
import { useState } from "react";
import { bmEdit, Bookmark } from "./bm";

export default function EditBookmarkForm({ bookmark }: { bookmark: Bookmark }) {
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
