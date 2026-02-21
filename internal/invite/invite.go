package invite

import (
	"fmt"

	"github.com/Work-Fort/Discord/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Run creates or retrieves a permanent server invite link
func Run(cfg *config.Config) error {
	// Create Discord session
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return fmt.Errorf("creating Discord session: %w", err)
	}

	// Open connection
	if err := session.Open(); err != nil {
		return fmt.Errorf("opening Discord connection: %w", err)
	}
	defer session.Close()

	fmt.Println("Connected to Discord")

	// Get existing invites
	invites, err := session.GuildInvites(cfg.GuildID)
	if err != nil {
		return fmt.Errorf("fetching guild invites: %w", err)
	}

	// Look for an existing permanent invite (maxAge=0, maxUses=0)
	for _, invite := range invites {
		if invite.MaxAge == 0 && invite.MaxUses == 0 {
			fmt.Println("Found existing permanent invite:")
			fmt.Printf("  https://discord.gg/%s\n", invite.Code)
			fmt.Printf("  Uses: %d\n", invite.Uses)
			fmt.Printf("  Created: %s\n", invite.CreatedAt)
			return nil
		}
	}

	// No permanent invite found, create one
	fmt.Println("No permanent invite found, creating one...")

	// Get the first text channel to create invite for
	channels, err := session.GuildChannels(cfg.GuildID)
	if err != nil {
		return fmt.Errorf("fetching guild channels: %w", err)
	}

	var targetChannelID string
	for _, ch := range channels {
		// Find first text channel (prefer #general or #welcome)
		if ch.Type == discordgo.ChannelTypeGuildText {
			if ch.Name == "general" || ch.Name == "welcome" {
				targetChannelID = ch.ID
				break
			}
			if targetChannelID == "" {
				targetChannelID = ch.ID
			}
		}
	}

	if targetChannelID == "" {
		return fmt.Errorf("no text channels found to create invite for")
	}

	// Create permanent invite (maxAge=0, maxUses=0, temporary=false)
	invite, err := session.ChannelInviteCreate(targetChannelID, discordgo.Invite{
		MaxAge:    0,     // Never expires
		MaxUses:   0,     // Unlimited uses
		Temporary: false, // Not temporary membership
		Unique:    true,  // Create a new unique code
	})
	if err != nil {
		return fmt.Errorf("creating invite: %w", err)
	}

	fmt.Println("Created new permanent invite:")
	fmt.Printf("  https://discord.gg/%s\n", invite.Code)
	fmt.Println("\nShare this link to invite people to your Discord server!")

	return nil
}
