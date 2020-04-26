package main

import (
  "fmt"
  u "utils"
  cu "clj_utils"
)

func printNode (n *u.TreeNode) {
  var d = (n.Data).(string)
  fmt.Printf("%s ", d)
}

type CoolThing struct {
  a int
  b int
}

func printCoolThingNode (n *u.TreeNode) {
  var d = (n.Data).(*CoolThing)
  fmt.Printf("(%d, %d) ", d.a, d.b)
}

func printTokenNode (n *u.TreeNode) {
  var d = (n.Data).(*cu.Token)
  fmt.Printf("(%s) ", d.Value)
}



func main() {
  var t = u.NewTree("root")
  t.AppendChild("c1")
  t.AppendChild("c2")
  t.Down()
  t.AppendChild("g1")
  t.DepthFirstTraverse(printNode)


  ct1 := CoolThing{1,2}
  ct2 := CoolThing{3,4}
  t2 := u.NewTree(&ct1)
  t2.AppendChild(&ct2)
  t2.DepthFirstTraverse(printCoolThingNode)

  tok1 := cu.NewToken("(")
  tok2 := cu.NewToken("ja")
  t3 := u.NewTree(&tok1)
  t3.AppendChild(&tok2)
  t3.DepthFirstTraverse(printTokenNode)


}
