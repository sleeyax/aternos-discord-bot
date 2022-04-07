package aternos_discord_bot

import "strings"

type MessageType int

const (
	success MessageType = iota
	warning
	info
	danger
	loading
	normal
)

func formatMessage(text string, messageType MessageType) string {
	b := strings.Builder{}

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
