package services

import (
	"errors"
	"my-api/app/models"
	"my-api/app/requests"
	"testing"
)

// ============================================================================
// Mock Repository
// ============================================================================
// mockUserRepository 實作 UserRepository interface，用於測試
// 這樣測試時不會碰到真正的資料庫
type mockUserRepository struct {
	users    map[uint]*models.User // 模擬資料庫中的使用者
	emailMap map[string]*models.User // Email 索引，方便 FindByEmail
	nextID   uint                     // 自動遞增 ID
}

// newMockUserRepository 建立新的 mock repository
func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:    make(map[uint]*models.User),
		emailMap: make(map[string]*models.User),
		nextID:   1,
	}
}

// Create 模擬新增使用者
func (m *mockUserRepository) Create(user *models.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.emailMap[user.Email] = user
	return nil
}

// FindAll 模擬取得所有使用者
func (m *mockUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

// FindByID 模擬根據 ID 查詢
func (m *mockUserRepository) FindByID(id uint) (*models.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("record not found")
}

// Update 模擬更新使用者
func (m *mockUserRepository) Update(user *models.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("record not found")
	}
	// 更新 email 索引
	oldUser := m.users[user.ID]
	if oldUser.Email != user.Email {
		delete(m.emailMap, oldUser.Email)
		m.emailMap[user.Email] = user
	}
	m.users[user.ID] = user
	return nil
}

// Delete 模擬刪除使用者
func (m *mockUserRepository) Delete(id uint) error {
	if user, ok := m.users[id]; ok {
		delete(m.emailMap, user.Email)
		delete(m.users, id)
		return nil
	}
	return errors.New("record not found")
}

// FindByEmail 模擬根據 Email 查詢
func (m *mockUserRepository) FindByEmail(email string) (*models.User, error) {
	if user, ok := m.emailMap[email]; ok {
		return user, nil
	}
	return nil, errors.New("record not found")
}

// ============================================================================
// 測試案例
// ============================================================================

// TestUserService_CreateUser 測試新增使用者
func TestUserService_CreateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewUserService(mockRepo)

	// 測試案例
	tests := []struct {
		name    string
		req     *requests.CreateUserRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功新增使用者",
			req: &requests.CreateUserRequest{
				Name:  "張三",
				Email: "zhangsan@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "Email 已存在",
			req: &requests.CreateUserRequest{
				Name:  "李四",
				Email: "zhangsan@example.com", // 重複的 Email
				Age:   30,
			},
			wantErr: true,
			errMsg:  "電子郵件已被使用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreateUser(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("預期應該發生錯誤，但沒有")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("錯誤訊息不符，got %q, want %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("不預期的錯誤: %v", err)
				}
				if resp == nil {
					t.Error("回應不應為 nil")
				}
				if resp != nil && resp.Name != tt.req.Name {
					t.Errorf("Name 不符，got %q, want %q", resp.Name, tt.req.Name)
				}
			}
		})
	}
}

// TestUserService_GetUserByID 測試根據 ID 取得使用者
func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewUserService(mockRepo)

	// 先新增一個使用者
	mockRepo.Create(&models.User{
		Name:  "測試使用者",
		Email: "test@example.com",
		Age:   20,
	})

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "成功取得使用者",
			id:      1,
			wantErr: false,
		},
		{
			name:    "使用者不存在",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.GetUserByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("預期應該發生錯誤，但沒有")
				}
			} else {
				if err != nil {
					t.Errorf("不預期的錯誤: %v", err)
				}
				if resp == nil {
					t.Error("回應不應為 nil")
				}
			}
		})
	}
}

// TestUserService_UpdateUser 測試更新使用者
func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewUserService(mockRepo)

	// 先新增兩個使用者
	mockRepo.Create(&models.User{
		Name:  "使用者A",
		Email: "a@example.com",
		Age:   20,
	})
	mockRepo.Create(&models.User{
		Name:  "使用者B",
		Email: "b@example.com",
		Age:   25,
	})

	tests := []struct {
		name    string
		id      uint
		req     *requests.UpdateUserRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功更新名稱",
			id:   1,
			req: &requests.UpdateUserRequest{
				Name: "使用者A（已更新）",
			},
			wantErr: false,
		},
		{
			name: "使用者不存在",
			id:   999,
			req: &requests.UpdateUserRequest{
				Name: "不存在",
			},
			wantErr: true,
			errMsg:  "使用者不存在",
		},
		{
			name: "Email 已被其他人使用",
			id:   1,
			req: &requests.UpdateUserRequest{
				Email: "b@example.com", // 使用者B 的 Email
			},
			wantErr: true,
			errMsg:  "電子郵件已被使用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.UpdateUser(tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("預期應該發生錯誤，但沒有")
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("錯誤訊息不符，got %q, want %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("不預期的錯誤: %v", err)
				}
				if resp == nil {
					t.Error("回應不應為 nil")
				}
			}
		})
	}
}

// TestUserService_DeleteUser 測試刪除使用者
func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewUserService(mockRepo)

	// 先新增一個使用者
	mockRepo.Create(&models.User{
		Name:  "待刪除使用者",
		Email: "delete@example.com",
		Age:   30,
	})

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "成功刪除使用者",
			id:      1,
			wantErr: false,
		},
		{
			name:    "使用者不存在",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteUser(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("預期應該發生錯誤，但沒有")
				}
			} else {
				if err != nil {
					t.Errorf("不預期的錯誤: %v", err)
				}
			}
		})
	}
}

// TestUserService_GetAllUsers 測試取得所有使用者
func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewUserService(mockRepo)

	// 測試空列表
	t.Run("空列表", func(t *testing.T) {
		users, err := service.GetAllUsers()
		if err != nil {
			t.Errorf("不預期的錯誤: %v", err)
		}
		if len(users) != 0 {
			t.Errorf("預期空列表，got %d 筆", len(users))
		}
	})

	// 新增一些使用者
	mockRepo.Create(&models.User{Name: "User1", Email: "user1@example.com", Age: 20})
	mockRepo.Create(&models.User{Name: "User2", Email: "user2@example.com", Age: 25})
	mockRepo.Create(&models.User{Name: "User3", Email: "user3@example.com", Age: 30})

	// 測試有資料
	t.Run("有資料", func(t *testing.T) {
		users, err := service.GetAllUsers()
		if err != nil {
			t.Errorf("不預期的錯誤: %v", err)
		}
		if len(users) != 3 {
			t.Errorf("預期 3 筆，got %d 筆", len(users))
		}
	})
}
