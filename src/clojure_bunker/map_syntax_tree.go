package main

import (
  termbox "github.com/nsf/termbox-go"
)

func getColor (what string) termbox.Attribute {
  switch what {
  case "built-in"  : return termbox.ColorYellow
  case "func-name" : return termbox.ColorCyan
  case "comment"   : return termbox.ColorBlue
  case "var-name"  : return termbox.ColorMagenta
  case "keyword"   : return termbox.ColorGreen
  case "string"    : return termbox.ColorRed
  case "symbol"    : return termbox.ColorWhite
  case "paren"     : return termbox.ColorWhite
  case "normal"    : return termbox.ColorWhite
  default          : panic("Not found")
  }
}

type Artist struct {
  window Window
  row int
  col int
}

func NewArtist (w Window) *Artist {
  return &Artist{w, 0, 0}
}

func (a Artist) Copy () *Artist {
  return &Artist{a.window, a.row, a.col}
}

func (a *Artist) Print (n *TreeNode, what string) {
  t := n.Data.(*Token)
  bg := termbox.ColorBlack
  fg := getColor(what)
  if t.Selected {
    bg = fg
    fg = termbox.ColorBlack
  }
  msg := n.Data.(*Token).Value
  a.window.Print(a.row, a.col, fg, bg, msg)
  a.col += len(msg)
}

func (a *Artist) Draw (n *TreeNode, role string) {
  t := n.Data.(*Token)
  t.Row = a.row
  t.Col = a.col
  t.Color = getColor(role)
  a.col += len(t.Value)
}

func (a *Artist) render (n *TreeNode, role string) {
  a.Draw(n, role)
  t := n.Data.(*Token)
  if len(n.Children) == 0 {
    // pass
  } else if t.Value == "[" {
    switch t.Style {
    case "alt": a.Vect(n)
    default: a.FlatCall(n)
    }
  } else {
    leader := n.Children[0].Data.(*Token).Value
    switch leader {
    case "let": a.Let(n)
    //case "def": a.Def(n)
    case "defn": a.Defn(n)
    default:
      switch t.Style {
      case "alt": a.FlatCall(n)
      default: a.Call(n)
      }
    }
  }
}

func (a *Artist) Root (n *TreeNode) {
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children[:nCs - 1] {
    switch i {
    case 0:
    default: a.row = b.row + 2
    }
    b = a.Copy()
    b.render(c, "normal")
  }
}

func (a *Artist) Call (n *TreeNode) {
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    switch i {
    case 0:
      a.row = b.row
      a.col = b.col
      role = "func-name"
    case nCs - 1:
      a.row = b.row
      a.col = b.col
      role = "paren"
    case 1:
      a.row = b.row
      a.col = b.col + 1
    default:
      a.row = b.row + 1
    }
    b = a.Copy()
    b.render(c, role)
  }
  a.row = b.row
  a.col = b.col
}

func (a *Artist) FlatCall (n *TreeNode) {
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    switch i {
    case 0:
      a.row = b.row
      a.col = b.col
      role = "func-name"
    case nCs - 1:
      a.row = b.row
      a.col = b.col
      role = "paren"
    default:
      a.row = b.row
      a.col = b.col + 1
    }
    b = a.Copy()
    b.render(c, role)
  }
  a.row = b.row
  a.col = b.col
}

func (a *Artist) Vect (n *TreeNode) {
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    switch i {
    case 0:
      a.row = b.row
      a.col = b.col
    case nCs - 1:
      a.row = b.row
      a.col = b.col
      role = "paren"
    default:
      a.row = b.row + 1
    }
    b = a.Copy()
    b.render(c, role)
  }
  a.row = b.row
  a.col = b.col
}

func (a *Artist) Binding (n *TreeNode) {
  a.Draw(n, "normal")
  b := a.Copy()
  var c0 int
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else if i == 0 {
      a.row = b.row
      a.col = b.col
      c0 = b.col
    } else if mod(i, 2) == 0 {
      a.row = b.row + 1
      a.col = c0
    } else {
      a.row = b.row
      a.col = b.col + 1
    }
    b = a.Copy()
    b.render(c, role)
  }
  a.row = b.row
  a.col = b.col
}

