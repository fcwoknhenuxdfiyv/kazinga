package kwin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"kazinga/cmd"
	"kazinga/config"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

var (
	displaysChan     chan displays
	windowsChan      chan windows
	disps            displays
	wins             windows
	conf             *config.Config
	badDimsError     = errors.New("bad dimensions")
	winNotFoundError = errors.New("window not found")
	windowWaitTime   = 10.0 // Wait for 10 seconds for window to be created
)

func Init() {
	startEventListener()
	displaysChan = make(chan displays)
	windowsChan = make(chan windows)
	getDisplays()
	disps = <-displaysChan
	getWindows()
	wins = <-windowsChan
	conf = config.Read("kazinga.conf")
}

func RunOrRaise(class, title, command string, always bool) (string, error) {
	win, err := find(wins, "", false, class, title)
	if (err != nil || always) && command != "" {
		if err := cmd.Run(command); err != nil {
			return "", err
		}
		counter := windowWaitTime * 5
		for err != nil && counter > 0 {
			counter--
			time.Sleep(200 * time.Millisecond)
			getWindows()
			wins = <-windowsChan
			win, err = find(wins, "", false, class, title)
		}
	}
	if err != nil {
		return "", err
	}
	script := `
        const id='%s';
        const clients = workspace.clientList();
        for (var i=0; i<clients.length; i++) {
            var client = clients[i];
            if (client.internalId == id) {
                workspace.activeClient = client;
            }
        }
        `
	_, err = runScript(fmt.Sprintf(script, win.Id))
	return win.Id, err
}

func ListWindows() []string {
	var lines []string
	for _, win := range wins {
		lines = append(lines, fmt.Sprintf("-class '%s' -title '%s'", win.Class, win.Title))
	}
	return lines
}

func MinimiseWindow(id string) error {
	active := true
	if id != "" {
		active = false
	}
	win, err := find(wins, id, active, "", "")
	if err != nil {
		return err
	}
	script := `
        const id='%s';
        const clients = workspace.clientList();
        for (var i=0; i<clients.length; i++) {
            var client = clients[i];
            if (client.internalId == id) {
                client.minimized = true;
                break;
            }
        }
        `
	_, err = runScript(fmt.Sprintf(script, win.Id))
	return err
}

func ResizeWindow(id, dims string) error {
	active := true
	if id != "" {
		active = false
	}
	win, err := find(wins, id, active, "", "")
	if err != nil {
		return err
	}
	pixels, err := calculatePixels(disps[win.Display], win, dims)
	if err != nil {
		return err
	}
	script := `
        const id='%s';
        const clients = workspace.clientList();
        for (var i=0; i<clients.length; i++) {
            var client = clients[i];
            if (client.internalId == id) {
                if (client.moveable) {
                    client.geometry = {
                        x: %.0f,
                        y: %.0f,
                        width: %.0f,
                        height: %.0f
                    };
                };
                break;
            }
        }
        `
	_, err = runScript(fmt.Sprintf(script, win.Id, pixels.X, pixels.Y, pixels.W, pixels.H))
	return err
}

// find finds a window based on it class and/or title
func find(windows []window, id string, findActive bool, class, title string) (window, error) {
	class = strings.ToLower(class)
	title = strings.ToLower(title)
	for _, win := range windows {
		winClass := strings.ToLower(win.Class)
		winTitle := strings.ToLower(win.Title)
		switch {
		case id != "":
			if win.Id == id {
				return win, nil
			}
		case findActive == true:
			if win.Active {
				return win, nil
			}
		case class != "" && title != "":
			if strings.Contains(winClass, class) && strings.Contains(winTitle, title) {
				return win, nil
			}
		case class != "":
			if strings.Contains(winClass, class) {
				return win, nil
			}
		case title != "":
			if strings.Contains(winTitle, title) {
				return win, nil
			}
		}
	}
	return window{}, winNotFoundError
}

func getDisplays() {
	script := Wrap("DisplayInfo", "JSON.stringify(workspace)")
	runScript(script)
}

func (f dbusResponse) DisplayInfo(data string) (string, *dbus.Error) {
	var d display
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return "", err.(*dbus.Error)
	}
	// tph, bph, err := getPanelHeights()
	// if err != nil {
	// 	return "", err.(*dbus.Error)
	// }
	// d.Dimensions.Y += tph
	// d.Dimensions.H -= (tph + bph)
	displays := displays{1: d}
	displaysChan <- displays
	return string(f), nil
}

func getWindows() {
	script := `
        var data = [];
        const clients = workspace.clientList();
        for (var i=0; i<clients.length; i++) {
            var client = clients[i];
            var cl = {
                id: client.internalId,
                class: client.resourceClass+"",
                title: client.caption+"",
                geometry: client.geometry,
                active: client.active,
                desktop: client.desktop,
                minimised: client.minimized,
            }
            data.push(cl);
        }
        ` +
		Wrap("WindowList", "JSON.stringify(data)")
	runScript(script)
}

func (f dbusResponse) WindowList(data string) (string, *dbus.Error) {
	var allWindows, realWindows windows
	err := json.Unmarshal([]byte(data), &allWindows)
	if err != nil {
		return "", err.(*dbus.Error)
	}
	for _, win := range allWindows {
		win.Class = strings.ToLower(win.Class)
		win.Title = strings.ToLower(win.Title)
		realWindows = append(realWindows, win)
	}
	windowsChan <- realWindows
	return string(f), nil
}

