package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"prettyql/web"
)

type formatRequest struct {
	Query string `json:"query"`
}

type formatResponse struct {
	Formatted string `json:"formatted,omitempty"`
	Error     *FormatError `json:"error,omitempty"`
}

func handleFormat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req formatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(formatResponse{Error: newAppError("invalid JSON request", "Send a POST body with a JSON object like {\"query\":\"...\"}.")})
		return
	}

	if strings.TrimSpace(req.Query) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(formatResponse{Error: newAppError("query cannot be empty", "Paste a PromQL query before formatting it.")})
		return
	}

	formatted, err := formatPromQL(req.Query)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if fe, ok := err.(*FormatError); ok {
			json.NewEncoder(w).Encode(formatResponse{Error: fe})
		} else {
			json.NewEncoder(w).Encode(formatResponse{Error: newAppError(err.Error(), "Try again or check the server logs for more details.")})
		}
		return
	}

	json.NewEncoder(w).Encode(formatResponse{Formatted: formatted})
}

func runServer(port int) {
	staticFiles, err := fs.Sub(web.StaticFiles, ".")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(staticFiles)))
	http.HandleFunc("/api/format", handleFormat)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("PrettyQL server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func runCLI() {
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

	input := strings.Join(lines, " ")
	result, err := formatPromQL(input)
	if err != nil {
		if fe, ok := err.(*FormatError); ok {
			fmt.Fprintf(os.Stderr, "Error parsing PromQL query: %s\n", fe.Message)
			if fe.Suggestion != "" {
				fmt.Fprintf(os.Stderr, "Suggestion: %s\n", fe.Suggestion)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error parsing PromQL query: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println(result)
}

func main() {
	serve := flag.Bool("serve", false, "Start web server instead of reading from stdin")
	port := flag.Int("port", 8080, "Port to listen on (used with --serve)")
	flag.Parse()

	if *serve {
		runServer(*port)
	} else {
		runCLI()
	}
}
