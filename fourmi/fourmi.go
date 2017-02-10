package fourmi

import (
	"fmt"

	// "github.com/b4b4r07/zsh-history"
	"github.com/b4b4r07/zsh-history/screen"
	"github.com/nsf/termbox-go"
)

type Fourmi struct {
}

func New() *Fourmi {
	return &Fourmi{}
}

func (f *Fourmi) Run() {
	output := ""
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	// candidates := []string{}
	// h := history.NewHistory()
	// rows, _ := h.Query("select * from history")
	// for _, row := range rows {
	// 	candidates = append(candidates, row.Command)
	// }
	s := screen.NewScreen()
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
		update_prompt := false
		update_all := false
		update_with_filtering := false
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlG: // Terminate.
				break loop
			case termbox.KeyCtrlA: // Move the cursor to the start.
				s.MoveCusorBegin()
				update_prompt = true
			case termbox.KeyCtrlE: // Move the cursor to the end.
				s.MoveCusorEnd()
				update_prompt = true
			case termbox.KeyArrowRight, termbox.KeyCtrlF: // Move the cursor to the next character.
				s.MoveCusorForward()
				update_prompt = true
			case termbox.KeyArrowLeft, termbox.KeyCtrlB: // Move the cursor to the previous character.
				s.MoveCusorBackward()
				update_prompt = true
			case termbox.KeyArrowDown, termbox.KeyCtrlN: // Move the cursor to the next line.
				s.SelectNext()
				update_all = true
			case termbox.KeyArrowUp, termbox.KeyCtrlP: // Move the cursor to the previous line.
				s.SelectPrevious()
				update_all = true
			case termbox.KeyEnter: // Terminate with output.
				output = s.Get_output()
				break loop
			case termbox.KeyBackspace, termbox.KeyBackspace2: // As backspace-key.
				s.DeleteBackwardChar()
				update_with_filtering = true
			case termbox.KeyDelete, termbox.KeyCtrlD: // As delete-key.
				s.DeleteChar()
				update_with_filtering = true
			case termbox.KeyCtrlU: // Make the input empty.
				s.ClearPrompt()
				update_with_filtering = true
			default:
				if ev.Key == termbox.KeySpace {
					ev.Ch = ' '
				}
				if ev.Ch > 0 {
					s.InsertChar(ev.Ch)
					update_with_filtering = true
				}
			}
		case termbox.EventResize:
			s.SetSize()
			update_all = true
		}
		if update_prompt {
			s.DrawPrompt()
		}
		if update_all {
			s.DrawScreen()
		}
		if update_with_filtering {
			// Update prompt immediately.
			s.DrawPrompt()
			// Gather filtered candidates.
			go func() {
				done := make(chan bool)
				s.Filter(done)
				// Draw result if it is not obsolete.
				if <-done {
					s.DrawScreen()
				}
			}()
		}
	}
}
