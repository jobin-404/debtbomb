package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Jira   JiraConfig        `toml:"jira"`
	Owners map[string]string `toml:"owners"`
	Notify []NotifyConfig    `toml:"notify"`
}

type JiraConfig struct {
	DefaultProject string `toml:"default_project"`
	IssueType      string `toml:"issue_type"`
	BaseURL        string
	Email          string
	APIToken       string
}

type NotifyConfig struct {
	On   string `toml:"on"`
	Via  string `toml:"via"`
	Days int    `toml:"days"`
}

func Load(rootPath string) (*Config, error) {
	configPath := filepath.Join(rootPath, ".debtbomb", "config.toml")

	var conf Config
	if _, err := os.Stat(configPath); err == nil {
		if _, err := toml.DecodeFile(configPath, &conf); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Load secrets from env vars
	conf.Jira.BaseURL = os.Getenv("JIRA_BASE_URL")
	conf.Jira.Email = os.Getenv("JIRA_EMAIL")
	conf.Jira.APIToken = os.Getenv("JIRA_API_TOKEN")

	return &conf, nil
}

func (c *Config) GetSlackWebhook() string {
	return os.Getenv("SLACK_WEBHOOK_URL")
}

func (c *Config) GetDiscordWebhook() string {
	return os.Getenv("DISCORD_WEBHOOK_URL")
}

func (c *Config) GetTeamsWebhook() string {
	return os.Getenv("TEAMS_WEBHOOK_URL")
}
