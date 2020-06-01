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

func (b *MiniBuffer) Delete () {
  l := len(b.data)
  if l > 0 {
    b.data = b.data[:l-1]
  }
}

func (b *MiniBuffer) Append (s string) {
  b.data += s
}

func (b *MiniBuffer) handle (cmd []string) {
}

func (b MiniBuffer) render(w Window) {
  w.Print(0, 0, fg1, bg1, b.prompt + b.data)
}

func (b MiniBuffer) stringify () string {
  return "miniBuffer: \n" + b.data + b.prompt + "\n\n"
}
