package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/admin/xagent/internal/agent"
	"github.com/admin/xagent/internal/appctx"
	"github.com/admin/xagent/internal/config"
)

// Execute 是 agent-cli 的命令行入口。
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		slog.Error("命令执行失败", "error", err)
		os.Exit(1)
	}
}

// NewRootCmd 创建 Cobra 根命令。这里仅做启动器职责，不放业务规则。
func NewRootCmd() *cobra.Command {
	var cfgFile string
	var apiKey string
	var apiBase string
	var model string
	var pluginDir string
	var sessionDBPath string
	var mockMode bool

	rootCmd := &cobra.Command{
		Use:   "agent-cli",
		Short: "AI Agent CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.WithConfigFile(cfgFile))
			if err != nil {
				return fmt.Errorf("加载配置失败: %w", err)
			}
			applyFlagOverrides(cmd, cfg, apiKey, apiBase, model, pluginDir, sessionDBPath, mockMode)

			rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			app, err := appctx.NewAppCtx(rootCtx, cfg)
			if err != nil {
				return fmt.Errorf("初始化应用上下文失败: %w", err)
			}
			defer func() {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				if err := app.Shutdown(shutdownCtx); err != nil {
					slog.Warn("应用关闭不完整", "error", err)
				}
			}()

			program := tea.NewProgram(app.TUI)
			app.TUI.SetSender(program.Send)
			app.TUI.SetSubmitHandler(func(input string) tea.Cmd {
				return func() tea.Msg {
					if err := agent.RunReActLoop(app, input); err != nil {
						app.TUI.ReportError(err)
					}
					return nil
				}
			})

			if _, err := program.Run(); err != nil {
				return fmt.Errorf("TUI 启动失败: %w", err)
			}
			return nil
		},
	}

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "配置文件路径")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "LLM API Key")
	rootCmd.Flags().StringVar(&apiBase, "api-base", "", "LLM API Base URL")
	rootCmd.Flags().StringVar(&model, "model", "", "LLM 模型名")
	rootCmd.Flags().StringVar(&pluginDir, "plugin-dir", "", "插件目录")
	rootCmd.Flags().StringVar(&sessionDBPath, "session-db-path", "", "会话数据库路径")
	rootCmd.Flags().BoolVar(&mockMode, "mock-mode", true, "启用本地假响应模式")

	return rootCmd
}

func applyFlagOverrides(cmd *cobra.Command, cfg *config.Config, apiKey, apiBase, model, pluginDir, sessionDBPath string, mockMode bool) {
	if cmd.Flags().Changed("api-key") {
		cfg.APIKey = apiKey
	}
	if cmd.Flags().Changed("api-base") {
		cfg.APIBase = apiBase
	}
	if cmd.Flags().Changed("model") {
		cfg.Model = model
	}
	if cmd.Flags().Changed("plugin-dir") {
		cfg.PluginDir = pluginDir
	}
	if cmd.Flags().Changed("session-db-path") {
		cfg.SessionDBPath = sessionDBPath
	}
	if cmd.Flags().Changed("mock-mode") {
		cfg.MockMode = mockMode
	}
}
