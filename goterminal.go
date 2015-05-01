package main

import (
	"github.com/nsf/termbox-go"
)

type Status struct {
	mode    string
	cur_pos CursorsPos
	cmd_pos int
	command string
}

var status Status

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	InitLayout(&status)

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

loop:
	for {
		ev := <-event_queue
		if ev.Type != termbox.EventKey {
			continue
		}

		if ev.Key == termbox.KeyEsc {
			status.mode = "IDLE"
			Reset(&status)
		}

		if status.mode == "IDLE" {
			foo_key, ok_key := idle_mode_map[rune(ev.Key)]
			foo, ok := idle_mode_map[rune(ev.Ch)]

			if ok_key {
				foo_key(&status)
			} else if ok {
				foo(&status)
			}
		}

		if status.mode == "CMD" {
			foo_key, ok_key := cmd_mode_map[rune(ev.Key)]
			foo, ok := cmd_mode_map[rune(ev.Ch)]

			if ok_key {
				foo_key(&status)
			} else if ok {
				foo(&status)
			} else {
				command_mode_line(&status, ev.Ch)
			}
		}

		if ev.Key == termbox.KeyCtrlQ {
			break loop
		}

		Draw()
	}
}
