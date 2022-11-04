package models

// LoginDTO defined the /login payload
type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"password"`
}

// SignupDTO defined the /login payload
type SignupDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"password"`
}
