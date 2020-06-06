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
  sendEventFn func([]string)
}

func NewEditor(sendEventFn func([]string)) *Editor {
  rows, cols := get_winsize()
  w_full := NewWindow(rows, cols)
  w_top, w_bottom := w_full.SplitHorizontally(rows - 3)
  w_bl, w_br := w_bottom.SplitVertically(cols / 2)
  miniBuffer := NewMiniBuffer()
  msgBuffer := NewMsgBuffer()
  w_bl.buffer = miniBuffer
  w_br.buffer = msgBuffer
  w_bl.canSelect = false
  w_br.canSelect = false
  return &Editor{
    windows: []*Window{w_top, w_bl, w_br},
    activeWindow: 0,
    buffers: []Buffer{miniBuffer, msgBuffer},
    replClient: nil,
    sendEventFn: sendEventFn,
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

func (e *Editor) setActiveBufferFile (file string) {
  b := e.getActiveBuffer()
  b.(*CodeBuffer).file = file
}

func (e *Editor) writeActiveBuffer() {
  b := e.getActiveBuffer()
  msg := stringifySubtree(b.(*CodeBuffer).tree.Root)
  err := ioutil.WriteFile(b.(*CodeBuffer).file, []byte(msg), 0644)
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
  b.file = fname
  b.tree = tree
  b.setCursor(true)
  e.buffers = append(e.buffers, b)
  e.getActiveWindow().buffer = b
}

func (e *Editor) replConnect(port string) {
  e.replClient = NewClient()
  e.replClient.Connect(port)
  replBuffer := NewReplBuffer()
  e.buffers = append(e.buffers, replBuffer)
}

func (e *Editor) replEval() {
  activeBuffer := e.getActiveBuffer().(*CodeBuffer)
  activeNode := activeBuffer.tree.GetActive()
  n := NewTreeNode(nil) // Awkwardly, stringifySubtree only uses
  // the tree strictly below the given node
  n.Children = append(n.Children, activeNode)
  code := stringifySubtree(n)
  id := e.replClient.nextId
  e.replClient.Send(code)
  go func() {
    text, _ := e.replClient.GetResponse(id)
    e.getReplBuffer().text += (text + "\n\n")
    e.sendEventFn([]string{"refresh"})
  }()
}

func (e *Editor) nextWindow () {
  n := len(e.windows)
  e.activeWindow = mod(e.activeWindow + 1, n)
  if e.getActiveWindow().canSelect == false {
    e.nextWindow()
  }
}

func (e *Editor) splitActiveWindowVertically () {
  w := e.getActiveWindow()
  iw := e.activeWindow // index of it
  ib := e.getActiveBufferIndex()
  nb := len(e.buffers)
  wl, wr := w.SplitVertically(w.cols / 2)
  wr.buffer = e.buffers[mod(ib + 1, nb)]
  e.windows = append(
    e.windows[:iw],
    append(
      []*Window{wl, wr},
      e.windows[iw + 1:]...)...)
}

func (e *Editor) handle (event []string) {
  e.getMsgBuffer().handle(event)
  switch event[0] {
  case "buffer": e.getActiveBuffer().handle(event[1:])
  case "new-buffer": e.newBuffer(event[1])
  case "next-buffer": e.nextBuffer()
  case "kill-buffer": e.killBuffer()
  case "set-file": e.setActiveBufferFile(event[1])
  case "write-file": e.writeActiveBuffer()
  case "load-file": e.loadFile(event[1])
  case "set-mode": e.getActiveBuffer().handle(event)
  case "center-window": e.CenterWindow()
  case "repl-eval": e.replEval()
  case "repl-connect": e.replConnect(event[1])
  case "window":
    switch event[1] {
    case "next": e.nextWindow()
    case "split-vertical": e.splitActiveWindowVertically()
    }
  default: panic("Not found: " + event[0])
  }
}
