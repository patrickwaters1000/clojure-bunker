package main

import (
  "fmt"
  //"errors"
  termbox "github.com/nsf/termbox-go"
)

type Buffer interface {
  handle([]string)
  render(Window)
  stringify() string
}

type CodeBuffer struct {
  name string
  mode string
  tree *Tree
}

func NewCodeBuffer () *CodeBuffer {
  rootToken := NewToken("root", "")
  tree := NewTree(rootToken)
  leafToken := NewToken("leaf", "")
  tree.AppendChild(leafToken)
  return &CodeBuffer{"", "normal", tree}
}

func stringifyTree (t *Tree) string {
  msg := ""
  traverseFn := func (n *TreeNode) {
    d := n.Data.(*Token)
    active := n == t.GetActive()
    msg += fmt.Sprintf(
      "class:%s value:%s children:%d selected:%v active:%v row:%d col:%d\n",
      d.Class, d.Value, len(n.Children), d.Selected, active, d.Row, d.Col)
  }
  t.DepthFirstTraverse(traverseFn)
  return msg
}

func logTree (t *Tree) {
  log(stringifyTree(t))
}

func (b CodeBuffer) stringify () string {
  return fmt.Sprintf("CodeBuffer:\nname: %s\nmode: %s\ntree:\n%s",
    b.name, b.mode, stringifyTree(b.tree)) + "\n\n"
}

func isBuiltIn(s string) bool {
  switch s {
  case "let":  return true
  case "def":  return true
  case "defn": return true
  case "for":  return true
  case "and":  return true
  case "or":   return true
  case "cond": return true
  default: return false
  }
}

func (b CodeBuffer) render (w Window) {
  w.Print(0, 0, fg1, bg1, b.name)
  logTree(b.tree)
  traverseFn := func (n *TreeNode) {
    t := n.Data.(*Token)
    bg := termbox.ColorBlack
    fg := t.Color
    if t.Selected {
      bg = fg
      fg = termbox.ColorBlack
    }
    w.Print(t.Row + 2, t.Col, fg, bg, t.Value)
  }
  a := NewArtist(w)
  a.Root(b.tree.Root) // sets the positions and colors of tokens
  b.tree.DepthFirstTraverseNoRoot(traverseFn)
}

// TODO refactor to use Tree methods
func (b *CodeBuffer) setCursor (v bool) {
  a := b.tree.GetActive()
  at := a.Data.(*Token)
  at.Selected = v
  if a.Data.(*Token).Class == "open" {
    c := a.Children[len(a.Children) - 1]
    ct := c.Data.(*Token)
    ct.Selected = v
  }
}

// Assumes we are in normal mode
func (b *CodeBuffer) enforceValidPoint (direction rune) {
  t := b.tree
  if t.GetActive().Data.(*Token).IsClosed() {
    i, err := t.GetActiveIndex()
    panicIfError(err)
    if i > 0 {
      switch direction {
      case 'l': t.Left()
      case 'r': t.Right()
      default: panic("Not found")
      }
      t.Right()
    } else {
      _ = t.Up()
    }
  }
}

func (b *CodeBuffer) moveRight () {
  b.tree.Right()
  b.enforceValidPoint('r')
}

func (b *CodeBuffer) moveLeft () {
  b.tree.Left()
  b.enforceValidPoint('l')
}

func (b *CodeBuffer) moveDown () {
  if len(b.tree.GetActive().Children) > 1 {
    _ = b.tree.DownFirst()
  }
}

func (b *CodeBuffer) deleteNode () {
  i, err := b.tree.DeleteActive()
  panicIfError(err)
  err = b.tree.Down(i)
  panicIfError(err)
  b.enforceValidPoint('l')
}

