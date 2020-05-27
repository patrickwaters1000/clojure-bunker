package main

import (
  s "strings"
)

var openChars = "([{"
var closedChars = ")]}"
var spaceChars = " \t\n"+string(byte(10)) // Getting "fail" token at end of file. Use byte slices instead?
var specialChars = openChars + closedChars + spaceChars

// Reads from data until non-whitespace character is found.
// First return value is whether a newline was read.
// Second return value is the remaining data.
func getSpace(data []byte) (bool, []byte) {
  newLine := false
  for i, b := range data {
    if string(b) == "\n" {
      newLine = true
    }
    if !s.Contains(spaceChars, string(b)) {
      return newLine, data[i:]
    }
  }
  return newLine, []byte{}
}

func getComment(data []byte) (string, []byte) {
  var length = 0
  for _,b := range data {
    if rune(b) == '\n' {
      break
    }
    length += 1
  }
  return string(data[:length]), data[length:]
}

func getString(data []byte) (string, []byte) {
  var length = 0
  for _,b := range data[1:] {
    if rune(b) == '"' {
      break
    }
    length += 1
  }
  return string(data[:length + 2]), data[length + 2:]
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
  } else if c := rune(data[0]); c == ';' {
    value, newData := getComment(data)
    return NewToken("comment", value), newData
  } else if c := rune(data[0]); c == '"' {
    value, newData := getString(data)
    return NewToken("string", value), newData
  } else {
    value, newData := getSymbol(data)
    tokenClass := "symbol"
    if value[0] == ':' {
      tokenClass = "keyword"
    }
    return NewToken(tokenClass, value), newData
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
  initToken := NewToken("root","normal")
  var tree = NewTree(initToken)
  var token *Token
  var l int = len(data)
  var err error
  for l > 0 {
    var newLine bool
    newLine, data = getSpace(data)
    token, data = getToken(data)
    if newLine {
      token.NewLine = true
    }
    if token.IsOpen() {
      token.Style = "custom"
    }
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

// type Parser struct {
//   tree Tree
//   data []byte
// }
// 
// func (p *Parser) TakeToken () {
// 
// 
// func (p *Parser) ParseToken (parent string, leader, string, idx int) {
//   newLine := p.TakeWhitespace()
//   t := p.TakeToken()
//   if newLine {
//     t.Newline = true
//   }
//   if t.IsOpen() {
//     leader := p.CheckToken() {
//     case "defn": p.ParseDefn()
//     case "let": p.ParseLet()
//     default: p.ParseCall()
//     }
//   } else if t.IsClosed() {
//     p.ParseClose()
//   } else {
//     p.tree.AppendChild(t)
//   }
// 