func (a *Artist) Let (n *TreeNode) {
  c0 := a.col
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == 0 {
      a.col = b.col
      role = "built-in"
    } else if i == 1 {
      a.row = b.row
      a.col = b.col + 1
    } else if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else {
      a.row = b.row + 1
      a.col = c0 + 2
    }
    b = a.Copy()
    if i == 1 {
      b.Binding(c)
    } else {
      b.render(c, role)
    }
  }
  a.row = b.row
  a.col = b.col
}

func (a *Artist) Defn (n *TreeNode) {
  c0 := a.col
  b := a.Copy()
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == 0 {
      a.row = b.row
      a.col = b.col
      role = "built-in"
    } else if i == 1 || i == 2 {
      a.row = b.row
      a.col = b.col + 1
    } else if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else {
      a.row = b.row + 1
      a.col = c0 + 1
    }
    b = a.Copy()
    b.render(c, role)
  }
  a.row = b.row
  a.col = b.col
}

// Returns whether a space should be rendered before the node's token.
func spaceRequired (n *TreeNode) bool {
  token := n.Data.(*Token)
  return n.GetIndex() > 0 && !token.IsClosed()
}

// Returns the token of a node's oldest sibling
func getLeader (n *TreeNode) string {
  if n.Parent == nil {
    panic("Node has no parent")
  }
  return n.Parent.Children[0].Data.(*Token).Value
}

func newLineRequiredForDefn (n *TreeNode) bool {
  if n.Parent == nil {
    return false
  }
  leader := getLeader(n)
  needNewLine := and(
    leader == "defn",
    n.GetIndex() > 2,
    !n.Data.(*Token).IsClosed(),
  )
  return needNewLine
}

func newLineRequiredForLet (n *TreeNode) bool {
  if n.Parent == nil {
    return false
  }
  leader := getLeader(n)
  return and(
    leader == "let",
    n.GetIndex() > 1,
    !n.Data.(*Token).IsClosed(),
  )
}

func newLineRequiredForLetBinding (n *TreeNode) bool {
  if n.Parent == nil || n.Parent.Parent == nil {
    return false
  }
  parent := n.Parent.Data.(*Token).Value
  parentLeader := getLeader(n.Parent)
  i := n.GetIndex()
  return and(
    parentLeader == "let",
    parent == "[",
    i != 0,
    i % 2 == 0,
    !n.Data.(*Token).IsClosed(),
  )
}

// Returns whether a newline should be rendered before the node's token.
func newLineRequired (n *TreeNode) bool {
  return or(
    newLineRequiredForDefn(n),
    newLineRequiredForLet(n),
    newLineRequiredForLetBinding(n),
  )
}

func doubleNewLineRequired (n *TreeNode) bool {
  return and(
    n.Parent.Data.(*Token).Class == "root",
    n.GetIndex() > 0)
}

func getRow(n *TreeNode, previousToken *Token) int {
  if previousToken == nil || n.Parent == nil {
    return 0
  }
  previousRow := previousToken.Row
  if previousRow == -1 {
    panic("Previous token's row hasn't been set")
  }
  if doubleNewLineRequired(n) {
    return previousRow + 2
  } else if newLineRequired(n) {
    return previousRow + 1
  } else {
    return previousRow
  }
}

func getCol(n *TreeNode, previousToken *Token) int {
  if previousToken == nil || n.Parent == nil {
    return 0
  } else if n.Parent.Data.(*Token).Class == "root" {
    return 0
  } else if newLineRequired(n) || doubleNewLineRequired(n) {
    parentToken := n.Parent.Data.(*Token)
    parentCol := parentToken.Col
    if parentCol == -1 {
      panic("Parent token col hasn't been set")
    }
    if parentToken.Value == "(" {
      return parentCol + 2
    } else {
      return parentCol + 1
    }
  } else { // no new line
    if previousToken.Col == -1 {
      panic("Previous token's col hasn't been set")
    }
    offset := len(previousToken.Value)
    if spaceRequired(n) {
      offset += 1
    }
    return previousToken.Col + offset
  }
}

func mapSyntaxTree (tree *Tree) {
  var previousToken *Token
  traverseFn := func (node *TreeNode) {
    token := node.Data.(*Token)
    token.Row = getRow(node, previousToken)
    token.Col = getCol(node, previousToken)
    previousToken = token
  }
  tree.DepthFirstTraverseNoRoot(traverseFn)
}
