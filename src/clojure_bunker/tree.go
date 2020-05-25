package main

import (
  "errors"
)
// NOTE Ideally only Tree methods would refer to tree nodes
//      That is why (for example) the insert methods take a data object,
//      not a *TreeNode. The exception is that when using the traverse
//      methods, it is necessary to define a traverse function that
//      knows about TreeNode's.

// TODO Replace TreeNode.GetIndex()
// TODO Replace TreeNode.GetSiblings()
// TODO Replace TreeNode.Up()


// By convention, methods that modify the tree structure are factored
// through TreeNode methods.  

type Tree struct {
  Root *TreeNode
  Path []*TreeNode // Path from `Root` to a distinguished "active" node
}

func NewTree (d interface{}) *Tree {
  var root = NewTreeNode(d)
  return &Tree{
    Root: root,
    Path: []*TreeNode{root},
  }
}

// NOTE The following cannot fail because `t.Path` always contains `t.Root`.
func (t Tree) GetActive () *TreeNode {
  return t.Path[len(t.Path) - 1]
}

func (t Tree) GetActiveParent () (*TreeNode, error) {
  l := len(t.Path)
  if l == 1 {
    return nil, errors.New("Can't get parent of root node")
  }
  return t.Path[len(t.Path) - 2], nil
}

func (t Tree) GetActiveIndex () (int, error) {
  a := t.GetActive()
  p, err := t.GetActiveParent()
  if err != nil {
    return 0, err
  }
  for i, c := range p.Children {
    if c == a {
      return i, nil
    }
  }
  panic("Not found")
}

func (t *Tree) AppendChild (d interface{}) {
  t.GetActive().AppendChild(d)
}

func (t *Tree) InsertChild (d interface{}, i int) error {
  a := t.GetActive()
  err := a.InsertChild(d, i)
  return err
}

// NOTE: I have wanted to get rid of this function on multiple
// occasions, and repeatedly concluded that's a bad idea!!!
// Inserts sibling in position relative to active node
func (t *Tree) InsertSibling (d interface{}, i int) error {
  p, err := t.GetActiveParent()
  if p == nil {
    return err
  }
  idx, _ := t.GetActiveIndex()
  err = p.InsertChild(d, idx + i + 1)
  return err
}

func (t *Tree) DeleteChild (i int) error {
  return t.GetActive().DeleteChild(i)
}

func (t *Tree) DeleteActive () (int, error) {
  i, err := t.GetActiveIndex()
  if err != nil {
    return 0, errors.New("Can't delete root node")
  }
  _ = t.Up()
  return i, t.DeleteChild(i)
}

// Returns a copy of `t` with a mutation applied to the active node
// Just enough of `t` is deeply copied that a sequence of `UpdateAvtiveNode`
// operations behave as if we were truly using deep copies.
func (t *Tree) UpdateActiveNode (updateFn func(*TreeNode)) *Tree {
  tNew := NewTree(t.Root.Data)
  l := len(t.Path)
  aOld := t.GetActive()
  for i, n := range t.Path[:l - 1] {
    cPath := t.Path[i + 1]
    var jCopy int
    for j, c := range n.Children {
      if c == cPath {
        tNew.AppendChild(c.Data)
        jCopy = j
        if c == aOld {
          updateFn(tNew.GetActive().Children[j])
        }
      } else {
        aNew := tNew.GetActive()
        aNew.Children = append(aNew.Children, c)
      }
    }
    newPathNode := tNew.GetActive().Children[jCopy]
    tNew.Path = append(tNew.Path, newPathNode)
  }
  return tNew
}

func (t *Tree) Down (i int) error {
  cs := t.GetActive().Children
  if len(cs) <= i {
    return errors.New("Not enough children")
  }
  t.Path = append(t.Path, cs[i])
  return nil
}

func (t *Tree) DownFirst () error {
  return t.Down(0)
}

func (t *Tree) DownLast () error {
  nCs := len(t.GetActive().Children)
  return t.Down(nCs - 1)
}

func (t *Tree) Up () error {
  if len(t.Path) == 1 {
    return errors.New("Can't get parent of root node")
  }
  t.Path = t.Path[:len(t.Path) - 1]
  return nil
}

func (t *Tree) Cycle(d int) {
  if len(t.Path) == 1 {
    return // `GetActiveIndex` fails for root node,
    // but in this case doing nothing is a reasonable behaviour
  }
  i1, _ := t.GetActiveIndex()
  p, _ := t.GetActiveParent()
  nCs := len(p.Children)
  i2 := mod(i1 + d, nCs)
  l := len(t.Path)
  t.Path[l - 1] = p.Children[i2]
}

func (t *Tree) Left () {
  t.Cycle(-1)
}

func (t *Tree) Right () {
  t.Cycle(1)
}

func (t *Tree) Move (move []rune) error {
  for _, d := range move {
    switch d {
    case 'l': t.Left()
    case 'r': t.Right()
    case 'd':
      err := t.DownFirst()
      if err != nil {
        return err
      }
    case 'D':
      err := t.DownLast()
      if err != nil {
        return err
      }
    case 'u':
      err := t.Up()
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (t *Tree) Swap (move []rune) error {

  // Capture info about active node
  d := t.GetActive().Data
  i, err := t.GetActiveIndex()
  p, _ := t.GetActiveParent()
  if err != nil { return err }
  // Update path
  err = t.Move(move)
  if err != nil { return err }
  // Delete old active node
  _ = p.DeleteChild(i)
  // Re-insert old active node in new location
  err = t.InsertSibling(d, -1)
  if err != nil { return err }
  // Make re-inserted node active
  t.Left()
  return nil
}

func (t *Tree) DepthFirstTraverse (f func (*TreeNode)) {
  var s = NewStack()
  var n *TreeNode
  s.Push(t.Root)
  for !s.Empty() {
    n = (s.Pop()).(*TreeNode)
    f(n)
    l := len(n.Children)
    for i:=l-1; i>=0; i-- {
      s.Push(n.Children[i])
    }
  }
}

func (t *Tree) DepthFirstTraverseNoRoot (f func (*TreeNode)) {
  var s = NewStack()
  var n *TreeNode
  l := len(t.Root.Children)
  for i:=l-1; i>=0; i-- {
    s.Push(t.Root.Children[i])
  }
  for !s.Empty() {
    n = (s.Pop()).(*TreeNode)
    f(n)
    l := len(n.Children)
    for i:=l-1; i>=0; i-- {
      s.Push(n.Children[i])
    }
  }
}
