package config

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dimensions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"width"`
	H float64 `json:"height"`
}

type display struct {
	Id         int64
	Dimensions dimensions `json:"displaySize"`
}

var testConfig = Config{
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

var testDisplays = map[int64]display{
	1: {
		Id:         0,
		Dimensions: dimensions{X: 0, Y: 0, W: 1920, H: 1080},
	},
}

func TestRead(t *testing.T) {
	fn := "kazinga_test.conf"
	os.Remove(path.Join(os.Getenv("HOME"), ".config/kazinga", fn))
	conf := Read(fn)
	assert.NotNil(t, conf)
	conf = Read(fn)
	assert.NotNil(t, conf)
}

func TestGetTweak(t *testing.T) {
	disp := testDisplays[1]
	key := fmt.Sprintf("%.0fx%.0f", disp.Dimensions.W, disp.Dimensions.H)
	_, exists := testConfig.GetTweak(key, "__ALL__")
	assert.Equal(t, false, exists, key+" __ALL__")

}
