package message

// POST /api/users/register
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

// POST /api/users/login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	ExpiresIn int64    `json:"expires_in"` // Token过期时间
	UserInfo  UserInfo `json:"user_info"`  // 用户信息
}
type UserInfo struct {
	ID       int32  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	CreateAt int64  `json:"create_at"`
}

// GET /api/users/refresh
type RefreshTokenResponse struct {
	ExpiresIn int64 `json:"expires_in"` // Token过期时间
}

// PUT /api/users/profile
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// PUT /api/users/password
type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// GET /api/users/{user_id}
type GetUserResponse struct {
	ID       int32  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	CreateAt int64  `json:"create_at"`
}
