package dto

// userLoginDto
type UserLoginDto struct {
	Identifier string `json:"identifier"` // username or email
	Password   string `json:"password"`
}
