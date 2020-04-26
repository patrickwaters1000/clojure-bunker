package main

import (
  "fmt"
  u "utils"
  cu "clj_utils"
  "io/ioutil"
)

func printTokenNode (n *u.TreeNode) {
  var d = (n.Data).(cu.Token)
  fmt.Printf("(%s, %s, %d) ", d.Class, d.Value, len(n.Children))
}

func main() {
  data, _ := ioutil.ReadFile("example.clj")
  var tree *u.Tree = cu.ParseClj(data)
  var s string = cu.UnparseClj(tree)
  fmt.Printf("%s", s)
}
