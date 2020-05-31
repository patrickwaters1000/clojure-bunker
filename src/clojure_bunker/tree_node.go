package main

import (
  "errors"
)

type TreeNode struct {
  Data interface{}
  Children []*TreeNode
}

func NewTreeNode (d interface{}) *TreeNode {
  return &TreeNode{
    Data: d,
    Children: []*TreeNode{},
  }
}

func (n *TreeNode) AppendChild (d interface{}) {
  c := NewTreeNode(d)
  n.Children = append(n.Children, c)
}

func (n *TreeNode) InsertChild (d interface{}, i int) error {
  if i > len(n.Children) {
    return errors.New("Insertion index is too large")
  }
  newCs := []*TreeNode{}
  newCs = append(newCs, n.Children[:i]...)
  newCs = append(newCs, NewTreeNode(d))
  for _,c := range n.Children[i:] {
    newCs = append(newCs, c)
  }
  n.Children = newCs
  return nil
}

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

func (r *TreeNode) DepthFirstTraverseNoRoot (f func (*TreeNode)) {
  var s = NewStack()
  var n *TreeNode
  l := len(r.Children)
  for i:=l-1; i>=0; i-- {
    s.Push(r.Children[i])
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
