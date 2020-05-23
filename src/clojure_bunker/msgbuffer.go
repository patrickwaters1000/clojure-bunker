package main

import (
  "strings"
)

type MsgBuffer struct {
  textRows []string
}

func NewMsgBuffer () *MsgBuffer {
  return &MsgBuffer{[]string{}}
}

func (b *MsgBuffer) handle (cmd []string) {
  b.textRows = cmd
}

func (b MsgBuffer) render (w Window) {
  for i, s := range b.textRows {
    w.Print(i, 0, fg1, bg1, s)
  }
}

func (b MsgBuffer) stringify () string {
  return "msgBuffer:\n" + strings.Join(b.textRows, "\n") + "\n\n"
}
