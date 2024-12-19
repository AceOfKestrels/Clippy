package contextCanceledError

import "strings"

func Is(err error) bool {
	return strings.Contains(err.Error(), "context canceled")
}
