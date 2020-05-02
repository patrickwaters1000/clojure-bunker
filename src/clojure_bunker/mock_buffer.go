package main

type MockBuffer struct {
  name string
}

func NewBuffer(name string) *MockBuffer {
  return &MockBuffer{name}
}

func (b *MockBuffer) handleEvent (ev []string) error {
  return nil
}

func (b MockBuffer) render () {
  tbPrint(0, 0, fg1, bg1, b.name)
}
