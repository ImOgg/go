package utils

import (
	"testing"
)

// TestHashPassword - 測試密碼加密功能
func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"

	// 測試加密是否成功
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() 發生錯誤: %v", err)
	}

	// 確認 hash 不為空
	if hash == "" {
		t.Error("HashPassword() 回傳空字串")
	}

	// 確認 hash 與原密碼不同
	if hash == password {
		t.Error("HashPassword() 回傳的 hash 與原密碼相同，加密失敗")
	}

	t.Logf("原密碼: %s", password)
	t.Logf("加密後: %s", hash)
}

// TestCheckPassword - 測試密碼驗證功能
func TestCheckPassword(t *testing.T) {
	password := "mySecretPassword123"
	wrongPassword := "wrongPassword456"

	// 先加密密碼
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() 發生錯誤: %v", err)
	}

	// 測試正確密碼
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword() 正確密碼驗證失敗")
	}

	// 測試錯誤密碼
	if CheckPassword(wrongPassword, hash) {
		t.Error("CheckPassword() 錯誤密碼驗證應該失敗，但卻成功了")
	}
}

// TestHashPassword_DifferentHashes - 測試同樣的密碼會產生不同的 hash（因為 salt）
func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// bcrypt 每次加密都會產生不同的 hash（因為隨機 salt）
	if hash1 == hash2 {
		t.Error("同樣的密碼應該產生不同的 hash（因為 salt），但產生了相同的 hash")
	}

	// 但兩個 hash 都應該能驗證成功
	if !CheckPassword(password, hash1) {
		t.Error("hash1 驗證失敗")
	}
	if !CheckPassword(password, hash2) {
		t.Error("hash2 驗證失敗")
	}
}

// TestCheckPassword_TableDriven - 使用 Table-Driven 測試（Go 推薦的測試模式）
func TestCheckPassword_TableDriven(t *testing.T) {
	// 準備一個已知的 hash
	originalPassword := "testPassword"
	hash, _ := HashPassword(originalPassword)

	// 定義測試案例
	tests := []struct {
		name     string // 測試名稱
		password string // 輸入的密碼
		want     bool   // 預期結果
	}{
		{
			name:     "正確密碼",
			password: originalPassword,
			want:     true,
		},
		{
			name:     "錯誤密碼",
			password: "wrongPassword",
			want:     false,
		},
		{
			name:     "空密碼",
			password: "",
			want:     false,
		},
		{
			name:     "相似但不同的密碼",
			password: "testPassword1",
			want:     false,
		},
		{
			name:     "大小寫不同",
			password: "TestPassword",
			want:     false,
		},
	}

	// 執行每個測試案例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.password, hash)
			if got != tt.want {
				t.Errorf("CheckPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

// TestHashPassword_EmptyPassword - 測試空密碼
func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() 處理空密碼時發生錯誤: %v", err)
	}

	// 空密碼也應該能正確加密和驗證
	if !CheckPassword("", hash) {
		t.Error("空密碼驗證失敗")
	}
}
