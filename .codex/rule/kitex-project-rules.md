---
inclusion: always
---

## 目标

本项目开发前先锁定框架与目录边界：后端 RPC 服务必须基于 CloudWeGo Kitex 脚手架生成与演进，网关层 HTTP 服务必须基于 CloudWeGo Hertz 脚手架生成与演进。项目目录严格参考 Kitex 与 Hertz 官方示例结构，不允许自行创建无明确用途、脱离官方示例的目录。

## 框架硬约束

1. 后端 RPC 框架固定使用 Kitex。
2. 网关层 HTTP 框架固定使用 Hertz。
3. RPC 服务代码必须通过 Kitex 脚手架生成，禁止手写一套自定义服务启动框架替代 Kitex。
4. HTTP 网关代码必须通过 Hertz 的 `hz` 脚手架生成，禁止手写一套自定义 HTTP 启动框架替代 Hertz。
5. 所有 HTTP 路由、参数绑定、参数校验、中间件、响应渲染、静态资源、文件上传下载、反向代理、HTTP client、单元测试等 HTTP 操作，必须优先使用 Hertz 官方框架或 Hertz 官方扩展能力。
6. 如果 Hertz 官方框架内没有提供当前要用到的功能、工具类或扩展能力，必须先提示用户，说明缺口、备选方案和影响，经过确认后才能新增自定义工具或第三方依赖。
7. 网关层禁止直接引入 Gin、Echo、Fiber、net/http 自建路由栈等替代 Hertz 的 HTTP 框架能力。
8. IDL 固定使用 Thrift，禁止混用 Protobuf；如需切换接口定义语言，必须先更新本规则并得到用户确认。
9. `kitex_gen/` 为 Kitex 生成代码目录，默认不手工编辑。
10. 每个 RPC 服务必须使用 `kitex -service` 生成服务脚手架。
11. 默认禁止使用自定义 Kitex scaffold template、custom layout 或 Hertz 自定义模板；如需启用，必须先更新本规则并得到用户确认。

## 官方示例目录约束

项目根目录只允许出现与 Kitex 官方示例一致或直接必要的目录与文件：

```text
.
├── go.mod
├── go.sum
├── idl/
├── kitex_gen/
├── api/                    # Hertz HTTP 网关层
└── rpc/
    └── <service>/
        ├── build.sh
        ├── handler.go
        ├── kitex_info.yaml
        ├── main.go
        └── script/
            └── bootstrap.sh
```

目录职责：

1. `idl/`：存放项目 IDL 文件，例如 `base.thrift`、`travel.thrift`。
2. `kitex_gen/`：存放 Kitex 根据 IDL 生成的序列化、反序列化、client、server stub 代码。
3. `rpc/<service>/`：存放单个 RPC 服务的 Kitex 脚手架代码。
4. `rpc/<service>/handler.go`：服务端业务逻辑主要入口。
5. `rpc/<service>/main.go`：服务启动与必要资源初始化入口。
6. `rpc/<service>/script/bootstrap.sh`、`build.sh`：Kitex 生成的构建与启动脚本，默认不做无关修改。
7. `api/`：Hertz HTTP 网关层目录，只放 HTTP 入口、路由、中间件、请求响应转换和 RPC client 调用编排，不放 RPC 服务内部实现。

## Hertz 网关层目录约束

`api/` 目录必须按照 Hertz 官方 `hz` 工具生成的标准结构开发。默认结构如下：

```text
api/
├── .hz
├── biz/
│   ├── dal/
│   ├── handler/
│   ├── model/
│   ├── router/
│   │   ├── <idl-related-router>/
│   │   └── register.go
│   └── service/
├── build.sh
├── main.go
├── router.go
├── router_gen.go
└── script/
    └── bootstrap.sh
```

目录职责：

1. `api/main.go`：Hertz 网关启动入口，只负责创建 Hertz server、注册路由和启动服务。
2. `api/router.go`：用户自定义路由注册入口，仅用于 Hertz 官方路由组织方式。
3. `api/router_gen.go`：`hz` 生成的路由注册代码，默认不手工编辑。
4. `api/biz/router/`：Hertz 路由分组与中间件挂载代码，必须使用 Hertz 的路由和 group 能力。
5. `api/biz/handler/`：HTTP 控制层，只负责接收 HTTP 请求、绑定参数、基础校验、调用网关 service、返回 HTTP 响应。
6. `api/biz/service/`：网关业务编排层，负责调用后端 Kitex RPC client、组合结果和处理网关侧流程。
7. `api/biz/model/`：`hz` 根据 IDL 或注解生成/维护的 HTTP 请求响应模型，默认不手工乱改。
8. `api/biz/dal/`：仅在网关层确实需要访问网关自有数据时使用；默认不允许网关直接读写业务服务私有数据表。
9. `api/script/bootstrap.sh`、`api/build.sh`：Hertz 生成的启动与构建脚本，默认不做无关修改。

