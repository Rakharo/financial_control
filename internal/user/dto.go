package user

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
