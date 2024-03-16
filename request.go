package main

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	userRequest
}

type chirpRequest struct {
	Body string `json:"body"`
}
