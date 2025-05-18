# Go-Stock 多平台架构设计优化方案

## 1. 引言

本文档是对Go-Stock多平台设计方案的进一步优化，基于软件设计模式对系统架构进行改进和完善。通过应用合适的设计模式，我们可以提高系统的可维护性、扩展性和灵活性，同时保持良好的性能和用户体验。

## 2. 现有架构评估

当前的Go-Stock多平台设计已经提出了良好的分层架构，包括前端层、核心层和平台实现层。然而，通过引入更多先进的设计模式，我们可以解决以下核心问题：

1. **组件间耦合度** - 减少业务逻辑与平台特定代码的耦合
2. **扩展性限制** - 改进系统以便更容易添加新功能和支持新平台
3. **测试难度** - 提高系统可测试性
4. **一致性保障** - 确保跨平台行为一致性

## 3. 设计模式应用

### 3.1 结构型模式

#### 3.1.1 适配器模式 (Adapter Pattern)

**应用场景**：统一不同平台API的接口差异

```go
// 定义统一的通知接口
type NotificationService interface {
    Send(title, message string) error
}

// Windows通知适配器
type WindowsNotificationAdapter struct {
    toaster *toast.Notification
}

func (w *WindowsNotificationAdapter) Send(title, message string) error {
    w.toaster.Title = title
    w.toaster.Message = message
    return w.toaster.Push()
}

// macOS通知适配器
type DarwinNotificationAdapter struct {}

func (d *DarwinNotificationAdapter) Send(title, message string) error {
    cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))
    return cmd.Run()
}
```

#### 3.1.2 桥接模式 (Bridge Pattern)

**应用场景**：分离平台实现与业务抽象

```go
// 定义抽象接口
type UIHandler interface {
    ShowDialog(title, message string) (string, error)
    ShowNotification(title, message string) error
}

// 实现接口
type WindowsUIHandler struct {
    platform *WindowsPlatform
}

type DarwinUIHandler struct {
    platform *DarwinPlatform
}

// 高级业务逻辑使用抽象接口
type AlertService struct {
    uiHandler UIHandler
}

func (a *AlertService) AlertPriceChange(stock *data.StockInfo) error {
    message := fmt.Sprintf("股票 %s 价格变动: %s", stock.Name, stock.Price)
    return a.uiHandler.ShowNotification("价格提醒", message)
}
```

#### 3.1.3 外观模式 (Facade Pattern)

**应用场景**：简化平台API的使用

```go
// 平台服务外观
type PlatformServices struct {
    notification NotificationService
    systemTray   SystemTrayService
    filesystem   FilesystemService
}

// 简化的API
func (p *PlatformServices) NotifyUser(stock *data.StockInfo, eventType string) error {
    title := fmt.Sprintf("%s 提醒", eventType)
    message := GenNotificationMsg(stock)
    return p.notification.Send(title, message)
}

func (p *PlatformServices) UpdateTrayWithProfit(profit float64) error {
    title := fmt.Sprintf("当日收益: %.2f", profit)
    return p.systemTray.UpdateTooltip(title)
}
```

### 3.2 创建型模式

#### 3.2.1 抽象工厂模式 (Abstract Factory Pattern)

**应用场景**：创建平台相关的服务族

```go
// 平台服务工厂接口
type PlatformServiceFactory interface {
    CreateNotificationService() NotificationService
    CreateSystemTrayService() SystemTrayService
    CreateDialogService() DialogService
    CreateFileSystemService() FileSystemService
}

// Windows实现
type WindowsServiceFactory struct {}

func (f *WindowsServiceFactory) CreateNotificationService() NotificationService {
    return &WindowsNotificationService{}
}

// 使用工厂
func NewAppWithPlatformServices(factory PlatformServiceFactory) *App {
    return &App{
        notification: factory.CreateNotificationService(),
        systemTray: factory.CreateSystemTrayService(),
        dialog: factory.CreateDialogService(),
        filesystem: factory.CreateFileSystemService(),
    }
}
```

#### 3.2.2 单例模式 (Singleton Pattern)

**应用场景**：确保平台服务单一实例

