package auth

import "time"

type Role string

const (
    RolePlatformAdmin Role = "platform_admin"
    RoleCompanyAdmin  Role = "company_admin"
    RoleSafetyManager Role = "safety_manager"
    RoleTrainer       Role = "trainer"
    RoleTrainee       Role = "trainee"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    AccessToken string    `json:"access_token"`
    ExpiresAt   time.Time `json:"expires_at"`
    UserID      string    `json:"user_id"`
    CompanyID   string    `json:"company_id"`
    Role        Role      `json:"role"`
}

type Principal struct {
    UserID    string
    CompanyID string
    Role      Role
}
