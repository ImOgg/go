package controllers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/app/models"
	"my-api/app/requests"
	"my-api/app/traits"
)

type PostController struct {
	app *app.App
}

func NewPostController(app *app.App) *PostController {
	return &PostController{app: app}
}

// Index - 取得所有文章
func (ctrl *PostController) Index(c *gin.Context) {
	posts, err := ctrl.app.PostRepository.FindAll()
	if err != nil {
		traits.RespondError(c, http.StatusInternalServerError, "取得文章列表失敗", err.Error())
		return
	}

	traits.RespondSuccess(c, posts, "成功取得文章列表")
}

// Show - 取得單一文章
func (ctrl *PostController) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的文章 ID", err.Error())
		return
	}

	post, err := ctrl.app.PostRepository.FindByID(uint(id))
	if err != nil {
		traits.RespondError(c, http.StatusNotFound, "文章不存在", err.Error())
		return
	}

	traits.RespondSuccess(c, post, "成功取得文章")
}

// Store - 建立新文章
func (ctrl *PostController) Store(c *gin.Context) {
	var req requests.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的請求格式", err.Error())
		return
	}

	post := &models.Post{
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		UserID:      req.UserID,
	}

	if err := ctrl.app.PostRepository.Create(post); err != nil {
		traits.RespondError(c, http.StatusInternalServerError, "建立文章失敗", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "成功建立文章",
		"data":    post,
	})
}

// Update - 更新文章
func (ctrl *PostController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的文章 ID", err.Error())
		return
	}

	var req requests.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的請求格式", err.Error())
		return
	}

	post, err := ctrl.app.PostRepository.FindByID(uint(id))
	if err != nil {
		traits.RespondError(c, http.StatusNotFound, "文章不存在", err.Error())
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Description != "" {
		post.Description = req.Description
	}

	if err := ctrl.app.PostRepository.Update(post); err != nil {
		traits.RespondError(c, http.StatusInternalServerError, "更新文章失敗", err.Error())
		return
	}

	traits.RespondSuccess(c, post, "成功更新文章")
}

// Delete - 刪除文章
func (ctrl *PostController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		traits.RespondError(c, http.StatusBadRequest, "無效的文章 ID", err.Error())
		return
	}

	if err := ctrl.app.PostRepository.Delete(uint(id)); err != nil {
		traits.RespondError(c, http.StatusInternalServerError, "刪除文章失敗", err.Error())
		return
	}

	traits.RespondSuccess(c, nil, "成功刪除文章")
}
