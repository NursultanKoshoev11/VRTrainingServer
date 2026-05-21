package companies

import "time"

type CompanyStatus string

const (
    CompanyStatusActive   CompanyStatus = "active"
    CompanyStatusDisabled CompanyStatus = "disabled"
)

type Company struct {
    ID        string        `json:"id"`
    Name      string        `json:"name"`
    Industry  string        `json:"industry"`
    Country   string        `json:"country"`
    Status    CompanyStatus `json:"status"`
    CreatedAt time.Time     `json:"created_at"`
    UpdatedAt time.Time     `json:"updated_at"`
}