```go
var (
    platformInstance Platform
    platformMutex    sync.Mutex
)

// 获取平台实例（线程安全）
func GetPlatform() Platform {
    if platformInstance == nil {
        platformMutex.Lock()
        defer platformMutex.Unlock()
        
        if platformInstance == nil {
            switch runtime.GOOS {
            case "windows":
                platformInstance = newWindowsPlatform()
            case "darwin":
                platformInstance = newDarwinPlatform()
            case "linux":
                platformInstance = newLinuxPlatform()
            default:
                platformInstance = newBasicPlatform()
            }
        }
    }
    return platformInstance
}
```

### 3.3 行为型模式

#### 3.3.1 策略模式 (Strategy Pattern)

**应用场景**：处理不同平台的通知策略

```go
// 通知策略接口
type NotificationStrategy interface {
    Notify(title, message string) error
}

// 策略实现
type ToastNotificationStrategy struct{}
type DarwinNotificationStrategy struct{}
type LinuxNotificationStrategy struct{}

// 通知上下文
type NotificationContext struct {
    strategy NotificationStrategy
}

func (c *NotificationContext) SetStrategy(strategy NotificationStrategy) {
    c.strategy = strategy
}

func (c *NotificationContext) SendNotification(title, message string) error {
    return c.strategy.Notify(title, message)
}
```

#### 3.3.2 观察者模式 (Observer Pattern)

**应用场景**：处理股票价格变化事件

```go
// 观察者接口
type StockObserver interface {
    Update(stock *data.StockInfo)
}

// 具体观察者
type PriceAlertObserver struct {
    notificationService NotificationService
}

func (o *PriceAlertObserver) Update(stock *data.StockInfo) {
    // 检查价格变化是否达到通知条件
    if shouldNotify(stock) {
        o.notificationService.Send("价格提醒", GenNotificationMsg(stock))
    }
}

// 被观察对象
type StockMonitor struct {
    observers []StockObserver
    stocks    map[string]*data.StockInfo
}

func (m *StockMonitor) AddObserver(observer StockObserver) {
    m.observers = append(m.observers, observer)
}

func (m *StockMonitor) UpdateStock(stock *data.StockInfo) {
    oldStock := m.stocks[stock.Code]
    m.stocks[stock.Code] = stock
    
    if hasSignificantChange(oldStock, stock) {
        for _, observer := range m.observers {
            observer.Update(stock)
        }
    }
}
```

#### 3.3.3 命令模式 (Command Pattern)

**应用场景**：封装系统托盘菜单操作

```go
// 命令接口
type Command interface {
    Execute() error
}

// 具体命令
type RefreshDataCommand struct {
    app *App
}

func (c *RefreshDataCommand) Execute() error {
    return c.app.refreshStockData()
}

type ShowWindowCommand struct {
    app *App
}

func (c *ShowWindowCommand) Execute() error {
    return c.app.showMainWindow()
}

// 系统托盘菜单使用命令
func setupSystemTrayMenu(app *App, tray SystemTrayService) {
    refreshCmd := &RefreshDataCommand{app: app}
    showCmd := &ShowWindowCommand{app: app}
    
    tray.AddMenuItem("刷新数据", func() { refreshCmd.Execute() })
    tray.AddMenuItem("显示窗口", func() { showCmd.Execute() })
}
```

## 4. 改进的架构设计

### 4.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                       用户界面层 (Frontend)                      │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                         应用层 (Application)                     │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   服务门面      │  │    命令处理器   │  │   观察者管理器  │  │
│  │   (Facade)      │  │   (Commands)    │  │   (Observers)   │  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  │
└───────────┼─────────────────────┼─────────────────────┼─────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                        领域层 (Domain)                          │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   业务服务      │  │    数据模型     │  │   领域事件      │  │
│  │  (Services)     │  │    (Models)     │  │    (Events)     │  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  │
└───────────┼─────────────────────┼─────────────────────┼─────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                      基础设施层 (Infrastructure)                 │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   平台抽象 (Platform Abstraction)           ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  ││
│  │  │ 策略接口    │  │ 适配器接口  │  │  抽象工厂接口       │  ││
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘  ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   平台实现 (Platform Implementation)        ││
│  │                                                             ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  ││
│  │  │ Windows实现 │  │ macOS实现   │  │  Linux实现          │  ││
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘  ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 主要改进点

1. **引入应用层**：
   - 应用层作为领域层和用户界面之间的中介
   - 包含应用服务、命令处理和事件观察者

