package entity

type User struct {
	Username string `json:"username" gorm:"unique"`
	Password string
}

type MsgResponse struct {
	Message string `json:"message" example:"user123 login successfully"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Invalid Username and Password"`
}
