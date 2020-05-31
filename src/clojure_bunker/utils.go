package main

import (
  "io/ioutil"
  termbox "github.com/nsf/termbox-go"
  "syscall"
  "os"
  "unsafe"
  "fmt"
)

const bg1 = termbox.ColorBlack
const fg1 = termbox.ColorWhite
const bg2 = termbox.ColorGreen
const fgh = termbox.ColorBlack
const bgh = termbox.ColorWhite


func panicIfError(err interface{}) {
  if err != nil {
    panic(err)
  }
}

func mod(n, d int) int {
  if d == 0 {
    return 0
  }
  m := n % d
  if m < 0 {
    m += d
  }
  return m
}

func and (conds ...bool) bool {
  for _, cond := range conds {
    if !cond { return false }
  }
  return true
}

func or (conds ...bool) bool {
  for _, cond := range conds {
    if cond { return true }
  }
  return false
}

func tbPrint(row, col int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(col, row, c, fg, bg)
    col++
  }
}

func log (msg string) {
  err := ioutil.WriteFile("log", []byte(msg), 0644)
  panicIfError(err)
}

type winsize struct {
  rows uint16
  cols uint16
  xpixels uint16
  ypixels uint16
}

func get_winsize() (int, int) {
  out, _ := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
  var size winsize
  _, _, _ = syscall.Syscall(
    syscall.SYS_IOCTL,
    out.Fd(),
    uintptr(syscall.TIOCGWINSZ),
    uintptr(unsafe.Pointer(&size)))
  _ = out.Close()
  return int(size.rows), int(size.cols)
}

func sprintTree (t *Tree) string {
  indentStr := ""
  msg := ""
  traverseFn := func (n *TreeNode) {
    t := n.Data.(*Token)
    msg += "\n" + indentStr
    msg += fmt.Sprintf("%s,%s", t.Class, t.Value)
    if t.IsOpen() {
      indentStr += "  "
    } else if t.IsClosed() {
      l := len(indentStr)
      indentStr = indentStr[:l-2]
    }
  }
  t.DepthFirstTraverse(traverseFn)
  return msg
}

func stringifySubtree (n *TreeNode) string {
  var msg string = ""
  var row int = 0
  var col int = 0
  traverseFn := func (n *TreeNode) {
    token := n.Data.(*Token)
    if token.Row > row {
      for i:=0; i<token.Row-row; i++ {
        msg += "\n"
        row += 1
      }
      col = 0
    }
    for j:=0; j<token.Col-col; j++ {
      msg += " "
      col += 1
    }
    msg += token.Value
    col += len(token.Value)
  }
  n.DepthFirstTraverseNoRoot(traverseFn)
  return msg
}
