//go:build js && wasm

package main

import (
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
	"syscall/js"
)

func formatPromQL(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return map[string]interface{}{
			"error": newAppError("no query provided", "Enter a PromQL query before clicking Format."),
		}
	}

	input := args[0].String()
	smushed := strings.Join(strings.Fields(input), " ")

	expr, err := parser.ParseExpr(smushed)
	if err != nil {
		return map[string]interface{}{
			"error": newPrettyQLError(err.Error()),
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
