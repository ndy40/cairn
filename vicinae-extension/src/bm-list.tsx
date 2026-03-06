import { List, Action, ActionPanel, open, Clipboard, showHUD } from "@vicinae/api";
import { useState, useEffect } from "react";
import { bmAvailable, bmList, Bookmark } from "./bm";

function formatDate(iso: string): string {
  return iso.slice(0, 10);
}

function BookmarkListItem({ bookmark }: { bookmark: Bookmark }) {
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
    if (!bmAvailable()) {
      setCliError("cairn is not installed. Install from: https://github.com/ndy40/bookmark-manager");
      setIsLoading(false);
      return;
    }
    const results = bmList();
    setBookmarks(results);
    setIsLoading(false);
  }, []);

  if (cliError) {
    return (
      <List>
        <List.EmptyView title="bm CLI not found" description={cliError} />
      </List>
    );
  }

  return (
    <List isLoading={isLoading} searchBarPlaceholder="Filter bookmarks...">
      {bookmarks.length === 0 && !isLoading ? (
        <List.EmptyView title="No bookmarks saved yet" />
      ) : (
        bookmarks.map((b) => <BookmarkListItem key={b.id} bookmark={b} />)
      )}
    </List>
  );
}
