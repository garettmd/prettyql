//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/prometheus/prometheus/promql/parser"
	"strings"
)

func formatPromQL(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return map[string]interface{}{
			"error": "no query provided",
		}
	}

	input := args[0].String()
	smushed := strings.Join(strings.Fields(input), " ")

	expr, err := parser.ParseExpr(smushed)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"formatted": expr.Pretty(0),
	}
}

func main() {
	js.Global().Set("formatPromQL", js.FuncOf(formatPromQL))

	// Block forever so the Go runtime stays alive
	select {}
}
