package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/app/pkg/logger"
	"my-api/app/requests"
	"my-api/app/traits"
)

// UserController - 使用者控制器
type UserController struct {
	app *app.App
}

// NewUserController - 建立新的使用者控制器
func NewUserController(app *app.App) *UserController {
	return &UserController{app: app}
}

// Index - 取得所有使用者
// GET /api/users
//
// Go 方法語法說明：
// func (ctrl *UserController) Index(c *gin.Context)
//      ─────────┬───────────  ──┬──  ──────┬──────
//            接收者           方法名      參數
//      （等於 Laravel 的 $this）      （Gin 的請求上下文）
//
// ctrl.app.UserService → 等於 Laravel 的 $this->userService
func (ctrl *UserController) Index(c *gin.Context) {
	users, err := ctrl.app.UserService.GetAllUsers()
	if err != nil {
		traits.RespondError(c, http.StatusInternalServerError, "取得使用者列表失敗", err.Error())
		return
	}

	traits.RespondSuccess(c, users, "成功取得使用者列表")
}

// Show - 取得單一使用者
// GET /api/users/:id
func (ctrl *UserController) Show(c *gin.Context) {
	// strconv.ParseUint(字串, 進位制, 位元大小)
	// 10 = 十進位，32 = 限制在 uint32 範圍（0 ~ 4,294,967,295）
	// 防止數值溢位（overflow），避免惡意輸入超大數字造成系統問題
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的使用者 ID", err.Error())
		return
	}

	user, err := ctrl.app.UserService.GetUserByID(uint(id))
	if err != nil {
		traits.RespondNotFound(c, err.Error())
		return
	}

	traits.RespondSuccess(c, user, "成功取得使用者資料")
}

// Store - 新增使用者
// POST /api/users
//
// 資料綁定流程：
// 1. var req → 宣告空的 struct（此時 Name="", Email="", Age=0）
// 2. req.Validate(c) → 內部呼叫 ShouldBindJSON，把 JSON body 填入 req
// 3. CreateUser(&req) → 用填好資料的 req 建立使用者
//
// 類似 Laravel：public function store(CreateUserRequest $request)
// 差別是 Laravel 自動注入並驗證，Go 需要手動宣告和呼叫
func (ctrl *UserController) Store(c *gin.Context) {
	// 從 context 取得 Logger（已帶有 request_id）
	log := logger.FromGinContext(c)

	// 步驟 1：宣告空的 request struct
	var req requests.CreateUserRequest

	// 步驟 2：驗證請求（內部會把 JSON 資料綁定到 req）
	if err := req.Validate(c); err != nil {
		log.Warning("使用者建立驗證失敗", map[string]interface{}{
			"errors": err.Error(),
		})
		validationErrors := requests.FormatValidationError(err)
		traits.RespondValidationError(c, validationErrors)
		return
	}

	log.Debug("開始建立使用者", map[string]interface{}{
		"email": req.Email,
		"name":  req.Name,
	})

	// 步驟 3：用填好資料的 req 建立使用者
	// &req = 傳遞 req 的記憶體位址（指標），避免複製整個 struct
	user, err := ctrl.app.UserService.CreateUser(&req)
	if err != nil {
		log.Error("使用者建立失敗", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		traits.RespondError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	log.Info("使用者建立成功", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	traits.RespondCreated(c, user, "使用者建立成功")
}

// Update - 更新使用者
// PUT /api/users/:id
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的使用者 ID", err.Error())
		return
	}

	var req requests.UpdateUserRequest
	
	// 驗證請求
	if err := req.Validate(c); err != nil {
		validationErrors := requests.FormatValidationError(err)
		traits.RespondValidationError(c, validationErrors)
		return
	}

	// 呼叫 Service 更新使用者
	user, err := ctrl.app.UserService.UpdateUser(uint(id), &req)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	traits.RespondSuccess(c, user, "使用者更新成功")
}

// Destroy - 刪除使用者
// DELETE /api/users/:id
func (ctrl *UserController) Destroy(c *gin.Context) {
	log := logger.FromGinContext(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Warning("刪除使用者失敗：無效的 ID", map[string]interface{}{
			"input": c.Param("id"),
			"error": err.Error(),
		})
		traits.RespondError(c, http.StatusBadRequest, "無效的使用者 ID", err.Error())
		return
	}

	// 呼叫 Service 刪除使用者
	if err := ctrl.app.UserService.DeleteUser(uint(id)); err != nil {
		log.Warning("刪除使用者失敗：找不到使用者", map[string]interface{}{
			"user_id": id,
			"error":   err.Error(),
		})
		traits.RespondNotFound(c, err.Error())
		return
	}

	log.Info("使用者刪除成功", map[string]interface{}{
		"user_id": id,
	})

	traits.RespondSuccess(c, nil, "使用者刪除成功")
}
