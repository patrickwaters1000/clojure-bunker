package main

import (
  u "utils"
  "fmt"
)

type thing struct {
  value int
}

func main() {
  s := u.NewStack()
  v := s.Top.(thing)
  fmt.Printf("%v", v == nil)
}
