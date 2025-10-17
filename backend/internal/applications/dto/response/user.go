package response

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type BaseUserResponse struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"`
	Data    UserResponse `json:"data"`
}
