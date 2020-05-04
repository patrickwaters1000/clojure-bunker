package main

import (
  //"bufio"
  //"fmt"
  //"os"
  //"unicode"
  "errors"
  termbox "github.com/nsf/termbox-go"
  //u "utils"
  //cu "clj_utils"
  //"io/ioutil"
)

const enterKey rune = 13
const deleteKey rune = 127
const rightKey rune = 65514
const leftKey rune = 65515
const downKey rune = 65516
const upKey rune = 65517
const bg1 = termbox.ColorBlack
const fg1 = termbox.ColorWhite
const bg2 = termbox.ColorGreen
const fgh = termbox.ColorBlack
const bgh = termbox.ColorWhite
const builtInColor = termbox.ColorYellow
const fnNameColor = termbox.ColorCyan
const commentColor = termbox.ColorBlue
const varNameColor = termbox.ColorMagenta
const keywordColor = termbox.ColorGreen
const stringColor = termbox.ColorRed
const symbolColor = termbox.ColorWhite

type Handler interface {
  handleEvent([]string) error
  render()
}

type App struct {
  editor Handler
  miniBuffer string
  miniBufferPrompt string
  miniBufferCallback func(string) []string
  mode string
}

func NewApp() *App {
  editor := NewEditor()
  return &App{
    editor: editor,
    miniBuffer: "",
    miniBufferPrompt: "",
    miniBufferCallback: nil,
    mode: "normal",
  }
}

func (app *App) finishCmdInMiniBuffer(prompt string, partialCmd []string) {
  app.mode = "miniBuffer"
  app.miniBuffer = ""
  app.miniBufferPrompt = prompt
  app.miniBufferCallback = func(input string) []string {
    return append(partialCmd, input)
  }
}

func (app *App) handleEvent(ev termbox.Event) error {
  var cmd []string
  var quit bool
  switch app.mode {
  case "normal":
    if ev.Ch != 0 {
      switch ev.Ch {
      case 'q': quit = true
      case 'h': cmd = []string{"buffer", "move", "left"}
      case 'j': cmd = []string{"buffer", "move", "down"}
      case 'k': cmd = []string{"buffer", "move", "up"}
      case 'l': cmd = []string{"buffer", "move", "right"}
      case 's':
        app.finishCmdInMiniBuffer(
          "symbol: ", []string{"buffer", "insert", "symbol", "after"})
      case 'S':
        app.finishCmdInMiniBuffer(
          "symbol: ", []string{"buffer", "insert", "symbol", "before"})
      case 'c': cmd = []string{"buffer", "insert", "call", "after"}
      case 'C': cmd = []string{"buffer", "insert", "call", "before"}
      case 'v': cmd = []string{"buffer", "insert", "vect", "after"}
      case 'V': cmd = []string{"buffer", "insert", "vect", "before"}
      case 'd': cmd = []string{"buffer", "delete"}
      case 'n':
      app.finishCmdInMiniBuffer(
        "buffer name: ", []string{"new-buffer"})
      case 'z': cmd = []string{"kill-buffer"}
      case 'i': cmd = []string{"buffer", "mode", "insert"}
      case 'w':
        app.finishCmdInMiniBuffer(
          "write buffer to: ", []string{"write"})
      }
    } else {
      if rune(ev.Key) == rightKey {
        cmd = []string{"next-buffer"}
      } else if rune(ev.Key) == 3 { // Ctrl + c
        cmd = []string{"buffer", "insert", "call", "below"}
      } else if rune(ev.Key) == 19 { // Ctrl + s
        app.finishCmdInMiniBuffer(
          "insert: ", []string{"buffer", "insert", "symbol", "below"})
      } else if rune(ev.Key) == 22 { // Ctrl + v
        cmd = []string{"buffer", "insert", "vect", "below"}
      }
    }
  case "miniBuffer":
    if rune(ev.Key) == enterKey {
      cmd = app.miniBufferCallback(app.miniBuffer)
      app.mode = "normal"
      app.miniBuffer = ""
      app.miniBufferPrompt = ""
      app.miniBufferCallback = nil
    } else if rune(ev.Key) == deleteKey {
      l := len(app.miniBuffer)
      app.miniBuffer = app.miniBuffer[:l-1]
    } else {
      app.miniBuffer += string(ev.Ch)
    }
  }
  if cmd != nil {
    err := app.editor.handleEvent(cmd)
    panicIfError(err)
  }
  if quit {
    return errors.New("quit")
  } else {
    return nil
  }
}

func (app *App) render() {
  termbox.Clear(bg1, bg1)
  app.editor.render()
  tbPrint(20, 0, fg1, bg1, app.miniBufferPrompt + app.miniBuffer)
  termbox.Flush()
}

func main() {
  err := termbox.Init()
  defer termbox.Close()
  panicIfError(err)

  app := NewApp()
  app.render()

  for {
    event := termbox.PollEvent()
    err := app.handleEvent(event)
    if err != nil {
      break
    } else {
      app.render()
    }
  }
}

