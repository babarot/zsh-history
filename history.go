package history

import (
	"fmt"
	"os"
	"regexp"
	// "sync"
	"text/tabwriter"

	"github.com/b4b4r07/zsh-history/db"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	DefaultY     int    = 1
	Prompt       string = "sqlite3> "
	DefaultQuery string = "SELECT DISTINCT(command) FROM history WHERE command LIKE '%%' AND status = 0 ORDER BY id DESC"
	InputPint    string = "%%"
)

var input = []rune{}
var width, height int
var cursor_x, cursor_y int

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

	input = []rune("select * from history")
	cursor_x = len(string(input))

	contents := []string{}
	for {
		update := false

		h.filterByQuery(string(input))
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
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				if cursor_y > 0 {
					cursor_y--
				}
				// update = true
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				if len(h.rows) > 0 {
					cursor_y++
				}
				// update = true
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				if cursor_x > 0 {
					cursor_x--
				}
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				if cursor_x < len([]rune(input)) {
					cursor_x++
				}
			case termbox.KeyHome, termbox.KeyCtrlA:
				cursor_x = 0
			case termbox.KeyEnd, termbox.KeyCtrlE:
				cursor_x = len(string(input))
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				// if i := len(input) - 1; i >= 0 {
				// 	cursor_x--
				// 	slice := input
				// 	input = slice[0:i]
				// }
				if cursor_x > 0 {
					input = append(input[0:cursor_x-1], input[cursor_x:len(input)]...)
					cursor_x--
					update = true
				}
			case termbox.KeyEnter:
				return true
			// case termbox.KeySpace:
			// 	input = append(input, rune(' '))
			// case 0:
			// 	input = append(input, rune(ev.Ch))
			// 	cursor_x++
			// default:
			case termbox.KeyCtrlW:
				part := string(input[0:cursor_x])
				rest := input[cursor_x:len(input)]
				pos := regexp.MustCompile(`\s+`).FindStringIndex(part)
				if len(pos) > 0 && pos[len(pos)-1] > 0 {
					input = []rune(part[0 : pos[len(pos)-1]-1])
					input = append(input, rest...)
				} else {
					input = []rune{}
				}
				cursor_x = len(input)
				// update = true
			default:
				cursor_y = 0
				if ev.Key == termbox.KeySpace {
					ev.Ch = ' '
				}
				if ev.Ch > 0 {
					out := []rune{}
					out = append(out, input[0:cursor_x]...)
					out = append(out, ev.Ch)
					input = append(out, input[cursor_x:len(input)]...)
					cursor_x++
					// update = true
				}
			}
		case termbox.EventError:
			panic(ev.Err)
			break
		default:
		}
		if update {
			draw(contents)
		}
	}
}

func (h *History) filterByQuery(q string) {
	// mutex.Lock()
	// defer mutex.Unlock()
	//
	// defer func() {
	// 	recover()
	// }()

	if len(h.rows) == 0 {
		rows, err := h.Query(q)
		if err == nil {
			h.rows = rows
		} else {
			h.rows = db.Records{}
		}
	}
}

func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range []rune(msg) {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func draw(rows []string) {
	width, height = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	fs := Prompt + string(input)
	if cursor_y < 0 {
		cursor_y = 0
	}
	if cursor_y >= height {
		cursor_y = height - 1
	}
	drawln(0, 0, fs)
	pos := len(Prompt)
	termbox.SetCursor(pos+runewidth.StringWidth(string(input[0:cursor_x])), 0)
	for idx, row := range rows {
		if idx == cursor_y {
			// drawln(0, idx+DefaultY, "hoge")
			print_tb(0, idx+DefaultY, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack, row)
		} else {
			print_tb(0, idx+DefaultY, termbox.ColorDefault, termbox.ColorDefault, row)
			// drawln(0, idx+DefaultY, fmt.Sprintf("%#v\n", row))
		}
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
