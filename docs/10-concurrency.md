# Go 併發處理與任務隊列

> 本文件說明 Go 語言的併發處理機制，包含 Goroutines、Channels 和 Worker Pool 模式。

## 目錄

- [基本概念](#基本概念)
- [Goroutines](#goroutines)
- [Channels](#channels)
- [Worker Pool 模式](#worker-pool-模式)
- [完整範例程式碼](#完整範例程式碼)
- [進階：優雅關閉](#進階優雅關閉)
- [排程任務（Task Scheduling）](#排程任務task-scheduling)
- [第三方任務隊列庫](#第三方任務隊列庫)
- [RabbitMQ 整合](#rabbitmq-整合)

---

## 基本概念

Go 語言的併發模型基於 CSP（Communicating Sequential Processes），核心理念是：

> **不要透過共享記憶體來通信，而是透過通信來共享記憶體。**

| 概念 | 說明 |
|------|------|
| Goroutine | 輕量級線程，由 Go runtime 管理 |
| Channel | Goroutines 之間的通信管道 |
| select | 多路 channel 選擇器 |
| sync 包 | 提供互斥鎖、WaitGroup 等同步原語 |

---

## Goroutines

Goroutine 是 Go 併發的基礎，比系統線程更輕量（初始只佔 2KB 記憶體）。

### 基本用法

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    // 一般呼叫
    sayHello("World")

    // 啟動 goroutine（非同步執行）
    go sayHello("Goroutine")

    // 匿名函數也可以
    go func() {
        fmt.Println("匿名 goroutine")
    }()

    // 等待 goroutines 完成（簡單方式）
    time.Sleep(time.Second)
}
```

### 使用 sync.WaitGroup

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1) // 增加計數

        go func(id int) {
            defer wg.Done() // 完成時減少計數
            fmt.Printf("Goroutine %d 執行中\n", id)
        }(i)
    }

    wg.Wait() // 等待所有 goroutines 完成
    fmt.Println("全部完成！")
}
```

---

## Channels

Channel 是 goroutines 之間安全通信的管道。

### 基本用法

```go
package main

import "fmt"

func main() {
    // 建立無緩衝 channel
    ch := make(chan string)

    // 在 goroutine 中發送資料
    go func() {
        ch <- "Hello from channel!"
    }()

    // 接收資料（會阻塞直到有資料）
    msg := <-ch
    fmt.Println(msg)
}
```

### 緩衝 Channel

```go
package main

import "fmt"

func main() {
    // 建立有緩衝的 channel（容量為 3）
    ch := make(chan int, 3)

    // 可以連續發送 3 個而不阻塞
    ch <- 1
    ch <- 2
    ch <- 3

    // 接收
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}
```

### Channel 方向

```go
// 只能發送的 channel
func sender(ch chan<- string) {
    ch <- "data"
}

// 只能接收的 channel
func receiver(ch <-chan string) {
    msg := <-ch
    fmt.Println(msg)
}
```

### 使用 range 遍歷 Channel

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 5)

    // 發送資料
    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
        }
        close(ch) // 關閉 channel
    }()

    // 使用 range 接收（直到 channel 關閉）
    for num := range ch {
        fmt.Println(num)
    }
}
```

### Select 多路選擇

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "來自 channel 1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "來自 channel 2"
    }()

    // 等待多個 channel
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        case <-time.After(3 * time.Second):
            fmt.Println("超時！")
        }
    }
}
```

---

## Worker Pool 模式

Worker Pool 是實現任務隊列的經典模式，適用於：
- 控制併發數量
- 處理大量任務
- 資源管理

### 架構圖

```
                    ┌──────────┐
                    │  Worker 1 │──┐
┌────────┐         ├──────────┤  │     ┌─────────┐
│  Jobs  │────────▶│  Worker 2 │──┼────▶│ Results │
│ Channel│         ├──────────┤  │     │ Channel │
└────────┘         │  Worker 3 │──┘     └─────────┘
                    └──────────┘
```

### 基本實現

```go
package main

import (
    "fmt"
    "time"
)

// Job 代表一個任務
type Job struct {
    ID   int
    Data string
}

// Result 代表任務結果
type Result struct {
    JobID  int
    Output string
}

// worker 處理任務的工作者
func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        fmt.Printf("Worker %d 開始處理任務 %d\n", id, job.ID)

        // 模擬處理時間
        time.Sleep(time.Second)

        // 回傳結果
        results <- Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("任務 %d 處理完成，資料: %s", job.ID, job.Data),
        }
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 10

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    // 啟動 workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // 發送任務
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{
            ID:   j,
            Data: fmt.Sprintf("資料_%d", j),
        }
    }
    close(jobs) // 關閉任務 channel

    // 收集結果
    for r := 1; r <= numJobs; r++ {
        result := <-results
        fmt.Println(result.Output)
    }
}
```

---

## 完整範例程式碼

以下是一個更完整的 Worker Pool 實現，包含錯誤處理：

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// ============================================
// 任務與結果定義
// ============================================

// Task 代表一個待處理的任務
type Task struct {
    ID      int
    Payload interface{}
}

// TaskResult 代表任務處理結果
type TaskResult struct {
    TaskID int
    Data   interface{}
    Err    error
}

// ============================================
// Worker Pool 實現
// ============================================

// WorkerPool 任務池
type WorkerPool struct {
    workerCount int
    tasks       chan Task
    results     chan TaskResult
    wg          sync.WaitGroup
}

// NewWorkerPool 建立新的 Worker Pool
func NewWorkerPool(workerCount, taskBufferSize int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        tasks:       make(chan Task, taskBufferSize),
        results:     make(chan TaskResult, taskBufferSize),
    }
}

// Start 啟動 Worker Pool
func (wp *WorkerPool) Start(ctx context.Context, processor func(Task) TaskResult) {
    for i := 1; i <= wp.workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx, i, processor)
    }
}

// worker 執行任務的工作者
func (wp *WorkerPool) worker(ctx context.Context, id int, processor func(Task) TaskResult) {
    defer wp.wg.Done()

    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: 收到停止信號\n", id)
            return
        case task, ok := <-wp.tasks:
            if !ok {
                fmt.Printf("Worker %d: 任務 channel 已關閉\n", id)
                return
            }
            fmt.Printf("Worker %d: 處理任務 %d\n", id, task.ID)
            result := processor(task)
            wp.results <- result
        }
    }
}

