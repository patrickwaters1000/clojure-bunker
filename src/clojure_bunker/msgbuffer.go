package main

import (
  "strings"
)

type MsgBuffer struct {
  text []string
}

func NewMsgBuffer () *MsgBuffer {
  return &MsgBuffer{[]string{}}
}

func (b *MsgBuffer) handle (cmd []string) {
  b.text = cmd
}

func (b MsgBuffer) render (w Window) {
  s := strings.Join(b.text, " ")
  w.Print(0, 0, fg1, bg1, s)
}

func (b MsgBuffer) stringify () string {
  return "msgBuffer:\n" + strings.Join(b.text, "\n") + "\n\n"
}
