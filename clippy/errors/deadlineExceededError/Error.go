package deadlineExceededError

import "strings"

func Is(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded")
}