// Submit 提交任務
func (wp *WorkerPool) Submit(task Task) {
    wp.tasks <- task
}

// Results 取得結果 channel
func (wp *WorkerPool) Results() <-chan TaskResult {
    return wp.results
}

// Close 關閉任務 channel
func (wp *WorkerPool) Close() {
    close(wp.tasks)
}

// Wait 等待所有 workers 完成
func (wp *WorkerPool) Wait() {
    wp.wg.Wait()
    close(wp.results)
}

// ============================================
// 使用範例
// ============================================

func main() {
    rand.Seed(time.Now().UnixNano())

    // 建立 Worker Pool（3 個 workers，緩衝 100 個任務）
    pool := NewWorkerPool(3, 100)

    // 定義任務處理邏輯
    processor := func(task Task) TaskResult {
        // 模擬處理時間（100ms - 500ms）
        processingTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
        time.Sleep(processingTime)

        // 模擬偶爾失敗
        if rand.Float32() < 0.1 { // 10% 機率失敗
            return TaskResult{
                TaskID: task.ID,
                Err:    errors.New("處理失敗：隨機錯誤"),
            }
        }

        return TaskResult{
            TaskID: task.ID,
            Data:   fmt.Sprintf("任務 %d 處理完成，耗時 %v", task.ID, processingTime),
        }
    }

    // 使用 context 控制（可選：設定超時）
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 啟動 Worker Pool
    pool.Start(ctx, processor)

    // 提交 20 個任務
    numTasks := 20
    go func() {
        for i := 1; i <= numTasks; i++ {
            pool.Submit(Task{
                ID:      i,
                Payload: fmt.Sprintf("任務資料 %d", i),
            })
        }
        pool.Close() // 提交完畢後關閉
    }()

    // 收集結果
    var successCount, failCount int
    go func() {
        for result := range pool.Results() {
            if result.Err != nil {
                failCount++
                fmt.Printf("❌ 任務 %d 失敗: %v\n", result.TaskID, result.Err)
            } else {
                successCount++
                fmt.Printf("✅ %v\n", result.Data)
            }
        }
    }()

    // 等待完成
    pool.Wait()

    fmt.Println("\n========== 統計 ==========")
    fmt.Printf("總任務數: %d\n", numTasks)
    fmt.Printf("成功: %d\n", successCount)
    fmt.Printf("失敗: %d\n", failCount)
}
```

---

## 進階：優雅關閉

使用 `context` 和 `os/signal` 實現優雅關閉：

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 建立可取消的 context
    ctx, cancel := context.WithCancel(context.Background())

    // 監聽系統信號
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // 啟動工作
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("收到停止信號，正在清理...")
                return
            default:
                fmt.Println("工作中...")
                time.Sleep(time.Second)
            }
        }
    }()

    // 等待停止信號
    sig := <-sigChan
    fmt.Printf("\n收到信號: %v\n", sig)

    // 取消 context，通知所有 goroutines 停止
    cancel()

    // 給予清理時間
    time.Sleep(2 * time.Second)
    fmt.Println("程式結束")
}
```

