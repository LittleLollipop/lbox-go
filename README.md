# lbox-go

lbox-go 是一个轻量级的 Go 工具库,提供了以下功能组件:

## 功能特性

### StateMachine (状态机)
- 线程安全的状态机实现
- 支持状态转换验证
- 支持状态进入/退出钩子
- 完整的状态生命周期管理

### Mission (任务流程控制器)
- 支持多步骤任务流程控制
- 支持任务步骤跳转
- 线程安全的任务执行
- 支持任务暂停和恢复

### ManifoldValve (多路阀门控制器)
- 支持多路数据流控制
- 自动触发数据流汇合处理
- 线程安全的数据处理

## 安装

使用 go get 安装:

    `go get github.com/LittleLollipop/lbox-go`

## 使用示例

### StateMachine 示例
```
    import "github.com/LittleLollipop/lbox-go/pkg/statemachine"

    // 创建游戏状态
    menuState := &GameState{name: "menu"}
    playState := &GameState{name: "play"}

    // 初始化状态机
    stateMap := statemachine.StateMap{
        States: map[string]State{
            "menu": menuState,
            "play": playState,
        },
    }

    sm := statemachine.NewStateMachine(gameStateMachine, stateMap)
    sm.Start("menu")

    // 状态切换
    sm.ChangeState("play")
```
### Mission 示例
```
    import "github.com/LittleLollipop/lbox-go/pkg/mission"

    // 创建任务步骤
    steps := []mission.StepDisposer{
        NewStep("step1"),
        NewStep("step2"),
    }

    // 初始化任务
    m := mission.NewMission(steps)
    m.Start()

    // 执行下一步
    m.GoNext()

    // 跳转到指定步骤
    m.Jump("step2")
```
### ManifoldValve 示例
```
    import "github.com/LittleLollipop/lbox-go/pkg/lbox"

    // 创建阀门控制器
    valves := []string{"valve1", "valve2"}
    mv := lbox.NewManifoldValve(valves, outputHandler)

    // 输入数据
    mv.Input("valve1", data1)
    mv.Input("valve2", data2)
```
## 文档

详细文档请访问: [GoDoc](https://pkg.go.dev/github.com/LittleLollipop/lbox-go)

## 测试

运行所有测试:

    go test ./...

## 贡献

欢迎提交 Pull Request 和 Issue!

1. Fork 本仓库
2. 创建您的特性分支
3. 提交您的更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 作者

- LittleLollipop


