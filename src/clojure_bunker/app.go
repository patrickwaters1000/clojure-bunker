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
  externalEventQueue chan termbox.Event
  internalEventQueue chan []string
}

func NewApp () *App {
  externalEventQueue := make(chan termbox.Event)
  internalEventQueue := make(chan []string)
  ed := NewEditor( // When something happens in the editor that requires
    // triggering a command, we have to put that command in the queue. A
    // callback for this must be supplied when building the editor. This is
    // really bad, and a compelling reason not to have separate app and
    // editor objects.
    func (cmd []string) {
      internalEventQueue <- cmd
    },
  )
  go func () {
    for {
      externalEventQueue <- termbox.PollEvent()
    }
  }()
  return &App{
    editor: ed,
    inMinibuffer: false,
    partialCmd: nil,
    externalEventQueue: externalEventQueue,
    internalEventQueue: internalEventQueue,
  }
}

// Translating the keyboard input to a command adds lines of code that
// aren't absolutely necessary. We could instead pass the termbox events
// to (for example) the active buffer's handler function. But I like that
// this way most of the interface with termbox is localized to this file. 
func getNormalCmd (ev termbox.Event) (bool, []string) {
  if ev.Ch != 0 {
    switch ev.Ch {
    case 'u': return true,  []string{"buffer", "undo"}
    case 'h': return true,  []string{"buffer", "move-normal-left"}
    case 'j': return true,  []string{"buffer", "move-normal-down"}
    case 'k': return true,  []string{"buffer", "move-normal-up"}
    case 'l': return true,  []string{"buffer", "move-normal-right"}
    case 'd': return true,  []string{"buffer", "delete"}
    case 'i': return true,  []string{"buffer", "enter-insert-mode"}
    case 'c': return true,  []string{"center-window"}
    case 'r': return false, []string{"repl-connect"}
    case 't': return true,  []string{"repl-eval"}
    case 'n': return false, []string{"new-buffer"}
    case 'z': return true,  []string{"kill-buffer"}
    case 'w': return true,  []string{"write-file"}
    case 'f': return false, []string{"set-file"}
    case 'e': return false, []string{"load-file"}
    }
  } else {
    switch ev.Key {
    case termbox.KeyCtrlQ:      return true,  []string{"quit"}
    case termbox.KeyArrowRight: return true,  []string{"next-buffer"}
    case termbox.KeyCtrlS:      return true,  []string{"buffer", "toggle-style"}
    case termbox.KeyCtrlW:      return false, []string{"window"}
    }
  }
  return false, []string{}
}

func getInsertCmd (ev termbox.Event) (bool, []string) {
  if ev.Key != 0 {
    switch ev.Key {
    case termbox.KeyCtrlQ: return true, []string{"quit"}
    case termbox.KeyCtrlU: return true, []string{"buffer", "undo"}
    case termbox.KeyEsc:   return true, []string{"buffer", "exit-insert-mode"}
    case termbox.KeyCtrlH: return true, []string{"buffer", "move-insert-left"}
    case termbox.KeyCtrlJ: return true, []string{"buffer", "move-insert-down"}
    case termbox.KeyCtrlK: return true, []string{"buffer", "move-insert-up"}
    case termbox.KeyCtrlL: return true, []string{"buffer", "move-insert-right"}
    case termbox.KeySpace: return true, []string{"buffer", "insert-token"}
    case 127:              return true, []string{"buffer", "insert-backspace"}
    case termbox.KeyCtrlC: return true, []string{"buffer", "insert-call"}
    case termbox.KeyCtrlV: return true, []string{"buffer", "insert-vect"}
    }
  } else {
    return true, []string{"buffer", "insert-string", string(ev.Ch)}
  }
  return false, []string{}
}


func (app *App) getMinibufferCmd (ev termbox.Event) (bool, []string) {
  switch ev.Key {
  case termbox.KeyCtrlQ: return true,  []string{"quit"}
  case termbox.KeyEsc:   return true,  []string{"abort"}
  case termbox.KeyEnter:
    cmd := app.editor.getMiniBuffer().data
    return true, []string{cmd}
  case 127:
    app.editor.getMiniBuffer().Delete()
    return false, []string{}
  default:
    app.editor.getMiniBuffer().Append(string(ev.Ch))
    return false, []string{}
  }
}

func getWindowCmd (ev termbox.Event) (bool, []string) {
  switch ev.Ch {
  case 'l': return true, []string{"next"}
  case 'r': return true, []string{"reset"}
  case 'v': return true, []string{"split-vertical"}
  case 'h': return true, []string{"split-horizontal"}
  }
  return false, []string{}
}

var minibufferCmdPrompts = map[string]string {
  "new-buffer":   "buffer name: ",
  "set-file":   "set file for buffer: ",
  "load-file":    "edit file: ",
  "repl-connect": "nrepl port: ",
}

func (app App) getMode () string {
  b := app.editor.getActiveBuffer()
  switch b.(type) {
  case *CodeBuffer:
    return b.(*CodeBuffer).tree.Root.Data.(*Token).Value
  default:
    return "normal"
  }
}

func (app App) getCmd (ev termbox.Event) (bool, []string) {
  if app.inMinibuffer {
    return app.getMinibufferCmd(ev)
  } else if len(app.partialCmd) > 0 {
    switch app.partialCmd[0] {
    case "window": return getWindowCmd(ev)
    default: panic("Unexpected incomplete command")
    }
  } else {
    switch app.getMode() {
    case "normal": return getNormalCmd(ev)
    case "insert": return getInsertCmd(ev)
    default: panic("Invalid mode")
    }
  }
}

func (app *App) exitMinibufferIfNecessary () {
  if app.inMinibuffer {
    app.inMinibuffer = false
    app.editor.getMiniBuffer().reset("")
  }
}

func (app *App) enterMiniBufferIfNecessary () {
  cmd := app.partialCmd
  if len(cmd) > 0 {
    prompt, useMinibuffer := minibufferCmdPrompts[cmd[0]]
    if !app.inMinibuffer && useMinibuffer {
      app.editor.getMiniBuffer().reset(prompt)
      app.inMinibuffer = true
    }
  }
}

func (app *App) handleCmd (done bool, cmd []string) bool {
  if !done {
    app.partialCmd = append(app.partialCmd, cmd...)
    app.enterMiniBufferIfNecessary()
    app.editor.render()
    return false
  } else if len(cmd) == 0 {
    return false
  } else if cmd[0] == "refresh" {
    app.editor.render()
    return false
  } else if cmd[0] == "abort" {
    app.exitMinibufferIfNecessary()
    app.partialCmd = []string{}
    app.editor.render()
    return false
  } else if cmd[0] == "quit" {
    return true
  } else {
    fullCmd := append(app.partialCmd, cmd...)
    app.editor.handle(fullCmd)
    app.exitMinibufferIfNecessary()
    app.partialCmd = []string{}
    app.editor.render()
    return false
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
    var done bool
    var cmd []string
    select {
    case ev := <-app.externalEventQueue:
      done, cmd = app.getCmd(ev)
    case internalCmd := <-app.internalEventQueue:
      done = true
      cmd = internalCmd
    }
    quit := app.handleCmd(done, cmd)
    if quit {
      break
    }
  }
}