---

## 排程任務（Task Scheduling）

Go 可以在程式內建排程功能，類似 Laravel 的 Task Scheduling。

### 方式一：time.Ticker（簡單間隔）

適合固定間隔執行的任務：

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 每 5 秒執行一次
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    // 一次性延遲執行
    time.AfterFunc(10*time.Second, func() {
        fmt.Println("10 秒後執行一次")
    })

    for {
        select {
        case t := <-ticker.C:
            fmt.Println("定時任務執行於:", t)
            // 執行你的任務
        }
    }
}
```

### 方式二：robfig/cron（推薦）

功能完整的 Cron 排程庫，語法與 Linux crontab 相同：

```bash
go get github.com/robfig/cron/v3
```

```go
package main

import (
    "fmt"
    "time"

    "github.com/robfig/cron/v3"
)

func main() {
    // 建立 cron 調度器
    c := cron.New()

    // ============================================
    // Cron 表達式語法：分 時 日 月 週
    // ============================================

    // 每分鐘執行
    c.AddFunc("* * * * *", func() {
        fmt.Println("[每分鐘] 執行於:", time.Now().Format("15:04:05"))
    })

    // 每小時的第 0 分鐘（整點）
    c.AddFunc("0 * * * *", func() {
        fmt.Println("[每小時] 整點任務")
    })

    // 每天凌晨 2:30
    c.AddFunc("30 2 * * *", func() {
        fmt.Println("[每天] 凌晨 2:30 執行")
    })

    // 每週一早上 9 點
    c.AddFunc("0 9 * * 1", func() {
        fmt.Println("[每週一] 早上 9 點執行")
    })

    // 每月 1 號凌晨 0 點
    c.AddFunc("0 0 1 * *", func() {
        fmt.Println("[每月] 1 號執行")
    })

    // ============================================
    // 擴展語法（更直觀）
    // ============================================

    // 每 30 秒
    c.AddFunc("@every 30s", func() {
        fmt.Println("[每30秒]")
    })

    // 每 5 分鐘
    c.AddFunc("@every 5m", func() {
        fmt.Println("[每5分鐘]")
    })

    // 每小時
    c.AddFunc("@hourly", func() {
        fmt.Println("[每小時] hourly")
    })

    // 每天午夜
    c.AddFunc("@daily", func() {
        fmt.Println("[每天] daily")
    })

    // 每週日午夜
    c.AddFunc("@weekly", func() {
        fmt.Println("[每週] weekly")
    })

    // 啟動調度器
    c.Start()

    fmt.Println("排程器已啟動，按 Ctrl+C 停止...")

    // 保持程式運行
    select {}
}
```

### 方式三：支援秒級的 Cron

預設 cron 只支援到分鐘，如需秒級精度：

```go
// 建立支援秒的調度器
c := cron.New(cron.WithSeconds())

// 語法變成：秒 分 時 日 月 週
c.AddFunc("*/5 * * * * *", func() {
    fmt.Println("每 5 秒執行")
})

