package main

import (
  termbox "github.com/nsf/termbox-go"
)


func panicIfError(err interface{}) {
  if err != nil {
    panic(err)
  }
}

func mod(n, d int) int {
  m := n % d
  if m < 0 {
    m += d
  }
  return m
}

func tbPrint(row, col int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(col, row, c, fg, bg)
    col++
  }
}
