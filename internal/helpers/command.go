package helpers

import "strings"

func GetCommandArgs(input string) (string, bool) {
	if input == "" {
		return "", false
	}

	parts := strings.SplitN(input, " ", 2)
	hasArgs := len(parts) > 1 && strings.TrimSpace(parts[1]) != ""

	if !hasArgs {
		return "", false
	}

	return strings.TrimSpace(parts[1]), true
}

func ParseCommand(input string) (cmd string) {
	if input == "" || input[0] != '/' {
		return ""
	}

	parts := strings.SplitN(input, " ", 2)
	cmd = strings.TrimPrefix(parts[0], "/")

	return strings.ToLower(cmd)
}
