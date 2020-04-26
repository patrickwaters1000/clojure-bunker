package main

import (
  "bufio"
  "fmt"
  "os"
  //"unicode"
  "golang.org/x/crypto/ssh/terminal"
  u "utils"
  cu "clj_utils"
  "io/ioutil"
)

var enterKey rune = 13
var deleteKey rune = 127
var clearLine = "\x1b[K"
var clearScreen = "\x1b[2J"
var cursorHome = "\x1b[;H"
var cursorJump = "\x1b[%d;%dH"

type editor struct {
  tree *u.Tree
}

func panicIfError(err interface{}) {
  if err != nil {
    panic(err)
  }
}

func readTokenFromStdin () *cu.Token {
  tokenStr := readInputInMiniBuffer("symbol: ")
  return cu.NewToken(tokenStr)
}


func newLeafFirst(tree *u.Tree) {
  token := readTokenFromStdin()
  tree.InsertSibling(token, -1)
}

func newParens(tree *u.Tree) {
  tree.InsertSibling(cu.NewToken("("), -1)
  err := tree.Left()
  tree.AppendChild(cu.NewToken(")"))
  err = tree.DownFirst()
  panicIfError(err)
}

func (ed *editor) refresh () {
  codeStr := cu.UnParseClj(ed.tree)
  active := ed.tree.Active.Data.(*cu.Token)
  printNode := fmt.Sprintf("\x1b[15;0H\x1b[K%v", active)
  cursorJumpStr := fmt.Sprintf(cursorJump, active.Row+1, active.Col+1)
  fmt.Print(clearScreen, cursorHome, codeStr, printNode, cursorJumpStr)
  //traverseFn := func (n *u.TreeNode) {
  //  token := *n.Data.(*cu.Token)
  //  fmt.Printf("\x1b[15;0H\x1b[K%v\x1b[%d;%dH", token, token.Row, token.Col)
  //  reader := bufio.NewReader(os.Stdin)
  //  reader.ReadRune()
  //}
  //ed.tree.DepthFirstTraverse(traverseFn)
}

func (ed *editor) writeFile() {
  fileName := readInputInMiniBuffer("file: ")
  data := []byte(cu.UnParseClj(ed.tree))
  ioutil.WriteFile(fileName, data, 0644)
}

func printMsg(msg string) {
  out := fmt.Sprintf(cursorJump, 21, 0)
  out += clearLine
  out += msg
  fmt.Print(out)
}

func (ed *editor) handleInput(c rune) {
  var err error
  switch c {
  case 'h': err = ed.tree.Left()
  case 'j': err = ed.tree.DownFirst()
  case 'k': err = ed.tree.Up()
  case 'l': err = ed.tree.Right()
  case 's': newLeafFirst(ed.tree)
  case 'd': newParens(ed.tree)
  case 'w': ed.writeFile()
  }
  if err != nil {
    panic(err)
    //printMsg("Fail")
  }
}

func main() {
  oldState, err := terminal.MakeRaw(0)
  if err != nil {
    panic(err)
  }
  defer terminal.Restore(0, oldState)

  data, _ := ioutil.ReadFile("example.clj")
  var tree *u.Tree = cu.ParseClj(data)
  ed := editor{tree}
  ed.refresh()

  reader := bufio.NewReader(os.Stdin)
  //fmt.Print("\x1b[2J")
  var c rune
  for err == nil {
    c, _, err = reader.ReadRune()
    if c == 'q' {
      break
    } else {
      ed.handleInput(c)
    }
    ed.refresh()
  }
}
