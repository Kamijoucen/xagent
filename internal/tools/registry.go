package tools

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Registry 维护可用工具集合。
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry 创建工具注册表。
func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

// Register 注册工具。
func (r *Registry) Register(name string, tool Tool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("工具名不能为空")
	}
	if tool == nil {
		return fmt.Errorf("工具 %s 不能为空", name)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[name] = tool
	return nil
}

// Get 根据名称获取工具。
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

// List 返回按名称排序的工具列表。
func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	items := make([]Tool, 0, len(names))
	for _, name := range names {
		items = append(items, r.tools[name])
	}
	return items
}
