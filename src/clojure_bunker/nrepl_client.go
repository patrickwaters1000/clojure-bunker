package main

import (
  "fmt"
  "net"
)

type List struct {
  vals []interface{}
}

func (l *List) conj (x interface{}) {
  l.vals = append(l.vals, x)
}

type Dict struct {
  keys []interface{}
  vals []interface{}
}

func (d *Dict) assoc (k, v interface{}) {
  d.keys = append(d.keys, k)
  d.vals = append(d.vals, v)
}

func (d *Dict) get (k_want interface{}) (interface{}, bool) {
  for i:=0; i<len(d.keys); i++ {
    k_got := d.keys[i]
    switch k_want.(type) {
    case string:
      switch k_got.(type) {
      case string:
        if k_got.(string) == k_want.(string) {
          return d.vals[i], true
        }
      }
    case int:
      switch k_got.(type) {
      case int:
        if k_got.(int) == k_want.(int) {
          return d.vals[i], true
        }
      }
    }
  }
  return nil, false
}

func isDigit(b byte) bool {
  return 48 <= b && b < 58
}

func toDigit (b byte) int {
  if !isDigit(b) {
    panic("Not a digit")
  }
  return int(b) - 48
}

func readBEncode(in chan byte) interface{} {
  b := <-in
  r := rune(b)
  if r == 'e' { // Only possible during recursive calls
    return nil
  } else if  r == 'i' {
    return readBEncodeInt(in)
  } else if isDigit(b) {
    return readBEncodeString(in, b)
  } else if r == 'l' {
    return readBEncodeList(in)
  } else if r == 'd' {
    return readBEncodeDict(in)
  } else {
    panic("Unexpected input")
  }
}

//func (bp *bEncodeParser) parseInt () int {
//  x := 1
//  b := <-bp.in
//  for rune(b) != 'e' {
//    switch rune(b) {
//    case '-': x *= -1
//    default: x = 10*x + toDigit(b)
//    }
//    b = <-bp.in
//  }
//  return x
//}

func readBEncodeInt (in chan byte) int {
  x := 1
  b := <-in
  if rune(b) == '-' {
    x = -1
  } else {
    x = toDigit(b)
  }
  for b := range in {
    if rune(b) == 'e' {
      break
    } else {
      x = 10 * x + toDigit(b)
    }
  }
  return x
}

func readBEncodeString (in chan byte, b byte) string {
  strlen := toDigit(b)
  for b := range in {
    if isDigit(b) {
      strlen = 10 * strlen + toDigit(b)
    } else if rune(b) == ':' {
      break
    } else {
      panic("Failed to parse string")
    }
  }
  buffer := []byte{}
  for i:=0; i<strlen; i++ {
    b := <-in
    buffer = append(buffer, b)
  }
  return string(buffer)
}

func readBEncodeList (in chan byte) *List {
  l := &List{[]interface{}{}}
  for {
    x := readBEncode(in)
    if x == nil {
      break
    } else {
      l.conj(x)
    }
  }
  return l
}

func readBEncodeDict (in chan byte) *Dict {
  d := &Dict{[]interface{}{}, []interface{}{}}
  for {
    k := readBEncode(in)
    if k == nil {
      break
    } else {
      v := readBEncode(in)
      d.assoc(k,v)
    }
  }
  return d
}

func pprints (d interface{}, newLine string) string {
  spaces := "    "
  switch d.(type) {
  case int:
    return fmt.Sprintf("%d", d.(int))
  case string:
    return d.(string)
  case *List:
    buffer := "[" + newLine + spaces
    for x := range d.(*List).vals {
      x_str := pprints(x, newLine + spaces)
      buffer += x_str + "," + newLine
    }
    return buffer + "]"
  case *Dict:
    buffer := "{" + newLine + spaces
    keys := d.(*Dict).keys
    vals := d.(*Dict).vals
    for i:=0; i<len(keys); i++ {
      k_str := pprints(keys[i], newLine + spaces)
      v_str := pprints(vals[i], newLine + spaces)
      buffer += k_str + ": " + v_str + "," + newLine
      if i<len(keys)-1 {
        buffer += spaces
      }
    }
    return buffer + "}"
  default:
    panic("Can't print input type")
  }
}

func writeBEncode(data interface{}) string {
  switch data.(type) {
  case int:
    return fmt.Sprintf("i%de", data.(int))
  case string:
    return fmt.Sprintf("%d:%s", len(data.(string)), data.(string))
  case *List:
    buffer := "l"
    for x := range data.(*List).vals {
      buffer += writeBEncode(x)
    }
    buffer += "e"
    return buffer
  case *Dict:
    buffer := "d"
    keys := data.(*Dict).keys
    vals := data.(*Dict).vals
    for i:=0; i<len(keys); i++ {
      buffer += writeBEncode(keys[i])
      buffer += writeBEncode(vals[i])
    }
    buffer += "e"
    return buffer
  default:
    panic("Failed to encode")
  }
}

type Client struct {
  in chan byte
  conn net.Conn
  session string
  nextId int
  results map[int]string
  failures map[int]string
}

func NewClient () *Client {
    return &Client{
    in: nil,
    conn: nil,
    session: "1234", // Maybe use a random number instead?
    nextId: 0,
    results: make(map[int]string),
    failures: make(map[int]string),
  }
}

func (c *Client) Connect (port string) {
  conn, err := net.Dial("tcp", "localhost:" + port)
  c.conn = conn
  if err != nil {
    panic(err)
  }
  c.in = make(chan byte)
  buffer := make([]byte, 1)
  go func() {
    for {
      n, err := c.conn.Read(buffer)
      if err != nil {
        panic(err)
      }
      if n != 1 {
        panic("Supposed to read 1 byte")
      }
      c.in <- buffer[0]
    }
  }()
}

func (c *Client) Close () {
  c.conn.Close()
}

func (c *Client) Send (code string) {
  msg := writeBEncode(
    &Dict{
      []interface{}{"code", "id", "op"},
      []interface{}{code, c.nextId, "eval"}})
  c.nextId++
  msgBytes := []byte(msg)
  n, err := c.conn.Write(msgBytes)
  if err != nil {
    panic(err)
  }
  if n != len(msgBytes) {
    panic(fmt.Sprintf(
      "Expected to write %d bytes, but only wrote %d.",
      len(msgBytes),
      n))
  }
}

func (c *Client) Receive () {
  msg := readBEncode(c.in).(*Dict)
  id, hasId := msg.get("id")
  idInt := id.(int)
  if hasId {
    value, hasValue := msg.get("value")
    err, hasErr := msg.get("err")
    if hasValue {
      valueStr := value.(string)
      c.results[idInt] = valueStr
    } else if hasErr {
      errStr := err.(string)
      c.failures[idInt] = errStr
    }
  }
}

func (c *Client) GetResponse (id int) (string, bool) {
  for {
    c.Receive()
    value, hasValue := c.results[id]
    err, hasErr := c.failures[id]
    if hasValue {
      return value, true
    } else if hasErr {
      return err, false
    }
  }
  return "", false // Do we need this?
}

//func main() {
//  port := os.Args[1]
//  session := "12345"
//  client := NewClient(port, session)
//  reader := bufio.NewReader(os.Stdin)
//  i := 0
//  for {
//    fmt.Print(">")
//    text, _ := reader.ReadString('\n')
//    client.Send(text)
//    resp, _ := client.GetResponse(i)
//    fmt.Println(resp)
//    i++
//  }
//}
