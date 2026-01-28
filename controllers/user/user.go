package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"my-api/database"
	"my-api/models"
)

// GetUsersGORM - GORM 範例：取得所有使用者
func GetUsersGORM(c *gin.Context) {
	var users []models.User
	result := database.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"method": "GORM",
		"count":  len(users),
		"data":   users,
	})
}

// CreateUserGORM - GORM 範例：新增使用者
func CreateUserGORM(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"method": "GORM",
		"data":   user,
	})
}

// GetUsersSQL - 原生 SQL 範例：取得所有使用者
func GetUsersSQL(c *gin.Context) {
	rows, err := database.SqlDB.Query("SELECT id, name, email, age FROM users WHERE deleted_at IS NULL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email string
		var age int

		if err := rows.Scan(&id, &name, &email, &age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
			"age":   age,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"method": "Raw SQL",
		"count":  len(users),
		"data":   users,
	})
}

// CreateUserSQL - 原生 SQL 範例：新增使用者
func CreateUserSQL(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
		Age   int    `json:"age"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.SqlDB.Exec(
		"INSERT INTO users (name, email, age, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())",
		input.Name, input.Email, input.Age,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"method": "Raw SQL",
		"data": map[string]interface{}{
			"id":    id,
			"name":  input.Name,
			"email": input.Email,
			"age":   input.Age,
		},
	})
}

// GetByName - 根據名稱取得使用者
func GetByName(c *gin.Context) {
	// 抓取網址路徑裡的 :name
	name := c.Param("name")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "歡迎回來！",
		"data":    name,
	})
}
