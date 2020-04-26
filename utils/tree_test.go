package utils

import (
  "testing"
)

func getExampleTree () *Tree {
  tree := NewTree("p")
  tree.AppendChild("c1")
  tree.AppendChild("c2")
  return tree
}

//func getTestDataTraverseFn (ns []string) func(*TreeNode) {
//  return func(n *TreeNode) {
//    ns = ns.append(n.Data.(string))
//  }
//}

func TestInsertSibling(t *testing.T) {
  tree := getExampleTree()
  tree.InsertChild("x", 0)
  wants := []string{"x","c1","c2"}
  for i, n := range tree.Root.Children {
    got := n.Data.(string)
    if got != wants[i] {
      t.Errorf(
        "Test case 1.%d: want %s, got %s",
        i, wants[i], got,
      )
    }
  }
  tree = getExampleTree()
  tree.InsertChild("x", 1)
  wants = []string{"c1","x","c2"}
  for i, n := range tree.Root.Children {
    got := n.Data.(string)
    if got != wants[i] {
      t.Errorf(
        "Test case 2.%d: want %s, got %s",
        i, wants[i], got,
      )
    }
  }
}
