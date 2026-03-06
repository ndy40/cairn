import {
  Form,
  Action,
  ActionPanel,
  Clipboard,
  showToast,
  Toast,
  popToRoot,
} from "@vicinae/api";
import { useState, useEffect } from "react";
import { bmAvailable, bmAdd } from "./bm";

export default function AddBookmark() {
  const [urlValue, setUrlValue] = useState("");
  const [urlError, setUrlError] = useState<string | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    Clipboard.readText().then((text) => {
      if (text && (text.startsWith("http://") || text.startsWith("https://"))) {
        setUrlValue(text);
      }
    });
  }, []);

  function validateUrl(value: string): boolean {
    if (!value || value.trim() === "") {
      setUrlError("URL is required");
      return false;
    }
    if (!value.startsWith("http://") && !value.startsWith("https://")) {
      setUrlError("URL must start with http:// or https://");
      return false;
    }
    setUrlError(undefined);
    return true;
  }

  async function handleSubmit(values: { url: string; tags: string }) {
    if (!validateUrl(values.url)) {
      return;
    }

    if (!bmAvailable()) {
      await showToast({
        style: Toast.Style.Failure,
        title: "bm CLI not found",
        message: "Install cairn from: https://github.com/ndy40/bookmark-manager",
      });
      return;
    }

    setIsSubmitting(true);
    const { exitCode, stderr } = bmAdd(values.url, values.tags);
    setIsSubmitting(false);

    if (exitCode === 0) {
      await showToast({ style: Toast.Style.Success, title: "Saved" });
      popToRoot();
    } else if (exitCode === 2) {
      await showToast({
        style: Toast.Style.Success,
        title: "Saved (title unavailable)",
      });
      popToRoot();
    } else if (exitCode === 1) {
      setUrlError("Already bookmarked");
    } else {
      setUrlError(stderr || "Failed to save bookmark");
    }
  }

  return (
    <Form
      isLoading={isSubmitting}
      actions={
        <ActionPanel>
          <Action.SubmitForm title="Save Bookmark" onSubmit={handleSubmit} />
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
      />
    </Form>
  );
}
