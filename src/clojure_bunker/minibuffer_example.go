package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	reader := bufio.NewReader(os.Stdin)

        s := readInputInMiniBuffer()
        fmt.Println(s)

	var c rune
	for err == nil {
		if c == 'q' {
			break
		}

		c, _, err = reader.ReadRune()

		if unicode.IsControl(c) {
			fmt.Printf("%d\r\n", c)
		} else {
			fmt.Printf("%d (%c)\r\n", c, c)
		}
	}

	if err != nil {
		panic(err)
	}
}
