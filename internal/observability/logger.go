package observability

import "log"

// Printf proxies formatted logs through the standard logger.
func Printf(format string, args ...any) {
	log.Printf(format, args...)
}
