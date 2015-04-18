package main

import (
	"github.com/nsf/termbox-go"
)

type Status struct {
	mode        string
	cur_pos     CursorsPos
	cmd_cur_pos CursorsPos
}

var status Status

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	init_layout(&status)

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
			termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
			// clean all..
		}

		if status.mode == "IDLE" {
			foo, ok := idle_mode_map[rune(ev.Key)]
			if ok {
				foo(&status)
			} else {
				foo, ok := idle_mode_map[rune(ev.Ch)]
				if ok {
					foo(&status)
				}
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

		//command_line(status, &cmd, ev.Ch)

		//if ev.Key == termbox.KeyEsc {
		//	editmode = true
		//	commandmode = false
		//}

		//if editmode {
		//	if ev.Ch == ':' && !commandmode {
		//		commandmode = true
		//	}

		//	if commandmode {
		//		command_line(&cmd, ev.Ch)
		//	} else {
		//		termbox.SetCell(10, 10, rune('t'), termbox.ColorDefault, termbox.ColorDefault)
		//	}
		//}

		if ev.Key == termbox.KeyCtrlQ {
			break loop
		}
		layout()
	}
}
