package screen

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/b4b4r07/zsh-history"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Matched struct {
	str string
	pos [][]int
}

func NewMatched(str string) *Matched {
	f := &Matched{str: str, pos: [][]int{}}
	return f
}

type Screen struct {
	width               int
	height              int
	cursor_x            int
	selected_line       int
	input               []rune
	candidates          []string
	filtered_candidates []Matched
	mutex               sync.Mutex
	h                   *history.History
}

func NewScreen() *Screen {
	input := []rune("select * from history")
	s := &Screen{cursor_x: 0, selected_line: 0, input: input, h: history.NewHistory()}
	// h := history.NewHistory()
	rows, _ := s.h.Query(string(input))
	for _, row := range rows {
		s.candidates = append(s.candidates, row.Command)
	}
	// s.candidates = gather_candidates_from_stdin(false)
	// s.candidates = []string{"hoge"}
	s.filtered_candidates = gather_filtered_candidates(s.input, s.candidates)
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
	if s.selected_line == len(s.filtered_candidates)-1 {
		s.selected_line = 0
	} else {
		s.selected_line++
	}
}

func (s *Screen) SelectPrevious() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.selected_line == 0 {
		s.selected_line = len(s.filtered_candidates) - 1
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
	set_prompt(s.width, s.cursor_x, s.selected_line, len(s.filtered_candidates), s.input)
	termbox.Flush()
}

func (s *Screen) DrawScreen() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw prompt.
	set_prompt(s.width, s.cursor_x, s.selected_line, len(s.filtered_candidates), s.input)

	// Draw filtered candidates.
	offset_page := s.selected_line / (s.height - 1)
	offset_line := offset_page * (s.height - 1)
	for i := 0; i < s.height-1 && i < len(s.filtered_candidates)-offset_line; i++ {
		str := s.filtered_candidates[i+offset_line].str
		set_line(0, i+1, termbox.ColorDefault, termbox.ColorDefault, str)
		for _, p := range s.filtered_candidates[i+offset_line].pos {
			begin := p[0]
			end := p[1]
			if end < len(p) {
				set_line(begin, i+1, termbox.ColorDefault|termbox.AttrBold|termbox.AttrUnderline, termbox.ColorDefault, string([]rune(str)[begin:end]))
			}
		}
	}

	// Draw selected lilne.
	select_line(s.width, s.height, s.selected_line, s.filtered_candidates)
	termbox.Flush()
}

func (s *Screen) Filter(done chan<- bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	input_snapshot := s.input
	go func() {
		// rows, _ := s.h.Query(string(s.input))
		// for _, row := range rows {
		// 	s.candidates = append(s.candidates, row.Command)
		// }
		f := gather_filtered_candidates(input_snapshot, s.candidates)
		s.mutex.Lock()
		if string(s.input) == string(input_snapshot) {
			s.filtered_candidates = f
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
	if len(s.filtered_candidates) != 0 {
		return s.filtered_candidates[s.selected_line].str
	}
	return ""
}

func gather_candidates_from_file(reverse bool, filenames []string) []string {
	var candidates = []string{}
	for _, filename := range filenames {
		read_file, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
		scanner := bufio.NewScanner(read_file)
		re, _ := regexp.Compile("^ *$")
		if reverse {
			// Gather in reverse order.
			for scanner.Scan() {
				str := scanner.Text()
				if scanner.Err() != nil {
					break
				}
				if !re.MatchString(str) {
					candidates = append([]string{str}, candidates...)
				}
			}
		} else {
			for scanner.Scan() {
				str := scanner.Text()
				if scanner.Err() != nil {
					break
				}
				if !re.MatchString(str) {
					candidates = append(candidates, str)
				}
			}
		}
	}
	return candidates
}

func gather_candidates_from_stdin(reverse bool) []string {
	var candidates = []string{}
	scanner := bufio.NewScanner(os.Stdin)
	re, _ := regexp.Compile("^ *$")
	if reverse {
		// Gather in reverse order.
		for scanner.Scan() {
			str := scanner.Text()
			if scanner.Err() != nil {
				break
			}
			if !re.MatchString(str) {
				candidates = append([]string{str}, candidates...)
			}
		}
	} else {
		for scanner.Scan() {
			str := scanner.Text()
			if scanner.Err() != nil {
				break
			}
			if !re.MatchString(str) {
				candidates = append(candidates, str)
			}
		}
	}
	return candidates
}

func set_line(x, y int, fg, bg termbox.Attribute, strs ...string) {
	for _, str := range strs {
		for _, c := range str {
			termbox.SetCell(x, y, c, fg, bg)
			x += runewidth.RuneWidth(c)
		}
	}
}

func set_prompt(width, cursor_x, selected_line, filtered_candidates_length int, input []rune) {
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	prompt := "> "
	set_line(0, 0, termbox.ColorDefault, termbox.ColorDefault, prompt, string(input))
	indicator := fmt.Sprintf("[%v/%v]", selected_line+1, filtered_candidates_length)
	set_line(width-len(indicator), 0, termbox.ColorDefault, termbox.ColorDefault, indicator)
	termbox.SetCursor(runewidth.StringWidth(prompt+string(input[0:cursor_x])), 0)
}

func select_line(width, height, selected_line int, filtered_candidates []Matched) {
	if len(filtered_candidates) != 0 {
		x := 0
		y := selected_line%(height-1) + 1
		str := filtered_candidates[selected_line].str
		set_line(0, y, termbox.ColorWhite, termbox.ColorBlue, str)
		x = runewidth.StringWidth(str)
		// Pad a string with whitespaces on the right to fill the line width.
		for x < width {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlue)
			x++
		}
		// Emphasize matches.
		for _, p := range filtered_candidates[selected_line].pos {
			begin := p[0]
			end := p[1]
			set_line(begin, y, termbox.ColorWhite|termbox.AttrBold|termbox.AttrUnderline, termbox.ColorBlue, string([]rune(str)[begin:end]))
		}
	}
}

func gather_filtered_candidates(input []rune, candidates []string) []Matched {
	input = []rune(strings.Trim(string(input), " "))
	filtered_candidates := []Matched{}

	if len(input) == 0 {
		for _, candidate := range candidates {
			f := *NewMatched(candidate)
			filtered_candidates = append(filtered_candidates, f)
		}
		return filtered_candidates
	}

	// Compile patterns.
	res := []*regexp.Regexp{}
	for _, pattern := range strings.Split(string(input), " ") {
		quoted_pattern := regexp.QuoteMeta(pattern)
		res = append(res, regexp.MustCompile(`(?i)`+quoted_pattern))
	}

	for _, candidate := range candidates {
		f := *NewMatched(candidate)
		for _, re := range res {
			matched_index := re.FindAllStringIndex(candidate, 10)
			if matched_index == nil {
				f.pos = [][]int{}
				break
			}
			for i, _ := range matched_index {
				f.pos = append(f.pos, matched_index[i])
			}
		}
		if len(f.pos) != 0 {
			filtered_candidates = append(filtered_candidates, f)
		}
	}
	return filtered_candidates
}
