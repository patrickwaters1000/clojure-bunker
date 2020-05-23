package main

import (
  "strings"
)

type ReplBuffer struct {
  text string
}

func NewReplBuffer () *ReplBuffer {
  return &ReplBuffer{""}
}

func (b ReplBuffer) render (w Window) {
  lines := strings.Split(b.text, "\n")
  for i, line := range lines {
    w.Print(i, 0, fg1, bg1, line)
  }
}

func (b *ReplBuffer) handle (event []string) {
}

func (b ReplBuffer) stringify () string {
  return "Repl buffer:\n" + b.text + "\n\n"
}
