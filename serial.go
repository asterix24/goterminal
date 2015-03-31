// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func openChannel(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("panel", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		fmt.Fprintln(v, "Open")
		if err := g.SetCurrentView("panel"); err != nil {
			return err
		}
	}
	return nil
}

func settings(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("settings", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Editable = true
		fmt.Fprintln(v, "settings")
		if err := g.SetCurrentView("settings"); err != nil {
			return err
		}
	}
	return nil
}

func returnLog(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("settings"); err != nil {
		return err
	}
	if err := g.SetCurrentView("log"); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, returnLog); err != nil {
		return err
	}

	if err := g.SetKeybinding("log", gocui.KeyCtrlO, gocui.ModNone, openChannel); err != nil {
		return err
	}

	if err := g.SetKeybinding("log", gocui.KeyCtrlS, gocui.ModNone, settings); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("log", 0, 2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		v.Editable = true
		v.Wrap = true
		v.Frame = false
		v.SelBgColor = gocui.ColorGreen
		if err := g.SetCurrentView("log"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("side", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Frame = false
		v.BgColor = gocui.ColorBlue
		fmt.Fprintln(v, "[Item 1] [Item 2] [Item 3]")
	}

	v, err := g.SetView("legend", maxX-25, 2, maxX-5, 8)
	if err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "Space: New View")
		fmt.Fprintln(v, "Tab: Next View")
		fmt.Fprintln(v, "← ↑ → ↓: Move View")
		fmt.Fprintln(v, "^C: Exit")
	}
	return nil
}

func main() {
	var err error

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.ShowCursor = true

	err = g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
