package main

import (
  termbox "github.com/nsf/termbox-go"
)

type Token struct {
  Class string
  Value string
  Selected bool
  NewLine bool
  Row int
  Col int
  Color termbox.Attribute
  Style string
}

func NewToken (class, value string) *Token {
  return &Token{
    Class: class,
    Value: value,
    Selected: false,
    NewLine: false, // used only with custom formatting
    Row: -1,
    Col: -1,
    Color: termbox.ColorWhite,
    Style: "",
  }
}

func (t Token) IsOpen() bool {
  switch t.Value {
  case "(": return true
  case "[": return true
  case "{": return true
  default: return false
  }
}

func (t Token) IsClosed() bool {
  switch t.Value {
  case ")": return true
  case "]": return true
  case "}": return true
  default: return false
  }
}
