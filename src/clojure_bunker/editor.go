package main

import (
  termbox "github.com/nsf/termbox-go"
  "io/ioutil"
  "strings"
)

type Editor struct {
  windows []*Window
  activeWindow int
  buffers []Buffer
  replClient *Client
}

func NewEditor() *Editor {
  rows, cols := get_winsize()
  w_full := NewWindow(rows, cols)
  w_left, w_right := w_full.SplitVertically(cols / 2)
  w_right_top, w_right_bottom := w_right.SplitHorizontally(rows / 2)
  miniBuffer := NewMiniBuffer()
  msgBuffer := NewMsgBuffer()
  w_right_top.buffer = miniBuffer
  w_right_bottom.buffer = msgBuffer
  return &Editor{
    windows: []*Window{w_left, w_right_top, w_right_bottom},
    activeWindow: 0,
    buffers: []Buffer{miniBuffer, msgBuffer},
    replClient: nil,
  }
}

func (e *Editor) getActiveWindow () *Window {
  return e.windows[e.activeWindow]
}

func (e *Editor) getActiveBuffer () Buffer {
  return e.getActiveWindow().buffer
}

func (e *Editor) getMiniBuffer () *MiniBuffer {
  for _, b := range e.buffers {
    switch b.(type) {
    case *MiniBuffer: return b.(*MiniBuffer)
    }
  }
  panic("Not found")
}

func (e *Editor) getMsgBuffer () *MsgBuffer {
  for _, b := range e.buffers {
    switch b.(type) {
    case *MsgBuffer: return b.(*MsgBuffer)
    }
  }
  panic("Not found")
}

func (e *Editor) getReplBuffer () *ReplBuffer {
  for _, b := range e.buffers {
    switch b.(type) {
    case *ReplBuffer: return b.(*ReplBuffer)
    }
  }
  panic("Not found")
}

func (e Editor) logState () {
  var s []string
  for _,b := range e.buffers {
    s = append(s, b.stringify())
  }
  joined := strings.Join(s,"\n")
  ioutil.WriteFile("state", []byte(joined), 0644)
}

func (e *Editor) CenterWindow () {
  b := e.getActiveBuffer().(*CodeBuffer)
  r := b.tree.GetActive().Data.(*Token).Row
  e.getActiveWindow().Center(r)
}

func (e *Editor) handle (event []string) error {
  e.getMsgBuffer().handle(event)
  cmd := event[0]
  switch cmd {
  case "buffer": e.getActiveBuffer().handle(event[1:])
  case "minibuffer": e.getMiniBuffer().handle(event[1:])
  // case "window": e.getActiveWindow().handle(event[1:])
  case "new-buffer": e.newBuffer(event[1])
  case "next-buffer": e.nextBuffer()
  case "kill-buffer": e.killBuffer()
  case "write-file": e.writeActiveBuffer(event[1])
  case "load-file": e.loadFile(event[1])
  case "set-mode": e.getActiveBuffer().handle(event)
  case "center-window": e.CenterWindow()
  case "repl":
    switch event[1] {
    case "eval": e.replEval()
    case "connect": e.replConnect(event[2])
    }
  default: panic("Not found")
  }
  return nil
}

func (e *Editor) render() {
  termbox.Clear(bg1, bg1)
  for _, w := range e.windows {
    b := w.buffer
    if b != nil {
      b.render(*w)
    }
  }
  termbox.Flush()
}

func (editor *Editor) newBuffer(name string) {
  b := NewCodeBuffer()
  b.name = name
  editor.buffers = append(editor.buffers, b)
  editor.getActiveWindow().buffer = b
}

func (e *Editor) getActiveBufferIndex () int {
  activeBuffer := e.getActiveBuffer()
  for i, b := range e.buffers {
    if b == activeBuffer {
      return i
    }
  }
  panic("Not found")
}

func (e *Editor) nextBuffer() {
  n := len(e.buffers)
  i1 := e.getActiveBufferIndex()
  i2 := mod(i1 + 1, n)
  e.getActiveWindow().buffer = e.buffers[i2]
}

func (e *Editor) killBuffer() {
  i := e.getActiveBufferIndex()
  e.buffers = append(e.buffers[:i], e.buffers[i+1:]...)
  e.getActiveWindow().buffer = nil
}

func (e *Editor) writeActiveBuffer(fname string) {
  msg := stringifySubtree(e.getActiveBuffer().(*CodeBuffer).tree.Root)
  err := ioutil.WriteFile(fname, []byte(msg), 0644)
  panicIfError(err)
}

func (e *Editor) loadFile(fname string) {
  data, err := ioutil.ReadFile(fname)
  panicIfError(err)
  tree := parseClj(data)
  tree.Path = []*TreeNode{tree.Root}
  _ = tree.DownFirst()
  b := NewCodeBuffer()
  b.name = fname
  b.tree = tree
  b.setCursor(true)
  e.buffers = append(e.buffers, b)
  e.getActiveWindow().buffer = b
}

func (e *Editor) replConnect(port string) {
  e.replClient.Connect(port)
  replBuffer := NewReplBuffer()
  e.buffers = append(e.buffers, replBuffer)
}

func (e *Editor) replEval() {
  activeBuffer := e.getActiveBuffer().(*CodeBuffer)
  activeNode := activeBuffer.tree.GetActive()
  code := stringifySubtree(activeNode)
  id := e.replClient.nextId
  e.replClient.Send(code)
  go func() {
    text, _ := e.replClient.GetResponse(id)
    e.getReplBuffer().text += (text + "\n\n")
  }()
}


