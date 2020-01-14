package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("asfgfasadfgfdas")
	for {
		fmt.Print("> ")
		rd := bufio.NewReader(os.Stdin)
		s, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("bye")
				return
			}
		}
		fmt.Println(s[:len(s)-1])
	}
}
