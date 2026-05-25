package appctx

import (
	"context"
	"errors"
	"fmt"

	"github.com/admin/xagent/internal/components"
	"github.com/admin/xagent/internal/config"
	"github.com/admin/xagent/internal/tools"
)

// AppCtx 是进程级 ApplicationContext，只负责组件组装和生命周期管理。
type AppCtx struct {
	Config        *config.Config
	LLMClient     *components.LLMClient
	PluginManager *components.PluginManager
	SessionStore  *components.SessionStore
	TUI           *components.TUIModel
	ToolRegistry  *tools.Registry

	rootCtx context.Context
	cancel  context.CancelFunc
}

// NewAppCtx 创建进程级应用上下文。
func NewAppCtx(parent context.Context, cfg *config.Config) (*AppCtx, error) {
	if parent == nil {
		parent = context.Background()
	}
	if cfg == nil {
		cfg = config.Default()
	}

	rootCtx, cancel := context.WithCancel(parent)
	registry := tools.NewRegistry()
	store, err := components.NewSessionStore(cfg.SessionDBPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("初始化会话存储失败: %w", err)
	}

	return &AppCtx{
		Config: cfg,
		LLMClient: components.NewLLMClient(components.LLMClientConfig{
			APIKey:      cfg.APIKey,
			APIBase:     cfg.APIBase,
			Model:       cfg.Model,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
			MockMode:    cfg.MockMode,
		}),
		PluginManager: components.NewPluginManager(cfg.PluginDir),
		SessionStore:  store,
		TUI:           components.NewTUIModel(),
		ToolRegistry:  registry,
		rootCtx:       rootCtx,
		cancel:        cancel,
	}, nil
}

// Context 返回应用生命周期上下文。
func (a *AppCtx) Context() context.Context {
	if a == nil || a.rootCtx == nil {
		return context.Background()
	}
	return a.rootCtx
}

// Shutdown 按依赖逆序释放组件资源。
func (a *AppCtx) Shutdown(ctx context.Context) error {
	if a == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if a.cancel != nil {
		a.cancel()
	}

	var errs []error
	if a.PluginManager != nil {
		errs = append(errs, a.PluginManager.Shutdown(ctx))
	}
	if a.SessionStore != nil {
		errs = append(errs, a.SessionStore.Shutdown(ctx))
	}
	if a.LLMClient != nil {
		errs = append(errs, a.LLMClient.Shutdown(ctx))
	}
	return errors.Join(errs...)
}