Hertz 开发约束：

1. 新建 HTTP 网关必须使用 `hz new`；更新 IDL 或路由必须优先使用 `hz update`。
2. 参数绑定与校验必须使用 Hertz 官方 binding、validator、RequestContext 能力，禁止手写重复解析框架。
3. JSON、HTML、Protobuf、文件、重定向等 HTTP 响应必须使用 Hertz 官方 render/response 能力。
4. 鉴权、阻断、限流、CORS、Recovery、日志、RequestID、压缩、安全头等横切能力必须优先使用 Hertz server-side middleware 或 Hertz 官方/官方示例扩展。
5. 对外 HTTP 调用必须优先使用 Hertz client 及其 middleware 能力；不得随意混用其他 HTTP client 封装。
6. 路由终止必须使用 Hertz `Abort`、`AbortWithMsg`、`AbortWithStatus` 等框架能力。
7. 单元测试优先使用 Hertz 提供的无网络传输测试接口和官方示例方式。
8. 需要 Swagger、JWT、Session、CSRF、Secure、ReverseProxy、WebSocket、SSE、Prometheus 等能力时，必须先查 Hertz 官方文档或 hertz-examples；官方已有示例或扩展时按官方方式接入。
9. Hertz 官方没有覆盖的工具类、公共 helper、第三方库或自定义 middleware，必须先说明为什么 Hertz 现有能力无法满足，再经用户确认后新增。

## 禁止事项

1. RPC 服务下禁止在未说明必要性的情况下新建 `cmd/`、`internal/`、`pkg/`、`app/`、`service/`、`controller/`、`repository/`、`biz/`、`configs/` 等非 Kitex 官方示例目录。
2. 网关层禁止在 `api/` 官方 Hertz 结构之外新建自定义框架目录；`api/biz/service`、`api/biz/handler`、`api/biz/router` 等 Hertz 官方目录除外。
3. 禁止把业务代码拆到自定义目录后留下空壳 handler。
4. 禁止手工修改 `kitex_gen/` 中的生成代码；需要变更接口时应修改 `idl/` 并重新运行 Kitex 生成。
5. 禁止手工修改 `api/router_gen.go` 和 Hertz 生成的路由代码；需要变更 HTTP 接口时应修改 IDL 或路由来源并运行 `hz update`。
6. 禁止提交无关 IDE、临时文件、构建产物或本地环境文件。
7. 禁止为了“预留扩展”创建空目录或占位目录。
8. 禁止绕过 Kitex 默认脚手架生成 main package；不得用 custom layout 生成 MVC、分层目录等非本项目规则允许的结构。
9. 禁止绕过 Hertz 默认脚手架生成网关结构；不得在 `api/` 下另起一套非 Hertz 官方结构。
10. 禁止把参数校验、鉴权、业务编排、数据处理、数据持久化、响应组装全部堆在同一个函数或同一个文件里。
11. 禁止跨服务随意引用其他 `rpc/<service>/` 下的内部实现；跨服务协作必须通过 IDL 暴露的 RPC 接口完成。

## 多服务并行开发约束

1. 每个业务服务独立放在 `rpc/<service>/` 下，服务之间以 IDL 和 RPC 调用作为边界。
2. 并行开发时，每个任务默认只修改自己负责的 `rpc/<service>/`、对应 `idl/` 和必要的生成代码。
3. 涉及共享 IDL、公共错误码、跨服务调用链或接口兼容性变更时，必须先统一接口方案，再分别实现服务。
4. 禁止某个服务为了开发方便直接读取或修改另一个服务的数据表、缓存 key、私有结构体或内部方法。
5. 服务之间传输的数据结构以 IDL 定义为准，不在业务代码里私自拼装未声明的跨服务协议。

## 单服务工程化分层约束

在不破坏 Kitex 官方脚手架目录的前提下，单个 `rpc/<service>/` 内必须按职责分层组织代码。默认优先使用同目录下的明确职责文件，而不是新增目录。

推荐文件组织：

```text
rpc/<service>/
├── main.go
├── handler.go
├── middleware.go
├── logic_<domain>.go
├── data_<domain>.go
├── repository_<domain>.go
├── converter_<domain>.go
├── errors.go
├── build.sh
├── kitex_info.yaml
└── script/
    └── bootstrap.sh
```

