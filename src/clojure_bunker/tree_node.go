package main

import (
  "errors"
)

type TreeNode struct {
  Data interface{}
  Index int // Deprecated
  Parent *TreeNode
  Children []*TreeNode
}

func (n *TreeNode) GetIndex () (int, error) {
  p := n.Parent
  if p == nil {
    return 0, errors.New("Node has no parent")
  }
  for i,c := range p.Children {
    if c == n {
      return i, nil
    }
  }
  panic("Something is wrong")
}

func (n *TreeNode) Up () (*TreeNode, error) {
  if n.Parent == nil {
    return nil, errors.New("Node has no parent")
  }
  return n.Parent, nil
}

func (n *TreeNode) DownFirst () (*TreeNode, error) {
  if len(n.Children) == 0 {
    return nil, errors.New("Node has no children")
  }
  return n.Children[0], nil
}

func (n *TreeNode) DownLast () (*TreeNode, error) {
  nCs := len(n.Children)
  if nCs == 0 {
    return nil, errors.New("Node has no children")
  }
  return n.Children[nCs-1], nil
}

func (n *TreeNode) Cycle (d int) (*TreeNode, error) {
  p, err := n.Up()
  if err != nil { return nil, err }
  i1, err := n.GetIndex()
  if err != nil { return nil, err }
  nCs := len(p.Children)
  i2 := mod(i1+d, nCs)
  return p.Children[i2], nil
}

func (n *TreeNode) Left () (*TreeNode, error) {
  return n.Cycle(-1)
}

func (n *TreeNode) Right () (*TreeNode, error) {
  return n.Cycle(1)
}

func (n *TreeNode) CycleNotLast (d int) (*TreeNode, error) {
  p, err := n.Up()
  if err != nil { return nil, err }
  nCs := len(p.Children)
  if nCs == 1 {
    return nil, errors.New("Only one child")
  }
  var r *TreeNode
  r, err = n.Cycle(d)
  if err != nil { return nil, err }
  var i int
  i, err = r.GetIndex()
  if err != nil { return nil, err }
  if i == nCs-1 {
    r, err = r.Cycle(d)
  }
  return r,err
}

func (n *TreeNode) LeftNotLast () (*TreeNode, error) {
  return n.CycleNotLast(-1)
}

func (n *TreeNode) RightNotLast () (*TreeNode, error) {
  return n.CycleNotLast(1)
}

func (n *TreeNode) InsertChild (d interface{}, i int) error {
  if i > len(n.Children) {
    return errors.New("Insertion index is too large")
  }
  newCs := []*TreeNode{}
  newCs = append(newCs, n.Children[:i]...)
  newCs = append(newCs, &TreeNode{d, i, n, []*TreeNode{}})
  for _,c := range n.Children[i:] {
    c.Index++
    newCs = append(newCs, c)
  }
  n.Children = newCs
  return nil
}

// We should NOT be manually reindexing the child nodes!

func (n *TreeNode) DeleteChild (i int) error {
  if i >= len(n.Children) {
    return errors.New("Child doesn't exist")
  }
  var cs = n.Children[i+1:]
  n.Children = n.Children[:i]
  for j, c := range cs {
    c.Index = i + j
    n.Children = append(n.Children, c)
  }
  return nil
}

// TODO call this fn instead of manually changing indices
// Better TODO remove Index attribute
func (n *TreeNode) IndexChildren () {
  for i,c := range n.Children {
    c.Index = i
  }
}

