package main

import (
  "io/ioutil"
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

func and (conds ...bool) bool {
  for _, cond := range conds {
    if !cond { return false }
  }
  return true
}

func or (conds ...bool) bool {
  for _, cond := range conds {
    if cond { return true }
  }
  return false
}

func tbPrint(row, col int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(col, row, c, fg, bg)
    col++
  }
}

func log (msg string) {
  err := ioutil.WriteFile("log", []byte(msg), 0644)
  panicIfError(err)
}
