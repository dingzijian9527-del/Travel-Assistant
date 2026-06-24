---
inclusion: always
---

# 后端开发规则遵循总则

## 目标

本规则用于明确后端开发过程中的前置约束：所有后端相关开发、修改、重构、测试、提交和问题排查，都必须先读取并遵循 `.codex/rule` 目录下的项目规则。

## 强制要求

1. 后端开发开始前，必须先检查 `.codex/rule` 下的规则文件。
2. 任何后端实现都必须同时遵循框架目录规则、代码规范、技术选型规则、工程治理补充规则、Code Review 规则和提交规范。
3. 如果用户临时要求与现有规则冲突，必须先提示冲突点，等待用户明确确认后再调整规则或执行变更。
4. 如果实现过程中发现规则缺失、不清晰或互相冲突，必须先暂停相关实现，补充或修正规则后再继续。
5. 不允许绕过规则直接新增目录、依赖、框架、基础设施、中间件或公共工具。

## 后端开发前检查清单

每次后端开发前至少确认：

1. 已阅读 `.codex/rule/kitex-project-rules.md`，确认 Kitex、Hertz 和目录结构约束。
2. 已阅读 `.codex/rule/project-code-style-rules.md`，确认编码、注释、分层、错误处理和测试规范。
3. 已阅读 `.codex/rule/project-tech-stack-rules.md`，确认 Go、MySQL、Milvus、Redis、Kafka 和 UniApp 等技术选型。
4. 已阅读 `.codex/rule/project-engineering-governance.md`，确认 IDL、依赖、配置、错误码、日志、AI 生成代码、韧性并发和版本治理要求。
5. 已阅读 `.codex/rule/code-review.md`，确认提交前审查要求。
6. 已阅读 `.codex/rule/commit-convention.md`，确认提交信息规范。

## 执行原则

1. 规则优先于个人习惯和默认工程模板。
2. 官方框架能力优先于自定义封装。
3. 项目已有规范优先于新增约定。
4. 小范围变更也必须遵守目录、分层、注释、测试和技术选型约束。
5. 规则更新后，后续开发必须以最新规则为准。
