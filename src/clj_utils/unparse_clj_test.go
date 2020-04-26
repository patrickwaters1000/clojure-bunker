package clj_utils

import (
  u "utils"
  //cmp "go-cmp"
  "testing"
  "fmt"
)

func TestLetBinding(t *testing.T) {
  clj := "(let [x 42\n" +
         "      y"
  tree := ParseClj([]byte(clj))

  //traverseFn := func (n *u.TreeNode) {
  //  fmt.Printf("%v\n", *n.Data.(*Token))
  //}
  //tree.DepthFirstTraverse(traverseFn)
  n := tree.Active.Children[2]
  v1 := n.Data.(*Token).Value
  if v1 != "y" {
    t.Errorf("Testing active node value. Want 'y' got %s", v1)
  }

  v2 := n.Parent.Data.(*Token).Value
  if v2 != "[" {
    t.Errorf("Testing active node value. Want '[' got %s", v2)
  }

  v3 := getLeader(n.Parent)
  if v3 != "let" {
    t.Errorf("Testing active node value. Want 'let' got %s", v3)
  }

  got4 := n.Index
  want4 := 2
  if want4 != got4 {
    t.Errorf("Testing let binding newline. Want %d got %d", want4, got4)
  }
}

type intPair struct {
  r int
  c int
}

func TestPositions (t *testing.T) {
  clj := "(let [x 42\n" +
         "      y 8]\n" +
         "  (+ x y))"
  tree := ParseClj([]byte(clj))
  treeStr := UnParseClj(tree)
  fmt.Println(treeStr)
  wants := []intPair{
    intPair{0, -4}, // case 0
    intPair{0, 0},
    intPair{0, 1},
    intPair{0, 5},
    intPair{0, 6},
    intPair{0, 8},
    intPair{1, 6}, // case 6
    intPair{1, 8},
    intPair{1, 9},
    intPair{2, 2}, // case 9
    intPair{2, 3},
    intPair{2, 5},
    intPair{2, 7}, // case 12
    intPair{2, 8},
    intPair{2, 9},
  }
  gots := []intPair{}
  traverseFn := func (n *u.TreeNode) {
    t := n.Data.(*Token)
    gots = append(gots, intPair{t.Row, t.Col})
  }
  tree.DepthFirstTraverse(traverseFn)
  for i:=0; i<len(wants); i++ {
    if wants[i] != gots[i] {
      t.Errorf(
        "Case %d; want %v; got %v",
        i, wants[i], gots[i],
      )
    }
  }
}

func TestDefn(t *testing.T) {
  //w := &CljWriter{""}
  clj := "(defn f [x]\n" +
         "  42)"
  tree := ParseClj([]byte(clj))
  _ = UnParseClj(tree)

  wantRows := [9]int{0, 0, 0, 0, 0, 0, 0, 1, 1}
  wantCols := [9]int{-4, 0, 1, 6, 8, 9, 10, 2, 4}
  gotRows := [9]int{}
  gotCols := [9]int{}
  var i int = 0
  traverseFn := func(n *u.TreeNode) {
    //fmt.Printf("%v\n", *n.Data.(*Token))
    t := n.Data.(*Token)
    gotRows[i] = t.Row
    gotCols[i] = t.Col
    i += 1
  }
  tree.DepthFirstTraverse(traverseFn)
  if wantRows != gotRows {
    t.Errorf(
      "Testing defn rows\n"+
      "want: %v\n"+
      "got: %v\n",
      wantRows, gotRows,
    )
  }
  if wantCols != gotCols {
    t.Errorf(
      "Testing defn cols\n"+
      "want: %d\n"+
      "got: %d\n",
      wantCols, gotCols,
    )
  }
}












