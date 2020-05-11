package main

import (
  "errors"
  termbox "github.com/nsf/termbox-go"
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

type Window struct {
  rows int
  cols int
  pos_r int
  pos_c int
  screen_r int
  screen_c int
}

func (w Window) Print (row, col int, fg, bg termbox.Attribute, msg string) {
  r := row - w.pos_r
  c := col - w.pos_c
  if and(0 <= r,
         0 <= c,
         r < w.rows,
         c < w.cols) {
    tbPrint(r + w.screen_r, c + w.screen_c, fg, bg, msg)
  }
}

type Handler interface {
  handleEvent([]string) error
  render(Window)
}

type MiniBuffer struct {
  data string
  prompt string
  callback func(string) []string
}

type App struct {
  editor Handler
  editorWindow Window
  miniBuffer MiniBuffer
  miniBufferWindow Window
  mode string
}

func NewApp() *App {
  editor := NewEditor()
  miniBuffer := MiniBuffer{"" ,"", nil}
  rows, cols := get_winsize()
  return &App{
    editor: editor,
    editorWindow: Window{rows-5, cols, 0, 0, 0, 0},
    miniBuffer: miniBuffer,
    miniBufferWindow: Window{5, cols, 0, 0, rows-5, 0},
    mode: "normal",
  }
}

func (app *App) finishCmdInMiniBuffer(prompt string, partialCmd []string) {
  app.mode = "miniBuffer"
  callback := func(input string) []string {
    return append(partialCmd, input)
  }
  app.miniBuffer = MiniBuffer{"", prompt, callback}
}

func (app *App) handleEvent(ev termbox.Event) error {
  var cmd []string
  var quit bool
  switch app.mode {
  case "normal":
    if ev.Ch != 0 {
      switch ev.Ch {
      case 'q': quit = true
      case 'r':
        app.finishCmdInMiniBuffer(
          "port: ", []string{"repl", "connect"})
      case 't': cmd = []string{"repl", "eval"}
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
      case 'e':
        app.finishCmdInMiniBuffer(
          "edit file: ", []string{"load"})
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
      cmd = app.miniBuffer.callback(app.miniBuffer.data)
      app.mode = "normal"
      app.miniBuffer = MiniBuffer{"", "", nil}
    } else if rune(ev.Key) == deleteKey {
      l := len(app.miniBuffer.data)
      app.miniBuffer.data = app.miniBuffer.data[:l-1]
    } else {
      app.miniBuffer.data += string(ev.Ch)
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

func (m MiniBuffer) render(w Window) {
  w.Print(0, 0, fg1, bg1, m.prompt + m.data)
}

func (app *App) render() {
  termbox.Clear(bg1, bg1)
  app.editor.render(app.editorWindow)
  app.miniBuffer.render(app.miniBufferWindow)
  //tbPrint(21, 0, fg1, bg1, fmt.Sprintf("Rows:%d, Cols:%d", app.height, app.width))
  termbox.Flush()
}

func main() {
  err := termbox.Init()
  defer termbox.Close()
  panicIfError(err)

  app := NewApp()
  //app.render()

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

