package util

import (
	"github.com/lirm/aeron-go/aeron"
)

func RetryPublicationResult(res int64) bool {
	switch res {
	case aeron.AdminAction, aeron.BackPressured:
		return true
	case aeron.NotConnected, aeron.MaxPositionExceeded, aeron.PublicationClosed:
		return false
	}
	return false
}

func PublicationErrorString(res int64) string {
	switch res {
	case aeron.AdminAction:
		return "ADMIN_ACTION"
	case aeron.BackPressured:
		return "BACK_PRESSURED"
	case aeron.PublicationClosed:
		return "CLOSED"
	case aeron.NotConnected:
		return "NOT_CONNECTED"
	case aeron.MaxPositionExceeded:
		return "MAX_POSITION_EXCEEDED"
	default:
		return "UNKNOWN"
	}
}
