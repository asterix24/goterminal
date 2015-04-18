package main

import (
	"github.com/nsf/termbox-go"
)

type CursorsPos struct {
	x int
	y int
}

const CMD_LINE = 1
const BAR_POS = 2
const FOOTER_HIGH = 4

func init_layout(status *Status) {
	message := "Asterix Emulatore Seriale!"

	w, h := termbox.Size()
	status.cur_pos = CursorsPos{w / 2, h / 2}
	status.mode = "IDLE"
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	status.cur_pos.x -= len(message) / 2
	for _, v := range message {
		termbox.SetCell(status.cur_pos.x, status.cur_pos.y, rune(v), termbox.ColorDefault, termbox.ColorDefault)
		status.cur_pos.x++
	}

	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
	termbox.HideCursor()

	draw()
}

func reset_view(status *Status) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(status.cur_pos.x, status.cur_pos.y)
}

func draw() {
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
	termbox.SetCursor(1, h)
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

func cmd_mode_up(status *Status) int {
	return 0
}
func cmd_mode_down(status *Status) int {
	return 0
}
func cmd_mode_right(status *Status) int {
	w, h := termbox.Size()
	if status.cmd_cur_pos.x >= w {
		return 0
	}
	status.cmd_cur_pos.x += 1
	termbox.SetCursor(status.cmd_cur_pos.x, h-CMD_LINE)
	return 0
}
func cmd_mode_left(status *Status) int {
	_, h := termbox.Size()
	if status.cmd_cur_pos.x <= 1 {
		return 0
	}
	status.cmd_cur_pos.x -= 1
	termbox.SetCursor(status.cmd_cur_pos.x, h-CMD_LINE)
	return 0
}

func cmd_mode_canc(status *Status) int {
	_, h := termbox.Size()
	if status.cmd_cur_pos.x <= 0 {
		return 0
	}
	termbox.SetCell(status.cmd_cur_pos.x, h-CMD_LINE, ' ', termbox.ColorDefault, termbox.ColorDefault)
	status.cmd_cur_pos.x -= 1
	termbox.SetCell(status.cmd_cur_pos.x, h-CMD_LINE, ' ', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(status.cmd_cur_pos.x, h-CMD_LINE)
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

var cmd_mode_map = map[rune]callback{
	rune(termbox.KeyArrowUp):    cmd_mode_up,
	rune(termbox.KeyArrowDown):  cmd_mode_down,
	rune(termbox.KeyArrowRight): cmd_mode_right,
	rune(termbox.KeyArrowLeft):  cmd_mode_left,
	rune(termbox.KeyBackspace):  cmd_mode_canc,
	rune(termbox.KeyBackspace2): cmd_mode_canc,
}
