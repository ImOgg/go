package responses

import "my-api/app/models"

// AuthResponse - 認證成功回應（登入/註冊後返回）
type AuthResponse struct {
	Token     string        `json:"token"`
	TokenType string        `json:"token_type"`
	ExpiresIn int           `json:"expires_in"` // 秒數
	User      *UserResponse `json:"user"`
}

// NewAuthResponse - 建立認證回應
func NewAuthResponse(token string, expiresIn int, user *models.User) *AuthResponse {
	return &AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		User:      NewUserResponse(user),
	}
}
