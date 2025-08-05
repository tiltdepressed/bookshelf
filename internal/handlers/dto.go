package handlers

type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Role     string `json:"role" example:"user"`
}

type BookResponse struct {
	ID          uint    `json:"id" example:"1"`
	Title       string  `json:"title" example:"The Go Programming Language"`
	Author      string  `json:"author" example:"Alan A. A. Donovan"`
	Genre       string  `json:"genre" example:"Programming"`
	Description string  `json:"description" example:"Definitive guide to Go programming"`
	Price       float64 `json:"price" example:"49.99"`
}

type BookBriefResponse struct {
	ID     uint    `json:"id" example:"1"`
	Title  string  `json:"title" example:"The Go Programming Language"`
	Author string  `json:"author" example:"Alan A. A. Donovan"`
	Genre  string  `json:"genre" example:"Programming"`
	Price  float64 `json:"price" example:"49.99"`
}

type PaginatedBooksResponse struct {
	Data []BookBriefResponse `json:"data"`
	Meta struct {
		Total      int64 `json:"total" example:"100"`
		Page       int   `json:"page" example:"1"`
		Limit      int   `json:"limit" example:"10"`
		TotalPages int   `json:"totalPages" example:"10"`
	} `json:"meta"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"new_user"`
	Password string `json:"password" example:"strong_password"`
}

type LoginRequest struct {
	Username string `json:"username" example:"existing_user"`
	Password string `json:"password" example:"user_password"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type UpdateRoleRequest struct {
	NewRole string `json:"new_role" example:"admin"`
}
