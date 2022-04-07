package aternos_discord_bot

import "strings"

type MessageType int

const (
	success MessageType = iota
	warning
	danger
	normal
	loading
)

func formatMessage(text string, messageType MessageType) string {
	b := strings.Builder{}

	switch messageType {
	case success:
		b.WriteString(":white_check_mark: ")
	case warning:
		b.WriteString(":warning: ")
	case danger:
		b.WriteString(":x: ")
	case loading:
		b.WriteString(":hourglass: ")
	case normal:
	default:
		break
	}

	b.WriteString(text)

	return b.String()
}
