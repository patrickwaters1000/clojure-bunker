package main

import (
  termbox "github.com/nsf/termbox-go"
)

type Window struct {
  canSelect bool
  rows int
  cols int
  pos_r int // Allows to control what part of the buffer
  pos_c int // if in the window
  screen_r int
  screen_c int
  buffer Buffer
}

func (w Window) Print (row, col int, fg, bg termbox.Attribute, msg string) {
  r := row - w.pos_r
  c := col - w.pos_c
  if and(0 <= r,
         0 <= c,
         r < w.rows,
         c < w.cols) {
    tbPrint(r + w.screen_r, c + w.screen_c, fg, bg, msg)
  }
}

func NewWindow (rows, cols int) *Window {
  return &Window{
    canSelect: true,
    rows: rows,
    cols: cols,
    pos_r: 0,
    pos_c: 0,
    screen_r: 0,
    screen_c: 0,
    buffer: nil,
  }
}

func (w *Window) SplitHorizontally (row int) (*Window, *Window) {
  if row < 0 || row >= w.rows {
    panic("Something is wrong")
  }
  w1 := &Window{
    canSelect: true,
    rows: row,
    cols: w.cols,
    pos_r: w.pos_r,
    pos_c: w.pos_c,
    screen_r: w.screen_r,
    screen_c: w.screen_c,
    buffer: w.buffer,
  }
  w2 := &Window{
    canSelect: true,
    rows: w.rows - row,
    cols: w.cols,
    pos_r: 0,
    pos_c: 0,
    screen_r: w.screen_r + row,
    screen_c: w.screen_c,
    buffer: nil,
  }
  return w1, w2
}

func (w *Window) SplitVertically(col int) (*Window, *Window) {
  if col < 0 || col >= w.cols {
    panic("Something is wrong")
  }
  w1 := &Window{
    canSelect: true,
    rows: w.rows,
    cols: col,
    pos_r: w.pos_r,
    pos_c: w.pos_c,
    screen_r: w.screen_r,
    screen_c: w.screen_c,
    buffer: w.buffer,
  }
  w2 := &Window{
    canSelect: true,
    rows: w.rows,
    cols: w.cols - col,
    pos_r: 0,
    pos_c: 0,
    screen_r: w.screen_r,
    screen_c: w.screen_c + col,
    buffer: nil,
  }
  return w1, w2
}

func (w *Window) Center (row int) {
  w.pos_r = row - w.rows / 2
}
