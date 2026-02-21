package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all Discord server configuration
type Config struct {
	Server       ServerConfig       `yaml:"-"`
	Channels     ChannelsConfig     `yaml:"-"`
	Roles        RolesConfig        `yaml:"-"`
	Integrations IntegrationsConfig `yaml:"-"`

	// Runtime configuration
	BotToken string `yaml:"-"`
	GuildID  string `yaml:"-"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Settings    struct {
		VerificationLevel         string `yaml:"verification_level"`
		DefaultNotificationLevel  string `yaml:"default_notification_level"`
		ExplicitContentFilter     string `yaml:"explicit_content_filter"`
	} `yaml:"settings"`
	Features struct {
		Community    bool `yaml:"community"`
		Discoverable bool `yaml:"discoverable"`
	} `yaml:"features"`
	VanityURL string `yaml:"vanity_url,omitempty"`
}

// ChannelsConfig holds channel structure
type ChannelsConfig struct {
	Categories []Category `yaml:"categories"`
}

type Category struct {
	Name     string    `yaml:"name"`
	Position int       `yaml:"position"`
	Channels []Channel `yaml:"channels"`
}

type Channel struct {
	Name        string                       `yaml:"name"`
	Type        string                       `yaml:"type"` // text, voice, forum
	Topic       string                       `yaml:"topic,omitempty"`
	Position    int                          `yaml:"position"`
	Permissions map[string]map[string]bool   `yaml:"permissions,omitempty"`
	Tags        []ForumTag                   `yaml:"available_tags,omitempty"`
}

type ForumTag struct {
	Name  string `yaml:"name"`
	Emoji string `yaml:"emoji"`
}

// RolesConfig holds role definitions
type RolesConfig struct {
	Roles []Role `yaml:"roles"`
}

type Role struct {
	Name        string   `yaml:"name"`
	Color       string   `yaml:"color"`
	Permissions []string `yaml:"permissions"`
	Hoist       bool     `yaml:"hoist"`
	Mentionable bool     `yaml:"mentionable"`
	Description string   `yaml:"description,omitempty"`
}

// IntegrationsConfig holds webhook and integration settings
type IntegrationsConfig struct {
	GitHub *GitHubIntegration `yaml:"github,omitempty"`
}

type GitHubIntegration struct {
	Enabled       bool     `yaml:"enabled"`
	TargetChannel string   `yaml:"target_channel"`
	Events        []string `yaml:"events"`
}

// Load reads all configuration files and environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Load environment variables
	cfg.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN environment variable is required")
	}

	cfg.GuildID = os.Getenv("DISCORD_GUILD_ID")
	if cfg.GuildID == "" {
		return nil, fmt.Errorf("DISCORD_GUILD_ID environment variable is required")
	}

	// Load YAML config files
	configDir := "config"

	if err := loadYAML(filepath.Join(configDir, "server.yaml"), &cfg.Server); err != nil {
		return nil, fmt.Errorf("loading server config: %w", err)
	}

	if err := loadYAML(filepath.Join(configDir, "channels.yaml"), &cfg.Channels); err != nil {
		return nil, fmt.Errorf("loading channels config: %w", err)
	}

	if err := loadYAML(filepath.Join(configDir, "roles.yaml"), &cfg.Roles); err != nil {
		return nil, fmt.Errorf("loading roles config: %w", err)
	}

	if err := loadYAML(filepath.Join(configDir, "integrations.yaml"), &cfg.Integrations); err != nil {
		return nil, fmt.Errorf("loading integrations config: %w", err)
	}

	return cfg, nil
}

func loadYAML(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	return nil
}
