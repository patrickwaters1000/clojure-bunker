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
  file string
  tree *Tree // By convention, the root node's value holds the buffer's mode
  history []*Tree
}

func NewCodeBuffer () *CodeBuffer {
  rootToken := NewToken("root", "normal")
  tree := NewTree(rootToken)
  leafToken := NewToken("leaf", "")
  tree.AppendChild(leafToken)
  return &CodeBuffer{
    name: "",
    file: "",
    tree: tree,
    history: []*Tree{},
  }
}

func stringifyTree (t *Tree) string {
  msg := ""
  traverseFn := func (n *TreeNode) {
    d := n.Data.(*Token)
    active := n == t.GetActive()
    msg += fmt.Sprintf(
      "class:%s value:%s children:%d selected:%v active:%v row:%d col:%d style:%s\n",
      d.Class, d.Value, len(n.Children), d.Selected, active, d.Row, d.Col, d.Style)
  }
  t.DepthFirstTraverse(traverseFn)
  return msg
}

func logTree (t *Tree) {
  log(stringifyTree(t))
}

// For debugging, not writing to disc
func (b CodeBuffer) stringify () string {
  mode := b.tree.Root.Data.(*Token).Value
  return fmt.Sprintf("CodeBuffer:\nname: %s\nmode: %s\ntree:\n%s",
    b.name, mode, stringifyTree(b.tree)) + "\n\n"
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

func (b *CodeBuffer) undo () {
  l := len(b.history)
  b.tree = b.history[l - 1]
  b.history = b.history[:l - 1]
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
// Should only be used after another movement operation,
// thus it's unnecessary to use deep copying
func (b *CodeBuffer) enforceValidPoint (direction rune) {
  t := b.tree
  if t.GetActive().Data.(*Token).IsClosed() {
    t.Move([]rune{direction})
  }
}

func (b *CodeBuffer) moveRightNormal () {
  tNew, err := b.tree.PersistentMove([]rune{'r'})
  panicIfError(err)
  b.history = append(b.history, b.tree)
  b.tree = tNew
  b.enforceValidPoint('r')
}

func (b *CodeBuffer) moveLeftNormal () {
  tNew, err := b.tree.PersistentMove([]rune{'l'})
  panicIfError(err)
  b.history = append(b.history, b.tree)
  b.tree = tNew
  b.enforceValidPoint('l')
}

func (b *CodeBuffer) moveDownNormal () {
  if len(b.tree.GetActive().Children) > 1 {
    tNew, err := b.tree.PersistentMove([]rune{'d'})
    panicIfError(err)
    b.history = append(b.history, b.tree)
    b.tree = tNew
  }
}

func (b *CodeBuffer) moveUpNormal () {
  t := b.tree
  if t.GetActive() != t.Root {
    tNew, err := t.PersistentMove([]rune{'u'})
    panicIfError(err)
    b.history = append(b.history, t)
    b.tree = tNew
  }
}

// Actually only mutates active node's parent
// totally OK to copy w.r.t. current path
func (b *CodeBuffer) moveLeftInsert () {
  t := b.tree
  i1, _ := t.GetActiveIndex() // Root cannot be active node in insert mode
  tNew := b.tree.PersistentCopy()
  p, _ := tNew.GetActiveParent()
  tempNode := p.Children[i1]
  i2 := mod(i1 - 1, len(p.Children) - 1)
  p.Children[i1] = p.Children[i2]
  p.Children[i2] = tempNode
  b.history = append(b.history, t)
  b.tree = tNew
}

func (b *CodeBuffer) moveRightInsert () {
  t := b.tree
  i1, _ := t.GetActiveIndex() // Root cannot be active node in insert mode
  tNew := b.tree.PersistentCopy()
  p, _ := tNew.GetActiveParent()
  tempNode := p.Children[i1]
  i2 := mod(i1 + 1, len(p.Children) - 1)
  p.Children[i1] = p.Children[i2]
  p.Children[i2] = tempNode
  b.history = append(b.history, t)
  b.tree = tNew
}

func (b *CodeBuffer) moveUpInsert () {
  t := b.tree
  pOld, err := t.GetActiveParent()
  panicIfError(err)
  if pOld != t.Root {
    tNew := b.tree.PersistentCopy()
    childIdx, _ := tNew.GetActiveIndex()
    _ = tNew.Up()
    p := tNew.GetActive() // Having moved up, this is parent node
    parentIdx, _ := tNew.GetActiveIndex()
    g, _ := tNew.GetActiveParent() // Grandparent
    childToken := p.Children[childIdx].Data
    _ = p.DeleteChild(childIdx)
    _ = g.InsertChild(childToken, parentIdx + 1)
    tNew.Right()
    b.history = append(b.history, t)
    b.tree = tNew
  }
}

func (b *CodeBuffer) moveDownInsert () {
  t := b.tree
  idx, _ := t.GetActiveIndex()
  data := t.GetActive().Data
  t.Right() // A bit sloppy here, after an undo, the cursor would incorrectly move
  if len(t.GetActive().Children) > 0 {
    tNew := b.tree.PersistentCopy()
    p, _ := tNew.GetActiveParent()
    _ = p.DeleteChild(idx)
    a := tNew.GetActive()
    a.InsertChild(data, 0)
    tNew.Path = append(tNew.Path, a.Children[0])
    b.history = append(b.history, t)
    b.tree = tNew
  }
}
func (b *CodeBuffer) deleteNode () {
  tNew, _, err := b.tree.PersistentDeleteActive()
  panicIfError(err)
  b.tree = tNew
}

func (b *CodeBuffer) enterInsertMode () {
  tNew := b.tree.PersistentCopy()
  tNew.Root.Data.(*Token).Value = "insert"
  cursorToken := NewToken("cursor", " ")
  if tNew.GetActive() == tNew.Root {
    _ = tNew.InsertChild(cursorToken, 0)
    _ = tNew.DownFirst()
  } else {
    _ = tNew.InsertSibling(cursorToken, -1)
    tNew.Left()
  }
  b.history = append(b.history, b.tree)
  b.tree = tNew
}

func (b *CodeBuffer) exitInsertMode () {
  tNew, _, err := b.tree.PersistentDeleteActive()
  panicIfError(err)
  tNew.Root.Data.(*Token).Value = "normal"
  b.history = append(b.history, b.tree)
  b.tree = tNew
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
  tNew := b.tree.PersistentCopy()
  a := tNew.GetActive()
  t := a.Data.(*Token)
  switch t.Value[0] {
  case '"': t.Class = "string"
  case ':': t.Class = "keyword"
  default: t.Class = "symbol"
  }
  err := tNew.InsertSibling(t, -1)
  panicIfError(err)
  a.Data = NewToken("cursor", " ")
  b.history = append(b.history, b.tree)
  b.tree = tNew
}

func (b *CodeBuffer) AppendOpen(what string) {
  // Construct tokens
  var openToken *Token
  var closeToken *Token
  switch what {
  case "call":
    openToken = NewToken("open", "(")
    openToken.Style = "long"
    closeToken = NewToken("close", ")")
  case "vect":
    openToken = NewToken("open", "[")
    openToken.Style = "short"
    closeToken = NewToken("close", "]")
  }
  cursorToken := NewToken("cursor", " ")
  // Insert tokens
  tNew := b.tree.PersistentCopy()
  i, err := tNew.DeleteActive() // Moves active up
  panicIfError(err)
  _ = tNew.InsertChild(openToken, i)
  _ = tNew.Down(i)
  tNew.AppendChild(cursorToken)
  tNew.AppendChild(closeToken)
  _ = tNew.DownFirst()
  b.history = append(b.history, b.tree)
  b.tree = tNew
}

func (b *CodeBuffer) toggleStyleAtPoint () {
  tNew := b.tree.PersistentCopy()
  p := tNew.GetActive()
  tok := p.Data.(*Token)
  switch tok.Style {
  case "long":       tok.Style = ""
  case "", "custom": tok.Style = "short"
  case "short":      tok.Style = "binding"
  case "binding":    tok.Style = "long"
  }
  b.history = append(b.history, b.tree)
  b.tree = tNew
}


func (b *CodeBuffer) handle (event []string) {
  b.setCursor(false)
  switch event[0] {
  case "undo": b.undo()
  case "move-normal-left":  b.moveLeftNormal()
  case "move-normal-right": b.moveRightNormal()
  case "move-normal-up":    b.moveUpNormal()
  case "move-normal-down":  b.moveDownNormal()
  case "move-insert-left":  b.moveLeftInsert()
  case "move-insert-right": b.moveRightInsert()
  case "move-insert-up":    b.moveUpInsert()
  case "move-insert-down":  b.moveDownInsert()
  case "insert-string":     b.AppendToToken(event[1])
  case "insert-backspace":  b.BackspaceToken()
  case "insert-token":      b.AppendToken()
  case "insert-call":       b.AppendOpen("call")
  case "insert-vect":       b.AppendOpen("vect")
  case "delete":            b.deleteNode()
  case "enter-insert-mode": b.enterInsertMode()
  case "exit-insert-mode":  b.exitInsertMode()
  case "toggle-style":      b.toggleStyleAtPoint()
  }
  b.setCursor(true)
}
