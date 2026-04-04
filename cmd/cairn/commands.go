package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ndy40/cairn/internal/store"
)

// hasHelpFlag reports whether the first argument is -h or --help.
func hasHelpFlag(args []string) bool {
	return len(args) > 0 && (args[0] == "-h" || args[0] == "--help")
}

func cmdAdd(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("add")
		os.Exit(0)
	}
	if len(ctx.args) < 1 {
		fatalf(3, "usage: cairn add <url>")
	}
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	tagsFlag := fs.String("tags", "", "comma-separated tags")
	if err := fs.Parse(ctx.args[1:]); err != nil {
		fatalf(3, "cairn add: %v", err)
	}
	runAdd(ctx.db, ctx.args[0], store.NormaliseTagsFromString(*tagsFlag))
}

func cmdEdit(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("edit")
		os.Exit(0)
	}

	if len(ctx.args) < 1 {
		fatalf(3, "usage: cairn edit <id> [--url=<url>] [--title=<title>] [--tags=<tags>]")
	}

	fs := flag.NewFlagSet("edit", flag.ContinueOnError)
	urlFlag := fs.String("url", "", "bookmark URL")
	title := fs.String("title", "", "bookmark title")
	tagsFlag := fs.String("tags", "", "comma-separated tags")

	if err := fs.Parse(ctx.args[1:]); err != nil {
		fatalf(3, "cairn edit: %v", err)
	}

	var urlSet, titleSet, tagsSet bool
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "url":
			urlSet = true
		case "title":
			titleSet = true
		case "tags":
			tagsSet = true
		}
	})

	runEdit(ctx.db, ctx.args[0], *urlFlag, urlSet, *title, titleSet, store.NormaliseTagsFromString(*tagsFlag), tagsSet)
}

func cmdList(ctx cmdContext) {
	runList(ctx.db, ctx.args)
}

func cmdSearch(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("search")
		os.Exit(0)
	}
	if len(ctx.args) < 1 {
		fatalf(3, "usage: cairn search <query>")
	}
	runSearch(ctx.db, ctx.args[0], ctx.args[1:])
}

func cmdDelete(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("delete")
		os.Exit(0)
	}
	if len(ctx.args) < 1 {
		fatalf(3, "usage: cairn delete <id>")
	}
	runDelete(ctx.db, ctx.args[0])
}

func cmdPin(ctx cmdContext) {
	if len(ctx.args) < 1 {
		fatalf(3, "usage: cairn pin <id>")
	}
	runPin(ctx.db, ctx.args[0])
}

func cmdSyncCmd(ctx cmdContext) {
	if len(ctx.args) == 0 || hasHelpFlag(ctx.args) {
		printCommandHelp("sync")
		os.Exit(0)
	}
	runSync(ctx.db, ctx.cfgManager, ctx.args)
}

func cmdUpdate(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("update")
		os.Exit(0)
	}
	runUpdate(ctx.args)
}

func cmdVersion(ctx cmdContext) {
	if hasHelpFlag(ctx.args) {
		printCommandHelp("version")
		os.Exit(0)
	}
	fmt.Printf("cairn version %s\n", version)
}

func cmdConfig(ctx cmdContext) {
	fmt.Printf("CAIRN_DB_PATH=%s\n", ctx.appCfg.DBPath)
	if ctx.appCfg.DropboxAppKey != "" {
		fmt.Println("CAIRN_DROPBOX_APP_KEY=(set)")
	}
}

func cmdHelpCmd(_ cmdContext) {
	printHelp()
}
