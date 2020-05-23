package main

import (
  "errors"
)

type TreeNode struct {
  Data interface{}
  Parent *TreeNode
  Children []*TreeNode
}

func NewTreeNode (d interface{}) *TreeNode {
  return &TreeNode{d, nil, []*TreeNode{}}
}

func (n *TreeNode) checkChildParentConsistency () bool {
  for _,c := range n.Children {
    if c.Parent != n {
      return false
    }
  }
  return true
}

func (n *TreeNode) GetIndex () int {
  p := n.Parent
  if p == nil {
    return 0
  }
  for i,c := range p.Children {
    if c == n {
      return i
    }
  }
  panic("Something is wrong")
}

func (n *TreeNode) GetSiblings () []*TreeNode {
  if n.Parent == nil {
    return []*TreeNode{n}
  }
  return n.Parent.Children
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

func (n *TreeNode) Cycle (d int) *TreeNode {
  cs := n.GetSiblings()
  i := mod(n.GetIndex() + d, len(cs))
  return cs[i]
}

func (n *TreeNode) Left () *TreeNode {
  return n.Cycle(-1)
}

func (n *TreeNode) Right () *TreeNode {
  return n.Cycle(1)
}

func (n *TreeNode) CycleNotLast (d int) (*TreeNode, error) {
  nCs := len(n.GetSiblings())
  if nCs == 1 {
    return nil, errors.New("Only one child")
  }
  r := n.Cycle(d)
  if r.GetIndex() == nCs-1 {
    r = r.Cycle(d)
  }
  return r, nil
}

func (n *TreeNode) LeftNotLast () (*TreeNode, error) {
  return n.CycleNotLast(-1)
}

func (n *TreeNode) RightNotLast () (*TreeNode, error) {
  return n.CycleNotLast(1)
}

func (n *TreeNode) AppendChild (c *TreeNode) {
  n.Children = append(n.Children, c)
  c.Parent = n
}

func (n *TreeNode) InsertChild (d interface{}, i int) error {
  if i > len(n.Children) {
    return errors.New("Insertion index is too large")
  }
  newCs := []*TreeNode{}
  newCs = append(newCs, n.Children[:i]...)
  newCs = append(newCs, &TreeNode{d, n, []*TreeNode{}})
  for _,c := range n.Children[i:] {
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
  for _, c := range cs {
    n.Children = append(n.Children, c)
  }
  return nil
}

