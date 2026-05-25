package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 保存进程级配置。第一版只加载配置，不做真实 LLM 对接。
type Config struct {
	APIKey                string
	APIBase               string
	Model                 string
	MaxTokens             int
	Temperature           float32
	PluginDir             string
	ConfirmDangerousTools bool
	SessionDBPath         string
	MockMode              bool
}

type loadOptions struct {
	configFile string
}

// LoadOption 调整配置加载行为，主要用于测试和命令行启动器。
type LoadOption func(*loadOptions)

// WithConfigFile 指定配置文件路径。空路径时使用默认搜索路径。
func WithConfigFile(path string) LoadOption {
	return func(opts *loadOptions) {
		opts.configFile = path
	}
}

// Default 返回不依赖外部文件的默认配置。
func Default() *Config {
	return &Config{
		APIBase:               "https://api.openai.com/v1",
		Model:                 "gpt-4o-mini",
		MaxTokens:             2048,
		Temperature:           0.2,
		PluginDir:             defaultPluginDir(),
		ConfirmDangerousTools: true,
		SessionDBPath:         defaultSessionDBPath(),
		MockMode:              true,
	}
}

// Load 使用 Viper 加载配置。配置文件不存在时不会报错。
func Load(options ...LoadOption) (*Config, error) {
	opts := loadOptions{}
	for _, option := range options {
		option(&opts)
	}

	v := viper.New()
	setDefaults(v)
	v.SetEnvPrefix("AGENT_CLI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if opts.configFile != "" {
		v.SetConfigFile(opts.configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(filepath.Join(homeDir(), ".config", "agent-cli"))
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
	}

	return &Config{
		APIKey:                v.GetString("api_key"),
		APIBase:               v.GetString("api_base"),
		Model:                 v.GetString("model"),
		MaxTokens:             v.GetInt("max_tokens"),
		Temperature:           float32(v.GetFloat64("temperature")),
		PluginDir:             v.GetString("plugin_dir"),
		ConfirmDangerousTools: v.GetBool("confirm_dangerous_tools"),
		SessionDBPath:         v.GetString("session_db_path"),
		MockMode:              v.GetBool("mock_mode"),
	}, nil
}

func setDefaults(v *viper.Viper) {
	defaults := Default()
	v.SetDefault("api_key", defaults.APIKey)
	v.SetDefault("api_base", defaults.APIBase)
	v.SetDefault("model", defaults.Model)
	v.SetDefault("max_tokens", defaults.MaxTokens)
	v.SetDefault("temperature", defaults.Temperature)
	v.SetDefault("plugin_dir", defaults.PluginDir)
	v.SetDefault("confirm_dangerous_tools", defaults.ConfirmDangerousTools)
	v.SetDefault("session_db_path", defaults.SessionDBPath)
	v.SetDefault("mock_mode", defaults.MockMode)
}

func defaultPluginDir() string {
	return filepath.Join(homeDir(), ".config", "agent-cli", "plugins")
}

func defaultSessionDBPath() string {
	return filepath.Join(homeDir(), ".local", "share", "agent-cli", "store.db")
}

func homeDir() string {
	if dir, err := filepath.Abs("."); err == nil && dir != "" {
		if home := strings.TrimSpace(homeFromEnv()); home != "" {
			return home
		}
	}
	return "."
}

func homeFromEnv() string {
	v := viper.New()
	v.AutomaticEnv()
	if home := v.GetString("HOME"); home != "" {
		return home
	}
	if profile := v.GetString("USERPROFILE"); profile != "" {
		return profile
	}
	return ""
}
