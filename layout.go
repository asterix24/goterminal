package main

import (
	"github.com/nsf/termbox-go"
)

type CursorsPos struct {
	x int
	y int
}

type Status struct {
	mode        string
	cur_pos     CursorsPos
	cmd_cur_pos CursorsPos
}

const CMD_LINE = 1
const BAR_POS = 2
const FOOTER_HIGH = 4

func layout() {
	w, h := termbox.Size()
	for x := 0; x < w; x++ {
		termbox.SetCell(x, h-BAR_POS, ' ', termbox.ColorDefault, termbox.ColorBlue)
	}

	termbox.Flush()
}

func idle_mode_up(status *Status) int {
	if status.cur_pos.y <= 0 {
		return 0
	}
	status.cur_pos.y -= 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}
func idle_mode_down(status *Status) int {
	_, h := termbox.Size()

	if status.cur_pos.y > h-FOOTER_HIGH {
		return 0
	}
	status.cur_pos.y += 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}
func idle_mode_right(status *Status) int {
	w, _ := termbox.Size()
	if status.cur_pos.x >= w {
		return 0
	}
	status.cur_pos.x += 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}
func idle_mode_left(status *Status) int {
	if status.cur_pos.x <= 0 {
		return 0
	}
	status.cur_pos.x -= 1
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	return 0
}

func command_mode(status *Status) int {
	status.mode = "CMD"
	_, h := termbox.Size()
	status.cmd_cur_pos.x = 0
	status.cmd_cur_pos.y = h
	termbox.SetCell(0, h, ':', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(0, h)
	termbox.Flush()

	return 0
}

func command_mode_line(status *Status, s rune) int {
	w, h := termbox.Size()
	if status.cmd_cur_pos.x <= w-1 {
		status.cmd_cur_pos.x += 1
		termbox.SetCell(status.cmd_cur_pos.x, h-CMD_LINE, s, termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCursor(status.cmd_cur_pos.x+1, h-CMD_LINE)
	}
	termbox.Flush()
	return 0
}

type callback func(status *Status) int

var idle_mode_map = map[rune]callback{
	rune(termbox.KeyArrowUp):    idle_mode_up,
	rune(termbox.KeyArrowDown):  idle_mode_down,
	rune(termbox.KeyArrowRight): idle_mode_right,
	rune(termbox.KeyArrowLeft):  idle_mode_left,
	rune('k'):                   idle_mode_up,
	rune('j'):                   idle_mode_down,
	rune('l'):                   idle_mode_right,
	rune('h'):                   idle_mode_left,
	rune(':'):                   command_mode,
}

var status Status

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	layout()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	status.mode = "IDLE"
	status.cur_pos = CursorsPos{0, 0}

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
			command_mode_line(&status, ev.Ch)
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