func getPanelHeights() (float64, float64, error) {
	script := "print(JSON.stringify(panels()))"
	obj := GetDbusObject("org.kde.plasmashell", "/PlasmaShell")
	call := obj.Call("org.kde.PlasmaShell.evaluateScript", 0, script)
	if call.Err != nil {
		return 0, 0, call.Err
	}
	type PlasmaPanel struct {
		Alignment string  `json:"alignment"`
		Hiding    string  `json:"hiding"`
		Location  string  `json:"location"`
		Height    float64 `json:"height"`
	}
	var topPanelHeight, bottomPanelHeight float64
	var panels []PlasmaPanel
	if len(call.Body) == 1 {
		var body []byte
		call.Store(&body)
		err := json.Unmarshal(body, &panels)
		if err != nil {
			return 0, 0, err
		}
		for _, panel := range panels {
			if panel.Hiding != "none" {
				continue
			}
			if panel.Location == "top" {
				topPanelHeight += panel.Height
			} else {
				bottomPanelHeight += panel.Height
			}
		}

	}
	// TODO: left and right panels
	return topPanelHeight, bottomPanelHeight, nil
}

func runScript(script string) ([]interface{}, error) {
	name, fn, err := writeScript(script)
	if err != nil {
		return nil, err
	}
	defer func() { os.Remove(fn) }()
	obj := GetDbusObject("org.kde.KWin", "/Scripting")
	call := obj.Call("org.kde.kwin.Scripting.loadScript", 0, fn, name)
	if call.Err != nil {
		return nil, call.Err
	}
	id := GetDbusObjectPathId(call)
	obj = GetDbusObject("org.kde.KWin", id)
	call = obj.Call("org.kde.kwin.Script.run", 0)
	if call.Err != nil {
		return nil, call.Err
	}
	obj = GetDbusObject("org.kde.KWin", "/Scripting")
	call = obj.Call("org.kde.kwin.Scripting.unloadScript", 0, name)
	if call.Err != nil {
		return nil, call.Err
	}
	return call.Body, nil
}

func writeScript(contents string) (string, string, error) {
	file, err := ioutil.TempFile("", "*.kazinga")
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = fmt.Fprintln(file, contents)
	if err != nil {
		return "", "", err
	}
	return path.Base(file.Name()), file.Name(), nil
}

func calculatePixels(disp display, win window, dims string) (dimensions, error) {
	tph, bph, err := getPanelHeights()
	if err != nil {
		return dimensions{}, err
	}
	disp.Dimensions.Y += tph
	disp.Dimensions.H -= (tph + bph)

	parts := strings.Split(dims, ":")
	if len(parts) < 2 {
		return dimensions{}, badDimsError
	}

	grid := strings.Split(parts[0], "x")
	if len(grid) < 2 {
		return dimensions{}, badDimsError
	}

	gridX, err := strconv.ParseFloat(grid[0], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	scaleX := disp.Dimensions.W / gridX

	gridY, err := strconv.ParseFloat(grid[1], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	scaleY := disp.Dimensions.H / gridY

	size := strings.Split(parts[1], ",")
	if len(size) < 4 {
		return dimensions{}, badDimsError
	}

	x, err := strconv.ParseFloat(size[0], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	x = (x-1)*scaleX + disp.Dimensions.X

	y, err := strconv.ParseFloat(size[1], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	y = (y-1)*scaleY + disp.Dimensions.Y

	w, err := strconv.ParseFloat(size[2], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	w = w * scaleX

	h, err := strconv.ParseFloat(size[3], 64)
	if err != nil {
		return dimensions{}, badDimsError
	}
	h = h * scaleY

	var brx float64
	brx = x + w
	if x <= disp.Dimensions.X+conf.EdgeBorderX {
		x = disp.Dimensions.X + conf.EdgeBorderX
	} else if x > disp.Dimensions.X+conf.EdgeBorderX {
		x += conf.BorderX / 2
	}
	if brx >= disp.Dimensions.W-conf.EdgeBorderX {
		brx = disp.Dimensions.W - conf.EdgeBorderX
		w = brx - x + disp.Dimensions.X
	} else if brx < disp.Dimensions.W-conf.EdgeBorderX {
		brx -= conf.BorderX / 2
		w = brx - x
	}

	var bry float64
	bry = y + h
	if y <= disp.Dimensions.Y+conf.EdgeBorderY {
		y = disp.Dimensions.Y + conf.EdgeBorderY
	} else if y > disp.Dimensions.Y+conf.EdgeBorderY {
		y += conf.BorderY / 2
	}
	if bry >= disp.Dimensions.H-conf.EdgeBorderY {
		bry = disp.Dimensions.H - conf.EdgeBorderY
		h = bry - y + disp.Dimensions.Y
	} else if bry < disp.Dimensions.H-conf.EdgeBorderY {
		bry -= conf.BorderY / 2
		h = bry - y
	}

	key := fmt.Sprintf("%.0fx%.0f", disp.Dimensions.W, disp.Dimensions.H)
	twk, ok := conf.GetTweak(key, win.Class)
	if !ok {
		twk, ok = conf.GetTweak(key, "__ALL__")
		if !ok {
			twk = config.Tweak{}
		}
	}

	x -= twk.NudgeX
	y -= twk.NudgeY
	w -= twk.ShrinkW
	h -= twk.ShrinkH

	return dimensions{
		X: math.Floor(x),
		Y: math.Floor(y),
		W: math.Floor(w),
		H: math.Floor(h),
	}, nil
}
