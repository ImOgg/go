# 初始化專案，名字可以隨便取（例如 myapi）
go mod init myapi

# 下載 Gin 框架
go get -u github.com/gin-gonic/gin

---

## 常見問題筆記

### 問題 1：main.go 缺少 package 聲明
**錯誤訊息：** `main.go:1:1: expected 'package', found 'func'`

**原因：** Go 的每個文件都必須以 `package` 開頭

**解決方案：**
```go
package main

import (
	"github.com/gin-gonic/gin"
	"my-api/routes"  // 注意：導入路徑要用 go.mod 中的 module 名稱
)

func main() {
	r := gin.Default()
	routes.InitRoutes(r)
	r.Run(":8080")
}
```

### 問題 2：導入路徑錯誤
**錯誤訊息：** `main.go:5:2: package go/routes is not in std`

**原因：** 導入路徑要用 `go.mod` 中定義的 module 名稱

**正確的導入方式：**
- go.mod 中是 `module my-api` → 導入時用 `"my-api/routes"`
- 不要用 `"go/routes"` 或其他自己造的路徑