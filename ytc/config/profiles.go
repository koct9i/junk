package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.ytsaurus.tech/yt/go/yson"
)

type ConfigFile struct {
	ConfigVersion  int                      `yson:"config_version" json:"config_version"`
	DefaultProfile string                   `yson:"default_profile" json:"default_profile"`
	Profiles       map[string]ConfigProfile `yson:"profiles" json:"profiles"`
}

type ConfigProfile struct {
	Proxy     ConfigProxy `yson:"proxy" json:"proxy"`
	Token     string      `yson:"token" json:"token"`
	TokenPath string      `yson:"token_path" json:"token_path"`
}

type ConfigProxy struct {
	URL                  string `yson:"url" json:"url"`
	HTTPProxyRole        string `yson:"http_proxy_role" json:"http_proxy_role"`
	RPCProxyRole         string `yson:"rpc_proxy_role" json:"rpc_proxy_role"`
	NetworkName          string `yson:"network_name" json:"network_name"`
	ProxyDiscoveryURL    string `yson:"proxy_discovery_url" json:"proxy_discovery_url"`
	EnableProxyDiscovery *bool  `yson:"enable_proxy_discovery" json:"enable_proxy_discovery"`
	PreferHTTPS          *bool  `yson:"prefer_https" json:"prefer_https"`
	TVMOnly              *bool  `yson:"tvm_only" json:"tvm_only"`
}

func resolveConfigPath() (string, bool) {
	configPaths := []string{
		os.Getenv("YT_CONFIG_PATH"),
		"",
		"/etc/ytclient.conf",
	}
	if homeDir, err := os.UserHomeDir(); err == nil && homeDir != "" {
		configPaths[1] = filepath.Join(homeDir, ".yt", "config")
	}
	for _, configPath := range configPaths {
		if configPath == "" {
			continue
		}
		if stat, err := os.Stat(configPath); err == nil && stat.Mode().IsRegular() {
			return configPath, true
		}
	}
	return "", false
}

func LoadProfile(profileName string) (*ConfigProfile, error) {
	configPath, ok := resolveConfigPath()
	if !ok {
		return &ConfigProfile{}, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YT config from %s: %w", configPath, err)
	}

	var parsed ConfigFile
	switch configFormat := os.Getenv("YT_CONFIG_FORMAT"); configFormat {
	case "yson", "":
		if err = yson.Unmarshal(content, &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse YT config from %s: %w", configPath, err)
		}
	case "json":
		if err = json.Unmarshal(content, &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse YT config from %s: %w", configPath, err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format %q (expected yson or json)", configFormat)
	}

	if parsed.ConfigVersion != 2 {
		return nil, fmt.Errorf("unsupported config version %d", parsed.ConfigVersion)
	}
	if parsed.Profiles == nil {
		return nil, fmt.Errorf("missing profiles key in YT config")
	}

	if profileName == "" {
		profileName = os.Getenv("YT_CONFIG_PROFILE")
		if profileName == "" {
			profileName = parsed.DefaultProfile
			if profileName == "" {
				return nil, fmt.Errorf("profile has not been set and there is no default profile in the config")
			}
		}
	}

	profile, ok := parsed.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("unknown profile %q", profileName)
	}
	return &profile, nil
}
