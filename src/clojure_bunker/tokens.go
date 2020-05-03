package main

type Token struct {
  Class string
  Value string
  Selected bool
  Row int
  Col int
}

func NewToken (class, value string) *Token {
  return &Token{class, value, false, -1, -1}
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
