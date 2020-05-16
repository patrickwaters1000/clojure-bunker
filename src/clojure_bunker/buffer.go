package main

import (
  "fmt"
  //"errors"
  termbox "github.com/nsf/termbox-go"
)

type Buffer struct {
  name string
  mode string
  tree *Tree
}

func NewBuffer (name string) *Buffer {
  token := NewToken("root", "")
  tree := NewTree(token)
  return &Buffer{name, "normal", tree}
}

func logTree (t *Tree) {
  msg := ""
  traverseFn := func (n *TreeNode) {
    d := n.Data.(*Token)
    msg += fmt.Sprintf("class:%s value:%s children:%d row:%d col:%d\n",
      d.Class, d.Value, len(n.Children), d.Row, d.Col)
  }
  t.DepthFirstTraverse(traverseFn)
  log(msg)
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

// Delete this because tree nodes know their Index values
func getPosition(n *TreeNode) int {
  for i, m := range n.Parent.Children {
    if m == n {
      return i
    }
  }
  return 0 // shouldn't happen
}

func getColor (n *TreeNode) termbox.Attribute {
  t := n.Data.(*Token)
  position := n.Index
  if t.Class == "open" || t.Class == "close" {
    return symbolColor
  } else if position == 0 {
    if isBuiltIn(t.Value) {
      return builtInColor
    } else {
      return fnNameColor
    }
  } else if position == 1 {
    lt := n.Parent.Children[0].Data.(*Token)
    switch lt.Value {
    case "defn": return fnNameColor
    case "def": return varNameColor
    default: return symbolColor
    }
  } else {
    return symbolColor
  }
}

func (b Buffer) render (w Window) {
  w.Print(0, 0, fg1, bg1, b.name)
  logTree(b.tree)
  traverseFn := func (node *TreeNode) {
    var bg, fg termbox.Attribute
    if node.Data.(*Token).Selected {
      fg = fgh
      bg = bgh
    } else {
      fg = getColor(node)
      bg = bg1
    }
    token := node.Data.(*Token)
    w.Print(token.Row + 2, token.Col, fg, bg, token.Value)
  }
  b.tree.DepthFirstTraverseNoRoot(traverseFn)
}

func (b *Buffer) insertCall (position string) {
  class := b.tree.Active.Data.(*Token).Class
  allowed := or(
    position != "below",
    class == "open" || class == "root")
  if allowed {
    var err error
    switch position {
    case "below":
      b.tree.InsertChild(NewToken("open", "("), 0)
      err = b.tree.DownFirst()
    case "before":
      b.tree.InsertSibling(NewToken("open", "("), -1)
      err = b.tree.Left()
    case "after":
      b.tree.InsertSibling(NewToken("open", "("), 0)
      err = b.tree.Right()
    }
    panicIfError(err)
    b.tree.AppendChild(NewToken("close", ")"))
    panicIfError(err)
  }
}

func (b *Buffer) insertVect (position string) {
  class := b.tree.Active.Data.(*Token).Class
  allowed := or(
    position != "below",
    class == "open" || class == "root")
  if allowed {
    var err error
    switch position {
    case "below":
      b.tree.InsertChild(NewToken("open", "["), 0)
      err = b.tree.DownFirst()
    case "before":
      b.tree.InsertSibling(NewToken("open", "["), -1)
      err = b.tree.Left()
    case "after":
      b.tree.InsertSibling(NewToken("open", "["), 0)
      err = b.tree.Right()
    }
    panicIfError(err)
    b.tree.AppendChild(NewToken("close", "]"))
    panicIfError(err)
  }
}

func (b *Buffer) insertSymbol (position, symbol string) {
  class := b.tree.Active.Data.(*Token).Class
  allowed := or(
    position != "below",
    class == "open" || class == "root")
  if allowed {
    token := NewToken("symbol", symbol)
    var err error
    switch position {
    case "below":
      b.tree.InsertChild(token, 0)
      err = b.tree.DownFirst()
    case "before":
      b.tree.InsertSibling(token, -1)
      err = b.tree.Left()
    case "after":
      b.tree.InsertSibling(token, 0)
      err = b.tree.Right()
    }
    panicIfError(err)
  }
}

func (b *Buffer) setCursor (v bool) {
  a := b.tree.Active
  at := a.Data.(*Token)
  at.Selected = v
  if a.Data.(*Token).Class == "open" {
    c := a.Children[len(a.Children) - 1]
    ct := c.Data.(*Token)
    ct.Selected = v
  }
}

func (b *Buffer) moveRight () error {
  err := b.tree.Right()
  token := b.tree.Active.Data.(*Token)
  if err == nil && token.Class == "close" {
    err = b.tree.Right()
  }
  return err
}

func (b *Buffer) moveLeft () error {
  err := b.tree.Left()
  token := b.tree.Active.Data.(*Token)
  if err == nil && token.Class == "close" {
    err = b.tree.Left()
  }
  return err
}

func (b *Buffer) moveDown () error {
  var err error
  if len(b.tree.Active.Children) > 1 {
    err = b.tree.DownFirst()
  }
  return err
}

func (b *Buffer) deleteNode () {
  active := b.tree.Active
  class := active.Data.(*Token).Class
  if class != "root" {
    idx := active.Index
    err := b.tree.Up()
    panicIfError(err)
    err = b.tree.DeleteChild(idx)
    panicIfError(err)
  }
}

func (b *Buffer) modeInsert () {
  token := NewToken("cursor", " ")
  err := b.tree.InsertSibling(token, 0)
  panicIfError(err)
  err = b.tree.Right()
  panicIfError(err)
}

func (b *Buffer) SwapUp() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    return n.Up()
  }
  err := b.tree.Swap(move)
  return err
}

func (b *Buffer) SwapDown() error {
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

func (b *Buffer) SwapLeft() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    n2, err := n.LeftNotLast()
    return n2, err
  }
  err := b.tree.Swap(move)
  return err
}

func (b *Buffer) SwapRight() error {
  move := func (n *TreeNode) (*TreeNode, error) {
    n2, err := n.RightNotLast()
    return n2, err
  }
  err := b.tree.Swap(move)
  return err
}

func (b *Buffer) handleEvent (event []string) error {
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
  case "insert":
    switch event[1] {
    case "call":   b.insertCall(event[2])
    case "vect":   b.insertVect(event[2])
    //case "map":
    case "symbol": b.insertSymbol(event[2], event[3])
    }
  case "delete": b.deleteNode()
  case "mode":
    switch event[1] {
    case "insert": b.modeInsert()
    //case "normal": b.modeNormal()
    }
  }
  panicIfError(err)
  b.setCursor(true)
  mapSyntaxTree(b.tree)
  return nil
}
