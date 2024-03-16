package main

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	ExpiresIn *int `json:"expires_in_seconds"`
	userRequest
}

type chirpRequest struct {
	Body string `json:"body"`
}
