// This file contains logic for translating termbox events into editor
// commands, respecting the app's mode. I also implemented at the app level
// that some commands must be completed in a minibuffer.
package main

import (
  termbox "github.com/nsf/termbox-go"
)

type App struct {
  editor *Editor
  inMinibuffer bool
  partialCmd []string
}

func NewApp () *App {
  ed := NewEditor()
  return &App{ed, false, nil}
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
      case 'u': cmd = []string{"buffer", "undo"}
      case 'c': cmd = []string{"center-window"}
      case 'r': cmd = []string{"partial", "repl", "connect"}
      case 't': cmd = []string{"repl", "eval"}
      case 'h': cmd = []string{"buffer", "move-normal", "left"}
      case 'j': cmd = []string{"buffer", "move-normal", "down"}
      case 'k': cmd = []string{"buffer", "move-normal", "up"}
      case 'l': cmd = []string{"buffer", "move-normal", "right"}
      case 'd': cmd = []string{"buffer", "delete"}
      case 'n': cmd = []string{"partial", "new-buffer"}
      case 'z': cmd = []string{"kill-buffer"}
      case 'i': cmd = []string{"set-mode", "insert"}
      case 'w': cmd = []string{"partial", "write-file"}
      case 'e': cmd = []string{"partial", "load-file"}
      }
    } else {
      switch ev.Key {
      case termbox.KeyCtrlQ:      cmd = []string{"quit"}
      case termbox.KeyArrowRight: cmd = []string{"next-buffer"}
      case termbox.KeyCtrlS:      cmd = []string{"buffer", "toggle-style"}
      }
    }
  case "insert":
    if ev.Key != 0 {
      switch ev.Key {
      case termbox.KeyCtrlQ:  cmd = []string{"quit"}
      case termbox.KeyCtrlU: cmd = []string{"buffer", "undo"}
      case termbox.KeyEsc:    cmd = []string{"set-mode", "not-insert"}
      case termbox.KeyCtrlH:  cmd = []string{"buffer", "move-insert", "left"}
      case termbox.KeyCtrlJ:  cmd = []string{"buffer", "move-insert", "down"}
      case termbox.KeyCtrlK:  cmd = []string{"buffer", "move-insert", "up"}
      case termbox.KeyCtrlL:  cmd = []string{"buffer", "move-insert", "right"}
      case termbox.KeySpace:  cmd = []string{"buffer", "append", "token"}
      case 127: cmd = []string{"buffer", "append", "backspace"}
      case termbox.KeyCtrlC:  cmd = []string{"buffer", "append", "open", "call"}
      case termbox.KeyCtrlV:  cmd = []string{"buffer", "append", "open", "vect"}
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
    case termbox.KeyEsc:   cmd = []string{"abort-partial"}
    case termbox.KeyEnter: cmd = []string{"finish-partial"}
    // 127 is backspace, but doesn't match termbox.KeyBackspace
    case 127: cmd = []string{"minibuffer", "delete"}
    default: cmd = []string{"minibuffer", "append", string(ev.Ch)}
    }
  default:
    panic("Invalid mode " + mode)
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
  case "repl":
    switch cmd[1] {
    case "connect": s = "nrepl port: "
    }
  default: panic("Not found")
  }
  return s
}

func (app *App) handle (cmd []string) {
  switch cmd[0] {
  case "partial":
    app.inMinibuffer = true
    app.partialCmd = cmd[1:]
    prompt := partialCmdPrompt(cmd[1:])
    app.editor.getMiniBuffer().reset(prompt)
  case "abort-partial":
    app.inMinibuffer = false
    app.editor.getMiniBuffer().reset("")
  case "finish-partial":
    app.inMinibuffer = false
    fullCmd := append(app.partialCmd, app.editor.getMiniBuffer().data)
    app.editor.getMiniBuffer().reset("")
    app.editor.handle(fullCmd)
  default: app.editor.handle(cmd)
  }
}

func (a App) getMode () string {
  if a.inMinibuffer {
    return "minibuffer"
  }
  b := a.editor.getActiveBuffer()
  switch b.(type) {
  case *CodeBuffer:
    return b.(*CodeBuffer).tree.Root.Data.(*Token).Value
  default:
    return "normal"
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
    currentMode := app.getMode()
    cmd := api(currentMode, event)
    if len(cmd) == 0 {
      // key not mapped to command; pass
    } else if cmd[0] == "quit" {
      break
    } else {
      app.handle(cmd)
      app.editor.render()
    }
  }
}

