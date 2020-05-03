package main

import (
)

type Editor struct {
  buffers []Handler
  lastCmd []string
  active int
}

func NewEditor() *Editor {
  return &Editor{
    buffers: []Handler{},
    active: 0,
    lastCmd: []string{},
  }
}

func (editor *Editor) newBuffer(name string) {
  b := NewBuffer(name)
  editor.buffers = append(editor.buffers, b)
  editor.active = len(editor.buffers) - 1
}

func (e *Editor) nextBuffer() {
  n := len(e.buffers)
  e.active = mod(e.active + 1, n)
}

func (e *Editor) killBuffer() {
  e.buffers = append(
    e.buffers[:e.active],
    e.buffers[e.active+1:]...)
  n := len(e.buffers)
  if n > 0 {
    e.active = mod(e.active, n)
  } else {
    e.active = 0
  }
}

func (e *Editor) handleEvent (event []string) error {
  e.lastCmd = event
  cmd := event[0]
  switch cmd {
  case "new-buffer": e.newBuffer(event[1])
  case "next-buffer": e.nextBuffer()
  case "kill-buffer": e.killBuffer()
  case "buffer": e.buffers[e.active].handleEvent(event[1:])
  }
  return nil
}

func (e Editor) render() {
  for i, s := range e.lastCmd {
    tbPrint(21+i, 0, fg1, bg1, s)
  }
  if len(e.buffers) > 0 {
    e.buffers[e.active].render()
  }
}
