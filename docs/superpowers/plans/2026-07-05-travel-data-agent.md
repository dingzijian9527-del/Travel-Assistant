# 旅行数据增强智能体实施计划

> **给自动化执行者：**需要按任务逐项执行。步骤使用复选框跟踪。

**目标：**把智能体快捷能力升级为可调用天气、地图路线、地点搜索和预算拆分的完整旅行规划能力。

**架构：**新增 `pkg/traveldata` 作为第三方旅行数据聚合层，封装高德地图、和风天气、预算拆分和兜底数据。`rpc/ai_agent` 只负责读取上下文、调用聚合层、把增强资料注入大模型提示词；`rpc/trip` 继续只负责保存结构化行程。

**技术栈：**Go、Hertz、Kitex、MySQL、高德开放平台、和风天气、Vue、uni-app。

---

### 任务一：配置第三方旅行数据能力

**文件：**
- 修改：`pkg/config/config.go`
- 修改：`conf/config.yaml`
- 修改：`conf/config.yaml.tpl`
- 测试：`pkg/config/config_test.go`

- [ ] 写配置测试：断言新增 `travel_data.amap_key`、`travel_data.qweather_key`、`travel_data.enabled`、`travel_data.timeout` 可读取。
- [ ] 运行 `go test ./pkg/config`，预期因为字段不存在失败。
- [ ] 新增 `TravelDataConfig`，并加入默认值。
- [ ] 补齐两个配置文件的示例字段。
- [ ] 运行 `go test ./pkg/config`，预期通过。

### 任务二：实现旅行数据聚合层

**文件：**
- 创建：`pkg/traveldata/model.go`
- 创建：`pkg/traveldata/budget.go`
- 创建：`pkg/traveldata/planner.go`
- 创建：`pkg/traveldata/amap.go`
- 创建：`pkg/traveldata/qweather.go`
- 测试：`pkg/traveldata/budget_test.go`
- 测试：`pkg/traveldata/planner_test.go`

- [ ] 写预算拆分测试：总预算三千元、三天、两人时，输出住宿、餐饮、交通、门票、机动费用，合计等于三千。
- [ ] 写聚合层测试：使用假的天气和地图客户端，断言输出包含天气、路线、地点和预算。
- [ ] 运行 `go test ./pkg/traveldata`，预期因为包未实现失败。
- [ ] 实现模型、预算拆分、聚合器和第三方客户端接口。
- [ ] 实现高德客户端：城市地理编码、周边地点搜索、路线距离和耗时。
- [ ] 实现和风天气客户端：城市经纬度天气查询。
- [ ] 运行 `go test ./pkg/traveldata`，预期通过。

### 任务三：智能体接入旅行数据增强

**文件：**
- 修改：`rpc/ai_agent/service.go`
- 修改：`rpc/ai_agent/handler.go`
- 修改：`rpc/ai_agent/stream_service.go`
- 测试：`rpc/ai_agent/service_test.go`

- [ ] 写服务测试：用户提问“成都三天预算3000”时，发送给模型的消息包含天气、预算拆分、路线建议和地点推荐。
- [ ] 运行 `go test ./rpc/ai_agent`，预期因为未注入旅行数据失败。
- [ ] 给 `aiAgentService` 增加 `travelPlanner` 依赖。
- [ ] 在普通对话和流式对话构建模型消息前调用旅行数据聚合层。
- [ ] 将旅行数据格式化为“参考资料”，和原有检索资料一起送入模型。
- [ ] 运行 `go test ./rpc/ai_agent`，预期通过。

### 任务四：前端快捷能力和推荐问题接入智能体

**文件：**
- 修改：`front/src/pages/ai/index.vue`

- [ ] 把能力卡片的提示词改为结构化语义：规划行程、选酒店、找美食、出行提醒都带上“请结合天气、路线、预算和本地推荐”。
- [ ] 实现“换一组”问题池轮换，不再只填一个固定问题。
- [ ] 点击能力卡片或推荐问题后直接发送到智能体。
- [ ] 运行 `npm run build:h5`，预期通过。

### 任务五：全量验证

**文件：**
- 全项目

- [ ] 运行 `go test ./...`，预期通过。
- [ ] 运行 `npm run build:h5`，预期通过。
- [ ] 手动验证：前端进入智能体页，点击“规划行程”“选酒店”“换一组”，能自动发送并得到带预算、天气、路线语境的回复。
