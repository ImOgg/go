# JWT 認證系統

本專案使用 JWT (JSON Web Token) 實作無狀態的身份驗證機制。

---

## 架構概覽

```
請求流程：

┌─────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client    │───▶│  AuthMiddleware  │───▶│   Controller    │
│  (Bearer)   │    │  (驗證 Token)     │    │  (處理請求)      │
└─────────────┘    └──────────────────┘    └─────────────────┘
                            │
                            ▼
                   ┌──────────────────┐
                   │   utils/jwt.go   │
                   │  (Token 驗證)     │
                   └──────────────────┘
```

### 相關檔案

| 檔案 | 說明 |
|------|------|
| `app/utils/jwt.go` | JWT 工具函數（產生/驗證 Token、密碼加密） |
| `app/middleware/auth.go` | JWT 驗證中間件 |
| `app/controllers/auth_controller.go` | 認證控制器 |
| `app/services/auth_service.go` | 認證業務邏輯 |
| `app/requests/auth_request.go` | 請求驗證結構 |
| `app/responses/auth_response.go` | 回應 DTO |

---

## 環境配置

在 `.env` 檔案中設定 JWT 相關參數：

```env
# JWT 設定
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
```

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `JWT_SECRET` | Token 簽名密鑰（務必更換為強密碼） | `your-super-secret-key-change-in-production` |
| `JWT_EXPIRY_HOURS` | Token 有效期（小時） | `24` |

---

## API 端點

### 公開路由（不需要驗證）

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/api/register` | 使用者註冊 |
| POST | `/api/login` | 使用者登入 |

### 受保護路由（需要驗證）

| 方法 | 路徑 | 說明 |
|------|------|------|
| POST | `/api/logout` | 使用者登出 |
| GET | `/api/me` | 取得當前用戶資訊 |

---

## 使用方式

### 1. 使用者註冊

**請求：**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "張三",
    "email": "zhang@example.com",
    "password": "password123",
    "password_confirm": "password123",
    "age": 25
  }'
```

**請求欄位驗證：**

| 欄位 | 類型 | 必填 | 規則 |
|------|------|------|------|
| `name` | string | 是 | 2-100 字元 |
| `email` | string | 是 | 有效的電子郵件格式 |
| `password` | string | 是 | 8-72 字元 |
| `password_confirm` | string | 是 | 必須與 password 相同 |
| `age` | int | 否 | 0-150 |

**成功回應：**
```json
{
  "success": true,
  "message": "註冊成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "name": "張三",
      "email": "zhang@example.com",
      "age": 25,
      "created_at": "2026-01-29T10:00:00Z",
      "updated_at": "2026-01-29T10:00:00Z"
    }
  }
}
```

### 2. 使用者登入

**請求：**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhang@example.com",
    "password": "password123"
  }'
```

**請求欄位驗證：**

| 欄位 | 類型 | 必填 | 規則 |
|------|------|------|------|
| `email` | string | 是 | 有效的電子郵件格式 |
| `password` | string | 是 | 必填 |

**成功回應：**
```json
{
  "success": true,
  "message": "登入成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "name": "張三",
      "email": "zhang@example.com",
      "age": 25,
      "created_at": "2026-01-29T10:00:00Z",
      "updated_at": "2026-01-29T10:00:00Z"
    }
  }
}
```

### 3. 存取受保護的路由

取得 Token 後，在請求 Header 中加入 `Authorization`：

```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**成功回應：**
```json
{
  "success": true,
  "message": "取得用戶資訊成功",
  "data": {
    "id": 1,
    "name": "張三",
    "email": "zhang@example.com",
    "age": 25,
    "created_at": "2026-01-29T10:00:00Z",
    "updated_at": "2026-01-29T10:00:00Z"
  }
}
```

### 4. 使用者登出

```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**成功回應：**
```json
{
  "success": true,
  "message": "登出成功",
  "data": null
}
```

> **注意：** JWT 是無狀態的，登出只需客戶端刪除 Token。如需實作 Token 黑名單，可將已登出的 Token 存入 Redis。

---

## 錯誤回應

### 缺少授權憑證

```json
{
  "success": false,
  "message": "缺少授權憑證"
}
```

### 無效的授權格式

```json
{
  "success": false,
  "message": "無效的授權格式"
}
```

### Token 過期或無效

```json
{
  "success": false,
  "message": "無效或已過期的授權憑證"
}
```

### 登入失敗

```json
{
  "success": false,
  "message": "帳號或密碼錯誤"
}
```

### 註冊失敗（Email 已存在）

```json
{
  "success": false,
  "message": "電子郵件已被使用"
}
```

---

## JWT Token 結構

本專案的 JWT Token 包含以下聲明（Claims）：

```go
type JWTClaims struct {
    UserID uint   `json:"user_id"`  // 使用者 ID
    Email  string `json:"email"`    // 使用者 Email
    jwt.RegisteredClaims            // 標準聲明
}
```

**標準聲明包含：**

| 欄位 | 說明 |
|------|------|
| `exp` | 過期時間 |
| `iat` | 簽發時間 |
| `nbf` | 生效時間 |
| `iss` | 簽發者（本專案為 `my-api`） |

---

## 密碼安全

密碼使用 **bcrypt** 演算法加密：

```go
// 加密密碼
hashedPassword, err := utils.HashPassword(password)

// 驗證密碼
isValid := utils.CheckPassword(password, hashedPassword)
```

- 使用 `bcrypt.DefaultCost`（目前為 10）
- 密碼長度限制：8-72 字元（bcrypt 的限制）

---

## 在 Controller 中取得用戶資訊

AuthMiddleware 會將解析後的用戶資訊存入 Gin Context：

```go
func SomeHandler(c *gin.Context) {
    // 取得用戶 ID
    userID, exists := c.Get("user_id")
    if !exists {
        // 未授權
    }

    // 取得用戶 Email
    userEmail, _ := c.Get("user_email")

    // 使用 userID（需要轉型）
    id := userID.(uint)
}
```

---

## 與 Laravel 的對比

| 功能 | Laravel (Sanctum/Passport) | Go (本專案) |
|------|---------------------------|-------------|
| Token 產生 | `$user->createToken()` | `utils.GenerateToken()` |
| Token 驗證 | `auth:sanctum` 中間件 | `AuthMiddleware()` |
| 密碼加密 | `Hash::make()` | `utils.HashPassword()` |
| 密碼驗證 | `Hash::check()` | `utils.CheckPassword()` |
| 取得用戶 | `auth()->user()` | `c.Get("user_id")` |

---

## 安全建議

1. **務必更換 JWT_SECRET**：在生產環境使用強隨機密鑰
2. **使用 HTTPS**：避免 Token 在傳輸中被竊取
3. **設定合理的過期時間**：根據需求調整 `JWT_EXPIRY_HOURS`
4. **考慮實作 Token 黑名單**：用 Redis 存儲已登出的 Token
5. **考慮 Refresh Token**：實作 Token 刷新機制延長會話

---

## 相關文檔

- [專案設置](01-setup.md)
- [架構說明](02-architecture.md)
- [Controller 結構](05-controller-structure.md)
- [常見問題](07-troubleshooting.md)

---

**最後更新：** 2026-01-29
