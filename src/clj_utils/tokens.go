package clj_utils

type Token struct {
  Class string
  Open bool
  Closed bool
  Value string
  Row int
  Col int
}

func NewToken (s string) *Token {
  switch s {
  case "(": return &Token{"call", true, false, "(", -1, -1}
  case ")": return &Token{"call", false, true, ")", -1, -1}
  case "[": return &Token{"vect", true, false, "[", -1, -1}
  case "]": return &Token{"vect", false, true, "]", -1, -1}
  case "{": return &Token{"map", true, false, "{", -1, -1}
  case "}": return &Token{"map", false, true, "}", -1, -1}
  default: return &Token{"symbol", false, false, s, -1, -1}
  }
}

