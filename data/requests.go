package data

// AuthenticateRequest ... Data received via http during an authentication request.
type AuthenticationParams struct {
	Email    string
	Password string
}

// CreateUserRequest ... Data received via http during a create user request.
type CreateUserRequest struct {
	Auth  AuthenticationParams
	FName string
	LName string
}

// CreateUserResponse ... Data sent back via http from a succesful create user request.
type CreateUserResponse struct {
	Success bool
	Reason  string
	Token   string
}

// LoginRequest ... Data received via http during a login request.
type LoginRequest struct {
	Auth AuthenticationParams
}

// LoginResponse ... Data sent back via http from a successful login request.
type LoginResponse struct {
	Token string
}
