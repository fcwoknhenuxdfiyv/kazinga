package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

var configFile string

func Read(fn string) *Config {
	configFile = path.Join(os.Getenv("HOME"), ".config/kazinga", fn)
	fi, err := os.Open(configFile)
	if err != nil && os.IsNotExist(err) {
		config := Config{
			EdgeBorderX: 4.0,
			EdgeBorderY: 4.0,
			BorderX:     4.0,
			BorderY:     4.0,
			Tweaks: Tweaks{
				"default": {
					{
						Class:   "__ALL__",
						NudgeX:  0,
						NudgeY:  0,
						ShrinkW: 0,
						ShrinkH: 0,
					},
					{
						Class:   "firefox",
						NudgeX:  0,
						NudgeY:  0,
						ShrinkW: 0,
						ShrinkH: 0,
					},
				},
			},
		}
		config.Save()
		return &config
	} else if err != nil {
		panic(err)
	}
	defer fi.Close()
	jsonParser := json.NewDecoder(fi)
	var config Config
	if err := jsonParser.Decode(&config); err != nil {
		panic(err)
	}
	return &config
}

func (c *Config) GetTweak(screenDims, class string) (Tweak, bool) {
	class = strings.ToLower(class)
	for _, t := range c.Tweaks[screenDims] {
		t.Class = strings.ReplaceAll(t.Class, ",", " ")
		t.Class = strings.ReplaceAll(t.Class, ";", " ")
		t.Class = strings.ReplaceAll(t.Class, ":", " ")
		classes := strings.Fields(t.Class)
		for _, c := range classes {
			if strings.ToLower(c) == class {
				return t, true
			}
		}
	}
	return Tweak{}, false
}

func (c Config) Save() {
	if err := os.MkdirAll(path.Dir(configFile), 0700); err != nil {
		panic(err)
	}
	configFile, err := os.Create(configFile)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	if out, err := json.MarshalIndent(c, "", "    "); err == nil {
		_, err = fmt.Fprintln(configFile, string(out))
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}
