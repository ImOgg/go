package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"my-api/app"
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
func (ctrl *UserController) Store(c *gin.Context) {
	var req requests.CreateUserRequest
	
	// 驗證請求
	if err := req.Validate(c); err != nil {
		validationErrors := requests.FormatValidationError(err)
		traits.RespondValidationError(c, validationErrors)
		return
	}

	// 呼叫 Service 建立使用者
	user, err := ctrl.app.UserService.CreateUser(&req)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的使用者 ID", err.Error())
		return
	}

	// 呼叫 Service 刪除使用者
	if err := ctrl.app.UserService.DeleteUser(uint(id)); err != nil {
		traits.RespondNotFound(c, err.Error())
		return
	}

	traits.RespondSuccess(c, nil, "使用者刪除成功")
}
