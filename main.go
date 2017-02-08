package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/b4b4r07/zsh-history/db"
)

var (
	append = flag.Bool("a", false, "Append to the history")
	list   = flag.Bool("l", false, "Show all histories")
	query  = flag.String("q", "", "Query string")
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	dbHandler := db.NewDBHandler()
	history := &History{dbHandler}

	if len(os.Args[1:]) == 0 {
		return msg(errors.New("too few arguments"))
	}

	if *append {
		if flag.NArg() < 2 {
			return msg(errors.New("too few arguments"))
		}
		cmd := flag.Arg(0)
		status, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			return msg(errors.New("string not"))
		}
		err = history.Append(cmd, status)
		if err != nil {
			return msg(err)
		}
	}

	if *list {
		if err := history.List(); err != nil {
			return msg(err)
		}
	}

	if *query != "" {
		rows, err := history.Query(*query)
		if err != nil {
			return msg(err)
		}
		for _, row := range rows {
			fmt.Printf("%s\n", row.Command)
		}
	}

	return 0
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}
