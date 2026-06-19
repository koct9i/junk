package config

import (
	"context"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	ytsdk "go.ytsaurus.tech/yt/go/yt"
)

var (
	Config      = ytsdk.Config{}
	ProfileName string
	TokenPath   string
	Format      string

	EnableProxyDiscovery *bool
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "profile",
			Usage:       "Name of configuration profile",
			Sources:     cli.EnvVars("YT_CONFIG_PROFILE"),
			Destination: &ProfileName,
		},
		&cli.StringFlag{
			Name:        "proxy",
			Usage:       "YT HTTP proxy; defaults to YT_PROXY when empty",
			Sources:     cli.EnvVars("YT_PROXY"),
			Destination: &Config.Proxy,
		},
		&cli.StringFlag{
			Name:        "token-path",
			Usage:       "Path to a file with user token",
			Sources:     cli.EnvVars("YT_TOKEN_PATH"),
			Destination: &TokenPath,
		},
		&cli.BoolWithInverseFlag{
			Name:    "discover-proxy",
			Aliases: []string{"use-hosts"},
			Usage:   "Enable proxy discovery",
			Sources: cli.EnvVars("YT_USE_HOSTS"),
			Action: func(ctx context.Context, c *cli.Command, enable bool) error {
				EnableProxyDiscovery = &enable
				return nil
			},
		},
		&cli.StringFlag{
			Name:        "format",
			Aliases:     []string{"f"},
			Usage:       "data format: json, yson, or yaml",
			Value:       "json",
			Destination: &Format,
		},
	}
}

func LoadConfig() error {
	profile, err := LoadProfile(ProfileName)
	if err != nil {
		return err
	}
	if Config.Token == "" && profile.Token != "" {
		Config.Token = profile.Token
	}
	if TokenPath == "" && profile.TokenPath != "" {
		TokenPath = profile.TokenPath
	}
	if Config.Proxy == "" && profile.Proxy.URL != "" {
		Config.Proxy = profile.Proxy.URL
	}
	if EnableProxyDiscovery != nil {
		Config.DisableProxyDiscovery = !*EnableProxyDiscovery
	} else if profile.Proxy.EnableProxyDiscovery != nil {
		Config.DisableProxyDiscovery = !*profile.Proxy.EnableProxyDiscovery
	}
	if TokenPath != "" {
		if strings.HasPrefix(TokenPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			TokenPath = homeDir + TokenPath[1:]
		}
		token, err := os.ReadFile(TokenPath)
		if err != nil {
			return err
		}
		Config.Token = strings.TrimSpace(string(token))
	}
	return nil
}
