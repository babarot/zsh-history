package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/b4b4r07/zsh-history"
	"github.com/b4b4r07/zsh-history/fourmi"
)

var (
	append      = flag.Bool("a", false, "Append to the history")
	list        = flag.Bool("l", false, "Show all histories")
	query       = flag.String("q", "", "Query string")
	interactive = flag.Bool("i", false, "")
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	h := history.NewHistory()

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
		err = h.Append(cmd, status)
		if err != nil {
			return msg(err)
		}
	}

	if *list {
		if err := h.List(); err != nil {
			return msg(err)
		}
	}

	if *query != "" {
		rows, err := h.Query(*query)
		if err != nil {
			return msg(err)
		}
		for _, row := range rows {
			fmt.Printf("%s\n", row.Command)
		}
	}

	if *interactive {
		// return h.Run()
		f := fourmi.New()
		f.Run()
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
