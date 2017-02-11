package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/b4b4r07/zsh-history"
)

var (
	insert = flag.Bool("i", false, "Insert to the history")
	query  = flag.String("q", "", "Query searching")
	screen = flag.Bool("s", false, "Start to screen searching")
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

	if *insert {
		if flag.NArg() < 2 {
			return msg(errors.New("Please give 'command', 'status' arguments"))
		}
		cmd := flag.Arg(0)
		status, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			return msg(errors.New("'status': not integer"))
		}
		err = h.Insert(cmd, status)
		if err != nil {
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

	if *screen {
		return h.Screen(flag.Args())
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
