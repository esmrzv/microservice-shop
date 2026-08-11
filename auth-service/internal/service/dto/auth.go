package dto



type RegisterRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


type RegisterResponse struct {
	Email string `json:"email"`
}


type LoginRequest struct {
	Email string 	`json:"email"`
	Password string `json:"password"`
}


type LoginResponse struct {
	Token string 	`json:"token"`
}