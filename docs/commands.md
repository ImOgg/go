# Go 常用命令

## 編譯與執行

### 編譯程式碼
編譯 test.go 程式碼，產生執行檔：
```bash
go build test.go
```

編譯整個專案（在專案根目錄執行）：
```bash
go build
```

### 清除編譯結果
```bash
go clean
```

## 依賴管理

### 下載依賴
```bash
go get -u github.com/gin-gonic/gin
```

### 整理依賴
```bash
go mod tidy
```

### 查看依賴
```bash
go mod graph
```

## 運行程式

### 直接運行（不編譯）
```bash
go run main.go
```

### 運行編譯後的執行檔
```bash
# Windows
.\myapi.exe

# Linux/Mac
./myapi
```

## 測試

### 運行測試
```bash
go test ./...
```

### 運行測試並顯示詳細信息
```bash
go test -v ./...
```
