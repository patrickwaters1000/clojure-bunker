package main

import (
  termbox "github.com/nsf/termbox-go"
)

// TODO: Clean up unused cases
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
  case "open"      : return termbox.ColorWhite
  case "close"     : return termbox.ColorWhite
  case "fail"      : return termbox.ColorWhite
  case "normal"    : return termbox.ColorWhite
  case "cursor"    : return termbox.ColorWhite
  default          : panic("Not found: " + what)
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

func (a *Artist) Draw (n *TreeNode, role string) {
  t := n.Data.(*Token)
  t.Row = a.row
  t.Col = a.col
  if t.Class == "symbol" {
    t.Color = getColor(role)
  } else {
    t.Color = getColor(t.Class)
  }
  a.col += len(t.Value)
}

func (a *Artist) render (n *TreeNode, role string) {
  a.Draw(n, role)
  t := n.Data.(*Token)
  if len(n.Children) == 0 {
    // pass
  } else if t.Value == "[" {
    a.Vect(n)
  } else {
    leader := n.Children[0].Data.(*Token).Value
    switch leader {
    case "let", "def", "defn", "if", "for": a.Let(n)
    default: a.Call(n)
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

func setNewLines (n *TreeNode, newLineFn func(int) bool) {
  for i, c := range n.Children {
    c.Data.(*Token).NewLine = newLineFn(i)
  }
}

func applyStyle (n *TreeNode) {
  switch n.Data.(*Token).Style {
  case "long":
    setNewLines(n, func(i int) bool { return i > 1 })
  case "short":
    setNewLines(n, func(i int) bool { return false })
  case "binding":
    setNewLines(n, func(i int) bool {
      return mod(i, 2) == 0 && i > 0
    })
  }
}

func (a *Artist) Call (n *TreeNode) {
  col0 := a.col + 1
  b := a.Copy() // Note: cannot be moved inside `for`
  applyStyle(n)
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == 0 {
      a.row = b.row
      a.col = b.col
      role = "func-name"
    } else if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else if c.Data.(*Token).NewLine {
      a.col = col0
      a.row = b.row + 1
    } else {
      if i == 1 {
        col0 = b.col + 1
      }
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
  col0 := a.col
  b := a.Copy()
  applyStyle(n)
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == 0 {
      a.row = b.row
      a.col = b.col
    } else if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else if c.Data.(*Token).NewLine {
      a.col = col0
      a.row = b.row + 1
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

// Shamefully, this is a copy of `Vect` with the first line changed.
// TODO: Dry this out!
func (a *Artist) Let (n *TreeNode) {
  col0 := a.col + 1
  b := a.Copy()
  applyStyle(n)
  nCs := len(n.Children)
  for i, c := range n.Children {
    role := "normal"
    if i == 0 {
      a.row = b.row
      a.col = b.col
      role = "built-in"
    } else if i == nCs - 1 {
      a.row = b.row
      a.col = b.col
      role = "paren"
    } else if c.Data.(*Token).NewLine {
      a.col = col0
      a.row = b.row + 1
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

