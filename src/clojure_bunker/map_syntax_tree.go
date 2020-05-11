package main

// Returns whether a space should be rendered before the node's token.
func spaceRequired (n *TreeNode) bool {
  token := n.Data.(*Token)
  return n.Index > 0 && !token.IsClosed()
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
    n.Index > 2,
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
    n.Index > 1,
    !n.Data.(*Token).IsClosed(),
  )
}

func newLineRequiredForLetBinding (n *TreeNode) bool {
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
    n.Index > 0)
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
