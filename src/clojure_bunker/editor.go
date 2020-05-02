package main

import (
)

type Editor struct {
  buffers []Handler
  active int
}

func NewEditor() *Editor {
  return &Editor{
    buffers: []Handler{},
    active: 0,
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
  e.active = mod(e.active, n)
}

func (e *Editor) handleEvent (event []string) error {
  cmd := event[0]
  switch cmd {
  case "new-buffer": e.newBuffer(event[1])
  case "next-buffer": e.nextBuffer()
  case "kill-buffer": e.killBuffer()
  }
  return nil
}

func (e Editor) render() {
  if len(e.buffers) > 0 {
    e.buffers[e.active].render()
  }
}
