package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sleeyax/aternos-discord-bot/message"
	"log"
)

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func respondWithText(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return respond(s, i, &discordgo.InteractionResponseData{
		Content: content,
	})
}

func respondWithHiddenText(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return respond(s, i, &discordgo.InteractionResponseData{
		Content: message.FormatDefault(content),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

func respondWithEmbeds(s *discordgo.Session, i *discordgo.InteractionCreate, embeds []*discordgo.MessageEmbed) error {
	return respond(s, i, &discordgo.InteractionResponseData{
		Embeds: embeds,
	})
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, content string, err error) error {
	log.Printf("%s: %s\n", content, err)
	return respondWithText(s, i, message.FormatError(content))
}

func optionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}
