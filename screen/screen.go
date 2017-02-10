package screen

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/b4b4r07/zsh-history"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	prompt string = "sqlite3> "
)

type Screen struct {
	width         int
	height        int
	cursor_x      int
	selected_line int
	input         []rune
	candidates    []string
	mutex         sync.Mutex
	h             *history.History
}

func NewScreen() *Screen {
	input := []rune(history.DefaultQuery)
	x := strings.Index(string(input), history.InputPint) + 1
	s := &Screen{
		cursor_x:      x, //len(string(input))
		selected_line: 0,
		input:         input,
		h:             history.NewHistory(),
	}
	rows, _ := s.h.Query(string(input))
	for _, row := range rows {
		s.candidates = append(s.candidates, row.Command)
	}
	s.width, s.height = termbox.Size()
	return s
}

func (s *Screen) MoveCusorBegin() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cursor_x = 0
}

func (s *Screen) MoveCusorEnd() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cursor_x = len(s.input)
}

func (s *Screen) MoveCusorForward() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursor_x < len(s.input) {
		s.cursor_x++
	}
}

func (s *Screen) MoveCusorBackward() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursor_x > 0 {
		s.cursor_x--
	}
}

func (s *Screen) SelectNext() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.selected_line == len(s.candidates)-1 {
		s.selected_line = 0
	} else {
		s.selected_line++
	}
}

func (s *Screen) SelectPrevious() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.selected_line == 0 {
		s.selected_line = len(s.candidates) - 1
	} else {
		s.selected_line--
	}
}

func (s *Screen) DeleteBackwardChar() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursor_x > 0 {
		s.input = append(s.input[0:s.cursor_x-1], s.input[s.cursor_x:len(s.input)]...)
		s.cursor_x--
		s.selected_line = 0
	}
}

func (s *Screen) DeleteChar() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursor_x < len(s.input) {
		s.input = append(s.input[0:s.cursor_x], s.input[s.cursor_x+1:len(s.input)]...)
		s.selected_line = 0
	}
}

func (s *Screen) DeleteBackwardWord() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursor_x > 0 {
		var words int
		pos := regexp.MustCompile(`\s+`).FindStringIndex(string(s.input))
		if len(pos) > 0 && pos[len(pos)-1] > 0 {
			words = pos[len(pos)-1]
			s.input = append(s.input[0:s.cursor_x-words], s.input[s.cursor_x:len(s.input)]...)
		}
		s.cursor_x -= words
		s.selected_line = 0
	}
}

func (s *Screen) ClearPrompt() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.input = []rune{}
	s.cursor_x = 0
}

func (s *Screen) InsertChar(ch rune) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tmp := []rune{}
	tmp = append(tmp, s.input[0:s.cursor_x]...)
	tmp = append(tmp, ch)
	s.input = append(tmp, s.input[s.cursor_x:len(s.input)]...)
	s.cursor_x++
	s.selected_line = 0
}

func (s *Screen) SetSize() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.width, s.height = termbox.Size()
}

func (s *Screen) DrawPrompt() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	setPrompt(s.width, s.cursor_x, s.selected_line, len(s.candidates), s.input)
	termbox.Flush()
}

func (s *Screen) DrawScreen() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	setPrompt(s.width, s.cursor_x, s.selected_line, len(s.candidates), s.input)

	offset_page := s.selected_line / (s.height - 1)
	offset_line := offset_page * (s.height - 1)
	for i := 0; i < s.height-1 && i < len(s.candidates)-offset_line; i++ {
		str := s.candidates[i+offset_line]
		setLine(0, i+1, termbox.ColorDefault, termbox.ColorDefault, str)
	}

	selectLine(s.width, s.height, s.selected_line, s.candidates)
	termbox.Flush()
}

func (s *Screen) Filter(done chan<- bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	input_snapshot := s.input
	go func() {
		rows, _ := s.h.Query(string(s.input))
		s.candidates = []string{}
		s.mutex.Lock()
		if string(s.input) == string(input_snapshot) {
			for _, row := range rows {
				s.candidates = append(s.candidates, row.Command)
			}
			s.selected_line = 0
			s.mutex.Unlock()
			done <- true
		}
		// Abort a drawing old results.
		done <- false
	}()
}

func (s *Screen) Get_output() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.candidates) != 0 {
		return s.candidates[s.selected_line]
	}
	return ""
}

func setLine(x, y int, fg, bg termbox.Attribute, strs ...string) {
	for _, str := range strs {
		for _, c := range str {
			termbox.SetCell(x, y, c, fg, bg)
			x += runewidth.RuneWidth(c)
		}
	}
}

func setPrompt(width, cursor_x, selected_line, filtered_candidates_length int, input []rune) {
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	setLine(0, 0, termbox.ColorDefault, termbox.ColorDefault, prompt, string(input))
	indicator := fmt.Sprintf("[%v/%v]", selected_line+1, filtered_candidates_length)
	setLine(width-len(indicator), 0, termbox.ColorDefault, termbox.ColorDefault, indicator)
	termbox.SetCursor(runewidth.StringWidth(prompt+string(input[0:cursor_x])), 0)
}

func selectLine(width, height, selected_line int, candidates []string) {
	if len(candidates) != 0 {
		x := 0
		y := selected_line%(height-1) + 1
		str := candidates[selected_line]
		setLine(0, y, termbox.ColorWhite, termbox.ColorBlue, str)
		x = runewidth.StringWidth(str)
		// Pad a string with whitespaces on the right to fill the line width.
		for x < width {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlue)
			x++
		}
	}
}
