# go-stock 系统架构与业务/数据流分析文档计划

## 1. 引言与目标

本计划旨在为 go-stock 跨平台股票分析桌面应用制定一份系统性技术文档，聚焦于业务/数据流分析，结合 MCP（多步链式思考）方法论和 Playwright 自动化测试工具，提升系统的可维护性、可扩展性和可验证性。

**约束：所有架构图、流程图、时序图等均需采用 drawio 格式绘制，并统一输出至 docs/system-arch-image 目录。文档中如有图示需求，均以 drawio 文件形式存放并引用。**

- [系统架构图（drawio）](../system-arch-image/system-arch.drawio)
- [添加自选股业务流程图（drawio）](../system-arch-image/business-flow-add-stock.drawio)
- [AI分析业务时序图（drawio）](../system-arch-image/sequence-ai-analyze.drawio)

## 2. 业务/数据流分析（MCP Sequential Thinking）

### 2.1 主要业务流程
- 用户启动应用（多平台支持：macOS/Windows/Linux）
- 用户登录/身份验证（如有）
- 股票分组与自选股管理
- 股票数据获取与实时监控
- AI 智能分析与推荐
- 通知与系统托盘交互
- 历史数据与趋势分析
- 用户设置与偏好管理

### 2.2 数据流动关键节点
1. **前端请求**：用户操作触发前端事件，通过 Wails 框架调用后端 Go 方法。
2. **核心逻辑处理**：后端核心模块处理业务逻辑，调用平台抽象层接口。
3. **平台实现层**：根据当前操作系统，调用对应平台实现（如数据目录、系统托盘、通知等）。
4. **数据持久化/外部服务**：本地数据库/文件存储，或外部股票/AI 服务接口。
5. **结果返回前端**：后端处理结果通过 Wails 返回前端，前端渲染展示。

### 2.3 典型业务流示意（伪代码/流程图建议）
- 用户添加自选股 → 前端事件 → Wails 调用 AddStock → 后端校验/存储 → 更新分组列表 → 前端刷新
- 用户请求 AI 分析 → 前端事件 → Wails 调用 AIAnalyze → 后端调用 AI 服务 → 返回分析结果 → 前端展示

> 注：所有流程图、时序图等建议以 drawio 格式绘制，存放于 docs/system-arch-image 目录。

## 3. 系统模块与业务流映射

| 业务流程         | 前端模块      | 后端核心      | 平台抽象层      | 平台实现层      | 外部服务/存储 |
|------------------|---------------|--------------|----------------|----------------|--------------|
| 启动/初始化      | App           | AppCore      | PlatformInit   | init_xxx.go    | 本地配置      |
| 股票分组管理     | GroupManager  | StockGroup   | DataDir        | app_xxx.go     | 本地DB        |
| 股票数据监控     | StockMonitor  | StockCore    | Notifier       | notify_xxx.go  | 股票API       |
| AI 智能分析      | AIAnalyze     | AICore       | AIAdapter      | ai_xxx.go      | AI服务        |
| 通知/托盘        | Tray/Notify   | TrayCore     | TrayAdapter    | tray_xxx.go    | 系统托盘      |
| 用户设置         | Settings      | SettingsCore | ConfigAdapter  | config_xxx.go  | 本地配置      |

> 注：xxx 代表不同平台（darwin, windows, linux）

## 4. 测试与验证策略（含 Playwright）

### 4.1 Playwright 测试集成
- 端到端（E2E）测试：覆盖主要用户操作流程（如添加自选股、AI 分析、通知弹窗等）。
- 跨平台 UI 验证：确保不同平台下界面一致性与功能可用性。
- 自动化回归测试：每次核心功能变更后自动运行。

### 4.2 其他验证点
- 单元测试：Go 后端核心逻辑、平台适配层。
- 集成测试：前后端接口、数据流动。
- 性能与稳定性测试：大数据量、长时间运行。

## 5. 下一步与责任分工

1. **业务/数据流详细建模**：绘制详细流程图、时序图，补充伪代码。（drawio 文件输出至 docs/system-arch-image）
2. **模块责任人指定**：明确每个模块的负责人。
3. **Playwright 脚本开发**：根据业务流编写自动化测试脚本。
4. **文档持续完善**：每阶段开发/重构后同步更新文档。
5. **风险与里程碑管理**：定期评审，跟踪多平台兼容性与关键技术难点。

---

> 本文档为 go-stock 系统架构与业务/数据流分析的顶层计划，后续将根据实际开发进展持续细化和完善。 