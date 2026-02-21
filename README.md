# WorkFort Discord Infrastructure

Infrastructure as Code for the WorkFort Discord community server.

## Architecture

- **Discord Bot**: Go application using discordgo library
- **Configuration**: YAML files defining server structure
- **Task Runner**: mise for tool management and commands
- **IaC Philosophy**: All server configuration lives in version control

## Prerequisites

- [mise](https://mise.jdx.dev/) - Tool version manager
- Discord bot token and server (guild) ID
- Discord server with bot invited (Guild Install)

## Quick Start

### 1. Clone and Install Tools

```bash
git clone https://github.com/Work-Fort/Discord.git
cd Discord
mise install
```

This installs Go and any other required tools.

### 2. Secrets are Already Configured

This repository uses SOPS (Secrets OPerationS) with age encryption to store secrets securely in git.

**Secrets are stored in `secrets.yaml` (encrypted and committed):**
- Discord bot token
- Discord guild (server) ID
- GitHub webhook URL

**The age private key (`age-key.txt`) is git-ignored** and must be obtained from the team lead.

To decrypt secrets (view only):
```bash
mise run secrets_decrypt
```

To edit secrets:
```bash
mise run secrets_edit
```

**For new contributors:** Contact the team lead for `age-key.txt` and place it in the repo root

### 3. Initial Setup

Run the initial setup to create channels, roles, and configure the server:

```bash
mise run setup
```

This reads `config/*.yaml` files and applies them to your Discord server.

## Configuration Files

All server configuration lives in `config/`:

- `server.yaml` - Server settings and metadata
- `channels.yaml` - Channel structure and permissions
- `roles.yaml` - Role definitions and permissions
- `integrations.yaml` - Webhooks and external integrations

## Available Commands

```bash
# Initial server setup from YAML configs
mise run setup

# Sync configuration changes to Discord
mise run sync

# Export current Discord state to YAML (backup/drift detection)
mise run backup

# Validate YAML configuration files
mise run validate

# Build the binary
mise run build

# Run tests
mise run test

# View encrypted secrets
mise run secrets_decrypt

# Edit encrypted secrets (opens editor)
mise run secrets_edit
```

## Secrets Management

This repository uses **SOPS** (Secrets OPerationS) with **age** encryption to securely store secrets in git.

- **`secrets.yaml`**: Encrypted secrets file (COMMITTED to git)
- **`age-key.txt`**: Private encryption key (GIT-IGNORED, never committed)
- **`.sops.yaml`**: SOPS configuration with age public key

All secrets (Discord bot token, guild ID, webhook URLs) are encrypted before being committed. The private age key is required to decrypt them.

## Development Workflow

1. Edit YAML configuration files in `config/`
2. Run `mise run validate` to check syntax
3. Run `mise run sync` to apply changes to Discord
4. Commit changes to git

## Backup and Drift Detection

Export current Discord server state:

```bash
mise run backup
```

This creates timestamped YAML snapshots. Compare with checked-in config to detect drift.

## Project Structure

```
.
├── .mise.toml              # Tool versions and task definitions
├── .sops.yaml              # SOPS encryption config (age public key)
├── secrets.yaml            # Encrypted secrets (COMMITTED)
├── age-key.txt             # age private key (GIT-IGNORED)
├── .env                    # Shared config (committed)
├── config/
│   ├── server.yaml         # Server settings
│   ├── channels.yaml       # Channel structure
│   ├── roles.yaml          # Roles and permissions
│   └── integrations.yaml   # Webhooks, bots
├── cmd/
│   └── discord-bot/
│       └── main.go         # CLI entry point
├── internal/
│   ├── setup/              # Initial setup logic
│   ├── sync/               # Config sync to Discord
│   ├── backup/             # Export Discord state
│   └── config/             # YAML config parsing
└── README.md               # This file
```

## License

GPL-2.0-only - See [LICENSE.md](LICENSE.md) for details.
