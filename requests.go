package dyn

// NewLoginRequest JSON struct
type NewLoginRequest struct {
	CustomerName string `json:"customer_name"`
	Username     string `json:"user_name"`
	Password     string `json:"password"`
}
