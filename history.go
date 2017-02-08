package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/b4b4r07/zsh-history/db"
)

type History struct {
	dbHandler *db.DBHandler
}

func (h *History) Append(cmd string, status int) error {
	return h.dbHandler.Insert(cmd, status)
}

func (h *History) List() error {
	list, err := h.dbHandler.QueryList()
	if err != nil {
		return err
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 0, '\t', 0)
	for _, l := range list {
		fmt.Fprintf(w, "%s\t%s\t\"%s\"\t%d\n",
			l.DateTime, l.Directory, l.Command, l.Status,
		)
	}
	w.Flush()
	return nil
}

func (h *History) Query(query string) (db.Records, error) {
	return h.dbHandler.Query(query)
}
