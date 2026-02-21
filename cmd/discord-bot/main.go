package main

import (
	"fmt"
	"os"

	"github.com/Work-Fort/Discord/internal/backup"
	"github.com/Work-Fort/Discord/internal/config"
	"github.com/Work-Fort/Discord/internal/invite"
	"github.com/Work-Fort/Discord/internal/setup"
	"github.com/Work-Fort/Discord/internal/sync"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "setup":
		runSetup()
	case "sync":
		runSync()
	case "backup":
		runBackup()
	case "validate":
		runValidate()
	case "create-invite":
		runCreateInvite()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("WorkFort Discord Infrastructure")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  discord-bot <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup          Initial Discord server setup from YAML configs")
	fmt.Println("  sync           Sync config changes to Discord server")
	fmt.Println("  backup         Export current Discord state to YAML")
	fmt.Println("  validate       Validate YAML configuration files")
	fmt.Println("  create-invite  Create or retrieve permanent server invite link")
	fmt.Println()
	fmt.Println("Environment variables:")
	fmt.Println("  DISCORD_BOT_TOKEN  Discord bot token (required)")
	fmt.Println("  DISCORD_GUILD_ID   Discord server/guild ID (required)")
}

func runSetup() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := setup.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error running setup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Discord server setup complete")
}

func runSync() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := sync.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error running sync: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Discord server sync complete")
}

func runBackup() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := backup.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error running backup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Discord server backup complete")
}

func runValidate() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Configuration is valid")
	fmt.Printf("  Server: %s\n", cfg.Server.Name)
	fmt.Printf("  Channels: %d categories\n", len(cfg.Channels.Categories))
	fmt.Printf("  Roles: %d roles\n", len(cfg.Roles.Roles))
}

func runCreateInvite() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := invite.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating invite: %v\n", err)
		os.Exit(1)
	}
}
