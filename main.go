package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"kazinga/kwin"
)

var version string

func main() {
	f := flag.FlagSet{}
	alwaysRun := f.Bool("always", false, "always run the command")
	class := f.String("class", "", "the class name of the window to match")
	dims := f.String("dims", "", "dimensions: grid:left,top,width,height ex. 12x10:1,1,12,10")
	list := f.Bool("list", false, "show a list of current windows classes and titles")
	minimise := f.Bool("minimise", false, "minimise the frontmost window")
	title := f.String("title", "", "the full or partial title of the window to match")
	showVersion := f.Bool("version", false, "show version number")
	f.Parse(os.Args[1:])
	commandArgumentsLength := len(f.Args())

	if *showVersion {
		fmt.Println(version)
		return
	}

	if *class == "" && *title == "" && commandArgumentsLength > 0 && !*alwaysRun {
		f.Usage()
		os.Exit(1)
	}

	kwin.Init()
	defer kwin.Close()

	if *list {
		for _, win := range kwin.ListWindows() {
			fmt.Println(win)
		}
		return
	}

	var id string
	var err error
	if *class != "" || *title != "" {
		var command string
		if commandArgumentsLength > 0 {
			command = strings.Join(f.Args(), " ")
		}
		id, err = kwin.RunOrRaise(*class, *title, command, *alwaysRun)
		if err != nil {
			fmt.Println("Error:\n" + err.Error())
			return
		}
		if id == "" {
			fmt.Println("Error:\ninvalid window id")
			return
		}
	}

	if *minimise {
		kwin.MinimiseWindow("")
		return
	}

	if *dims != "" {
		kwin.ResizeWindow(id, *dims)
		return
	}
}
