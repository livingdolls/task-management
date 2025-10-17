package response

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type BaseAuthResponse struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"`
	Data    AuthResponse `json:"data"`
}
