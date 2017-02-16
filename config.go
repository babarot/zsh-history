package history

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type config struct {
	Prompt        string   `toml:"prompt"`
	InitQuery     string   `toml:"init_query"`
	InitCursor    string   `toml:"init_cursor"`
	ScreenColumns []string `toml:"screen_columns"`
	VimModePrompt string   `toml:"vim_mode_prompt"`
	IgnoreWords   []string `toml:"ignore_words"`
}

const tomlDir = "zhist"

func (cfg *config) load() error {
	var dir string
	if runtime.GOOS == "windows" {
		base := os.Getenv("APPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
		dir = filepath.Join(base, tomlDir)
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", tomlDir)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	tomlFile := filepath.Join(dir, "config.toml")

	_, err := os.Stat(tomlFile)
	if err == nil {
		_, err := toml.DecodeFile(tomlFile, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(tomlFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Set default value
	cfg.InitQuery = DefaultQuery
	cfg.InitCursor = Wildcard
	cfg.Prompt = Prompt
	cfg.ScreenColumns = []string{"command"}
	cfg.VimModePrompt = "VIM-MODE"
	cfg.IgnoreWords = []string{}

	return toml.NewEncoder(f).Encode(cfg)
}
