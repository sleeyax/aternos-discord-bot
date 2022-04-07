package message

import (
	"fmt"
	"strings"
)

type Type int

const (
	success Type = iota
	warning
	info
	danger
	loading
	normal
)

// Format creates a message string optionally prefixed with an emoji based on the given message type.
//
// Defaults to a plain text message.
func Format(text string, messageTypes ...Type) string {
	b := strings.Builder{}
	messageType := normal

	if len(messageTypes) > 0 {
		messageType = messageTypes[0]
	}

	switch messageType {
	case success:
		b.WriteString(":white_check_mark: ")
	case warning:
		b.WriteString(":warning: ")
	case info:
		b.WriteString(":information_source: ")
	case danger:
		b.WriteString(":x: ")
	case loading:
		b.WriteString(":hourglass: ")
	case normal:
		fallthrough
	default:
		break
	}

	b.WriteString(text)

	return b.String()
}

func FormatSuccess(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), success)
}

func FormatWarning(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), warning)
}

func FormatInfo(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), info)
}

func FormatError(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), danger)
}

func FormatLoading(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), loading)
}

func FormatDefault(format string, a ...any) string {
	return Format(fmt.Sprintf(format, a...), normal)
}
