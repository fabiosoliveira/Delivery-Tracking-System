package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	name := "Fabio"
	name2 := "Fábio"

	fmt.Println(utf8.RuneCountInString(name))
	fmt.Println(utf8.RuneCountInString(name2))
}