c.AddFunc("0 0 * * * *", func() {
    fmt.Println("每小時整點")
})
```

### 完整排程器範例（含優雅關閉）

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/robfig/cron/v3"
)

// Scheduler 排程管理器
type Scheduler struct {
    cron *cron.Cron
}

// NewScheduler 建立排程器
func NewScheduler() *Scheduler {
    return &Scheduler{
        cron: cron.New(cron.WithSeconds()), // 支援秒級
    }
}

// RegisterJobs 註冊所有排程任務
func (s *Scheduler) RegisterJobs() {
    // 每分鐘清理暫存檔
    s.cron.AddFunc("0 * * * * *", func() {
        fmt.Printf("[%s] 清理暫存檔...\n", time.Now().Format("15:04:05"))
        cleanTempFiles()
    })

    // 每 5 分鐘同步資料
    s.cron.AddFunc("0 */5 * * * *", func() {
        fmt.Printf("[%s] 同步資料...\n", time.Now().Format("15:04:05"))
        syncData()
    })

    // 每天凌晨 3 點備份資料庫
    s.cron.AddFunc("0 0 3 * * *", func() {
        fmt.Printf("[%s] 備份資料庫...\n", time.Now().Format("15:04:05"))
        backupDatabase()
    })

    // 每週日凌晨 4 點產生週報
    s.cron.AddFunc("0 0 4 * * 0", func() {
        fmt.Printf("[%s] 產生週報...\n", time.Now().Format("15:04:05"))
        generateWeeklyReport()
    })
}

// Start 啟動排程器
func (s *Scheduler) Start() {
    s.cron.Start()
    fmt.Println("排程器已啟動")
}

// Stop 停止排程器
func (s *Scheduler) Stop() context.Context {
    fmt.Println("正在停止排程器...")
    return s.cron.Stop()
}

// 任務函數
func cleanTempFiles()       { fmt.Println("  → 暫存檔已清理") }
func syncData()             { fmt.Println("  → 資料同步完成") }
func backupDatabase()       { fmt.Println("  → 資料庫備份完成") }
func generateWeeklyReport() { fmt.Println("  → 週報已產生") }

func main() {
    scheduler := NewScheduler()
    scheduler.RegisterJobs()
    scheduler.Start()

    // 監聽系統信號
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("排程器運行中，按 Ctrl+C 停止...")
    <-sigChan

    // 優雅關閉
    ctx := scheduler.Stop()
    <-ctx.Done()
    fmt.Println("排程器已停止")
}
```

### Cron 表達式速查表

```
┌──────────── 秒 (0-59)（需啟用 WithSeconds）
│ ┌────────── 分 (0-59)
│ │ ┌──────── 時 (0-23)
│ │ │ ┌────── 日 (1-31)
│ │ │ │ ┌──── 月 (1-12)
│ │ │ │ │ ┌── 週 (0-6, 0=週日)
│ │ │ │ │ │
* * * * * *
```

| 符號 | 說明 | 範例 |
|------|------|------|
| `*` | 任意值 | `* * * * *` 每分鐘 |
| `,` | 列舉 | `0,30 * * * *` 第0和30分 |
| `-` | 範圍 | `0-5 * * * *` 0到5分 |
| `/` | 間隔 | `*/5 * * * *` 每5分鐘 |

### 常用排程對照表

| 需求 | Cron 表達式 | 擴展語法 |
|------|-------------|----------|
| 每分鐘 | `* * * * *` | `@every 1m` |
| 每 5 分鐘 | `*/5 * * * *` | `@every 5m` |
| 每小時 | `0 * * * *` | `@hourly` |
| 每天午夜 | `0 0 * * *` | `@daily` |
| 每天早上 8 點 | `0 8 * * *` | - |
| 每週日午夜 | `0 0 * * 0` | `@weekly` |
| 每月 1 號 | `0 0 1 * *` | `@monthly` |
| 每年 1/1 | `0 0 1 1 *` | `@yearly` |

### 與 Laravel 排程對照

| Laravel | Go (robfig/cron) |
|---------|------------------|
| `->everyMinute()` | `* * * * *` |
| `->everyFiveMinutes()` | `*/5 * * * *` |
| `->hourly()` | `@hourly` |
| `->daily()` | `@daily` |
| `->dailyAt('13:00')` | `0 13 * * *` |
| `->weekly()` | `@weekly` |
| `->monthly()` | `@monthly` |
| `->cron('* * * * *')` | 直接使用表達式 |
| `php artisan schedule:run` | 程式內建，持續運行 |

