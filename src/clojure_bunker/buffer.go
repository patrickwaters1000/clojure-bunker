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
    active := n == t.Active
    consistent := n.checkChildParentConsistency()
    msg += fmt.Sprintf(
      "class:%s value:%s children:%d consistent:%v selected:%v active:%v row:%d col:%d\n",
      d.Class, d.Value, len(n.Children), consistent, d.Selected, active, d.Row, d.Col)
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



// func getColor (n *TreeNode) termbox.Attribute {
//   t := n.Data.(*Token)
//   position := n.GetIndex()
//   if t.Class == "open" || t.Class == "close" {
//     return symbolColor
//   } else if position == 0 {
//     if isBuiltIn(t.Value) {
//       return builtInColor
//     } else {
//       return fnNameColor
//     }
//   } else if position == 1 {
//     lt := n.Parent.Children[0].Data.(*Token)
//     switch lt.Value {
//     case "defn": return fnNameColor
//     case "def": return varNameColor
//     default: return symbolColor
//     }
//   } else {
//     return symbolColor
//   }
// }

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

func (b *CodeBuffer) setCursor (v bool) {
  a := b.tree.Active
  at := a.Data.(*Token)
  at.Selected = v
  if a.Data.(*Token).Class == "open" {
    c := a.Children[len(a.Children) - 1]
    ct := c.Data.(*Token)
    ct.Selected = v
  }
}

func (b *CodeBuffer) moveRight () error {
  err := b.tree.Right()
  token := b.tree.Active.Data.(*Token)
  if err == nil && token.Class == "close" {
    err = b.tree.Right()
  }
  return err
}

func (b *CodeBuffer) moveLeft () error {
  err := b.tree.Left()
  token := b.tree.Active.Data.(*Token)
  if err == nil && token.Class == "close" {
    err = b.tree.Left()
  }
  return err
}

func (b *CodeBuffer) moveDown () error {
  var err error
  if len(b.tree.Active.Children) > 1 {
    err = b.tree.DownFirst()
  }
  return err
}

func (b *CodeBuffer) deleteNode () {
  active := b.tree.Active
  class := active.Data.(*Token).Class
  if class != "root" {
    idx := active.GetIndex()
    err := b.tree.Up()
    panicIfError(err)
    err = b.tree.DeleteChild(idx)
    panicIfError(err)
  }
}

func (b *CodeBuffer) modeInsert () {
  a := b.tree.Active
  token := NewToken("cursor", " ")
  if a == b.tree.Root {
    err := b.tree.InsertChild(token, 0)
    panicIfError(err)
    err = b.tree.DownFirst()
    panicIfError(err)
  } else {
    err := b.tree.InsertSibling(token, -1)
    panicIfError(err)
    err = b.tree.Left()
    panicIfError(err)
  }
}

func (b *CodeBuffer) modeNotInsert () {
  a := b.tree.Active
  cs := a.GetSiblings()
  if len(cs) == 1 { panic("Something is wrong") }
  i := a.GetIndex()
  p := a.Parent
  if p == nil { panic("Something is wrong") }
  err := p.DeleteChild(i)
  panicIfError(err)
  if len(cs) == 2 { // Active node and close token
    b.tree.Active = p
  } else {
    b.tree.Active = p.Children[i]
  }
}

func (b *CodeBuffer) SwapUp() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    return n.Up()
  }
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) SwapDown() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    r, err := n.RightNotLast()
    if err != nil { return nil, err }
    fc, err := r.DownFirst()
    if err != nil { return nil, err }
    return fc, nil
  }
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) SwapLeft() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    n2, err := n.LeftNotLast()
    return n2, err
  }
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) SwapRight() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    n2, err := n.RightNotLast()
    return n2, err
  }
  err := b.tree.Swap(move)
  return err
}

func (b *CodeBuffer) AppendToToken(s string) {
  t := b.tree.Active.Data.(*Token)
  if t.Value == " " {
    t.Value = ""
  }
  t.Value += s
}

func (b *CodeBuffer) BackspaceToken() {
  t := b.tree.Active.Data.(*Token)
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
  a := b.tree.Active
  t := a.Data.(*Token)
  t.Class = "symbol"
  err := b.tree.InsertSibling(t, -1)
  panicIfError(err)
  a.Data = NewToken("cursor", " ")
}

func (b *CodeBuffer) AppendOpen(what string) {
  a := b.tree.Active
  // p, err := a.Up()
  // panicIfError(err)
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
  cursorNode := NewTreeNode(cursorToken)
  closeNode := NewTreeNode(closeToken)
  a.Data = openToken
  a.AppendChild(cursorNode)
  a.AppendChild(closeNode)
  err := b.tree.DownFirst()
  panicIfError(err)
}

func (b *CodeBuffer) toggleStyleAtPoint () {
  t := b.tree.Active.Parent.Data.(*Token)
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
    case "left":  err = b.moveLeft()
    case "right": err = b.moveRight()
    case "up":    err = b.tree.Up()
    case "down":  err = b.moveDown()
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