func (b *CodeBuffer) modeInsert () {
  a := b.tree.GetActive()
  token := NewToken("cursor", " ")
  if a == b.tree.Root {
    err := b.tree.InsertChild(token, 0)
    panicIfError(err)
    err = b.tree.DownFirst()
    panicIfError(err)
  } else {
    err := b.tree.InsertSibling(token, -1)
    panicIfError(err)
    b.tree.Left()
  }
}

func (b *CodeBuffer) modeNotInsert () {
  b.deleteNode()
}

func (b *CodeBuffer) SwapUp() error {
  move := []rune{'u', 'r'}
  err := b.tree.Swap(move)
  return err
}

// Fails if cursor is at end of group, but that seems OK
func (b *CodeBuffer) SwapDown() error {
  move := []rune{'r', 'd'}
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) SwapLeft() error {
  move := []rune{'l'}
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) SwapRight() error {
  move := []rune{'r', 'r'}
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) AppendToToken(s string) {
  t := b.tree.GetActive().Data.(*Token)
  if t.Value == " " {
    t.Value = ""
  }
  t.Value += s
}

func (b *CodeBuffer) BackspaceToken() {
  t := b.tree.GetActive().Data.(*Token)
  l := len(t.Value)
  if l==0 {
    panic("Something is wrong")
  } else if l == 1 {
    t.Value = " "
  } else {
    t.Value = t.Value[:l-1]
  }
}

func (b *CodeBuffer) AppendToken() {
  a := b.tree.GetActive()
  t := a.Data.(*Token)
  t.Class = "symbol"
  err := b.tree.InsertSibling(t, -1)
  panicIfError(err)
  a.Data = NewToken("cursor", " ")
}

func (b *CodeBuffer) AppendOpen(what string) {
  // Construct tokens
  var openToken *Token
  var closeToken *Token
  switch what {
  case "call":
    openToken = NewToken("open", "(")
    closeToken = NewToken("close", ")")
  case "vect":
    openToken = NewToken("open", "[")
    closeToken = NewToken("close", "]")
  }
  cursorToken := NewToken("cursor", " ")
  // Insert tokens
  t := b.tree
  i, err := t.DeleteActive() // Moves active up
  panicIfError(err)
  _ = t.InsertChild(openToken, i)
  _ = t.Down(i)
  t.AppendChild(cursorToken)
  t.AppendChild(closeToken)
  _ = t.DownFirst()
}

func (b *CodeBuffer) toggleStyleAtPoint () {
  p, err := b.tree.GetActiveParent()
  panicIfError(err)
  t := p.Data.(*Token)
  if t.Style == "" {
    t.Style = "alt"
  } else {
    t.Style = ""
  }
}

func (b *CodeBuffer) handle (event []string) {
  var err error
  b.setCursor(false)
  switch event[0] {
  case "move":
    switch event[1] {
    case "left":  b.moveLeft()
    case "right": b.moveRight()
    case "up":    err = b.tree.Up()
    case "down":  b.moveDown()
    }
  case "swap":
    switch event[1] {
      case "left": err = b.SwapLeft()
      case "right": err = b.SwapRight()
      case "up": err = b.SwapUp()
      case "down": err = b.SwapDown()
    }
  case "append":
    switch event[1] {
    case "string": b.AppendToToken(event[2])
    case "backspace": b.BackspaceToken()
    case "token": b.AppendToken()
    case "open": b.AppendOpen(event[2])
    }
  case "insert":
    switch event[1] {
    //case "call":   b.insertCall(event[2])
    //case "vect":   b.insertVect(event[2])
    //case "map":
    //case "symbol": b.insertSymbol(event[2], event[3])
    }
  case "delete": b.deleteNode()
  case "set-mode":
    switch event[1] {
    case "insert": b.modeInsert()
    case "not-insert": b.modeNotInsert()
    //case "normal": b.modeNormal()
    }
  case "toggle-style": b.toggleStyleAtPoint()
  }
  panicIfError(err)
  b.setCursor(true)
  //mapSyntaxTree(b.tree)
}
