package utils

import (
  "errors"
)

type TreeNode struct {
  Data interface{}
  Index int
  Parent *TreeNode
  Children []*TreeNode
}

/*func (p *TreeNode) ReIndexChildren () {
  for i, c := range p.Children {
    c.Index = i
  }
}*/

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

// We should NOT be manually reindexing the child nodes!

func (t *Tree) DeleteChild (i int) error {
  var a = t.Active
  if i >= len(a.Children) {
    return errors.New("Child doesn't exist")
  }
  var cs = a.Children[i+1:]
  a.Children = a.Children[:i]
  for j, c := range cs {
    c.Index = i + j
    a.Children = append(a.Children, c)
  }
  return nil
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

/*func (t *Tree) Print (strFn func(interface{}) string) {
  t.DepthFirstTraverse(
    func(n *TreeNode) {
      fmt.Printf("%s\n", strFn(n.Data))
    },
  )
}*/











