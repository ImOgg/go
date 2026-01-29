package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/app/requests"
	"my-api/app/traits"
)

// AuthController - 認證控制器
type AuthController struct {
	app *app.App
}

// NewAuthController - 建立新的認證控制器
func NewAuthController(app *app.App) *AuthController {
	return &AuthController{app: app}
}

// Register - 使用者註冊
// POST /api/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req requests.RegisterRequest

	// 驗證請求
	if err := req.Validate(c); err != nil {
		validationErrors := requests.FormatValidationError(err)
		traits.RespondValidationError(c, validationErrors)
		return
	}

	// 呼叫 Service 註冊使用者
	response, err := ctrl.app.AuthService.Register(&req)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	traits.RespondCreated(c, response, "註冊成功")
}

// Login - 使用者登入
// POST /api/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest

	// 驗證請求
	if err := req.Validate(c); err != nil {
		validationErrors := requests.FormatValidationError(err)
		traits.RespondValidationError(c, validationErrors)
		return
	}

	// 呼叫 Service 登入
	response, err := ctrl.app.AuthService.Login(&req)
	if err != nil {
		traits.RespondUnauthorized(c, err.Error())
		return
	}

	traits.RespondSuccess(c, response, "登入成功")
}

// Logout - 使用者登出
// POST /api/logout
func (ctrl *AuthController) Logout(c *gin.Context) {
	// JWT 是無狀態的，登出只需客戶端刪除 Token
	// 如果需要 Token 黑名單，可以將 Token 存入 Redis
	traits.RespondSuccess(c, nil, "登出成功")
}

// Me - 取得當前用戶資訊
// GET /api/me
func (ctrl *AuthController) Me(c *gin.Context) {
	// 從 context 取得 user_id（由 AuthMiddleware 設定）
	userID, exists := c.Get("user_id")
	if !exists {
		traits.RespondUnauthorized(c, "未授權")
		return
	}

	// 呼叫 Service 取得用戶資訊
	response, err := ctrl.app.AuthService.GetCurrentUser(userID.(uint))
	if err != nil {
		traits.RespondNotFound(c, err.Error())
		return
	}

	traits.RespondSuccess(c, response, "取得用戶資訊成功")
}