---

## 第三方任務隊列庫

對於生產環境的複雜任務隊列需求，推薦使用以下第三方庫：

### 1. Asynq（推薦）

基於 Redis 的分散式任務隊列，類似 Ruby 的 Sidekiq。

```bash
go get github.com/hibiken/asynq
```

```go
// 定義任務
func NewEmailTask(to string) (*asynq.Task, error) {
    payload, _ := json.Marshal(map[string]string{"to": to})
    return asynq.NewTask("email:send", payload), nil
}

// 發送任務
client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
task, _ := NewEmailTask("user@example.com")
client.Enqueue(task)

// 處理任務（另一個程序）
srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: "localhost:6379"},
    asynq.Config{Concurrency: 10},
)
mux := asynq.NewServeMux()
mux.HandleFunc("email:send", handleEmailTask)
srv.Run(mux)
```

### 2. Machinery

支援多種 broker（Redis、RabbitMQ、MongoDB）。

```bash
go get github.com/RichardKnop/machinery/v2
```

### 3. Go-queue（go-zero 生態）

```bash
go get github.com/zeromicro/go-queue
```

---

## RabbitMQ 整合

RabbitMQ 是企業級的訊息佇列，結合 Go 的多線程可以實現高效能的分散式任務處理。

### 架構圖

```
                           RabbitMQ Server
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│   │  Exchange   │───▶│   Queue 1   │    │   Queue 2   │     │
│   │  (路由器)    │    │ (email任務) │    │ (圖片處理)  │     │
│   └─────────────┘    └─────────────┘    └─────────────┘     │
│         ▲                   │                  │             │
└─────────│───────────────────│──────────────────│─────────────┘
          │                   │                  │
          │                   ▼                  ▼
┌─────────┴───────┐   ┌───────────────────────────────────────┐
│    Producer     │   │          Consumer (Go App)            │
│   (發送任務)     │   │                                       │
│                 │   │  ┌─────────────────────────────────┐  │
│  - Web API      │   │  │       Message Handler           │  │
│  - 排程任務      │   │  │       (接收訊息)                 │  │
│  - 其他服務      │   │  └────────────┬────────────────────┘  │
│                 │   │               │                       │
└─────────────────┘   │               ▼                       │
                      │  ┌─────────────────────────────────┐  │
                      │  │         Worker Pool             │  │
                      │  │  ┌────────┐ ┌────────┐          │  │
                      │  │  │Worker 1│ │Worker 2│ ...      │  │
                      │  │  └────────┘ └────────┘          │  │
                      │  └─────────────────────────────────┘  │
                      │               │                       │
                      │               ▼                       │
                      │         處理完成 → ACK                │
                      └───────────────────────────────────────┘
```

### 訊息流程

```
1. Producer 發送訊息到 Exchange
              │
              ▼
2. Exchange 根據 routing key 路由到對應 Queue
              │
              ▼
3. Queue 暫存訊息（持久化）
              │
              ▼
4. Consumer 從 Queue 取得訊息
              │
              ▼
5. 訊息分發到 Worker Pool（多個 Goroutines 並行處理）
              │
              ▼
6. Worker 處理完成後發送 ACK
              │
              ▼
7. RabbitMQ 移除該訊息
```

### 安裝依賴

```bash
go get github.com/rabbitmq/amqp091-go
```

### 基本 Consumer 範例

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    // 連接 RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal("連接失敗:", err)
    }
    defer conn.Close()

    // 建立 Channel
    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("建立 Channel 失敗:", err)
    }
    defer ch.Close()

    // 宣告 Queue
    q, err := ch.QueueDeclare(
        "task_queue", // Queue 名稱
        true,         // 持久化
        false,        // 自動刪除
        false,        // 獨占
        false,        // 不等待
        nil,          // 參數
    )
    if err != nil {
        log.Fatal("宣告 Queue 失敗:", err)
    }

    // 設定 prefetch（每個 worker 一次取幾個訊息）
    // 這是控制併發的關鍵！
    err = ch.Qos(
        1,     // prefetch count（每個 consumer 一次取 1 個）
        0,     // prefetch size
        false, // global
    )
    if err != nil {
        log.Fatal("設定 Qos 失敗:", err)
    }

    // 開始消費訊息
    msgs, err := ch.Consume(
        q.Name, // Queue
        "",     // Consumer tag
        false,  // Auto-ack（設為 false，手動確認）
        false,  // Exclusive
        false,  // No-local
        false,  // No-wait
        nil,    // Args
    )
    if err != nil {
        log.Fatal("註冊 Consumer 失敗:", err)
    }

    // 啟動多個 Workers
    numWorkers := 5
    var wg sync.WaitGroup

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, msgs, &wg)
    }

    fmt.Printf("🚀 等待訊息中，啟動 %d 個 workers...\n", numWorkers)
    fmt.Println("按 Ctrl+C 停止")

    wg.Wait()
}

