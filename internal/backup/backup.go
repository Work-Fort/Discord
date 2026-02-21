package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Work-Fort/Discord/internal/config"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

// Run exports current Discord server state to YAML files
func Run(cfg *config.Config) error {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return fmt.Errorf("creating Discord session: %w", err)
	}

	if err := session.Open(); err != nil {
		return fmt.Errorf("opening Discord connection: %w", err)
	}
	defer session.Close()

	fmt.Println("Connected to Discord")
	fmt.Println("Exporting server state...")

	// Create backup directory with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join("backups", timestamp)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("creating backup directory: %w", err)
	}

	// Export roles
	if err := exportRoles(session, cfg.GuildID, backupDir); err != nil {
		return fmt.Errorf("exporting roles: %w", err)
	}

	// Export channels
	if err := exportChannels(session, cfg.GuildID, backupDir); err != nil {
		return fmt.Errorf("exporting channels: %w", err)
	}

	fmt.Printf("✓ Backup saved to: %s\n", backupDir)

	return nil
}

func exportRoles(session *discordgo.Session, guildID, backupDir string) error {
	roles, err := session.GuildRoles(guildID)
	if err != nil {
		return fmt.Errorf("fetching roles: %w", err)
	}

	rolesConfig := config.RolesConfig{
		Roles: make([]config.Role, 0),
	}

	for _, role := range roles {
		// Skip @everyone role
		if role.Name == "@everyone" {
			continue
		}

		rolesConfig.Roles = append(rolesConfig.Roles, config.Role{
			Name:        role.Name,
			Color:       fmt.Sprintf("#%06x", role.Color),
			Permissions: []string{}, // TODO: Convert int64 permissions back to names
			Hoist:       role.Hoist,
			Mentionable: role.Mentionable,
		})
	}

	data, err := yaml.Marshal(rolesConfig)
	if err != nil {
		return fmt.Errorf("marshaling roles: %w", err)
	}

	path := filepath.Join(backupDir, "roles.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing roles backup: %w", err)
	}

	fmt.Printf("  ✓ Exported roles (%d)\n", len(rolesConfig.Roles))

	return nil
}

func exportChannels(session *discordgo.Session, guildID, backupDir string) error {
	channels, err := session.GuildChannels(guildID)
	if err != nil {
		return fmt.Errorf("fetching channels: %w", err)
	}

	channelsConfig := config.ChannelsConfig{
		Categories: make([]config.Category, 0),
	}

	// Build category map
	categoryMap := make(map[string]*config.Category)

	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory {
			category := config.Category{
				Name:     ch.Name,
				Position: ch.Position,
				Channels: make([]config.Channel, 0),
			}
			categoryMap[ch.ID] = &category
			channelsConfig.Categories = append(channelsConfig.Categories, category)
		}
	}

	// Add channels to categories
	for _, ch := range channels {
		if ch.ParentID != "" {
			if category, ok := categoryMap[ch.ParentID]; ok {
				channelType := "text"
				if ch.Type == discordgo.ChannelTypeGuildVoice {
					channelType = "voice"
				} else if ch.Type == discordgo.ChannelTypeGuildForum {
					channelType = "forum"
				}

				channel := config.Channel{
					Name:     ch.Name,
					Type:     channelType,
					Topic:    ch.Topic,
					Position: ch.Position,
				}

				category.Channels = append(category.Channels, channel)
			}
		}
	}

	data, err := yaml.Marshal(channelsConfig)
	if err != nil {
		return fmt.Errorf("marshaling channels: %w", err)
	}

	path := filepath.Join(backupDir, "channels.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing channels backup: %w", err)
	}

	fmt.Printf("  ✓ Exported channels (%d categories)\n", len(channelsConfig.Categories))

	return nil
}
