package main

type StackNode struct {
  Data interface{}
  Next *StackNode
}

type Stack struct {
  Top *StackNode
}

func NewStack () *Stack {
  return &Stack{nil}
}

func (s *Stack) Empty () bool {
  return s.Top == nil
}

func (s *Stack) Push (d interface{}) {
  var n = &StackNode{d, s.Top}
  s.Top = n
}

func (s *Stack) Pop () interface{} {
  var d = s.Top.Data
  s.Top = s.Top.Next
  return d
}

