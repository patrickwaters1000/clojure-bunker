package main

import (
  s "strings"
)

var openChars = "([{"
var closedChars = ")]}"
var spaceChars = " \t\n"+string(rune(10)) // Getting "fail" token at end of file. Use byte slices instead?
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

func getToken(data []byte) (*Token, []byte) {
  if len(data) == 0 {
    return  NewToken("fail",""), data
  } else if c := string(data[0]); s.Contains(closedChars, c) {
    return NewToken("close", c), data[1:]
  } else if c := string(data[0]); s.Contains(openChars, c) {
    return NewToken("open", c), data[1:]
  } else {
    symbol, newData := getSymbol(data)
    return NewToken("symbol", symbol), newData
  }
}

func handleToken(tree *Tree, token *Token) error {
  tree.AppendChild(token)
  var err error
  if token.Class == "open" {
    err = tree.DownLast()
  } else if token.Class == "close" {
    err = tree.Up()
  }
  return err
}

func parseClj(data []byte) *Tree {
  initToken := NewToken("root","")
  var tree = NewTree(initToken)
  var token *Token
  var l int = len(data)
  var err error
  for l > 0 {
    data = getSpace(data)
    token, data = getToken(data)
    if len(data) == l {
      panic("Data must get shorter " + string(data))
    }
    l = len(data)
    err = handleToken(tree, token)
    panicIfError(err)
  }
  leafToken := NewToken("leaf", "")
  tree.AppendChild(leafToken)
  return tree
}


