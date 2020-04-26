package clj_utils

import (
  u "utils"
  "testing"
)

type getSymbolTestCase struct {
  input string
  wantSymbol string
  wantDataLength int
}

func TestGetSymbol(t *testing.T) {
  cases := []getSymbolTestCase{
    getSymbolTestCase{"defn", "defn", 0},
    getSymbolTestCase{"defn ", "defn", 1},
    getSymbolTestCase{"y", "y", 0},
    getSymbolTestCase{"defn)", "defn", 1},
  }
  var gotSymbol string
  var gotData []byte
  for i,c := range cases {
    gotSymbol, gotData = getSymbol([]byte(c.input))
    if gotSymbol != c.wantSymbol || len(gotData) != c.wantDataLength {
      t.Errorf(
        "Case %d, want %s and %d, got %s and %d",
        i, c.wantSymbol, c.wantDataLength, gotSymbol, len(gotData),
      )
    }
  }
}

type testCase struct {
  symbol string
  children int
}

func TestParseClj(t *testing.T) {
  clj := "(defn f [x]\n" +
         "  (let [y (* x x)]\n" +
         "    (inc y)))"
  tree := ParseClj([]byte(clj))
  var gots = []testCase{}
  tree.DepthFirstTraverse(func(n *u.TreeNode) {
    d := n.Data.(*Token)
    gots = append(gots, testCase{d.Value, len(n.Children)})
  })
  wants := []testCase{
    testCase{"root",1},
    testCase{"(",5},
    testCase{"defn",0},
    testCase{"f",0},
    testCase{"[",2},
    testCase{"x",0},
    testCase{"]",0},
    testCase{"(",4},
    testCase{"let",0},
    testCase{"[",3},
    testCase{"y",0},
    testCase{"(",4},
    testCase{"*",0},
    testCase{"x",0},
    testCase{"x",0},
    testCase{")",0},
    testCase{"]",0},
    testCase{"(",3},
    testCase{"inc",0},
    testCase{"y",0},
    testCase{")",0},
    testCase{")",0},
    testCase{")",0},
  }
  for i, want := range wants {
    got := gots[i]
    if want != got {
      t.Errorf(
        "Case %d: want %v, got %v",
        i, want, got,
      )
    }
  }
}
