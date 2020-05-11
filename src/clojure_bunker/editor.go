package main

import (
  "io/ioutil"
  "strings"
)

type Editor struct {
  buffers []Handler
  lastCmd []string
  active int
  replClient *Client
  replBuffer *ReplBuffer
}

func NewEditor() *Editor {
  return &Editor{
    buffers: []Handler{},
    active: 0,
    lastCmd: []string{},
    replClient: NewClient(),
    replBuffer: nil,
  }
}

type ReplBuffer struct {
  text string
}

func (b ReplBuffer) render (w Window) {
  lines := strings.Split(b.text, "\n")
  for i, line := range lines {
    w.Print(i, 0, fg1, bg1, line)
  }
}

func (b *ReplBuffer) handleEvent (event []string) error {
  // Add movement commands, etc. here.
  return nil
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

func stringifySubtree (n *TreeNode) string {
  token := NewToken("root","")
  subtree := NewTree(token)
  subtree.Root.Children = []*TreeNode{n}
  var msg string = ""
  var row int = 0
  var col int = 0
  traverseFn := func (n *TreeNode) {
    token := n.Data.(*Token)
    if token.Row > row {
      for i:=0; i<token.Row-row; i++ {
        msg += "\n"
        row += 1
      }
      col = 0
    }
    for j:=0; j<token.Col-col; j++ {
      msg += " "
      col += 1
    }
    msg += token.Value
    col += len(token.Value)
  }
  subtree.DepthFirstTraverseNoRoot(traverseFn)
  return msg
}

func (e *Editor) writeActiveBuffer(fname string) {
  msg := stringifySubtree(e.buffers[e.active].(*Buffer).tree.Root)
  err := ioutil.WriteFile(fname, []byte(msg), 0644)
  panicIfError(err)
}

func (e *Editor) loadFile(fname string) {
  data, err := ioutil.ReadFile(fname)
  panicIfError(err)
  tree := parseClj(data)
  tree.Active = tree.Root.Children[0]
  mapSyntaxTree(tree)
  b := &Buffer{fname, "normal", tree}
  b.setCursor(true)
  e.buffers = append(e.buffers, b)
  e.active = len(e.buffers) - 1
}

func (e *Editor) replConnect(port string) {
  e.replClient.Connect(port)
  e.replBuffer = &ReplBuffer{""}
  e.buffers = append(e.buffers, e.replBuffer)
}

func (e *Editor) replEval() {
  code := stringifySubtree(e.buffers[e.active].(*Buffer).tree.Active)
  id := e.replClient.nextId
  e.replClient.Send(code)
  go func() {
    text, _ := e.replClient.GetResponse(id)
    e.replBuffer.text += (text + "\n\n")
  }()
}

func (e *Editor) handleEvent (event []string) error {
  e.lastCmd = event
  cmd := event[0]
  switch cmd {
  case "new-buffer": e.newBuffer(event[1])
  case "next-buffer": e.nextBuffer()
  case "kill-buffer": e.killBuffer()
  case "buffer": e.buffers[e.active].handleEvent(event[1:])
  case "write": e.writeActiveBuffer(event[1])
  case "load": e.loadFile(event[1])
  case "repl":
    switch event[1] {
    case "eval": e.replEval()
    case "connect": e.replConnect(event[2])
    }
  }
  return nil
}

func (e Editor) render(w Window) {
  for i, s := range e.lastCmd {
    w.Print(w.rows - 5 + i, 0, fg1, bg1, s)
  }
  if len(e.buffers) > 0 {
    e.buffers[e.active].render(w)
  }
}
