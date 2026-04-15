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
	const [titleValue, setTitleValue] = useState(bookmark.Title);
	const [urlValue, setUrlValue] = useState(bookmark.URL);
	const [tagsValue, setTagsValue] = useState(
		bookmark.Tags ? bookmark.Tags.join(", ") : "",
	);
	const [titleError, setTitleError] = useState<string | undefined>(undefined);
	const [urlError, setUrlError] = useState<string | undefined>(undefined);
	const [isSubmitting, setIsSubmitting] = useState(false);

	function validateTitle(value: string): boolean {
		if (!value || value.trim() === "") {
			setTitleError("Title cannot be empty");
			return false;
		}
		setTitleError(undefined);
		return true;
	}

	function validateUrl(value: string): boolean {
		if (!value || value.trim() === "") {
			setUrlError("URL is required");
			return false;
		}
		setUrlError(undefined);
		return true;
	}

	async function handleSubmit(values: {
		title: string;
		url: string;
		tags: string;
	}) {
		if (!validateTitle(values.title) || !validateUrl(values.url)) {
			return;
		}

		setIsSubmitting(true);

		const titleChanged = values.title.trim() !== bookmark.Title;
		const urlChanged = values.url !== bookmark.URL;
		const tagsChanged =
			values.tags !== (bookmark.Tags ? bookmark.Tags.join(", ") : "");

		if (!titleChanged && !urlChanged && !tagsChanged) {
			await showToast({ style: Toast.Style.Success, title: "No changes" });
			popToRoot();
			return;
		}

		const { exitCode, stderr } = await bmEdit(
			bookmark.ID,
			urlChanged ? values.url : undefined,
			tagsChanged ? values.tags : undefined,
			titleChanged ? values.title.trim() : undefined,
		);
		setIsSubmitting(false);

		if (exitCode === 0) {
			await showToast({
				style: Toast.Style.Success,
				title: "Bookmark updated",
			});
			popToRoot();
		} else if (exitCode === 1) {
			if (stderr.includes("Duplicate URL")) {
				setUrlError("Duplicate URL");
			} else if (stderr.toLowerCase().includes("not found")) {
				setUrlError("Bookmark not found");
			} else {
				setUrlError(stderr.trim() || "Failed to update bookmark");
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
				id="title"
				title="Title"
				placeholder="Bookmark title"
				value={titleValue}
				onChange={setTitleValue}
				error={titleError}
				onBlur={(e) => validateTitle(e.target.value)}
			/>
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
