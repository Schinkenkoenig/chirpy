package main

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type chirpRequest struct {
	Body string `json:"body"`
}
