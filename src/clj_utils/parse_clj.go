package clj_utils

import (
  "errors"
  u "utils"
  s "strings"
)

var openChars = "([{"
var closedChars = ")]}"
var spaceChars = " \n"
var specialChars = openChars + closedChars + spaceChars

func getSpace(data []byte) []byte {
  for i, b := range data {
    if !s.Contains(spaceChars, string(b)) {
      return data[i:]
    }
  }
  return []byte{}
}

func getSymbol(data []byte) (string, []byte) {
  var length = 0
  for _,b := range data {
    if s.Contains(specialChars, string(b)) {
      break
    }
    length += 1
  }
  return string(data[:length]), data[length:]
}

func getTokenString(data []byte) (string, []byte) {
  if len(data) == 0 {
    return  "", data
  } else if c := string(data[0]); s.Contains(specialChars, c) {
    return c, data[1:]
  }
  return getSymbol(data)
}

func handleToken(tree *u.Tree, token *Token) error {
  tree.AppendChild(token)
  var err error
  if token.Class == "symbol" {
    // pass
  } else if token.Open {
    err = tree.DownLast()
  } else if token.Closed {
    err = tree.Up()
  } else {
    err = errors.New("Non-symbol token must be open or closed")
  }
  return err
}

func ParseClj(data []byte) *u.Tree {
  initToken := NewToken("root")
  var tree = u.NewTree(initToken)
  var s string
  var t *Token
  var l int = len(data)
  var err error
  for l > 0 {
    data = getSpace(data)
    s, data = getTokenString(data)
    if len(data) == l {
      panic("Data must get shorter " + string(data))
    }
    l = len(data)
    t = NewToken(s)
    err = handleToken(tree, t)
    if err != nil {
      panic(err)
    }
  }
  return tree
}