2. **领域层聚焦业务逻辑**：
   - 纯粹的业务规则和业务对象
   - 不依赖于平台特定实现

3. **基础设施层优化**：
   - 包含策略、适配器和工厂
   - 平台实现与抽象完全分离

4. **依赖方向**：
   - 确保依赖方向向下流动
   - 上层依赖下层的抽象，而非实现

## 5. 关键接口和实现

### 5.1 平台服务抽象工厂

```go
// 平台服务抽象工厂
type PlatformServiceFactory interface {
    // 基础服务
    CreateLogger() Logger
    CreateFileSystem() FileSystem
    CreateNetwork() Network
    
    // UI服务
    CreateNotification() NotificationService
    CreateDialog() DialogService
    CreateSystemTray() SystemTrayService
    
    // 系统服务
    CreateScreenInfo() ScreenInfoService
    CreateProcessManager() ProcessManagerService
    CreateBrowserDetector() BrowserDetectorService
}

// 创建平台特定工厂
func CreatePlatformFactory() PlatformServiceFactory {
    switch runtime.GOOS {
    case "windows":
        return &WindowsServiceFactory{}
    case "darwin":
        return &DarwinServiceFactory{}
    case "linux":
        return &LinuxServiceFactory{}
    default:
        return &BasicServiceFactory{}
    }
}
```

### 5.2 应用层服务门面

```go
// 应用服务门面
type AppServiceFacade struct {
    stockService       *StockService
    notificationService NotificationService
    alertService       *AlertService
    platformServices   map[string]interface{}
}

// 创建服务门面
func NewAppServiceFacade(factory PlatformServiceFactory) *AppServiceFacade {
    notification := factory.CreateNotification()
    fileSystem := factory.CreateFileSystem()
    systemTray := factory.CreateSystemTray()
    
    // 创建应用服务
    stockService := NewStockService()
    alertService := NewAlertService(notification)
    
    return &AppServiceFacade{
        stockService: stockService,
        notificationService: notification,
        alertService: alertService,
        platformServices: map[string]interface{}{
            "notification": notification,
            "fileSystem": fileSystem,
            "systemTray": systemTray,
        },
    }
}

// 暴露给前端的方法
func (a *AppServiceFacade) MonitorStock(stockCode string) error {
    stock, err := a.stockService.GetStockInfo(stockCode)
    if err != nil {
        return err
    }
    
    // 添加监控
    a.stockService.AddMonitor(stock)
    return nil
}
```

### 5.3 命令处理器

```go
// 命令处理器
type CommandHandler struct {
    commands map[string]Command
}

// 注册命令
func (h *CommandHandler) Register(name string, cmd Command) {
    h.commands[name] = cmd
}

// 执行命令
func (h *CommandHandler) Execute(name string, params ...interface{}) error {
    cmd, exists := h.commands[name]
    if !exists {
        return fmt.Errorf("command not found: %s", name)
    }
    
    if paramCmd, ok := cmd.(ParameterizedCommand); ok {
        return paramCmd.ExecuteWithParams(params...)
    }
    
    return cmd.Execute()
}
```

### 5.4 事件观察者管理器

```go
// 事件类型
type EventType string

const (
    StockPriceChanged EventType = "stock.price.changed"
    MarketOpened      EventType = "market.opened"
    MarketClosed      EventType = "market.closed"
)

// 事件观察者
type EventObserver interface {
    Handle(event Event)
}

// 事件
type Event interface {
    Type() EventType
    Payload() interface{}
    Timestamp() time.Time
}

// 观察者管理器
type EventManager struct {
    observers map[EventType][]EventObserver
}

// 注册观察者
func (m *EventManager) RegisterObserver(eventType EventType, observer EventObserver) {
    m.observers[eventType] = append(m.observers[eventType], observer)
}

// 触发事件
func (m *EventManager) Emit(event Event) {
    observers, exists := m.observers[event.Type()]
    if !exists {
        return
    }
    
    for _, observer := range observers {
        observer.Handle(event)
    }
}
```

## 6. 优化后的核心服务实现

### 6.1 股票服务

