// package models

// type Record map[string]interface{}
package models

type Record struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"county"`
	Postal      string `json:"postal"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Web         string `json:"web"`
}
func (Record) TableName() string {
	return "records"
}