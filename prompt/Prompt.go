package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

const promptMultipleCoice string = "Welche Antwort ist die richtige? Gib nur die richtige Zahl bzw den richtigen Buchstaben."
const promptOpen string = "Beantworte die folgende Frage in einem kurzen Stichpunkt."

func New(question string) string {
	if IsMultipleChoice(question) {
		return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", promptMultipleCoice+"\n\n"+question)
	}
	return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", promptOpen+"\n\n"+question)

}

func IsMultipleChoice(question string) bool {
	// Check if it ends with a question mark
	if !strings.Contains(question, "?") {
		return false
	}

	// Regular expression to find options like A., B., 1., etc.
	optionPattern := regexp.MustCompile(`(?m)^\s*([A-Z]\.|[1-9]\.|\d\))\s+`)
	options := optionPattern.FindAllStringSubmatch(question, -1)

	// Multiple choice questions typically have at least 2 options
	return len(options) >= 2
}
