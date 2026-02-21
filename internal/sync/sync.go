package sync

import (
	"fmt"

	"github.com/Work-Fort/Discord/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Run syncs configuration changes to Discord server
// This is similar to setup but handles updates to existing resources
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
	fmt.Println("Syncing configuration...")

	// For now, sync is a simplified implementation
	// Full implementation would:
	// 1. Fetch existing server state
	// 2. Compare with desired config
	// 3. Apply only the differences (update/create/delete)

	fmt.Println("  Note: Full sync implementation pending")
	fmt.Println("  For now, use 'setup' command or manually update via Discord")

	return nil
}
