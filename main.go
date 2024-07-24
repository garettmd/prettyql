package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	expr, err := parser.ParseExpr(scanner.Text())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PromQL query: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(expr.Pretty(0))
}
