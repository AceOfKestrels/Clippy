package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

const promptMultipleCoice string = "Welche Antwort ist die richtige? Gib nur die richtige Zahl bzw den richtigen Buchstaben.\n\n"
const promptOpen string = "Beantworte die folgende Frage in einem kurzen Stichpunkt.\n\n"

func New(question string) string {
	var prompt string
	if isMultipleChoice(question) {
		prompt = promptMultipleCoice + question
	} else {
		prompt = promptOpen + question
	}
	return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", prompt)
}

func isMultipleChoice(question string) bool {
	if !strings.Contains(question, "?") {
		return false
	}

	optionPattern := regexp.MustCompile(`(?m)^\s*([A-Z]\.|[1-9]\.|\d\))\s+`)
	options := optionPattern.FindAllStringSubmatch(question, -1)

	return len(options) >= 2
}
