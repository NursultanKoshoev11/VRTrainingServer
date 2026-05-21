package users

import "time"

type User struct {
    ID         string    `json:"id"`
    CompanyID  string    `json:"company_id"`
    Role       string    `json:"role"`
    FullName   string    `json:"full_name"`
    Email      string    `json:"email"`
    EmployeeID string    `json:"employee_id"`
    Language   string    `json:"language"`
    IsActive   bool      `json:"is_active"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
