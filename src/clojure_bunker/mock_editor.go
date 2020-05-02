package main

import (
  termbox "github.com/nsf/termbox-go"
)

type MockEditor struct {
  lastCmd []string
}

func (editor *MockEditor) handleEvent (event []string) error {
  editor.lastCmd = event
  return nil
}

func NewEditor() *MockEditor {
  return &MockEditor{[]string{}}
}

func (editor MockEditor) render() {
  for i, s := range editor.lastCmd {
    tbPrint(i, 0, termbox.ColorWhite, termbox.ColorBlack, s)
  }
}
