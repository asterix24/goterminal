package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

var editmode bool
var commandmode bool

func layout() {
	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '_', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, h-2, ' ', termbox.ColorDefault, termbox.ColorBlue)
	}

	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func foo(s string) int {
	fmt.Println(s)
	return 0
}

type callback func(s string) int

func main() {

	m := make(map[rune]callback)
	call := foo
	call("sdflsd")

	m['a'] = foo
	m['a']("mappa..")

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	_, h := termbox.Size()
	layout()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

loop:
	for {
		ev := <-event_queue
		if ev.Type != termbox.EventKey {
			continue
		}

		if ev.Key == termbox.KeyEsc {
			editmode = true
			commandmode = false
		}

		if editmode {
			if ev.Ch == ':' && !commandmode {
				commandmode = true
			}

			if commandmode {
				termbox.SetCell(0, h-1, ev.Ch, termbox.ColorDefault, termbox.ColorDefault)
				termbox.SetCursor(0, h-1)
			} else {
				termbox.SetCell(10, 10, rune('t'), termbox.ColorDefault, termbox.ColorDefault)
			}
		}

		if ev.Key == termbox.KeyCtrlQ {
			break loop
		}
		layout()
	}
}
