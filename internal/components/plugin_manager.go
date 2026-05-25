package components

import (
	"context"
	"errors"

	"github.com/admin/xagent/internal/agent/types"
)

// ErrPluginUnsupported 表示第一版插件进程池尚未启用。
var ErrPluginUnsupported = errors.New("插件系统第一版仅保留接口，尚未启用")

// PluginManager 是插件进程池的占位 Component。
type PluginManager struct {
	pluginDir string
}

// NewPluginManager 创建插件管理器占位实现。
func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{pluginDir: pluginDir}
}

// Scan 第一版不扫描插件，仅保留扩展点。
func (m *PluginManager) Scan(ctx context.Context) error {
	_ = ctx
	return nil
}

// Execute 第一版不执行插件，返回明确的 unsupported 结果。
func (m *PluginManager) Execute(ctx context.Context, pluginName string, args map[string]interface{}) (types.PluginResult, error) {
	_ = ctx
	_ = args
	return types.PluginResult{Success: false, Error: "插件未启用: " + pluginName}, ErrPluginUnsupported
}

// Shutdown 关闭插件管理器。当前没有后台进程。
func (m *PluginManager) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}
