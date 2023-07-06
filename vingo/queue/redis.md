### 初始化
```go
    // 在主进程中执行初始化
    queue.InitRedisQueue(nil)
```

### 定义消费类
```go

type Handle struct{}
func (s *Handle) HandleMessage(message *string) {
    // 如果处理失败抛出panic异常，默认5秒后重新消费
    fmt.Println(*message)
}

// 使用方法
h := Handle{}
queue.Redis.StartMonitor("test", &h) // 实时队列协程
queue.Redis.StartMonitorDelay("test") // 延迟队列协程

// 可以开启多个不同的消费主题队列协程

```