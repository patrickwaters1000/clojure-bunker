package clj_utils

import (
  u "utils"
  "fmt"
)

type CljWriter struct {
  Buffer string
}

func (w *CljWriter) write (s string) {
  w.Buffer += s
}

func (w *CljWriter) jump (row int, col int) {
  w.Buffer += fmt.Sprintf("\x1b[%d;%dH", row, col)
}

func (w *CljWriter) newLine (col int) {
  w.Buffer += "\r\n"
  for i:=0; i<col; i++ {
    w.Buffer += " "
  }
}

func and (conds ...bool) bool {
  for _, cond := range conds {
    if !cond { return false }
  }
  return true
}

func or (conds ...bool) bool {
  for _, cond := range conds {
    if cond { return true }
  }
  return false
}

// Returns whether a space should be rendered before the node's token.
func spaceRequired (n *u.TreeNode) bool {
  token := n.Data.(*Token)
  return n.Index > 0 && !token.Closed
}

// Returns the token of a node's oldest sibling
func getLeader (n *u.TreeNode) string {
  if n.Parent == nil {
    panic("Node has no parent")
  }
  return n.Parent.Children[0].Data.(*Token).Value
}

func newLineRequiredForDefn (n *u.TreeNode) bool {
  if n.Parent == nil {
    return false
  }
  leader := getLeader(n)
  needNewLine := and(
    leader == "defn",
    n.Index > 2,
    !n.Data.(*Token).Closed,
  )
  return needNewLine
}

func newLineRequiredForLet (n *u.TreeNode) bool {
  if n.Parent == nil {
    return false
  }
  leader := getLeader(n)
  return and(
    leader == "let",
    n.Index > 1,
    !n.Data.(*Token).Closed,
  )
}

func newLineRequiredForLetBinding (n *u.TreeNode) bool {
  if n.Parent == nil || n.Parent.Parent == nil {
    return false
  }
  parent := n.Parent.Data.(*Token).Value
  parentLeader := getLeader(n.Parent)
  return and(
    parentLeader == "let",
    parent == "[",
    n.Index != 0,
    n.Index % 2 == 0,
    !n.Data.(*Token).Closed,
  )
}

// Returns whether a newline should be rendered before the node's token.
func newLineRequired (n *u.TreeNode) bool {
  return or(
    newLineRequiredForDefn(n),
    newLineRequiredForLet(n),
    newLineRequiredForLetBinding(n),
  )
}

//func getPreviousToken(n *u.TreeNode) *Token {
//  if n.Parent == nil {
//    panic("No previous token")
//  } else if n.Index == 0 {
//    return n.Parent.Data.(*Token)
//  } else {
//    return n.Parent.Children[n.Index - 1].Data.(*Token)
//  }
//}

func getRow(n *u.TreeNode, previousToken *Token) int {
  if n.Parent == nil {
    return 0
  }
  previousRow := previousToken.Row
  if previousRow == -1 {
    panic("Previous token's row hasn't been set")
  }
  if newLineRequired(n) {
    return previousRow + 1
  } else {
    return previousRow
  }
}

func getCol(n *u.TreeNode, previousToken *Token) int {
  if n.Parent == nil {
    return 0
  }
  if newLineRequired(n) {
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

func UnParseClj(tree *u.Tree) string {
  w := &CljWriter{""}
  var previousToken *Token
  traverseFn := func (n *u.TreeNode) {
    token := n.Data.(*Token)
    if (token.Class=="symbol" && token.Value=="root") {
      token.Row = 0
      token.Col = -4
    } else {
      token.Row = getRow(n, previousToken)
      token.Col = getCol(n, previousToken)
      if spaceRequired(n) {
        w.write(" ")
      }
      if newLineRequired(n) {
        w.newLine(token.Col)
      }
      w.write(token.Value)
    }
    previousToken = token
  }
  tree.DepthFirstTraverse(traverseFn)
  return w.Buffer
}


