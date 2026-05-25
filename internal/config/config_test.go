package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if cfg.Model == "" {
		t.Fatal("默认模型不能为空")
	}
	if !cfg.MockMode {
		t.Fatal("MVP 默认应启用 mock 模式")
	}
	if cfg.PluginDir == "" || cfg.SessionDBPath == "" {
		t.Fatal("默认路径不能为空")
	}
}

func TestLoadWithoutConfigFile(t *testing.T) {
	cfg, err := Load(WithConfigFile("/path/not/exist/config.yaml"))
	if err == nil {
		t.Fatal("显式指定不存在的配置文件应返回错误")
	}
	if cfg != nil {
		t.Fatal("加载失败时不应返回配置")
	}
}
