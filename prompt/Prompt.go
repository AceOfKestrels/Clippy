package prompt

import (
	"encoding/json"
	"fmt"
	"strings"
)

const prompt string = "Beantworte die folgende Frage in einem kurzen Stichpunkt.\n\n"

func New(question string) string {
	return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", escapeJsonCharacters(prompt+question))
}

func escapeJsonCharacters(str string) string {
	escaped, _ := json.Marshal(str)
	escapedString, _ := strings.CutPrefix(string(escaped), "\"")
	escapedString, _ = strings.CutSuffix(escapedString, "\"")
	return escapedString
}
