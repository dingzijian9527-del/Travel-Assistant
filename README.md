# Travel-Assistant

## 后端服务说明

本项目后端 RPC 服务使用 CloudWeGo Kitex 脚手架生成，接口统一使用 Thrift IDL 定义。

### 目录结构

- `idl/base.thrift`：公共错误码与通用返回结构。
- `idl/user.thrift`：用户账号与资料服务接口。
- `idl/ai_agent.thrift`：AI 智能旅行助手服务接口。
- `kitex_gen/`：Kitex 根据 IDL 生成的代码，不手工修改。
- `cmd/api/`：Hertz 网关启动入口；`api/`：网关业务、路由注册和处理逻辑。
- `cmd/user/`：用户 RPC 启动入口；`rpc/user/`：用户 RPC 业务实现；`rpc/scaffold/user/`：用户 RPC 脚手架构建与启动脚本。
- `cmd/ai-agent/`：智能体 RPC 启动入口；`rpc/ai_agent/`：AI 智能体 RPC 业务实现；`rpc/scaffold/ai_agent/`：智能体 RPC 脚手架构建与启动脚本。

### 用户服务 UserService

IDL 文件：`idl/user.thrift`

接口：

- `Register`：用户注册，入参包含手机号、密码、可选昵称、头像、常住城市和当前位置城市。
- `Login`：用户使用手机号和密码登录，返回令牌和用户资料。
- `GetProfile`：按用户字符串主键查询用户资料。
- `UpdateProfile`：更新昵称、头像、常住城市和当前位置城市。

### AI 智能体服务 AIAgentService

IDL 文件：`idl/ai_agent.thrift`

定位：

- 该服务是旅游专用 AI 智能体。
- 主要负责出行计划制定、旅游攻略制作、旅游目的地问答、交通住宿建议、天气穿搭建议、当地美食推荐等旅游相关问题。
- 对非旅游问题不展开回答，只引导用户回到旅游咨询场景。

接口：

- `Chat`：发送旅行咨询消息，返回 AI 助手回复和后续建议。
- `GetPromptSuggestions`：根据场景返回推荐提示词。

### 当前实现说明

- 当前持久层使用内存仓储，便于本地编译和单元测试。
- 后续接入数据库时，仅替换 `repo.go` 持久层实现。
- AI 服务当前使用旅游领域规则回复占位，不硬编码第三方模型地址或密钥。
- 对外错误统一封装为 `BaseResp`，不直接透传原生错误字符串。

### 验证命令

```bash
go test ./...
```

