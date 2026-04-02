package main

import (
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
)

// formatPromQL takes a raw PromQL query string and returns the prettified version.
func formatPromQL(input string) (string, error) {
	smushed := strings.Join(strings.Fields(input), " ")
	expr, err := parser.ParseExpr(smushed)
	if err != nil {
		return "", err
	}
	return expr.Pretty(0), nil
}
