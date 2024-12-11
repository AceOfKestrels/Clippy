package prompt

import (
	"fmt"
)

const prompt string = "Welche Antwort ist die richtige? Gib nur die richtige Zahl bzw den richtigen Buchstaben.\n"

func New(question string) string {
	return fmt.Sprintf("{ \"contents\": [ { \"parts\": [ {\"text\": \"%v\" } ] } ] }", prompt+"\n\n"+question)
}
