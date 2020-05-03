package main

import (
  "fmt"
  termbox "github.com/nsf/termbox-go"
)

type Buffer struct {
  name string
  tree *Tree
}

func NewBuffer (name string) *Buffer {
  token := NewToken("root", "")
  tree := NewTree(token)
  return &Buffer{name, tree}
}

func logTree (t *Tree) {
  msg := ""
  traverseFn := func (n *TreeNode) {
    d := n.Data.(*Token)
    msg += fmt.Sprintf("%s %s %d\n", d.Class, d.Value, len(n.Children))
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
  position := getPosition(n)
  if t.Class == "open" || t.Class == "closed" {
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

func (b Buffer) render () {
  tbPrint(0, 0, fg1, bg1, b.name)
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
    tbPrint(token.Row + 2, token.Col, fg, bg, token.Value)
  }
  b.tree.DepthFirstTraverseNoRoot(traverseFn)
}

func (b *Buffer) insertCall () {
  b.tree.AppendChild(NewToken("open", "("))
  err := b.tree.DownFirst()
  panicIfError(err)
  b.tree.AppendChild(NewToken("close", ")"))
  err = b.tree.DownFirst()
  panicIfError(err)
}

func (b *Buffer) insertVect () {
  b.tree.InsertSibling(NewToken("open", "["), -1)
  err := b.tree.Left()
  panicIfError(err)
  b.tree.AppendChild(NewToken("close", "]"))
  err = b.tree.DownFirst()
  panicIfError(err)
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
  if err != nil && token.Class == "closed" {
    err = b.tree.Right()
  }
  return err
}

func (b *Buffer) moveLeft () error {
  err := b.tree.Left()
  token := b.tree.Active.Data.(*Token)
  if err != nil && token.Class == "closed" {
    err = b.tree.Left()
  }
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
    case "down":  err = b.tree.DownFirst()
    }
  case "insert":
    switch event[1] {
    case "call":   b.insertCall()
    case "vect":   b.insertVect()
    //case "map":
    case "symbol":
      token := NewToken("symbol", event[2])
      b.tree.InsertSibling(token, -1)
    }
  case "delete":
    switch event[1] {
    case "symbol":
    case "group":
    }
  }
  panicIfError(err)
  b.setCursor(true)
  mapSyntaxTree(b.tree)
  return nil
}
