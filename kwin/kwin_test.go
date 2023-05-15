package kwin

import (
	"kazinga/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig = config.Config{
	EdgeBorderX: 4.0,
	EdgeBorderY: 4.0,
	BorderX:     4.0,
	BorderY:     4.0,
	Tweaks: config.Tweaks{
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

var testWindows = []window{
	{
		Id:         "528ec59d-90ca-4f28-92aa-104d98d2762a",
		Class:      "org.kde.plasmashell",
		Title:      "desktop @ qrect(0,0 1920x1080) — plasma",
		Active:     false,
		Minimised:  false,
		Display:    -1,
		Dimensions: dimensions{X: 0, Y: 0, W: 1920, H: 1080},
	},
	{
		Id:         "3bf8faa3-9a1c-468a-8975-2612e2b94643",
		Class:      "org.kde.plasmashell",
		Title:      "plasma",
		Active:     false,
		Minimised:  false,
		Display:    -1,
		Dimensions: dimensions{X: 0, Y: 0, W: 1920, H: 30},
	},
	{
		Id:         "dc2d6f2d-5db3-42f7-a071-00d8de90d449",
		Class:      "org.kde.plasmashell",
		Title:      "plasma",
		Active:     false,
		Minimised:  false,
		Display:    -1,
		Dimensions: dimensions{X: 590, Y: 1000, W: 740, H: 80},
	},
	{
		Id:         "a515d902-4a82-421e-88dc-677c66d2941a",
		Class:      "org.kde.neochat",
		Title:      "neochat",
		Active:     false,
		Minimised:  true,
		Display:    1,
		Dimensions: dimensions{X: 4, Y: 34, W: 1912, H: 1042},
	},
	{
		Id:         "c0c91b2a-5348-4d85-bb06-b1e3da174ed2",
		Class:      "org.kde.discover",
		Title:      "updates — discover",
		Active:     false,
		Minimised:  true,
		Display:    1,
		Dimensions: dimensions{X: 4, Y: 34, W: 1912, H: 1042},
	},
}

func TestListWindows(t *testing.T) {
	// Should return at least 1 window
	wins = testWindows
	wins := ListWindows()
	assert.NotNil(t, wins)
	if assert.NotNil(t, wins) {
		found := false
		for _, win := range wins {
			if strings.Contains(win, "plasmashell") {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func TestCalculatePixels(t *testing.T) {
	Init()
	testData := map[string]dimensions{
		"2x2:1,1,1,1": {4, 34, 954, 519},
		"2x2:1,1,2,1": {4, 34, 1912, 519},
		"2x2:2,1,1,1": {962, 34, 954, 519},
		"2x2:1,1,1,2": {4, 34, 954, 1042},
		"1x1:1,1,1,1": {4, 34, 1912, 1042},
		"2x2:2,1,1,2": {962, 34, 954, 1042},
		"2x2:1,2,1,1": {4, 557, 954, 519},
		"2x2:1,2,2,1": {4, 557, 1912, 519},
		"2x2:2,2,1,1": {962, 557, 954, 519},
	}
	conf = &testConfig
	win := testWindows[len(testWindows)-1]
	disp := testDisplays[1]
	for dims, expected := range testData {
		pixels, err := calculatePixels(disp, win, dims)
		assert.Nil(t, err)
		assert.Greater(t, pixels.Y, 0.0)
		assert.Equal(t, expected.X, pixels.X, "x position of "+dims)
		assert.Equal(t, expected.Y, pixels.Y, "y position of "+dims)
		assert.Equal(t, expected.W, pixels.W, "width of "+dims)
		assert.Equal(t, expected.H, pixels.H, "height of "+dims)
	}
	for _, badDims := range []string{
		"",
		"1:",
		"a:",
		"axb:",
		"1xb:",
		"1x1:a,b,c",
		"1x1:a,1,1,1",
		"1x1:1,b,1,1",
		"1x1:1,1,c,1",
		"1x1:1,1,1,d",
	} {
		_, err := calculatePixels(disp, win, badDims)
		assert.NotNil(t, err)
	}
}
func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, dbusConnection)
	assert.NotNil(t, disps)
	assert.NotNil(t, wins)
}

func TestClose(t *testing.T) {
	assert.NotNil(t, dbusConnection)
	Close()
	assert.False(t, dbusConnection.Connected())
}

func TestRunOrRaise(t *testing.T) {
	id, err := RunOrRaise("blah", "blah", "", false)
	assert.NotNil(t, err)
	assert.Equal(t, id, "")
	Init()
	testData := []struct {
		class  string
		title  string
		cmd    string
		always bool
	}{
		{"blah", "blah", "", false},
		{"blah", "blah", "echo hello", false},
		{"blah", "blah", "blahsplat!", false},
	}
	windowWaitTime = 0.5
	for _, data := range testData {
		id, err = RunOrRaise(data.class, data.title, data.cmd, data.always)
		assert.NotNil(t, err, data)
	}
	id, err = RunOrRaise("org.kde.discover", "", "plasma-discover", false)
	assert.Nil(t, err, "run or raise plasma discover")
	Close()
}

func TestResizeWindow(t *testing.T) {
	Init()
	id, err := RunOrRaise("org.kde.discover", "", "plasma-discover", false)
	assert.Nil(t, err, "run or raise plasma discover")
	err = ResizeWindow(id, "1x1:1,1,1,1")
	assert.Nil(t, err, "resize plasma discover")
	err = MinimiseWindow(id)
	assert.Nil(t, err, "minimise plasma discover")
	Close()
}
