package history

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/b4b4r07/zsh-history/db"
	"github.com/nsf/termbox-go"
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

func (h *History) Insert(cmd string, status int) error {
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

func (h *History) Screen(args []string) int {
	output := ""
	err := termbox.Init()
	if err != nil {
		log.Println(err)
		return 1
	}

	s := NewScreen(strings.Join(args, " "))
	s.DrawScreen()

	defer func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Close()
		if output != "" {
			fmt.Println(output)
		}
	}()

loop:
	for {
		var (
			updatePrompt     bool
			updateAll        bool
			updateWithFilter bool
		)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				s.ToggleVimMode()
			case termbox.KeyCtrlC, termbox.KeyCtrlG:
				break loop
			case termbox.KeyCtrlA:
				s.MoveCusorBegin()
				updatePrompt = true
			case termbox.KeyCtrlE:
				s.MoveCusorEnd()
				updatePrompt = true
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				s.MoveCusorForward()
				updatePrompt = true
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				s.MoveCusorBackward()
				updatePrompt = true
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				s.SelectNext()
				updateAll = true
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				s.SelectPrevious()
				updateAll = true
			case termbox.KeyEnter:
				output = s.GetOutput()
				break loop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				s.DeleteBackwardChar()
				updateWithFilter = true
			case termbox.KeyDelete, termbox.KeyCtrlD:
				s.DeleteChar()
				updateWithFilter = true
			case termbox.KeyCtrlU:
				s.ClearPrompt()
				updateWithFilter = true
			case termbox.KeyCtrlW:
				s.DeleteBackwardWord()
				updateWithFilter = true
			default:
				if ev.Key == termbox.KeySpace {
					ev.Ch = ' '
				}
				if ev.Ch > 0 {
					if s.IsVimMode() {
						switch ev.Ch {
						case 'j':
							s.SelectNext()
							updateAll = true
						case 'k':
							s.SelectPrevious()
							updateAll = true
						case 'l':
							s.MoveCusorForward()
							updatePrompt = true
						case 'h':
							s.MoveCusorBackward()
							updatePrompt = true
						case '0', '^':
							s.MoveCusorBegin()
							updatePrompt = true
						case '$':
							s.MoveCusorEnd()
							updatePrompt = true
						case 'i':
							s.ToggleVimMode()
						case 'a':
							s.ToggleVimMode()
							s.MoveCusorForward()
							updatePrompt = true
						case 'I':
							s.ToggleVimMode()
							s.MoveCusorBegin()
							updatePrompt = true
						case 'A':
							s.ToggleVimMode()
							s.MoveCusorEnd()
							updatePrompt = true
						}
					} else {
						s.InsertChar(ev.Ch)
						updateWithFilter = true
					}
				}
			}
		case termbox.EventResize:
			s.SetSize()
			updateAll = true
		}
		if updatePrompt {
			s.DrawPrompt()
		}
		if updateAll {
			s.DrawScreen()
		}
		if updateWithFilter {
			s.DrawPrompt()
			go func() {
				done := make(chan bool)
				s.Filter(done)
				if <-done {
					s.DrawScreen()
				}
			}()
		}
	}
	return 0
}
