# 常見問題與解決方案

## 問題 1：main.go 缺少 package 聲明

**錯誤訊息：** 
```
main.go:1:1: expected 'package', found 'func'
```

**原因：** 
Go 的每個文件都必須以 `package` 開頭

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

---

## 問題 2：導入路徑錯誤

**錯誤訊息：** 
```
main.go:5:2: package go/routes is not in std
```

**原因：** 
導入路徑要用 `go.mod` 中定義的 module 名稱

**正確的導入方式：**
- go.mod 中是 `module my-api` → 導入時用 `"my-api/routes"`
- 不要用 `"go/routes"` 或其他自己造的路徑

**檢查步驟：**
1. 打開 `go.mod` 文件，查看第一行的 module 名稱
2. 確保所有導入的內部包都使用這個 module 名稱作為前綴
3. 例如：如果 module 是 `myapi`，導入 routes 包應該寫成 `"myapi/routes"`
