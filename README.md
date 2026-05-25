# agent-cli

`agent-cli` 是一个 Go 版终端 AI Agent 客户端模板。第一版是 MVP：保留 Atomic Architecture 的清晰分层，但只实现可启动 TUI 和本地假交互。

## 当前范围

- 可以启动 Bubble Tea TUI。
- 可以输入文本并按 Enter 发送。
- Agent 会返回一条本地假响应。
- 配置、LLM 客户端、插件管理器、会话存储、工具注册表都有清晰扩展点。
- 真实 LLM、插件进程池、高级长期记忆、完整 ReAct 循环暂未实现。

## 架构

```text
cmd/                  Service 启动器，只负责配置、AppCtx、TUI 和信号处理
internal/appctx/      ApplicationContext，持有并管理组件生命周期
internal/agent/       Agent 编排层，当前是本地假响应
internal/agent/logic/ Logic 原子函数，当前是简单实现和占位
internal/components/  TUI、LLM、插件、会话存储组件
internal/tools/       Tool 接口与注册表
internal/plugins/     JSON-RPC 协议结构
plugins/example/      独立可编译的示例插件
```

依赖方向遵循：`cmd -> appctx`、`cmd -> agent`、`agent/logic -> appctx`、`appctx -> components/tools`、`components -> infrastructure`。

## 运行

```bash
go mod tidy
make run
```

进入 TUI 后输入文本并按 Enter。使用 `Esc` 或 `Ctrl+C` 退出；输入为空时也可以按 `q` 退出。

## 构建

```bash
make build
./bin/agent-cli
```

Linux ARM64 交叉编译：

```bash
make build-linux-arm64
```

## 配置

配置加载优先级为：命令行 flag > 环境变量 > `~/.config/agent-cli/config.yaml` > 默认值。

可用字段：

```yaml
api_key: ""
api_base: "https://api.openai.com/v1"
model: "gpt-4o-mini"
max_tokens: 2048
temperature: 0.2
plugin_dir: "~/.config/agent-cli/plugins"
confirm_dangerous_tools: true
session_db_path: "~/.local/share/agent-cli/store.db"
mock_mode: true
```

第一版默认 `mock_mode: true`，不会请求真实 LLM。

## 插件示例

示例插件是独立 Go module：

```bash
cd plugins/example
go build .
```

它读取 JSON-RPC 2.0 stdin 请求，并返回 echo 响应：

```json
{"jsonrpc":"2.0","id":1,"method":"execute","params":{"args":{"text":"hello"}}}
```

插件进程池暂未接入主程序，后续会由 `internal/components/PluginManager` 扩展。

## 测试

```bash
make test
```

当前测试覆盖配置默认值、会话内存存储、工具注册表和 Logic 格式化函数。