package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/ndy40/cairn/internal/store"
)

func printHelp()                  { fmt.Print(helpTexts[""]) }
func printCommandHelp(cmd string) { fmt.Print(helpTexts[cmd]) }

func printBookmarks(bookmarks []*store.Bookmark, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(bookmarks)
		return
	}
	printBookmarkTable(bookmarks)
}

func printBookmarkTable(bookmarks []*store.Bookmark) {
	const (
		colID     = 6
		colDomain = 22
		colTitle  = 38
		colURL    = 55
	)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Padding(0, 1)

	evenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	oddStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Padding(0, 1)

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	rows := make([][]string, 0, len(bookmarks))
	for _, b := range bookmarks {
		title := b.Title
		if len([]rune(title)) > colTitle-2 {
			title = string([]rune(title)[:colTitle-3]) + "…"
		}
		url := b.URL
		if len([]rune(url)) > colURL-2 {
			url = string([]rune(url)[:colURL-3]) + "…"
		}
		domain := b.Domain
		if len([]rune(domain)) > colDomain-2 {
			domain = string([]rune(domain)[:colDomain-3]) + "…"
		}
		rows = append(rows, []string{
			strconv.FormatInt(b.ID, 10),
			title,
			url,
			domain,
		})
	}

	tbl := ltable.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle).
		Headers("ID", "Title", "URL", "Domain").
		Width(colID + colDomain + colTitle + colURL + 10).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return headerStyle
			}
			if row%2 == 0 {
				return evenStyle
			}
			return oddStyle
		})

	for _, r := range rows {
		tbl.Row(r...)
	}

	fmt.Println(tbl.Render())
}
