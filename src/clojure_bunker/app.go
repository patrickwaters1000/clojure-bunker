// This file contains logic for translating termbox events into editor
// commands, respecting the app's mode. I also implemented at the app level
// that some commands must be completed in a minibuffer.
package main

import (
  termbox "github.com/nsf/termbox-go"
)

type App struct {
  editor *Editor
  mode string
  partialCmd []string
}

func NewApp () *App {
  ed := NewEditor()
  return &App{ed, "normal", nil}
}

// Translating the keyboard input to a command adds lines of code that
// aren't absolutely necessary. We could instead pass the termbox events
// to (for example) the active buffer's handler function. But I like that
// this way most of the interface with termbox is localized to this file. 
func api (mode string, ev termbox.Event) []string {
  var cmd []string
  switch mode {
  case "normal":
    if ev.Ch != 0 {
      switch ev.Ch {
      case 'c': cmd = []string{"center-window"}
      case 'r': cmd = []string{"partial", "repl", "connect"}
      case 't': cmd = []string{"repl", "eval"}
      case 'h': cmd = []string{"buffer", "move", "left"}
      case 'j': cmd = []string{"buffer", "move", "down"}
      case 'k': cmd = []string{"buffer", "move", "up"}
      case 'l': cmd = []string{"buffer", "move", "right"}
      case 's': cmd = []string{"partial", "buffer", "insert", "symbol", "after"}
      case 'S': cmd = []string{"partial", "buffer", "insert", "symbol", "before"}
      case 'd': cmd = []string{"buffer", "delete"}
      case 'n': cmd = []string{"partial", "new-buffer"}
      case 'z': cmd = []string{"kill-buffer"}
      case 'i': cmd = []string{"set-mode", "insert"}
      case 'w': cmd = []string{"partial", "write-file"}
      case 'e': cmd = []string{"partial", "load-file"}
      }
    } else {
      switch ev.Key {
      case termbox.KeyCtrlQ: cmd = []string{"quit"}
      case termbox.KeyArrowRight: cmd = []string{"next-buffer"}
      case termbox.KeyCtrlC: cmd = []string{"buffer", "insert", "call", "below"}
      case termbox.KeyCtrlS:
        cmd = []string{"partial", "buffer", "insert", "symbol", "below"}
      case termbox.KeyCtrlV: cmd = []string{"buffer", "insert", "vect", "below"}
      }
    }
  case "insert":
    if ev.Key != 0 {
      switch ev.Key {
      case termbox.KeyCtrlQ:  cmd = []string{"quit"}
      case termbox.KeyEsc:    cmd = []string{"buffer", "set-mode", "not-insert"}
      case termbox.KeyCtrlH:  cmd = []string{"buffer", "swap", "left"}
      case termbox.KeyCtrlJ:  cmd = []string{"buffer", "swap", "down"}
      case termbox.KeyCtrlK:  cmd = []string{"buffer", "swap", "up"}
      case termbox.KeyCtrlL:  cmd = []string{"buffer", "swap", "right"}
      case termbox.KeySpace:  cmd = []string{"buffer", "append", "token"}
      case termbox.KeyDelete: cmd = []string{"buffer", "append", "backspace"}
      case termbox.KeyCtrlC:  cmd = []string{"buffer", "append", "open", "call"}
      case termbox.KeyCtrlV:  cmd = []string{"buffer", "append", "open", "vect"}
      case termbox.KeyCtrlT:  cmd = []string{"buffer", "toggle-style"}
      }
    } else {
      switch ev.Ch {
      case ' ': cmd = []string{"buffer", "append", "token"}
      default:  cmd = []string{"buffer", "append", "string", string(ev.Ch)}
      }
    }
  case "minibuffer":
    switch ev.Key {
    case termbox.KeyCtrlQ: cmd = []string{"quit"}
    case termbox.KeyEnter:  cmd = []string{"finish-partial"}
    case termbox.KeyDelete: cmd = []string{"minibuffer", "delete"}
    default: cmd = []string{"minibuffer", "append", string(ev.Ch)}
    }
  }
  return cmd
}

func partialCmdPrompt(cmd []string) string {
  var s string
  switch cmd[0] {
  case "insert":     s = "symbol: "
  case "new-buffer": s = "buffer name: "
  case "write-file": s = "write buffer to: "
  case "load-file":  s = "edit file: "
  default: panic("Not found")
  }
  return s
}

func (app *App) handle (cmd []string) {
  switch cmd[0] {
  case "partial":
    app.mode = "minibuffer"
    app.partialCmd = cmd[1:]
    prompt := partialCmdPrompt(cmd[1:])
    app.editor.getMiniBuffer().reset(prompt)
  case "finish-partial":
    app.mode = "normal"
    fullCmd := append(app.partialCmd, app.editor.getMiniBuffer().data)
    app.editor.getMiniBuffer().reset("")
    app.editor.handle(fullCmd)
  case "set-mode":
    app.mode = cmd[1]
    app.editor.handle(cmd)
  default: app.editor.handle(cmd)
  }
}

func main() {
  err := termbox.Init()
  defer termbox.Close()
  panicIfError(err)

  app := NewApp()
  app.editor.render()

  for {
    app.editor.logState()
    event := termbox.PollEvent()
    cmd := api(app.mode, event)
    if cmd[0] == "quit" {
      break
    } else {
      app.handle(cmd)
      app.editor.render()
    }
  }
}

