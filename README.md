## Purpose

This small app allows you to place a window where you want it based on a grid system. It can also raise existing windows based on their class and/or window title.

I have only tested on KDE Neon.

## Usage

```
Usage:
  -always
    	always run the command
  -class string
    	the class name of the window to match
  -dims string
    	dimensions: grid:left,top,width,height ex. 12x10:1,1,12,10
  -list
    	show a list of current windows classes and titles
  -minimise
    	minimise the frontmost window
  -title string
    	the full or partial title of the window to match

```

To place the currently active window on the top left quarter of the screen add the following command as a custom keyboard shortcut in the Settings app.

```
kazinga -dims 2x2:1,1,1,1
```

The grid is a flexible way of dividing up the screen. For instance I use 12x2 to have 3:4:3 column layout for my windows. The grid and positions all start at 1 not 0.

To launch (or raise an existing window) add this command as a custom keyboard shortcut.

```
kazinga -class qterminal -title vim launch_vim_in_tmux.sh
```

I keep my notes in Vimwiki in a tmux session and I have a script which switches to the session. So, to raise or launch the window with my notes, I have a custom keyboard shortcut for it.

```
kazinga -class qterminal -title notes -always switch_to_notes.sh
```

## Installation

Simply download the binary from the releases page and move it to somewhere in your PATH.

## Configuration

By default a configuration file is created in $HOME/.config/kazinga/kazinga.conf. It's possible to adjust the borders generally and by class name. Just look at the file that's generated on first run, it's pretty straight-forward.

Here is an excerpt from my KMonad config to get an idea of what it can do:

```
  sz1 (cmd-button "kazinga -dims 2x2:1,1,1,1")
  sz2 (cmd-button "kazinga -dims 2x2:1,1,2,1")
  sz3 (cmd-button "kazinga -dims 2x2:2,1,1,1")
  sz4 (cmd-button "kazinga -dims 2x2:1,1,1,2")
  sz5 (cmd-button "kazinga -dims 1x1:1,1,1,1")
  sz6 (cmd-button "kazinga -dims 2x2:2,1,1,2")
  sz7 (cmd-button "kazinga -dims 2x2:1,2,1,1")
  sz8 (cmd-button "kazinga -dims 2x2:1,2,2,1")
  sz9 (cmd-button "kazinga -dims 2x2:2,2,1,1")
  sz0 (cmd-button "kazinga -dims 1x1:1,1,1,1")
  hid (cmd-button "kazinga -minimise")
  btw (cmd-button "kazinga -class bitwarden bitwarden")
  fil (cmd-button "kazinga -class dolphin dolphin")
  ffx (cmd-button "kazinga -class firefox firefox")
  kty (cmd-button "kazinga -class kitty kitty")
  mus (cmd-button "kazinga -class elisa elisa")
  pod (cmd-button "kazinga -class kasts kasts")
  prf (cmd-button "kazinga -class systemsettings systemsettings")
  thu (cmd-button "kazinga -class thunderbird-nightly thunderbird")

```

## Building

Install a recent version of Go and run

```
go install
```

Make sure that $GOPATH/bin (usually ~/go/bin) is in your PATH

## Bugs

If you think you have found a bug please file an issue. I will do my best to fix it in a timely manner.

Duplicate issues will be closed without comment.

## Feature requests

The program is feature complete enough for my needs so I most likely won't implement anything new but I'm always ready to review merge requests.

