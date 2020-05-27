package main

import (
  "errors"
)

type Copier interface {
  Copy() Copier
}

// Returns a copy of `t` that is JUST deep enough that mutating the active
// node doesn't alter the original tree.
// operations behave as if we were truly using deep copies.
func (tOld *Tree) PersistentCopy () *Tree {
  tNew := NewTree(tOld.Root.Data)
  l := len(tOld.Path)
  for i, n := range tOld.Path[:l - 1] {
    cPath := tOld.Path[i + 1] // child of n on path to active node
    var jPath int // will be index of cPath
    for j, c := range n.Children {
      if c == cPath {
        token := c.Data.(Copier).Copy()
        tNew.AppendChild(token)
        jPath = j
      } else {
        aNew := tNew.GetActive()
        aNew.Children = append(aNew.Children, c)
      }
    }
    newPathNode := tNew.GetActive().Children[jPath]
    tNew.Path = append(tNew.Path, newPathNode)
  }
  aNew := tNew.GetActive()
  for _, c := range tOld.GetActive().Children {
    aNew.Children = append(aNew.Children, c)
  }
  return tNew
}

func (t *Tree) UpdateActive (updateFn func(*TreeNode)) *Tree {
  tNew := t.PersistentCopy()
  a := tNew.GetActive()
  updateFn(a)
  return tNew
}

func (t *Tree) PersistentMove (move []rune) (*Tree, error) {
  pathCopy := append(
    []*TreeNode{},
    t.Path...)
  tNew := &Tree{
    Root: t.Root,
    Path: pathCopy,
  }
  err := tNew.Move(move)
  return tNew, err
}

func (t *Tree) PersistentAppend (d interface{}) *Tree {
  tNew := t.PersistentCopy()
  tNew.GetActive().AppendChild(d)
  return tNew
}

func (t *Tree) PersistentInsert (d interface{}, i int) (*Tree, error) {
  tNew := t.PersistentCopy()
  err := tNew.GetActive().InsertChild(d, i)
  return tNew, err
}

func (t *Tree) PersistentInsertSibling (d interface{}, i int) (*Tree, error) {
  tNew, err := t.PersistentMove([]rune{'u'})
  if err != nil {
    return tNew, err
  }
  tNew, err = tNew.PersistentInsert(d, i)
  if err != nil {
    return tNew, err
  }
  move := []rune{'d'}
  for j:=0; j<i; j++ {
    move = append(move, 'r')
  }
  tNew, err = t.PersistentMove(move)
  return tNew, err
}

func (t *Tree) PersistentDelete (i int) (*Tree, error) {
  tNew := t.PersistentCopy()
  err := tNew.GetActive().DeleteChild(i)
  return tNew, err
}

func (t *Tree) PersistentDeleteActive () (*Tree, int, error) {
  i, err := t.GetActiveIndex()
  if err != nil {
    return t, 0, errors.New("Can't delete root node")
  }
  tNew, err := t.PersistentMove([]rune{'u'})
  if err != nil {
    return tNew, 0, err
  }
  tNew, err = tNew.PersistentDelete(i)
  return tNew, i, err
}


