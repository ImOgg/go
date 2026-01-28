# 專案設置

## 初始化專案

```bash
# 初始化專案，名字可以隨便取（例如 myapi）
go mod init myapi

# 下載 Gin 框架
go get -u github.com/gin-gonic/gin
```

## 專案結構

```
.
├── main.go           # 主程式入口
├── go.mod            # 模組依賴管理
├── Dockerfile        # Docker 配置
├── docker-compose.yml
├── controllers/      # 控制器
│   └── hello.go
├── routes/           # 路由設定
│   └── router.go
└── docs/             # 文檔
    ├── setup.md
    ├── troubleshooting.md
    └── commands.md
```

## 參考資源

- [使用 Golang 打造 Web 應用程式](https://willh.gitbook.io/build-web-application-with-golang-zhtw/01.0)
- [Golang教學筆記](https://hackmd.io/@action/rk8R2cuAU)