```go
// 股票服务（领域服务）
type StockService struct {
    stockRepository StockRepository
    eventManager    *EventManager
}

// 获取股票信息
func (s *StockService) GetStockInfo(code string) (*data.StockInfo, error) {
    return s.stockRepository.GetByCode(code)
}

// 更新股票价格
func (s *StockService) UpdatePrice(code string, price float64) error {
    stock, err := s.stockRepository.GetByCode(code)
    if err != nil {
        return err
    }
    
    oldPrice := stock.Price
    stock.Price = price
    stock.UpdateTime = time.Now()
    
    // 保存更新
    if err := s.stockRepository.Save(stock); err != nil {
        return err
    }
    
    // 如果价格有变化，发出事件
    if oldPrice != price {
        s.eventManager.Emit(&StockPriceChangedEvent{
            stock: stock,
            oldPrice: oldPrice,
            newPrice: price,
        })
    }
    
    return nil
}
```

### 6.2 通知服务策略

```go
// 通知配置
type NotificationConfig struct {
    EnableSound bool
    Timeout     time.Duration
    Priority    int
}

// 通知策略接口
type NotificationStrategy interface {
    ShowNotification(title, message string, config NotificationConfig) error
}

// Windows通知策略
type WindowsNotificationStrategy struct{}

func (s *WindowsNotificationStrategy) ShowNotification(title, message string, config NotificationConfig) error {
    notification := toast.Notification{
        AppID:   "go-stock",
        Title:   title,
        Message: message,
        Audio:   toast.Silent,
    }
    
    if config.EnableSound {
        notification.Audio = toast.Default
    }
    
    return notification.Push()
}

// macOS通知策略
type DarwinNotificationStrategy struct{}

func (s *DarwinNotificationStrategy) ShowNotification(title, message string, config NotificationConfig) error {
    scriptArgs := []string{"-e"}
    
    scriptCmd := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
    if config.EnableSound {
        scriptCmd += " sound name \"default\""
    }
    
    cmd := exec.Command("osascript", scriptArgs...)
    return cmd.Run()
}
```

## 7. 实现建议

### 7.1 实现路径

1. **核心抽象层实现**：
   - 首先实现平台抽象接口和工厂
   - 实现基础的领域模型和服务

2. **平台适配实现**：
   - 为每个目标平台实现具体适配器
   - 实现平台特定的策略类

3. **应用服务实现**：
   - 实现命令处理器和服务门面
   - 实现事件管理和观察者

### 7.2 代码组织

```
go-stock/
├── cmd/
│   └── go-stock/
│       └── main.go                 // 主程序入口
├── internal/
│   ├── app/                        // 应用层
│   │   ├── commands/               // 命令模式实现
│   │   ├── facade/                 // 门面模式实现
│   │   ├── observers/              // 观察者实现
│   │   └── services/               // 应用服务
│   ├── domain/                     // 领域层
│   │   ├── models/                 // 领域模型
│   │   ├── events/                 // 领域事件
│   │   ├── repositories/           // 仓储接口
│   │   └── services/               // 领域服务
│   └── platform/                   // 平台层
│       ├── factory.go              // 平台工厂
│       ├── interfaces/             // 平台接口
│       ├── common/                 // 通用实现
│       ├── windows/                // Windows实现
│       ├── darwin/                 // macOS实现
│       └── linux/                  // Linux实现
├── pkg/
│   ├── logger/                     // 日志包
│   ├── utils/                      // 通用工具
│   └── wails/                      // Wails相关
└── frontend/                       // 前端代码
    ├── src/
    └── dist/
```

### 7.3 测试策略

1. **单元测试**：
   - 使用接口模拟测试领域服务
   - 使用依赖注入简化测试

2. **集成测试**：
   - 使用测试工厂创建完整的服务图
   - 使用测试双代替外部依赖

3. **平台特定测试**：
   - 为每个平台创建特定的测试套件
   - 使用条件编译隔离平台特定测试

## 8. 总结

本设计优化方案通过应用多种设计模式，从根本上改进了Go-Stock的多平台架构设计。主要改进包括：

1. **清晰的职责分离**：应用层、领域层和基础设施层各司其职
2. **灵活的扩展机制**：使用策略模式和工厂模式支持新平台
3. **低耦合设计**：通过接口分离原则降低组件间耦合
4. **提高可测试性**：便于模拟和测试的依赖注入设计

通过这些改进，Go-Stock将获得更高的代码质量、更好的可维护性和更强的扩展能力，同时保持良好的性能和用户体验。 