func worker(id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup) {
    defer wg.Done()

    for msg := range msgs {
        fmt.Printf("Worker %d: 收到訊息 [%s]\n", id, msg.Body)

        // 模擬處理時間
        processTime := time.Duration(len(msg.Body)) * 100 * time.Millisecond
        time.Sleep(processTime)

        fmt.Printf("Worker %d: 處理完成 ✓\n", id)

        // 確認訊息（ACK）
        msg.Ack(false)
    }
}
```

### Producer 範例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    defer conn.Close()

    ch, _ := conn.Channel()
    defer ch.Close()

    q, _ := ch.QueueDeclare("task_queue", true, false, false, false, nil)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 發送 10 個任務
    for i := 1; i <= 10; i++ {
        body := fmt.Sprintf("任務 #%d", i)
        err := ch.PublishWithContext(ctx,
            "",     // Exchange
            q.Name, // Routing key
            false,  // Mandatory
            false,  // Immediate
            amqp.Publishing{
                DeliveryMode: amqp.Persistent, // 訊息持久化
                ContentType:  "text/plain",
                Body:         []byte(body),
            },
        )
        if err != nil {
            log.Printf("發送失敗: %v", err)
        } else {
            fmt.Printf("已發送: %s\n", body)
        }
    }
}
```

### 完整 Worker Pool + RabbitMQ 範例

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

// ============================================
// 任務定義
// ============================================

// EmailTask 郵件任務
type EmailTask struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// ============================================
// RabbitMQ Consumer
// ============================================

type Consumer struct {
    conn       *amqp.Connection
    channel    *amqp.Channel
    queueName  string
    numWorkers int
    wg         sync.WaitGroup
}

func NewConsumer(url, queueName string, numWorkers int) (*Consumer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("連接失敗: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("建立 Channel 失敗: %w", err)
    }

    // 宣告 Queue
    _, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, fmt.Errorf("宣告 Queue 失敗: %w", err)
    }

    // 設定 QoS
    err = ch.Qos(1, 0, false)
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, fmt.Errorf("設定 Qos 失敗: %w", err)
    }

    return &Consumer{
        conn:       conn,
        channel:    ch,
        queueName:  queueName,
        numWorkers: numWorkers,
    }, nil
}

func (c *Consumer) Start(ctx context.Context, handler func([]byte) error) error {
    msgs, err := c.channel.Consume(
        c.queueName, "", false, false, false, false, nil,
    )
    if err != nil {
        return fmt.Errorf("註冊 Consumer 失敗: %w", err)
    }

    // 啟動 Workers
    for i := 1; i <= c.numWorkers; i++ {
        c.wg.Add(1)
        go c.worker(ctx, i, msgs, handler)
    }

    fmt.Printf("🚀 Consumer 已啟動，%d 個 workers 等待訊息...\n", c.numWorkers)
    return nil
}

func (c *Consumer) worker(ctx context.Context, id int, msgs <-chan amqp.Delivery, handler func([]byte) error) {
    defer c.wg.Done()

    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: 收到停止信號\n", id)
            return
        case msg, ok := <-msgs:
            if !ok {
                fmt.Printf("Worker %d: Channel 已關閉\n", id)
                return
            }

            fmt.Printf("Worker %d: 處理訊息...\n", id)

            if err := handler(msg.Body); err != nil {
                fmt.Printf("Worker %d: 處理失敗 - %v\n", id, err)
                // 可選：重新入列或發送到死信隊列
                msg.Nack(false, true) // requeue
            } else {
                fmt.Printf("Worker %d: 處理成功 ✓\n", id)
                msg.Ack(false)
            }
        }
    }
}

