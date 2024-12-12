package prompt

import (
	"fmt"
)

const prompt string = "Beantworte die folgende Frage in einem kurzen Stichpunkt.\n\n"

func New(question string) string {
	return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", prompt+question)
}
