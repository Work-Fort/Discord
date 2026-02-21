package setup

import (
	"fmt"

	"github.com/Work-Fort/Discord/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Run performs initial Discord server setup
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

	// Setup roles first (they're referenced in channel permissions)
	if err := setupRoles(session, cfg); err != nil {
		return fmt.Errorf("setting up roles: %w", err)
	}

	// Setup channels and categories
	if err := setupChannels(session, cfg); err != nil {
		return fmt.Errorf("setting up channels: %w", err)
	}

	// Setup integrations (webhooks)
	if err := setupIntegrations(session, cfg); err != nil {
		return fmt.Errorf("setting up integrations: %w", err)
	}

	return nil
}

func setupRoles(session *discordgo.Session, cfg *config.Config) error {
	fmt.Println("Setting up roles...")

	existingRoles, err := session.GuildRoles(cfg.GuildID)
	if err != nil {
		return fmt.Errorf("fetching existing roles: %w", err)
	}

	// Build map of existing roles
	roleMap := make(map[string]*discordgo.Role)
	for _, role := range existingRoles {
		roleMap[role.Name] = role
	}

	for _, roleCfg := range cfg.Roles.Roles {
		if _, exists := roleMap[roleCfg.Name]; exists {
			fmt.Printf("  ⊙ Role already exists: %s\n", roleCfg.Name)
			continue
		}

		// Convert color hex to int
		var color int
		fmt.Sscanf(roleCfg.Color, "#%x", &color)

		// Convert permission strings to int64
		permissions := int64(0)
		for _, perm := range roleCfg.Permissions {
			permissions |= permissionValue(perm)
		}

		params := &discordgo.RoleParams{
			Name:        roleCfg.Name,
			Color:       &color,
			Permissions: &permissions,
			Hoist:       &roleCfg.Hoist,
			Mentionable: &roleCfg.Mentionable,
		}

		_, err := session.GuildRoleCreate(cfg.GuildID, params)
		if err != nil {
			return fmt.Errorf("creating role %s: %w", roleCfg.Name, err)
		}

		fmt.Printf("  ✓ Created role: %s\n", roleCfg.Name)
	}

	return nil
}

func setupChannels(session *discordgo.Session, cfg *config.Config) error {
	fmt.Println("Setting up channels...")

	for _, category := range cfg.Channels.Categories {
		// Create category
		categoryChannel, err := session.GuildChannelCreateComplex(cfg.GuildID, discordgo.GuildChannelCreateData{
			Name:     category.Name,
			Type:     discordgo.ChannelTypeGuildCategory,
			Position: category.Position,
		})
		if err != nil {
			return fmt.Errorf("creating category %s: %w", category.Name, err)
		}

		fmt.Printf("  ✓ Created category: %s\n", category.Name)

		// Create channels in category
		for _, ch := range category.Channels {
			channelType := discordgo.ChannelTypeGuildText
			if ch.Type == "voice" {
				channelType = discordgo.ChannelTypeGuildVoice
			} else if ch.Type == "forum" {
				channelType = discordgo.ChannelTypeGuildForum
			}

			channelData := discordgo.GuildChannelCreateData{
				Name:     ch.Name,
				Type:     channelType,
				Topic:    ch.Topic,
				Position: ch.Position,
				ParentID: categoryChannel.ID,
			}

			channel, err := session.GuildChannelCreateComplex(cfg.GuildID, channelData)
			if err != nil {
				return fmt.Errorf("creating channel %s: %w", ch.Name, err)
			}

			// Apply channel-specific permissions
			if ch.Permissions != nil {
				if err := applyChannelPermissions(session, cfg.GuildID, channel.ID, ch.Permissions); err != nil {
					return fmt.Errorf("applying permissions to %s: %w", ch.Name, err)
				}
			}

			// Add forum tags if it's a forum channel
			if ch.Type == "forum" && len(ch.Tags) > 0 {
				if err := addForumTags(session, channel.ID, ch.Tags); err != nil {
					return fmt.Errorf("adding tags to forum %s: %w", ch.Name, err)
				}
			}

			fmt.Printf("    ✓ Created channel: %s\n", ch.Name)
		}
	}

	return nil
}

func setupIntegrations(session *discordgo.Session, cfg *config.Config) error {
	fmt.Println("Setting up integrations...")

	if cfg.Integrations.GitHub != nil && cfg.Integrations.GitHub.Enabled {
		// Find the target channel
		channels, err := session.GuildChannels(cfg.GuildID)
		if err != nil {
			return fmt.Errorf("fetching channels: %w", err)
		}

		var targetChannelID string
		for _, ch := range channels {
			if ch.Name == cfg.Integrations.GitHub.TargetChannel {
				targetChannelID = ch.ID
				break
			}
		}

		if targetChannelID == "" {
			return fmt.Errorf("target channel not found: %s", cfg.Integrations.GitHub.TargetChannel)
		}

		// Create webhook for GitHub
		webhook, err := session.WebhookCreate(targetChannelID, "GitHub", "")
		if err != nil {
			return fmt.Errorf("creating GitHub webhook: %w", err)
		}

		fmt.Printf("  ✓ Created GitHub webhook\n")
		fmt.Printf("    Add this URL to GitHub repo webhooks:\n")
		fmt.Printf("    https://discord.com/api/webhooks/%s/%s/github\n", webhook.ID, webhook.Token)
	}

	return nil
}

func applyChannelPermissions(session *discordgo.Session, guildID, channelID string, perms map[string]map[string]bool) error {
	// Get @everyone role ID (it's the same as guild ID)
	everyonePerms, hasEveryone := perms["everyone"]
	if !hasEveryone {
		return nil
	}

	// Calculate permission overwrite
	allow := int64(0)
	deny := int64(0)

	for perm, value := range everyonePerms {
		permValue := permissionValue(perm)
		if value {
			allow |= permValue
		} else {
			deny |= permValue
		}
	}

	err := session.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, allow, deny)
	if err != nil {
		return fmt.Errorf("setting @everyone permissions: %w", err)
	}

	return nil
}

func addForumTags(session *discordgo.Session, channelID string, tags []config.ForumTag) error {
	// Discord API for forum tags requires channel edit
	// This is a simplified implementation - full implementation would use ChannelEdit
	// with AvailableTags field
	return nil
}

func permissionValue(name string) int64 {
	perms := map[string]int64{
		"administrator":        discordgo.PermissionAdministrator,
		"send_messages":        discordgo.PermissionSendMessages,
		"embed_links":          discordgo.PermissionEmbedLinks,
		"attach_files":         discordgo.PermissionAttachFiles,
		"read_message_history": discordgo.PermissionReadMessageHistory,
		"use_external_emojis":  discordgo.PermissionUseExternalEmojis,
		"add_reactions":        discordgo.PermissionAddReactions,
		"manage_messages":      discordgo.PermissionManageMessages,
	}

	if val, ok := perms[name]; ok {
		return val
	}

	return 0
}
