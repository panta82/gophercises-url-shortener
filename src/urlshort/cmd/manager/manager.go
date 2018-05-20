package main

import (
	"flag"
	"fmt"
	"urlshort/database"
)

const (
	OP_LIST = iota
	OP_ADD = iota
	OP_REMOVE = iota
)

type CommandLineArgs struct {
	Op int
	Path string
	Url string
}

func parseCommandLineArgs() CommandLineArgs {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "URL Shortener manager utility\n\n" +
			"Usage:");
		flag.PrintDefaults()
	}

	flag.Bool("list", false, "List existing redirects (default command)")
	addOp := flag.Bool("add", false, "Add redirect")
	removeOp := flag.Bool("remove", false, "Remove redirect")
	path := flag.String("path", "", "Path to add or remove")
	url := flag.String("url", "", "Url to add")

	flag.Parse()

	result := CommandLineArgs {
		Path: *path,
		Url: *url,
	}

	if addOp != nil && *addOp {
		result.Op = OP_ADD
	} else if removeOp != nil && *removeOp {
		result.Op = OP_REMOVE
	} else {
		result.Op = OP_LIST
	}

	return result
}

func listOp(db database.Database) {
	redirects, err := db.ListAllRedirects()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%-20s | %s\n", "Path", "Url")
	fmt.Println("-------------------------------------------------------")
	for _, redirect := range redirects {
		fmt.Printf("%-20s | %s\n", redirect.Path, redirect.Url)
	}
	fmt.Println("-------------------------------------------------------")
	fmt.Println( len(redirects), "Redirects")
}

func addOp(db database.Database, path string, url string) {
	if path == "" {
		panic("You must provide path")
	}

	if url == "" {
		panic("You must provide url")
	}

	err := db.SetUrlForPath(path, url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Added redirect %s -> %s\n", path, url)
}

func removeOp(db database.Database, path string) {
	if path == "" {
		panic("You must provide path")
	}

	err := db.RemoveUrlForPath(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Removed redirect for %s\n", path)
}

func main()  {
	args := parseCommandLineArgs()

	db, err := database.NewDatabase(database.GetDefaultDatabasePath())
	if err != nil {
		panic(err)
	}

	if args.Op == OP_LIST {
		listOp(db)
	}
	if args.Op == OP_ADD {
		addOp(db, args.Path, args.Url)
	}
	if args.Op == OP_REMOVE {
		removeOp(db, args.Path)
	}
}
