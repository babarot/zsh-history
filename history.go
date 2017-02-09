package history

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/b4b4r07/zsh-history/db"
	"github.com/nsf/termbox-go"
)

const (
	DefaultY int    = 1
	Prompt   string = "> "
)

var (
	f *[]rune
)

type History struct {
	DB   *db.DBHandler
	rows db.Records
}

func NewHistory() *History {
	return &History{
		DB:   db.NewDBHandler(),
		rows: db.Records{},
	}
}

func (h *History) Run() int {
	if !h.render() {
		return 1
	}
	fmt.Println(h.rows)
	return 0
}

func (h *History) render() bool {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	f = &[]rune{}

	contents := []string{}
	for {
		h.filterByQuery(string(*f))
		for _, row := range h.rows {
			contents = append(contents, fmt.Sprintf("%#v", row))
		}
		if len(h.rows) == 0 {
			contents = []string{}
		}
		draw(contents)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return false
			case termbox.KeySpace:
				*f = append(*f, rune(' '))
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
			case termbox.KeyHome, termbox.KeyCtrlA:
			case termbox.KeyEnd, termbox.KeyCtrlE:
			case termbox.KeyCtrlW:
				//delete whole word to period
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if i := len(*f) - 1; i >= 0 {
					slice := *f
					*f = slice[0:i]
				}
			case termbox.KeyEnter:
				return true
			case 0:
				*f = append(*f, rune(ev.Ch))
			default:
			}
		case termbox.EventError:
			panic(ev.Err)
			break
		default:
		}
	}
}

func (h *History) filterByQuery(q string) {
	rows, err := h.Query(q)
	if err == nil {
		h.rows = rows
	} else {
		h.rows = db.Records{}
	}
}

func draw(rows []string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	fs := Prompt + string(*f)
	drawln(0, 0, fs)
	termbox.SetCursor(len(fs), 0)
	for idx, row := range rows {
		drawln(0, idx+DefaultY, fmt.Sprintf("%#v\n", row))
	}
	termbox.Flush()
}

func drawln(x int, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault

	var c termbox.Attribute
	for i, s := range str {
		c = color
		termbox.SetCell(x+i, y, s, c, backgroundColor)
	}
}

func (h *History) Append(cmd string, status int) error {
	return h.DB.Insert(cmd, status)
}

func (h *History) List() error {
	list, err := h.DB.QueryList()
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
	return h.DB.Query(query)
}