func (c *Consumer) Stop() {
    c.channel.Close()
    c.conn.Close()
    c.wg.Wait()
    fmt.Println("Consumer 已停止")
}

// ============================================
// 主程式
// ============================================

func main() {
    // 建立 Consumer（5 個 workers）
    consumer, err := NewConsumer(
        "amqp://guest:guest@localhost:5672/",
        "email_queue",
        5,
    )
    if err != nil {
        log.Fatal(err)
    }

    // 定義任務處理邏輯
    handler := func(body []byte) error {
        var task EmailTask
        if err := json.Unmarshal(body, &task); err != nil {
            return fmt.Errorf("解析失敗: %w", err)
        }

        // 模擬發送郵件
        fmt.Printf("  📧 發送郵件到: %s\n", task.To)
        fmt.Printf("     主題: %s\n", task.Subject)
        time.Sleep(500 * time.Millisecond) // 模擬處理時間

        return nil
    }

    // 使用 context 控制生命週期
    ctx, cancel := context.WithCancel(context.Background())

    // 啟動 Consumer
    if err := consumer.Start(ctx, handler); err != nil {
        log.Fatal(err)
    }

    // 監聽系統信號
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan
    fmt.Println("\n正在優雅關閉...")

    cancel()
    consumer.Stop()
}
```

### Docker Compose（RabbitMQ）

```yaml
version: '3.8'
services:
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # 管理介面
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  rabbitmq_data:
```

### RabbitMQ vs 其他方案比較

| 特性 | RabbitMQ | Redis (Asynq) | Channel |
|------|----------|---------------|---------|
| 分散式 | ✅ | ✅ | ❌ |
| 持久化 | ✅ | ✅ | ❌ |
| 訊息確認 | ✅ ACK/NACK | ✅ | ❌ |
| 死信隊列 | ✅ 內建 | ✅ | 需自建 |
| 延遲訊息 | ✅ 插件 | ✅ 內建 | 需自建 |
| 複雜度 | 中 | 低 | 極低 |
| 適用場景 | 企業級、高可靠 | 快速開發 | 單機、簡單任務 |

### 與 Laravel 對照

| Laravel | Go + RabbitMQ |
|---------|---------------|
| `QUEUE_CONNECTION=rabbitmq` | `amqp.Dial()` |
| `dispatch(new Job)` | `ch.PublishWithContext()` |
| `php artisan queue:work` | 程式內建 Consumer |
| `Job::onQueue('emails')` | 指定 Queue 名稱 |
| `$tries = 3` | 自行實現重試或用 `msg.Nack(requeue)` |
| 失敗任務表 | Dead Letter Queue |

---

## 與 Laravel 對照表

| Laravel | Go | 說明 |
|---------|-----|------|
| `dispatch(new Job)` | `pool.Submit(task)` | 發送任務 |
| `php artisan queue:work` | 內建於程式中 | Worker 執行 |
| `Job::dispatch()->delay(60)` | `time.AfterFunc()` 或 Asynq | 延遲任務 |
| Redis/Database Driver | Channel 或 Redis | 任務存儲 |
| `Queue::failing()` | 自定義錯誤處理 | 失敗處理 |
| `$tries = 3` | 自行實現重試邏輯 | 重試機制 |

---

## 注意事項

1. **避免 Goroutine 洩漏**：確保所有 goroutines 有退出條件
2. **Channel 關閉**：只由發送方關閉 channel，接收方不要關閉
3. **競態條件**：使用 `go run -race` 檢測
4. **資源限制**：控制 goroutine 數量，避免系統資源耗盡

```bash
# 競態檢測
go run -race main.go

# 或測試時
go test -race ./...
```

---

## 參考資源

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Asynq GitHub](https://github.com/hibiken/asynq)
- [Go by Example: Worker Pools](https://gobyexample.com/worker-pools)
- [RabbitMQ Go Client](https://github.com/rabbitmq/amqp091-go)
- [RabbitMQ Tutorials - Go](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
- [robfig/cron GitHub](https://github.com/robfig/cron)