职责边界：

1. 控制层：Kitex handler 是 RPC 控制层入口，只负责接收请求、提取参数、基础参数校验、调用业务逻辑、返回响应。
2. 中间件：`middleware.go` 负责鉴权、阻断、限流、日志、链路追踪、panic recovery 等横切逻辑，不写具体业务流程。
3. 业务逻辑层：`logic_<domain>.go` 负责编排业务流程、领域规则判断、调用数据处理和持久化能力，不直接拼 SQL 或封装底层存储细节。
4. 数据处理层：`data_<domain>.go` 负责清洗、计算、聚合、去重、排序、状态流转等纯数据处理逻辑，尽量保持无外部副作用。
5. 数据持久化层：`repository_<domain>.go` 负责数据库、缓存、外部存储读写，必须封装参数化查询、事务边界和存储错误转换。
6. 数据返回层：`converter_<domain>.go` 负责领域对象、存储对象与 IDL response 之间的转换，禁止在 handler 中散落响应拼装逻辑。
7. 错误处理：`errors.go` 负责本服务错误定义、错误码映射和对外错误转换，禁止在各层随意硬编码错误消息。

调用方向：

```text
handler -> logic -> data/repository -> converter -> response
middleware -> handler
```

约束细则：

1. `handler.go` 中单个方法只保留控制流，不承载复杂业务规则；当逻辑超过清晰可读范围时必须下沉到 `logic_<domain>.go`。
2. 鉴权与阻断逻辑必须进入中间件或明确的 guard 方法，禁止散落在多个 handler 方法中重复判断。
3. 数据处理和数据持久化必须分离；纯计算逻辑不得依赖数据库连接、缓存客户端或 RPC client。
4. 数据持久化方法只表达存储语义，不负责拼装最终响应。
5. 响应组装必须集中在 converter 或明确的返回构造方法中，保证 IDL 字段映射可追踪。
6. 同一业务能力的参数校验、业务规则、存储访问、响应转换必须各在其位，禁止为了省事写成一个大函数。
7. 若单服务复杂度确实需要子目录分层，必须先更新本规则并得到用户确认；默认不创建 `controller/`、`service/`、`repository/` 等目录。

## 生成流程

### Kitex RPC 服务

首次生成项目骨架时：

```bash
go mod init <module>
mkdir idl
```

生成 IDL 相关代码：

```bash
kitex -module <module> idl/<service>.thrift
```

生成单个 RPC 服务脚手架：

```bash
mkdir -p rpc/<service>
cd rpc/<service>
kitex -module <module> -service <service-name> -use <module>/kitex_gen ../../idl/<service>.thrift
```

生成后回到项目根目录执行依赖整理：

```bash
go mod tidy
```

### Hertz HTTP 网关

首次生成网关层时：

```bash
cd api
hz new -module <module>/api -idl ../idl/<api>.thrift
```

更新 HTTP 接口、路由或模型时：

```bash
cd api
hz update -idl ../idl/<api>.thrift
```

生成或更新后回到项目根目录执行依赖整理：

```bash
go mod tidy
```

## 新增目录审批规则

如果后续确实需要新增不在官方示例中的目录，必须同时满足：

1. 有明确业务或工程必要性。
2. 不能由现有 `idl/`、`kitex_gen/`、`rpc/<service>/`、`api/` 官方结构承担。
3. 在变更前更新本规则，写清新增目录名称、职责、边界和禁止放入的内容。
4. 得到用户明确确认后再创建。

## 依据

1. Kitex 官方 Code Generation Tool 文档说明：`kitex` 是 Kitex 提供的代码生成命令，支持通过 `-service` 生成服务端项目脚手架。
2. Kitex 官方高级教程示例结构包含 `go.mod`、`go.sum`、`idl/`、`kitex_gen/`、`rpc/<service>/build.sh`、`handler.go`、`kitex_info.yaml`、`main.go`、`script/bootstrap.sh`。
3. 官方参考链接：
   - https://www.cloudwego.io/docs/kitex/tutorials/code-gen/code_generation/
   - https://www.cloudwego.io/docs/kitex/getting-started/tutorial/
   - https://www.cloudwego.io/docs/kitex/tutorials/code-gen/custom_tpl/
   - https://www.cloudwego.io/docs/hertz/
   - https://www.cloudwego.io/docs/hertz/tutorials/toolkit/
   - https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/
   - https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/middleware/
   - https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/binding-and-validate/
