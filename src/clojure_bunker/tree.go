package main

import (
  "errors"
)

type Tree struct {
  Root *TreeNode
  Active *TreeNode
}

func NewTree (d interface{}) *Tree {
  var root = &TreeNode{d, 0, nil, []*TreeNode{}}
  return &Tree{root, root}
}

func (t *Tree) AppendChild (d interface{}) error {
  a := t.Active
  if a == nil {
    return errors.New("Active node can't be `nil`")
  } else {
    i := len(a.Children)
    c := &TreeNode{d, i, a, []*TreeNode{}}
    a.Children = append(a.Children, c)
    return nil
  }
}

// Should define get fns that do error handling

func (t *Tree) InsertChild (d interface{}, i int) error {
  a := t.Active
  if a == nil {
    return errors.New("Active node is `nil`")
  }
  err := a.InsertChild(d, i)
  return err
}

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
  err := p.InsertChild(d, a.Index + i + 1)
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

func (t *Tree) Left () error {
  var a = t.Active
  if a == t.Root {
    return errors.New("Root node has no siblings")
  }
  var numChildren = len(a.Parent.Children)
  var newIndex = (a.Index + numChildren - 1) % numChildren
  t.Active = t.Active.Parent.Children[newIndex]
  return nil
}

func (t *Tree) Right () error {
  var a = t.Active
  if a == t.Root {
    return errors.New("Root node has no siblings")
  }
  var numChildren = len(a.Parent.Children)
  var newIndex = (a.Index + 1) % numChildren
  t.Active = t.Active.Parent.Children[newIndex]
  return nil
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
  i1 := n1.Index
  i2 := n2.Index
  err = p1.DeleteChild(i1)
  if err != nil { return err }
  err = p2.InsertChild(n1.Data, i2)
  if err != nil { return err }
  t.Active = p2.Children[i2]
  p1.IndexChildren()
  p2.IndexChildren()
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
