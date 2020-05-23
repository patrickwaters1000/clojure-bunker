package main

import (
)

type MiniBuffer struct {
  data string
  prompt string
}

func NewMiniBuffer () *MiniBuffer {
  return &MiniBuffer{"", ""}
}

func (b *MiniBuffer) reset (prompt string) {
  b.data = ""
  b.prompt = prompt
}

func (b *MiniBuffer) handle (cmd []string) {
  switch cmd[0] {
  case "delete":
    l := len(b.data)
    b.data = b.data[:l-1]
  case "append":
    b.data += cmd[1]
  }
}

func (b MiniBuffer) render(w Window) {
  w.Print(0, 0, fg1, bg1, b.prompt + b.data)
}

func (b MiniBuffer) stringify () string {
  return "miniBuffer: \n" + b.data + b.prompt + "\n\n"
}
