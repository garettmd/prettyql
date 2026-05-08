package main

import "strings"

type FormatError struct {
	Source     string `json:"source"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

func (e *FormatError) Error() string {
	return e.Message
}

func newPrettyQLError(message string) *FormatError {
	return &FormatError{
		Source:     "prettyql",
		Message:    message,
		Suggestion: suggestFix(message),
	}
}

func newAppError(message, suggestion string) *FormatError {
	return &FormatError{
		Source:     "app",
		Message:    message,
		Suggestion: suggestion,
	}
}

func suggestFix(message string) string {
	ml := strings.ToLower(message)

	switch {
	case strings.Contains(ml, "unexpected end of input"):
		return "Check for a missing closing parenthesis, brace, or operator near the end of the query."
	case strings.Contains(ml, "unexpected") && strings.Contains(ml, "got"):
		return "Look for a typo, stray token, or missing comma in the expression."
	case strings.Contains(ml, "unknown function"):
		return "Verify the function name and argument count against PromQL syntax."
	case strings.Contains(ml, "parse") || strings.Contains(ml, "syntax"):
		return "Try formatting a smaller sub-expression and check PromQL syntax step by step."
	default:
		return "Double-check the PromQL syntax and try isolating the part of the query that fails."
	}
}
