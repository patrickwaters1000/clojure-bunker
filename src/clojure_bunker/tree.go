package main

import (
  "errors"
)

// By convention, methods that modify the tree structure are factored
// through TreeNode methods.  

type Tree struct {
  Root *TreeNode
  Active *TreeNode
}

func NewTree (d interface{}) *Tree {
  var root = NewTreeNode(d)
  return &Tree{root, root}
}

func (t *Tree) AppendChild (d interface{}) {
  a := t.Active
  if a == nil {
    panic("Can't append child to active node `nil`")
  }
  c := NewTreeNode(d)
  a.AppendChild(c)
}

func (t *Tree) InsertChild (d interface{}, i int) error {
  a := t.Active
  if a == nil {
    return errors.New("Active node is `nil`")
  }
  err := a.InsertChild(d, i)
  return err
}

// Deprecated because it returns an error
// Inserts sibling in position relative to active node
func (t *Tree) InsertSibling (d interface{}, i int) error {
  a := t.Active
  if a == nil {
    return errors.New("Active node can't be `nil`")
  }
  p := a.Parent
  if p == nil {
    return errors.New("Active node has no parent")
  }
  idx := a.GetIndex()
  err := p.InsertChild(d, idx + i + 1)
  return err
}

func (t *Tree) DeleteChild (i int) error {
  return t.Active.DeleteChild(i)
}

func (t *Tree) DownFirst () error {
  cs := t.Active.Children
  if len(cs) == 0 {
    return errors.New("Active node has no children")
  }
  t.Active = cs[0]
  return nil
}

func (t *Tree) DownLast () error {
  cs := t.Active.Children
  if len(cs) == 0 {
    return errors.New("Active node has no children")
  }
  t.Active = cs[len(cs)-1]
  return nil
}

func (t *Tree) Up () error {
  if t.Active == t.Root {
    return errors.New("Root node has no parent")
  }
  t.Active = t.Active.Parent
  return nil
}

func (t *Tree) Cycle(d int) error {
  a := t.Active
  if a == t.Root {
    return errors.New("Root node has no siblings")
  }
  nCs := len(a.Parent.Children)
  i := mod(a.GetIndex() + d, nCs)
  t.Active = a.Parent.Children[i]
  return nil
}

func (t *Tree) Left () error {
  return t.Cycle(-1)
}

func (t *Tree) Right () error {
  return t.Cycle(1)
}

func (t *Tree) Swap (move func(*TreeNode) (*TreeNode, error)) error {
  n1 := t.Active
  n2, err := move(n1)
  if err != nil { return err }
  if n1.Parent == nil {
      return errors.New("Current node has no parent")
  } else if n2.Parent == nil {
    return errors.New("Target node has no parent")
  }
  p1 := n1.Parent
  p2 := n2.Parent
  i1 := n1.GetIndex()
  i2 := n2.GetIndex()
  err = p1.DeleteChild(i1)
  if err != nil { return err }
  err = p2.InsertChild(n1.Data, i2)
  if err != nil { return err }
  t.Active = p2.Children[i2]
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
