package main

import (
  "bufio"
  "fmt"
  "os"
)

//var enterKey rune = 13
//var deleteKey rune = 127
//var clearLine = "\x1b[K"
//var clearScreen = "\x1b[2J"
//var cursorHome = "\x1b[;H"
//var cursorJump = "\x1b[%d;%dH"

//func panicIfError(err interface{}) {
//  if err != nil {
//    panic(err)
//  }
//}


func readFromStdin(
  stopCond func(rune) bool,
  handleChar func(rune),
) {
  var err error
  var c rune
  reader := bufio.NewReader(os.Stdin)
  for err == nil {
    c, _, err = reader.ReadRune()
    if stopCond(c) {
      break
    } else {
      handleChar(c)
    }
  }
  panicIfError(err)
}

type miniBuffer struct {
  data []rune
  pos [2]int
  prompt string
}

func (m miniBuffer) start() {
  out := fmt.Sprintf(cursorJump, m.pos[0], m.pos[1])
  out += clearLine
  out += m.prompt
  fmt.Print(out)
}

func (m miniBuffer) printToScreen() {
  out := fmt.Sprintf(cursorJump, m.pos[0], m.pos[1])
  out += clearLine
  out += m.prompt
  for _,c := range(m.data) {
    out += string(c)
  }
  fmt.Print(out)
}

func (m *miniBuffer) handleChar (c rune) {
  if c == deleteKey {
    m.data = m.data[:len(m.data)-1]
  } else {
    m.data = append(m.data, c)
  }
  m.printToScreen()
}

func (m miniBuffer) quit () {
  out := fmt.Sprintf(cursorJump, m.pos[0], m.pos[1])
  out += clearLine
  fmt.Print(out)
}

func readInputInMiniBuffer(prompt string) string {
  m := miniBuffer{[]rune{}, [2]int{20,1}, prompt}
  m.start()
  readFromStdin(
    func(c rune) bool {
      if c == enterKey {
        return true
      } else {
        return false
      }
    },
    func(c rune) {
      m.handleChar(c)
    },
  )
  v := ""
  for _,c := range(m.data) {
    v += string(c)
  }
  m.quit()
  return v
}
