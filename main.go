package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	var lines []string

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return
		}
	}

	// Smush multiline strings into one line. The promql parser doesn't handle multiline strings.very well.
	smushed := strings.Join(lines, " ")

	expr, err := parser.ParseExpr(smushed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PromQL query: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(expr.Pretty(0))
}
