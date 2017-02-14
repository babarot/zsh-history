package history

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/b4b4r07/zsh-history/db"
	"github.com/fatih/structs"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Screen struct {
	width, height int
	cursorX       int
	selectedLine  int
	input         []rune
	candidates    []string
	mutex         sync.Mutex
	history       *History
	vimMode       bool
}

func getKeyTagMaps() map[string]string {
	r := &db.Record{}
	keyTags := make(map[string]string)

	rec := structs.New(r)
	for _, name := range rec.Names() {
		f := rec.Field(name)
		tagValue := f.Tag("db")
		keyTags[tagValue] = name
	}

	return keyTags
}

func convertToKeys(tags []string) []string {
	var keys []string
	keyTagMap := getKeyTagMaps()
	for _, tag := range tags {
		if key, ok := keyTagMap[tag]; ok {
			keys = append(keys, key)
		}
	}
	return keys
}

func NewScreen(buffer string) *Screen {
	var cfg config
	err := cfg.load()
	if err != nil {
		panic(err)
	}

	query := cfg.InitQuery
	if buffer != "" {
		// check SQL syntax
		if strings.Contains(strings.ToUpper(query), "WHERE") && strings.Contains(strings.ToUpper(query), "LIKE") {
			query = strings.Replace(query, "%%", "%"+buffer+"%", -1)
		}
	}
	input := []rune(query)
	x := strings.LastIndex(string(input), cfg.InitCursor)
	if x < 0 {
		x = len(string(input))
	}

	s := &Screen{
		cursorX:      x,
		selectedLine: 0,
		input:        input,
		history:      NewHistory(),
	}

	rows, _ := s.history.Query(string(input))
	for _, row := range rows {
		// TODO: imple
		// msg := ""
		// rowMap := structs.Map(row)
		// for _, key := range convertToKeys(cfg.ScreenColumns) {
		// 	msg += "|" + rowMap[key].(string)
		// }
		// s.candidates = append(s.candidates, msg)
		s.candidates = append(s.candidates, row.Command)
	}

	s.width, s.height = termbox.Size()
	return s
}

func (s *Screen) ToggleVimMode() {
	s.vimMode = !s.vimMode
}

func (s *Screen) IsVimMode() bool {
	return s.vimMode
}

func (s *Screen) MoveCusorBegin() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cursorX = 0
}

func (s *Screen) MoveCusorEnd() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cursorX = len(s.input)
}

func (s *Screen) MoveCusorForward() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursorX < len(s.input) {
		s.cursorX++
	}
}

func (s *Screen) MoveCusorBackward() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursorX > 0 {
		s.cursorX--
	}
}

func (s *Screen) SelectNext() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.selectedLine == len(s.candidates)-1 {
		s.selectedLine = 0
	} else {
		s.selectedLine++
	}
}

func (s *Screen) SelectPrevious() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.selectedLine == 0 {
		s.selectedLine = len(s.candidates) - 1
	} else {
		s.selectedLine--
	}
}

func (s *Screen) DeleteBackwardChar() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursorX > 0 {
		s.input = append(s.input[0:s.cursorX-1], s.input[s.cursorX:len(s.input)]...)
		s.cursorX--
		s.selectedLine = 0
	}
}

func (s *Screen) DeleteChar() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursorX < len(s.input) {
		s.input = append(s.input[0:s.cursorX], s.input[s.cursorX+1:len(s.input)]...)
		s.selectedLine = 0
	}
}

func (s *Screen) DeleteBackwardWord() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cursorX > 0 {
		var words int
		pos := regexp.MustCompile(`\s+`).FindStringIndex(string(s.input))
		if len(pos) > 0 && pos[len(pos)-1] > 0 {
			words = pos[len(pos)-1]
			s.input = append(s.input[0:s.cursorX-words], s.input[s.cursorX:len(s.input)]...)
		}
		s.cursorX -= words
		s.selectedLine = 0
	}
}

func (s *Screen) ClearPrompt() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.input = []rune{}
	s.cursorX = 0
}

func (s *Screen) InsertChar(ch rune) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tmp := []rune{}
	tmp = append(tmp, s.input[0:s.cursorX]...)
	tmp = append(tmp, ch)
	s.input = append(tmp, s.input[s.cursorX:len(s.input)]...)
	s.cursorX++
	s.selectedLine = 0
}

func (s *Screen) SetSize() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.width, s.height = termbox.Size()
}

func (s *Screen) DrawPrompt() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var vimMode bool
	if s.IsVimMode() {
		vimMode = true
	}
	setPrompt(s.width, s.cursorX, s.selectedLine, len(s.candidates), s.input, vimMode)
	termbox.Flush()
}

func (s *Screen) DrawScreen() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	var vimMode bool
	if s.IsVimMode() {
		vimMode = true
	}
	setPrompt(s.width, s.cursorX, s.selectedLine, len(s.candidates), s.input, vimMode)

	offsetPage := s.selectedLine / (s.height - 1)
	offsetLine := offsetPage * (s.height - 1)
	for i := 0; i < s.height-1 && i < len(s.candidates)-offsetLine; i++ {
		str := s.candidates[i+offsetLine]
		setLine(0, i+1, termbox.ColorDefault, termbox.ColorDefault, str)
	}

	selectLine(s.width, s.height, s.selectedLine, s.candidates)
	termbox.Flush()
}

func (s *Screen) Filter(done chan<- bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	go func() {
		s.candidates = []string{}
		s.mutex.Lock()
		rows, err := s.history.Query(string(s.input))
		if err != nil {
			done <- false
		}
		for _, row := range rows {
			s.candidates = append(s.candidates, row.Command)
		}
		s.selectedLine = 0
		s.mutex.Unlock()
		done <- true
	}()
}

func (s *Screen) GetOutput() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.candidates) != 0 {
		return s.candidates[s.selectedLine]
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

func setPrompt(width, cursorX, selectedLine, selectedLineLength int, input []rune, vimMode bool) {
	var cfg config
	err := cfg.load()
	if err != nil {
		panic(err)
	}
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	setLine(0, 0, termbox.ColorDefault, termbox.ColorDefault, cfg.Prompt, string(input))
	indicator := fmt.Sprintf("[%v/%v]", selectedLine+1, selectedLineLength)
	if vimMode {
		vimIndicator := "VIM-MODE"
		setLine(width-len(indicator)-len(vimIndicator)-1, 0, termbox.ColorGreen, termbox.ColorDefault, vimIndicator)
	}
	setLine(width-len(indicator), 0, termbox.ColorDefault, termbox.ColorDefault, indicator)
	termbox.SetCursor(runewidth.StringWidth(cfg.Prompt+string(input[0:cursorX])), 0)
}

func selectLine(width, height, selectedLine int, candidates []string) {
	if len(candidates) != 0 {
		x := 0
		y := selectedLine%(height-1) + 1
		str := candidates[selectedLine]
		setLine(0, y, termbox.ColorWhite, termbox.ColorBlue, str)
		x = runewidth.StringWidth(str)
		for x < width {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlue)
			x++
		}
	}
}
