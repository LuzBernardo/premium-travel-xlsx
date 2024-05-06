package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Test1() {
	/* data := os.Stdin */

	scanner := bufio.NewScanner(strings.NewReader(data))
	lines := make(map[int]string)
	in := make([]string, 0)

	/* 	ioReader := io.Reader(data) */
	/* ioReadAll, _ := io.ReadAll(data)
	fmt.Printf("\nall \n%+v\n", string(ioReadAll)) */

	buf := make([]byte, 10)
	scanner.Buffer(buf, 20)

	lineNumber := 0

	for scanner.Scan() {
		text := scanner.Text()
		lines[lineNumber] = text
		in = append(in, text)
		lineNumber++
	}
	fmt.Printf("\ndata %+v\n", data)
	fmt.Printf("\n%+v\n", lines)
	fmt.Printf("\n%+v\n", in)
	/* fmt.Printf("%+v", scanner) */
}
